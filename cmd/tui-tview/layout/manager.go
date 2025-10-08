// Package layout provides responsive layout management for the TUI.
package layout

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"pass-cli/cmd/tui-tview/models"
)

// LayoutMode represents the current layout configuration based on terminal width.
type LayoutMode int

const (
	// LayoutSmall is for terminals < 80 columns (table only, no sidebar)
	LayoutSmall LayoutMode = iota
	// LayoutMedium is for terminals 80-120 columns (sidebar + table)
	LayoutMedium
	// LayoutLarge is for terminals > 120 columns (sidebar + table + detail)
	LayoutLarge
)

// LayoutManager manages responsive layout composition and terminal resize handling.
// It adapts the UI based on terminal width using breakpoints.
type LayoutManager struct {
	app      *tview.Application
	appState *models.AppState

	// Current terminal dimensions
	width  int
	height int

	// Current layout mode
	currentMode LayoutMode

	// Layout primitives
	mainLayout *tview.Flex // Root layout (vertical: contentRow + statusBar)
	contentRow *tview.Flex // Main content area (horizontal: sidebar + table + detail)

	// Component references (retrieved from AppState, not created new)
	sidebar    *tview.TreeView
	table      *tview.Table
	detailView *tview.TextView
	statusBar  *tview.TextView

	// Breakpoints (configurable)
	mediumBreakpoint int // Default: 80
	largeBreakpoint  int // Default: 120
}

// NewLayoutManager creates a new layout manager with default breakpoints.
// Components are not retrieved until CreateMainLayout is called.
func NewLayoutManager(app *tview.Application, appState *models.AppState) *LayoutManager {
	return &LayoutManager{
		app:              app,
		appState:         appState,
		mediumBreakpoint: 80,
		largeBreakpoint:  120,
		currentMode:      LayoutSmall, // Start with small, will adjust on first draw
	}
}

// CreateMainLayout builds the complete layout structure.
// Structure:
//
//	┌───────────────────────────────┐
//	│   contentRow (FlexColumn)     │  ← Main content area (flex)
//	│   ┌──────┬───────┬──────────┐ │
//	│   │ Side │ Table │  Detail  │ │
//	│   │ bar  │       │          │ │
//	│   └──────┴───────┴──────────┘ │
//	├───────────────────────────────┤
//	│       statusBar (1 row)       │  ← Fixed height
//	└───────────────────────────────┘
func (lm *LayoutManager) CreateMainLayout() *tview.Flex {
	// Get component references from appState (don't create new ones)
	lm.sidebar = lm.appState.GetSidebar()
	lm.table = lm.appState.GetTable()
	lm.detailView = lm.appState.GetDetailView()
	lm.statusBar = lm.appState.GetStatusBar()

	// Create content row (horizontal layout for main panels)
	lm.contentRow = tview.NewFlex().SetDirection(tview.FlexColumn)

	// Create main layout (vertical: content + status bar)
	lm.mainLayout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(lm.contentRow, 0, 1, true). // Content area (flex height)
		AddItem(lm.statusBar, 1, 0, false)  // Status bar (fixed 1 row)

	// Setup resize detection using SetDrawFunc
	// This detects the initial terminal size and subsequent resizes
	lm.mainLayout.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// Always check for size changes (not just initial)
		termWidth, termHeight := screen.Size()
		if lm.width != termWidth || lm.height != termHeight {
			lm.HandleResize(termWidth, termHeight)
		}
		return x, y, width, height
	})

	return lm.mainLayout
}

// HandleResize responds to terminal size changes.
// It determines if the layout mode needs to change and triggers a rebuild if necessary.
func (lm *LayoutManager) HandleResize(width, height int) {
	lm.width = width
	lm.height = height

	// Determine the new layout mode based on width
	newMode := lm.determineLayoutMode(width)

	// Only rebuild layout if the mode changed
	if newMode != lm.currentMode {
		lm.currentMode = newMode
		lm.rebuildLayout()
	}
}

// rebuildLayout reconstructs the layout based on the current mode.
// It clears the content row and adds components according to breakpoint rules:
//   - Small: Table only (full width)
//   - Medium: Sidebar (20 cols) + Table (flex)
//   - Large: Sidebar (20 cols) + Table (flex) + Detail (40 cols)
func (lm *LayoutManager) rebuildLayout() {
	// Skip rebuild if layout hasn't been initialized yet
	if lm.contentRow == nil {
		return
	}

	// Clear existing content
	lm.contentRow.Clear()

	// Build layout based on current mode
	switch lm.currentMode {
	case LayoutSmall:
		// Table only (full width)
		// size=0, proportion=1, focus=true means flex width, takes all space, can receive focus
		lm.contentRow.AddItem(lm.table, 0, 1, true)

	case LayoutMedium:
		// Sidebar + Table
		// Sidebar: size=20, proportion=0, focus=false means fixed 20 cols, no flex, no focus
		// Table: size=0, proportion=1, focus=true means flex width, takes remaining space, can receive focus
		lm.contentRow.
			AddItem(lm.sidebar, 20, 0, false).
			AddItem(lm.table, 0, 1, true)

	case LayoutLarge:
		// Sidebar + Table + Detail
		// Sidebar: fixed 20 cols
		// Table: flex width (takes remaining space between sidebar and detail)
		// Detail: fixed 40 cols
		lm.contentRow.
			AddItem(lm.sidebar, 20, 0, false).
			AddItem(lm.table, 0, 1, true).
			AddItem(lm.detailView, 40, 0, false)
	}
}

// determineLayoutMode calculates the appropriate layout mode based on terminal width.
// Breakpoint rules:
//   - width < mediumBreakpoint (default 80): LayoutSmall
//   - width >= mediumBreakpoint && width < largeBreakpoint (default 80-119): LayoutMedium
//   - width >= largeBreakpoint (default 120+): LayoutLarge
func (lm *LayoutManager) determineLayoutMode(width int) LayoutMode {
	if width < lm.mediumBreakpoint {
		return LayoutSmall
	}
	if width < lm.largeBreakpoint {
		return LayoutMedium
	}
	return LayoutLarge
}

// GetCurrentMode returns the current layout mode.
func (lm *LayoutManager) GetCurrentMode() LayoutMode {
	return lm.currentMode
}

// SetBreakpoints allows customizing the layout breakpoints.
// Default values are medium=80, large=120.
func (lm *LayoutManager) SetBreakpoints(medium, large int) {
	lm.mediumBreakpoint = medium
	lm.largeBreakpoint = large
}
