package tui

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/internal/vault"
)

// setupTestVault creates a test vault and returns cleanup function
func setupTestVault(t *testing.T) (*vault.VaultService, string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "tui-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	vaultPath := filepath.Join(tempDir, "test.vault")
	v, err := vault.New(vaultPath)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		t.Fatalf("Failed to create vault service: %v", err)
	}

	// Initialize and unlock vault
	password := "test-password-12345"
	if err := v.Initialize(password, false); err != nil {
		_ = os.RemoveAll(tempDir)
		t.Fatalf("Failed to initialize vault: %v", err)
	}

	if err := v.Unlock(password); err != nil {
		_ = os.RemoveAll(tempDir)
		t.Fatalf("Failed to unlock vault: %v", err)
	}

	cleanup := func() {
		v.Lock()
		_ = os.RemoveAll(tempDir)
	}

	return v, vaultPath, cleanup
}

func TestNewModel(t *testing.T) {
	_, vaultPath, cleanup := setupTestVault(t)
	defer cleanup()

	model, err := NewModel(vaultPath)
	if err != nil {
		t.Fatalf("NewModel() failed: %v", err)
	}

	if model == nil {
		t.Fatal("NewModel() returned nil")
	}

	if model.state != StateUnlocking {
		t.Errorf("Expected initial state to be StateUnlocking, got %v", model.state)
	}

	if model.vaultService == nil {
		t.Error("VaultService is nil")
	}

	if model.statusBar == nil {
		t.Error("StatusBar is nil")
	}

	if model.keychainService == nil {
		t.Error("KeychainService is nil")
	}
}

func TestModelStateTransitions(t *testing.T) {
	_, vaultPath, cleanup := setupTestVault(t)
	defer cleanup()

	model, err := NewModel(vaultPath)
	if err != nil {
		t.Fatalf("NewModel() failed: %v", err)
	}

	tests := []struct {
		name          string
		initialState  AppState
		msg           tea.Msg
		expectedState AppState
	}{
		{
			name:          "Help key opens help from list",
			initialState:  StateList,
			msg:           tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}},
			expectedState: StateHelp,
		},
		{
			name:          "F1 opens help from list",
			initialState:  StateList,
			msg:           tea.KeyMsg{Type: tea.KeyF1},
			expectedState: StateHelp,
		},
		{
			name:          "Any key closes help",
			initialState:  StateHelp,
			msg:           tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			expectedState: StateList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model.state = tt.initialState
			if tt.initialState == StateHelp {
				model.previousState = StateList
			}

			updatedModel, _ := model.Update(tt.msg)
			m := updatedModel.(Model)

			if m.state != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, m.state)
			}
		})
	}
}

func TestModelQuitHandling(t *testing.T) {
	_, vaultPath, cleanup := setupTestVault(t)
	defer cleanup()

	model, err := NewModel(vaultPath)
	if err != nil {
		t.Fatalf("NewModel() failed: %v", err)
	}

	model.state = StateList

	// Test 'q' key
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(quitMsg)

	if cmd == nil {
		t.Error("Expected quit command, got nil")
	}

	// Verify it's a quit command (we can't directly compare commands, but we can check it's not nil)
	// The actual quit behavior would be tested in integration tests
}

func TestModelWindowResize(t *testing.T) {
	_, vaultPath, cleanup := setupTestVault(t)
	defer cleanup()

	model, err := NewModel(vaultPath)
	if err != nil {
		t.Fatalf("NewModel() failed: %v", err)
	}

	resizeMsg := tea.WindowSizeMsg{Width: 100, Height: 40}
	updatedModel, _ := model.Update(resizeMsg)
	m := updatedModel.(Model)

	if m.width != 100 {
		t.Errorf("Expected width 100, got %d", m.width)
	}

	if m.height != 40 {
		t.Errorf("Expected height 40, got %d", m.height)
	}
}

func TestUpdateStatusBar(t *testing.T) {
	_, vaultPath, cleanup := setupTestVault(t)
	defer cleanup()

	model, err := NewModel(vaultPath)
	if err != nil {
		t.Fatalf("NewModel() failed: %v", err)
	}

	tests := []struct {
		state        AppState
		expectedView string
	}{
		{StateList, "List"},
		{StateDetail, "Detail"},
		{StateAdd, "Add"},
		{StateEdit, "Edit"},
		{StateConfirmDelete, "Confirm"},
		{StateConfirmDiscard, "Confirm"},
	}

	for _, tt := range tests {
		t.Run(tt.expectedView, func(t *testing.T) {
			model.state = tt.state
			model.updateStatusBar()

			// We can't easily inspect the current view without adding a getter,
			// but we can at least verify the method doesn't panic
		})
	}
}

func TestModelView(t *testing.T) {
	_, vaultPath, cleanup := setupTestVault(t)
	defer cleanup()

	model, err := NewModel(vaultPath)
	if err != nil {
		t.Fatalf("NewModel() failed: %v", err)
	}

	// Test unlocking view
	model.unlocking = true
	view := model.View()
	if view != "Unlocking vault...\n" {
		t.Errorf("Expected unlocking message, got: %s", view)
	}

	// Test error view
	model.unlocking = false
	model.err = os.ErrNotExist
	model.errMsg = "test error"
	view = model.View()
	if view != "Error: test error\n" {
		t.Errorf("Expected error message, got: %s", view)
	}
}
