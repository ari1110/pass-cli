package components

import (
	"errors"
	"sync"
	"testing"
	"time"

	"pass-cli/cmd/tui-tview/models"
	"pass-cli/internal/vault"
)

// MockVaultService for component tests
type MockVaultService struct {
	mu          sync.Mutex
	credentials []vault.CredentialMetadata
}

func NewMockVaultService() *MockVaultService {
	return &MockVaultService{credentials: make([]vault.CredentialMetadata, 0)}
}

func (m *MockVaultService) ListCredentialsWithMetadata() ([]vault.CredentialMetadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.credentials, nil
}

func (m *MockVaultService) AddCredential(service, username, password, category string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.credentials = append(m.credentials, vault.CredentialMetadata{
		Service: service, Username: username, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	return nil
}

func (m *MockVaultService) UpdateCredential(service, username, password, category string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, cred := range m.credentials {
		if cred.Service == service {
			m.credentials[i].Username = username
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MockVaultService) DeleteCredential(service string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, cred := range m.credentials {
		if cred.Service == service {
			m.credentials = append(m.credentials[:i], m.credentials[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MockVaultService) GetCredential(service string, trackUsage bool) (*vault.Credential, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, cred := range m.credentials {
		if cred.Service == service {
			return &vault.Credential{Service: cred.Service, Username: cred.Username, Password: "mock"}, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockVaultService) SetCredentials(creds []vault.CredentialMetadata) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.credentials = creds
}

// TestNewSidebar verifies Sidebar creation.
func TestNewSidebar(t *testing.T) {
	mockVault := NewMockVaultService()
	state := models.NewAppState(mockVault)

	sidebar := NewSidebar(state)

	if sidebar == nil {
		t.Fatal("NewSidebar returned nil")
	}

	// Verify root node exists
	if sidebar.rootNode == nil {
		t.Error("Root node is nil")
	}

	// Verify root node text
	if sidebar.rootNode.GetText() != "All Credentials" {
		t.Errorf("Expected root text 'All Credentials', got '%s'", sidebar.rootNode.GetText())
	}

	// Verify root is expanded
	if !sidebar.rootNode.IsExpanded() {
		t.Error("Root node should be expanded")
	}
}

// TestSidebarRefresh verifies tree rebuilding.
func TestSidebarRefresh(t *testing.T) {
	mockVault := NewMockVaultService()
	state := models.NewAppState(mockVault)

	// Setup mock credentials with different services (categories)
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
		{Service: "GitHub", Username: "user", CreatedAt: time.Now()},
		{Service: "Database", Username: "dbuser", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()

	sidebar := NewSidebar(state)

	// Refresh should rebuild tree with categories
	sidebar.Refresh()

	// Verify categories added as children
	children := sidebar.rootNode.GetChildren()
	if len(children) != 3 {
		t.Errorf("Expected 3 category nodes, got %d", len(children))
	}

	// Verify category names (should be sorted)
	expectedCategories := []string{"AWS", "Database", "GitHub"}
	for i, child := range children {
		if child.GetText() != expectedCategories[i] {
			t.Errorf("Expected category '%s' at index %d, got '%s'", expectedCategories[i], i, child.GetText())
		}
	}

	// Verify root still expanded after refresh
	if !sidebar.rootNode.IsExpanded() {
		t.Error("Root node should remain expanded after refresh")
	}
}

// TestSidebarRefresh_EmptyCategories verifies handling of empty categories.
func TestSidebarRefresh_EmptyCategories(t *testing.T) {
	mockVault := NewMockVaultService()
	state := models.NewAppState(mockVault)

	sidebar := NewSidebar(state)

	// Refresh with no credentials
	sidebar.Refresh()

	// Verify no children added
	children := sidebar.rootNode.GetChildren()
	if len(children) != 0 {
		t.Errorf("Expected 0 category nodes for empty state, got %d", len(children))
	}

	// Verify root still exists and is expanded
	if sidebar.rootNode == nil {
		t.Error("Root node should still exist")
	}
	if !sidebar.rootNode.IsExpanded() {
		t.Error("Root node should be expanded even when empty")
	}
}

// TestSidebarSelection_RootNode verifies root node selection behavior.
func TestSidebarSelection_RootNode(t *testing.T) {
	mockVault := NewMockVaultService()
	state := models.NewAppState(mockVault)

	// Setup categories
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
		{Service: "GitHub", Username: "user", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()

	sidebar := NewSidebar(state)

	// Select root node (simulates user clicking "All Credentials")
	sidebar.onSelect(sidebar.rootNode)

	// Verify selected category is empty (shows all credentials)
	selectedCategory := state.GetSelectedCategory()
	if selectedCategory != "" {
		t.Errorf("Expected empty category (show all), got '%s'", selectedCategory)
	}
}

// TestSidebarSelection_CategoryNode verifies category node selection behavior.
func TestSidebarSelection_CategoryNode(t *testing.T) {
	mockVault := NewMockVaultService()
	state := models.NewAppState(mockVault)

	// Setup categories
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
		{Service: "GitHub", Username: "user", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()

	sidebar := NewSidebar(state)

	// Get first category node (AWS)
	children := sidebar.rootNode.GetChildren()
	if len(children) == 0 {
		t.Fatal("Expected category nodes, got none")
	}

	categoryNode := children[0] // AWS (sorted first)

	// Select category node
	sidebar.onSelect(categoryNode)

	// Verify selected category updated in state
	selectedCategory := state.GetSelectedCategory()
	if selectedCategory != "AWS" {
		t.Errorf("Expected selected category 'AWS', got '%s'", selectedCategory)
	}
}

// TestSidebarSelection_UpdatesAppState verifies AppState is updated on selection.
func TestSidebarSelection_UpdatesAppState(t *testing.T) {
	mockVault := NewMockVaultService()
	state := models.NewAppState(mockVault)

	// Track selection changes
	selectionChanged := false
	state.SetOnSelectionChanged(func() {
		selectionChanged = true
	})

	// Setup categories
	mockCreds := []vault.CredentialMetadata{
		{Service: "GitHub", Username: "user", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()

	sidebar := NewSidebar(state)

	// Select category
	children := sidebar.rootNode.GetChildren()
	if len(children) > 0 {
		sidebar.onSelect(children[0])
	}

	// Verify callback invoked
	if !selectionChanged {
		t.Error("Selection change callback was not invoked")
	}

	// Verify state updated
	if state.GetSelectedCategory() != "GitHub" {
		t.Errorf("Expected selected category 'GitHub', got '%s'", state.GetSelectedCategory())
	}
}

// TestSidebarRefresh_PreservesRootExpansion verifies root remains expanded after refresh.
func TestSidebarRefresh_PreservesRootExpansion(t *testing.T) {
	mockVault := NewMockVaultService()
	state := models.NewAppState(mockVault)

	sidebar := NewSidebar(state)

	// Verify initial expansion
	if !sidebar.rootNode.IsExpanded() {
		t.Error("Root should be expanded initially")
	}

	// Manually collapse root (simulates user collapsing)
	sidebar.rootNode.SetExpanded(false)

	// Refresh should re-expand root
	sidebar.Refresh()

	// Verify root is expanded again
	if !sidebar.rootNode.IsExpanded() {
		t.Error("Root should be re-expanded after refresh")
	}
}

// TestSidebarRefresh_ClearsOldCategories verifies old categories are removed.
func TestSidebarRefresh_ClearsOldCategories(t *testing.T) {
	mockVault := NewMockVaultService()
	state := models.NewAppState(mockVault)

	sidebar := NewSidebar(state)

	// Setup initial categories
	mockCreds := []vault.CredentialMetadata{
		{Service: "AWS", Username: "admin", CreatedAt: time.Now()},
		{Service: "GitHub", Username: "user", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(mockCreds)
	state.LoadCredentials()
	sidebar.Refresh()

	// Verify 2 categories
	if len(sidebar.rootNode.GetChildren()) != 2 {
		t.Errorf("Expected 2 categories initially, got %d", len(sidebar.rootNode.GetChildren()))
	}

	// Update to new categories (different set)
	newCreds := []vault.CredentialMetadata{
		{Service: "Database", Username: "dbuser", CreatedAt: time.Now()},
	}
	mockVault.SetCredentials(newCreds)
	state.LoadCredentials()
	sidebar.Refresh()

	// Verify old categories cleared, only new one present
	children := sidebar.rootNode.GetChildren()
	if len(children) != 1 {
		t.Errorf("Expected 1 category after refresh, got %d", len(children))
	}
	if children[0].GetText() != "Database" {
		t.Errorf("Expected category 'Database', got '%s'", children[0].GetText())
	}
}
