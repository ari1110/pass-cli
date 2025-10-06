# main.go

## Purpose
Entry point for the tview-based TUI. Handles initialization, vault unlocking, and launching the main application.

## Responsibilities

1. **Vault service initialization**: Create vault service instance
2. **Keychain check**: Attempt automatic unlock via system keychain
3. **Password prompt fallback**: If keychain fails, prompt for master password
4. **Application launch**: Once unlocked, start the TUI application
5. **Error handling**: Display errors and exit gracefully

## Dependencies

### Internal Packages
- `pass-cli/internal/vault` - Vault service for credential management
- `pass-cli/internal/keychain` - System keychain integration
- `pass-cli/cmd/tui-tview/models` - Application state
- `pass-cli/cmd/tui-tview/layout` - Layout manager
- `pass-cli/cmd/tui-tview/events` - Event handlers
- `pass-cli/cmd/tui-tview/components` - UI components

### External Dependencies
- `github.com/rivo/tview` - TUI framework
- `github.com/gdamore/tcell/v2` - Terminal library

## Flow

```
main()
  ↓
Initialize vault service
  ↓
Check for keychain (auto-unlock)
  ↓
  ├─ Success → Load credentials → Launch app
  ↓
  └─ Failure → Prompt for password → Unlock → Launch app
  ↓
Setup tview.Application
  ↓
Create AppState (models/state.go)
  ↓
Create UI components (components/*.go)
  ↓
Create layout (layout/manager.go)
  ↓
Setup event handlers (events/handlers.go)
  ↓
Run application (blocking)
  ↓
Exit
```

## Key Functions

### `main()`
**Purpose**: Entry point, orchestrates initialization sequence

**Steps**:
1. Create vault service instance
2. Attempt keychain unlock
3. If keychain fails, prompt for password
4. Call `launchTUI()` with unlocked vault
5. Handle errors and exit codes

### `launchTUI(vaultService *vault.Vault)`
**Purpose**: Initialize and run the TUI application

**Steps**:
1. Create tview.Application
2. Initialize AppState with vault service
3. Load initial credentials
4. Create all UI components
5. Build layout
6. Setup keyboard shortcuts
7. Set root primitive
8. Run application (blocking)

### `promptForPassword() (string, error)`
**Purpose**: Securely prompt user for master password

**Uses**:
- `howeyc/gopass` or similar for masked input
- Returns password string or error

## Error Handling

- **Vault service creation failure**: Exit with code 1
- **Unlock failure**: Display error, prompt again (max 3 attempts)
- **Credential loading failure**: Show error modal, allow retry
- **TUI initialization failure**: Exit with error message

## Example Structure

```go
package main

import (
    "fmt"
    "os"

    "github.com/rivo/tview"
    "pass-cli/internal/vault"
    "pass-cli/cmd/tui-tview/models"
    "pass-cli/cmd/tui-tview/components"
    "pass-cli/cmd/tui-tview/layout"
    "pass-cli/cmd/tui-tview/events"
)

func main() {
    // 1. Initialize vault service
    vaultService, err := vault.NewVault()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // 2. Try keychain unlock
    err = vaultService.UnlockWithKeychain()
    if err != nil {
        // 3. Fallback to password prompt
        password, err := promptForPassword()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: %v\n", err)
            os.Exit(1)
        }

        err = vaultService.Unlock(password)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Unlock failed: %v\n", err)
            os.Exit(1)
        }
    }

    // 4. Launch TUI
    if err := launchTUI(vaultService); err != nil {
        fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
        os.Exit(1)
    }
}

func launchTUI(vaultService *vault.Vault) error {
    // Create tview application
    app := tview.NewApplication()

    // Create application state
    appState := models.NewAppState(vaultService)

    // Load credentials
    if err := appState.LoadCredentials(); err != nil {
        return fmt.Errorf("failed to load credentials: %w", err)
    }

    // Create UI components
    // (components are created and stored in appState)

    // Create layout manager
    layoutMgr := layout.NewLayoutManager(app, appState)
    mainLayout := layoutMgr.CreateMainLayout()

    // Setup event handlers
    eventHandler := events.NewEventHandler(app, appState)
    eventHandler.SetupGlobalShortcuts()

    // Run application
    return app.SetRoot(mainLayout, true).Run()
}

func promptForPassword() (string, error) {
    fmt.Print("Enter master password: ")
    // Use howeyc/gopass or similar for masked input
    password, err := gopass.GetPasswd()
    return string(password), err
}
```

## Notes

- **No event loop before Run()**: All setup must happen before `app.Run()` is called
- **Run() blocks**: The application runs until user quits (SetInputCapture handles 'q')
- **Clean exit**: app.Stop() gracefully shuts down the application
- **Alt screen**: Run with `app.Run()` automatically uses alternate screen (no flag needed)

## Testing Considerations

- **Mock vault service**: For testing, inject a mock vault with test credentials
- **Headless mode**: Can't easily test TUI rendering, focus on state and logic
- **Unit test components**: Test component creation and refresh logic separately
- **Integration test flow**: Test that unlock → launch → credential display works

## Future Enhancements

- **Session restore**: Remember last viewed credential/category
- **Multi-vault support**: Allow switching between different vault files
- **Theme selection**: Allow user to choose color scheme
- **Help screen on first launch**: Show keyboard shortcuts for new users
