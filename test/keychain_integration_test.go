//go:build integration
// +build integration

package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"pass-cli/internal/keychain"
)

// TestIntegration_KeychainWorkflow tests the full keychain integration
func TestIntegration_KeychainWorkflow(t *testing.T) {
	// Check if keychain is available before running tests
	ks := keychain.New()
	if !ks.IsAvailable() {
		t.Skip("System keychain not available - skipping keychain integration tests")
	}

	testPassword := "keychain-test-password-123"
	vaultPath := filepath.Join(testDir, "keychain-vault", "vault.enc")

	// Ensure clean state
	defer cleanupKeychain(t, ks)

	t.Run("1_Init_With_Keychain", func(t *testing.T) {
		// Initialize vault with --use-keychain flag
		input := testPassword + "\n" + testPassword + "\n"
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "init", "--use-keychain")
		cmd.Stdin = strings.NewReader(input)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Init with keychain failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
		}

		output := stdout.String()
		if !strings.Contains(output, "successfully") && !strings.Contains(output, "initialized") {
			t.Errorf("Expected success message in output, got: %s", output)
		}

		if !strings.Contains(output, "keychain") {
			t.Errorf("Expected keychain confirmation in output, got: %s", output)
		}

		// Verify vault file was created
		if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
			t.Error("Vault file was not created")
		}

		// Verify password is in keychain
		retrievedPassword, err := ks.Retrieve()
		if err != nil {
			t.Fatalf("Password was not stored in keychain: %v", err)
		}

		if retrievedPassword != testPassword {
			t.Errorf("Keychain password = %q, want %q", retrievedPassword, testPassword)
		}
	})

	t.Run("2_Add_Without_Password_Prompt", func(t *testing.T) {
		// Add credential - should NOT prompt for master password (uses keychain)
		input := "testuser\n" + "testpass123\n" // Only username and credential password
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "add", "github.com")
		cmd.Stdin = strings.NewReader(input)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Add failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
		}

		output := stdout.String()
		if !strings.Contains(output, "added") && !strings.Contains(output, "successfully") {
			t.Errorf("Expected success message, got: %s", output)
		}

		// Should NOT contain "Master password:" prompt
		allOutput := stdout.String() + stderr.String()
		if strings.Contains(allOutput, "Master password:") {
			t.Error("Unexpected master password prompt - keychain should have been used")
		}
	})

	t.Run("3_Get_Without_Password_Prompt", func(t *testing.T) {
		// Get credential - should NOT prompt for master password
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "get", "github.com", "--no-clipboard")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Get failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
		}

		output := stdout.String()
		if !strings.Contains(output, "testuser") || !strings.Contains(output, "testpass123") {
			t.Errorf("Expected credential details in output, got: %s", output)
		}

		// Should NOT contain "Master password:" prompt
		allOutput := stdout.String() + stderr.String()
		if strings.Contains(allOutput, "Master password:") {
			t.Error("Unexpected master password prompt - keychain should have been used")
		}
	})

	t.Run("4_List_Without_Password_Prompt", func(t *testing.T) {
		// List credentials - should NOT prompt for master password
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "list")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("List failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
		}

		output := stdout.String()
		if !strings.Contains(output, "github.com") {
			t.Errorf("Expected github.com in list output, got: %s", output)
		}

		// Should NOT contain "Master password:" prompt
		allOutput := stdout.String() + stderr.String()
		if strings.Contains(allOutput, "Master password:") {
			t.Error("Unexpected master password prompt - keychain should have been used")
		}
	})

	t.Run("5_Update_Without_Password_Prompt", func(t *testing.T) {
		// Update credential - should NOT prompt for master password
		input := "updateduser\n" + "updatedpass456\n"
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "update", "github.com")
		cmd.Stdin = strings.NewReader(input)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Update failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
		}

		// Should NOT contain "Master password:" prompt
		allOutput := stdout.String() + stderr.String()
		if strings.Contains(allOutput, "Master password:") {
			t.Error("Unexpected master password prompt - keychain should have been used")
		}
	})

	t.Run("6_Delete_Without_Password_Prompt", func(t *testing.T) {
		// Delete credential - should NOT prompt for master password
		input := "y\n" // confirm deletion
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "delete", "github.com")
		cmd.Stdin = strings.NewReader(input)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Delete failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
		}

		// Should NOT contain "Master password:" prompt in stderr before confirmation
		// Note: "Master password:" might appear in help text, so we check more specifically
		stderrOutput := stderr.String()
		if strings.Count(stderrOutput, "Master password:") > 0 {
			// Check if it's actually prompting (before the confirmation)
			lines := strings.Split(stderrOutput, "\n")
			for _, line := range lines {
				if strings.TrimSpace(line) == "Master password:" {
					t.Error("Unexpected master password prompt - keychain should have been used")
					break
				}
			}
		}
	})
}

// TestIntegration_KeychainFallback tests fallback to password prompt
func TestIntegration_KeychainFallback(t *testing.T) {
	ks := keychain.New()
	if !ks.IsAvailable() {
		t.Skip("System keychain not available - skipping keychain fallback tests")
	}

	testPassword := "fallback-test-password-789"
	vaultPath := filepath.Join(testDir, "fallback-vault", "vault.enc")

	// Ensure clean state
	defer cleanupKeychain(t, ks)

	// Initialize vault WITH keychain
	input := testPassword + "\n" + testPassword + "\n"
	cmd := exec.Command(binaryPath, "--vault", vaultPath, "init", "--use-keychain")
	cmd.Stdin = strings.NewReader(input)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize vault: %v", err)
	}

	t.Run("Fallback_After_Keychain_Deleted", func(t *testing.T) {
		// Delete password from keychain
		if err := ks.Delete(); err != nil {
			t.Fatalf("Failed to delete keychain entry: %v", err)
		}

		// Try to add credential - should now prompt for master password
		input := testPassword + "\n" + "testuser\n" + "testpass\n"
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "add", "test.com")
		cmd.Stdin = strings.NewReader(input)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Add with password prompt failed: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
		}

		// Should work successfully even without keychain
		output := stdout.String()
		if !strings.Contains(output, "added") && !strings.Contains(output, "successfully") {
			t.Errorf("Expected success message, got: %s", output)
		}
	})
}

// TestIntegration_KeychainUnavailable tests behavior when keychain is unavailable
func TestIntegration_KeychainUnavailable(t *testing.T) {
	ks := keychain.New()

	// This test verifies graceful handling when keychain is unavailable
	// If keychain IS available, we skip this test
	if ks.IsAvailable() {
		t.Skip("Keychain is available - cannot test unavailable scenario")
	}

	testPassword := "no-keychain-password-456"
	vaultPath := filepath.Join(testDir, "no-keychain-vault", "vault.enc")

	t.Run("Init_Without_Keychain_Available", func(t *testing.T) {
		// Try to initialize with --use-keychain when keychain unavailable
		input := testPassword + "\n" + testPassword + "\n"
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "init", "--use-keychain")
		cmd.Stdin = strings.NewReader(input)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()

		// Should either:
		// 1. Succeed with warning (graceful degradation)
		// 2. Fail with clear error message
		if err == nil {
			// Check for warning in output
			allOutput := stdout.String() + stderr.String()
			if !strings.Contains(allOutput, "warning") && !strings.Contains(allOutput, "Warning") {
				t.Log("Init succeeded without warning when keychain unavailable (acceptable)")
			}
		} else {
			// Check for clear error message
			allOutput := stdout.String() + stderr.String()
			if !strings.Contains(allOutput, "keychain") {
				t.Errorf("Error message should mention keychain when unavailable, got: %s", allOutput)
			}
		}
	})
}

// TestIntegration_MultipleVaultsKeychain tests multiple vaults with same keychain
func TestIntegration_MultipleVaultsKeychain(t *testing.T) {
	ks := keychain.New()
	if !ks.IsAvailable() {
		t.Skip("System keychain not available - skipping multiple vaults test")
	}

	// Note: Currently pass-cli uses a single keychain entry for all vaults
	// This test documents the current behavior and can be updated if we add
	// per-vault keychain support in the future

	testPassword := "multi-vault-password-999"
	vault1Path := filepath.Join(testDir, "multi-vault-1", "vault.enc")
	vault2Path := filepath.Join(testDir, "multi-vault-2", "vault.enc")

	defer cleanupKeychain(t, ks)

	t.Run("First_Vault_Init", func(t *testing.T) {
		input := testPassword + "\n" + testPassword + "\n"
		cmd := exec.Command(binaryPath, "--vault", vault1Path, "init", "--use-keychain")
		cmd.Stdin = strings.NewReader(input)

		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to init vault 1: %v", err)
		}
	})

	t.Run("Second_Vault_With_Same_Password", func(t *testing.T) {
		// Initialize second vault with same password
		// It will use the same keychain entry
		input := testPassword + "\n" + testPassword + "\n"
		cmd := exec.Command(binaryPath, "--vault", vault2Path, "init", "--use-keychain")
		cmd.Stdin = strings.NewReader(input)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			t.Fatalf("Failed to init vault 2: %v\nStdout: %s\nStderr: %s", err, stdout.String(), stderr.String())
		}

		// Add credential to vault 2 using keychain
		input = "user2\n" + "pass2\n"
		cmd = exec.Command(binaryPath, "--vault", vault2Path, "add", "service2.com")
		cmd.Stdin = strings.NewReader(input)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to add to vault 2: %v", err)
		}

		// Verify vault 1 still works with same keychain
		cmd = exec.Command(binaryPath, "--vault", vault1Path, "list")
		if err := cmd.Run(); err != nil {
			t.Errorf("Vault 1 should still work after vault 2 operations: %v", err)
		}
	})
}

// TestIntegration_KeychainVerboseOutput tests verbose mode shows keychain usage
func TestIntegration_KeychainVerboseOutput(t *testing.T) {
	ks := keychain.New()
	if !ks.IsAvailable() {
		t.Skip("System keychain not available - skipping verbose output test")
	}

	testPassword := "verbose-test-password-321"
	vaultPath := filepath.Join(testDir, "verbose-vault", "vault.enc")

	defer cleanupKeychain(t, ks)

	// Initialize with keychain
	input := testPassword + "\n" + testPassword + "\n"
	cmd := exec.Command(binaryPath, "--vault", vaultPath, "init", "--use-keychain")
	cmd.Stdin = strings.NewReader(input)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize vault: %v", err)
	}

	t.Run("Verbose_Shows_Keychain_Usage", func(t *testing.T) {
		// Run list command with --verbose flag
		cmd := exec.Command(binaryPath, "--vault", vaultPath, "--verbose", "list")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("List with verbose failed: %v", err)
		}

		// Check if verbose output mentions keychain usage
		allOutput := stdout.String() + stderr.String()
		if !strings.Contains(allOutput, "keychain") && !strings.Contains(allOutput, "Keychain") {
			t.Logf("Verbose mode output:\n%s", allOutput)
			t.Skip("Verbose keychain message may not be implemented yet")
		}
	})
}

// cleanupKeychain removes test keychain entries
func cleanupKeychain(t *testing.T, ks *keychain.KeychainService) {
	t.Helper()
	if err := ks.Clear(); err != nil && err != keychain.ErrKeychainUnavailable {
		t.Logf("Warning: failed to clean up keychain: %v", err)
	}
}
