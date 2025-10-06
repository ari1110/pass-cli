# components/forms.go

## Purpose
Modal forms for adding and editing credentials using tview.Form and tview.Modal. Provides input validation and submission handling.

## Responsibilities

1. **Add credential form**: Collect service, username, password, category, notes
2. **Edit credential form**: Pre-populate and update existing credential
3. **Input validation**: Ensure required fields are filled
4. **Password generation**: Optional password generator integration
5. **Modal display**: Show form as modal dialog over main UI

## Dependencies

### Internal Dependencies
- `pass-cli/cmd/tui-tview/models` - For AppState
- `pass-cli/internal/vault` - For Credential type

### External Dependencies
- `github.com/rivo/tview` - Form and Modal primitives
- `github.com/gdamore/tcell/v2` - For colors

## Key Types

### `CredentialForm`
**Purpose**: Wrapper around tview.Form for credential input

**Fields**:
```go
type CredentialForm struct {
    *tview.Form  // Embedded Form

    appState       *models.AppState
    mode           FormMode  // Add or Edit
    credentialID   string    // For edit mode
    onSubmit       func()    // Callback on successful submit
    onCancel       func()    // Callback on cancel
}

type FormMode int
const (
    FormModeAdd FormMode = iota
    FormModeEdit
)
```

## Key Functions

### Constructors

#### `NewAddForm(appState *models.AppState) *CredentialForm`
**Purpose**: Create form for adding new credential

**Fields**:
- Service (required)
- Username (required)
- Password (required) + "Generate" button
- Category (optional dropdown)
- URL (optional)
- Notes (optional textarea)

**Buttons**:
- Save
- Cancel

#### `NewEditForm(appState *models.AppState, credential *vault.Credential) *CredentialForm`
**Purpose**: Create form for editing existing credential

**Pre-populated fields**: All fields filled with current values

**Buttons**:
- Update
- Cancel

### Form Building

#### `buildForm()`
**Purpose**: Construct form fields based on mode

**Field types**:
```go
form.AddInputField("Service", defaultValue, 30, nil, nil)
form.AddInputField("Username", defaultValue, 30, nil, nil)
form.AddPasswordField("Password", defaultValue, 30, '*', nil)
form.AddDropDown("Category", categories, selectedIndex, nil)
form.AddInputField("URL", defaultValue, 50, nil, nil)
form.AddTextArea("Notes", defaultValue, 50, 5, 0, nil)
```

#### `populateCategories() []string`
**Purpose**: Get list of existing categories for dropdown

**Returns**: Sorted unique categories from vault + "New Category..." option

### Validation

#### `validateForm() error`
**Purpose**: Validate all required fields before submission

**Rules**:
- Service: Required, min 1 character
- Username: Required, min 1 character
- Password: Required, min 8 characters (or configurable)
- Category: Optional
- URL: Optional, must be valid URL if provided
- Notes: Optional

**Returns**: Error describing first validation failure, or nil if valid

### Submission

#### `onSave()`
**Purpose**: Handle form submission

**Steps**:
1. Validate form fields
2. If validation fails, show error and return
3. Extract field values
4. Call appropriate AppState method:
   - Add mode: `appState.AddCredential()`
   - Edit mode: `appState.UpdateCredential()`
5. If error, show error message
6. If success, call onSubmit callback (closes modal)

#### `onCancelPressed()`
**Purpose**: Handle cancel button

**Steps**:
1. Discard changes
2. Call onCancel callback (closes modal without saving)

### Password Generation

#### `generatePassword()`
**Purpose**: Generate strong random password

**Behavior**:
1. Open password generator dialog
2. User configures: length, character types
3. Generate password
4. Fill password field with generated value

**Integration**: Button next to password field

### Modal Wrapping

#### `WrapInModal(form *CredentialForm) *tview.Modal`
**Purpose**: Wrap form in modal dialog

**Why**: Forms should appear as popups over main UI, not replace it

**Configuration**:
- Semi-transparent background
- Centered on screen
- Fixed size or auto-size based on content
- Escape key to cancel

## Example Structure

```go
package components

import (
    "fmt"
    "net/url"

    "github.com/rivo/tview"
    "pass-cli/cmd/tui-tview/models"
    "pass-cli/cmd/tui-tview/styles"
    "pass-cli/internal/vault"
)

type FormMode int

const (
    FormModeAdd FormMode = iota
    FormModeEdit
)

type CredentialForm struct {
    *tview.Form
    appState     *models.AppState
    mode         FormMode
    credentialID string
    onSubmit     func()
    onCancel     func()
}

func NewAddForm(appState *models.AppState) *CredentialForm {
    form := tview.NewForm()

    cf := &CredentialForm{
        Form:     form,
        appState: appState,
        mode:     FormModeAdd,
    }

    cf.buildForm()
    cf.applyStyles()

    return cf
}

func NewEditForm(appState *models.AppState, credential *vault.Credential) *CredentialForm {
    form := tview.NewForm()

    cf := &CredentialForm{
        Form:         form,
        appState:     appState,
        mode:         FormModeEdit,
        credentialID: credential.ID,
    }

    cf.buildFormWithValues(credential)
    cf.applyStyles()

    return cf
}

func (cf *CredentialForm) buildForm() {
    categories := cf.populateCategories()

    cf.AddInputField("Service", "", 40, nil, nil)
    cf.AddInputField("Username", "", 40, nil, nil)
    cf.AddPasswordField("Password", "", 40, '*', nil)
    cf.AddButton("Generate", cf.generatePassword)
    cf.AddDropDown("Category", categories, 0, nil)
    cf.AddInputField("URL", "", 50, nil, nil)
    cf.AddTextArea("Notes", "", 50, 5, 0, nil)

    cf.AddButton("Save", cf.onSave)
    cf.AddButton("Cancel", cf.onCancelPressed)
}

func (cf *CredentialForm) buildFormWithValues(cred *vault.Credential) {
    categories := cf.populateCategories()
    categoryIndex := findCategoryIndex(categories, cred.Category)

    cf.AddInputField("Service", cred.Service, 40, nil, nil)
    cf.AddInputField("Username", cred.Username, 40, nil, nil)
    cf.AddPasswordField("Password", cred.Password, 40, '*', nil)
    cf.AddDropDown("Category", categories, categoryIndex, nil)
    cf.AddInputField("URL", cred.URL, 50, nil, nil)
    cf.AddTextArea("Notes", cred.Notes, 50, 5, 0, nil)

    cf.AddButton("Update", cf.onSave)
    cf.AddButton("Cancel", cf.onCancelPressed)
}

func (cf *CredentialForm) onSave() {
    // Validate
    if err := cf.validateForm(); err != nil {
        // Show error (need status bar reference or modal error)
        return
    }

    // Extract values
    service := cf.GetFormItem(0).(*tview.InputField).GetText()
    username := cf.GetFormItem(1).(*tview.InputField).GetText()
    password := cf.GetFormItem(2).(*tview.InputField).GetText()
    _, category := cf.GetFormItem(4).(*tview.DropDown).GetCurrentOption()
    url := cf.GetFormItem(5).(*tview.InputField).GetText()
    notes := cf.GetFormItem(6).(*tview.TextArea).GetText()

    // Save based on mode
    var err error
    if cf.mode == FormModeAdd {
        err = cf.appState.AddCredential(service, username, password, category, url, notes)
    } else {
        err = cf.appState.UpdateCredential(cf.credentialID, service, username, password, category, url, notes)
    }

    if err != nil {
        // Show error
        return
    }

    // Success - call callback to close modal
    if cf.onSubmit != nil {
        cf.onSubmit()
    }
}

func (cf *CredentialForm) onCancelPressed() {
    if cf.onCancel != nil {
        cf.onCancel()
    }
}

func (cf *CredentialForm) validateForm() error {
    service := cf.GetFormItem(0).(*tview.InputField).GetText()
    if service == "" {
        return fmt.Errorf("service is required")
    }

    username := cf.GetFormItem(1).(*tview.InputField).GetText()
    if username == "" {
        return fmt.Errorf("username is required")
    }

    password := cf.GetFormItem(2).(*tview.InputField).GetText()
    if len(password) < 8 {
        return fmt.Errorf("password must be at least 8 characters")
    }

    urlField := cf.GetFormItem(5).(*tview.InputField).GetText()
    if urlField != "" {
        if _, err := url.Parse(urlField); err != nil {
            return fmt.Errorf("invalid URL format")
        }
    }

    return nil
}

func (cf *CredentialForm) populateCategories() []string {
    categories := cf.appState.GetCategories()
    // Add "New Category..." option
    return append(categories, "New Category...")
}

func (cf *CredentialForm) generatePassword() {
    // TODO: Implement password generator
    // For now, generate a simple random password
    password := generateRandomPassword(16)
    cf.GetFormItem(2).(*tview.InputField).SetText(password)
}

func (cf *CredentialForm) applyStyles() {
    cf.SetBorder(true).
        SetTitle(" Add Credential ").
        SetTitleAlign(tview.AlignLeft).
        SetBorderColor(styles.BorderColor).
        SetBackgroundColor(styles.BackgroundColor)

    cf.SetButtonsAlign(tview.AlignRight)
    cf.SetButtonBackgroundColor(styles.ButtonBackground)
    cf.SetButtonTextColor(styles.ButtonText)
}

func (cf *CredentialForm) SetOnSubmit(callback func()) {
    cf.onSubmit = callback
}

func (cf *CredentialForm) SetOnCancel(callback func()) {
    cf.onCancel = callback
}
```

## Usage Example

```go
// In event handler for 'n' key (new credential):
func showAddCredentialForm(app *tview.Application, pages *tview.Pages, appState *models.AppState) {
    form := components.NewAddForm(appState)

    form.SetOnSubmit(func() {
        pages.RemovePage("add-form")  // Close modal
        statusBar.ShowSuccess("Credential added successfully!")
    })

    form.SetOnCancel(func() {
        pages.RemovePage("add-form")  // Close modal
    })

    // Create modal
    modal := tview.NewModal().
        SetText("").
        AddButtons([]string{}).
        SetBackgroundColor(tcell.ColorBlack)

    // Add form to modal (hack: use modal as container)
    pages.AddPage("add-form", form, true, true)
}
```

## Interaction Flow

```
User presses 'n' (new credential)
    ↓
Event handler creates NewAddForm
    ↓
Form displayed as modal over main UI
    ↓
User fills fields
    ↓
User clicks "Generate" for password
    ↓
Password field populated
    ↓
User clicks "Save"
    ↓
Form validates fields
    ↓
If valid: appState.AddCredential()
    ↓
onSubmit callback fires
    ↓
Modal closes
    ↓
Status bar shows success message
    ↓
Table refreshes with new credential
```

## Visual Design

### Form Layout
```
┌─ Add Credential ────────────────────┐
│                                      │
│ Service:    [___________________________] │
│ Username:   [___________________________] │
│ Password:   [***************************] [Generate] │
│ Category:   [AWS              ▼] │
│ URL:        [___________________________] │
│ Notes:      [                          ] │
│             [                          ] │
│             [                          ] │
│                                      │
│                      [Save] [Cancel] │
└──────────────────────────────────────┘
```

### Colors
- Border: Accent color (cyan)
- Labels: Gray
- Input fields: White background
- Buttons: Highlighted background
- Modal backdrop: Semi-transparent black

## Keyboard Navigation

### Within Form
- **Tab**: Next field
- **Shift+Tab**: Previous field
- **Enter**: Submit (when on button)
- **Esc**: Cancel

### Form-specific
- **Arrow keys**: Navigate dropdowns
- **Space**: Toggle checkboxes (if any)

## Integration with Pages

Forms are displayed using tview.Pages for modal behavior:

```go
pages := tview.NewPages()
pages.AddPage("main", mainLayout, true, true)

// Show form modal
pages.AddPage("form", formModal, true, true)

// Close form
pages.RemovePage("form")
```

## Validation Messages

### Field-level
- Show validation errors inline (red text below field)
- Highlight invalid field with red border

### Form-level
- Show error modal if submission fails
- Prevent submission until all fields valid

## Edge Cases

1. **Empty vault**: No categories in dropdown, show "Uncategorized" only
2. **Long field values**: Truncate or scroll within field
3. **Special characters in password**: Handle all UTF-8 characters
4. **Duplicate service names**: Allow duplicates (user's choice)
5. **Cancel with unsaved changes**: Optionally confirm discard

## Testing Considerations

- **Test validation**: Verify all validation rules
- **Test submission**: Verify correct AppState method called
- **Test cancellation**: Verify no changes made
- **Test pre-population**: Verify edit mode fills fields correctly
- **Mock AppState**: Use mock for isolated testing

## Future Enhancements

- **Password generator dialog**: Full-featured generator with options
- **Password strength meter**: Visual indicator of password strength
- **Field templates**: Pre-fill common credential types
- **Bulk import**: CSV/JSON import form
- **Custom fields**: User-defined additional fields
- **Autofill suggestions**: Suggest usernames/categories based on service
- **Validation rules**: Custom validation per field type
- **Async validation**: Check for duplicates while typing
