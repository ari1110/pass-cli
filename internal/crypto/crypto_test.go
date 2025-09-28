package crypto

import (
	"bytes"
	"testing"
)

func TestCryptoService_GenerateSalt(t *testing.T) {
	cs := NewCryptoService()

	salt, err := cs.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt failed: %v", err)
	}

	if len(salt) != SaltLength {
		t.Errorf("Expected salt length %d, got %d", SaltLength, len(salt))
	}

	// Generate another salt to ensure they're different
	salt2, err := cs.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt failed: %v", err)
	}

	if bytes.Equal(salt, salt2) {
		t.Error("Two generated salts should not be equal")
	}
}

func TestCryptoService_DeriveKey(t *testing.T) {
	cs := NewCryptoService()
	password := "test-password"
	salt := make([]byte, SaltLength)

	key, err := cs.DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	if len(key) != KeyLength {
		t.Errorf("Expected key length %d, got %d", KeyLength, len(key))
	}

	// Same password and salt should produce same key
	key2, err := cs.DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	if !bytes.Equal(key, key2) {
		t.Error("Same password and salt should produce same key")
	}

	// Different salt should produce different key
	salt2 := make([]byte, SaltLength)
	salt2[0] = 1 // Make it different
	key3, err := cs.DeriveKey(password, salt2)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	if bytes.Equal(key, key3) {
		t.Error("Different salts should produce different keys")
	}
}

func TestCryptoService_EncryptDecrypt(t *testing.T) {
	cs := NewCryptoService()

	// Generate key
	salt, err := cs.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt failed: %v", err)
	}

	key, err := cs.DeriveKey("test-password", salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	// Test data
	testData := []byte("Hello, World! This is a test message.")

	// Encrypt
	encrypted, err := cs.Encrypt(testData, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Verify encrypted data is different
	if bytes.Equal(testData, encrypted) {
		t.Error("Encrypted data should be different from original")
	}

	// Decrypt
	decrypted, err := cs.Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// Verify decrypted data matches original
	if !bytes.Equal(testData, decrypted) {
		t.Error("Decrypted data should match original")
	}
}

func TestCryptoService_EncryptDecryptEmpty(t *testing.T) {
	cs := NewCryptoService()

	salt, err := cs.GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt failed: %v", err)
	}

	key, err := cs.DeriveKey("test-password", salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	// Test empty data
	testData := []byte("")

	encrypted, err := cs.Encrypt(testData, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := cs.Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(testData, decrypted) {
		t.Error("Decrypted empty data should match original")
	}
}

func TestCryptoService_SecureRandom(t *testing.T) {
	cs := NewCryptoService()

	// Test different lengths
	lengths := []int{1, 16, 32, 64, 128}
	for _, length := range lengths {
		randomBytes, err := cs.SecureRandom(length)
		if err != nil {
			t.Fatalf("SecureRandom failed for length %d: %v", length, err)
		}

		if len(randomBytes) != length {
			t.Errorf("Expected length %d, got %d", length, len(randomBytes))
		}
	}

	// Test that two calls produce different results
	random1, err := cs.SecureRandom(32)
	if err != nil {
		t.Fatalf("SecureRandom failed: %v", err)
	}

	random2, err := cs.SecureRandom(32)
	if err != nil {
		t.Fatalf("SecureRandom failed: %v", err)
	}

	if bytes.Equal(random1, random2) {
		t.Error("Two random byte arrays should not be equal")
	}
}

func TestCryptoService_InvalidInputs(t *testing.T) {
	cs := NewCryptoService()

	// Test invalid key length for encryption
	shortKey := make([]byte, 16) // Too short
	data := []byte("test")
	_, err := cs.Encrypt(data, shortKey)
	if err != ErrInvalidKeyLength {
		t.Errorf("Expected ErrInvalidKeyLength, got %v", err)
	}

	// Test invalid key length for decryption
	_, err = cs.Decrypt(data, shortKey)
	if err != ErrInvalidKeyLength {
		t.Errorf("Expected ErrInvalidKeyLength, got %v", err)
	}

	// Test invalid salt length for key derivation
	shortSalt := make([]byte, 16) // Too short
	_, err = cs.DeriveKey("password", shortSalt)
	if err != ErrInvalidSaltLength {
		t.Errorf("Expected ErrInvalidSaltLength, got %v", err)
	}

	// Test invalid ciphertext length for decryption
	validKey := make([]byte, KeyLength)
	shortCiphertext := make([]byte, 5) // Too short
	_, err = cs.Decrypt(shortCiphertext, validKey)
	if err != ErrInvalidCiphertext {
		t.Errorf("Expected ErrInvalidCiphertext, got %v", err)
	}

	// Test invalid length for SecureRandom
	_, err = cs.SecureRandom(0)
	if err == nil {
		t.Error("Expected error for invalid length 0")
	}

	_, err = cs.SecureRandom(-1)
	if err == nil {
		t.Error("Expected error for negative length")
	}
}

func TestCryptoService_ClearMethods(t *testing.T) {
	cs := NewCryptoService()

	// Test clearing key
	key := make([]byte, KeyLength)
	copy(key, "test-key-data-here-32-bytes-long")
	cs.ClearKey(key)

	// Verify key is cleared
	emptyKey := make([]byte, KeyLength)
	if !bytes.Equal(key, emptyKey) {
		t.Error("Key should be cleared to zeros")
	}

	// Test clearing data
	data := []byte("sensitive data")
	cs.ClearData(data)

	// Verify data is cleared
	emptyData := make([]byte, len(data))
	if !bytes.Equal(data, emptyData) {
		t.Error("Data should be cleared to zeros")
	}

	// Test clearing nil values (should not panic)
	cs.ClearKey(nil)
	cs.ClearData(nil)
}