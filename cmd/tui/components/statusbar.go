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

	// Calculate spacing
	leftLen := lipgloss.Width(left)
	centerLen := lipgloss.Width(center)
	rightLen := lipgloss.Width(right)

	// Total space needed
	totalContent := leftLen + centerLen + rightLen

	// If content fits, distribute evenly
	if totalContent < s.width {
		// Calculate spacing
		totalSpacing := s.width - totalContent
		leftSpacing := totalSpacing / 2
		rightSpacing := totalSpacing - leftSpacing

		return styles.StatusBarStyle.Render(
			left +
				lipgloss.NewStyle().Width(leftSpacing).Render(" ") +
				center +
				lipgloss.NewStyle().Width(rightSpacing).Render(" ") +
				right,
		)
	}

	// If content doesn't fit, prioritize left and right, truncate center if needed
	available := s.width - leftLen - rightLen - 4 // 4 for spacing
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

	spacing := s.width - leftLen - lipgloss.Width(truncatedCenter) - rightLen
	if spacing < 2 {
		spacing = 2
	}

	return styles.StatusBarStyle.Render(
		left +
			lipgloss.NewStyle().Width(spacing).Render(" ") +
			truncatedCenter +
			right,
	)
}
