package cmd

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// readPassword reads a password from stdin without echoing
func readPassword() (string, error) {
	// Get file descriptor for stdin
	fd := int(os.Stdin.Fd())

	// Check if stdin is a terminal
	if !term.IsTerminal(fd) {
		// Not a terminal, read normally (for testing/scripts)
		var password string
		_, err := fmt.Scanln(&password)
		return password, err
	}

	// Read password without echoing
	passwordBytes, err := term.ReadPassword(fd)
	if err != nil {
		return "", err
	}

	return string(passwordBytes), nil
}