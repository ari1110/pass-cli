package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/howeyc/gopass"
	"pass-cli/cmd/tui-tview/components"
	"pass-cli/cmd/tui-tview/events"
	"pass-cli/cmd/tui-tview/layout"
	"pass-cli/cmd/tui-tview/models"
	"pass-cli/cmd/tui-tview/styles"
	"pass-cli/internal/vault"
)

const maxPasswordAttempts = 3

func main() {
	// 1. Get default vault path
	vaultPath := getDefaultVaultPath()

	// 2. Initialize vault service
	vaultService, err := vault.New(vaultPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize vault service: %v\n", err)
		os.Exit(1)
	}

	// 3. Try keychain unlock
	err = vaultService.UnlockWithKeychain()
	if err != nil {
		// Keychain unlock failed, fall back to password prompt
		fmt.Println("Keychain unlock unavailable, prompting for password...")

		// 4. Prompt for password with max attempts
		unlocked := false
		for attempt := 1; attempt <= maxPasswordAttempts; attempt++ {
			password, err := promptForPassword()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Failed to read password: %v\n", err)
				os.Exit(1)
			}

			// Try to unlock with provided password
			err = vaultService.Unlock(password)
			if err == nil {
				unlocked = true
				break
			}

			// Show error for failed attempt
			if attempt < maxPasswordAttempts {
				fmt.Fprintf(os.Stderr, "Unlock failed: %v. Please try again (%d/%d attempts).\n",
					err, attempt, maxPasswordAttempts)
			}
		}

		// Check if unlock was successful
		if !unlocked {
			fmt.Fprintf(os.Stderr, "Error: Failed to unlock vault after %d attempts.\n", maxPasswordAttempts)
			os.Exit(1)
		}
	}

	// 5. Launch TUI
	if err := launchTUI(vaultService); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}

// getDefaultVaultPath returns the default vault file path
func getDefaultVaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home not available
		return ".pass-cli/vault.enc"
	}

	return filepath.Join(home, ".pass-cli", "vault.enc")
}

// promptForPassword securely prompts user for master password
func promptForPassword() (string, error) {
	fmt.Print("Enter master password: ")

	// Use gopass for masked input
	passwordBytes, err := gopass.GetPasswdMasked()
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}

	return string(passwordBytes), nil
}

// LaunchTUI initializes and runs the TUI application
// This function is exported to be called from cmd/root.go
func LaunchTUI(vaultService *vault.VaultService) error {
	// Panic recovery to restore terminal
	defer RestoreTerminal()

	// Set rounded borders globally
	styles.SetRoundedBorders()

	// 1. Create tview.Application
	app := NewApp()

	// 2. Initialize AppState with vault service
	appState := models.NewAppState(vaultService)

	// 3. Load credentials
	if err := appState.LoadCredentials(); err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	// 4. Create UI components
	sidebar := components.NewSidebar(appState)
	table := components.NewCredentialTable(appState)
	detailView := components.NewDetailView(appState)
	statusBar := components.NewStatusBar(app, appState)

	// 5. Store components in AppState
	appState.SetSidebar(sidebar.TreeView)
	appState.SetTable(table.Table)
	appState.SetDetailView(detailView.TextView)
	appState.SetStatusBar(statusBar.TextView)

	// 6. Register callbacks
	appState.SetOnCredentialsChanged(func() {
		// Refresh all components that depend on credentials
		sidebar.Refresh()
		table.Refresh()
		detailView.Refresh()
	})

	appState.SetOnSelectionChanged(func() {
		// Refresh detail view when selection changes
		detailView.Refresh()
	})

	appState.SetOnError(func(err error) {
		// Display error in status bar
		statusBar.ShowError(err)
	})

	// 7. Create NavigationState
	nav := models.NewNavigationState(app, appState)

	// Register focus change callback to update statusbar
	nav.SetOnFocusChanged(func(focus models.FocusableComponent) {
		events.OnFocusChanged(focus, statusBar)
	})

	// 8. Create LayoutManager and build layout
	layoutMgr := layout.NewLayoutManager(app, appState)
	mainLayout := layoutMgr.CreateMainLayout()

	// 9. Create PageManager
	pageManager := layout.NewPageManager(app)

	// 10. Create EventHandler and setup shortcuts
	eventHandler := events.NewEventHandler(app, appState, nav, pageManager, statusBar, detailView)
	eventHandler.SetupGlobalShortcuts()

	// 11. Set root primitive (use pages for modal support over main layout)
	pageManager.ShowPage("main", mainLayout)
	app.SetRoot(pageManager.Pages, true)

	// 12. Run application (blocking)
	return app.Run()
}

// launchTUI is kept as a private wrapper for backward compatibility if needed
func launchTUI(vaultService *vault.VaultService) error {
	return LaunchTUI(vaultService)
}
