# events/handlers.go

## Purpose
Global keyboard shortcut handling with focus-aware input protection. Prevents shortcuts from interfering with form input.

## Responsibilities

1. **Global shortcuts**: Handle app-wide keyboard shortcuts (q, n, e, d, ?, etc.)
2. **Focus protection**: Don't intercept input for Form, InputField, TextArea
3. **Context-aware actions**: Execute different actions based on current view/focus
4. **Modifier keys**: Handle Ctrl+, Alt+, Shift+ combinations
5. **Help system**: Display keyboard shortcut help

## CRITICAL: Form Input Protection

**The Problem**: Previous implementation intercepted keys like 'e', 'n', 'd' globally, preventing users from typing these characters in forms.

**The Solution**: Check focused component type before intercepting shortcuts.

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For AppState, NavigationState
- `pass-cli/cmd/tui-tview/components` - For accessing components
- `pass-cli/cmd/tui-tview/layout` - For PageManager

### External Dependencies
- `github.com/rivo/tview` - Application
- `github.com/gdamore/tcell/v2` - Event types

## Key Types

### `EventHandler`
**Purpose**: Centralized keyboard event management

**Fields**:
```go
type EventHandler struct {
    app         *tview.Application
    appState    *models.AppState
    nav         *models.NavigationState
    pageManager *layout.PageManager
    statusBar   *components.StatusBar
}
```

## Key Functions

### Constructor

#### `NewEventHandler(app *tview.Application, appState *models.AppState, nav *models.NavigationState, pageManager *layout.PageManager, statusBar *components.StatusBar) *EventHandler`
**Purpose**: Create event handler with all dependencies

### Setup

#### `SetupGlobalShortcuts()`
**Purpose**: Install global keyboard shortcut handler

**CRITICAL IMPLEMENTATION**:
```go
func (eh *EventHandler) SetupGlobalShortcuts() {
    eh.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        // ✅ CRITICAL: Check if focused component should handle input
        focused := eh.app.GetFocus()
        if focused != nil {
            switch focused.(type) {
            case *tview.Form, *tview.InputField, *tview.TextArea:
                // ✅ Let input components handle their own keys
                // Only intercept Ctrl+C for quit
                if event.Key() == tcell.KeyCtrlC {
                    eh.handleQuit()
                    return nil
                }
                return event  // ✅ Pass all other keys to input component
            }
        }

        // Handle global shortcuts
        return eh.handleGlobalKey(event)
    })
}
```

### Shortcut Handlers

#### `handleGlobalKey(event *tcell.EventKey) *tcell.EventKey`
**Purpose**: Route keyboard events to appropriate handlers

**Shortcuts**:
```go
switch event.Key() {
case tcell.KeyRune:
    switch event.Rune() {
    case 'q':
        eh.handleQuit()
        return nil
    case 'n':
        eh.handleNewCredential()
        return nil
    case 'e':
        eh.handleEditCredential()
        return nil
    case 'd':
        eh.handleDeleteCredential()
        return nil
    case 'p':
        eh.handleTogglePassword()
        return nil
    case 'c':
        eh.handleCopyPassword()
        return nil
    case '?':
        eh.handleShowHelp()
        return nil
    }

case tcell.KeyTab:
    eh.handleTabFocus()
    return nil

case tcell.KeyBacktab:  // Shift+Tab
    eh.handleShiftTabFocus()
    return nil

case tcell.KeyCtrlC:
    eh.handleQuit()
    return nil
}

return event  // Pass through unhandled keys
```

### Action Handlers

#### `handleQuit()`
**Purpose**: Quit application

**Behavior**:
- If modal open: Close modal instead
- Otherwise: Confirm quit (optional) and exit

#### `handleNewCredential()`
**Purpose**: Show add credential form

**Steps**:
1. Create NewAddForm
2. Set callbacks for submit/cancel
3. Show via pageManager
4. Update status bar

#### `handleEditCredential()`
**Purpose**: Edit selected credential

**Steps**:
1. Get selected credential from appState
2. If none selected, show error
3. Create NewEditForm with credential
4. Set callbacks
5. Show via pageManager

#### `handleDeleteCredential()`
**Purpose**: Delete selected credential with confirmation

**Steps**:
1. Get selected credential
2. If none selected, show error
3. Show confirm dialog
4. If confirmed, delete via appState
5. Update status bar

#### `handleTogglePassword()`
**Purpose**: Toggle password visibility in detail view

**Context**: Only works when detail view is visible

**Steps**:
1. Get detail view from appState
2. Call TogglePasswordVisibility()
3. Refresh detail view

#### `handleCopyPassword()`
**Purpose**: Copy password to clipboard

**Steps**:
1. Get selected credential
2. If none selected, show error
3. Copy password via detail view
4. Show success message

#### `handleShowHelp()`
**Purpose**: Display keyboard shortcuts help screen

**Implementation**: Show modal with help text

### Focus Handling

#### `handleTabFocus()`
**Purpose**: Cycle focus to next component

**Delegates to**: `nav.CycleFocus()`

#### `handleShiftTabFocus()`
**Purpose**: Cycle focus to previous component

**Delegates to**: `nav.CycleFocusReverse()`

## Example Structure

```go
package events

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/components"
    "pass-cli/cmd/tui-tview/layout"
    "pass-cli/cmd/tui-tview/models"
)

type EventHandler struct {
    app         *tview.Application
    appState    *models.AppState
    nav         *models.NavigationState
    pageManager *layout.PageManager
    statusBar   *components.StatusBar
}

func NewEventHandler(app *tview.Application, appState *models.AppState, nav *models.NavigationState, pm *layout.PageManager, sb *components.StatusBar) *EventHandler {
    return &EventHandler{
        app:         app,
        appState:    appState,
        nav:         nav,
        pageManager: pm,
        statusBar:   sb,
    }
}

func (eh *EventHandler) SetupGlobalShortcuts() {
    eh.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        // ✅ CRITICAL: Protect input components
        focused := eh.app.GetFocus()
        if focused != nil {
            switch focused.(type) {
            case *tview.Form, *tview.InputField, *tview.TextArea:
                if event.Key() == tcell.KeyCtrlC {
                    eh.handleQuit()
                    return nil
                }
                return event  // Let input handle its keys
            }
        }

        // Handle global shortcuts
        return eh.handleGlobalKey(event)
    })
}

func (eh *EventHandler) handleGlobalKey(event *tcell.EventKey) *tcell.EventKey {
    switch event.Key() {
    case tcell.KeyRune:
        switch event.Rune() {
        case 'q':
            eh.handleQuit()
            return nil
        case 'n':
            eh.handleNewCredential()
            return nil
        case 'e':
            eh.handleEditCredential()
            return nil
        case 'd':
            eh.handleDeleteCredential()
            return nil
        case 'p':
            eh.handleTogglePassword()
            return nil
        case 'c':
            eh.handleCopyPassword()
            return nil
        case '?':
            eh.handleShowHelp()
            return nil
        }

    case tcell.KeyTab:
        eh.nav.CycleFocus()
        return nil

    case tcell.KeyBacktab:
        eh.nav.CycleFocusReverse()
        return nil

    case tcell.KeyCtrlC:
        eh.handleQuit()
        return nil
    }

    return event
}

func (eh *EventHandler) handleQuit() {
    // Check if modal is open
    if eh.pageManager.HasModals() {
        eh.pageManager.CloseTopModal()
        return
    }

    // Quit application
    eh.app.Stop()
}

func (eh *EventHandler) handleNewCredential() {
    form := components.NewAddForm(eh.appState)

    form.SetOnSubmit(func() {
        eh.pageManager.CloseModal("add-form")
        eh.statusBar.ShowSuccess("Credential added!")
    })

    form.SetOnCancel(func() {
        eh.pageManager.CloseModal("add-form")
    })

    eh.pageManager.ShowForm(form, "Add Credential")
}

func (eh *EventHandler) handleEditCredential() {
    cred := eh.appState.GetSelectedCredential()
    if cred == nil {
        eh.statusBar.ShowError(fmt.Errorf("no credential selected"))
        return
    }

    form := components.NewEditForm(eh.appState, cred)

    form.SetOnSubmit(func() {
        eh.pageManager.CloseModal("edit-form")
        eh.statusBar.ShowSuccess("Credential updated!")
    })

    form.SetOnCancel(func() {
        eh.pageManager.CloseModal("edit-form")
    })

    eh.pageManager.ShowForm(form, "Edit Credential")
}

func (eh *EventHandler) handleDeleteCredential() {
    cred := eh.appState.GetSelectedCredential()
    if cred == nil {
        eh.statusBar.ShowError(fmt.Errorf("no credential selected"))
        return
    }

    eh.pageManager.ShowConfirmDialog(
        "Delete Credential",
        fmt.Sprintf("Delete %s?", cred.Service),
        func() {
            err := eh.appState.DeleteCredential(cred.ID)
            if err != nil {
                eh.statusBar.ShowError(err)
            } else {
                eh.statusBar.ShowSuccess("Credential deleted")
            }
        },
        func() {
            // Cancelled
        },
    )
}

func (eh *EventHandler) handleTogglePassword() {
    detailView := eh.appState.GetDetailView()
    if detailView != nil {
        detailView.TogglePasswordVisibility()
    }
}

func (eh *EventHandler) handleCopyPassword() {
    detailView := eh.appState.GetDetailView()
    if detailView == nil {
        return
    }

    err := detailView.CopyPasswordToClipboard()
    if err != nil {
        eh.statusBar.ShowError(err)
    } else {
        eh.statusBar.ShowSuccess("Password copied to clipboard!")
    }
}

func (eh *EventHandler) handleShowHelp() {
    helpText := `
[yellow]Keyboard Shortcuts[-]

[cyan]Navigation[-]
  Tab          - Next component
  Shift+Tab    - Previous component
  ↑/↓         - Navigate lists
  Enter        - Select / View details

[cyan]Actions[-]
  n            - New credential
  e            - Edit credential
  d            - Delete credential
  p            - Toggle password visibility
  c            - Copy password

[cyan]General[-]
  ?            - Show this help
  q            - Quit
  Esc          - Close modal / Go back
  Ctrl+C       - Quit
`

    modal := tview.NewModal().
        SetText(helpText).
        AddButtons([]string{"Close"}).
        SetDoneFunc(func(buttonIndex int, buttonLabel string) {
            eh.pageManager.CloseModal("help")
        })

    eh.pageManager.ShowModal("help", modal, 60, 25)
}
```

## Keyboard Shortcut Summary

### Global (Work Everywhere)
- **q**: Quit application (or close modal if open)
- **Ctrl+C**: Force quit
- **Tab**: Next component
- **Shift+Tab**: Previous component
- **?**: Show help

### Context-Aware (Based on Focus)
- **n**: New credential (works everywhere)
- **e**: Edit selected credential (table/detail)
- **d**: Delete selected credential (table/detail)
- **p**: Toggle password (detail view)
- **c**: Copy password (detail view)
- **Esc**: Close modal / Go back

### Component-Specific (Handled by Component)
- **Arrow keys**: Navigate in lists/tables/trees
- **Enter**: Select item / Submit form
- **Page Up/Down**: Scroll in viewports

## Input Protection Pattern

### ✅ Correct (Current Implementation)
```go
focused := app.GetFocus()
if focused != nil {
    switch focused.(type) {
    case *tview.Form, *tview.InputField, *tview.TextArea:
        // Let component handle its own input
        return event
    }
}

// Only handle shortcuts if NOT in input component
handleShortcuts(event)
```

### ❌ Wrong (Previous Implementation)
```go
// This intercepted 'e' even when typing in a form!
switch event.Rune() {
case 'e':
    handleEdit()
    return nil
}
```

## Context-Aware Behavior

### When Table Has Focus
- **e**: Edit selected credential
- **d**: Delete selected credential
- **Enter**: View details

### When Detail View Has Focus
- **e**: Edit credential
- **d**: Delete credential
- **p**: Toggle password visibility
- **c**: Copy password

### When Modal Is Open
- **q**: Close modal, not quit app
- **Esc**: Close modal
- Other shortcuts: Disabled (modal has priority)

## Testing Considerations

- **Test input protection**: Type 'e', 'n', 'd' in forms
- **Test shortcut execution**: Verify actions fire correctly
- **Test context awareness**: Verify shortcuts work in right contexts
- **Test modal priority**: Verify 'q' closes modal before quitting
- **Mock components**: Use test harness with mock primitives

## Future Enhancements

- **Customizable shortcuts**: User-defined key bindings
- **Shortcut conflicts**: Detect and warn about conflicts
- **Shortcut hints**: Tooltip showing available shortcuts
- **Recording mode**: Record keyboard macros
- **Vim-style commands**: `:` command mode for power users
- **Mouse support**: Click handlers for actions
