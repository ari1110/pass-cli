package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/howeyc/gopass"
	"golang.org/x/term"
)

// readPassword reads a password from stdin with asterisk masking.
// Returns []byte for secure memory handling (no string conversion).
func readPassword() ([]byte, error) {
	// Get file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Check if stdin is a terminal
	if !term.IsTerminal(fd) {
		// Not a terminal, read normally (for testing/scripts)
		var password string
		_, err := fmt.Scanln(&password)
		return []byte(password), err
	}

	// Read password with asterisk masking using gopass
	passwordBytes, err := gopass.GetPasswdMasked()
	if err != nil {
		return nil, err
	}

	return passwordBytes, nil
}

// T072: getAuditLogPath returns the audit log path from environment variable or default
// Per FR-023: PASS_AUDIT_LOG environment variable for custom log location
func getAuditLogPath(vaultPath string) string {
	// Check environment variable first
	if auditPath := os.Getenv("PASS_AUDIT_LOG"); auditPath != "" {
		return auditPath
	}

	// Default: <vault-dir>/audit.log
	vaultDir := filepath.Dir(vaultPath)
	return filepath.Join(vaultDir, "audit.log")
}

// T072: getVaultID returns a unique identifier for the vault (used for keychain)
// Uses vault file path as unique identifier
func getVaultID(vaultPath string) string {
	// Use absolute path as vault ID for keychain
	absPath, err := filepath.Abs(vaultPath)
	if err != nil {
		return vaultPath // Fallback to relative path
	}
	return absPath
}
