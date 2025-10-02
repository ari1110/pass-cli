package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/cmd/tui/views"
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

	// Views
	listView   *views.ListView
	detailView *views.DetailView

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
		if m.listView != nil {
			m.listView.SetSize(m.width, m.height)
		}
		if m.detailView != nil {
			m.detailView.SetSize(m.width, m.height)
		}

	case vaultUnlockedMsg:
		m.unlocking = false
		m.state = StateList
		return m, loadCredentialsCmd(m.vaultService)

	case vaultUnlockErrorMsg:
		m.unlocking = false
		m.err = msg.err
		m.errMsg = msg.err.Error()
		return m, tea.Quit

	case credentialsLoadedMsg:
		m.credentials = msg.credentials
		m.listView = views.NewListView(msg.credentials)
		m.listView.SetSize(m.width, m.height)

	case credentialLoadedMsg:
		m.state = StateDetail
		m.detailView = views.NewDetailView(msg.credential)
		m.detailView.SetSize(m.width, m.height)
	}

	// Update active view
	if m.state == StateList && m.listView != nil {
		// Check for Enter key to navigate to detail
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			selected := m.listView.SelectedCredential()
			if selected != nil {
				// Load full credential details
				return m, loadCredentialDetailsCmd(m.vaultService, selected.Service)
			}
		}

		var cmd tea.Cmd
		m.listView, cmd = m.listView.Update(msg)
		return m, cmd
	}

	if m.state == StateDetail && m.detailView != nil {
		// Check for Escape to go back to list
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
			m.state = StateList
			m.detailView = nil
			return m, nil
		}

		var cmd tea.Cmd
		m.detailView, cmd = m.detailView.Update(msg)
		return m, cmd
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
		if m.listView != nil {
			return m.listView.View()
		}
		return "Loading credentials...\n"
	case StateDetail:
		if m.detailView != nil {
			return m.detailView.View()
		}
		return "Loading credential...\n"
	default:
		return "TUI - coming soon!\n"
	}
}
