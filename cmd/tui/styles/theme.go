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

	// Panel styles for dashboard
	ActivePanelBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(PrimaryColor). // Cyan border for active panel
				Padding(0, 1)

	InactivePanelBorderStyle = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(SubtleColor). // Gray border for inactive panel
					Padding(0, 1)

	// Panel title styles
	PanelTitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	InactivePanelTitleStyle = lipgloss.NewStyle().
				Foreground(SubtleColor).
				Bold(true).
				Padding(0, 1)

	// Breadcrumb style
	BreadcrumbStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	BreadcrumbSeparatorStyle = lipgloss.NewStyle().
					Foreground(SubtleColor)
)

// CategoryIcons maps category types to their display icons
var CategoryIcons = map[string]string{
	"APIs & Services":      "🔑",
	"Cloud Infrastructure": "☁️",
	"Databases":            "💾",
	"Version Control":      "📦",
	"Communication":        "📧",
	"Payment Processing":   "💰",
	"AI Services":          "🤖",
	"Uncategorized":        "📁",
}

// CategoryIconsASCII provides ASCII fallback for terminals without emoji support
var CategoryIconsASCII = map[string]string{
	"APIs & Services":      "[API]",
	"Cloud Infrastructure": "[CLD]",
	"Databases":            "[DB]",
	"Version Control":      "[GIT]",
	"Communication":        "[MSG]",
	"Payment Processing":   "[PAY]",
	"AI Services":          "[AI]",
	"Uncategorized":        "[???]",
}

// StatusIcons maps status types to their display icons
var StatusIcons = map[string]string{
	"pending":   "⏳",
	"running":   "⏳",
	"success":   "✓",
	"failed":    "✗",
	"collapsed": "▶",
	"expanded":  "▼",
}

// StatusIconsASCII provides ASCII fallback for status icons
var StatusIconsASCII = map[string]string{
	"pending":   "[.]",
	"running":   "[.]",
	"success":   "[+]",
	"failed":    "[X]",
	"collapsed": "[>]",
	"expanded":  "[v]",
}

// SupportsUnicode checks if the terminal supports Unicode icons
// This is a simple heuristic - you can enhance it based on TERM environment variable
var SupportsUnicode = true // Default to true, can be overridden based on terminal detection

// GetCategoryIcon returns the appropriate category icon based on Unicode support
func GetCategoryIcon(category string) string {
	if SupportsUnicode {
		if icon, ok := CategoryIcons[category]; ok {
			return icon
		}
		return CategoryIcons["Uncategorized"]
	}

	if icon, ok := CategoryIconsASCII[category]; ok {
		return icon
	}
	return CategoryIconsASCII["Uncategorized"]
}

// GetStatusIcon returns the appropriate status icon based on Unicode support
func GetStatusIcon(status string) string {
	if SupportsUnicode {
		if icon, ok := StatusIcons[status]; ok {
			return icon
		}
		return ""
	}

	if icon, ok := StatusIconsASCII[status]; ok {
		return icon
	}
	return ""
}
