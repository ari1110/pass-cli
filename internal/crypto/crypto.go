package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"

	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)

const (
	KeyLength   = 32     // AES-256 key length
	NonceLength = 12     // GCM nonce length
	SaltLength  = 32     // PBKDF2 salt length
	Iterations  = 100000 // PBKDF2 iterations for key derivation
)

var (
	ErrInvalidKeyLength   = errors.New("invalid key length")
	ErrInvalidNonceLength = errors.New("invalid nonce length")
	ErrInvalidSaltLength  = errors.New("invalid salt length")
	ErrDecryptionFailed   = errors.New("decryption failed")
	ErrInvalidCiphertext  = errors.New("invalid ciphertext length")
)

type CryptoService struct{}

func NewCryptoService() *CryptoService {
	return &CryptoService{}
}

func (c *CryptoService) GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func (c *CryptoService) DeriveKey(password []byte, salt []byte) ([]byte, error) {
	if len(salt) != SaltLength {
		return nil, ErrInvalidSaltLength
	}

	key := pbkdf2.Key(password, salt, Iterations, KeyLength, sha256.New)
	return key, nil
}

func (c *CryptoService) Encrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != KeyLength {
		return nil, ErrInvalidKeyLength
	}

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate nonce
	nonce := make([]byte, NonceLength)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	// #nosec G407 -- Nonce is randomly generated via crypto/rand (line 75), not hardcoded
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Prepend nonce to ciphertext
	result := make([]byte, NonceLength+len(ciphertext))
	copy(result[:NonceLength], nonce)
	copy(result[NonceLength:], ciphertext)

	return result, nil
}

func (c *CryptoService) Decrypt(encryptedData []byte, key []byte) ([]byte, error) {
	if len(key) != KeyLength {
		return nil, ErrInvalidKeyLength
	}

	if len(encryptedData) < NonceLength {
		return nil, ErrInvalidCiphertext
	}

	// Extract nonce and ciphertext
	nonce := encryptedData[:NonceLength]
	ciphertext := encryptedData[NonceLength:]

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	return plaintext, nil
}

func (c *CryptoService) SecureRandom(length int) ([]byte, error) {
	if length <= 0 {
		return nil, errors.New("invalid length")
	}

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return randomBytes, nil
}

func (c *CryptoService) ClearKey(key []byte) {
	if key != nil {
		ClearBytes(key)
	}
}

func (c *CryptoService) ClearData(data []byte) {
	if data != nil {
		ClearBytes(data)
	}
}

// ClearBytes securely zeros a byte array by overwriting with zeros.
// Uses crypto/subtle.ConstantTimeCompare as a compiler barrier to prevent
// the compiler from optimizing away the zeroing operation.
// Exposed publicly for use in vault and CLI packages.
func ClearBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}

	// Use subtle.ConstantTimeByteEq to prevent compiler optimizations
	// that might remove the clearing operation
	dummy := make([]byte, len(data))
	subtle.ConstantTimeCompare(data, dummy)
}
