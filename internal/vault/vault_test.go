package vault

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"pass-cli/internal/storage"
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

func setupTestVaultWithStorage(t *testing.T) (*VaultService, *storage.StorageService, func()) {
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

	return vault, vault.storageService, cleanup
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
	err := vault.Initialize([]byte(password), false)
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
	err := vault.Initialize([]byte("short"), false)
	if err == nil {
		t.Error("Initialize() should fail with short password")
	}
}

func TestInitializeExistingVault(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize once
	err := vault.Initialize([]byte(password), false)
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Try to initialize again
	err = vault.Initialize([]byte(password), false)
	if err == nil {
		t.Error("Initialize() should fail on existing vault")
	}
}

func TestUnlock(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize and unlock
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	if err := vault.Unlock([]byte(password)); err != nil {
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
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Try to unlock with wrong password
	err := vault.Unlock([]byte("wrong-password"))
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
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
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
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Add credential
	err := vault.AddCredential("github", "user@example.com", []byte("secret123"), "Work", "https://github.com", "My GitHub account")
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
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Try to add credential
	err := vault.AddCredential("github", "user@example.com", []byte("secret123"), "", "", "")
	if err != ErrVaultLocked {
		t.Errorf("AddCredential() error = %v, want %v", err, ErrVaultLocked)
	}
}

func TestAddDuplicateCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize and unlock
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Add credential
	if err := vault.AddCredential("github", "user", []byte("pass"), "", "", ""); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Try to add duplicate
	err := vault.AddCredential("github", "user2", []byte("pass2"), "", "", "")
	if err == nil {
		t.Error("AddCredential() should return error for duplicate")
	}
}

func TestGetCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("github", "user@example.com", []byte("secret123"), "Personal", "https://github.com", "My GitHub"); err != nil {
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
	// T020d: Compare []byte using bytes.Equal
	if !bytes.Equal(cred.Password, []byte("secret123")) {
		t.Errorf("Password = %s, want secret123", cred.Password)
	}
	if cred.Category != "Personal" {
		t.Errorf("Category = %s, want Personal", cred.Category)
	}
	if cred.URL != "https://github.com" {
		t.Errorf("URL = %s, want https://github.com", cred.URL)
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
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("github", "user", []byte("pass"), "", "", ""); err != nil {
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
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("github", "old-user", []byte("old-pass"), "Old Category", "https://old.com", "old notes"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Wait a moment to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	// Update credential using UpdateOpts
	newUser := "new-user"
	newPass := []byte("new-pass") // T020d: []byte for password
	newCategory := "New Category"
	newURL := "https://new.com"
	newNotes := "new notes"

	err := vault.UpdateCredential("github", UpdateOpts{
		Username: &newUser,
		Password: &newPass,
		Category: &newCategory,
		URL:      &newURL,
		Notes:    &newNotes,
	})
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
	// T020d: Convert []byte to string for comparison
	if string(cred.Password) != "new-pass" {
		t.Errorf("Password = %s, want new-pass", string(cred.Password))
	}
	if cred.Category != "New Category" {
		t.Errorf("Category = %s, want New Category", cred.Category)
	}
	if cred.URL != "https://new.com" {
		t.Errorf("URL = %s, want https://new.com", cred.URL)
	}
	if cred.Notes != "new notes" {
		t.Errorf("Notes = %s, want new notes", cred.Notes)
	}

	// Verify UpdatedAt was changed
	if !cred.UpdatedAt.After(cred.CreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt")
	}
}

func TestUpdateCredentialClearFields(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential with category and URL
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("test-service", "user", []byte("pass"), "Work", "https://example.com", "notes"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Clear category and URL by passing pointers to empty strings
	emptyCategory := ""
	emptyURL := ""

	err := vault.UpdateCredential("test-service", UpdateOpts{
		Category: &emptyCategory,
		URL:      &emptyURL,
	})
	if err != nil {
		t.Fatalf("UpdateCredential() failed: %v", err)
	}

	// Verify category and URL are cleared
	cred, err := vault.GetCredential("test-service", false)
	if err != nil {
		t.Fatalf("GetCredential() failed: %v", err)
	}

	if cred.Category != "" {
		t.Errorf("Category = %s, want empty string", cred.Category)
	}
	if cred.URL != "" {
		t.Errorf("URL = %s, want empty string", cred.URL)
	}
	// Verify other fields were not changed
	if cred.Username != "user" {
		t.Errorf("Username = %s, want user", cred.Username)
	}
	if cred.Notes != "notes" {
		t.Errorf("Notes = %s, want notes", cred.Notes)
	}
}

func TestUpdateCredentialPartial(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("test-service", "old-user", []byte("old-pass"), "Old Category", "https://old.com", "old notes"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Update only password (leave other fields unchanged)
	newPassword := []byte("new-pass") // T020d: []byte for password
	err := vault.UpdateCredential("test-service", UpdateOpts{
		Password: &newPassword,
	})
	if err != nil {
		t.Fatalf("UpdateCredential() failed: %v", err)
	}

	// Verify only password changed
	cred, err := vault.GetCredential("test-service", false)
	if err != nil {
		t.Fatalf("GetCredential() failed: %v", err)
	}

	// T020d: Convert []byte to string for comparison
	if string(cred.Password) != "new-pass" {
		t.Errorf("Password = %s, want new-pass", string(cred.Password))
	}
	// Verify other fields remain unchanged
	if cred.Username != "old-user" {
		t.Errorf("Username = %s, want old-user", cred.Username)
	}
	if cred.Category != "Old Category" {
		t.Errorf("Category = %s, want Old Category", cred.Category)
	}
	if cred.URL != "https://old.com" {
		t.Errorf("URL = %s, want https://old.com", cred.URL)
	}
	if cred.Notes != "old notes" {
		t.Errorf("Notes = %s, want old notes", cred.Notes)
	}
}

func TestUpdateCredentialFields(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("test-service", "old-user", []byte("old-pass"), "Old", "https://old.com", "old"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Update using convenience wrapper (6-parameter API)
	err := vault.UpdateCredentialFields("test-service", "new-user", "new-pass", "New", "https://new.com", "new")
	if err != nil {
		t.Fatalf("UpdateCredentialFields() failed: %v", err)
	}

	// Verify all fields updated
	cred, err := vault.GetCredential("test-service", false)
	if err != nil {
		t.Fatalf("GetCredential() failed: %v", err)
	}

	if cred.Username != "new-user" {
		t.Errorf("Username = %s, want new-user", cred.Username)
	}
	// T020d: Convert []byte to string for comparison
	if string(cred.Password) != "new-pass" {
		t.Errorf("Password = %s, want new-pass", string(cred.Password))
	}
	if cred.Category != "New" {
		t.Errorf("Category = %s, want New", cred.Category)
	}
	if cred.URL != "https://new.com" {
		t.Errorf("URL = %s, want https://new.com", cred.URL)
	}
	if cred.Notes != "new" {
		t.Errorf("Notes = %s, want new", cred.Notes)
	}
}

func TestUpdateCredentialFieldsEmptyMeansNoChange(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("test-service", "old-user", []byte("old-pass"), "Old", "https://old.com", "old"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Update only password using convenience wrapper (empty strings = no change)
	err := vault.UpdateCredentialFields("test-service", "", "new-pass", "", "", "")
	if err != nil {
		t.Fatalf("UpdateCredentialFields() failed: %v", err)
	}

	// Verify only password changed, others remain
	cred, err := vault.GetCredential("test-service", false)
	if err != nil {
		t.Fatalf("GetCredential() failed: %v", err)
	}

	// T020d: Convert []byte to string for comparison
	if string(cred.Password) != "new-pass" {
		t.Errorf("Password = %s, want new-pass", string(cred.Password))
	}
	// Verify others unchanged
	if cred.Username != "old-user" {
		t.Errorf("Username = %s, want old-user", cred.Username)
	}
	if cred.Category != "Old" {
		t.Errorf("Category = %s, want Old", cred.Category)
	}
	if cred.URL != "https://old.com" {
		t.Errorf("URL = %s, want https://old.com", cred.URL)
	}
	if cred.Notes != "old" {
		t.Errorf("Notes = %s, want old", cred.Notes)
	}
}

func TestListCredentialsWithMetadataIncludesCategoryAndURL(t *testing.T) {
	v, _, cleanup := setupTestVault(t)
	defer cleanup()

	pw := "test-password-12345"
	// T020d: Convert to []byte
	if err := v.Initialize([]byte(pw), false); err != nil {
		t.Fatal(err)
	}
	if err := v.Unlock([]byte(pw)); err != nil {
		t.Fatal(err)
	}
	if err := v.AddCredential("svc", "user", []byte("pass"), "Work", "https://ex", "notes"); err != nil {
		t.Fatal(err)
	}

	metas, err := v.ListCredentialsWithMetadata()
	if err != nil {
		t.Fatal(err)
	}
	if len(metas) != 1 {
		t.Fatalf("want 1 meta, got %d", len(metas))
	}
	m := metas[0]
	if m.Category != "Work" {
		t.Errorf("Category=%q, want Work", m.Category)
	}
	if m.URL != "https://ex" {
		t.Errorf("URL=%q, want https://ex", m.URL)
	}
}

func TestDeleteCredential(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize, unlock, and add credential
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault.AddCredential("github", "user", []byte("pass"), "", "", ""); err != nil {
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
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
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
	// T020d: Convert to []byte
	if err := vault1.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault1.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}
	if err := vault1.AddCredential("github", "user", []byte("pass"), "Test", "https://test.com", "notes"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}
	vault1.Lock()

	// Create second vault instance pointing to same file
	vault2, err := New(vaultPath)
	if err != nil {
		t.Fatalf("Failed to create vault2: %v", err)
	}

	// Unlock and verify credential exists
	if err := vault2.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() vault2 failed: %v", err)
	}

	cred, err := vault2.GetCredential("github", false)
	if err != nil {
		t.Fatalf("GetCredential() from vault2 failed: %v", err)
	}

	// T020d: Compare []byte using bytes.Equal
	if cred.Username != "user" || !bytes.Equal(cred.Password, []byte("pass")) {
		t.Error("Credential data not persisted correctly")
	}
	if cred.Category != "Test" || cred.URL != "https://test.com" {
		t.Error("Category and URL not persisted correctly")
	}
}

func TestChangePassword(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	oldPassword := "old-password-12345"
	newPassword := "new-password-67890"

	// Initialize and unlock
	// T020d: Convert to []byte
	if err := vault.Initialize([]byte(oldPassword), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(oldPassword)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Add a credential
	if err := vault.AddCredential("test", "user", []byte("pass"), "", "", ""); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Change password
	// T020d: Convert to []byte
	if err := vault.ChangePassword([]byte(newPassword)); err != nil {
		t.Fatalf("ChangePassword() failed: %v", err)
	}

	// Lock and try to unlock with old password (should fail)
	vault.Lock()
	// T020d: Convert to []byte
	err := vault.Unlock([]byte(oldPassword))
	if err == nil {
		t.Error("Should not be able to unlock with old password")
	}

	// Unlock with new password (should succeed)
	// T020d: Convert to []byte
	if err := vault.Unlock([]byte(newPassword)); err != nil {
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

func TestBackwardCompatibility(t *testing.T) {
	vault, _, cleanup := setupTestVault(t)
	defer cleanup()

	password := "test-password-12345"

	// Initialize and unlock
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Add credential without category and URL (empty strings)
	if err := vault.AddCredential("legacy-service", "user", []byte("pass"), "", "", "notes"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Lock and unlock to force serialization/deserialization
	vault.Lock()
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Verify credential loads with empty string defaults
	cred, err := vault.GetCredential("legacy-service", false)
	if err != nil {
		t.Fatalf("GetCredential() failed: %v", err)
	}

	if cred.Category != "" {
		t.Errorf("Category should default to empty string, got %s", cred.Category)
	}
	if cred.URL != "" {
		t.Errorf("URL should default to empty string, got %s", cred.URL)
	}
	if cred.Username != "user" {
		t.Errorf("Username = %s, want user", cred.Username)
	}
	if cred.Notes != "notes" {
		t.Errorf("Notes = %s, want notes", cred.Notes)
	}
}

// T023 [US2]: Test automatic migration from 100k to 600k iterations on password change
// FR-010: System MUST automatically upgrade legacy vaults to 600k iterations
func TestIterationsMigrationOnPasswordChange(t *testing.T) {
	// T023/T036: Test iteration count migration from 100k to 600k during password change
	vault, storageService, cleanup := setupTestVaultWithStorage(t)
	defer cleanup()

	password := "test-password-12345"
	newPassword := "new-password-67890"

	// Initialize vault (will use 600k iterations by default)
	if err := vault.Initialize([]byte(password), false); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Add a credential to ensure data survives migration
	if err := vault.AddCredential("test", "user", []byte("pass"), "", "", "test migration"); err != nil {
		t.Fatalf("AddCredential() failed: %v", err)
	}

	// Lock and manually downgrade iterations to simulate legacy vault
	vault.Lock()

	// Simulate legacy vault by saving with 100k iterations directly
	if err := vault.Unlock([]byte(password)); err != nil {
		t.Fatalf("Unlock() failed: %v", err)
	}

	// Save with downgraded iterations to simulate legacy vault
	data, err := json.Marshal(vault.vaultData)
	if err != nil {
		t.Fatalf("Failed to marshal vault data: %v", err)
	}

	if err := storageService.SaveVaultWithIterationsUnsafe(data, password, 100000); err != nil {
		t.Fatalf("Failed to save legacy vault: %v", err)
	}

	// Verify starting with 100k iterations
	currentIterations := storageService.GetIterations()
	if currentIterations != 100000 {
		t.Fatalf("Expected initial iterations 100000, got %d", currentIterations)
	}

	// Change password - should trigger migration to 600k iterations (T033)
	if err := vault.ChangePassword([]byte(newPassword)); err != nil {
		t.Fatalf("ChangePassword() failed: %v", err)
	}

	// Verify iterations were upgraded to 600k
	newIterations := storageService.GetIterations()
	if newIterations != 600000 {
		t.Errorf("Expected iterations upgraded to 600000, got %d", newIterations)
	}

	// Lock and unlock with new password to verify migration worked
	vault.Lock()
	if err := vault.Unlock([]byte(newPassword)); err != nil {
		t.Fatalf("Unlock() with new password failed: %v", err)
	}

	// Verify credential still exists after migration
	cred, err := vault.GetCredential("test", false)
	if err != nil {
		t.Fatalf("GetCredential() failed after migration: %v", err)
	}

	if cred.Username != "user" {
		t.Errorf("Username = %s, want user", cred.Username)
	}
	if string(cred.Password) != "pass" {
		t.Errorf("Password = %s, want pass", string(cred.Password))
	}
	if cred.Notes != "test migration" {
		t.Errorf("Notes = %s, want 'test migration'", cred.Notes)
	}

	t.Log("Migration from 100k to 600k iterations successful")
}
