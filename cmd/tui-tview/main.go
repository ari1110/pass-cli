package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/howeyc/gopass"
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

// launchTUI initializes and runs the TUI application
// This is a stub that will be implemented in subsequent tasks
func launchTUI(vaultService *vault.VaultService) error {
	// TODO: Implement TUI initialization in next tasks
	// This will include:
	// 1. Create tview.Application
	// 2. Initialize AppState with vault service
	// 3. Load credentials
	// 4. Create UI components
	// 5. Build layout
	// 6. Setup keyboard shortcuts
	// 7. Run application (blocking)

	fmt.Println("TUI launch successful! (stub implementation)")
	fmt.Println("Press Ctrl+C to exit")

	// For now, just return success
	// The actual implementation will call app.Run() which blocks until quit
	return nil
}
