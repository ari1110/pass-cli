package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/cmd/tui/components"
	"pass-cli/cmd/tui/views"
	"pass-cli/internal/keychain"
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
	StateConfirmDiscard
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
	listView    *views.ListView
	detailView  *views.DetailView
	addForm     *views.AddFormView
	editForm    *views.EditFormView
	confirmView *views.ConfirmView

	// Components
	statusBar *components.StatusBar

	// Services
	keychainService *keychain.KeychainService

	// State for confirmation dialogs
	previousState AppState

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

	// Initialize keychain service
	keychainService := keychain.New()

	// Initialize status bar
	statusBar := components.NewStatusBar(
		keychainService.IsAvailable(),
		0,
		"Unlocking",
	)

	return &Model{
		state:           StateUnlocking,
		vaultService:    vaultService,
		vaultPath:       vaultPath,
		keychainService: keychainService,
		statusBar:       statusBar,
		keys:            DefaultKeyMap(),
		unlocking:       true,
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
		if m.statusBar != nil {
			m.statusBar.SetSize(m.width)
		}
		if m.listView != nil {
			m.listView.SetSize(m.width, m.height)
		}
		if m.detailView != nil {
			m.detailView.SetSize(m.width, m.height)
		}
		if m.addForm != nil {
			m.addForm.SetSize(m.width, m.height)
		}
		if m.editForm != nil {
			m.editForm.SetSize(m.width, m.height)
		}
		if m.confirmView != nil {
			m.confirmView.SetSize(m.width, m.height)
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
		// Update status bar with credential count
		if m.statusBar != nil {
			m.statusBar.SetCredentialCount(len(msg.credentials))
		}
		// If we were in add/edit form, return to list
		if m.state == StateAdd || m.state == StateEdit {
			m.state = StateList
			m.addForm = nil
			m.editForm = nil
		}

	case credentialLoadedMsg:
		m.state = StateDetail
		m.detailView = views.NewDetailView(msg.credential)
		m.detailView.SetSize(m.width, m.height)
	}

	// Update active view
	if m.state == StateList && m.listView != nil {
		// Check for special keys
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":
				selected := m.listView.SelectedCredential()
				if selected != nil {
					// Load full credential details
					return m, loadCredentialDetailsCmd(m.vaultService, selected.Service)
				}
			case "a":
				// Open add form
				m.state = StateAdd
				m.addForm = views.NewAddFormView()
				m.addForm.SetSize(m.width, m.height)
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.listView, cmd = m.listView.Update(msg)
		return m, cmd
	}

	if m.state == StateDetail && m.detailView != nil {
		// Check for special keys
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				// Go back to list
				m.state = StateList
				m.detailView = nil
				return m, nil
			case "e":
				// Open edit form with current credential
				cred := m.detailView.GetCredential()
				if cred != nil {
					m.state = StateEdit
					m.editForm = views.NewEditFormView(cred)
					m.editForm.SetSize(m.width, m.height)
					return m, nil
				}
			case "d":
				// Open delete confirmation
				cred := m.detailView.GetCredential()
				if cred != nil {
					m.previousState = StateDetail
					m.state = StateConfirmDelete
					m.confirmView = views.NewDeleteConfirmView(cred)
					m.confirmView.SetSize(m.width, m.height)
					return m, nil
				}
			}
		}

		var cmd tea.Cmd
		m.detailView, cmd = m.detailView.Update(msg)
		return m, cmd
	}

	if m.state == StateAdd && m.addForm != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				// Check if there are changes
				if m.addForm.HasChanges() {
					// Show confirmation dialog
					m.previousState = StateAdd
					m.state = StateConfirmDiscard
					m.confirmView = views.NewConfirmView("Discard new credential?")
					m.confirmView.SetSize(m.width, m.height)
					return m, nil
				}
				// No changes, just go back
				m.state = StateList
				m.addForm = nil
				return m, nil
			case "ctrl+s":
				// Save the credential
				return m, m.saveNewCredential()
			case "g":
				// Generate password
				if password, err := generatePassword(20); err == nil {
					m.addForm.SetPassword(password)
					m.addForm.SetNotification("Password generated (20 characters)")
				}
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.addForm, cmd = m.addForm.Update(msg)
		return m, cmd
	}

	if m.state == StateEdit && m.editForm != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				// Check if there are changes
				if m.editForm.HasChanges() {
					// Show confirmation dialog
					m.previousState = StateEdit
					m.state = StateConfirmDiscard
					m.confirmView = views.NewConfirmView("Discard changes?")
					m.confirmView.SetSize(m.width, m.height)
					return m, nil
				}
				// No changes, go back to detail view
				m.state = StateDetail
				m.editForm = nil
				return m, nil
			case "ctrl+s":
				// Save the credential
				return m, m.updateCredential()
			case "g":
				// Generate password
				if password, err := generatePassword(20); err == nil {
					m.editForm.SetPassword(password)
					m.editForm.SetNotification("Password generated (20 characters)")
				}
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.editForm, cmd = m.editForm.Update(msg)
		return m, cmd
	}

	if m.state == StateConfirmDiscard && m.confirmView != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "y":
				// User confirmed discard
				if m.previousState == StateEdit {
					m.state = StateDetail
					m.editForm = nil
				} else {
					m.state = StateList
					m.addForm = nil
				}
				m.confirmView = nil
				return m, nil
			case "n", "esc":
				// User cancelled, go back to form
				m.state = m.previousState
				m.confirmView = nil
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.confirmView, cmd = m.confirmView.Update(msg)
		return m, cmd
	}

	if m.state == StateConfirmDelete && m.confirmView != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			// Check if this is a typed confirmation
			if m.confirmView.IsTypedConfirmation() {
				switch keyMsg.String() {
				case "enter":
					// Validate typed service name
					cred := m.confirmView.GetCredential()
					typedValue := m.confirmView.GetTypedValue()
					if typedValue != cred.Service {
						m.confirmView.SetError("Service name doesn't match")
						return m, nil
					}
					// Confirmed, delete the credential
					return m, m.deleteCredential(cred.Service)
				case "esc":
					// User cancelled, go back
					m.state = m.previousState
					m.confirmView = nil
					return m, nil
				}
			} else {
				// Simple y/n confirmation
				switch keyMsg.String() {
				case "y":
					// Confirmed, delete the credential
					cred := m.confirmView.GetCredential()
					return m, m.deleteCredential(cred.Service)
				case "n", "esc":
					// User cancelled, go back
					m.state = m.previousState
					m.confirmView = nil
					return m, nil
				}
			}
		}

		var cmd tea.Cmd
		m.confirmView, cmd = m.confirmView.Update(msg)
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

	// Update status bar based on current state
	m.updateStatusBar()

	// Render main content
	var mainContent string
	switch m.state {
	case StateList:
		if m.listView != nil {
			mainContent = m.listView.View()
		} else {
			mainContent = "Loading credentials...\n"
		}
	case StateDetail:
		if m.detailView != nil {
			mainContent = m.detailView.View()
		} else {
			mainContent = "Loading credential...\n"
		}
	case StateAdd:
		if m.addForm != nil {
			mainContent = m.addForm.View()
		} else {
			mainContent = "Loading form...\n"
		}
	case StateEdit:
		if m.editForm != nil {
			mainContent = m.editForm.View()
		} else {
			mainContent = "Loading form...\n"
		}
	case StateConfirmDiscard, StateConfirmDelete:
		if m.confirmView != nil {
			mainContent = m.confirmView.View()
		} else {
			mainContent = "Loading...\n"
		}
	default:
		mainContent = "TUI - coming soon!\n"
	}

	// Append status bar
	if m.statusBar != nil {
		return mainContent + "\n" + m.statusBar.Render()
	}

	return mainContent
}

// updateStatusBar updates the status bar with current view information
func (m Model) updateStatusBar() {
	if m.statusBar == nil {
		return
	}

	// Update current view name and shortcuts
	switch m.state {
	case StateList:
		m.statusBar.SetCurrentView("List")
		m.statusBar.SetShortcuts("/: search | a: add | ?: help | q: quit")
	case StateDetail:
		m.statusBar.SetCurrentView("Detail")
		m.statusBar.SetShortcuts("m: toggle | c: copy | e: edit | d: delete | esc: back | q: quit")
	case StateAdd:
		m.statusBar.SetCurrentView("Add")
		m.statusBar.SetShortcuts("tab: next | g: generate | ctrl+s: save | esc: cancel")
	case StateEdit:
		m.statusBar.SetCurrentView("Edit")
		m.statusBar.SetShortcuts("tab: next | g: generate | ctrl+s: save | esc: cancel")
	case StateConfirmDelete, StateConfirmDiscard:
		m.statusBar.SetCurrentView("Confirm")
		m.statusBar.SetShortcuts("y/n or enter/esc")
	default:
		m.statusBar.SetCurrentView("TUI")
		m.statusBar.SetShortcuts("?: help | q: quit")
	}
}

// saveNewCredential validates and saves a new credential
func (m *Model) saveNewCredential() tea.Cmd {
	service, username, password, notes := m.addForm.GetValues()

	// Validate service name
	if service == "" {
		m.addForm.SetError("Service name is required")
		return nil
	}

	// Check if service already exists
	for _, cred := range m.credentials {
		if cred.Service == service {
			m.addForm.SetError("Service already exists")
			return nil
		}
	}

	// Generate password if empty
	if password == "" {
		var err error
		password, err = generatePassword(20)
		if err != nil {
			m.addForm.SetError("Failed to generate password")
			return nil
		}
	}

	// Add credential
	return func() tea.Msg {
		if err := m.vaultService.AddCredential(service, username, password, notes); err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		// Reload credentials and return to list
		credentials, err := m.vaultService.ListCredentialsWithMetadata()
		if err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		return credentialsLoadedMsg{credentials: credentials}
	}
}

// updateCredential validates and updates an existing credential
func (m *Model) updateCredential() tea.Cmd {
	service, username, password, notes := m.editForm.GetValues()

	// Service name should not be empty (though it's read-only)
	if service == "" {
		m.editForm.SetError("Service name is required")
		return nil
	}

	// Update credential
	return func() tea.Msg {
		if err := m.vaultService.UpdateCredential(service, username, password, notes); err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		// Reload credentials
		credentials, err := m.vaultService.ListCredentialsWithMetadata()
		if err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		// Also reload the current credential details to show in detail view
		cred, err := m.vaultService.GetCredential(service, false)
		if err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		// Update credentials list
		m.credentials = credentials

		// Return to detail view with updated credential
		return credentialLoadedMsg{credential: cred}
	}
}

// deleteCredential deletes a credential from the vault
func (m *Model) deleteCredential(service string) tea.Cmd {
	return func() tea.Msg {
		// Delete the credential
		if err := m.vaultService.DeleteCredential(service); err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		// Reload credentials and return to list
		credentials, err := m.vaultService.ListCredentialsWithMetadata()
		if err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		return credentialsLoadedMsg{credentials: credentials}
	}
}
