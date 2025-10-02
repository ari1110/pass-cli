package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors - Using ANSI colors for broad terminal compatibility
	PrimaryColor   = lipgloss.Color("6")   // Cyan (Go brand color)
	SecondaryColor = lipgloss.Color("4")   // Blue
	SuccessColor   = lipgloss.Color("2")   // Green
	WarningColor   = lipgloss.Color("3")   // Yellow/Orange
	ErrorColor     = lipgloss.Color("1")   // Red
	SubtleColor    = lipgloss.Color("240") // Gray (256-color, degrades to white/default)
	TextColor      = lipgloss.Color("15")  // White
	BackgroundDark = lipgloss.Color("235") // Dark gray
	FocusedColor   = lipgloss.Color("6")   // Cyan for focused elements

	// Text styles
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(SecondaryColor)

	LabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(SubtleColor)

	ValueStyle = lipgloss.NewStyle().
		Foreground(TextColor)

	FocusedLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(FocusedColor)

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

	// Password and sensitive data styles
	PasswordStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Background(BackgroundDark).
		Padding(0, 1)

	// Selected/Focused item styles
	SelectedStyle = lipgloss.NewStyle().
		Foreground(TextColor).
		Background(PrimaryColor).
		Bold(true)

	// Border styles
	BorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SubtleColor).
		Padding(1, 2)

	ModalBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2)

	WarningBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(WarningColor).
		Padding(1, 2)

	ErrorBorderStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ErrorColor).
		Padding(1, 2)

	// Notification styles
	NotificationStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true).
		Padding(0, 1)

	ErrorNotificationStyle = lipgloss.NewStyle().
		Foreground(ErrorColor).
		Bold(true).
		Padding(0, 1)

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
		Foreground(SubtleColor).
		Padding(0, 1)

	// Help text styles
	HelpStyle = lipgloss.NewStyle().
		Foreground(SubtleColor)

	KeyStyle = lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true)

	// Table/List styles
	TableHeaderStyle = lipgloss.NewStyle().
		Foreground(PrimaryColor).
		Bold(true)

	TableRowStyle = lipgloss.NewStyle().
		Foreground(TextColor)

	TableDividerStyle = lipgloss.NewStyle().
		Foreground(SubtleColor)
)
