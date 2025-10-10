package components

import (
	"errors"
	"sync"
	"testing"
	"time"

	"pass-cli/cmd/tui/models"
	"pass-cli/internal/vault"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// mockVaultServiceForForms implements models.VaultService for form tests
type mockVaultServiceForForms struct {
	mu          sync.Mutex
	credentials []vault.CredentialMetadata
}

func newMockVaultServiceForForms() *mockVaultServiceForForms {
	return &mockVaultServiceForForms{credentials: make([]vault.CredentialMetadata, 0)}
}

func (m *mockVaultServiceForForms) ListCredentialsWithMetadata() ([]vault.CredentialMetadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.credentials, nil
}

func (m *mockVaultServiceForForms) AddCredential(service, username, password, category, url, notes string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.credentials = append(m.credentials, vault.CredentialMetadata{
		Service: service, Username: username, Category: category, URL: url, Notes: notes,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	})
	return nil
}

func (m *mockVaultServiceForForms) UpdateCredential(service string, opts vault.UpdateOpts) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, cred := range m.credentials {
		if cred.Service == service {
			if opts.Username != nil {
				m.credentials[i].Username = *opts.Username
			}
			if opts.Category != nil {
				m.credentials[i].Category = *opts.Category
			}
			if opts.URL != nil {
				m.credentials[i].URL = *opts.URL
			}
			if opts.Notes != nil {
				m.credentials[i].Notes = *opts.Notes
			}
			return nil
		}
	}
	return errors.New("not found")
}

func (m *mockVaultServiceForForms) DeleteCredential(service string) error {
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

func (m *mockVaultServiceForForms) GetCredential(service string, trackUsage bool) (*vault.Credential, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, cred := range m.credentials {
		if cred.Service == service {
			return &vault.Credential{Service: cred.Service, Username: cred.Username, Password: "mock"}, nil
		}
	}
	return nil, errors.New("not found")
}

// TestAddFormPasswordVisibilityToggle verifies the toggle changes label
// T004: Unit test for AddForm password visibility toggle functionality
// NOTE: tview InputField doesn't expose GetMaskCharacter(), so we test via label changes
func TestAddFormPasswordVisibilityToggle(t *testing.T) {
	// Setup
	mockVault := newMockVaultServiceForForms()
	appState := models.NewAppState(mockVault)
	form := NewAddForm(appState)

	// Get password field (index 2: Service=0, Username=1, Password=2)
	passwordField := form.GetFormItem(2).(*tview.InputField)

	// Test initial state - label should show "Password"
	t.Run("InitialStateMasked", func(t *testing.T) {
		if passwordField.GetLabel() != "Password" {
			t.Errorf("Expected initial label 'Password', got '%s'", passwordField.GetLabel())
		}
	})

	// Test toggle to visible - this will FAIL until implementation
	t.Run("ToggleToVisible", func(t *testing.T) {
		// Simulate Ctrl+H key event
		event := tcell.NewEventKey(tcell.KeyCtrlH, 0, tcell.ModNone)
		result := form.GetInputCapture()(event)

		if result != nil {
			t.Error("Expected Ctrl+H to be consumed (return nil)")
		}

		expectedLabel := "Password [VISIBLE]"
		if passwordField.GetLabel() != expectedLabel {
			t.Errorf("Expected label '%s', got '%s'", expectedLabel, passwordField.GetLabel())
		}
	})

	// Test toggle back to masked
	t.Run("ToggleBackToMasked", func(t *testing.T) {
		// Simulate Ctrl+H again
		event := tcell.NewEventKey(tcell.KeyCtrlH, 0, tcell.ModNone)
		form.GetInputCapture()(event)

		if passwordField.GetLabel() != "Password" {
			t.Errorf("Expected label 'Password', got '%s'", passwordField.GetLabel())
		}
	})
}

// TestAddFormCtrlHShortcut verifies Ctrl+H key event is consumed
// T005: Unit test for Ctrl+H keyboard shortcut handling
func TestAddFormCtrlHShortcut(t *testing.T) {
	mockVault := newMockVaultServiceForForms()
	appState := models.NewAppState(mockVault)
	form := NewAddForm(appState)

	t.Run("CtrlHConsumed", func(t *testing.T) {
		event := tcell.NewEventKey(tcell.KeyCtrlH, 0, tcell.ModNone)
		result := form.GetInputCapture()(event)

		if result != nil {
			t.Errorf("Expected Ctrl+H to be consumed (return nil), but event was passed through")
		}
	})

	t.Run("OtherKeysNotAffected", func(t *testing.T) {
		// Test that Tab still works
		tabEvent := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		result := form.GetInputCapture()(tabEvent)

		if result == nil {
			t.Errorf("Expected Tab to pass through, but it was consumed")
		}
	})
}

// TestAddFormCursorPreservation validates text preservation after toggle
// T006: Integration test - verifies SetMaskCharacter doesn't clear text
func TestAddFormCursorPreservation(t *testing.T) {
	mockVault := newMockVaultServiceForForms()
	appState := models.NewAppState(mockVault)
	form := NewAddForm(appState)
	passwordField := form.GetFormItem(2).(*tview.InputField)

	// Type "test"
	passwordField.SetText("test")

	// Note: tview doesn't expose SetCursorPosition or GetCursorPosition in InputField
	// This test validates that toggling doesn't clear the text
	t.Run("TextPreservedAfterToggle", func(t *testing.T) {
		originalText := passwordField.GetText()

		// Toggle to visible
		event := tcell.NewEventKey(tcell.KeyCtrlH, 0, tcell.ModNone)
		form.GetInputCapture()(event)

		if passwordField.GetText() != originalText {
			t.Errorf("Expected text '%s' to be preserved, got '%s'", originalText, passwordField.GetText())
		}

		// Toggle back to masked
		form.GetInputCapture()(event)

		if passwordField.GetText() != originalText {
			t.Errorf("Expected text '%s' to be preserved after second toggle, got '%s'", originalText, passwordField.GetText())
		}
	})
}

// TestEditFormPasswordVisibilityToggle verifies EditForm toggle functionality
// T014: Unit test for EditForm password visibility toggle
func TestEditFormPasswordVisibilityToggle(t *testing.T) {
	// Setup
	mockVault := newMockVaultServiceForForms()
	appState := models.NewAppState(mockVault)

	credential := &vault.CredentialMetadata{
		Service:  "test-service",
		Username: "test-user",
		Category: "Test",
	}
	form := NewEditForm(appState, credential)

	// Get password field (index 2)
	passwordField := form.GetFormItem(2).(*tview.InputField)

	// Test initial state - label should show "Password"
	t.Run("InitialStateMasked", func(t *testing.T) {
		if passwordField.GetLabel() != "Password" {
			t.Errorf("Expected initial label 'Password', got '%s'", passwordField.GetLabel())
		}
	})

	// Test toggle to visible - this will FAIL until implementation
	t.Run("ToggleToVisible", func(t *testing.T) {
		event := tcell.NewEventKey(tcell.KeyCtrlH, 0, tcell.ModNone)
		form.GetInputCapture()(event)

		expectedLabel := "Password [VISIBLE]"
		if passwordField.GetLabel() != expectedLabel {
			t.Errorf("Expected label '%s', got '%s'", expectedLabel, passwordField.GetLabel())
		}
	})

	// Test toggle back to masked
	t.Run("ToggleBackToMasked", func(t *testing.T) {
		event := tcell.NewEventKey(tcell.KeyCtrlH, 0, tcell.ModNone)
		form.GetInputCapture()(event)

		if passwordField.GetLabel() != "Password" {
			t.Errorf("Expected label 'Password', got '%s'", passwordField.GetLabel())
		}
	})
}

// TestPasswordDefaultsMasked validates FR-009: passwords default to hidden
// T024: Unit test for password field initialization
func TestPasswordDefaultsMasked(t *testing.T) {
	t.Run("AddFormDefaultsMasked", func(t *testing.T) {
		mockVault := newMockVaultServiceForForms()
		appState := models.NewAppState(mockVault)
		form := NewAddForm(appState)
		passwordField := form.GetFormItem(2).(*tview.InputField)

		// Verify label is not in VISIBLE state
		if passwordField.GetLabel() == "Password [VISIBLE]" {
			t.Error("AddForm: Password should not be visible by default")
		}
	})

	t.Run("EditFormDefaultsMasked", func(t *testing.T) {
		mockVault := newMockVaultServiceForForms()
		appState := models.NewAppState(mockVault)
		credential := &vault.CredentialMetadata{
			Service:  "test",
			Username: "user",
		}
		form := NewEditForm(appState, credential)
		passwordField := form.GetFormItem(2).(*tview.InputField)

		// Verify label is not in VISIBLE state
		if passwordField.GetLabel() == "Password [VISIBLE]" {
			t.Error("EditForm: Password should not be visible by default")
		}
	})
}

// TestEmptyPasswordFieldToggle validates toggle works on empty password field
// T029: Edge case test for empty password
func TestEmptyPasswordFieldToggle(t *testing.T) {
	mockVault := newMockVaultServiceForForms()
	appState := models.NewAppState(mockVault)
	form := NewAddForm(appState)
	passwordField := form.GetFormItem(2).(*tview.InputField)

	// Ensure field is empty
	passwordField.SetText("")

	// Toggle visibility on empty field - should not crash
	t.Run("ToggleEmptyField", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Toggle on empty field caused panic: %v", r)
			}
		}()

		event := tcell.NewEventKey(tcell.KeyCtrlH, 0, tcell.ModNone)
		form.GetInputCapture()(event)

		// Label should still update
		expectedLabel := "Password [VISIBLE]"
		if passwordField.GetLabel() != expectedLabel {
			t.Errorf("Expected label '%s', got '%s'", expectedLabel, passwordField.GetLabel())
		}
	})
}
