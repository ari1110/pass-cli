package events

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"pass-cli/cmd/tui-tview/components"
	"pass-cli/cmd/tui-tview/layout"
	"pass-cli/cmd/tui-tview/models"
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
}

// NewEventHandler creates a new event handler with all required dependencies.
func NewEventHandler(
	app *tview.Application,
	appState *models.AppState,
	nav *models.NavigationState,
	pageManager *layout.PageManager,
	statusBar *components.StatusBar,
	detailView *components.DetailView,
) *EventHandler {
	return &EventHandler{
		app:         app,
		appState:    appState,
		nav:         nav,
		pageManager: pageManager,
		statusBar:   statusBar,
		detailView:  detailView,
	}
}

// SetupGlobalShortcuts installs the global keyboard shortcut handler.
// CRITICAL: Implements input protection to prevent intercepting form input.
func (eh *EventHandler) SetupGlobalShortcuts() {
	eh.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ✅ CRITICAL: Check if focused component should handle input
		focused := eh.app.GetFocus()
		if focused != nil {
			switch focused.(type) {
			case *tview.Form, *tview.InputField, *tview.TextArea, *tview.DropDown, *tview.TextView, *tview.Button:
				// ✅ Let form components handle their own keys (including navigation)
				// Only intercept Ctrl+C for quit
				if event.Key() == tcell.KeyCtrlC {
					eh.handleQuit()
					return nil
				}
				return event // ✅ Pass all other keys to form component
			}
		}

		// Handle global shortcuts for non-input components
		return eh.handleGlobalKey(event)
	})
}

// handleGlobalKey routes keyboard events to appropriate action handlers.
// Only called when focus is NOT on an input component.
func (eh *EventHandler) handleGlobalKey(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'q':
			eh.handleQuit()
			return nil
		case 'n':
			eh.handleNewCredential()
			return nil
		case 'e':
			eh.handleEditCredential()
			return nil
		case 'd':
			eh.handleDeleteCredential()
			return nil
		case 'p':
			eh.handleTogglePassword()
			return nil
		case 'c':
			eh.handleCopyPassword()
			return nil
		case '?':
			eh.handleShowHelp()
			return nil
		}

	case tcell.KeyTab:
		eh.handleTabFocus()
		return nil

	case tcell.KeyBacktab: // Shift+Tab
		eh.handleShiftTabFocus()
		return nil

	case tcell.KeyCtrlC:
		eh.handleQuit()
		return nil
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

// handleShowHelp displays a modal with keyboard shortcuts help.
func (eh *EventHandler) handleShowHelp() {
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
	addShortcut("n", "New credential")
	addShortcut("e", "Edit credential")
	addShortcut("d", "Delete credential")
	addShortcut("p", "Toggle password visibility")
	addShortcut("c", "Copy password to clipboard")
	row++ // Blank line (just skip row, don't add cells)

	// General section
	addSection("General")
	addShortcut("?", "Show this help")
	addShortcut("q", "Quit application")
	addShortcut("Esc", "Close modal / Go back")
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
