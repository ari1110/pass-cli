package components

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"pass-cli/cmd/tui/styles"
)

// StatusBar displays system status and keyboard shortcuts
type StatusBar struct {
	keychainAvailable bool
	credentialCount   int
	currentView       string
	shortcuts         string
	width             int
}

// NewStatusBar creates a new status bar
func NewStatusBar(keychainAvailable bool, credentialCount int, currentView string) *StatusBar {
	return &StatusBar{
		keychainAvailable: keychainAvailable,
		credentialCount:   credentialCount,
		currentView:       currentView,
	}
}

// SetSize updates the status bar width
func (s *StatusBar) SetSize(width int) {
	s.width = width
}

// SetKeychainStatus updates keychain availability
func (s *StatusBar) SetKeychainStatus(available bool) {
	s.keychainAvailable = available
}

// SetCredentialCount updates the credential count
func (s *StatusBar) SetCredentialCount(count int) {
	s.credentialCount = count
}

// SetCurrentView updates the current view name
func (s *StatusBar) SetCurrentView(view string) {
	s.currentView = view
}

// SetShortcuts sets the keyboard shortcuts to display
func (s *StatusBar) SetShortcuts(shortcuts string) {
	s.shortcuts = shortcuts
}

// Render returns the rendered status bar
func (s *StatusBar) Render() string {
	// Use actual terminal width without adjustment
	availableWidth := s.width
	if availableWidth < 10 {
		availableWidth = 10
	}

	// Left side: keychain status and credential count
	var keychainIndicator string
	if s.keychainAvailable {
		keychainIndicator = styles.SuccessStyle.Render("🔓 Keychain")
	} else {
		keychainIndicator = styles.WarningStyle.Render("🔒 Password")
	}

	credCount := fmt.Sprintf("%d credential", s.credentialCount)
	if s.credentialCount != 1 {
		credCount += "s"
	}

	left := fmt.Sprintf("%s  %s", keychainIndicator, credCount)

	// Center: current view
	center := s.currentView

	// Right side: shortcuts
	right := s.shortcuts
	if right == "" {
		right = "?: help | q: quit"
	}

	// Calculate lengths (accounting for styled characters)
	leftLen := lipgloss.Width(left)
	centerLen := lipgloss.Width(center)
	rightLen := lipgloss.Width(right)

	// Reserve space for padding (2 chars)
	contentWidth := availableWidth - 2

	// If content fits comfortably, distribute evenly
	totalContent := leftLen + centerLen + rightLen
	if totalContent+4 <= contentWidth { // +4 for minimum spacing
		// Calculate spacing
		totalSpacing := contentWidth - totalContent
		leftSpacing := totalSpacing / 2
		rightSpacing := totalSpacing - leftSpacing

		content := left +
			lipgloss.NewStyle().Width(leftSpacing).Render("") +
			center +
			lipgloss.NewStyle().Width(rightSpacing).Render("") +
			right

		// Use MaxWidth to prevent overflow
		return lipgloss.NewStyle().
			MaxWidth(availableWidth).
			Foreground(styles.SubtleColor).
			Padding(0, 1).
			Render(content)
	}

	// Content doesn't fit - need to truncate intelligently
	// Priority: left (fixed) > right (truncate) > center (truncate)

	// Reserve space for left + minimum spacing
	remainingWidth := contentWidth - leftLen - 2 // 2 for spacing

	// Truncate right side if needed (shortcuts)
	truncatedRight := right
	maxRightWidth := remainingWidth - centerLen - 2 // 2 for spacing
	if rightLen > maxRightWidth && maxRightWidth > 10 {
		// Truncate shortcuts intelligently - keep as many complete shortcuts as possible
		truncatedRight = s.truncateShortcuts(right, maxRightWidth)
		rightLen = lipgloss.Width(truncatedRight)
	}

	// Recalculate remaining width for center
	remainingForCenter := contentWidth - leftLen - rightLen - 4 // 4 for spacing

	// Truncate center if needed
	truncatedCenter := center
	if centerLen > remainingForCenter && remainingForCenter > 0 {
		if remainingForCenter > 3 {
			// Use lipgloss truncation for proper handling
			truncatedCenter = lipgloss.NewStyle().MaxWidth(remainingForCenter - 3).Render(center) + "..."
		} else {
			truncatedCenter = ""
		}
	}

	// Calculate final spacing
	finalLeftLen := lipgloss.Width(left)
	finalCenterLen := lipgloss.Width(truncatedCenter)
	finalRightLen := lipgloss.Width(truncatedRight)
	spacing := contentWidth - finalLeftLen - finalCenterLen - finalRightLen
	if spacing < 2 {
		spacing = 2
	}

	content := left +
		lipgloss.NewStyle().Width(spacing).Render("") +
		truncatedCenter +
		truncatedRight

	// Use MaxWidth to prevent overflow
	return lipgloss.NewStyle().
		MaxWidth(availableWidth).
		Foreground(styles.SubtleColor).
		Padding(0, 1).
		Render(content)
}

// truncateShortcuts intelligently truncates shortcuts to fit width
// Tries to keep complete shortcut entries (e.g., "q: quit") rather than cutting mid-word
func (s *StatusBar) truncateShortcuts(shortcuts string, maxWidth int) string {
	if lipgloss.Width(shortcuts) <= maxWidth {
		return shortcuts
	}

	// Split by " | " separator
	parts := []string{}
	current := ""
	for i, r := range shortcuts {
		current += string(r)
		if i < len(shortcuts)-3 && shortcuts[i:i+3] == " | " {
			parts = append(parts, current[:len(current)-3])
			current = ""
			i += 2 // Skip past " | "
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	// Add parts until we exceed width
	result := ""
	for _, part := range parts {
		test := result
		if test != "" {
			test += " | "
		}
		test += part

		if lipgloss.Width(test) > maxWidth-3 { // -3 for "..."
			break
		}
		result = test
	}

	if result == "" {
		// If even one shortcut doesn't fit, just truncate hard
		if maxWidth > 3 {
			return lipgloss.NewStyle().MaxWidth(maxWidth - 3).Render(shortcuts) + "..."
		}
		return ""
	}

	return result
}
