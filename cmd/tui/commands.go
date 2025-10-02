package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/internal/vault"
)

// unlockVaultCmd attempts to unlock the vault using keychain, then password prompt
func unlockVaultCmd(vaultService *vault.VaultService) tea.Cmd {
	return func() tea.Msg {
		// Try keychain first
		if err := vaultService.UnlockWithKeychain(); err == nil {
			return vaultUnlockedMsg{}
		}

		// TODO: For now, just return error if keychain fails
		// In future tasks, we'll add password prompt UI
		return vaultUnlockErrorMsg{
			err: fmt.Errorf("keychain unlock failed - password prompt UI not yet implemented"),
		}
	}
}

// loadCredentialsCmd loads credentials from the vault
func loadCredentialsCmd(vaultService *vault.VaultService) tea.Cmd {
	return func() tea.Msg {
		credentials, err := vaultService.ListCredentialsWithMetadata()
		if err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		return credentialsLoadedMsg{credentials: credentials}
	}
}
