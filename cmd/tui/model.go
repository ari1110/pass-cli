package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/internal/vault"
)

// AppState represents the current state of the TUI
type AppState int

const (
	StateUnlocking AppState = iota
	StateList
	StateDetail
	StateAdd
	StateEdit
	StateConfirmDelete
	StateHelp
)

// Model is the main Bubble Tea model for the TUI
type Model struct {
	// Application state
	state        AppState
	vaultService *vault.VaultService
	vaultPath    string

	// Data
	credentials   []vault.CredentialMetadata
	selectedIndex int

	// UI state
	width  int
	height int
	keys   keyMap

	// Error handling
	err       error
	errMsg    string
	unlocking bool
}

// NewModel creates a new TUI model
func NewModel(vaultPath string) (*Model, error) {
	vaultService, err := vault.New(vaultPath)
	if err != nil {
		return nil, err
	}

	return &Model{
		state:        StateUnlocking,
		vaultService: vaultService,
		vaultPath:    vaultPath,
		keys:         DefaultKeyMap(),
		unlocking:    true,
	}, nil
}

// Init initializes the model (Bubble Tea Init method)
func (m Model) Init() tea.Cmd {
	return unlockVaultCmd(m.vaultService)
}

// Update handles messages and updates the model (Bubble Tea Update method)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle quit keys globally
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case vaultUnlockedMsg:
		m.unlocking = false
		m.state = StateList
		// TODO: Load credentials

	case vaultUnlockErrorMsg:
		m.unlocking = false
		m.err = msg.err
		m.errMsg = msg.err.Error()
		return m, tea.Quit

	case credentialsLoadedMsg:
		m.credentials = msg.credentials
	}

	return m, nil
}

// View renders the UI (Bubble Tea View method)
func (m Model) View() string {
	if m.unlocking {
		return "Unlocking vault...\n"
	}

	if m.err != nil {
		return "Error: " + m.errMsg + "\n"
	}

	switch m.state {
	case StateList:
		return "Credential list view - coming soon!\n\nPress q to quit"
	default:
		return "TUI - coming soon!\n"
	}
}
