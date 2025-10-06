package layout

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// PageManager manages modal dialogs and page switching using tview.Pages.
// It handles showing/hiding forms and dialogs over the main UI with proper
// stack management for nested modals.
//
// Responsibilities:
// - Page stacking: Layer modals over main UI
// - Modal display: Show forms, dialogs, confirmations
// - Page removal: Close modals and return to main UI
// - Focus management: Restore focus when closing modals
// - Escape handling: Close topmost modal on Escape key
type PageManager struct {
	*tview.Pages

	app        *tview.Application
	modalStack []string // Track modal names for proper close operations
}

// NewPageManager creates a new page manager for handling modals and page switching.
// The Pages primitive serves as the root of the application.
func NewPageManager(app *tview.Application) *PageManager {
	pages := tview.NewPages()

	pm := &PageManager{
		Pages:      pages,
		app:        app,
		modalStack: []string{},
	}

	pm.setupEscapeHandler()

	return pm
}

// ShowPage adds a non-modal page to the page manager.
// This is typically used to add the main layout as the base page.
func (pm *PageManager) ShowPage(name string, primitive tview.Primitive) *PageManager {
	pm.AddPage(name, primitive, true, true)
	return pm
}

// SwitchToPage changes the active page without modal management.
func (pm *PageManager) SwitchToPage(name string) *PageManager {
	pm.Pages.SwitchToPage(name)
	return pm
}

// ShowModal displays a modal over the current page with specified dimensions.
// The modal is centered on screen and added to the modal stack.
//
// Parameters:
//   - name: Unique identifier for this modal
//   - modal: The primitive to display (form, dialog, etc.)
//   - width, height: Modal dimensions (use 0 for proportional sizing)
func (pm *PageManager) ShowModal(name string, modal tview.Primitive, width, height int) *PageManager {
	// Center the modal using Flex layouts
	centered := pm.centerModal(modal, width, height)

	pm.AddPage(name, centered, true, true)
	pm.modalStack = append(pm.modalStack, name)

	return pm
}

// ShowForm displays a credential form as a centered modal dialog.
// This is a convenience wrapper around ShowModal specifically for forms.
func (pm *PageManager) ShowForm(form *tview.Form, title string) *PageManager {
	// Set form title
	form.SetTitle(" " + title + " ")
	form.SetBorder(true)

	// Forms are shown with fixed dimensions (60 width, 20 height)
	return pm.ShowModal("form", form, 60, 20)
}

// ShowConfirmDialog displays a yes/no confirmation dialog with callbacks.
// The dialog automatically closes when a button is pressed and calls the
// appropriate callback.
//
// Parameters:
//   - title: Dialog title (not used with tview.Modal, but kept for API consistency)
//   - message: Confirmation message to display
//   - onYes: Callback to execute when "Yes" is pressed
//   - onNo: Callback to execute when "No" is pressed
func (pm *PageManager) ShowConfirmDialog(title, message string, onYes, onNo func()) *PageManager {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pm.CloseTopModal()
			if buttonIndex == 0 {
				if onYes != nil {
					onYes()
				}
			} else {
				if onNo != nil {
					onNo()
				}
			}
		})

	return pm.ShowModal("confirm", modal, 60, 10)
}

// CloseModal removes a modal by name and pops it from the stack.
// If the modal is not found, this is a safe no-op.
func (pm *PageManager) CloseModal(name string) {
	pm.RemovePage(name)

	// Remove from stack
	for i, page := range pm.modalStack {
		if page == name {
			pm.modalStack = append(pm.modalStack[:i], pm.modalStack[i+1:]...)
			break
		}
	}

	// If no more modals, ensure we're back on main page
	if len(pm.modalStack) == 0 {
		pm.SwitchToPage("main")
	}
}

// CloseTopModal closes the most recently opened modal.
// If no modals are open, this is a safe no-op.
func (pm *PageManager) CloseTopModal() {
	if len(pm.modalStack) > 0 {
		topModal := pm.modalStack[len(pm.modalStack)-1]
		pm.CloseModal(topModal)
	}
}

// HasModals returns true if any modals are currently displayed.
func (pm *PageManager) HasModals() bool {
	return len(pm.modalStack) > 0
}

// centerModal wraps a modal primitive in Flex layouts to center it on screen.
// Uses the width and height to determine fixed or proportional sizing.
func (pm *PageManager) centerModal(modal tview.Primitive, width, height int) tview.Primitive {
	// Create horizontal centering
	hFlex := tview.NewFlex().
		AddItem(nil, 0, 1, false).      // Left spacer (flex)
		AddItem(modal, width, 0, true). // Modal (fixed width)
		AddItem(nil, 0, 1, false)       // Right spacer (flex)

	// Create vertical + horizontal centering
	vFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).       // Top spacer
		AddItem(hFlex, height, 0, true). // Middle row
		AddItem(nil, 0, 1, false)        // Bottom spacer

	return vFlex
}

// setupEscapeHandler configures Escape key to close the topmost modal.
// If no modals are open, the Escape key passes through to underlying components.
func (pm *PageManager) setupEscapeHandler() {
	pm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			if len(pm.modalStack) > 0 {
				pm.CloseTopModal()
				return nil // Consume event
			}
		}
		return event // Pass through
	})
}
