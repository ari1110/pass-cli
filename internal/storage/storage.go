package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"pass-cli/internal/crypto"
)

const (
	VaultPermissions = 0600 // Read/write for owner only
	DefaultVaultName = "vault.enc"
	BackupSuffix     = ".backup"
	TempSuffix       = ".tmp"
)

var (
	ErrVaultNotFound     = errors.New("vault file not found")
	ErrVaultCorrupted    = errors.New("vault file corrupted")
	ErrInvalidVaultPath  = errors.New("invalid vault path")
	ErrBackupFailed      = errors.New("backup operation failed")
	ErrAtomicWriteFailed = errors.New("atomic write operation failed")
)

type VaultMetadata struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Salt      []byte    `json:"salt"`
}

type EncryptedVault struct {
	Metadata VaultMetadata `json:"metadata"`
	Data     []byte        `json:"data"`
}

type StorageService struct {
	cryptoService *crypto.CryptoService
	vaultPath     string
}

func NewStorageService(cryptoService *crypto.CryptoService, vaultPath string) (*StorageService, error) {
	if cryptoService == nil {
		return nil, errors.New("crypto service cannot be nil")
	}

	if vaultPath == "" {
		return nil, ErrInvalidVaultPath
	}

	// Ensure the directory exists
	dir := filepath.Dir(vaultPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create vault directory: %w", err)
	}

	return &StorageService{
		cryptoService: cryptoService,
		vaultPath:     vaultPath,
	}, nil
}

func (s *StorageService) InitializeVault(password string) error {
	// Check if vault already exists
	if s.VaultExists() {
		return errors.New("vault already exists")
	}

	// Generate salt for key derivation
	salt, err := s.cryptoService.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Create initial empty vault data
	emptyVault := []byte("{}")

	// Create vault metadata
	metadata := VaultMetadata{
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Salt:      salt,
	}

	// Encrypt and save vault
	if err := s.saveEncryptedVault(emptyVault, metadata, password); err != nil {
		return fmt.Errorf("failed to initialize vault: %w", err)
	}

	return nil
}

func (s *StorageService) LoadVault(password string) ([]byte, error) {
	encryptedVault, err := s.loadEncryptedVault()
	if err != nil {
		return nil, err
	}

	// Derive key from password and salt
	key, err := s.cryptoService.DeriveKey([]byte(password), encryptedVault.Metadata.Salt)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}
	defer s.cryptoService.ClearKey(key)

	// Decrypt vault data
	plaintext, err := s.cryptoService.Decrypt(encryptedVault.Data, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault (invalid password?): %w", err)
	}

	return plaintext, nil
}

func (s *StorageService) SaveVault(data []byte, password string) error {
	// Load existing vault to get metadata
	encryptedVault, err := s.loadEncryptedVault()
	if err != nil {
		return err
	}

	// Update metadata
	encryptedVault.Metadata.UpdatedAt = time.Now()

	// Create backup before saving
	if err := s.createBackup(); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Save encrypted vault
	if err := s.saveEncryptedVault(data, encryptedVault.Metadata, password); err != nil {
		// Restore from backup on failure
		if restoreErr := s.restoreFromBackup(); restoreErr != nil {
			return fmt.Errorf("save failed and backup restore failed: %v (original error: %w)", restoreErr, err)
		}
		return fmt.Errorf("failed to save vault: %w", err)
	}

	return nil
}

func (s *StorageService) VaultExists() bool {
	_, err := os.Stat(s.vaultPath)
	return err == nil
}

func (s *StorageService) GetVaultInfo() (*VaultMetadata, error) {
	encryptedVault, err := s.loadEncryptedVault()
	if err != nil {
		return nil, err
	}

	// Return a copy of metadata (without the salt for security)
	info := VaultMetadata{
		Version:   encryptedVault.Metadata.Version,
		CreatedAt: encryptedVault.Metadata.CreatedAt,
		UpdatedAt: encryptedVault.Metadata.UpdatedAt,
		Salt:      nil, // Don't expose salt
	}

	return &info, nil
}

func (s *StorageService) ValidateVault() error {
	encryptedVault, err := s.loadEncryptedVault()
	if err != nil {
		return err
	}

	// Basic validation checks
	if encryptedVault.Metadata.Version <= 0 {
		return ErrVaultCorrupted
	}

	if len(encryptedVault.Metadata.Salt) != 32 {
		return ErrVaultCorrupted
	}

	if len(encryptedVault.Data) == 0 {
		return ErrVaultCorrupted
	}

	if encryptedVault.Metadata.CreatedAt.IsZero() {
		return ErrVaultCorrupted
	}

	if encryptedVault.Metadata.UpdatedAt.Before(encryptedVault.Metadata.CreatedAt) {
		return ErrVaultCorrupted
	}

	return nil
}

func (s *StorageService) CreateBackup() error {
	return s.createBackup()
}

func (s *StorageService) RestoreFromBackup() error {
	return s.restoreFromBackup()
}

func (s *StorageService) RemoveBackup() error {
	backupPath := s.vaultPath + BackupSuffix
	err := os.Remove(backupPath)
	if os.IsNotExist(err) {
		return nil // Backup doesn't exist, which is fine
	}
	return err
}

// Private helper methods

func (s *StorageService) loadEncryptedVault() (*EncryptedVault, error) {
	if !s.VaultExists() {
		return nil, ErrVaultNotFound
	}

	data, err := os.ReadFile(s.vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	var encryptedVault EncryptedVault
	if err := json.Unmarshal(data, &encryptedVault); err != nil {
		return nil, fmt.Errorf("failed to parse vault file: %w", err)
	}

	return &encryptedVault, nil
}

func (s *StorageService) saveEncryptedVault(data []byte, metadata VaultMetadata, password string) error {
	// Derive key from password and salt
	key, err := s.cryptoService.DeriveKey([]byte(password), metadata.Salt)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}
	defer s.cryptoService.ClearKey(key)

	// Encrypt vault data
	encryptedData, err := s.cryptoService.Encrypt(data, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault data: %w", err)
	}

	// Create encrypted vault structure
	encryptedVault := EncryptedVault{
		Metadata: metadata,
		Data:     encryptedData,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(encryptedVault)
	if err != nil {
		return fmt.Errorf("failed to marshal vault data: %w", err)
	}

	// Atomic write using temporary file
	return s.atomicWrite(s.vaultPath, jsonData)
}

func (s *StorageService) atomicWrite(path string, data []byte) error {
	tempPath := path + TempSuffix

	// Write to temporary file
	// #nosec G304 -- Vault path is user-controlled by design for CLI tool
	tempFile, err := os.OpenFile(tempPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, VaultPermissions)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	// Ensure temp file is cleaned up on error
	defer func() {
		if tempFile != nil {
			_ = tempFile.Close()
			_ = os.Remove(tempPath)
		}
	}()

	// Write data
	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	// Sync to ensure data is written to disk
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync data: %w", err)
	}

	// Close file
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	tempFile = nil // Prevent cleanup in defer

	// Atomic move (rename) to final location
	if err := os.Rename(tempPath, path); err != nil {
		_ = os.Remove(tempPath) // Clean up on failure
		return fmt.Errorf("failed to move temp file to final location: %w", err)
	}

	return nil
}

func (s *StorageService) createBackup() error {
	if !s.VaultExists() {
		return nil // No vault to backup
	}

	backupPath := s.vaultPath + BackupSuffix

	// Copy vault file to backup
	src, err := os.Open(s.vaultPath)
	if err != nil {
		return fmt.Errorf("failed to open vault for backup: %w", err)
	}
	defer func() { _ = src.Close() }()

	// #nosec G304 -- Backup path is user-controlled by design for CLI tool
	dst, err := os.OpenFile(backupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, VaultPermissions)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer func() { _ = dst.Close() }()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy vault to backup: %w", err)
	}

	if err := dst.Sync(); err != nil {
		return fmt.Errorf("failed to sync backup file: %w", err)
	}

	return nil
}

func (s *StorageService) restoreFromBackup() error {
	backupPath := s.vaultPath + BackupSuffix

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return ErrBackupFailed
	}

	// Copy backup to vault location
	// #nosec G304 -- Backup path is user-controlled by design for CLI tool
	src, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() { _ = src.Close() }()

	dst, err := os.OpenFile(s.vaultPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, VaultPermissions)
	if err != nil {
		return fmt.Errorf("failed to create vault file: %w", err)
	}
	defer func() { _ = dst.Close() }()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}

	if err := dst.Sync(); err != nil {
		return fmt.Errorf("failed to sync restored vault: %w", err)
	}

	return nil
}
