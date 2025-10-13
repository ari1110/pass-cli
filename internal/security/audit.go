package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// T057: AuditLogEntry represents a single security event with tamper-evident HMAC signature
// Per data-model.md:256-262
type AuditLogEntry struct {
	Timestamp      time.Time `json:"timestamp"`       // Event time (FR-019, FR-020)
	EventType      string    `json:"event_type"`      // Type of operation (see constants below)
	Outcome        string    `json:"outcome"`         // "success" or "failure"
	CredentialName string    `json:"credential_name"` // Service name (NOT password, FR-021)
	HMACSignature  []byte    `json:"hmac_signature"`  // Tamper detection (FR-022)
}

// T058: Event type constants for audit logging
// Per data-model.md:268-277
const (
	EventVaultUnlock         = "vault_unlock"          // FR-019
	EventVaultLock           = "vault_lock"            // FR-019
	EventVaultPasswordChange = "vault_password_change" // FR-019
	EventCredentialAccess    = "credential_access"     // FR-020 (get)
	EventCredentialAdd       = "credential_add"        // FR-020
	EventCredentialUpdate    = "credential_update"     // FR-020
	EventCredentialDelete    = "credential_delete"     // FR-020
)

// Outcome constants
const (
	OutcomeSuccess = "success"
	OutcomeFailure = "failure"
)

// T059: AuditLogger manages tamper-evident audit logging
// Per data-model.md:332-337
type AuditLogger struct {
	filePath     string
	maxSizeBytes int64  // Default: 10MB (FR-024)
	currentSize  int64  // Current log file size
	auditKey     []byte // HMAC key for signing entries
}

// T060: Sign calculates HMAC signature for audit log entry
// Per data-model.md:291-305
func (e *AuditLogEntry) Sign(key []byte) error {
	// Canonical serialization (order matters!)
	data := fmt.Sprintf("%s|%s|%s|%s",
		e.Timestamp.Format(time.RFC3339Nano),
		e.EventType,
		e.Outcome,
		e.CredentialName,
	)

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	e.HMACSignature = mac.Sum(nil)

	return nil
}

// T061: Verify validates HMAC signature for audit log entry
// Per data-model.md:307-326
func (e *AuditLogEntry) Verify(key []byte) error {
	// Recalculate HMAC
	data := fmt.Sprintf("%s|%s|%s|%s",
		e.Timestamp.Format(time.RFC3339Nano),
		e.EventType,
		e.Outcome,
		e.CredentialName,
	)

	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	expected := mac.Sum(nil)

	// Constant-time comparison to prevent timing attacks
	if !hmac.Equal(e.HMACSignature, expected) {
		return fmt.Errorf("HMAC verification failed at %s", e.Timestamp)
	}

	return nil
}

// T062: ShouldRotate checks if log rotation is needed
// Per data-model.md:339-341
func (l *AuditLogger) ShouldRotate() bool {
	return l.currentSize >= l.maxSizeBytes
}

// T063: Rotate renames current log to .old and creates new empty log
// Per data-model.md:343-347
func (l *AuditLogger) Rotate() error {
	// Rename current log to .old
	oldPath := l.filePath + ".old"
	if err := os.Rename(l.filePath, oldPath); err != nil {
		// If file doesn't exist, that's OK
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to rotate log: %w", err)
		}
	}

	// Create new empty log
	f, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create new log: %w", err)
	}
	f.Close()

	// Reset size counter
	l.currentSize = 0

	return nil
}

// T064: Log writes an audit entry with HMAC signature and handles rotation
func (l *AuditLogger) Log(entry *AuditLogEntry) error {
	// Sign the entry
	if err := entry.Sign(l.auditKey); err != nil {
		return fmt.Errorf("failed to sign entry: %w", err)
	}

	// Check if rotation needed
	if l.ShouldRotate() {
		if err := l.Rotate(); err != nil {
			return fmt.Errorf("failed to rotate log: %w", err)
		}
	}

	// Serialize entry to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	// Append to log file
	f, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer f.Close()

	// Write JSON entry with newline
	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	// Update size counter
	l.currentSize += int64(len(data) + 1)

	return nil
}
