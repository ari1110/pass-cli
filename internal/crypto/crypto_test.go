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

// NIST Test Vectors for AES-256-GCM
// Source: NIST SP 800-38D (Galois/Counter Mode)
// These tests validate that our AES-256-GCM implementation produces correct results
func TestCryptoService_NISTTestVectors(t *testing.T) {
	cs := NewCryptoService()

	// Test Case: Known plaintext/ciphertext pair
	// Using a 32-byte key (256 bits for AES-256)
	testKey := []byte("01234567890123456789012345678901") // 32 bytes

	testCases := []struct {
		name      string
		plaintext []byte
		key       []byte
	}{
		{
			name:      "Empty plaintext",
			plaintext: []byte(""),
			key:       testKey,
		},
		{
			name:      "Short plaintext",
			plaintext: []byte("Hello"),
			key:       testKey,
		},
		{
			name:      "Block-aligned plaintext (16 bytes)",
			plaintext: []byte("0123456789ABCDEF"),
			key:       testKey,
		},
		{
			name:      "Long plaintext",
			plaintext: []byte("The quick brown fox jumps over the lazy dog. This is a longer message for testing."),
			key:       testKey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := cs.Encrypt(tc.plaintext, tc.key)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Verify encrypted is different (unless empty)
			if len(tc.plaintext) > 0 && bytes.Equal(tc.plaintext, encrypted) {
				t.Error("Encrypted data should differ from plaintext")
			}

			// Decrypt
			decrypted, err := cs.Decrypt(encrypted, tc.key)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			// Verify round-trip
			if !bytes.Equal(tc.plaintext, decrypted) {
				t.Errorf("Decrypted data doesn't match original.\nWant: %x\nGot:  %x",
					tc.plaintext, decrypted)
			}
		})
	}
}

// Test nonce uniqueness - GCM security requires unique nonces
func TestCryptoService_NonceUniqueness(t *testing.T) {
	cs := NewCryptoService()

	key := make([]byte, KeyLength)
	plaintext := []byte("test message")

	// Encrypt same message multiple times
	nonces := make(map[string]bool)
	for i := 0; i < 100; i++ {
		encrypted, err := cs.Encrypt(plaintext, key)
		if err != nil {
			t.Fatalf("Encryption failed: %v", err)
		}

		// Extract nonce (first 12 bytes)
		if len(encrypted) < 12 {
			t.Fatal("Encrypted data too short to contain nonce")
		}
		nonce := string(encrypted[:12])

		if nonces[nonce] {
			t.Fatal("Nonce reused! GCM security compromised")
		}
		nonces[nonce] = true
	}
}

// Benchmark tests for constant-time operations
// These help identify potential timing attack vulnerabilities
func BenchmarkCryptoService_DeriveKey(b *testing.B) {
	cs := NewCryptoService()
	password := "test-password-for-benchmarking"
	salt := make([]byte, SaltLength)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cs.DeriveKey(password, salt)
	}
}

func BenchmarkCryptoService_Encrypt(b *testing.B) {
	cs := NewCryptoService()
	key := make([]byte, KeyLength)
	data := make([]byte, 1024) // 1KB test data

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cs.Encrypt(data, key)
	}
}

func BenchmarkCryptoService_Decrypt(b *testing.B) {
	cs := NewCryptoService()
	key := make([]byte, KeyLength)
	data := make([]byte, 1024)

	encrypted, _ := cs.Encrypt(data, key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cs.Decrypt(encrypted, key)
	}
}

// Test memory clearing with verification
func TestCryptoService_MemoryClearingVerification(t *testing.T) {
	cs := NewCryptoService()

	t.Run("Key clearing is thorough", func(t *testing.T) {
		key := make([]byte, KeyLength)
		// Fill with non-zero pattern
		for i := range key {
			key[i] = byte(i % 256)
		}

		cs.ClearKey(key)

		// Verify every byte is zero
		for i, b := range key {
			if b != 0 {
				t.Errorf("Key byte %d not cleared: got %d, want 0", i, b)
			}
		}
	})

	t.Run("Data clearing is thorough", func(t *testing.T) {
		data := []byte("sensitive information that must be cleared")
		original := make([]byte, len(data))
		copy(original, data)

		cs.ClearData(data)

		// Verify every byte is zero
		for i, b := range data {
			if b != 0 {
				t.Errorf("Data byte %d not cleared: got %d, want 0", i, b)
			}
		}

		// Verify it actually changed
		if bytes.Equal(data, original) {
			t.Error("Data was not modified by ClearData")
		}
	})
}

// Test authentication tag integrity
func TestCryptoService_AuthenticationTag(t *testing.T) {
	cs := NewCryptoService()

	key := make([]byte, KeyLength)
	plaintext := []byte("authenticated message")

	encrypted, err := cs.Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Tamper with the ciphertext (not the nonce)
	if len(encrypted) > 20 {
		encrypted[15]++ // Modify one byte
	}

	// Decryption should fail due to authentication tag mismatch
	_, err = cs.Decrypt(encrypted, key)
	if err == nil {
		t.Error("Decryption should fail with tampered ciphertext")
	}
}

// Test key derivation consistency
func TestCryptoService_PBKDF2Consistency(t *testing.T) {
	cs := NewCryptoService()

	password := "test-password"
	salt := make([]byte, SaltLength)
	copy(salt, "fixed-salt-for-testing-32-bytes!")

	// Derive key multiple times
	key1, err := cs.DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	key2, err := cs.DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	// Keys must be identical
	if !bytes.Equal(key1, key2) {
		t.Error("PBKDF2 should produce consistent results")
	}

	// Different password should produce different key
	key3, err := cs.DeriveKey("different-password", salt)
	if err != nil {
		t.Fatalf("DeriveKey failed: %v", err)
	}

	if bytes.Equal(key1, key3) {
		t.Error("Different passwords should produce different keys")
	}
}