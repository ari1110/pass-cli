package components

import (
	"fmt"
	"time"

	"github.com/rivo/tview"

	"pass-cli/cmd/tui-tview/models"
	"pass-cli/cmd/tui-tview/styles"
)

// FocusContext represents the current focus context for determining which shortcuts to display.
type FocusContext int

const (
	// FocusSidebar indicates the sidebar is focused
	FocusSidebar FocusContext = iota
	// FocusTable indicates the credential table is focused
	FocusTable
	// FocusDetail indicates the detail view is focused
	FocusDetail
	// FocusModal indicates a modal (form or dialog) is focused
	FocusModal
)

// StatusBar displays context-aware keyboard shortcuts and temporary status messages.
type StatusBar struct {
	*tview.TextView

	appState     *models.AppState
	currentFocus FocusContext
	messageTimer *time.Timer
}

// NewStatusBar creates and initializes a new status bar.
func NewStatusBar(appState *models.AppState) *StatusBar {
	theme := styles.GetCurrentTheme()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	// Configure styling: no borders, dark background, fixed height
	textView.SetBackgroundColor(theme.StatusBarBg).
		SetBorder(false)

	sb := &StatusBar{
		TextView:     textView,
		appState:     appState,
		currentFocus: FocusSidebar, // Default focus
	}

	// Set initial shortcuts display
	sb.UpdateForContext(FocusSidebar)

	return sb
}

// UpdateForContext updates the displayed shortcuts based on the current focus context.
func (sb *StatusBar) UpdateForContext(focus FocusContext) {
	sb.currentFocus = focus
	shortcuts := sb.getShortcutsForContext(focus)
	sb.SetText(shortcuts)
}

// ShowSuccess displays a temporary success message (green text, 3 seconds).
func (sb *StatusBar) ShowSuccess(message string) {
	formatted := fmt.Sprintf("[green]%s[-]", message)
	sb.showTemporaryMessage(formatted, 3*time.Second)
}

// ShowError displays a temporary error message (red text, 5 seconds).
func (sb *StatusBar) ShowError(err error) {
	formatted := fmt.Sprintf("[red]Error: %s[-]", err.Error())
	sb.showTemporaryMessage(formatted, 5*time.Second)
}

// showTemporaryMessage displays a message for the specified duration, then restores shortcuts.
func (sb *StatusBar) showTemporaryMessage(message string, duration time.Duration) {
	// Cancel previous message timer if it exists
	if sb.messageTimer != nil {
		sb.messageTimer.Stop()
	}

	// Display the message
	sb.SetText(message)

	// Schedule restoration of shortcuts after duration
	sb.messageTimer = time.AfterFunc(duration, func() {
		sb.UpdateForContext(sb.currentFocus)
	})
}

// getShortcutsForContext returns the appropriate shortcut text for the given focus context.
func (sb *StatusBar) getShortcutsForContext(focus FocusContext) string {
	switch focus {
	case FocusSidebar:
		return "[gray][Tab] Next  [↑↓] Navigate  [Enter] Select  [n] New  [?] Help  [q] Quit[-]"

	case FocusTable:
		return "[gray][n] New  [e] Edit  [d] Delete  [c] Copy  [?] Help  [q] Quit[-]"

	case FocusDetail:
		return "[gray][e] Edit  [d] Delete  [p] Toggle  [c] Copy  [?] Help  [q] Quit[-]"

	case FocusModal:
		return "[gray][Tab] Next Field  [Enter] Submit  [Esc] Cancel[-]"

	default:
		return "[gray][?] Help  [q] Quit[-]"
	}
}
