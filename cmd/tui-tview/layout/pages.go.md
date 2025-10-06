# layout/pages.go

## Purpose
Modal and page management using tview.Pages. Handles showing/hiding forms and dialogs over the main UI.

## Responsibilities

1. **Page stacking**: Layer modals over main UI
2. **Modal display**: Show forms, dialogs, confirmations
3. **Page removal**: Close modals and return to main UI
4. **Focus management**: Restore focus when closing modals
5. **Escape handling**: Close topmost modal on Escape key

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/components` - For forms
- `pass-cli/cmd/tui-tview/models` - For NavigationState

### External Dependencies
- `github.com/rivo/tview` - Pages primitive
- `github.com/gdamore/tcell/v2` - For event handling

## Key Types

### `PageManager`
**Purpose**: Wrapper around tview.Pages with modal management

**Fields**:
```go
type PageManager struct {
    *tview.Pages  // Embedded Pages

    nav       *models.NavigationState
    pageStack []string  // Track page names for history
}
```

## Key Functions

### Constructor

#### `NewPageManager(mainLayout tview.Primitive, nav *models.NavigationState) *PageManager`
**Purpose**: Create page manager with main UI as base

**Steps**:
1. Create tview.Pages
2. Add main layout as "main" page
3. Setup Escape key handler
4. Return wrapper

### Page Management

#### `ShowModal(name string, modal tview.Primitive, width, height int) *PageManager`
**Purpose**: Display modal over current page

**Steps**:
1. Add modal to pages with name
2. Center modal at specified size
3. Push name to page stack
4. Switch to modal page
5. Set focus to modal

**Parameters**:
- `name`: Unique identifier for this modal
- `modal`: The primitive to display (form, dialog, etc.)
- `width, height`: Modal dimensions (-1 for auto-size)

**Returns**: Self for chaining

#### `CloseModal(name string)`
**Purpose**: Remove modal and return to previous page

**Steps**:
1. Remove page by name
2. Pop from page stack
3. Restore focus to previous page
4. If empty stack, focus on main

**Error handling**: Safe if page doesn't exist

#### `CloseTopModal()`
**Purpose**: Close the topmost modal

**Convenience**: Uses page stack to determine which modal to close

#### `ShowForm(form *components.CredentialForm, title string) *PageManager`
**Purpose**: Display credential form as modal

**Wrapper**: Calls ShowModal with form-appropriate sizing

#### `ShowConfirmDialog(title, message string, onConfirm, onCancel func())`
**Purpose**: Display yes/no confirmation dialog

**Example**:
```go
pm.ShowConfirmDialog(
    "Delete Credential",
    "Are you sure you want to delete this credential?",
    func() {
        appState.DeleteCredential(id)
        pm.CloseTopModal()
    },
    func() {
        pm.CloseTopModal()
    },
)
```

### Escape Key Handling

#### `setupEscapeHandler()`
**Purpose**: Close modal on Escape key

**Behavior**:
- If modal is open: Close topmost modal
- If no modal: Pass event to underlying page

**Implementation**:
```go
pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    if event.Key() == tcell.KeyEscape {
        if len(pm.pageStack) > 0 {
            pm.CloseTopModal()
            return nil  // Consume event
        }
    }
    return event  // Pass through
})
```

## Example Structure

```go
package layout

import (
    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/components"
    "pass-cli/cmd/tui-tview/models"
)

type PageManager struct {
    *tview.Pages
    nav       *models.NavigationState
    pageStack []string
}

func NewPageManager(mainLayout tview.Primitive, nav *models.NavigationState) *PageManager {
    pages := tview.NewPages()
    pages.AddPage("main", mainLayout, true, true)

    pm := &PageManager{
        Pages:     pages,
        nav:       nav,
        pageStack: []string{},
    }

    pm.setupEscapeHandler()

    return pm
}

func (pm *PageManager) ShowModal(name string, modal tview.Primitive, width, height int) *PageManager {
    pm.AddPage(name, modal, true, true)
    pm.pageStack = append(pm.pageStack, name)
    return pm
}

func (pm *PageManager) CloseModal(name string) {
    pm.RemovePage(name)

    // Remove from stack
    for i, page := range pm.pageStack {
        if page == name {
            pm.pageStack = append(pm.pageStack[:i], pm.pageStack[i+1:]...)
            break
        }
    }

    // Restore focus if stack is empty
    if len(pm.pageStack) == 0 {
        // Back to main page
        pm.nav.SetFocus(models.FocusSidebar)
    }
}

func (pm *PageManager) CloseTopModal() {
    if len(pm.pageStack) > 0 {
        topModal := pm.pageStack[len(pm.pageStack)-1]
        pm.CloseModal(topModal)
    }
}

func (pm *PageManager) ShowForm(form *components.CredentialForm, title string) *PageManager {
    // Wrap form in a centered box
    wrapper := tview.NewFlex().
        AddItem(nil, 0, 1, false).
        AddItem(tview.NewFlex().
            SetDirection(tview.FlexRow).
            AddItem(nil, 0, 1, false).
            AddItem(form, 20, 0, true).
            AddItem(nil, 0, 1, false),
            60, 0, true).
        AddItem(nil, 0, 1, false)

    return pm.ShowModal("form", wrapper, -1, -1)
}

func (pm *PageManager) ShowConfirmDialog(title, message string, onConfirm, onCancel func()) {
    modal := tview.NewModal().
        SetText(message).
        AddButtons([]string{"Confirm", "Cancel"}).
        SetDoneFunc(func(buttonIndex int, buttonLabel string) {
            pm.CloseTopModal()
            if buttonIndex == 0 {
                onConfirm()
            } else {
                onCancel()
            }
        })

    pm.ShowModal("confirm", modal, 60, 10)
}

func (pm *PageManager) setupEscapeHandler() {
    pm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyEscape {
            if len(pm.pageStack) > 0 {
                pm.CloseTopModal()
                return nil
            }
        }
        return event
    })
}

func (pm *PageManager) GetPageStack() []string {
    return append([]string{}, pm.pageStack...)  // Return copy
}

func (pm *PageManager) HasModals() bool {
    return len(pm.pageStack) > 0
}
```

## Usage Examples

### Show Add Form
```go
// In event handler for 'n' key:
form := components.NewAddForm(appState)

form.SetOnSubmit(func() {
    pageManager.CloseModal("add-form")
    statusBar.ShowSuccess("Credential added!")
})

form.SetOnCancel(func() {
    pageManager.CloseModal("add-form")
})

pageManager.ShowForm(form, "Add Credential")
```

### Show Delete Confirmation
```go
// In event handler for 'd' key:
cred := appState.GetSelectedCredential()
if cred == nil {
    return
}

pageManager.ShowConfirmDialog(
    "Delete Credential",
    fmt.Sprintf("Delete %s?", cred.Service),
    func() {
        err := appState.DeleteCredential(cred.ID)
        if err != nil {
            statusBar.ShowError(err)
        } else {
            statusBar.ShowSuccess("Credential deleted")
        }
    },
    func() {
        // User cancelled, do nothing
    },
)
```

### Show Edit Form
```go
// In event handler for 'e' key:
cred := appState.GetSelectedCredential()
if cred == nil {
    return
}

form := components.NewEditForm(appState, cred)

form.SetOnSubmit(func() {
    pageManager.CloseModal("edit-form")
    statusBar.ShowSuccess("Credential updated!")
})

form.SetOnCancel(func() {
    pageManager.CloseModal("edit-form")
})

pageManager.ShowForm(form, "Edit Credential")
```

## Interaction Flow

```
User presses 'n' to add credential
    ↓
Event handler creates AddForm
    ↓
pageManager.ShowForm(form, "Add Credential")
    ↓
Form added to pages as "form"
    ↓
Form name pushed to pageStack
    ↓
Focus set to form
    ↓
User fills form and presses Save
    ↓
form.onSubmit() callback fires
    ↓
Callback calls pageManager.CloseModal("form")
    ↓
Form removed from pages
    ↓
"form" popped from pageStack
    ↓
Focus returns to main page
    ↓
Table refreshes with new credential
```

## Page Stack Management

### Why Track Stack?
- Enables Escape key to close topmost modal
- Allows nested modals (e.g., confirm dialog over form)
- Proper focus restoration on close

### Stack Example
```
Initial: []
Add form: ["add-form"]
Show confirm: ["add-form", "confirm"]
Close confirm: ["add-form"]
Close form: []
```

## Modal Centering

### Using Flex for Centering
```go
// Horizontal centering
wrapper := tview.NewFlex().
    AddItem(nil, 0, 1, false).          // Left spacer (flex)
    AddItem(modal, width, 0, true).     // Modal (fixed width)
    AddItem(nil, 0, 1, false)           // Right spacer (flex)

// Vertical + Horizontal centering
fullWrapper := tview.NewFlex().
    SetDirection(tview.FlexRow).
    AddItem(nil, 0, 1, false).                    // Top spacer
    AddItem(horizontalWrapper, height, 0, true).  // Middle row
    AddItem(nil, 0, 1, false)                     // Bottom spacer
```

## Focus Restoration

When closing modal, restore focus appropriately:
- **After add**: Focus on new credential in table
- **After edit**: Focus on updated credential
- **After delete**: Focus on next credential
- **Cancel**: Restore previous focus

## Edge Cases

1. **Multiple modals**: Stack supports nested modals
2. **Escape on main page**: Pass through, don't consume
3. **Close non-existent modal**: Safe, no error
4. **Close while stack empty**: Safe, no action
5. **Form validation fails**: Don't close modal, keep open

## Integration with Navigation

### Focus Management
```go
// When showing modal:
previousFocus := nav.GetCurrentFocus()
nav.SetFocus(models.FocusForm)

// When closing modal:
nav.SetFocus(previousFocus)
```

### Breadcrumb Updates
```go
// When showing detail as modal (if needed):
nav.PushView("list")
pageManager.ShowModal("detail", detailView, -1, -1)

// When closing:
pageManager.CloseModal("detail")
nav.PopView()
```

## Testing Considerations

- **Test modal stacking**: Verify multiple modals work
- **Test Escape key**: Verify closes topmost
- **Test focus restoration**: Verify focus returns correctly
- **Test page removal**: Verify clean removal
- **Mock Pages**: Use test harness with mock pages

## Future Enhancements

- **Modal animations**: Fade in/out (challenging with tview)
- **Modal backdrop**: Semi-transparent overlay
- **Drag modals**: Allow repositioning (very advanced)
- **Modal sizing**: Smart auto-sizing based on content
- **Modal history**: Remember last shown modal for quick re-open
- **Keyboard shortcuts**: Number keys to select button (1=Confirm, 2=Cancel)
