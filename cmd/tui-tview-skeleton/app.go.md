# app.go

## Purpose
Manages the tview.Application lifecycle, including setup, configuration, and graceful shutdown.

## Responsibilities

1. **Application configuration**: Set up tview.Application with proper settings
2. **Screen management**: Handle alternate screen, mouse support, colors
3. **Root primitive management**: Set and update the root UI element
4. **Shutdown handling**: Clean exit and resource cleanup
5. **Panic recovery**: Catch panics and restore terminal state

## Dependencies

### External Dependencies
- `github.com/rivo/tview` - Application instance
- `github.com/gdamore/tcell/v2` - Screen and terminal control

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For accessing AppState
- None - This is a low-level utility module

## Key Functions

### `NewApp() *tview.Application`
**Purpose**: Create and configure a new tview.Application instance

**Configuration**:
- Enable alternate screen (automatic with Run())
- Enable mouse support (optional, can be disabled)
- Set color support (auto-detect terminal capabilities)
- Configure panic handler to restore terminal

**Returns**: Configured tview.Application ready to use

### `SetRootSafely(app *tview.Application, root tview.Primitive)`
**Purpose**: Safely update the root primitive without race conditions

**Why needed**: Changing root while app is running can cause issues

**Approach**: Use app.QueueUpdateDraw() to schedule root change

### `Quit(app *tview.Application)`
**Purpose**: Gracefully shut down the application

**Steps**:
1. Call app.Stop() to exit event loop
2. No need for manual screen restoration (tview handles it)

### `RestoreTerminal()`
**Purpose**: Emergency terminal restoration if app panics

**When used**: Defer this in main() to ensure terminal is restored on panic

## Example Structure

```go
package main

import (
    "fmt"
    "os"

    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

// NewApp creates and configures a new tview.Application
func NewApp() *tview.Application {
    app := tview.NewApplication()

    // Enable mouse support (optional, can disable for pure keyboard)
    app.EnableMouse(true)

    // Set up panic handler to restore terminal
    defer func() {
        if r := recover(); r != nil {
            app.Stop()
            fmt.Fprintf(os.Stderr, "Application panic: %v\n", r)
            os.Exit(1)
        }
    }()

    return app
}

// SetRootSafely updates the root primitive safely
func SetRootSafely(app *tview.Application, root tview.Primitive, fullscreen bool) {
    app.QueueUpdateDraw(func() {
        app.SetRoot(root, fullscreen)
    })
}

// Quit gracefully shuts down the application
func Quit(app *tview.Application) {
    app.Stop()
}

// RestoreTerminal ensures terminal is restored on panic
// Call this with defer in main()
func RestoreTerminal() {
    if r := recover(); r != nil {
        // Try to restore terminal
        screen, err := tcell.NewScreen()
        if err == nil {
            screen.Fini()
        }
        fmt.Fprintf(os.Stderr, "Panic: %v\n", r)
        os.Exit(1)
    }
}
```

## Configuration Options

### Mouse Support
```go
app.EnableMouse(true)  // Enable mouse clicks and scrolling
app.EnableMouse(false) // Pure keyboard navigation
```

**Decision**: Enable by default, but provide flag to disable for pure keyboard users

### Color Support
tview automatically detects terminal color capabilities:
- 256 colors: Most modern terminals
- True color (24-bit): Windows Terminal, iTerm2, modern Linux terminals
- 16 colors: Fallback for older terminals

No manual configuration needed - tview handles it

### Alternate Screen
```go
app.Run() // Automatically uses alternate screen
```

This means:
- TUI takes over entire terminal
- Previous terminal content preserved
- On exit, terminal returns to previous state

## Error Handling

### Initialization Errors
```go
app := tview.NewApplication()
// NewApplication() doesn't return errors
// Screen initialization happens on Run()

if err := app.Run(); err != nil {
    // Handle runtime errors here
    return fmt.Errorf("application error: %w", err)
}
```

### Runtime Errors
- **Screen lost**: app.Run() returns error if screen is lost
- **Panic**: Defer recovery in main() to restore terminal
- **User quit**: Normal exit via app.Stop() returns nil

## Integration with Main

```go
// In main.go:
func launchTUI(vaultService *vault.Vault) error {
    // 1. Create app
    app := NewApp()

    // 2. Ensure terminal restoration on panic
    defer RestoreTerminal()

    // 3. Create state and components
    appState := models.NewAppState(vaultService)
    // ... component creation ...

    // 4. Create layout
    layout := createLayout(appState)

    // 5. Setup events
    setupEventHandlers(app, appState)

    // 6. Run (blocks until quit)
    if err := app.SetRoot(layout, true).Run(); err != nil {
        return fmt.Errorf("application error: %w", err)
    }

    return nil
}
```

## Best Practices

### Do:
✅ Create application once in main()
✅ Call Run() only once
✅ Use app.Stop() for clean exit
✅ Use QueueUpdateDraw() for updates from goroutines
✅ Set up panic recovery

### Don't:
❌ Call Run() multiple times
❌ Call SetRoot() before Run() multiple times
❌ Update UI from goroutines without QueueUpdateDraw()
❌ Forget to call app.Stop() on quit
❌ Create multiple application instances

## Testing Considerations

- **No headless testing**: tview requires actual terminal
- **Test initialization**: Can test NewApp() returns non-nil
- **Mock screen**: Can create mock tcell.Screen for testing (advanced)
- **Focus on logic**: Test state and components, not app itself

## Terminal Compatibility

### Tested Terminals
- ✅ Windows Terminal
- ✅ iTerm2 (macOS)
- ✅ gnome-terminal (Linux)
- ✅ Alacritty
- ✅ Kitty

### Known Issues
- Some terminals don't support mouse
- Some terminals have limited color support
- tview handles these automatically with fallbacks

## Future Enhancements

- **Session management**: Save/restore window size, position
- **Theme switching**: Allow runtime theme changes
- **Performance monitoring**: Add FPS counter for debugging
- **Debug mode**: Optional logging of all events
