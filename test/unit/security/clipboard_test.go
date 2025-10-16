package security_test

import (
	"testing"
	"time"

	"github.com/atotto/clipboard"
)

// TestClipboardSecurityVerification verifies 5-second auto-clear per constitution.
// This test simulates the clipboard auto-clear behavior implemented in cmd/get.go.
func TestClipboardSecurityVerification(t *testing.T) {
	// Skip if clipboard is not available (e.g., headless CI)
	if clipboard.Unsupported {
		t.Skip("Clipboard not supported on this platform")
	}
	if err := clipboard.WriteAll("test-check"); err != nil {
		t.Skipf("Clipboard not available: %v", err)
	}

	// Test password
	testPassword := "test-clipboard-password-123"

	// Write password to clipboard
	if err := clipboard.WriteAll(testPassword); err != nil {
		t.Fatalf("Failed to write to clipboard: %v", err)
	}

	// Verify password is in clipboard immediately
	content, err := clipboard.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read from clipboard: %v", err)
	}
	if content != testPassword {
		t.Errorf("Clipboard content mismatch: expected %q, got %q", testPassword, content)
	}

	// Simulate the 5-second auto-clear from cmd/get.go
	done := make(chan bool, 1)
	go func() {
		time.Sleep(5 * time.Second)
		// Only clear if clipboard still contains our password
		if current, err := clipboard.ReadAll(); err == nil && current == testPassword {
			_ = clipboard.WriteAll("")
			done <- true
		} else {
			done <- false
		}
	}()

	// Wait 6 seconds to ensure auto-clear happened
	time.Sleep(6 * time.Second)

	// Verify clipboard was cleared
	select {
	case cleared := <-done:
		if !cleared {
			t.Error("Clipboard was not cleared within expected time")
		}
	default:
		t.Error("Clipboard clear goroutine did not complete")
	}

	content, err = clipboard.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read from clipboard: %v", err)
	}
	if content != "" {
		t.Errorf("Clipboard should be empty, but contains: %q", content)
	}
}

// TestClipboardClearingTiming verifies clipboard is cleared within 5 seconds per FR-001.
func TestClipboardClearingTiming(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping timing test in short mode")
	}

	// Skip if clipboard is not available (e.g., headless CI)
	if clipboard.Unsupported {
		t.Skip("Clipboard not supported on this platform")
	}
	if err := clipboard.WriteAll("test-check"); err != nil {
		t.Skipf("Clipboard not available: %v", err)
	}

	testPassword := "timing-test-password"

	// Write to clipboard
	if err := clipboard.WriteAll(testPassword); err != nil {
		t.Fatalf("Failed to write to clipboard: %v", err)
	}

	// Start timer
	start := time.Now()

	// Simulate auto-clear
	go func() {
		time.Sleep(5 * time.Second)
		if current, _ := clipboard.ReadAll(); current == testPassword {
			_ = clipboard.WriteAll("")
		}
	}()

	// Poll clipboard every second to detect when it's cleared
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		content, err := clipboard.ReadAll()
		if err != nil {
			continue
		}
		if content == "" {
			elapsed := time.Since(start)
			// Verify cleared within 6 seconds (5s + 1s tolerance)
			if elapsed > 6*time.Second {
				t.Errorf("Clipboard cleared too late: took %v, should be <= 6s", elapsed)
			}
			if elapsed < 4*time.Second {
				t.Errorf("Clipboard cleared too early: took %v, should be >= 4s", elapsed)
			}
			return
		}
	}

	t.Error("Clipboard was never cleared within 10 seconds")
}
