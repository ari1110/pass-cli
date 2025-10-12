package cmd

import (
	"fmt"
	"os"

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
