# components/statusbar.go

## Purpose
Context-aware status bar using tview.TextView. Displays relevant keyboard shortcuts and hints based on current focus.

## Responsibilities

1. **Keyboard shortcut display**: Show context-appropriate shortcuts
2. **Status messages**: Display temporary messages (copied, saved, error)
3. **Focus-aware**: Update shortcuts based on focused component
4. **Minimal height**: Single line at bottom of screen

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For NavigationState

### External Dependencies
- `github.com/rivo/tview` - TextView primitive
- `github.com/gdamore/tcell/v2` - For colors

## Key Types

### `StatusBar`
**Purpose**: Wrapper around tview.TextView with context-aware content

**Fields**:
```go
type StatusBar struct {
    *tview.TextView  // Embedded TextView

    currentFocus models.FocusableComponent
    messageTimer *time.Timer  // For temporary messages
}
```

## Key Functions

### Constructor

#### `NewStatusBar() *StatusBar`
**Purpose**: Create and configure status bar

**Steps**:
1. Create tview.TextView
2. Set dynamic colors enabled
3. Configure styling (no borders, background color)
4. Set initial shortcuts
5. Return wrapped status bar

### Content Updates

#### `UpdateForFocus(focus models.FocusableComponent)`
**Purpose**: Update shortcuts based on focused component

**Content by focus**:
- **Sidebar**: `[Tab] Next | [↑↓] Navigate | [Enter] Select | [n] New | [?] Help | [q] Quit`
- **Table**: `[Tab] Next | [↑↓] Navigate | [Enter] View | [n] New | [e] Edit | [d] Delete | [q] Quit`
- **Detail**: `[Tab] Next | [e] Edit | [d] Delete | [p] Toggle Password | [c] Copy | [Esc] Back | [q] Quit`
- **Form**: `[Tab] Next Field | [Enter] Submit | [Esc] Cancel`

#### `ShowMessage(message string, duration time.Duration)`
**Purpose**: Display temporary status message

**Behavior**:
1. Clear current content
2. Show message with formatting
3. After duration, restore context shortcuts
4. Cancel previous timer if exists

**Examples**:
- `ShowMessage("[green]Password copied to clipboard![-]", 3*time.Second)`
- `ShowMessage("[yellow]Credential saved successfully[-]", 2*time.Second)`
- `ShowMessage("[red]Error: Failed to delete credential[-]", 5*time.Second)`

#### `ShowError(err error)`
**Purpose**: Display error message

**Convenience wrapper**: Calls `ShowMessage` with red color and 5-second duration

### Styling

#### `applyStyles()`
**Purpose**: Apply theme colors

**Configuration**:
- No border (clean look)
- Background color from theme
- Text centered
- Fixed height (1 row)

## Example Structure

```go
package components

import (
    "fmt"
    "time"

    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/models"
    "pass-cli/cmd/tui-tview/styles"
)

type StatusBar struct {
    *tview.TextView
    currentFocus models.FocusableComponent
    messageTimer *time.Timer
}

func NewStatusBar() *StatusBar {
    textView := tview.NewTextView().
        SetDynamicColors(true).
        SetTextAlign(tview.AlignCenter)

    sb := &StatusBar{
        TextView: textView,
    }

    sb.applyStyles()
    sb.UpdateForFocus(models.FocusSidebar)  // Default

    return sb
}

func (sb *StatusBar) UpdateForFocus(focus models.FocusableComponent) {
    sb.currentFocus = focus

    var shortcuts string
    switch focus {
    case models.FocusSidebar:
        shortcuts = "[gray][Tab] Next | [↑↓] Navigate | [Enter] Select | [n] New | [?] Help | [q] Quit[-]"

    case models.FocusTable:
        shortcuts = "[gray][Tab] Next | [↑↓] Navigate | [Enter] View | [n] New | [e] Edit | [d] Delete | [q] Quit[-]"

    case models.FocusDetail:
        shortcuts = "[gray][Tab] Next | [e] Edit | [d] Delete | [p] Show/Hide | [c] Copy | [Esc] Back | [q] Quit[-]"

    case models.FocusForm:
        shortcuts = "[gray][Tab] Next Field | [Enter] Submit | [Esc] Cancel[-]"

    default:
        shortcuts = "[gray][q] Quit | [?] Help[-]"
    }

    sb.SetText(shortcuts)
}

func (sb *StatusBar) ShowMessage(message string, duration time.Duration) {
    // Cancel previous message timer
    if sb.messageTimer != nil {
        sb.messageTimer.Stop()
    }

    // Show message
    sb.SetText(message)

    // Restore shortcuts after duration
    sb.messageTimer = time.AfterFunc(duration, func() {
        sb.UpdateForFocus(sb.currentFocus)
    })
}

func (sb *StatusBar) ShowError(err error) {
    message := fmt.Sprintf("[red]Error: %s[-]", err.Error())
    sb.ShowMessage(message, 5*time.Second)
}

func (sb *StatusBar) ShowSuccess(message string) {
    formatted := fmt.Sprintf("[green]%s[-]", message)
    sb.ShowMessage(formatted, 3*time.Second)
}

func (sb *StatusBar) applyStyles() {
    sb.SetBackgroundColor(styles.StatusBarBackground)
    sb.SetBorder(false)
}
```

## Interaction Flow

```
Focus changes to Table
    ↓
NavigationState calls onFocusChanged callback
    ↓
StatusBar.UpdateForFocus(FocusTable)
    ↓
Display table-specific shortcuts
    ↓
User performs action (e.g., copy password)
    ↓
StatusBar.ShowSuccess("Password copied!")
    ↓
Message displays for 3 seconds
    ↓
Shortcuts automatically restore
```

## Visual Design

### Layout
- Single line at bottom of screen
- Full width
- No borders
- Centered text

### Colors
- Default shortcuts: Gray text
- Success messages: Green text
- Error messages: Red text
- Warning messages: Yellow text
- Background: Darker than main background

### Content Format
```
[gray][Tab] Next | [↑↓] Navigate | [Enter] View | [q] Quit[-]
```

Pipe (|) separators between shortcut groups

## Integration Points

### With NavigationState
```go
// In main setup:
nav.SetOnFocusChanged(func(focus models.FocusableComponent) {
    statusBar.UpdateForFocus(focus)
})
```

### With Components
```go
// After copying password:
err := detailView.CopyPasswordToClipboard()
if err != nil {
    statusBar.ShowError(err)
} else {
    statusBar.ShowSuccess("Password copied to clipboard!")
}
```

### With Forms
```go
// After saving credential:
err := appState.AddCredential(service, username, password)
if err != nil {
    statusBar.ShowError(err)
} else {
    statusBar.ShowSuccess("Credential saved successfully!")
}
```

## Keyboard Shortcuts by Context

### Sidebar Focus
- Tab: Move to next component
- ↑↓: Navigate categories
- Enter: Select category
- n: New credential
- ?: Help
- q: Quit

### Table Focus
- Tab: Move to next component
- ↑↓: Navigate credentials
- Enter: View details
- n: New credential
- e: Edit selected
- d: Delete selected
- q: Quit

### Detail Focus
- Tab: Move to next component
- e: Edit credential
- d: Delete credential
- p: Toggle password visibility
- c: Copy password
- Esc: Back to list
- q: Quit

### Form Focus
- Tab: Next form field
- Enter: Submit form
- Esc: Cancel and close

## Message Types

### Success Messages
- Green color
- 2-3 second duration
- Examples:
  - "Credential saved successfully!"
  - "Password copied to clipboard!"
  - "Credential deleted"

### Error Messages
- Red color
- 5 second duration (longer to read)
- Examples:
  - "Error: Failed to save credential"
  - "Error: Invalid password format"
  - "Error: Clipboard unavailable"

### Warning Messages
- Yellow color
- 4 second duration
- Examples:
  - "Warning: This credential is used in 3 projects"
  - "Warning: Password strength is weak"

### Info Messages
- Cyan color
- 3 second duration
- Examples:
  - "Refreshing credentials..."
  - "Loading vault..."

## Edge Cases

1. **Rapid messages**: Cancel previous timer, show new message
2. **Long error messages**: Truncate if too long
3. **Component not focused**: Show default shortcuts
4. **Modal open**: Show modal-specific shortcuts

## Testing Considerations

- **Test context switching**: Verify shortcuts update correctly
- **Test message timing**: Verify messages disappear after duration
- **Test message cancellation**: Verify rapid messages work
- **Mock timer**: Use mock time for testing
- **Visual testing**: Manual verification of display

## Future Enhancements

- **Progress indicator**: Show loading progress
- **Multiple zones**: Left (context), center (message), right (status)
- **Scrolling messages**: For very long messages
- **Customizable shortcuts**: User-defined key bindings
- **Help tooltip**: Hover for detailed shortcut explanations
- **Animation**: Fade in/out for messages
