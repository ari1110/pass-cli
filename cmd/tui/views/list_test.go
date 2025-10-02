package views

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/internal/vault"
)

func createTestCredentials() []vault.CredentialMetadata {
	return []vault.CredentialMetadata{
		{
			Service:   "github.com",
			Username:  "user1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Service:   "gitlab.com",
			Username:  "user2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Service:   "example.com",
			Username:  "admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func TestNewListView(t *testing.T) {
	credentials := createTestCredentials()
	view := NewListView(credentials)

	if view == nil {
		t.Fatal("NewListView() returned nil")
	}

	if len(view.allItems) != len(credentials) {
		t.Errorf("Expected %d items, got %d", len(credentials), len(view.allItems))
	}

	if view.searchFocused {
		t.Error("Search should not be focused initially")
	}
}

func TestListViewSetSize(t *testing.T) {
	credentials := createTestCredentials()
	view := NewListView(credentials)

	view.SetSize(100, 40)

	if view.width != 100 {
		t.Errorf("Expected width 100, got %d", view.width)
	}

	if view.height != 40 {
		t.Errorf("Expected height 40, got %d", view.height)
	}
}

func TestListViewSearchFocus(t *testing.T) {
	credentials := createTestCredentials()
	view := NewListView(credentials)
	view.SetSize(80, 24)

	// Test slash key focuses search
	slashKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	view, _ = view.Update(slashKey)

	if !view.searchFocused {
		t.Error("Search should be focused after pressing '/'")
	}

	// Test escape clears search when focused
	escKey := tea.KeyMsg{Type: tea.KeyEsc}
	view, _ = view.Update(escKey)

	if view.searchFocused {
		t.Error("Search should not be focused after pressing Esc")
	}

	if view.searchInput.Value() != "" {
		t.Error("Search input should be cleared after pressing Esc")
	}
}

func TestListViewNavigation(t *testing.T) {
	credentials := createTestCredentials()
	view := NewListView(credentials)
	view.SetSize(80, 24)

	tests := []struct {
		name string
		key  tea.KeyMsg
	}{
		{
			name: "Down arrow",
			key:  tea.KeyMsg{Type: tea.KeyDown},
		},
		{
			name: "Up arrow",
			key:  tea.KeyMsg{Type: tea.KeyUp},
		},
		{
			name: "j key (vim down)",
			key:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
		},
		{
			name: "k key (vim up)",
			key:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify these keys don't cause panics
			_, _ = view.Update(tt.key)
		})
	}
}

func TestListViewTabSwitchesFocus(t *testing.T) {
	credentials := createTestCredentials()
	view := NewListView(credentials)
	view.SetSize(80, 24)

	// Focus search first
	slashKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	view, _ = view.Update(slashKey)

	if !view.searchFocused {
		t.Fatal("Search should be focused")
	}

	// Press Tab to switch focus back to list
	tabKey := tea.KeyMsg{Type: tea.KeyTab}
	view, _ = view.Update(tabKey)

	if view.searchFocused {
		t.Error("Tab should switch focus from search to list")
	}

	// Press Tab again to switch focus back to search
	view, _ = view.Update(tabKey)

	if !view.searchFocused {
		t.Error("Tab should switch focus from list to search")
	}
}

func TestListViewSelectedCredential(t *testing.T) {
	credentials := createTestCredentials()
	view := NewListView(credentials)
	view.SetSize(80, 24)

	// Get selected credential (should be first one initially)
	selected := view.SelectedCredential()
	if selected == nil {
		t.Fatal("SelectedCredential() returned nil")
	}

	if selected.Service != credentials[0].Service {
		t.Errorf("Expected selected service %s, got %s", credentials[0].Service, selected.Service)
	}
}

func TestListViewEmptyCredentials(t *testing.T) {
	view := NewListView([]vault.CredentialMetadata{})

	if view == nil {
		t.Fatal("NewListView() returned nil for empty credentials")
	}

	if len(view.allItems) != 0 {
		t.Errorf("Expected 0 items, got %d", len(view.allItems))
	}

	selected := view.SelectedCredential()
	if selected != nil {
		t.Error("SelectedCredential() should return nil for empty list")
	}
}

func TestCredentialItemFilterValue(t *testing.T) {
	item := credentialItem{
		metadata: vault.CredentialMetadata{
			Service:  "github.com",
			Username: "testuser",
		},
	}

	filterValue := item.FilterValue()
	if filterValue != "github.com testuser" {
		t.Errorf("Expected 'github.com testuser', got '%s'", filterValue)
	}
}

func TestCredentialItemDescription(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected string
	}{
		{
			name:     "With username",
			username: "testuser",
			expected: "testuser",
		},
		{
			name:     "Without username",
			username: "",
			expected: "(no username)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := credentialItem{
				metadata: vault.CredentialMetadata{
					Service:  "github.com",
					Username: tt.username,
				},
			}

			desc := item.Description()
			if desc != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, desc)
			}
		})
	}
}
