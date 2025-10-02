package views

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/internal/vault"
)

func createTestCredential() *vault.Credential {
	return &vault.Credential{
		Service:  "github.com",
		Username: "testuser",
		Password: "testpassword123",
		Notes:    "Test notes",
		CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
		UsageRecord: map[string]vault.UsageRecord{
			"/home/user/project": {
				Location:  "/home/user/project",
				GitRepo:   "github.com/user/repo",
				Timestamp: time.Date(2024, 1, 3, 12, 0, 0, 0, time.UTC),
				Count:     5,
			},
		},
	}
}

func TestNewDetailView(t *testing.T) {
	cred := createTestCredential()
	view := NewDetailView(cred)

	if view == nil {
		t.Fatal("NewDetailView() returned nil")
	}

	if view.credential != cred {
		t.Error("Credential was not set correctly")
	}

	if !view.passwordMasked {
		t.Error("Password should be masked by default")
	}

	if view.notification != "" {
		t.Error("Notification should be empty initially")
	}
}

func TestDetailViewSetSize(t *testing.T) {
	cred := createTestCredential()
	view := NewDetailView(cred)

	view.SetSize(100, 40)

	if view.width != 100 {
		t.Errorf("Expected width 100, got %d", view.width)
	}

	if view.height != 40 {
		t.Errorf("Expected height 40, got %d", view.height)
	}
}

func TestDetailViewPasswordToggle(t *testing.T) {
	cred := createTestCredential()
	view := NewDetailView(cred)
	view.SetSize(80, 24)

	// Initially masked
	if !view.passwordMasked {
		t.Fatal("Password should be masked initially")
	}

	// Press 'm' to unmask
	mKey := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}}
	view, _ = view.Update(mKey)

	if view.passwordMasked {
		t.Error("Password should be unmasked after pressing 'm'")
	}

	// Press 'm' again to mask
	view, _ = view.Update(mKey)

	if !view.passwordMasked {
		t.Error("Password should be masked after pressing 'm' again")
	}
}

func TestDetailViewGetCredential(t *testing.T) {
	cred := createTestCredential()
	view := NewDetailView(cred)

	retrieved := view.GetCredential()
	if retrieved != cred {
		t.Error("GetCredential() did not return the correct credential")
	}
}

func TestDetailViewRendersCorrectly(t *testing.T) {
	cred := createTestCredential()
	view := NewDetailView(cred)
	view.SetSize(80, 24)

	output := view.View()

	// Verify key elements are present
	if !strings.Contains(output, "Credential Details") {
		t.Error("View should contain title 'Credential Details'")
	}

	if !strings.Contains(output, cred.Service) {
		t.Error("View should contain service name")
	}

	if !strings.Contains(output, cred.Username) {
		t.Error("View should contain username")
	}

	// Password should be masked
	if strings.Contains(output, cred.Password) {
		t.Error("View should not contain raw password when masked")
	}

	if !strings.Contains(output, strings.Repeat("*", len(cred.Password))) {
		t.Error("View should contain masked password")
	}

	if !strings.Contains(output, cred.Notes) {
		t.Error("View should contain notes")
	}
}

func TestDetailViewPasswordVisibility(t *testing.T) {
	cred := createTestCredential()
	view := NewDetailView(cred)
	view.SetSize(80, 24)

	// Unmask password
	view.passwordMasked = false
	view.updateContent()
	output := view.View()

	// Raw password should be visible
	if !strings.Contains(output, cred.Password) {
		t.Error("View should contain raw password when unmasked")
	}

	// Should not contain masked password
	if strings.Contains(output, strings.Repeat("*", len(cred.Password))) {
		t.Error("View should not contain masked password when unmasked")
	}
}

func TestDetailViewUsageRecords(t *testing.T) {
	cred := createTestCredential()
	view := NewDetailView(cred)
	view.SetSize(80, 24)

	output := view.View()

	// Verify usage records are displayed
	if !strings.Contains(output, "Usage Records") {
		t.Error("View should contain 'Usage Records' header")
	}

	if !strings.Contains(output, "/home/user/project") {
		t.Error("View should contain usage location")
	}

	if !strings.Contains(output, "github.com/user/repo") {
		t.Error("View should contain git repo")
	}
}

func TestDetailViewNoUsageRecords(t *testing.T) {
	cred := &vault.Credential{
		Service:     "github.com",
		Username:    "testuser",
		Password:    "testpassword123",
		Notes:       "Test notes",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UsageRecord: map[string]vault.UsageRecord{},
	}

	view := NewDetailView(cred)
	view.SetSize(80, 24)

	output := view.View()

	// Verify "None" is displayed for usage records
	if !strings.Contains(output, "Usage Records") {
		t.Error("View should contain 'Usage Records' label")
	}

	if !strings.Contains(output, "None") {
		t.Error("View should contain 'None' for empty usage records")
	}
}

func TestDetailViewEmptyUsername(t *testing.T) {
	cred := &vault.Credential{
		Service:     "github.com",
		Username:    "",
		Password:    "testpassword123",
		Notes:       "Test notes",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UsageRecord: map[string]vault.UsageRecord{},
	}

	view := NewDetailView(cred)
	view.SetSize(80, 24)

	output := view.View()

	// Verify "(not set)" is displayed for empty username
	if !strings.Contains(output, "(not set)") {
		t.Error("View should contain '(not set)' for empty username")
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Valid time",
			time:     time.Date(2024, 1, 1, 12, 30, 45, 0, time.UTC),
			expected: "2024-01-01 12:30:45",
		},
		{
			name:     "Zero time",
			time:     time.Time{},
			expected: "Never",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTime(tt.time)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
