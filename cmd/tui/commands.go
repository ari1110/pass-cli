package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/internal/vault"
)

// unlockVaultCmd attempts to unlock the vault using keychain, then falls back to password prompt
func unlockVaultCmd(vaultService *vault.VaultService) tea.Cmd {
	return func() tea.Msg {
		// Try keychain first
		if err := vaultService.UnlockWithKeychain(); err == nil {
			return vaultUnlockedMsg{}
		}

		// Keychain failed, request password from user
		return needPasswordMsg{}
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

// loadCredentialDetailsCmd loads a single credential's full details
func loadCredentialDetailsCmd(vaultService *vault.VaultService, service string) tea.Cmd {
	return func() tea.Msg {
		// Don't track usage when just viewing details in TUI
		credential, err := vaultService.GetCredential(service, false)
		if err != nil {
			return vaultUnlockErrorMsg{err: err}
		}

		return credentialLoadedMsg{credential: credential}
	}
}
