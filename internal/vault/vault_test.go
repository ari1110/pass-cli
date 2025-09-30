package vault

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestVault(t *testing.T) (*VaultService, string, func()) {
	t.Helper()

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "vault-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	vaultPath := filepath.Join(tempDir, "test.vault")
	vault, err := New(vaultPath)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		t.Fatalf("Failed to create vault service: %v", err)
	}

	cleanup := func() {
		vault.Lock()
		_ = os.RemoveAll(tempDir)
	}

	return vault, vaultPath, cleanup
}

func TestNew(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	if vault == nil {
		t.Fatal("New() returned nil")
	}

	if vault.IsUnlocked() {
		t.Error("New vault should be locked")
	}
}

func TestInitialize(t *testing.T) {
	vault, vaultPath, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize vault
	err := vault.Initialize(password, false)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Verify vault file was created
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		t.Error("Vault file was not created")
	}
}

func TestInitializeWithShortPassword(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	// Try with password < 8 characters
	err := vault.Initialize("short", false)
	if err == nil {
		t.Error("Initialize() should fail with short password")
	}
}

func TestInitializeExistingVault(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize once
	err := vault.Initialize(password, false)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Try to initialize again
	err = vault.Initialize(password, false)
	if err == nil {
		t.Error("Initialize() should fail on existing vault")
	}
}

func TestUnlock(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize and unlock
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	if !vault.IsUnlocked() {
		t.Error("Vault should be unlocked")
	}
}

func TestUnlockWithWrongPassword(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Try to unlock with wrong password
	err := vault.Unlock("wrong-password")
	if err == nil {
		t.Error("Unlock() should fail with wrong password")
	}

	if vault.IsUnlocked() {
		t.Error("Vault should not be unlocked")
	}
}

func TestLock(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize and unlock
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Lock
	vault.Lock()

	if vault.IsUnlocked() {
		t.Error("Vault should be locked")
	}
}

func TestAddCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize and unlock
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Add credential
	err := vault.AddCredential("github", "user@example.com", "secret123", "My GitHub account")
	if err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Verify it was added
	services, err := vault.ListCredentials()
	if err != nil {
		t.Fatalf("ListCredentials() failed: %v", err)
	}

	if len(services) != 1 || services[0] != "github" {
		t.Errorf("Expected [github], got %v", services)
	}
}

func TestAddCredentialWhenLocked(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize but don't unlock
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Try to add credential
	err := vault.AddCredential("github", "user@example.com", "secret123", "")
	if err != ErrVaultLocked {
		t.Errorf("AddCredential() error = %v, want %v", err, ErrVaultLocked)
	}
}

func TestAddDuplicateCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize and unlock
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Add credential
	if err := vault.AddCredential("github", "user", "pass", ""); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Try to add duplicate
	err := vault.AddCredential("github", "user2", "pass2", "")
	if err == nil {
		t.Error("AddCredential() should return error for duplicate")
	}
}

func TestGetCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("github", "user@example.com", "secret123", "My GitHub"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Get credential (without usage tracking for this test)
	cred, err := vault.GetCredential("github", false)
	if err != nil {
		t.Fatalf("GetCredential() failed: %v", err)
	}

	if cred.Service != "github" {
		t.Errorf("Service = %s, want github", cred.Service)
	}
	if cred.Username != "user@example.com" {
		t.Errorf("Username = %s, want user@example.com", cred.Username)
	}
	if cred.Password != "secret123" {
		t.Errorf("Password = %s, want secret123", cred.Password)
	}
	if cred.Notes != "My GitHub" {
		t.Errorf("Notes = %s, want My GitHub", cred.Notes)
	}
}

func TestGetCredentialWithUsageTracking(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("github", "user", "pass", ""); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Get credential with usage tracking
	_, err := vault.GetCredential("github", true)
	if err != nil {
		t.Fatalf("GetCredential() failed: %v", err)
	}

	// Check usage stats
	stats, err := vault.GetUsageStats("github")
	if err != nil {
		t.Fatalf("GetUsageStats() failed: %v", err)
	}

	if len(stats) == 0 {
		t.Error("Expected usage record, got none")
	}

	// Access again to increment count
	_, err = vault.GetCredential("github", true)
	if err != nil {
		t.Fatalf("Second GetCredential() failed: %v", err)
	}

	stats, err = vault.GetUsageStats("github")
	if err != nil {
		t.Fatalf("GetUsageStats() failed: %v", err)
	}

	// Should have count of 2 now
	for _, record := range stats {
		if record.Count != 2 {
			t.Errorf("Usage count = %d, want 2", record.Count)
		}
	}
}

func TestUpdateCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("github", "old-user", "old-pass", "old notes"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Wait a moment to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	// Update credential
	err := vault.UpdateCredential("github", "new-user", "new-pass", "new notes")
	if err != nil {
		t.Fatalf("UpdateCredential() failed: %v", err)
	}

	// Verify update
	cred, err := vault.GetCredential("github", false)
	if err != nil {
		t.Fatalf("GetCredential() failed: %v", err)
	}

	if cred.Username != "new-user" {
		t.Errorf("Username = %s, want new-user", cred.Username)
	}
	if cred.Password != "new-pass" {
		t.Errorf("Password = %s, want new-pass", cred.Password)
	}
	if cred.Notes != "new notes" {
		t.Errorf("Notes = %s, want new notes", cred.Notes)
	}

	// Verify UpdatedAt was changed
	if !cred.UpdatedAt.After(cred.CreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt")
	}
}

func TestDeleteCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("github", "user", "pass", ""); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Delete credential
	err := vault.DeleteCredential("github")
	if err != nil {
		t.Fatalf("DeleteCredential() failed: %v", err)
	}

	// Verify it's gone
	services, err := vault.ListCredentials()
	if err != nil {
		t.Fatalf("ListCredentials() failed: %v", err)
	}

	if len(services) != 0 {
		t.Errorf("Expected empty list, got %v", services)
	}
}

func TestDeleteNonExistentCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize and unlock
	if err := vault.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Try to delete non-existent credential
	err := vault.DeleteCredential("nonexistent")
	if err == nil {
		t.Error("DeleteCredential() should return error for non-existent credential")
	}
}

func TestPersistence(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "vault-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	vaultPath := filepath.Join(tempDir, "test.vault")
	password := "test-password-12345"

	// Create first vault instance
	vault1, err := New(vaultPath)
	if err != nil {
		t.Fatalf("Failed to create vault1: %v", err)
	}

	// Initialize and add credential
	if err := vault1.Initialize(password, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault1.Unlock(password); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault1.AddCredential("github", "user", "pass", "notes"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}
	vault1.Lock()

	// Create second vault instance pointing to same file
	vault2, err := New(vaultPath)
	if err != nil {
		t.Fatalf("Failed to create vault2: %v", err)
	}

	// Unlock and verify credential exists
	if err := vault2.Unlock(password); err != nil {
		t.Fatalf("Unlock() vault2 failed: %v", err)
	}

	cred, err := vault2.GetCredential("github", false)
	if err != nil {
		t.Fatalf("GetCredential() from vault2 failed: %v", err)
	}

	if cred.Username != "user" || cred.Password != "pass" {
		t.Error("Credential data not persisted correctly")
	}
}

func TestChangePassword(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	oldPassword := "old-password-12345"
	newPassword := "new-password-67890"

	// Initialize and unlock
	if err := vault.Initialize(oldPassword, false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock(oldPassword); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Add a credential
	if err := vault.AddCredential("test", "user", "pass", ""); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Change password
	if err := vault.ChangePassword(newPassword); err != nil {
		t.Fatalf("ChangePassword() failed: %v", err)
	}

	// Lock and try to unlock with old password (should fail)
	vault.Lock()
	err := vault.Unlock(oldPassword)
	if err == nil {
		t.Error("Should not be able to unlock with old password")
	}

	// Unlock with new password (should succeed)
	if err := vault.Unlock(newPassword); err != nil {
		t.Fatalf("Failed to unlock with new password: %v", err)
	}

	// Verify credential still exists
	cred, err := vault.GetCredential("test", false)
	if err != nil {
		t.Fatalf("GetCredential() failed after password change: %v", err)
	}
	if cred.Username != "user" {
		t.Error("Credential data corrupted after password change")
	}
}