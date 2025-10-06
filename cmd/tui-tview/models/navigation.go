package models

import (
	"github.com/rivo/tview"
)

// FocusableComponent represents components that can receive focus.
type FocusableComponent int

const (
	FocusSidebar FocusableComponent = iota
	FocusTable
	FocusDetail
)

// NavigationState manages focus navigation between components.
// Provides Tab cycling and focus management with proper order tracking.
type NavigationState struct {
	app *tview.Application

	focusOrder   []tview.Primitive
	currentIndex int

	onFocusChanged func(FocusableComponent)
}

// NewNavigationState creates a new navigation state manager.
// Initialize with empty focus order - components will be registered later.
func NewNavigationState(app *tview.Application) *NavigationState {
	return &NavigationState{
		app:          app,
		focusOrder:   make([]tview.Primitive, 0),
		currentIndex: 0,
	}
}

// SetFocusOrder sets the components that can receive focus in tab order.
// Order: Sidebar -> Table -> Detail -> (back to Sidebar)
func (ns *NavigationState) SetFocusOrder(order []tview.Primitive) {
	ns.focusOrder = order
	ns.currentIndex = 0
}

// CycleFocus moves focus to the next component in the focus order.
// Used for Tab key navigation.
func (ns *NavigationState) CycleFocus() {
	if len(ns.focusOrder) == 0 {
		return
	}

	ns.currentIndex = (ns.currentIndex + 1) % len(ns.focusOrder)
	ns.setFocus(ns.currentIndex)
}

// CycleFocusReverse moves focus to the previous component in the focus order.
// Used for Shift+Tab navigation.
func (ns *NavigationState) CycleFocusReverse() {
	if len(ns.focusOrder) == 0 {
		return
	}

	ns.currentIndex--
	if ns.currentIndex < 0 {
		ns.currentIndex = len(ns.focusOrder) - 1
	}
	ns.setFocus(ns.currentIndex)
}

// SetFocus directly sets focus to a specific component.
func (ns *NavigationState) SetFocus(target FocusableComponent) {
	if int(target) < len(ns.focusOrder) {
		ns.currentIndex = int(target)
		ns.setFocus(ns.currentIndex)
	}
}

// GetCurrentFocus returns the currently focused component.
func (ns *NavigationState) GetCurrentFocus() FocusableComponent {
	if ns.currentIndex < 0 || ns.currentIndex >= len(ns.focusOrder) {
		return FocusSidebar
	}
	return FocusableComponent(ns.currentIndex)
}

// SetOnFocusChanged registers a callback to be invoked when focus changes.
func (ns *NavigationState) SetOnFocusChanged(callback func(FocusableComponent)) {
	ns.onFocusChanged = callback
}

// setFocus is an internal helper that updates focus and triggers callbacks.
func (ns *NavigationState) setFocus(index int) {
	if index < 0 || index >= len(ns.focusOrder) {
		return
	}

	primitive := ns.focusOrder[index]
	ns.app.SetFocus(primitive)

	if ns.onFocusChanged != nil {
		ns.onFocusChanged(FocusableComponent(index))
	}
}
