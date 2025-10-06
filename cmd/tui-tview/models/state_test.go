package models

import (
	"errors"
	"sync"
	"testing"
	"time"

	"pass-cli/internal/vault"

	"github.com/rivo/tview"
)

// MockVaultService is a mock implementation of VaultService for testing.
type MockVaultService struct {
	mu sync.Mutex

	// Mock data
	credentials []vault.CredentialMetadata

	// Mock behaviors
	listError   error
	addError    error
	updateError error
	deleteError error
	getError    error

	// Call tracking
	listCalled   int
	addCalled    int
	updateCalled int
	deleteCalled int
	getCalled    int
}

// NewMockVaultService creates a new mock vault service.
func NewMockVaultService() *MockVaultService {
	return &MockVaultService{
		credentials: make([]vault.CredentialMetadata, 0),
	}
}

// ListCredentialsWithMetadata returns the mock credentials.
func (m *MockVaultService) ListCredentialsWithMetadata() ([]vault.CredentialMetadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.listCalled++
	if m.listError != nil {
		return nil, m.listError
	}
	return m.credentials, nil
}

// AddCredential adds a mock credential.
func (m *MockVaultService) AddCredential(service, username, password, category string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.addCalled++
	if m.addError != nil {
		return m.addError
	}

	// Add credential to mock storage
	cred := vault.CredentialMetadata{
		Service:      service,
		Username:     username,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastAccessed: time.Time{},
	}
	m.credentials = append(m.credentials, cred)
	return nil
}

// UpdateCredential updates a mock credential.
func (m *MockVaultService) UpdateCredential(service, username, password, category string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.updateCalled++
	if m.updateError != nil {
		return m.updateError
	}

	// Find and update credential
	for i, cred := range m.credentials {
		if cred.Service == service {
			m.credentials[i].Username = username
			m.credentials[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return errors.New("credential not found")
}

// DeleteCredential deletes a mock credential.
func (m *MockVaultService) DeleteCredential(service string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.deleteCalled++
	if m.deleteError != nil {
		return m.deleteError
	}

	// Find and remove credential
	for i, cred := range m.credentials {
		if cred.Service == service {
			m.credentials = append(m.credentials[:i], m.credentials[i+1:]...)
			return nil
		}
	}
	return errors.New("credential not found")
}

// GetCredential returns a mock full credential.
func (m *MockVaultService) GetCredential(service string, trackUsage bool) (*vault.Credential, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.getCalled++
	if m.getError != nil {
		return nil, m.getError
	}

	// Find credential
	for _, cred := range m.credentials {
		if cred.Service == service {
			return &vault.Credential{
				Service:     cred.Service,
				Username:    cred.Username,
				Password:    "mock-password",
				Notes:       cred.Notes,
				CreatedAt:   cred.CreatedAt,
				UpdatedAt:   cred.UpdatedAt,
				UsageRecord: make(map[string]vault.UsageRecord),
			}, nil
		}
	}
	return nil, errors.New("credential not found")
}

// SetCredentials sets the mock credentials for testing.
func (m *MockVaultService) SetCredentials(creds []vault.CredentialMetadata) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.credentials = creds
}

// TestNewAppState verifies AppState creation.
func TestNewAppState(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	if state == nil {
		t.Fatal("NewAppState returned nil")
	}

	// Verify initial state
	if len(state.GetCredentials()) != 0 {
		t.Errorf("Expected empty credentials, got %d", len(state.GetCredentials()))
	}
	if len(state.GetCategories()) != 0 {
		t.Errorf("Expected empty categories, got %d", len(state.GetCategories()))
	}
	if state.GetSelectedCategory() != "" {
		t.Errorf("Expected empty selected category, got %s", state.GetSelectedCategory())
	}
	if state.GetSelectedCredential() != nil {
		t.Error("Expected nil selected credential")
	}
}

// TestLoadCredentials verifies credential loading from vault.
func TestLoadCredentials(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Setup mock data
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
		{Service: "GitHub", Username: "user", CreatedAt: time.Now()},
		{Service: "Database", Username: "dbuser", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)

	// Load credentials
	err := state.LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials failed: %v", err)
	}

	// Verify credentials loaded
	creds := state.GetCredentials()
	if len(creds) != 3 {
		t.Errorf("Expected 3 credentials, got %d", len(creds))
	}

	// Verify categories extracted
	categories := state.GetCategories()
	if len(categories) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(categories))
	}
}

// TestLoadCredentials_Error verifies error handling in LoadCredentials.
func TestLoadCredentials_Error(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Setup error
	expectedErr := errors.New("vault error")
	mockVault.listError = expectedErr

	// Setup error callback
	var callbackErr error
	state.SetOnError(func(err error) {
		callbackErr = err
	})

	// Load credentials (should fail)
	err := state.LoadCredentials()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Verify error callback invoked
	if callbackErr == nil {
		t.Error("Error callback was not invoked")
	}
}

// TestAddCredential verifies adding a credential.
func TestAddCredential(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Track callback invocation
	callbackInvoked := false
	state.SetOnCredentialsChanged(func() {
		callbackInvoked = true
	})

	// Add credential
	err := state.AddCredential("AWS", "admin", "password123")
	if err != nil {
		t.Fatalf("AddCredential failed: %v", err)
	}

	// Verify callback invoked
	if !callbackInvoked {
		t.Error("onCredentialsChanged callback was not invoked")
	}

	// Verify credential added to mock vault
	if mockVault.addCalled != 1 {
		t.Errorf("Expected AddCredential called 1 time, got %d", mockVault.addCalled)
	}

	// Verify state updated
	creds := state.GetCredentials()
	if len(creds) != 1 {
		t.Errorf("Expected 1 credential, got %d", len(creds))
	}
	if creds[0].Service != "AWS" {
		t.Errorf("Expected service 'AWS', got '%s'", creds[0].Service)
	}
}

// TestUpdateCredential verifies updating a credential.
func TestUpdateCredential(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Setup existing credential
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()

	// Track callback invocation
	callbackInvoked := false
	state.SetOnCredentialsChanged(func() {
		callbackInvoked = true
	})

	// Update credential
	err := state.UpdateCredential("AWS", "newuser", "newpass")
	if err != nil {
		t.Fatalf("UpdateCredential failed: %v", err)
	}

	// Verify callback invoked
	if !callbackInvoked {
		t.Error("onCredentialsChanged callback was not invoked")
	}

	// Verify credential updated in mock vault
	if mockVault.updateCalled != 1 {
		t.Errorf("Expected UpdateCredential called 1 time, got %d", mockVault.updateCalled)
	}

	// Verify state updated
	creds := state.GetCredentials()
	if len(creds) != 1 {
		t.Errorf("Expected 1 credential, got %d", len(creds))
	}
	if creds[0].Username != "newuser" {
		t.Errorf("Expected username 'newuser', got '%s'", creds[0].Username)
	}
}

// TestDeleteCredential verifies deleting a credential.
func TestDeleteCredential(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Setup existing credentials
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
		{Service: "GitHub", Username: "user", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()

	// Track callback invocation
	callbackInvoked := false
	state.SetOnCredentialsChanged(func() {
		callbackInvoked = true
	})

	// Delete credential
	err := state.DeleteCredential("AWS")
	if err != nil {
		t.Fatalf("DeleteCredential failed: %v", err)
	}

	// Verify callback invoked
	if !callbackInvoked {
		t.Error("onCredentialsChanged callback was not invoked")
	}

	// Verify credential deleted from mock vault
	if mockVault.deleteCalled != 1 {
		t.Errorf("Expected DeleteCredential called 1 time, got %d", mockVault.deleteCalled)
	}

	// Verify state updated
	creds := state.GetCredentials()
	if len(creds) != 1 {
		t.Errorf("Expected 1 credential remaining, got %d", len(creds))
	}
	if creds[0].Service != "GitHub" {
		t.Errorf("Expected remaining service 'GitHub', got '%s'", creds[0].Service)
	}
}

// TestCallbackInvocation_AfterUnlock is the CRITICAL deadlock prevention test.
// It verifies that callbacks are invoked AFTER releasing locks.
func TestCallbackInvocation_AfterUnlock(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Setup mock data
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)

	// CRITICAL TEST: Callback tries to read state (would deadlock if lock not released)
	callbackExecuted := false
	state.SetOnCredentialsChanged(func() {
		// This read would deadlock if callback was invoked while holding lock
		creds := state.GetCredentials()
		if len(creds) > 0 {
			callbackExecuted = true
		}
	})

	// Load credentials (should not deadlock)
	err := state.LoadCredentials()
	if err != nil {
		t.Fatalf("LoadCredentials failed: %v", err)
	}

	// Verify callback executed successfully
	if !callbackExecuted {
		t.Error("Callback was not executed or failed to read state")
	}
}

// TestConcurrentAccess verifies thread-safety with concurrent operations.
func TestConcurrentAccess(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Setup initial data
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()

	// Track callback invocations
	var callbackCount int
	var mu sync.Mutex
	state.SetOnCredentialsChanged(func() {
		mu.Lock()
		callbackCount++
		mu.Unlock()
	})

	// Run concurrent operations
	var wg sync.WaitGroup
	operations := 10

	// Concurrent reads
	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = state.GetCredentials()
			_ = state.GetCategories()
			_ = state.GetSelectedCategory()
		}()
	}

	// Concurrent writes
	for i := 0; i < operations; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			service := "Service" + string(rune(n))
			state.AddCredential(service, "user", "pass")
		}(i)
	}

	// Wait for all operations to complete
	wg.Wait()

	// Verify callbacks were invoked (at least once per add)
	mu.Lock()
	if callbackCount < operations {
		t.Errorf("Expected at least %d callback invocations, got %d", operations, callbackCount)
	}
	mu.Unlock()
}

// TestSetSelectedCategory verifies category selection with callback.
func TestSetSelectedCategory(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Track callback invocation
	callbackInvoked := false
	state.SetOnSelectionChanged(func() {
		callbackInvoked = true
	})

	// Set selected category
	state.SetSelectedCategory("AWS")

	// Verify callback invoked
	if !callbackInvoked {
		t.Error("onSelectionChanged callback was not invoked")
	}

	// Verify category set
	category := state.GetSelectedCategory()
	if category != "AWS" {
		t.Errorf("Expected category 'AWS', got '%s'", category)
	}
}

// TestSetSelectedCredential verifies credential selection with callback.
func TestSetSelectedCredential(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Track callback invocation
	callbackInvoked := false
	state.SetOnSelectionChanged(func() {
		callbackInvoked = true
	})

	// Set selected credential
	cred := &vault.CredentialMetadata{
		Service:  "AWS",
		Username: "admin",
	}
	state.SetSelectedCredential(cred)

	// Verify callback invoked
	if !callbackInvoked {
		t.Error("onSelectionChanged callback was not invoked")
	}

	// Verify credential set
	selected := state.GetSelectedCredential()
	if selected == nil {
		t.Fatal("Expected selected credential, got nil")
	}
	if selected.Service != "AWS" {
		t.Errorf("Expected service 'AWS', got '%s'", selected.Service)
	}
}

// TestComponentStorage verifies component storage and retrieval.
func TestComponentStorage(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Create mock components
	sidebar := tview.NewTreeView()
	table := tview.NewTable()
	detailView := tview.NewTextView()
	statusBar := tview.NewTextView()

	// Store components
	state.SetSidebar(sidebar)
	state.SetTable(table)
	state.SetDetailView(detailView)
	state.SetStatusBar(statusBar)

	// Retrieve and verify components
	if state.GetSidebar() != sidebar {
		t.Error("Sidebar component mismatch")
	}
	if state.GetTable() != table {
		t.Error("Table component mismatch")
	}
	if state.GetDetailView() != detailView {
		t.Error("DetailView component mismatch")
	}
	if state.GetStatusBar() != statusBar {
		t.Error("StatusBar component mismatch")
	}
}

// TestGetFullCredential verifies full credential fetching.
func TestGetFullCredential(t *testing.T) {
	mockVault := NewMockVaultService()
	state := NewAppState(mockVault)

	// Setup mock data
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()

	// Get full credential
	fullCred, err := state.GetFullCredential("AWS")
	if err != nil {
		t.Fatalf("GetFullCredential failed: %v", err)
	}

	// Verify credential data
	if fullCred.Service != "AWS" {
		t.Errorf("Expected service 'AWS', got '%s'", fullCred.Service)
	}
	if fullCred.Password != "mock-password" {
		t.Errorf("Expected password 'mock-password', got '%s'", fullCred.Password)
	}

	// Verify vault method called
	if mockVault.getCalled != 1 {
		t.Errorf("Expected GetCredential called 1 time, got %d", mockVault.getCalled)
	}
}
