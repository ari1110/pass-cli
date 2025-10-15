package events

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"pass-cli/cmd/tui/components"
	"pass-cli/cmd/tui/layout"
	"pass-cli/cmd/tui/models"
	"pass-cli/internal/config"
)

// EventHandler manages global keyboard shortcuts with focus-aware input protection.
// Prevents shortcuts from interfering with form input while enabling app-wide navigation.
type EventHandler struct {
	app         *tview.Application
	appState    *models.AppState
	nav         *models.NavigationState
	pageManager *layout.PageManager
	statusBar   *components.StatusBar
	detailView  *components.DetailView // Direct reference for password operations
	layoutMgr   *layout.LayoutManager  // Reference for layout manipulation
	config      *config.Config         // User configuration for keybindings
}

// NewEventHandler creates a new event handler with all required dependencies.
func NewEventHandler(
	app *tview.Application,
	appState *models.AppState,
	nav *models.NavigationState,
	pageManager *layout.PageManager,
	statusBar *components.StatusBar,
	detailView *components.DetailView,
	layoutMgr *layout.LayoutManager,
	cfg *config.Config,
) *EventHandler {
	return &EventHandler{
		app:         app,
		appState:    appState,
		nav:         nav,
		pageManager: pageManager,
		statusBar:   statusBar,
		detailView:  detailView,
		layoutMgr:   layoutMgr,
		config:      cfg,
	}
}

// SetupGlobalShortcuts installs the global keyboard shortcut handler.
// CRITICAL: Implements input protection to prevent intercepting form input.
func (eh *EventHandler) SetupGlobalShortcuts() {
	eh.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ✅ CRITICAL: If size warning is active, block ALL input except Ctrl+C
		// User must resize terminal - no interaction allowed
		if eh.pageManager.IsSizeWarningActive() {
			if event.Key() == tcell.KeyCtrlC {
				eh.handleQuit()
				return nil
			}
			// Block all other keys while size warning is displayed
			return nil
		}

		// ✅ CRITICAL: When a modal is open, only intercept Ctrl+C
		// Let modals handle Escape (for custom close logic like confirmation dialogs)
		// All other keys go to the modal/form
		if eh.pageManager.HasModals() {
			switch event.Key() {
			case tcell.KeyCtrlC:
				eh.handleQuit() // Closes modal or quits app
				return nil
			}
			// Let modal handle all other keys (including Escape, Tab, input, etc.)
			return event
		}

		// ✅ CRITICAL: Check if focused component should handle input
		focused := eh.app.GetFocus()
		if focused != nil {
			switch focused.(type) {
			case *tview.Form, *tview.InputField, *tview.TextArea, *tview.DropDown, *tview.Button:
				// ✅ Let form/input components handle their own keys (including Tab navigation)
				// Only intercept Ctrl+C for quit
				if event.Key() == tcell.KeyCtrlC {
					eh.handleQuit()
					return nil
				}
				return event // ✅ Pass all other keys to form component
			}
			// Note: TextView, Table, TreeView are NOT in the list above
			// They are read-only/navigation components that should allow global shortcuts
		}

		// Handle global shortcuts for non-input components
		return eh.handleGlobalKey(event)
	})
}

// handleGlobalKey routes keyboard events to appropriate action handlers.
// Only called when focus is NOT on an input component.
func (eh *EventHandler) handleGlobalKey(event *tcell.EventKey) *tcell.EventKey {
	// Check against configured keybindings
	if eh.config.MatchesKeybinding(event, "quit") {
		eh.handleQuit()
		return nil
	}
	if eh.config.MatchesKeybinding(event, "add_credential") {
		eh.handleNewCredential()
		return nil
	}
	if eh.config.MatchesKeybinding(event, "edit_credential") {
		eh.handleEditCredential()
		return nil
	}
	if eh.config.MatchesKeybinding(event, "delete_credential") {
		eh.handleDeleteCredential()
		return nil
	}
	if eh.config.MatchesKeybinding(event, "toggle_detail") {
		eh.handleToggleDetailPanel()
		return nil
	}
	if eh.config.MatchesKeybinding(event, "toggle_sidebar") {
		eh.handleToggleSidebar()
		return nil
	}
	if eh.config.MatchesKeybinding(event, "help") {
		eh.handleShowHelp()
		return nil
	}
	if eh.config.MatchesKeybinding(event, "search") {
		eh.handleSearchActivate()
		return nil
	}

	// Handle hardcoded keys that are not configurable
	switch event.Key() {
	case tcell.KeyTab:
		eh.handleTabFocus()
		return nil

	case tcell.KeyBacktab: // Shift+Tab
		eh.handleShiftTabFocus()
		return nil

	case tcell.KeyEscape:
		// Check if search is active, deactivate if so
		searchState := eh.appState.GetSearchState()
		if searchState != nil && searchState.Active {
			eh.handleSearchDeactivate()
			return nil
		}

	case tcell.KeyCtrlC:
		eh.handleQuit()
		return nil
	}

	// Additional non-configurable shortcuts (password operations)
	// These are not in the config spec but should stay as they are context-specific
	if event.Key() == tcell.KeyRune {
		switch event.Rune() {
		case 'p':
			eh.handleTogglePassword()
			return nil
		case 'c':
			eh.handleCopyPassword()
			return nil
		}
	}

	return event // Pass through unhandled keys
}

// handleQuit quits the application or closes the topmost modal.
func (eh *EventHandler) handleQuit() {
	// If modal is open, close it instead of quitting
	if eh.pageManager.HasModals() {
		eh.pageManager.CloseTopModal()
		return
	}

	// Quit application
	eh.app.Stop()
}

// handleNewCredential shows the add credential form modal.
func (eh *EventHandler) handleNewCredential() {
	form := components.NewAddForm(eh.appState)

	form.SetOnSubmit(func() {
		eh.pageManager.CloseModal("add-form")
		eh.statusBar.ShowSuccess("Credential added!")
	})

	form.SetOnCancel(func() {
		eh.pageManager.CloseModal("add-form")
	})

	form.SetOnCancelConfirm(func(message string, onYes func(), onNo func()) {
		eh.pageManager.ShowConfirmDialog("Confirm", message, onYes, onNo)
	})

	eh.pageManager.ShowModal("add-form", form, layout.FormModalWidth, layout.FormModalHeight)
}

// handleEditCredential shows the edit credential form for the selected credential.
func (eh *EventHandler) handleEditCredential() {
	cred := eh.appState.GetSelectedCredential()
	if cred == nil {
		eh.statusBar.ShowError(fmt.Errorf("no credential selected"))
		return
	}

	form := components.NewEditForm(eh.appState, cred)

	form.SetOnSubmit(func() {
		eh.pageManager.CloseModal("edit-form")
		eh.statusBar.ShowSuccess("Credential updated!")
	})

	form.SetOnCancel(func() {
		eh.pageManager.CloseModal("edit-form")
	})

	form.SetOnCancelConfirm(func(message string, onYes func(), onNo func()) {
		eh.pageManager.ShowConfirmDialog("Confirm", message, onYes, onNo)
	})

	eh.pageManager.ShowModal("edit-form", form, layout.FormModalWidth, layout.FormModalHeight)
}

// handleDeleteCredential shows a confirmation dialog before deleting the selected credential.
func (eh *EventHandler) handleDeleteCredential() {
	cred := eh.appState.GetSelectedCredential()
	if cred == nil {
		eh.statusBar.ShowError(fmt.Errorf("no credential selected"))
		return
	}

	message := fmt.Sprintf("Delete credential '%s'?\nThis action cannot be undone.", cred.Service)

	eh.pageManager.ShowConfirmDialog(
		"Delete Credential",
		message,
		func() {
			// Yes - delete credential
			err := eh.appState.DeleteCredential(cred.Service)
			if err != nil {
				eh.statusBar.ShowError(err)
			} else {
				eh.statusBar.ShowSuccess("Credential deleted")
			}
		},
		func() {
			// No - cancelled
		},
	)
}

// handleTogglePassword toggles password visibility in the detail view.
func (eh *EventHandler) handleTogglePassword() {
	if eh.detailView == nil {
		return
	}

	eh.detailView.TogglePasswordVisibility()
}

// handleCopyPassword copies the password of the selected credential to clipboard.
func (eh *EventHandler) handleCopyPassword() {
	if eh.detailView == nil {
		return
	}

	err := eh.detailView.CopyPasswordToClipboard()
	if err != nil {
		eh.statusBar.ShowError(err)
	} else {
		eh.statusBar.ShowSuccess("Password copied to clipboard!")
	}
}

// handleToggleDetailPanel toggles the detail panel visibility through three states.
// Cycles: Auto (responsive) -> Hide -> Show -> Auto
// Displays status bar message showing the new state.
func (eh *EventHandler) handleToggleDetailPanel() {
	if eh.layoutMgr == nil {
		return
	}

	message := eh.layoutMgr.ToggleDetailPanel()
	eh.statusBar.ShowInfo(message)
}

// handleToggleSidebar toggles the sidebar visibility through three states.
// Cycles: Auto (responsive) -> Hide -> Show -> Auto
// Displays status bar message showing the new state.
func (eh *EventHandler) handleToggleSidebar() {
	if eh.layoutMgr == nil {
		return
	}

	message := eh.layoutMgr.ToggleSidebar()
	eh.statusBar.ShowInfo(message)
}

// handleShowHelp displays a modal with keyboard shortcuts help.
func (eh *EventHandler) handleShowHelp() {
	// Helper to get display string for keybinding
	getKey := func(action string) string {
		keyStr := eh.config.GetKeybindingForAction(action)
		if keyStr == "" {
			return "?"
		}
		return config.GetDisplayString(keyStr)
	}

	// Create table for properly aligned shortcuts (scrollable with arrow keys)
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false). // Rows selectable for scrolling, columns not
		SetFixed(1, 0).             // Fix title row at top when scrolling
		SetSelectedStyle(tcell.StyleDefault.
			Background(tcell.ColorBlue).
			Foreground(tcell.ColorWhite)) // Keep selection invisible (same colors)

	row := 0

	// Title
	titleCell := tview.NewTableCell("Keyboard Shortcuts").
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorBlue).
		SetAlign(tview.AlignCenter).
		SetExpansion(1).
		SetAttributes(tcell.AttrBold)
	table.SetCell(row, 0, titleCell)
	table.SetCell(row, 1, tview.NewTableCell("").SetBackgroundColor(tcell.ColorBlue))
	row++

	// Separator
	separatorCell := tview.NewTableCell("══════════════════").
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorBlue).
		SetAlign(tview.AlignCenter).
		SetExpansion(1)
	table.SetCell(row, 0, separatorCell)
	table.SetCell(row, 1, tview.NewTableCell("").SetBackgroundColor(tcell.ColorBlue))
	row++
	row++ // Skip blank line row (will just be empty space)

	// Helper to add section header
	addSection := func(title string) {
		table.SetCell(row, 0, tview.NewTableCell(title).
			SetTextColor(tcell.ColorYellow).
			SetBackgroundColor(tcell.ColorBlue).
			SetAttributes(tcell.AttrBold).
			SetExpansion(1))
		table.SetCell(row, 1, tview.NewTableCell("").
			SetBackgroundColor(tcell.ColorBlue))
		row++
	}

	// Helper to add shortcut row
	addShortcut := func(key, description string) {
		table.SetCell(row, 0, tview.NewTableCell("  "+key).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorBlue).
			SetAlign(tview.AlignLeft))
		table.SetCell(row, 1, tview.NewTableCell(description).
			SetTextColor(tcell.ColorWhite).
			SetBackgroundColor(tcell.ColorBlue).
			SetAlign(tview.AlignLeft))
		row++
	}

	// Navigation section
	addSection("Navigation")
	addShortcut("Tab", "Next component")
	addShortcut("Shift+Tab", "Previous component")
	addShortcut("↑/↓", "Navigate lists")
	addShortcut("Enter", "Select / View details")
	row++ // Blank line (just skip row, don't add cells)

	// Actions section
	addSection("Actions")
	addShortcut(getKey("add_credential"), "New credential")
	addShortcut(getKey("edit_credential"), "Edit credential")
	addShortcut(getKey("delete_credential"), "Delete credential")
	addShortcut("p", "Toggle password visibility")
	addShortcut("c", "Copy password to clipboard")
	row++ // Blank line (just skip row, don't add cells)

	// View section
	addSection("View")
	addShortcut(getKey("toggle_detail"), "Toggle detail panel")
	addShortcut(getKey("toggle_sidebar"), "Toggle sidebar")
	addShortcut(getKey("search"), "Search / Filter credentials")
	row++ // Blank line (just skip row, don't add cells)

	// General section
	addSection("General")
	addShortcut(getKey("help"), "Show this help")
	addShortcut(getKey("quit"), "Quit application")
	addShortcut("Esc", "Close modal / Cancel search")
	addShortcut("Ctrl+C", "Quit application")

	// Set table background color (after all cells are set)
	table.SetBackgroundColor(tcell.ColorBlue)

	// Create TextView for close button (no SetTextAlign - it clips text!)
	closeButtonText := tview.NewTextView()
	closeButtonText.SetText("	PgUp/PgDn or Mouse Wheel to scroll  •  Esc to close")
	closeButtonText.SetTextColor(tcell.ColorWhite)
	closeButtonText.SetBackgroundColor(tcell.ColorBlue)

	// Make it close modal on Enter
	closeButtonText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			eh.pageManager.CloseModal("help")
			return nil
		}
		return event
	})

	// Add padding around table for better visual appearance
	paddedTable := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlue), 2, 0, false). // Left padding
		AddItem(table, 0, 1, true).                                               // Table (flex width, focusable)
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlue), 2, 0, false)  // Right padding

	// Combine padded table and button in vertical layout
	helpContent := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlue), 1, 0, false). // Top padding
		AddItem(paddedTable, 0, 1, true).                                         // Table (flex height, gets focus for scrolling)
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlue), 1, 0, false). // Spacer
		AddItem(closeButtonText, 1, 0, false).                                    // Close text (fixed 1 height)
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlue), 1, 0, false)  // Bottom padding

	helpContent.SetBackgroundColor(tcell.ColorBlue).
		SetBorder(true).
		SetTitle(" Help ").
		SetBorderColor(tcell.ColorWhite)

	// Add input capture to handle Escape key to close modal
	helpContent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			eh.pageManager.CloseModal("help")
			return nil
		}
		return event
	})

	eh.pageManager.ShowModal("help", helpContent, layout.HelpModalWidth, layout.HelpModalHeight)
}

// handleTabFocus cycles focus to the next component in tab order.
func (eh *EventHandler) handleTabFocus() {
	eh.nav.CycleFocus()
}

// handleShiftTabFocus cycles focus to the previous component in reverse tab order.
func (eh *EventHandler) handleShiftTabFocus() {
	eh.nav.CycleFocusReverse()
}

// handleSearchActivate activates search mode.
func (eh *EventHandler) handleSearchActivate() {
	searchState := eh.appState.GetSearchState()
	if searchState == nil {
		return
	}

	// Activate search (creates InputField)
	searchState.Activate()

	// Get table reference for arrow key forwarding
	table := eh.appState.GetTable()

	// Setup input capture to forward navigation keys to table
	searchState.InputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp, tcell.KeyDown, tcell.KeyEnter:
			// Forward navigation keys to table
			if table != nil {
				// Simulate key press on table
				table.InputHandler()(event, nil)
				return nil // Consume event
			}
		case tcell.KeyEscape:
			// Exit search input but keep filter active
			// User can press p/c to view/copy password of filtered results
			table := eh.appState.GetTable()
			if table != nil {
				eh.app.SetFocus(table)
			}
			eh.statusBar.ShowInfo("Filter active. Press '/' to edit search, Esc again to clear.")
			return nil
		}
		// Let all other keys (typing) pass through to InputField
		return event
	})

	// Setup real-time filtering callback (T035)
	searchState.InputField.SetChangedFunc(func(text string) {
		// Update query in real-time
		searchState.Query = text

		// Trigger filter changed callback to refresh table only (not detail view)
		eh.appState.TriggerFilterChanged()
	})

	// Setup done function to handle Escape (redundant but safe)
	searchState.InputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			eh.handleSearchDeactivate()
		}
	})

	// Rebuild layout to show InputField (T034)
	eh.layoutMgr.RebuildLayout()

	// Set focus to the InputField
	eh.app.SetFocus(searchState.InputField)

	// Update statusbar to show search shortcuts
	eh.statusBar.UpdateForContext(components.FocusTable)

	eh.statusBar.ShowInfo("Search mode activated. Type to filter, Esc to exit.")
}

// handleSearchDeactivate deactivates search mode and clears the filter.
func (eh *EventHandler) handleSearchDeactivate() {
	searchState := eh.appState.GetSearchState()
	if searchState == nil {
		return
	}

	// Deactivate search (clears query and destroys InputField)
	searchState.Deactivate()

	// Rebuild layout to remove InputField
	eh.layoutMgr.RebuildLayout()

	// Trigger table refresh to clear filter
	eh.appState.TriggerRefresh()

	// Return focus to table
	table := eh.appState.GetTable()
	if table != nil {
		eh.app.SetFocus(table)
	}

	// Update statusbar to restore normal shortcuts
	eh.statusBar.UpdateForContext(components.FocusTable)

	eh.statusBar.ShowInfo("Search cleared")
}
