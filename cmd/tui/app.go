package tui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NewApp creates and configures a new tview.Application instance with mouse support
// and panic handling. The application uses tview's built-in alternate screen support.
func NewApp() *tview.Application {
	app := tview.NewApplication()

	// Enable mouse support by default for click and scroll interactions
	app.EnableMouse(true)

	return app
}

// SetRootSafely updates the root primitive safely using QueueUpdateDraw to avoid
// race conditions when changing the root while the application is running.
func SetRootSafely(app *tview.Application, root tview.Primitive, fullscreen bool) {
	app.QueueUpdateDraw(func() {
		app.SetRoot(root, fullscreen)
	})
}

// Quit gracefully shuts down the application by stopping the event loop.
// tview automatically handles terminal restoration and alternate screen cleanup.
func Quit(app *tview.Application) {
	app.Stop()
}

// RestoreTerminal performs emergency terminal restoration in case of panic.
// This should be called with defer in main() to ensure the terminal is properly
// restored even if the application crashes.
func RestoreTerminal() {
	if r := recover(); r != nil {
		// Attempt to restore terminal state
		screen, err := tcell.NewScreen()
		if err == nil {
			screen.Fini()
		}
		fmt.Fprintf(os.Stderr, "Panic: %v\n", r)
		os.Exit(1)
	}
}
