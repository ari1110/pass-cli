package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/cmd"
)

// Run starts the TUI interface
func Run() error {
	vaultPath := cmd.GetVaultPath()

	// Check if vault exists
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault not found at %s\nRun 'pass-cli init' to create a vault first", vaultPath)
	}

	// Create model
	model, err := NewModel(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to create TUI model: %w", err)
	}

	// Create Bubble Tea program
	p := tea.NewProgram(model)

	// Run the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
