package tui

import (
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

// PanelFocus represents which panel has focus in dashboard mode
type PanelFocus int

const (
	FocusSidebar PanelFocus = iota
	FocusMain
	FocusMetadata
	FocusCommandBar
)

// Model is the main Bubble Tea model for the TUI
type Model struct {
	// Application state
	state        AppState
	vaultService *vault.VaultService
	vaultPath    string

	// Data
	credentials []vault.CredentialMetadata

	// Views
	listView    *views.ListView
	detailView  *views.DetailView
	addForm     *views.AddFormView
	editForm    *views.EditFormView
	confirmView *views.ConfirmView
	helpView    *views.HelpView

	// Components
	statusBar *components.StatusBar

	// Services
	keychainService *keychain.KeychainService

	// State for confirmation dialogs
	previousState AppState

	// UI state
	width  int
	height int

	// Error handling
	err       error
	errMsg    string
	unlocking bool

	// === DASHBOARD COMPONENTS (new) ===
	layoutManager *components.LayoutManager
	sidebar       *components.SidebarPanel
	metadataPanel *components.MetadataPanel
	processPanel  *components.ProcessPanel
	commandBar    *components.CommandBar
	breadcrumb    *components.Breadcrumb

	// Dashboard state
	panelFocus      PanelFocus
	sidebarVisible  bool
	metadataVisible bool
	processVisible  bool
	commandBarOpen  bool

	// Category state
	categories      []components.Category
	currentCategory string
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

	// Initialize dashboard components
	layoutManager := components.NewLayoutManager()
	sidebar := components.NewSidebarPanel([]vault.CredentialMetadata{})
	metadataPanel := components.NewMetadataPanel()
	processPanel := components.NewProcessPanel()
	commandBar := components.NewCommandBar()
	breadcrumb := components.NewBreadcrumb()

	return &Model{
		state:           StateUnlocking,
		vaultService:    vaultService,
		vaultPath:       vaultPath,
		keychainService: keychainService,
		statusBar:       statusBar,
		unlocking:       true,

		// Initialize dashboard components
		layoutManager: layoutManager,
		sidebar:       sidebar,
		metadataPanel: metadataPanel,
		processPanel:  processPanel,
		commandBar:    commandBar,
		breadcrumb:    breadcrumb,

		// Initialize dashboard state (sidebar visible by default)
		panelFocus:      FocusSidebar,
		sidebarVisible:  true,
		metadataVisible: false,
		processVisible:  false,
		commandBarOpen:  false,
		categories:      []components.Category{},
		currentCategory: "",
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
		// Check if any view has input focus - if so, let it handle ALL keys first
		hasInputFocus := false

		switch m.state {
		case StateList:
			hasInputFocus = m.listView != nil && m.listView.IsSearchFocused()
		case StateAdd:
			hasInputFocus = m.addForm != nil // Forms always have input focus
		case StateEdit:
			hasInputFocus = m.editForm != nil // Forms always have input focus
		case StateConfirmDelete, StateConfirmDiscard:
			hasInputFocus = m.confirmView != nil && m.confirmView.IsTypedConfirmation()
		}

		// If view has input focus, skip global key handling (except Ctrl+C for emergency quit)
		if !hasInputFocus {
			// Handle help overlay - allow scrolling with arrow keys
			if m.state == StateHelp {
				// Only close on specific keys, not arrow keys
				switch msg.String() {
				case "up", "down", "pgup", "pgdown", "home", "end":
					// Let help view handle scrolling
				case "q", "esc", "?", "f1", "enter", " ":
					// These keys close help
					m.state = m.previousState
					m.helpView = nil
					return m, nil
				default:
					// Any other key closes help
					m.state = m.previousState
					m.helpView = nil
					return m, nil
				}
			}

			// Handle help key (? or F1) from any view (except when typing)
			if msg.String() == "?" || msg.String() == "f1" {
				m.previousState = m.state
				m.state = StateHelp
				m.helpView = views.NewHelpView()
				m.helpView.SetSize(m.width, m.height)
				return m, nil
			}

			// Handle quit keys (but only when not typing)
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		} else {
			// Even with input focus, allow Ctrl+C for emergency quit
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Calculate layout with panel visibility
		m.recalculateLayout()

		// Update status bar
		if m.statusBar != nil {
			m.statusBar.SetSize(m.width)
		}

		// Views will be sized according to the layout's main panel dimensions
		// This happens in recalculateLayout() for views in List/Detail states
		// For forms/help/confirm, use full dimensions
		if m.addForm != nil {
			m.addForm.SetSize(m.width, m.height)
		}
		if m.editForm != nil {
			m.editForm.SetSize(m.width, m.height)
		}
		if m.confirmView != nil {
			m.confirmView.SetSize(m.width, m.height)
		}
		if m.helpView != nil {
			m.helpView.SetSize(m.width, m.height)
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
		// Update sidebar with categorized credentials
		if m.sidebar != nil {
			m.sidebar.UpdateCredentials(msg.credentials)
			m.categories = components.CategorizeCredentials(msg.credentials)
		}
		// Return to list from add/edit/delete states
		if m.state == StateAdd || m.state == StateEdit || m.state == StateConfirmDelete {
			m.state = StateList
			m.addForm = nil
			m.editForm = nil
			m.confirmView = nil
			m.detailView = nil // Clear detail view since credential may have been deleted
		}

	case credentialLoadedMsg:
		m.state = StateDetail
		m.detailView = views.NewDetailView(msg.credential)
		m.detailView.SetSize(m.width, m.height)
		// Also update metadata panel with credential
		if m.metadataPanel != nil {
			m.metadataPanel.SetCredential(msg.credential)
		}
	}

	// Check if any view has input focus for panel toggle keys
	hasInputFocusForToggle := false
	switch m.state {
	case StateList:
		hasInputFocusForToggle = m.listView != nil && m.listView.IsSearchFocused()
	case StateAdd, StateEdit:
		hasInputFocusForToggle = true // Forms always have input focus
	case StateConfirmDelete, StateConfirmDiscard:
		hasInputFocusForToggle = m.confirmView != nil && m.confirmView.IsTypedConfirmation()
	}

	// Handle panel toggle keys (only when no input focus and not in forms/confirmations)
	if keyMsg, ok := msg.(tea.KeyMsg); ok && !hasInputFocusForToggle && !m.commandBarOpen {
		switch keyMsg.String() {
		case "s":
			// Toggle sidebar (only in List/Detail states)
			if m.state == StateList || m.state == StateDetail {
				m.sidebarVisible = !m.sidebarVisible
				m.recalculateLayout()
				return m, nil
			}
		case "m":
			// Toggle metadata panel (only in Detail state, avoid conflict with mask in DetailView)
			if m.state == StateDetail && m.detailView != nil {
				// Let detail view handle 'm' for password masking if it has focus
				if m.panelFocus != FocusMetadata {
					// If metadata panel is focused, let it handle the key
					// Otherwise, toggle metadata panel visibility
					if m.panelFocus == FocusMain {
						// Main panel has focus, toggle metadata
						m.metadataVisible = !m.metadataVisible
						m.recalculateLayout()
						return m, nil
					}
				}
			}
		case "p":
			// Toggle process panel
			if m.state == StateList || m.state == StateDetail {
				m.processVisible = !m.processVisible
				m.recalculateLayout()
				return m, nil
			}
		case "f":
			// Toggle all footer panels (process + command bar if open)
			if m.state == StateList || m.state == StateDetail {
				m.processVisible = !m.processVisible
				m.recalculateLayout()
				return m, nil
			}
		case "tab":
			// Switch panel focus (only in List/Detail states)
			if m.state == StateList || m.state == StateDetail {
				m.panelFocus = m.nextPanelFocus()
				m.updatePanelFocus()
				return m, nil
			}
		case "shift+tab":
			// Previous panel focus
			if m.state == StateList || m.state == StateDetail {
				m.panelFocus = m.previousPanelFocus()
				m.updatePanelFocus()
				return m, nil
			}
		case ":":
			// Open command bar (only in List/Detail states)
			if m.state == StateList || m.state == StateDetail {
				m.commandBarOpen = true
				m.panelFocus = FocusCommandBar
				m.updatePanelFocus()
				return m, m.commandBar.Focus()
			}
		}
	}

	// Handle command bar when open
	if m.commandBarOpen {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc":
				// Close command bar
				m.commandBarOpen = false
				m.commandBar.Blur()
				m.panelFocus = FocusMain
				m.updatePanelFocus()
				return m, nil
			case "enter":
				// Execute command
				cmd, err := m.commandBar.GetCommand()
				if err != nil {
					return m, nil
				}
				if cmd != nil {
					return m.executeCommand(cmd)
				}
				// Close command bar after execution
				m.commandBarOpen = false
				m.commandBar.Blur()
				m.panelFocus = FocusMain
				m.updatePanelFocus()
				return m, nil
			}
		}
		// Update command bar
		var cmd tea.Cmd
		m.commandBar, cmd = m.commandBar.Update(msg)
		return m, cmd
	}

	// Update sidebar if it has focus (in List or Detail states)
	if (m.state == StateList || m.state == StateDetail) && m.panelFocus == FocusSidebar && m.sidebar != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "enter":
				// Get selected credential from sidebar
				selected := m.sidebar.GetSelectedCredential()
				if selected != nil {
					// Load full credential details
					m.panelFocus = FocusMain // Switch focus to main panel
					m.updatePanelFocus()
					// Update breadcrumb with category path
					category := m.sidebar.GetSelectedCategory()
					if category != "" {
						m.breadcrumb.SetPath([]string{"Home", category, selected.Service})
					}
					return m, loadCredentialDetailsCmd(m.vaultService, selected.Service)
				}
			}
		}
		// Update sidebar with message
		var cmd tea.Cmd
		m.sidebar, cmd = m.sidebar.Update(msg)
		return m, cmd
	}

	// Update metadata panel if it has focus (in Detail state)
	if m.state == StateDetail && m.panelFocus == FocusMetadata && m.metadataPanel != nil {
		var cmd tea.Cmd
		m.metadataPanel, cmd = m.metadataPanel.Update(msg)
		return m, cmd
	}

	// Update active view
	if m.state == StateList && m.listView != nil {
		// Check for special keys (but only if search is not focused)
		if keyMsg, ok := msg.(tea.KeyMsg); ok && !m.listView.IsSearchFocused() {
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
			case "ctrl+g":
				// Generate password (Ctrl+G to avoid conflict with typing 'g')
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
			case "ctrl+g":
				// Generate password (Ctrl+G to avoid conflict with typing 'g')
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

	if m.state == StateHelp && m.helpView != nil {
		var cmd tea.Cmd
		m.helpView, cmd = m.helpView.Update(msg)
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

	// Check if we should render dashboard layout (List or Detail states only)
	if m.state == StateList || m.state == StateDetail {
		return m.renderDashboardView()
	}

	// For other states (forms, help, confirmations), render full screen
	var mainContent string
	switch m.state {
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
	case StateHelp:
		if m.helpView != nil {
			// Help overlay doesn't show status bar
			return m.helpView.View()
		}
		mainContent = "Loading help...\n"
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
		shortcuts := "/: search | a: add | :: cmd | ?: help | q: quit"
		// Add panel shortcuts
		panelShortcuts := m.getPanelShortcuts()
		if panelShortcuts != "" {
			shortcuts = panelShortcuts + " | " + shortcuts
		}
		m.statusBar.SetShortcuts(shortcuts)
	case StateDetail:
		m.statusBar.SetCurrentView("Detail")
		shortcuts := "m: toggle | c: copy | e: edit | d: delete | esc: back | q: quit"
		// Add panel shortcuts
		panelShortcuts := m.getPanelShortcuts()
		if panelShortcuts != "" {
			shortcuts = panelShortcuts + " | " + shortcuts
		}
		m.statusBar.SetShortcuts(shortcuts)
	case StateAdd:
		m.statusBar.SetCurrentView("Add")
		m.statusBar.SetShortcuts("tab: next | ctrl+g: generate | ctrl+s: save | esc: cancel")
	case StateEdit:
		m.statusBar.SetCurrentView("Edit")
		m.statusBar.SetShortcuts("tab: next | ctrl+g: generate | ctrl+s: save | esc: cancel")
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
