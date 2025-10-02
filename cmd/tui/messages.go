package tui

import "pass-cli/internal/vault"

// Bubble Tea messages for state updates

// vaultUnlockedMsg is sent when vault is successfully unlocked
type vaultUnlockedMsg struct{}

// vaultUnlockErrorMsg is sent when vault unlock fails
type vaultUnlockErrorMsg struct {
	err error
}

// credentialsLoadedMsg is sent when credentials are loaded
type credentialsLoadedMsg struct {
	credentials []vault.CredentialMetadata
}

// credentialLoadedMsg is sent when a single credential is loaded
type credentialLoadedMsg struct {
	credential *vault.Credential
}
