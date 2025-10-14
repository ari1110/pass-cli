package layout

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Modal dimension constants to ensure consistent sizing across all modals.
const (
	FormModalWidth  = 60 // Standard width for credential forms (add, edit)
	FormModalHeight = 27 // Standard height for 6-field forms + buttons + keyboard hints

	ConfirmDialogWidth  = 60 // Width for confirmation dialogs
	ConfirmDialogHeight = 10 // Height for yes/no confirmation dialogs

	HelpModalWidth  = 60 // Width for help screen modal
	HelpModalHeight = 25 // Height for help screen content
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

	sizeWarningActive bool // Track whether terminal size warning is currently displayed
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

	// Use standard form dimensions from constants
	return pm.ShowModal("form", form, FormModalWidth, FormModalHeight)
}

// ShowModalWithAutoHeight displays a form modal with auto-calculated height.
// Computes height based on form field count and caps to available screen size.
// This prevents overflow on small terminals while adapting to form complexity.
//
// Height calculation:
//   - Each field/button requires ~2 rows (label + input/spacing)
//   - Add 6 rows for borders, padding, title
//   - Cap at terminalHeight - 4 to leave breathing room
//
// Use this for forms with variable field counts or when targeting small terminals.
func (pm *PageManager) ShowModalWithAutoHeight(name string, form *tview.Form, width int) *PageManager {
	// Get terminal dimensions using tview's Box primitive as a proxy
	// Note: tview.Application doesn't expose screen directly, so we estimate
	// based on typical terminal sizes. In practice, forms with fixed height=25
	// work well for terminals >= 30 rows (vast majority of terminals).
	// This method is provided as an optional enhancement for future use cases.

	// Calculate height from form item count
	itemCount := form.GetFormItemCount() + form.GetButtonCount()
	calculatedHeight := itemCount*2 + 6 // 2 rows per item, 6 for chrome

	// Use a conservative maximum (assume 40-row terminal minimum)
	// For very small terminals (<30 rows), the modal may still touch edges,
	// but tview's rendering will gracefully degrade.
	maxHeight := 30
	if calculatedHeight > maxHeight {
		calculatedHeight = maxHeight
	}

	return pm.ShowModal(name, form, width, calculatedHeight)
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

	return pm.ShowModal("confirm", modal, ConfirmDialogWidth, ConfirmDialogHeight)
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

// setupEscapeHandler is intentionally left empty.
// Escape key handling is delegated to individual modals/forms via their InputCapture handlers,
// allowing them to implement custom logic (e.g., confirmation dialogs before closing).
// The global handler in handlers.go (handleQuit) provides the fallback Escape behavior.
func (pm *PageManager) setupEscapeHandler() {
	// No-op: Let forms and global handler manage Escape key
}

// ShowSizeWarning displays a blocking warning overlay when terminal is too small.
// Shows current dimensions vs. minimum required dimensions in plain language.
// The warning uses a dark red background for visual distinction.
//
// Parameters:
//   - currentWidth, currentHeight: Current terminal dimensions in columns/rows
//   - minWidth, minHeight: Minimum required dimensions in columns/rows
func (pm *PageManager) ShowSizeWarning(currentWidth, currentHeight, minWidth, minHeight int) {
	message := fmt.Sprintf(
		"Terminal too small!\n\nCurrent: %dx%d\nMinimum required: %dx%d\n\nPlease resize your terminal window.",
		currentWidth, currentHeight, minWidth, minHeight,
	)

	modal := tview.NewModal().
		SetText(message).
		SetBackgroundColor(tcell.ColorDarkRed)

	// Add as page with centering (reuse centerModal for consistency)
	centered := pm.centerModal(modal, 60, 12)
	pm.AddPage("size-warning", centered, true, true)

	pm.sizeWarningActive = true
	pm.app.Draw()
}

// HideSizeWarning removes the terminal size warning overlay.
// Safe to call even if warning is not currently showing (idempotent).
func (pm *PageManager) HideSizeWarning() {
	if !pm.sizeWarningActive {
		return
	}

	pm.RemovePage("size-warning")
	pm.sizeWarningActive = false
	pm.app.Draw()
}

// IsSizeWarningActive returns true if the terminal size warning is currently displayed.
func (pm *PageManager) IsSizeWarningActive() bool {
	return pm.sizeWarningActive
}
