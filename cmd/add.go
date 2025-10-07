package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"pass-cli/internal/vault"
)

var (
	addUsername string
	addPassword string
	addNotes    string
)

var addCmd = &cobra.Command{
	Use:   "add <service>",
	Short: "Add a new credential to the vault",
	Long: `Add stores a new credential (username and password) for a service in your vault.

You will be prompted for the username and password. The password input will be
hidden for security. If you want to provide these values via flags, use:
  --username (-u) for the username
  --password (-p) for the password (not recommended for security)

The service name should be descriptive and unique (e.g., "github", "aws-prod", "db-staging").`,
	Example: `  # Add a credential with prompts
  pass-cli add github

  # Add with username flag
  pass-cli add github --username user@example.com

  # Add with notes
  pass-cli add github --notes "My GitHub account"

  # Add with all flags (password in flag not recommended)
  pass-cli add github -u user@example.com -p secret123 --notes "Work account"`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addUsername, "username", "u", "", "username for the credential")
	addCmd.Flags().StringVarP(&addPassword, "password", "p", "", "password for the credential (not recommended, use prompt instead)")
	addCmd.Flags().StringVar(&addNotes, "notes", "", "optional notes about the credential")
}

func runAdd(cmd *cobra.Command, args []string) error {
	service := args[0]

	// Validate service name
	service = strings.TrimSpace(service)
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

	// Get username if not provided
	if addUsername == "" {
		fmt.Print("Username: ")
		if _, err := fmt.Scanln(&addUsername); err != nil {
			return fmt.Errorf("failed to read username: %w", err)
		}
		addUsername = strings.TrimSpace(addUsername)
	}

	// Get password if not provided
	if addPassword == "" {
		fmt.Print("Password: ")
		password, err := readPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println() // newline after password input
		addPassword = password
	}

	// Validate password is not empty
	if addPassword == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Add credential to vault (category and URL are empty for now)
	if err := vaultService.AddCredential(service, addUsername, addPassword, "", "", addNotes); err != nil {
		return fmt.Errorf("failed to add credential: %w", err)
	}

	// Success message
	fmt.Printf("✅ Credential added successfully!\n")
	fmt.Printf("📝 Service: %s\n", service)
	if addUsername != "" {
		fmt.Printf("👤 Username: %s\n", addUsername)
	}
	if addNotes != "" {
		fmt.Printf("📋 Notes: %s\n", addNotes)
	}

	return nil
}

// unlockVault attempts to unlock the vault with keychain or prompts for password
func unlockVault(vaultService *vault.VaultService) error {
	// Try keychain first
	if err := vaultService.UnlockWithKeychain(); err == nil {
		if IsVerbose() {
			fmt.Fprintln(os.Stderr, "🔓 Unlocked vault using keychain")
		}
		return nil
	}

	// Prompt for master password
	fmt.Print("Master password: ")
	password, err := readPassword()
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // newline after password input

	if err := vaultService.Unlock(password); err != nil {
		return fmt.Errorf("failed to unlock vault: %w", err)
	}

	return nil
}