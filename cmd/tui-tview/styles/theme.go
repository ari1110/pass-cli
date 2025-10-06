package styles

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ColorScheme defines the color palette for the TUI.
type ColorScheme struct {
	// Background colors
	Background      tcell.Color // Main background (dark)
	BackgroundLight tcell.Color // Slightly lighter background
	BackgroundDark  tcell.Color // Darker background (for contrast)

	// Border colors
	BorderColor    tcell.Color // Default border color (cyan)
	BorderInactive tcell.Color // Inactive/unfocused borders (gray)

	// Text colors
	TextPrimary   tcell.Color // Main text (white)
	TextSecondary tcell.Color // Secondary text (gray)
	TextAccent    tcell.Color // Accent text (cyan/yellow)

	// Status colors
	Success tcell.Color // Success messages (green)
	Error   tcell.Color // Error messages (red)
	Warning tcell.Color // Warning messages (yellow)
	Info    tcell.Color // Info messages (cyan)

	// Component-specific
	TableHeader      tcell.Color // Table header text
	TableSelected    tcell.Color // Selected row highlight
	SidebarSelected  tcell.Color // Selected tree node
	StatusBarBg      tcell.Color // Status bar background
	ButtonBackground tcell.Color // Button background
	ButtonText       tcell.Color // Button text
}

// DefaultTheme provides a Dracula-inspired color scheme.
var DefaultTheme = ColorScheme{
	// Backgrounds
	Background:      tcell.NewRGBColor(40, 42, 54),   // #282a36
	BackgroundLight: tcell.NewRGBColor(68, 71, 90),   // #44475a
	BackgroundDark:  tcell.NewRGBColor(30, 32, 44),   // #1e202c

	// Borders
	BorderColor:    tcell.NewRGBColor(139, 233, 253), // #8be9fd (cyan)
	BorderInactive: tcell.NewRGBColor(98, 114, 164),  // #6272a4 (gray)

	// Text
	TextPrimary:   tcell.NewRGBColor(248, 248, 242), // #f8f8f2 (white)
	TextSecondary: tcell.NewRGBColor(98, 114, 164),  // #6272a4 (gray)
	TextAccent:    tcell.NewRGBColor(241, 250, 140), // #f1fa8c (yellow)

	// Status
	Success: tcell.NewRGBColor(80, 250, 123),  // #50fa7b (green)
	Error:   tcell.NewRGBColor(255, 85, 85),   // #ff5555 (red)
	Warning: tcell.NewRGBColor(241, 250, 140), // #f1fa8c (yellow)
	Info:    tcell.NewRGBColor(139, 233, 253), // #8be9fd (cyan)

	// Components
	TableHeader:      tcell.NewRGBColor(189, 147, 249), // #bd93f9 (purple)
	TableSelected:    tcell.NewRGBColor(68, 71, 90),    // #44475a (lighter bg)
	SidebarSelected:  tcell.NewRGBColor(255, 121, 198), // #ff79c6 (pink)
	StatusBarBg:      tcell.NewRGBColor(30, 32, 44),    // #1e202c (dark)
	ButtonBackground: tcell.NewRGBColor(68, 71, 90),    // #44475a
	ButtonText:       tcell.NewRGBColor(248, 248, 242), // #f8f8f2
}

// GetCurrentTheme returns the currently active color scheme.
func GetCurrentTheme() ColorScheme {
	return DefaultTheme
}

// SetRoundedBorders configures tview to use rounded border characters.
func SetRoundedBorders() {
	tview.Borders.Horizontal = '─'
	tview.Borders.Vertical = '│'
	tview.Borders.TopLeft = '╭'
	tview.Borders.TopRight = '╮'
	tview.Borders.BottomLeft = '╰'
	tview.Borders.BottomRight = '╯'
}

// ApplyBorderedStyle applies consistent border styling to a component.
// Uses type switch to handle all tview primitive types.
func ApplyBorderedStyle(p tview.Primitive, title string, active bool) {
	theme := GetCurrentTheme()
	borderColor := theme.BorderInactive
	if active {
		borderColor = theme.BorderColor
	}

	switch v := p.(type) {
	case *tview.Box:
		v.SetBorder(true).
			SetTitle(" " + title + " ").
			SetTitleAlign(tview.AlignLeft).
			SetBorderColor(borderColor).
			SetBackgroundColor(theme.Background)

	case *tview.Table:
		v.SetBorder(true).
			SetTitle(" " + title + " ").
			SetTitleAlign(tview.AlignLeft).
			SetBorderColor(borderColor).
			SetBackgroundColor(theme.Background)

	case *tview.TreeView:
		v.SetBorder(true).
			SetTitle(" " + title + " ").
			SetTitleAlign(tview.AlignLeft).
			SetBorderColor(borderColor).
			SetBackgroundColor(theme.Background)

	case *tview.TextView:
		v.SetBorder(true).
			SetTitle(" " + title + " ").
			SetTitleAlign(tview.AlignLeft).
			SetBorderColor(borderColor).
			SetBackgroundColor(theme.Background)

	case *tview.Form:
		v.SetBorder(true).
			SetTitle(" " + title + " ").
			SetTitleAlign(tview.AlignLeft).
			SetBorderColor(borderColor).
			SetBackgroundColor(theme.Background)

	case *tview.Modal:
		v.SetBorder(true).
			SetTitle(" " + title + " ").
			SetTitleAlign(tview.AlignLeft).
			SetBorderColor(borderColor).
			SetBackgroundColor(theme.Background)

	case *tview.List:
		v.SetBorder(true).
			SetTitle(" " + title + " ").
			SetTitleAlign(tview.AlignLeft).
			SetBorderColor(borderColor).
			SetBackgroundColor(theme.Background)
	}
}

// ApplyTableStyle applies consistent styling to table components.
func ApplyTableStyle(table *tview.Table) {
	theme := GetCurrentTheme()

	table.SetBackgroundColor(theme.Background)
	table.SetSelectedStyle(tcell.StyleDefault.
		Background(theme.TableSelected).
		Foreground(theme.TextPrimary).
		Bold(true))
}

// ApplyFormStyle applies consistent styling to form components.
func ApplyFormStyle(form *tview.Form) {
	theme := GetCurrentTheme()

	form.SetBackgroundColor(theme.Background)
	form.SetButtonBackgroundColor(theme.ButtonBackground)
	form.SetButtonTextColor(theme.ButtonText)
	form.SetLabelColor(theme.TextSecondary)
	form.SetFieldBackgroundColor(theme.BackgroundLight)
	form.SetFieldTextColor(theme.TextPrimary)
}

// Lighten makes a color lighter by the given percentage.
func Lighten(color tcell.Color, amount float64) tcell.Color {
	r, g, b := color.RGB()
	r = clampUint8(float64(r) * (1.0 + amount))
	g = clampUint8(float64(g) * (1.0 + amount))
	b = clampUint8(float64(b) * (1.0 + amount))
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

// Darken makes a color darker by the given percentage.
func Darken(color tcell.Color, amount float64) tcell.Color {
	r, g, b := color.RGB()
	r = clampUint8(float64(r) * (1.0 - amount))
	g = clampUint8(float64(g) * (1.0 - amount))
	b = clampUint8(float64(b) * (1.0 - amount))
	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

// clampUint8 clamps a float64 value to the uint8 range [0, 255].
func clampUint8(v float64) int32 {
	if v > 255 {
		return 255
	}
	if v < 0 {
		return 0
	}
	return int32(v)
}
