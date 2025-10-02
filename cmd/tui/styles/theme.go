package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("6")   // Cyan (Go brand color)
	SuccessColor   = lipgloss.Color("2")   // Green
	WarningColor   = lipgloss.Color("3")   // Yellow/Orange
	ErrorColor     = lipgloss.Color("1")   // Red
	SubtleColor    = lipgloss.Color("240") // Gray
	TextColor      = lipgloss.Color("15")  // White
	BackgroundDark = lipgloss.Color("235") // Dark gray

	// Text styles
	SuccessStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true)

	WarningStyle = lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true)

	SubtleStyle = lipgloss.NewStyle().
		Foreground(SubtleColor)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
		Foreground(SubtleColor).
		Padding(0, 1)
)
