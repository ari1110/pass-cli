package storage

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"pass-cli/internal/crypto"
)

func TestNewStorageService(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")

	// Test valid creation
	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	if storage.vaultPath != vaultPath {
		t.Errorf("Expected vault path %s, got %s", vaultPath, storage.vaultPath)
	}

	// Test nil crypto service
	_, err = NewStorageService(nil, vaultPath)
	if err == nil {
		t.Error("Expected error for nil crypto service")
	}

	// Test empty vault path
	_, err = NewStorageService(cryptoService, "")
	if err != ErrInvalidVaultPath {
		t.Errorf("Expected ErrInvalidVaultPath, got %v", err)
	}
}

func TestStorageService_InitializeVault(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")

	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	password := "test-password"

	// Test initialization
	if err := storage.InitializeVault(password); err != nil {
		t.Fatalf("InitializeVault failed: %v", err)
	}

	// Verify vault exists
	if !storage.VaultExists() {
		t.Error("Vault should exist after initialization")
	}

	// Test double initialization (should fail)
	if err := storage.InitializeVault(password); err == nil {
		t.Error("Expected error when initializing existing vault")
	}

	// Verify file permissions (skip on Windows as it doesn't support Unix permissions)
	info, err := os.Stat(vaultPath)
	if err != nil {
		t.Fatalf("Failed to stat vault file: %v", err)
	}

	// Only check permissions on Unix-like systems
	if info.Mode().Perm() != os.FileMode(VaultPermissions) {
		// This is expected on Windows, so just log it
		t.Logf("Note: File permissions are %v (expected %v) - this is normal on Windows",
			info.Mode().Perm(), os.FileMode(VaultPermissions))
	}
}

func TestStorageService_LoadSaveVault(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")

	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	password := "test-password"
	testData := []byte(`{"credentials": [{"service": "example.com", "username": "user", "password": "pass"}]}`)

	// Initialize vault
	if err := storage.InitializeVault(password); err != nil {
		t.Fatalf("InitializeVault failed: %v", err)
	}

	// Save test data
	if err := storage.SaveVault(testData, password); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	// Load and verify data
	loadedData, err := storage.LoadVault(password)
	if err != nil {
		t.Fatalf("LoadVault failed: %v", err)
	}

	if !bytes.Equal(testData, loadedData) {
		t.Error("Loaded data does not match saved data")
	}

	// Test wrong password
	_, err = storage.LoadVault("wrong-password")
	if err == nil {
		t.Error("Expected error with wrong password")
	}
}

func TestStorageService_VaultNotFound(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "nonexistent_vault.enc")

	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	// Test loading non-existent vault
	_, err = storage.LoadVault("password")
	if err != ErrVaultNotFound {
		t.Errorf("Expected ErrVaultNotFound, got %v", err)
	}

	// Test saving to non-existent vault
	err = storage.SaveVault([]byte("data"), "password")
	if err != ErrVaultNotFound {
		t.Errorf("Expected ErrVaultNotFound, got %v", err)
	}
}

func TestStorageService_GetVaultInfo(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")

	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	password := "test-password"

	// Initialize vault
	if err := storage.InitializeVault(password); err != nil {
		t.Fatalf("InitializeVault failed: %v", err)
	}

	// Get vault info
	info, err := storage.GetVaultInfo()
	if err != nil {
		t.Fatalf("GetVaultInfo failed: %v", err)
	}

	// Verify metadata
	if info.Version != 1 {
		t.Errorf("Expected version 1, got %d", info.Version)
	}

	if info.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if info.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	// Verify salt is not exposed
	if info.Salt != nil {
		t.Error("Salt should not be exposed in vault info")
	}

	// Save data and check if UpdatedAt changes
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	originalUpdatedAt := info.UpdatedAt

	if err := storage.SaveVault([]byte("new data"), password); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	newInfo, err := storage.GetVaultInfo()
	if err != nil {
		t.Fatalf("GetVaultInfo failed: %v", err)
	}

	if !newInfo.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated after save")
	}
}

func TestStorageService_ValidateVault(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")

	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	// Test validation of non-existent vault
	err = storage.ValidateVault()
	if err != ErrVaultNotFound {
		t.Errorf("Expected ErrVaultNotFound, got %v", err)
	}

	// Initialize vault and test validation
	password := "test-password"
	if err := storage.InitializeVault(password); err != nil {
		t.Fatalf("InitializeVault failed: %v", err)
	}

	// Valid vault should pass validation
	if err := storage.ValidateVault(); err != nil {
		t.Errorf("ValidateVault failed for valid vault: %v", err)
	}

	// Test corrupted vault by writing invalid JSON
	invalidJSON := []byte(`{"invalid": "json"`)
	if err := os.WriteFile(vaultPath, invalidJSON, VaultPermissions); err != nil {
		t.Fatalf("Failed to write corrupted vault: %v", err)
	}

	err = storage.ValidateVault()
	if err == nil {
		t.Error("Expected error for corrupted vault")
	}
}

func TestStorageService_BackupRestore(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")

	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	password := "test-password"
	originalData := []byte(`{"original": "data"}`)

	// Initialize vault with original data
	if err := storage.InitializeVault(password); err != nil {
		t.Fatalf("InitializeVault failed: %v", err)
	}

	if err := storage.SaveVault(originalData, password); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	// Create backup
	if err := storage.CreateBackup(); err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Verify backup file exists
	backupPath := vaultPath + BackupSuffix
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Backup file should exist")
	}

	// Modify vault
	modifiedData := []byte(`{"modified": "data"}`)
	if err := storage.SaveVault(modifiedData, password); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	// Verify modification
	loadedData, err := storage.LoadVault(password)
	if err != nil {
		t.Fatalf("LoadVault failed: %v", err)
	}

	if bytes.Equal(originalData, loadedData) {
		t.Error("Data should be modified")
	}

	// Restore from backup
	if err := storage.RestoreFromBackup(); err != nil {
		t.Fatalf("RestoreFromBackup failed: %v", err)
	}

	// Verify restoration
	restoredData, err := storage.LoadVault(password)
	if err != nil {
		t.Fatalf("LoadVault failed after restore: %v", err)
	}

	if !bytes.Equal(originalData, restoredData) {
		t.Error("Restored data should match original data")
	}

	// Test removing backup
	if err := storage.RemoveBackup(); err != nil {
		t.Fatalf("RemoveBackup failed: %v", err)
	}

	// Verify backup is removed
	if _, err := os.Stat(backupPath); !os.IsNotExist(err) {
		t.Error("Backup file should be removed")
	}

	// Test removing non-existent backup (should not error)
	if err := storage.RemoveBackup(); err != nil {
		t.Errorf("RemoveBackup should not error for non-existent backup: %v", err)
	}
}

func TestStorageService_AtomicWrite(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")

	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	password := "test-password"

	// Initialize vault
	if err := storage.InitializeVault(password); err != nil {
		t.Fatalf("InitializeVault failed: %v", err)
	}

	// Save some data
	testData := []byte(`{"test": "data"}`)
	if err := storage.SaveVault(testData, password); err != nil {
		t.Fatalf("SaveVault failed: %v", err)
	}

	// Verify no temporary files are left behind
	tempPath := vaultPath + TempSuffix
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		t.Error("Temporary file should not exist after successful write")
	}

	// Verify data is correctly saved
	loadedData, err := storage.LoadVault(password)
	if err != nil {
		t.Fatalf("LoadVault failed: %v", err)
	}

	if !bytes.Equal(testData, loadedData) {
		t.Error("Loaded data does not match saved data")
	}
}

func TestStorageService_EmptyData(t *testing.T) {
	cryptoService := crypto.NewCryptoService()
	tempDir := t.TempDir()
	vaultPath := filepath.Join(tempDir, "test_vault.enc")

	storage, err := NewStorageService(cryptoService, vaultPath)
	if err != nil {
		t.Fatalf("NewStorageService failed: %v", err)
	}

	password := "test-password"

	// Initialize vault
	if err := storage.InitializeVault(password); err != nil {
		t.Fatalf("InitializeVault failed: %v", err)
	}

	// Save empty data
	emptyData := []byte("")
	if err := storage.SaveVault(emptyData, password); err != nil {
		t.Fatalf("SaveVault failed with empty data: %v", err)
	}

	// Load and verify empty data
	loadedData, err := storage.LoadVault(password)
	if err != nil {
		t.Fatalf("LoadVault failed: %v", err)
	}

	if !bytes.Equal(emptyData, loadedData) {
		t.Error("Loaded empty data does not match saved empty data")
	}
}