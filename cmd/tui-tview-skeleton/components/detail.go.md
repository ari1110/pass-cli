# components/detail.go

## Purpose
Credential detail view using tview.TextView. Displays full information for the selected credential with formatting and copy-to-clipboard support.

## Responsibilities

1. **Display credential details**: Service, username, password, category, notes
2. **Password masking**: Show/hide password toggle
3. **Metadata display**: Created date, last modified, last used
4. **Copy to clipboard**: Quick copy password
5. **Formatted output**: Clean, readable layout with sections

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For AppState
- `pass-cli/internal/vault` - For Credential type

### External Dependencies
- `github.com/rivo/tview` - TextView primitive
- `github.com/gdamore/tcell/v2` - For colors
- `github.com/atotto/clipboard` - For clipboard operations

## Key Types

### `DetailView`
**Purpose**: Wrapper around tview.TextView with credential display logic

**Fields**:
```go
type DetailView struct {
    *tview.TextView  // Embedded TextView

    appState       *models.AppState
    passwordMasked bool  // Toggle for password visibility
}
```

## Key Functions

### Constructor

#### `NewDetailView(appState *models.AppState) *DetailView`
**Purpose**: Create and configure detail view

**Steps**:
1. Create tview.TextView
2. Enable dynamic colors (for formatting)
3. Configure scrolling
4. Apply styles (borders, colors)
5. Set initial state (password masked)
6. Return wrapped view

### Content Rendering

#### `Refresh()`
**Purpose**: Rebuild detail view from selected credential

**Steps**:
1. Get selected credential from appState
2. If nil, show "No credential selected" message
3. Otherwise, format credential details
4. Set text with formatted content
5. Scroll to top

**Called when**: Credential selection changes

#### `formatCredential(cred *vault.Credential) string`
**Purpose**: Format credential into readable text

**Layout**:
```
═══════════════════════════════════
          AWS Production
═══════════════════════════════════

Service:    AWS Production
Username:   admin@company.com
Password:   ••••••••••••        [Press 'p' to reveal]
Category:   AWS
URL:        https://aws.amazon.com

Notes:
  Production environment credentials.
  Rotate every 90 days.

═══════════════════════════════════
Metadata
═══════════════════════════════════

Created:     2024-01-15 10:30 AM
Modified:    2024-03-20 02:45 PM
Last Used:   2 hours ago
Usage Count: 47 times

[e] Edit  [d] Delete  [p] Show/Hide Password  [c] Copy Password  [Esc] Back
```

**Formatting**:
- Use tview color tags: `[cyan]`, `[white]`, `[gray]`, `[yellow]`
- Use box drawing characters for sections
- Align labels and values
- Handle multi-line notes

### Password Handling

#### `TogglePasswordVisibility()`
**Purpose**: Show or hide password

**Behavior**:
```go
func (dv *DetailView) TogglePasswordVisibility() {
    dv.passwordMasked = !dv.passwordMasked
    dv.Refresh()  // Redraw with new password state
}
```

**Display**:
- Masked: `••••••••••••`
- Revealed: `ActualP@ssw0rd!`

#### `CopyPasswordToClipboard() error`
**Purpose**: Copy password to clipboard

**Steps**:
1. Get selected credential
2. Extract password
3. Copy to clipboard using atotto/clipboard
4. Show temporary "Copied!" message
5. Clear clipboard after 30 seconds (optional)

**Error handling**: Return error if clipboard unavailable

### Styling

#### `applyStyles()`
**Purpose**: Apply theme colors and borders

**Configuration**:
- Rounded borders
- Border color from theme
- Background color from theme
- Dynamic colors enabled for tags
- Scrollable content
- Word wrap enabled

## Example Structure

```go
package components

import (
    "fmt"
    "strings"
    "time"

    "github.com/atotto/clipboard"
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/models"
    "pass-cli/cmd/tui-tview/styles"
    "pass-cli/internal/vault"
)

type DetailView struct {
    *tview.TextView
    appState       *models.AppState
    passwordMasked bool
}

func NewDetailView(appState *models.AppState) *DetailView {
    textView := tview.NewTextView().
        SetDynamicColors(true).
        SetScrollable(true).
        SetWordWrap(true)

    dv := &DetailView{
        TextView:       textView,
        appState:       appState,
        passwordMasked: true,
    }

    dv.applyStyles()
    dv.Refresh()

    return dv
}

func (dv *DetailView) Refresh() {
    cred := dv.appState.GetSelectedCredential()

    if cred == nil {
        dv.showEmptyState()
        return
    }

    content := dv.formatCredential(cred)
    dv.SetText(content)
    dv.ScrollToBeginning()
}

func (dv *DetailView) formatCredential(cred *vault.Credential) string {
    var b strings.Builder

    // Header
    b.WriteString("[cyan]═══════════════════════════════════[-]\n")
    b.WriteString(fmt.Sprintf("          [yellow]%s[-]\n", cred.Service))
    b.WriteString("[cyan]═══════════════════════════════════[-]\n\n")

    // Main fields
    b.WriteString(fmt.Sprintf("[gray]Service:[-]    [white]%s[-]\n", cred.Service))
    b.WriteString(fmt.Sprintf("[gray]Username:[-]   [white]%s[-]\n", cred.Username))

    // Password with masking
    password := cred.Password
    if dv.passwordMasked {
        password = strings.Repeat("•", len(password))
    }
    hint := ""
    if dv.passwordMasked {
        hint = "  [gray](Press 'p' to reveal)[-]"
    }
    b.WriteString(fmt.Sprintf("[gray]Password:[-]   [white]%s[-]%s\n", password, hint))

    if cred.Category != "" {
        b.WriteString(fmt.Sprintf("[gray]Category:[-]   [white]%s[-]\n", cred.Category))
    }

    if cred.URL != "" {
        b.WriteString(fmt.Sprintf("[gray]URL:[-]        [cyan]%s[-]\n", cred.URL))
    }

    // Notes
    if cred.Notes != "" {
        b.WriteString("\n[gray]Notes:[-]\n")
        b.WriteString(fmt.Sprintf("[white]  %s[-]\n", strings.ReplaceAll(cred.Notes, "\n", "\n  ")))
    }

    // Metadata
    b.WriteString("\n[cyan]═══════════════════════════════════[-]\n")
    b.WriteString("[gray]Metadata[-]\n")
    b.WriteString("[cyan]═══════════════════════════════════[-]\n\n")

    b.WriteString(fmt.Sprintf("[gray]Created:[-]     [white]%s[-]\n", cred.CreatedAt.Format("2006-01-02 03:04 PM")))
    b.WriteString(fmt.Sprintf("[gray]Modified:[-]    [white]%s[-]\n", cred.ModifiedAt.Format("2006-01-02 03:04 PM")))

    if !cred.LastUsed.IsZero() {
        relativeTime := formatRelativeTime(cred.LastUsed)
        b.WriteString(fmt.Sprintf("[gray]Last Used:[-]   [white]%s[-]\n", relativeTime))
    }

    if cred.UsageCount > 0 {
        b.WriteString(fmt.Sprintf("[gray]Usage Count:[-] [white]%d times[-]\n", cred.UsageCount))
    }

    // Keyboard shortcuts
    b.WriteString("\n[gray][e] Edit  [d] Delete  [p] Show/Hide  [c] Copy  [Esc] Back[-]")

    return b.String()
}

func (dv *DetailView) showEmptyState() {
    content := `
[cyan]═══════════════════════════════════[-]

        [gray]No Credential Selected[-]

    Select a credential from the list
    to view its details.

[cyan]═══════════════════════════════════[-]
`
    dv.SetText(content)
}

func (dv *DetailView) TogglePasswordVisibility() {
    dv.passwordMasked = !dv.passwordMasked
    dv.Refresh()
}

func (dv *DetailView) CopyPasswordToClipboard() error {
    cred := dv.appState.GetSelectedCredential()
    if cred == nil {
        return fmt.Errorf("no credential selected")
    }

    err := clipboard.WriteAll(cred.Password)
    if err != nil {
        return fmt.Errorf("failed to copy to clipboard: %w", err)
    }

    // Show temporary "Copied!" message (would need status bar integration)
    return nil
}

func (dv *DetailView) applyStyles() {
    dv.SetBorder(true).
        SetTitle(" Details ").
        SetTitleAlign(tview.AlignLeft).
        SetBorderColor(styles.BorderColor).
        SetBackgroundColor(styles.BackgroundColor)
}

func formatRelativeTime(t time.Time) string {
    duration := time.Since(t)
    switch {
    case duration < time.Minute:
        return "just now"
    case duration < time.Hour:
        return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
    case duration < 24*time.Hour:
        return fmt.Sprintf("%d hours ago", int(duration.Hours()))
    case duration < 7*24*time.Hour:
        return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
    default:
        return t.Format("Jan 2, 2006")
    }
}
```

## Interaction Flow

```
User selects credential in table
    ↓
appState.SetSelectedCredential(cred)
    ↓
AppState notifies selection changed
    ↓
DetailView.Refresh() called
    ↓
Format credential with current mask state
    ↓
Display formatted text
    ↓
User presses 'p' to toggle password
    ↓
DetailView.TogglePasswordVisibility()
    ↓
Refresh with new mask state
```

## Visual Design

### Layout
- Full-width text display
- Scrollable for long content
- Word-wrapped for narrow terminals
- Sections separated by lines

### Colors
- Headers: Cyan
- Labels: Gray
- Values: White
- Service name: Yellow (highlighted)
- URLs: Cyan (link-like)

### Typography
- Box drawing characters for section dividers
- Bullets (•) for masked passwords
- Indentation for notes (preserve formatting)

## Keyboard Shortcuts

Custom shortcuts (handled in event handler):
- **p**: Toggle password visibility
- **c**: Copy password to clipboard
- **e**: Edit credential
- **d**: Delete credential
- **Esc**: Return to list view

## State Integration

### Reads from AppState
- `GetSelectedCredential()` - Current credential

### Writes to AppState
- None (read-only view)

### Responds to AppState callbacks
- `onSelectionChanged` - Triggers Refresh()

## Security Considerations

1. **Password masking**: Default to masked for screen privacy
2. **Clipboard timeout**: Clear clipboard after 30 seconds (optional)
3. **Shoulder surfing**: Large, visible password reveal hint
4. **Screen capture**: Masked by default protects screenshots

## Edge Cases

1. **No credential selected**: Show empty state message
2. **Long passwords**: Handle display without breaking layout
3. **Multi-line notes**: Preserve formatting with indentation
4. **Missing fields**: Only show fields that have values
5. **Clipboard unavailable**: Handle gracefully with error message

## Testing Considerations

- **Test formatting**: Verify layout for various credential types
- **Test masking**: Verify toggle works correctly
- **Test clipboard**: Mock clipboard for testing copy
- **Test empty state**: Verify message displays
- **Mock AppState**: Use mock for isolated testing

## Future Enhancements

- **Syntax highlighting**: Color-code passwords by strength
- **QR code display**: Show QR code for 2FA seeds
- **Edit inline**: Edit fields directly in detail view
- **History view**: Show password change history
- **Export single credential**: Export to file format
- **Share securely**: Generate secure share link
