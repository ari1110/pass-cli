package models

import (
	"fmt"
	"sort"
	"sync"

	"pass-cli/internal/vault"

	"github.com/rivo/tview"
)

// VaultService interface defines the vault operations needed by AppState.
// This interface enables testing with mock implementations.
type VaultService interface {
	ListCredentialsWithMetadata() ([]vault.CredentialMetadata, error)
	AddCredential(service, username, password, category, url, notes string) error
	UpdateCredential(service string, opts vault.UpdateOpts) error
	DeleteCredential(service string) error
	GetCredential(service string, trackUsage bool) (*vault.Credential, error)
}

// UpdateCredentialOpts mirrors vault.UpdateOpts for AppState layer.
// Using pointer fields allows distinguishing "don't update" (nil) from "clear to empty" (non-nil empty string).
type UpdateCredentialOpts struct {
	Username *string
	Password *string
	Category *string
	URL      *string
	Notes    *string
}

// AppState holds all application state with thread-safe access.
// This is the single source of truth for the entire TUI.
type AppState struct {
	// Concurrency control
	mu sync.RWMutex // Protects all fields below

	// Vault service (interface for testability)
	vault VaultService

	// Credential data
	credentials []vault.CredentialMetadata
	categories  []string

	// Current selections
	selectedCategory   string
	selectedCredential *vault.CredentialMetadata

	// UI components (single instances, created once)
	sidebar    *tview.TreeView
	table      *tview.Table
	detailView *tview.TextView
	statusBar  *tview.TextView

	// Notification callbacks
	onCredentialsChanged func()      // Called when credentials are loaded/modified
	onSelectionChanged   func()      // Called when selection changes
	onError              func(error) // Called when errors occur
}

// NewAppState creates a new AppState with the given vault service.
func NewAppState(vaultService VaultService) *AppState {
	return &AppState{
		vault:       vaultService,
		credentials: make([]vault.CredentialMetadata, 0),
		categories:  make([]string, 0),
	}
}

// GetCredentials returns a copy of the credentials slice (thread-safe read).
func (s *AppState) GetCredentials() []vault.CredentialMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.credentials
}

// GetCategories returns a copy of the categories slice (thread-safe read).
func (s *AppState) GetCategories() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external mutation
	categories := make([]string, len(s.categories))
	copy(categories, s.categories)
	return categories
}

// GetSelectedCredential returns a copy of the selected credential (thread-safe read).
func (s *AppState) GetSelectedCredential() *vault.CredentialMetadata {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.selectedCredential
}

// GetSelectedCategory returns the selected category (thread-safe read).
func (s *AppState) GetSelectedCategory() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.selectedCategory
}

// GetFullCredential fetches the full credential including password from vault.
// This is used when password access is needed (display, clipboard).
// SECURITY: Only call when password is actually needed (on-demand fetching).
func (s *AppState) GetFullCredential(service string) (*vault.Credential, error) {
	return s.GetFullCredentialWithTracking(service, true)
}

// GetFullCredentialWithTracking retrieves a credential with optional usage tracking.
// Set track=false to avoid incrementing usage statistics (e.g., for form pre-population).
func (s *AppState) GetFullCredentialWithTracking(service string, track bool) (*vault.Credential, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.vault.GetCredential(service, track)
}

// LoadCredentials loads all credentials from the vault.
// CRITICAL: Follows Lock→Mutate→Unlock→Notify pattern to prevent deadlocks.
func (s *AppState) LoadCredentials() error {
	s.mu.Lock()

	// Load credentials from vault
	creds, err := s.vault.ListCredentialsWithMetadata()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to load credentials: %w", err)
		s.mu.Unlock()             // ✅ RELEASE LOCK FIRST
		s.notifyError(wrappedErr) // ✅ THEN notify
		return wrappedErr
	}

	// Update state
	s.credentials = creds
	s.updateCategories() // Internal helper, safe to call while locked

	s.mu.Unlock()                // ✅ RELEASE LOCK
	s.notifyCredentialsChanged() // ✅ THEN notify

	return nil
}

// AddCredential adds a new credential to the vault.
// CRITICAL: Minimizes lock duration by releasing lock during vault I/O operations.
func (s *AppState) AddCredential(service, username, password, category, url, notes string) error {
	// Perform vault I/O without holding lock (vault has its own synchronization)
	err := s.vault.AddCredential(service, username, password, category, url, notes)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to add credential: %w", err)
		s.notifyError(wrappedErr)
		return wrappedErr
	}

	// Reload credentials without holding lock
	creds, err := s.vault.ListCredentialsWithMetadata()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to reload credentials: %w", err)
		s.notifyError(wrappedErr)
		return wrappedErr
	}

	// Only lock to update state
	s.mu.Lock()
	s.credentials = creds
	s.updateCategories() // Update categories while locked
	s.mu.Unlock()

	// Notify after releasing lock
	s.notifyCredentialsChanged()

	return nil
}

// UpdateCredential updates an existing credential in the vault.
// CRITICAL: Minimizes lock duration by releasing lock during vault I/O operations.
// Accepts UpdateCredentialOpts to allow clearing fields to empty strings (non-nil pointer to empty string).
func (s *AppState) UpdateCredential(service string, opts UpdateCredentialOpts) error {
	// Convert AppState UpdateCredentialOpts to vault.UpdateOpts
	vaultOpts := vault.UpdateOpts{
		Username: opts.Username,
		Password: opts.Password,
		Category: opts.Category,
		URL:      opts.URL,
		Notes:    opts.Notes,
	}

	// Perform vault I/O without holding lock (vault has its own synchronization)
	err := s.vault.UpdateCredential(service, vaultOpts)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to update credential: %w", err)
		s.notifyError(wrappedErr)
		return wrappedErr
	}

	// Reload credentials without holding lock
	creds, err := s.vault.ListCredentialsWithMetadata()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to reload credentials: %w", err)
		s.notifyError(wrappedErr)
		return wrappedErr
	}

	// Only lock to update state
	s.mu.Lock()
	s.credentials = creds
	s.updateCategories() // Update categories while locked
	s.mu.Unlock()

	// Notify after releasing lock
	s.notifyCredentialsChanged()

	return nil
}

// DeleteCredential deletes a credential from the vault.
// CRITICAL: Minimizes lock duration by releasing lock during vault I/O operations.
func (s *AppState) DeleteCredential(service string) error {
	// Perform vault I/O without holding lock (vault has its own synchronization)
	err := s.vault.DeleteCredential(service)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to delete credential: %w", err)
		s.notifyError(wrappedErr)
		return wrappedErr
	}

	// Reload credentials without holding lock
	creds, err := s.vault.ListCredentialsWithMetadata()
	if err != nil {
		wrappedErr := fmt.Errorf("failed to reload credentials: %w", err)
		s.notifyError(wrappedErr)
		return wrappedErr
	}

	// Only lock to update state
	s.mu.Lock()
	s.credentials = creds
	s.updateCategories() // Update categories while locked
	s.mu.Unlock()

	// Notify after releasing lock
	s.notifyCredentialsChanged()

	return nil
}

// SetSelectedCategory updates the selected category.
// CRITICAL: Follows Lock→Mutate→Unlock→Notify pattern.
func (s *AppState) SetSelectedCategory(category string) {
	s.mu.Lock()
	s.selectedCategory = category
	s.mu.Unlock() // ✅ RELEASE LOCK

	s.notifySelectionChanged() // ✅ THEN notify
}

// SetSelectedCredential updates the selected credential.
// CRITICAL: Follows Lock→Mutate→Unlock→Notify pattern.
func (s *AppState) SetSelectedCredential(credential *vault.CredentialMetadata) {
	s.mu.Lock()
	s.selectedCredential = credential
	s.mu.Unlock() // ✅ RELEASE LOCK

	s.notifySelectionChanged() // ✅ THEN notify
}

// SetSidebar stores the sidebar component reference.
func (s *AppState) SetSidebar(sidebar *tview.TreeView) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sidebar = sidebar
}

// GetSidebar retrieves the sidebar component reference.
func (s *AppState) GetSidebar() *tview.TreeView {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sidebar
}

// SetTable stores the table component reference.
func (s *AppState) SetTable(table *tview.Table) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.table = table
}

// GetTable retrieves the table component reference.
func (s *AppState) GetTable() *tview.Table {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.table
}

// SetDetailView stores the detail view component reference.
func (s *AppState) SetDetailView(view *tview.TextView) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.detailView = view
}

// GetDetailView retrieves the detail view component reference.
func (s *AppState) GetDetailView() *tview.TextView {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.detailView
}

// SetStatusBar stores the status bar component reference.
func (s *AppState) SetStatusBar(bar *tview.TextView) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statusBar = bar
}

// GetStatusBar retrieves the status bar component reference.
func (s *AppState) GetStatusBar() *tview.TextView {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.statusBar
}

// SetOnCredentialsChanged registers a callback for credential changes.
func (s *AppState) SetOnCredentialsChanged(callback func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onCredentialsChanged = callback
}

// SetOnSelectionChanged registers a callback for selection changes.
func (s *AppState) SetOnSelectionChanged(callback func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onSelectionChanged = callback
}

// SetOnError registers a callback for errors.
func (s *AppState) SetOnError(callback func(error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onError = callback
}

// notifyCredentialsChanged invokes the credentials changed callback.
// CRITICAL: Must be called AFTER releasing locks to prevent deadlocks.
func (s *AppState) notifyCredentialsChanged() {
	// Read callback without holding lock
	s.mu.RLock()
	callback := s.onCredentialsChanged
	s.mu.RUnlock()

	if callback != nil {
		callback()
	}
}

// notifySelectionChanged invokes the selection changed callback.
// CRITICAL: Must be called AFTER releasing locks to prevent deadlocks.
func (s *AppState) notifySelectionChanged() {
	// Read callback without holding lock
	s.mu.RLock()
	callback := s.onSelectionChanged
	s.mu.RUnlock()

	if callback != nil {
		callback()
	}
}

// notifyError invokes the error callback.
// CRITICAL: Must be called AFTER releasing locks to prevent deadlocks.
func (s *AppState) notifyError(err error) {
	// Read callback without holding lock
	s.mu.RLock()
	callback := s.onError
	s.mu.RUnlock()

	if callback != nil {
		callback(err)
	}
}

// updateCategories extracts unique categories from credentials.
// CRITICAL: Must be called while holding a write lock.
func (s *AppState) updateCategories() {
	categoryMap := make(map[string]bool)

	for _, cred := range s.credentials {
		// Extract category from credential's Category field
		if cred.Category != "" {
			categoryMap[cred.Category] = true
		} else {
			// Empty category becomes "Uncategorized"
			categoryMap["Uncategorized"] = true
		}
	}

	// Convert map to sorted slice
	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	s.categories = categories
}
