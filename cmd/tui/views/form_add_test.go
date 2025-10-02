package views

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewAddFormView(t *testing.T) {
	view := NewAddFormView()

	if view == nil {
		t.Fatal("NewAddFormView() returned nil")
	}

	if view.focusedField != FieldService {
		t.Error("Initial focus should be on service field")
	}

	if view.errorMsg != "" {
		t.Error("Error message should be empty initially")
	}

	if view.notification != "" {
		t.Error("Notification should be empty initially")
	}
}

func TestAddFormViewSetSize(t *testing.T) {
	view := NewAddFormView()

	view.SetSize(100, 40)

	if view.width != 100 {
		t.Errorf("Expected width 100, got %d", view.width)
	}

	if view.height != 40 {
		t.Errorf("Expected height 40, got %d", view.height)
	}
}

func TestAddFormViewFieldNavigation(t *testing.T) {
	view := NewAddFormView()
	view.SetSize(80, 24)

	tests := []struct {
		name          string
		key           tea.KeyMsg
		startField    FormField
		expectedField FormField
	}{
		{
			name:          "Tab from service to username",
			key:           tea.KeyMsg{Type: tea.KeyTab},
			startField:    FieldService,
			expectedField: FieldUsername,
		},
		{
			name:          "Tab from username to password",
			key:           tea.KeyMsg{Type: tea.KeyTab},
			startField:    FieldUsername,
			expectedField: FieldPassword,
		},
		{
			name:          "Tab from password to notes",
			key:           tea.KeyMsg{Type: tea.KeyTab},
			startField:    FieldPassword,
			expectedField: FieldNotes,
		},
		{
			name:          "Tab from notes wraps to service",
			key:           tea.KeyMsg{Type: tea.KeyTab},
			startField:    FieldNotes,
			expectedField: FieldService,
		},
		{
			name:          "Down arrow navigates forward",
			key:           tea.KeyMsg{Type: tea.KeyDown},
			startField:    FieldService,
			expectedField: FieldUsername,
		},
		{
			name:          "Up arrow navigates backward",
			key:           tea.KeyMsg{Type: tea.KeyUp},
			startField:    FieldUsername,
			expectedField: FieldService,
		},
		{
			name:          "Shift+Tab navigates backward",
			key:           tea.KeyMsg{Type: tea.KeyShiftTab},
			startField:    FieldUsername,
			expectedField: FieldService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view.focusedField = tt.startField
			updatedView, _ := view.Update(tt.key)

			if updatedView.focusedField != tt.expectedField {
				t.Errorf("Expected field %v, got %v", tt.expectedField, updatedView.focusedField)
			}
		})
	}
}

func TestAddFormViewGetValues(t *testing.T) {
	view := NewAddFormView()
	view.SetSize(80, 24)

	// Set values directly
	view.serviceInput.SetValue("github.com")
	view.usernameInput.SetValue("testuser")
	view.passwordInput.SetValue("testpass123")
	view.notesInput.SetValue("Test notes")

	service, username, password, notes := view.GetValues()

	if service != "github.com" {
		t.Errorf("Expected service 'github.com', got '%s'", service)
	}

	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}

	if password != "testpass123" {
		t.Errorf("Expected password 'testpass123', got '%s'", password)
	}

	if notes != "Test notes" {
		t.Errorf("Expected notes 'Test notes', got '%s'", notes)
	}
}

func TestAddFormViewSetPassword(t *testing.T) {
	view := NewAddFormView()

	view.SetPassword("generated-password-123")

	_, _, password, _ := view.GetValues()
	if password != "generated-password-123" {
		t.Errorf("Expected password 'generated-password-123', got '%s'", password)
	}
}

func TestAddFormViewSetError(t *testing.T) {
	view := NewAddFormView()

	view.SetError("Test error message")

	if view.errorMsg != "Test error message" {
		t.Errorf("Expected error 'Test error message', got '%s'", view.errorMsg)
	}
}

func TestAddFormViewSetNotification(t *testing.T) {
	view := NewAddFormView()

	view.SetNotification("Test notification")

	if view.notification != "Test notification" {
		t.Errorf("Expected notification 'Test notification', got '%s'", view.notification)
	}
}

func TestAddFormViewHasChanges(t *testing.T) {
	view := NewAddFormView()

	// Initially no changes
	if view.HasChanges() {
		t.Error("HasChanges() should return false for empty form")
	}

	// Add service
	view.serviceInput.SetValue("github.com")
	if !view.HasChanges() {
		t.Error("HasChanges() should return true when service is set")
	}

	// Clear and add username
	view = NewAddFormView()
	view.usernameInput.SetValue("testuser")
	if !view.HasChanges() {
		t.Error("HasChanges() should return true when username is set")
	}

	// Clear and add password
	view = NewAddFormView()
	view.passwordInput.SetValue("testpass")
	if !view.HasChanges() {
		t.Error("HasChanges() should return true when password is set")
	}

	// Clear and add notes
	view = NewAddFormView()
	view.notesInput.SetValue("test notes")
	if !view.HasChanges() {
		t.Error("HasChanges() should return true when notes are set")
	}
}

func TestAddFormViewRendersCorrectly(t *testing.T) {
	view := NewAddFormView()
	view.SetSize(80, 24)

	output := view.View()

	// Verify key elements are present
	if !strings.Contains(output, "Add New Credential") {
		t.Error("View should contain title 'Add New Credential'")
	}

	if !strings.Contains(output, "Service") {
		t.Error("View should contain 'Service' label")
	}

	if !strings.Contains(output, "Username") {
		t.Error("View should contain 'Username' label")
	}

	if !strings.Contains(output, "Password") {
		t.Error("View should contain 'Password' label")
	}

	if !strings.Contains(output, "Notes") {
		t.Error("View should contain 'Notes' label")
	}
}

func TestAddFormViewErrorDisplay(t *testing.T) {
	view := NewAddFormView()
	view.SetSize(80, 24)
	view.SetError("Service name is required")

	output := view.View()

	if !strings.Contains(output, "Service name is required") {
		t.Error("View should display error message")
	}
}

func TestAddFormViewNotificationDisplay(t *testing.T) {
	view := NewAddFormView()
	view.SetSize(80, 24)
	view.SetNotification("Password generated")

	output := view.View()

	if !strings.Contains(output, "Password generated") {
		t.Error("View should display notification message")
	}
}
