package components

import (
	"fmt"

	"github.com/rivo/tview"
	"pass-cli/cmd/tui-tview/models"
	"pass-cli/cmd/tui-tview/styles"
	"pass-cli/internal/vault"
)

// AddForm provides a modal form for adding new credentials.
// Embeds tview.Form and manages validation and submission.
type AddForm struct {
	*tview.Form

	appState *models.AppState

	onSubmit func()
	onCancel func()
}

// EditForm provides a modal form for editing existing credentials.
// Embeds tview.Form and pre-populates fields from credential.
type EditForm struct {
	*tview.Form

	appState   *models.AppState
	credential *vault.CredentialMetadata

	onSubmit func()
	onCancel func()
}

// NewAddForm creates a new form for adding credentials.
// Creates input fields for Service, Username, Password, Category, URL, Notes.
func NewAddForm(appState *models.AppState) *AddForm {
	form := tview.NewForm()

	af := &AddForm{
		Form:     form,
		appState: appState,
	}

	af.buildFormFields()
	af.applyStyles()

	return af
}

// buildFormFields constructs all input fields for the add form.
func (af *AddForm) buildFormFields() {
	categories := af.getCategories()

	// Core credential fields
	af.AddInputField("Service", "", 40, nil, nil)
	af.AddInputField("Username", "", 40, nil, nil)
	af.AddPasswordField("Password", "", 40, '*', nil)

	// Optional metadata fields
	af.AddDropDown("Category", categories, 0, nil)
	af.AddInputField("URL", "", 50, nil, nil)
	af.AddTextArea("Notes", "", 50, 5, 0, nil)

	// Action buttons
	af.AddButton("Add", af.onAddPressed)
	af.AddButton("Cancel", af.onCancelPressed)
}

// onAddPressed handles the Add button submission.
// Validates inputs, calls AppState.AddCredential(), invokes onSubmit callback.
func (af *AddForm) onAddPressed() {
	// Validate inputs before submission
	if err := af.validate(); err != nil {
		// Validation failed - error will be shown via status bar or modal
		// Form stays open for correction
		return
	}

	// Extract field values (using form item index)
	service := af.GetFormItem(0).(*tview.InputField).GetText()
	username := af.GetFormItem(1).(*tview.InputField).GetText()
	password := af.GetFormItem(2).(*tview.InputField).GetText()

	// Extract category from dropdown (index 3)
	categoryDropdown := af.GetFormItem(3).(*tview.DropDown)
	_, category := categoryDropdown.GetCurrentOption()

	url := af.GetFormItem(4).(*tview.InputField).GetText()
	notes := af.GetFormItem(5).(*tview.TextArea).GetText()

	// Call AppState to add credential with all 6 fields
	err := af.appState.AddCredential(service, username, password, category, url, notes)
	if err != nil {
		// Error already handled by AppState onError callback
		// Form stays open for correction
		return
	}

	// Success - invoke callback to close modal
	if af.onSubmit != nil {
		af.onSubmit()
	}
}

// onCancelPressed handles the Cancel button.
// Invokes onCancel callback to close modal without saving.
func (af *AddForm) onCancelPressed() {
	if af.onCancel != nil {
		af.onCancel()
	}
}

// validate checks that required fields are filled.
// Returns error describing first validation failure, or nil if valid.
func (af *AddForm) validate() error {
	// Service is required (cannot be empty)
	service := af.GetFormItem(0).(*tview.InputField).GetText()
	if service == "" {
		return fmt.Errorf("service is required")
	}

	// Username is required (minimum validation)
	username := af.GetFormItem(1).(*tview.InputField).GetText()
	if username == "" {
		return fmt.Errorf("username is required")
	}

	// Password validation (basic check)
	password := af.GetFormItem(2).(*tview.InputField).GetText()
	if password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

// getCategories retrieves available categories from AppState.
// Returns default "Uncategorized" if no categories exist.
func (af *AddForm) getCategories() []string {
	categories := af.appState.GetCategories()
	if len(categories) == 0 {
		return []string{"Uncategorized"}
	}
	return categories
}

// applyStyles applies theme colors and border styling to the form.
func (af *AddForm) applyStyles() {
	theme := styles.GetCurrentTheme()

	// Apply form-level styling
	styles.ApplyFormStyle(af.Form)

	// Set border and title
	af.SetBorder(true).
		SetTitle(" Add Credential ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(theme.BorderColor)

	// Button alignment
	af.SetButtonsAlign(tview.AlignRight)
}

// SetOnSubmit registers a callback to be invoked after successful add.
func (af *AddForm) SetOnSubmit(callback func()) {
	af.onSubmit = callback
}

// SetOnCancel registers a callback to be invoked when cancel is pressed.
func (af *AddForm) SetOnCancel(callback func()) {
	af.onCancel = callback
}

// NewEditForm creates a new form for editing an existing credential.
// Pre-populates all fields with current credential values.
func NewEditForm(appState *models.AppState, credential *vault.CredentialMetadata) *EditForm {
	form := tview.NewForm()

	ef := &EditForm{
		Form:       form,
		appState:   appState,
		credential: credential,
	}

	ef.buildFormFieldsWithValues()
	ef.applyStyles()

	return ef
}

// buildFormFieldsWithValues constructs form fields pre-populated with credential data.
func (ef *EditForm) buildFormFieldsWithValues() {
	categories := ef.getCategories()
	categoryIndex := ef.findCategoryIndex(categories)

	// Pre-populate fields with existing credential data
	ef.AddInputField("Service", ef.credential.Service, 40, nil, nil)
	ef.AddInputField("Username", ef.credential.Username, 40, nil, nil)

	// Password field - need to fetch full credential to get password
	// For now, leave empty (user can enter new password or keep existing)
	ef.AddPasswordField("Password", "", 40, '*', nil)

	// Optional metadata fields
	// TODO: Pre-populate from credential once Category, URL, Notes are available
	ef.AddDropDown("Category", categories, categoryIndex, nil)
	ef.AddInputField("URL", "", 50, nil, nil)
	ef.AddTextArea("Notes", "", 50, 5, 0, nil)

	// Action buttons
	ef.AddButton("Save", ef.onSavePressed)
	ef.AddButton("Cancel", ef.onCancelPressed)
}

// onSavePressed handles the Save button submission.
// Validates inputs, calls AppState.UpdateCredential(), invokes onSubmit callback.
func (ef *EditForm) onSavePressed() {
	// Validate inputs before submission
	if err := ef.validate(); err != nil {
		// Validation failed - form stays open for correction
		return
	}

	// Extract field values
	service := ef.GetFormItem(0).(*tview.InputField).GetText()
	username := ef.GetFormItem(1).(*tview.InputField).GetText()
	password := ef.GetFormItem(2).(*tview.InputField).GetText()

	// Extract category from dropdown
	categoryDropdown := ef.GetFormItem(3).(*tview.DropDown)
	_, category := categoryDropdown.GetCurrentOption()

	url := ef.GetFormItem(4).(*tview.InputField).GetText()
	notes := ef.GetFormItem(5).(*tview.TextArea).GetText()

	// Build UpdateCredentialOpts with only non-empty fields
	opts := models.UpdateCredentialOpts{}

	if username != "" {
		opts.Username = &username
	}

	// If password is empty, don't update it (nil pointer = don't change)
	// If password is non-empty, update it
	if password != "" {
		opts.Password = &password
	}

	// Always set category (even if empty, to allow clearing)
	opts.Category = &category

	// Always set URL (even if empty, to allow clearing)
	opts.URL = &url

	// Always set notes (even if empty, to allow clearing)
	opts.Notes = &notes

	// Call AppState to update credential with options struct
	err := ef.appState.UpdateCredential(service, opts)
	if err != nil {
		// Error already handled by AppState onError callback
		// Form stays open for correction
		return
	}

	// Success - invoke callback to close modal
	if ef.onSubmit != nil {
		ef.onSubmit()
	}
}

// onCancelPressed handles the Cancel button.
// Invokes onCancel callback to close modal without saving.
func (ef *EditForm) onCancelPressed() {
	if ef.onCancel != nil {
		ef.onCancel()
	}
}

// validate checks that required fields are filled.
// Returns error describing first validation failure, or nil if valid.
func (ef *EditForm) validate() error {
	// Service is required (cannot be empty)
	service := ef.GetFormItem(0).(*tview.InputField).GetText()
	if service == "" {
		return fmt.Errorf("service is required")
	}

	// Username is required (minimum validation)
	username := ef.GetFormItem(1).(*tview.InputField).GetText()
	if username == "" {
		return fmt.Errorf("username is required")
	}

	// Password not required in edit form (can keep existing)

	return nil
}

// getCategories retrieves available categories from AppState.
// Returns default "Uncategorized" if no categories exist.
func (ef *EditForm) getCategories() []string {
	categories := ef.appState.GetCategories()
	if len(categories) == 0 {
		return []string{"Uncategorized"}
	}
	return categories
}

// findCategoryIndex finds the index of the credential's category in the dropdown.
// Returns 0 if category not found.
func (ef *EditForm) findCategoryIndex(categories []string) int {
	// TODO: Once credential has Category field, match it here
	// For now, return 0 (first category)
	return 0
}

// applyStyles applies theme colors and border styling to the form.
func (ef *EditForm) applyStyles() {
	theme := styles.GetCurrentTheme()

	// Apply form-level styling
	styles.ApplyFormStyle(ef.Form)

	// Set border and title
	ef.SetBorder(true).
		SetTitle(" Edit Credential ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(theme.BorderColor)

	// Button alignment
	ef.SetButtonsAlign(tview.AlignRight)
}

// SetOnSubmit registers a callback to be invoked after successful update.
func (ef *EditForm) SetOnSubmit(callback func()) {
	ef.onSubmit = callback
}

// SetOnCancel registers a callback to be invoked when cancel is pressed.
func (ef *EditForm) SetOnCancel(callback func()) {
	ef.onCancel = callback
}
