package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewCommandBar(t *testing.T) {
	cb := NewCommandBar()

	if cb == nil {
		t.Fatal("NewCommandBar returned nil")
	}

	if cb.width != 0 {
		t.Errorf("Expected initial width 0, got %d", cb.width)
	}

	if cb.height != 0 {
		t.Errorf("Expected initial height 0, got %d", cb.height)
	}

	if len(cb.history) != 0 {
		t.Errorf("Expected empty history, got %d items", len(cb.history))
	}

	if cb.historyIdx != -1 {
		t.Errorf("Expected historyIdx -1, got %d", cb.historyIdx)
	}
}

func TestCommandBar_SetSize(t *testing.T) {
	cb := NewCommandBar()

	cb.SetSize(80, 3)

	if cb.width != 80 {
		t.Errorf("Expected width 80, got %d", cb.width)
	}

	if cb.height != 3 {
		t.Errorf("Expected height 3, got %d", cb.height)
	}
}

func TestCommandBar_Focus(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)

	cb.Focus()

	// Check that input is focused by checking if it has focus
	if !cb.input.Focused() {
		t.Error("Input should be focused after Focus()")
	}
}

func TestCommandBar_Blur(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)

	cb.Focus()
	cb.Blur()

	// Check that input is blurred
	if cb.input.Focused() {
		t.Error("Input should not be focused after Blur()")
	}
}

func TestCommandBar_GetCommand_Valid(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Simulate typing a command
	cb.input.SetValue(":help")

	cmd, err := cb.GetCommand()

	if err != nil {
		t.Errorf("GetCommand returned error: %v", err)
	}

	if cmd == nil {
		t.Fatal("GetCommand returned nil command")
	}

	if cmd.Name != "help" {
		t.Errorf("Expected command name 'help', got '%s'", cmd.Name)
	}

	if len(cmd.Args) != 0 {
		t.Errorf("Expected 0 args, got %d", len(cmd.Args))
	}
}

func TestCommandBar_GetCommand_WithArgs(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Simulate typing a command with args
	cb.input.SetValue(":add github myuser")

	cmd, err := cb.GetCommand()

	if err != nil {
		t.Errorf("GetCommand returned error: %v", err)
	}

	if cmd == nil {
		t.Fatal("GetCommand returned nil command")
	}

	if cmd.Name != "add" {
		t.Errorf("Expected command name 'add', got '%s'", cmd.Name)
	}

	if len(cmd.Args) != 2 {
		t.Fatalf("Expected 2 args, got %d", len(cmd.Args))
	}

	if cmd.Args[0] != "github" {
		t.Errorf("Expected first arg 'github', got '%s'", cmd.Args[0])
	}

	if cmd.Args[1] != "myuser" {
		t.Errorf("Expected second arg 'myuser', got '%s'", cmd.Args[1])
	}
}

func TestCommandBar_GetCommand_Empty(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Simulate typing just colon
	cb.input.SetValue(":")

	cmd, err := cb.GetCommand()

	if err != nil {
		t.Errorf("GetCommand returned error: %v", err)
	}

	if cmd != nil {
		t.Error("Expected nil command for empty input")
	}

	if cb.error == "" {
		t.Error("Expected error message for empty command")
	}
}

func TestCommandBar_GetCommand_NoColon(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Simulate typing without colon
	cb.input.SetValue("help")

	cmd, err := cb.GetCommand()

	if err != nil {
		t.Errorf("GetCommand returned error: %v", err)
	}

	if cmd != nil {
		t.Error("Expected nil command for input without colon")
	}

	if cb.error == "" {
		t.Error("Expected error message for command without colon")
	}
}

func TestCommandBar_SetError(t *testing.T) {
	cb := NewCommandBar()

	cb.SetError("Test error message")

	if cb.error != "Test error message" {
		t.Errorf("Expected error 'Test error message', got '%s'", cb.error)
	}
}

func TestCommandBar_ClearError(t *testing.T) {
	cb := NewCommandBar()

	cb.SetError("Test error")
	cb.ClearError()

	if cb.error != "" {
		t.Errorf("Expected empty error, got '%s'", cb.error)
	}
}

func TestCommandBar_GetHistory(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Add some commands to history
	cb.input.SetValue(":help")
	cb.GetCommand()

	cb.input.SetValue(":add test user")
	cb.GetCommand()

	history := cb.GetHistory()

	if len(history) != 2 {
		t.Errorf("Expected 2 history items, got %d", len(history))
	}

	if history[0] != "help" {
		t.Errorf("Expected first history item 'help', got '%s'", history[0])
	}

	if history[1] != "add test user" {
		t.Errorf("Expected second history item 'add test user', got '%s'", history[1])
	}
}

func TestCommandBar_View(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)

	output := cb.View()

	if output == "" {
		t.Error("View() should return non-empty output")
	}

	// Should contain the prompt
	if !strings.Contains(output, ":") {
		t.Error("View output should contain command prompt")
	}
}

func TestCommandBar_View_WithError(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.SetError("Invalid command")

	output := cb.View()

	if output == "" {
		t.Error("View() should return non-empty output")
	}

	// Should contain the error message
	if !strings.Contains(output, "Invalid command") {
		t.Error("View output should contain error message")
	}
}

func TestCommandBar_Update(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Test that Update processes messages
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")}
	updatedCb, _ := cb.Update(msg)

	if updatedCb == nil {
		t.Error("Update should return updated CommandBar")
	}
}

func TestCommandBar_GetCommand_WhitespaceHandling(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Test with extra whitespace
	cb.input.SetValue(":  help  ")

	cmd, err := cb.GetCommand()

	if err != nil {
		t.Errorf("GetCommand returned error: %v", err)
	}

	if cmd == nil {
		t.Fatal("GetCommand returned nil command")
	}

	if cmd.Name != "help" {
		t.Errorf("Expected command name 'help', got '%s'", cmd.Name)
	}
}

func TestCommandBar_GetCommand_MultipleArgs(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Test with many arguments
	cb.input.SetValue(":command arg1 arg2 arg3 arg4")

	cmd, err := cb.GetCommand()

	if err != nil {
		t.Errorf("GetCommand returned error: %v", err)
	}

	if cmd == nil {
		t.Fatal("GetCommand returned nil command")
	}

	if cmd.Name != "command" {
		t.Errorf("Expected command name 'command', got '%s'", cmd.Name)
	}

	if len(cmd.Args) != 4 {
		t.Errorf("Expected 4 args, got %d", len(cmd.Args))
	}

	expectedArgs := []string{"arg1", "arg2", "arg3", "arg4"}
	for i, expected := range expectedArgs {
		if cmd.Args[i] != expected {
			t.Errorf("Expected arg %d to be '%s', got '%s'", i, expected, cmd.Args[i])
		}
	}
}

func TestCommand_Structure(t *testing.T) {
	// Test Command struct
	cmd := Command{
		Name: "test",
		Args: []string{"arg1", "arg2"},
	}

	if cmd.Name != "test" {
		t.Errorf("Expected name 'test', got '%s'", cmd.Name)
	}

	if len(cmd.Args) != 2 {
		t.Errorf("Expected 2 args, got %d", len(cmd.Args))
	}
}

func TestCommandBar_HistoryNavigation(t *testing.T) {
	cb := NewCommandBar()
	cb.SetSize(80, 3)
	cb.Focus()

	// Add commands to history
	cb.input.SetValue(":help")
	cb.GetCommand()

	cb.input.SetValue(":add test")
	cb.GetCommand()

	// Check history was saved
	history := cb.GetHistory()
	if len(history) != 2 {
		t.Errorf("Expected 2 history items, got %d", len(history))
	}

	// Check historyIdx was reset
	if cb.historyIdx != -1 {
		t.Errorf("Expected historyIdx -1 after command, got %d", cb.historyIdx)
	}
}
