package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"

	"pass-cli/internal/vault"
)

var (
	getQuiet       bool
	getField       string
	getNoClipboard bool
	getMasked      bool
)

var getCmd = &cobra.Command{
	Use:   "get <service>",
	Short: "Retrieve a credential from the vault",
	Long: `Get retrieves a credential from your vault and copies the password to clipboard.

By default, the password is copied to the clipboard and credential details
are displayed. Use flags to customize the output:

  --quiet      Output only the requested value (for scripts)
  --field      Extract a specific field (username, password, category, url, notes, service)
  --no-clipboard  Skip copying to clipboard
  --masked     Display password as asterisks (default shows full password)

Automatic usage tracking records where credentials are accessed based on
your current working directory.`,
	Example: `  # Get credential with clipboard copy
  pass-cli get github

  # Get for scripts (outputs only password)
  pass-cli get github --quiet

  # Get specific field for scripts
  pass-cli get github --field username --quiet

  # Get without clipboard
  pass-cli get github --no-clipboard

  # Get with masked password display
  pass-cli get github --masked`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVarP(&getQuiet, "quiet", "q", false, "output only the requested value (script-friendly)")
	getCmd.Flags().StringVarP(&getField, "field", "f", "password", "field to extract (username, password, category, url, notes, service)")
	getCmd.Flags().BoolVar(&getNoClipboard, "no-clipboard", false, "do not copy to clipboard")
	getCmd.Flags().BoolVar(&getMasked, "masked", false, "display password as asterisks")
}

func runGet(cmd *cobra.Command, args []string) error {
	service := strings.TrimSpace(args[0])
	if service == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	vaultPath := GetVaultPath()

	// Check if vault exists
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return fmt.Errorf("vault not found at %s\nRun 'pass-cli init' to create a vault first", vaultPath)
	}

	// Create vault service
	vaultService, err := vault.New(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to create vault service: %w", err)
	}

	// Unlock vault
	if err := unlockVault(vaultService); err != nil {
		return err
	}
	defer vaultService.Lock()

	// Get credential with usage tracking (unless quiet mode for scripts)
	trackUsage := !getQuiet
	cred, err := vaultService.GetCredential(service, trackUsage)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	// Quiet mode - output only requested field
	if getQuiet {
		return outputQuietMode(cred)
	}

	// Normal mode - display credential details
	return outputNormalMode(cred)
}

func outputQuietMode(cred *vault.Credential) error {
	field := strings.ToLower(getField)
	var value string

	switch field {
	case "username", "user", "u":
		value = cred.Username
	case "password", "pass", "p":
		value = cred.Password
	case "category", "cat", "c":
		value = cred.Category
	case "url":
		value = cred.URL
	case "notes", "note", "n":
		value = cred.Notes
	case "service", "s":
		value = cred.Service
	default:
		return fmt.Errorf("invalid field: %s (valid: username, password, category, url, notes, service)", getField)
	}

	fmt.Println(value)
	return nil
}

func outputNormalMode(cred *vault.Credential) error {
	// Display credential details
	fmt.Printf("📝 Service: %s\n", cred.Service)

	if cred.Username != "" {
		fmt.Printf("👤 Username: %s\n", cred.Username)
	}

	// Display password (masked or full)
	if getMasked {
		fmt.Printf("🔑 Password: %s\n", strings.Repeat("*", len(cred.Password)))
	} else {
		fmt.Printf("🔑 Password: %s\n", cred.Password)
	}

	if cred.Category != "" {
		fmt.Printf("🏷️  Category: %s\n", cred.Category)
	}

	if cred.URL != "" {
		fmt.Printf("🔗 URL: %s\n", cred.URL)
	}

	if cred.Notes != "" {
		fmt.Printf("📋 Notes: %s\n", cred.Notes)
	}

	// Display timestamps
	fmt.Printf("📅 Created: %s\n", cred.CreatedAt.Format("2006-01-02 15:04:05"))
	if !cred.UpdatedAt.Equal(cred.CreatedAt) {
		fmt.Printf("📅 Updated: %s\n", cred.UpdatedAt.Format("2006-01-02 15:04:05"))
	}

	// Copy to clipboard unless disabled
	if !getNoClipboard {
		if err := clipboard.WriteAll(cred.Password); err != nil {
			fmt.Fprintf(os.Stderr, "\n⚠️  Warning: failed to copy to clipboard: %v\n", err)
		} else {
			fmt.Println("\n✅ Password copied to clipboard!")

			// Schedule clipboard clear in background (30 seconds)
			go func() {
				time.Sleep(30 * time.Second)
				// Only clear if the clipboard still contains our password
				if current, err := clipboard.ReadAll(); err == nil && current == cred.Password {
					_ = clipboard.WriteAll("")
					if IsVerbose() {
						fmt.Fprintln(os.Stderr, "🧹 Clipboard cleared")
					}
				}
			}()
		}
	}

	return nil
}
