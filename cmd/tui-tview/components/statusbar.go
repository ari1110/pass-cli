package components

import (
	"fmt"
	"time"

	"github.com/rivo/tview"

	"pass-cli/cmd/tui/models"
	"pass-cli/cmd/tui/styles"
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

	app          *tview.Application // For forcing redraws
	appState     *models.AppState
	currentFocus FocusContext
	messageTimer *time.Timer
}

// NewStatusBar creates and initializes a new status bar.
func NewStatusBar(app *tview.Application, appState *models.AppState) *StatusBar {
	theme := styles.GetCurrentTheme()

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	// Configure styling: no borders, dark background, fixed height
	textView.SetBackgroundColor(theme.StatusBarBg).
		SetBorder(false)

	sb := &StatusBar{
		TextView:     textView,
		app:          app,
		appState:     appState,
		currentFocus: FocusSidebar, // Default focus
	}

	// Set initial shortcuts display (direct SetText, no queue - app not running yet)
	shortcuts := sb.getShortcutsForContext(FocusSidebar)
	sb.SetText(shortcuts)

	return sb
}

// UpdateForContext updates the displayed shortcuts based on the current focus context.
func (sb *StatusBar) UpdateForContext(focus FocusContext) {
	sb.currentFocus = focus
	shortcuts := sb.getShortcutsForContext(focus)

	// Direct SetText is sufficient - tview redraws automatically on next frame
	sb.SetText(shortcuts)
}

// ShowSuccess displays a temporary success message (green text, 3 seconds).
func (sb *StatusBar) ShowSuccess(message string) {
	formatted := fmt.Sprintf("[green]%s[-]", message)
	sb.showTemporaryMessage(formatted, 3*time.Second)
}

// ShowInfo displays a temporary info message (cyan text, 3 seconds).
func (sb *StatusBar) ShowInfo(message string) {
	formatted := fmt.Sprintf("[cyan]%s[-]", message)
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
		return "[yellow]Tab[white]/[yellow]Shift+Tab[-]:Switch  [yellow]↑↓[-]:Nav  [yellow]Enter[-]:Select  [yellow]n[-]:New  [yellow]i[-]:Details  [yellow]?[-]:Help  [yellow]q[-]:Quit"

	case FocusTable:
		return "[yellow]Tab[white]/[yellow]Shift+Tab[-]:Switch  [yellow]↑↓[-]:Nav  [yellow]n[-]:New  [yellow]e[-]:Edit  [yellow]d[-]:Del  [yellow]c[-]:Copy  [yellow]i[-]:Details  [yellow]?[-]:Help  [yellow]q[-]:Quit"

	case FocusDetail:
		return "[yellow]Tab[white]/[yellow]Shift+Tab[-]:Switch  [yellow]e[-]:Edit  [yellow]d[-]:Del  [yellow]p[-]:Toggle  [yellow]c[-]:Copy  [yellow]i[-]:Details  [yellow]?[-]:Help  [yellow]q[-]:Quit"

	case FocusModal:
		return "[yellow]Tab[white]/[yellow]Shift+Tab[-]:Field  [yellow]Enter[-]:Submit  [yellow]Esc[-]:Cancel"

	default:
		return "[yellow]Tab[white]/[yellow]Shift+Tab[-]:Switch  [yellow]i[-]:Details  [yellow]?[-]:Help  [yellow]q[-]:Quit"
	}
}
