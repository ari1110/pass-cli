package components

import (
	"testing"
	"time"

	"pass-cli/internal/vault"
)

func TestNewSidebarPanel(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "test", Username: "user"},
	}
	sidebar := NewSidebarPanel(credentials)

	if sidebar == nil {
		t.Fatal("NewSidebarPanel returned nil")
	}

	if sidebar.width != 0 {
		t.Errorf("Expected initial width 0, got %d", sidebar.width)
	}

	if sidebar.height != 0 {
		t.Errorf("Expected initial height 0, got %d", sidebar.height)
	}

	if sidebar.focused {
		t.Error("Sidebar should not be focused initially")
	}

	if sidebar.stats.Total != 1 {
		t.Errorf("Expected 1 credential in stats, got %d", sidebar.stats.Total)
	}
}

func TestSidebarPanel_SetSize(t *testing.T) {
	sidebar := NewSidebarPanel([]vault.CredentialMetadata{})

	sidebar.SetSize(30, 20)

	if sidebar.width != 30 {
		t.Errorf("Expected width 30, got %d", sidebar.width)
	}

	if sidebar.height != 20 {
		t.Errorf("Expected height 20, got %d", sidebar.height)
	}
}

func TestSidebarPanel_SetFocus(t *testing.T) {
	sidebar := NewSidebarPanel([]vault.CredentialMetadata{})

	// Initially not focused
	if sidebar.focused {
		t.Error("Sidebar should not be focused initially")
	}

	// Set focus
	sidebar.SetFocus(true)
	if !sidebar.focused {
		t.Error("Sidebar should be focused after SetFocus(true)")
	}

	// Remove focus
	sidebar.SetFocus(false)
	if sidebar.focused {
		t.Error("Sidebar should not be focused after SetFocus(false)")
	}
}

func TestSidebarPanel_UpdateCredentials(t *testing.T) {
	sidebar := NewSidebarPanel([]vault.CredentialMetadata{})
	sidebar.SetSize(30, 20)

	credentials := []vault.CredentialMetadata{
		{Service: "aws-prod", Username: "admin", CreatedAt: time.Now()},
		{Service: "github-repo", Username: "dev", CreatedAt: time.Now()},
		{Service: "postgres-db", Username: "user", CreatedAt: time.Now()},
	}

	sidebar.UpdateCredentials(credentials)

	// Should have categorized credentials
	if len(sidebar.categories) == 0 {
		t.Error("Categories should not be empty after UpdateCredentials")
	}

	// Should have stats
	if sidebar.stats.Total != 3 {
		t.Errorf("Expected total stats 3, got %d", sidebar.stats.Total)
	}
}

func TestSidebarPanel_GetSelectedCategory_Empty(t *testing.T) {
	sidebar := NewSidebarPanel([]vault.CredentialMetadata{})

	category := sidebar.GetSelectedCategory()

	if category != "" {
		t.Errorf("Expected empty category for empty sidebar, got %s", category)
	}
}

func TestSidebarPanel_GetSelectedCredential_Empty(t *testing.T) {
	sidebar := NewSidebarPanel([]vault.CredentialMetadata{})

	cred := sidebar.GetSelectedCredential()

	if cred != nil {
		t.Error("Expected nil credential for empty sidebar")
	}
}

func TestSidebarPanel_GetSelectedCategory_WithData(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "aws-prod", Username: "admin"},
		{Service: "github-repo", Username: "dev"},
	}
	sidebar := NewSidebarPanel(credentials)
	sidebar.SetSize(30, 20)

	// Should have at least one category
	if len(sidebar.categories) == 0 {
		t.Fatal("Expected categories after NewSidebarPanel")
	}

	category := sidebar.GetSelectedCategory()

	// Should return the name of the first category
	if category == "" {
		t.Error("Expected non-empty category name")
	}

	if category != sidebar.categories[0].Name {
		t.Errorf("Expected category %s, got %s", sidebar.categories[0].Name, category)
	}
}

func TestSidebarPanel_Navigation_MoveDown(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "aws-prod", Username: "admin"},
		{Service: "aws-dev", Username: "user"},
		{Service: "github-repo", Username: "dev"},
	}
	sidebar := NewSidebarPanel(credentials)
	sidebar.SetSize(30, 20)

	// Initially at first category/credential
	initialCategory := sidebar.selectedCategory
	initialCred := sidebar.selectedCred

	// Move down
	sidebar.moveDown()

	// Should have moved
	if sidebar.selectedCategory == initialCategory && sidebar.selectedCred == initialCred {
		t.Error("moveDown should change selection")
	}
}

func TestSidebarPanel_Navigation_MoveUp(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "aws-prod", Username: "admin"},
		{Service: "github-repo", Username: "dev"},
	}
	sidebar := NewSidebarPanel(credentials)
	sidebar.SetSize(30, 20)

	// Move down first to have something to move up from
	sidebar.moveDown()
	afterDownCategory := sidebar.selectedCategory
	afterDownCred := sidebar.selectedCred

	// Move up
	sidebar.moveUp()

	// Should have moved back (or wrapped)
	if sidebar.selectedCategory == afterDownCategory && sidebar.selectedCred == afterDownCred {
		t.Error("moveUp should change selection")
	}
}

func TestSidebarPanel_View(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "aws-prod", Username: "admin"},
	}
	sidebar := NewSidebarPanel(credentials)
	sidebar.SetSize(30, 20)

	output := sidebar.View()

	if output == "" {
		t.Error("View() should return non-empty output")
	}
}

func TestCountUsedCredentials(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "service1", UsageCount: 5},
		{Service: "service2", UsageCount: 10},
		{Service: "service3", UsageCount: 0}, // Never used
	}

	count := countUsedCredentials(credentials)

	if count != 2 {
		t.Errorf("Expected 2 used credentials, got %d", count)
	}
}

func TestCountUsedCredentials_Empty(t *testing.T) {
	credentials := []vault.CredentialMetadata{}

	count := countUsedCredentials(credentials)

	if count != 0 {
		t.Errorf("Expected 0 used credentials, got %d", count)
	}
}

func TestCountUsedCredentials_NoneUsed(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "service1"}, // Never used
		{Service: "service2"}, // Never used
	}

	count := countUsedCredentials(credentials)

	if count != 0 {
		t.Errorf("Expected 0 used credentials, got %d", count)
	}
}

func TestCountRecentCredentials(t *testing.T) {
	// Currently countRecentCredentials just returns count of used credentials
	// This test verifies the current implementation
	credentials := []vault.CredentialMetadata{
		{Service: "used1", UsageCount: 5},
		{Service: "used2", UsageCount: 3},
		{Service: "unused", UsageCount: 0},
	}

	count := countRecentCredentials(credentials)

	if count != 2 {
		t.Errorf("Expected 2 recent credentials (used credentials), got %d", count)
	}
}

func TestCountRecentCredentials_Empty(t *testing.T) {
	credentials := []vault.CredentialMetadata{}

	count := countRecentCredentials(credentials)

	if count != 0 {
		t.Errorf("Expected 0 recent credentials, got %d", count)
	}
}

func TestCountRecentCredentials_NoneRecent(t *testing.T) {
	// Currently countRecentCredentials returns count of used credentials
	credentials := []vault.CredentialMetadata{
		{Service: "unused1", UsageCount: 0},
		{Service: "unused2", UsageCount: 0},
	}

	count := countRecentCredentials(credentials)

	if count != 0 {
		t.Errorf("Expected 0 recent credentials, got %d", count)
	}
}

func TestTruncate(t *testing.T) {
	testCases := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"short", 10, "short"},
		{"exactly ten", 11, "exactly ten"},
		{"this is a very long string", 10, "this is..."},
		{"test", 3, "..."},
		{"", 10, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := truncate(tc.input, tc.maxLen)
			if result != tc.expected {
				t.Errorf("truncate(%q, %d) = %q, expected %q", tc.input, tc.maxLen, result, tc.expected)
			}
		})
	}
}

func TestSidebarPanel_GetSelectedCredential_WithData(t *testing.T) {
	credentials := []vault.CredentialMetadata{
		{Service: "aws-prod", Username: "admin"},
		{Service: "aws-dev", Username: "user"},
	}
	sidebar := NewSidebarPanel(credentials)
	sidebar.SetSize(30, 20)

	// Should have categories
	if len(sidebar.categories) == 0 {
		t.Fatal("Expected categories after NewSidebarPanel")
	}

	// Initial selection is on category header (selectedCred == -1)
	// Move down to select first credential
	sidebar.moveDown()

	// Get selected credential
	cred := sidebar.GetSelectedCredential()

	if cred == nil {
		t.Error("Expected non-nil credential after moveDown")
	}

	// Verify it's one of the credentials we added
	if cred != nil && cred.Service != "aws-prod" && cred.Service != "aws-dev" {
		t.Errorf("Expected credential from our test data, got service: %s", cred.Service)
	}
}

func TestSidebarPanel_Stats(t *testing.T) {
	now := time.Now()
	credentials := []vault.CredentialMetadata{
		{Service: "service1", Username: "user1", CreatedAt: now.Add(-1 * time.Hour), UsageCount: 5},
		{Service: "service2", Username: "user2", CreatedAt: now.Add(-2 * time.Hour), UsageCount: 3},
		{Service: "service3", Username: "user3", CreatedAt: now.Add(-10 * 24 * time.Hour), UsageCount: 0}, // Old, never used
	}
	sidebar := NewSidebarPanel(credentials)
	sidebar.SetSize(30, 20)

	// Check stats
	if sidebar.stats.Total != 3 {
		t.Errorf("Expected total 3, got %d", sidebar.stats.Total)
	}

	if sidebar.stats.Used != 2 {
		t.Errorf("Expected 2 used credentials, got %d", sidebar.stats.Used)
	}

	if sidebar.stats.Recent != 2 {
		t.Errorf("Expected 2 recent credentials, got %d", sidebar.stats.Recent)
	}
}
