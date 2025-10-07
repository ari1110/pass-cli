package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"pass-cli/internal/vault"
)

var (
	updateUsername string
	updatePassword string
	updateNotes    string
	updateForce    bool
)

var updateCmd = &cobra.Command{
	Use:   "update <service>",
	Short: "Update an existing credential",
	Long: `Update modifies an existing credential in your vault.

You can selectively update individual fields (username, password, notes) without
affecting the others. Empty values mean "don't change".

By default, you'll see a usage warning if the credential has been accessed before,
showing where and when it was last used. Use --force to skip the confirmation.`,
	Example: `  # Update password only (interactive prompt)
  pass-cli update github

  # Update username only
  pass-cli update github --username new-user@example.com

  # Update password only
  pass-cli update github --password newpass123

  # Update notes
  pass-cli update github --notes "Updated account"

  # Update multiple fields
  pass-cli update github -u user -p pass --notes "New info"

  # Skip confirmation
  pass-cli update github --force`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&updateUsername, "username", "u", "", "new username")
	updateCmd.Flags().StringVarP(&updatePassword, "password", "p", "", "new password")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "new notes")
	updateCmd.Flags().BoolVar(&updateForce, "force", false, "skip confirmation prompt")
}

func runUpdate(cmd *cobra.Command, args []string) error {
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

	// Check if credential exists
	cred, err := vaultService.GetCredential(service, false)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	// If no flags provided, prompt for what to update
	if updateUsername == "" && updatePassword == "" && updateNotes == "" {
		fmt.Println("What would you like to update? (leave empty to keep current value)")
		fmt.Println()

		// Prompt for username
		fmt.Printf("Username [%s]: ", cred.Username)
		var username string
		_, _ = fmt.Scanln(&username)
		updateUsername = strings.TrimSpace(username)

		// Prompt for password
		fmt.Print("Password (hidden): ")
		password, err := readPassword()
		if err != nil {
			return fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println()
		updatePassword = password

		// Prompt for notes
		fmt.Printf("Notes [%s]: ", cred.Notes)
		var notes string
		_, _ = fmt.Scanln(&notes)
		updateNotes = strings.TrimSpace(notes)
	}

	// Check if anything is being updated
	if updateUsername == "" && updatePassword == "" && updateNotes == "" {
		fmt.Println("No changes specified.")
		return nil
	}

	// Show usage warning if credential has been accessed
	stats, _ := vaultService.GetUsageStats(service)
	if len(stats) > 0 && !updateForce {
		fmt.Println("\n⚠️  Usage Warning:")

		totalCount := 0
		var lastAccessed string
		for _, record := range stats {
			totalCount += record.Count
			if lastAccessed == "" || record.Timestamp.After(cred.UpdatedAt) {
				lastAccessed = formatRelativeTime(record.Timestamp)
			}
		}

		fmt.Printf("   Used in %d location(s), last used %s\n", len(stats), lastAccessed)
		fmt.Printf("   Total access count: %d\n\n", totalCount)

		// Ask for confirmation
		fmt.Print("Continue with update? (y/N): ")
		var confirm string
		_, _ = fmt.Scanln(&confirm)
		confirm = strings.ToLower(strings.TrimSpace(confirm))

		if confirm != "y" && confirm != "yes" {
			fmt.Println("Update cancelled.")
			return nil
		}
	}

	// Perform update using UpdateOpts (only update non-empty fields)
	opts := vault.UpdateOpts{}
	if updateUsername != "" {
		opts.Username = &updateUsername
	}
	if updatePassword != "" {
		opts.Password = &updatePassword
	}
	if updateNotes != "" {
		opts.Notes = &updateNotes
	}

	if err := vaultService.UpdateCredential(service, opts); err != nil {
		return fmt.Errorf("failed to update credential: %w", err)
	}

	// Success message
	fmt.Printf("✅ Credential updated successfully!\n")
	fmt.Printf("📝 Service: %s\n", service)

	if updateUsername != "" {
		fmt.Printf("👤 New username: %s\n", updateUsername)
	}
	if updatePassword != "" {
		fmt.Printf("🔑 Password updated\n")
	}
	if updateNotes != "" {
		fmt.Printf("📋 New notes: %s\n", updateNotes)
	}

	return nil
}