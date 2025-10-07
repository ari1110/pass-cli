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
			case *tview.Form, *tview.InputField, *tview.TextArea:
				// ✅ Let input components handle their own keys
				// Only intercept Ctrl+C for quit
				if event.Key() == tcell.KeyCtrlC {
					eh.handleQuit()
					return nil
				}
				return event // ✅ Pass all other keys to input component
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
	helpText := `[yellow]═══════════════════════════════════════════════════[-]
[yellow]              Keyboard Shortcuts[-]
[yellow]═══════════════════════════════════════════════════[-]

[cyan]Navigation[-]
  [white]Tab[-]          [gray]Next component[-]
  [white]Shift+Tab[-]    [gray]Previous component[-]
  [white]↑/↓[-]          [gray]Navigate lists[-]
  [white]Enter[-]        [gray]Select / View details[-]

[cyan]Actions[-]
  [white]n[-]            [gray]New credential[-]
  [white]e[-]            [gray]Edit credential[-]
  [white]d[-]            [gray]Delete credential[-]
  [white]p[-]            [gray]Toggle password visibility[-]
  [white]c[-]            [gray]Copy password to clipboard[-]

[cyan]General[-]
  [white]?[-]            [gray]Show this help[-]
  [white]q[-]            [gray]Quit application[-]
  [white]Esc[-]          [gray]Close modal / Go back[-]
  [white]Ctrl+C[-]       [gray]Quit application[-]
`

	modal := tview.NewModal().
		SetText(helpText).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			eh.pageManager.CloseModal("help")
		})

	eh.pageManager.ShowModal("help", modal, layout.HelpModalWidth, layout.HelpModalHeight)
}

// handleTabFocus cycles focus to the next component in tab order.
func (eh *EventHandler) handleTabFocus() {
	eh.nav.CycleFocus()
}

// handleShiftTabFocus cycles focus to the previous component in reverse tab order.
func (eh *EventHandler) handleShiftTabFocus() {
	eh.nav.CycleFocusReverse()
}
