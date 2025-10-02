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
	// Account for status bar padding (2 chars total from padding 0, 1)
	availableWidth := s.width - 2
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

	// Total space needed
	totalContent := leftLen + centerLen + rightLen

	// If content fits, distribute evenly
	if totalContent+4 < availableWidth { // +4 for minimum spacing
		// Calculate spacing
		totalSpacing := availableWidth - totalContent
		leftSpacing := totalSpacing / 2
		rightSpacing := totalSpacing - leftSpacing

		content := left +
			lipgloss.NewStyle().Width(leftSpacing).Render("") +
			center +
			lipgloss.NewStyle().Width(rightSpacing).Render("") +
			right

		// Ensure we don't exceed width by using MaxWidth
		return lipgloss.NewStyle().
			MaxWidth(availableWidth).
			Foreground(styles.SubtleColor).
			Padding(0, 1).
			Render(content)
	}

	// If content doesn't fit, prioritize left and right, truncate center
	available := availableWidth - leftLen - rightLen - 4 // 4 for spacing
	if available < 0 {
		available = 0
	}

	truncatedCenter := center
	if centerLen > available {
		if available > 3 {
			truncatedCenter = center[:available-3] + "..."
		} else {
			truncatedCenter = ""
		}
	}

	spacing := availableWidth - leftLen - lipgloss.Width(truncatedCenter) - rightLen
	if spacing < 2 {
		spacing = 2
	}

	content := left +
		lipgloss.NewStyle().Width(spacing).Render("") +
		truncatedCenter +
		right

	// Ensure we don't exceed width by using MaxWidth
	return lipgloss.NewStyle().
		MaxWidth(availableWidth).
		Foreground(styles.SubtleColor).
		Padding(0, 1).
		Render(content)
}
