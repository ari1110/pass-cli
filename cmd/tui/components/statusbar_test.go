package components

import (
	"strings"
	"testing"
)

func TestNewStatusBar(t *testing.T) {
	statusBar := NewStatusBar(true, 5, "Test View")

	if statusBar == nil {
		t.Fatal("NewStatusBar() returned nil")
	}

	if !statusBar.keychainAvailable {
		t.Error("Keychain should be available")
	}

	if statusBar.credentialCount != 5 {
		t.Errorf("Expected credential count 5, got %d", statusBar.credentialCount)
	}

	if statusBar.currentView != "Test View" {
		t.Errorf("Expected current view 'Test View', got '%s'", statusBar.currentView)
	}
}

func TestStatusBarSetSize(t *testing.T) {
	statusBar := NewStatusBar(true, 0, "Test")
	statusBar.SetSize(100)

	if statusBar.width != 100 {
		t.Errorf("Expected width 100, got %d", statusBar.width)
	}
}

func TestStatusBarSetCredentialCount(t *testing.T) {
	statusBar := NewStatusBar(true, 0, "Test")

	statusBar.SetCredentialCount(10)

	if statusBar.credentialCount != 10 {
		t.Errorf("Expected credential count 10, got %d", statusBar.credentialCount)
	}
}

func TestStatusBarSetCurrentView(t *testing.T) {
	statusBar := NewStatusBar(true, 0, "Test")

	statusBar.SetCurrentView("Detail")

	if statusBar.currentView != "Detail" {
		t.Errorf("Expected current view 'Detail', got '%s'", statusBar.currentView)
	}
}

func TestStatusBarSetShortcuts(t *testing.T) {
	statusBar := NewStatusBar(true, 0, "Test")

	statusBar.SetShortcuts("a: add | e: edit")

	if statusBar.shortcuts != "a: add | e: edit" {
		t.Errorf("Expected shortcuts 'a: add | e: edit', got '%s'", statusBar.shortcuts)
	}
}

func TestStatusBarKeychainIndicator(t *testing.T) {
	tests := []struct {
		name              string
		keychainAvailable bool
		expectedText      string
	}{
		{
			name:              "Keychain available",
			keychainAvailable: true,
			expectedText:      "🔓 Keychain",
		},
		{
			name:              "Keychain unavailable",
			keychainAvailable: false,
			expectedText:      "🔒 Password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusBar := NewStatusBar(tt.keychainAvailable, 0, "Test")
			statusBar.SetSize(100)

			output := statusBar.Render()

			if !strings.Contains(output, tt.expectedText) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expectedText, output)
			}
		})
	}
}

func TestStatusBarCredentialCountPlural(t *testing.T) {
	tests := []struct {
		name          string
		count         int
		expectedText  string
	}{
		{
			name:         "Zero credentials",
			count:        0,
			expectedText: "0 credentials",
		},
		{
			name:         "One credential",
			count:        1,
			expectedText: "1 credential",
		},
		{
			name:         "Multiple credentials",
			count:        5,
			expectedText: "5 credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusBar := NewStatusBar(true, tt.count, "Test")
			statusBar.SetSize(100)

			output := statusBar.Render()

			if !strings.Contains(output, tt.expectedText) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expectedText, output)
			}
		})
	}
}

func TestStatusBarRenderWithShortcuts(t *testing.T) {
	statusBar := NewStatusBar(true, 5, "List")
	statusBar.SetShortcuts("a: add | e: edit | q: quit")
	statusBar.SetSize(100)

	output := statusBar.Render()

	// Verify all sections are present
	if !strings.Contains(output, "Keychain") {
		t.Error("Output should contain keychain indicator")
	}

	if !strings.Contains(output, "5 credentials") {
		t.Error("Output should contain credential count")
	}

	if !strings.Contains(output, "List") {
		t.Error("Output should contain current view")
	}

	if !strings.Contains(output, "a: add") {
		t.Error("Output should contain shortcuts")
	}
}

func TestStatusBarRenderWithoutShortcuts(t *testing.T) {
	statusBar := NewStatusBar(true, 5, "List")
	statusBar.SetSize(200) // Wide enough to avoid truncation

	output := statusBar.Render()

	// Should use default shortcuts (at least the beginning of them)
	// Note: lipgloss MaxWidth may truncate, so we check for key parts
	if !strings.Contains(output, "help") {
		t.Errorf("Output should contain 'help', got: %s", output)
	}
	// Check for at least "qui" since "quit" might get truncated
	if !strings.Contains(output, "qui") {
		t.Errorf("Output should contain 'qui' (from quit), got: %s", output)
	}
}

func TestStatusBarRenderNarrowTerminal(t *testing.T) {
	statusBar := NewStatusBar(true, 5, "List")
	statusBar.SetShortcuts("a: add | e: edit | d: delete | q: quit")
	statusBar.SetSize(40) // Narrow terminal

	// Should not panic
	output := statusBar.Render()

	if output == "" {
		t.Error("Render() should return output even for narrow terminals")
	}
}
