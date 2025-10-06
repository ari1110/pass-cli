# styles/theme.go

## Purpose
Centralized color palette and styling definitions using tcell colors. Provides consistent theming across all components.

## Responsibilities

1. **Color palette**: Define all colors used in the TUI
2. **Component styles**: Pre-configured styles for common components
3. **Border characters**: Rounded border definitions
4. **Color conversion**: Helper functions for color manipulation
5. **Theme variants**: Support multiple color schemes (future)

## Dependencies

### External Dependencies
- `github.com/gdamore/tcell/v2` - Color types
- `github.com/rivo/tview` - For applying colors to components

## Key Constants

### Color Palette

#### `ColorScheme`
**Purpose**: Struct holding all theme colors

```go
type ColorScheme struct {
    // Background colors
    Background         tcell.Color  // Main background (dark)
    BackgroundLight    tcell.Color  // Slightly lighter background
    BackgroundDark     tcell.Color  // Darker background (for contrast)

    // Border colors
    BorderColor        tcell.Color  // Default border color (cyan)
    BorderInactive     tcell.Color  // Inactive/unfocused borders (gray)

    // Text colors
    TextPrimary        tcell.Color  // Main text (white)
    TextSecondary      tcell.Color  // Secondary text (gray)
    TextAccent         tcell.Color  // Accent text (cyan/yellow)

    // Status colors
    Success            tcell.Color  // Success messages (green)
    Error              tcell.Color  // Error messages (red)
    Warning            tcell.Color  // Warning messages (yellow)
    Info               tcell.Color  // Info messages (cyan)

    // Component-specific
    TableHeader        tcell.Color  // Table header text
    TableSelected      tcell.Color  // Selected row highlight
    SidebarSelected    tcell.Color  // Selected tree node
    StatusBarBg        tcell.Color  // Status bar background
    ButtonBackground   tcell.Color  // Button background
    ButtonText         tcell.Color  // Button text
}
```

### Default Theme (Dracula-inspired)

```go
var DefaultTheme = ColorScheme{
    // Backgrounds
    Background:      tcell.NewRGBColor(40, 42, 54),    // #282a36
    BackgroundLight: tcell.NewRGBColor(68, 71, 90),    // #44475a
    BackgroundDark:  tcell.NewRGBColor(30, 32, 44),    // #1e202c

    // Borders
    BorderColor:     tcell.NewRGBColor(139, 233, 253), // #8be9fd (cyan)
    BorderInactive:  tcell.NewRGBColor(98, 114, 164),  // #6272a4 (gray)

    // Text
    TextPrimary:     tcell.NewRGBColor(248, 248, 242), // #f8f8f2 (white)
    TextSecondary:   tcell.NewRGBColor(98, 114, 164),  // #6272a4 (gray)
    TextAccent:      tcell.NewRGBColor(241, 250, 140), // #f1fa8c (yellow)

    // Status
    Success:         tcell.NewRGBColor(80, 250, 123),  // #50fa7b (green)
    Error:           tcell.NewRGBColor(255, 85, 85),   // #ff5555 (red)
    Warning:         tcell.NewRGBColor(241, 250, 140), // #f1fa8c (yellow)
    Info:            tcell.NewRGBColor(139, 233, 253), // #8be9fd (cyan)

    // Components
    TableHeader:     tcell.NewRGBColor(189, 147, 249), // #bd93f9 (purple)
    TableSelected:   tcell.NewRGBColor(68, 71, 90),    // #44475a (lighter bg)
    SidebarSelected: tcell.NewRGBColor(255, 121, 198), // #ff79c6 (pink)
    StatusBarBg:     tcell.NewRGBColor(30, 32, 44),    // #1e202c (dark)
    ButtonBackground: tcell.NewRGBColor(68, 71, 90),   // #44475a
    ButtonText:      tcell.NewRGBColor(248, 248, 242), // #f8f8f2
}
```

## Key Functions

### Theme Access

#### `GetCurrentTheme() ColorScheme`
**Purpose**: Get active color scheme

**Returns**: Current theme (for now, always DefaultTheme)

#### `SetTheme(theme ColorScheme)`
**Purpose**: Switch to different color scheme

**Future**: Support multiple themes

### Color Helpers

#### `Lighten(color tcell.Color, amount float64) tcell.Color`
**Purpose**: Make color lighter by percentage

**Use case**: Create hover states or highlights

#### `Darken(color tcell.Color, amount float64) tcell.Color`
**Purpose**: Make color darker by percentage

**Use case**: Create shadows or dimmed states

#### `WithAlpha(color tcell.Color, alpha uint8) tcell.Color`
**Purpose**: Set alpha transparency (if supported)

**Note**: tcell doesn't support alpha, but placeholder for future

### Component Styling

#### `ApplyBorderedStyle(primitive tview.Primitive, title string, active bool)`
**Purpose**: Apply consistent border styling to component

**Parameters**:
- `primitive`: Component to style
- `title`: Border title text
- `active`: Whether component is focused

**Effect**:
```go
func ApplyBorderedStyle(p tview.Primitive, title string, active bool) {
    borderColor := DefaultTheme.BorderInactive
    if active {
        borderColor = DefaultTheme.BorderColor
    }

    switch v := p.(type) {
    case *tview.Box:
        v.SetBorder(true).
            SetTitle(" " + title + " ").
            SetTitleAlign(tview.AlignLeft).
            SetBorderColor(borderColor).
            SetBackgroundColor(DefaultTheme.Background)
    }
}
```

#### `ApplyTableStyle(table *tview.Table)`
**Purpose**: Apply consistent table styling

**Styles**:
- Header row with accent color
- Selected row with highlight background
- Alternating row colors (optional)

#### `ApplyFormStyle(form *tview.Form)`
**Purpose**: Apply consistent form styling

**Styles**:
- Input field backgrounds
- Button backgrounds and colors
- Border and title

### Border Characters

#### Rounded Borders

```go
func SetRoundedBorders() {
    tview.Borders.Horizontal = '─'
    tview.Borders.Vertical = '│'
    tview.Borders.TopLeft = '╭'
    tview.Borders.TopRight = '╮'
    tview.Borders.BottomLeft = '╰'
    tview.Borders.BottomRight = '╯'
}
```

#### Box Drawing Characters

```go
const (
    BoxHorizontal = '─'
    BoxVertical   = '│'
    BoxTopLeft    = '╭'
    BoxTopRight   = '╮'
    BoxBottomLeft = '╰'
    BoxBottomRight = '╯'
    BoxDoubleHorizontal = '═'
)
```

## Example Structure

```go
package styles

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

type ColorScheme struct {
    Background      tcell.Color
    BackgroundLight tcell.Color
    BackgroundDark  tcell.Color
    BorderColor     tcell.Color
    BorderInactive  tcell.Color
    TextPrimary     tcell.Color
    TextSecondary   tcell.Color
    TextAccent      tcell.Color
    Success         tcell.Color
    Error           tcell.Color
    Warning         tcell.Color
    Info            tcell.Color
    TableHeader     tcell.Color
    TableSelected   tcell.Color
    SidebarSelected tcell.Color
    StatusBarBg     tcell.Color
    ButtonBackground tcell.Color
    ButtonText      tcell.Color
}

var DefaultTheme = ColorScheme{
    Background:       tcell.NewRGBColor(40, 42, 54),
    BackgroundLight:  tcell.NewRGBColor(68, 71, 90),
    BackgroundDark:   tcell.NewRGBColor(30, 32, 44),
    BorderColor:      tcell.NewRGBColor(139, 233, 253),
    BorderInactive:   tcell.NewRGBColor(98, 114, 164),
    TextPrimary:      tcell.NewRGBColor(248, 248, 242),
    TextSecondary:    tcell.NewRGBColor(98, 114, 164),
    TextAccent:       tcell.NewRGBColor(241, 250, 140),
    Success:          tcell.NewRGBColor(80, 250, 123),
    Error:            tcell.NewRGBColor(255, 85, 85),
    Warning:          tcell.NewRGBColor(241, 250, 140),
    Info:             tcell.NewRGBColor(139, 233, 253),
    TableHeader:      tcell.NewRGBColor(189, 147, 249),
    TableSelected:    tcell.NewRGBColor(68, 71, 90),
    SidebarSelected:  tcell.NewRGBColor(255, 121, 198),
    StatusBarBg:      tcell.NewRGBColor(30, 32, 44),
    ButtonBackground: tcell.NewRGBColor(68, 71, 90),
    ButtonText:       tcell.NewRGBColor(248, 248, 242),
}

func GetCurrentTheme() ColorScheme {
    return DefaultTheme
}

func SetRoundedBorders() {
    tview.Borders.Horizontal = '─'
    tview.Borders.Vertical = '│'
    tview.Borders.TopLeft = '╭'
    tview.Borders.TopRight = '╮'
    tview.Borders.BottomLeft = '╰'
    tview.Borders.BottomRight = '╯'
}

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
    }
}

func ApplyTableStyle(table *tview.Table) {
    theme := GetCurrentTheme()

    table.SetBackgroundColor(theme.Background).
        SetSelectedStyle(tcell.StyleDefault.
            Background(theme.TableSelected).
            Foreground(theme.TextPrimary).
            Bold(true))
}

func ApplyFormStyle(form *tview.Form) {
    theme := GetCurrentTheme()

    form.SetBackgroundColor(theme.Background).
        SetButtonBackgroundColor(theme.ButtonBackground).
        SetButtonTextColor(theme.ButtonText).
        SetLabelColor(theme.TextSecondary).
        SetFieldBackgroundColor(theme.BackgroundLight).
        SetFieldTextColor(theme.TextPrimary)
}

func Lighten(color tcell.Color, amount float64) tcell.Color {
    r, g, b := color.RGB()
    r = uint8(float64(r) * (1.0 + amount))
    g = uint8(float64(g) * (1.0 + amount))
    b = uint8(float64(b) * (1.0 + amount))
    return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}

func Darken(color tcell.Color, amount float64) tcell.Color {
    r, g, b := color.RGB()
    r = uint8(float64(r) * (1.0 - amount))
    g = uint8(float64(g) * (1.0 - amount))
    b = uint8(float64(b) * (1.0 - amount))
    return tcell.NewRGBColor(int32(r), int32(g), int32(b))
}
```

## Color Palette Visualization

```
Backgrounds:
  Background      ███ #282a36 (dark gray)
  BackgroundLight ███ #44475a (medium gray)
  BackgroundDark  ███ #1e202c (darker gray)

Borders:
  BorderColor     ███ #8be9fd (cyan)
  BorderInactive  ███ #6272a4 (gray)

Text:
  TextPrimary     ███ #f8f8f2 (white)
  TextSecondary   ███ #6272a4 (gray)
  TextAccent      ███ #f1fa8c (yellow)

Status:
  Success         ███ #50fa7b (green)
  Error           ███ #ff5555 (red)
  Warning         ███ #f1fa8c (yellow)
  Info            ███ #8be9fd (cyan)

Components:
  TableHeader     ███ #bd93f9 (purple)
  SidebarSelected ███ #ff79c6 (pink)
```

## Usage Examples

### Apply Theme to Component
```go
// In component constructor:
sidebar := tview.NewTreeView()
styles.ApplyBorderedStyle(sidebar, "Categories", true)
styles.SetRoundedBorders()
```

### Use Theme Colors Directly
```go
theme := styles.GetCurrentTheme()

textView.SetTextColor(theme.TextPrimary).
    SetBackgroundColor(theme.Background)

// For status messages:
statusBar.SetText(fmt.Sprintf("[%s]Success![-]", theme.Success.Hex()))
```

### Switch Active State
```go
// When component gains focus:
styles.ApplyBorderedStyle(component, "Title", true)  // Active

// When component loses focus:
styles.ApplyBorderedStyle(component, "Title", false) // Inactive
```

## Alternative Themes (Future)

### Gruvbox Theme
```go
var GruvboxTheme = ColorScheme{
    Background:  tcell.NewRGBColor(40, 40, 40),   // #282828
    // ... more colors
}
```

### Solarized Theme
```go
var SolarizedTheme = ColorScheme{
    Background:  tcell.NewRGBColor(0, 43, 54),    // #002b36
    // ... more colors
}
```

## Terminal Compatibility

### Color Support Levels
- **True Color (24-bit)**: Full RGB support (Windows Terminal, iTerm2)
- **256 Colors**: Good approximation (most modern terminals)
- **16 Colors**: Basic fallback (rare, older terminals)

tview/tcell automatically detects and uses best available

### Testing Colors
Test in multiple terminals to ensure consistent appearance:
- Windows Terminal
- iTerm2 (macOS)
- gnome-terminal (Linux)
- Alacritty
- Kitty

## Future Enhancements

- **Theme selection**: User-selectable themes via config
- **Custom themes**: User-defined color schemes
- **Light mode**: Light background theme variant
- **High contrast**: Accessibility-focused theme
- **Color preview**: TUI for previewing theme changes
- **Theme export/import**: Share themes via JSON
- **Dynamic theming**: Change theme without restart
