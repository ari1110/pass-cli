// Package components provides TUI form components for credential management.
// All forms support the complete credential model: service, username, password, category, URL, and notes.
package components

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"pass-cli/cmd/tui/models"
	"pass-cli/cmd/tui/styles"
	"pass-cli/internal/security"
	"pass-cli/internal/vault"
)

// readOnlyLabel wraps TextView as a non-focusable FormItem for display-only content.
// Used for password strength indicators and keyboard hints.
type readOnlyLabel struct {
	*tview.TextView
}

func newReadOnlyLabel() *readOnlyLabel {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true)
	tv.SetWrap(true)
	tv.SetWordWrap(true)
	tv.SetScrollable(false) // Don't allow scrolling in single-line labels

	return &readOnlyLabel{TextView: tv}
}

// FormItem interface implementation
func (r *readOnlyLabel) GetLabel() string                        { return "" }
func (r *readOnlyLabel) SetLabel(label string) tview.FormItem    { return r }
func (r *readOnlyLabel) GetLabelWidth() int                      { return 0 }
func (r *readOnlyLabel) SetLabelWidth(width int) tview.FormItem  { return r }
func (r *readOnlyLabel) GetFieldWidth() int                      { return 0 }
func (r *readOnlyLabel) SetFieldWidth(width int) tview.FormItem  { return r }
func (r *readOnlyLabel) GetFieldHeight() int                     { return 1 }
func (r *readOnlyLabel) SetDisabled(disabled bool) tview.FormItem { return r }
func (r *readOnlyLabel) GetDisabled() bool                       { return true }
func (r *readOnlyLabel) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	return r
}
func (r *readOnlyLabel) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	r.SetBackgroundColor(bgColor)
	r.SetTextColor(fieldTextColor)
	return r
}

// normalizeCategory converts the "Uncategorized" UI label to empty string for storage.
// Prevents the UI label from leaking into credential data.
func normalizeCategory(c string) string {
	if c == "Uncategorized" {
		return ""
	}
	return c
}

// AddForm provides a modal form for adding new credentials.
// Embeds tview.Frame (which contains the Form) and manages validation and submission.
type AddForm struct {
	*tview.Frame
	form *tview.Form

	appState *models.AppState

	passwordVisible bool // Track password visibility state for toggle

	onSubmit func()
	onCancel func()
}

// EditForm provides a modal form for editing existing credentials.
// Embeds tview.Frame (which contains the Form) and pre-populates fields from credential.
type EditForm struct {
	*tview.Frame
	form *tview.Form

	appState   *models.AppState
	credential *vault.CredentialMetadata

	originalPassword string // Track original password to detect changes
	passwordFetched  bool   // Track if password has been fetched (lazy loading)
	passwordVisible  bool   // Track password visibility state for toggle

	onSubmit func()
	onCancel func()
}

// NewAddForm creates a new form for adding credentials.
// Creates input fields for Service, Username, Password, Category, URL, Notes.
func NewAddForm(appState *models.AppState) *AddForm {
	form := tview.NewForm()

	af := &AddForm{
		form:     form,
		appState: appState,
	}

	af.buildFormFields()
	af.applyStyles()
	af.setupKeyboardShortcuts()
	af.wrapInFrame()

	return af
}

// buildFormFields constructs all input fields for the add form.
func (af *AddForm) buildFormFields() {
	categories := af.getCategories()

	// Ensure "Uncategorized" is always present and find its index
	uncategorizedIndex := -1
	for i, cat := range categories {
		if cat == "Uncategorized" {
			uncategorizedIndex = i
			break
		}
	}

	// If "Uncategorized" is not present, prepend it to the list
	if uncategorizedIndex == -1 {
		categories = append([]string{"Uncategorized"}, categories...)
		uncategorizedIndex = 0
	}

	// Core credential fields
	// Use 0 width to make fields fill available space (prevents black rectangles)
	af.form.AddInputField("Service (UID)", "", 0, nil, nil)
	af.form.AddInputField("Username", "", 0, nil, nil)

	// T048, T049: Password field with real-time strength indicator
	passwordField := tview.NewInputField().
		SetLabel("Password").
		SetFieldWidth(0).
		SetMaskCharacter('*')

	// T049: Update password label with strength indicator inline
	passwordField.SetChangedFunc(func(text string) {
		af.updatePasswordLabel(passwordField, []byte(text))
	})

	af.form.AddFormItem(passwordField)

	// Optional metadata fields - default to "Uncategorized"

	// Original dropdown approach (commented out for autocomplete field)
	// af.form.AddDropDown("Category", categories, uncategorizedIndex, nil)

	// New autocomplete input field for Category
	categoryField := tview.NewInputField().
      SetLabel("Category").
      SetFieldWidth(0)

  categoryField.SetAutocompleteFunc(func(currentText string) []string {
      if currentText == "" {
          return categories
      }
      var matches []string
      lowerText := strings.ToLower(currentText)
      for _, cat := range categories {
          if strings.HasPrefix(strings.ToLower(cat), lowerText) {
              matches = append(matches, cat)
          }
      }
      return matches
  })

  af.form.AddFormItem(categoryField)
	af.form.AddInputField("URL", "", 0, nil, nil)
	af.form.AddTextArea("Notes", "", 0, 5, 0, nil)

	// Action buttons
	af.form.AddButton("Add", af.onAddPressed)
	af.form.AddButton("Cancel", af.onCancelPressed)
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
	service := af.form.GetFormItem(0).(*tview.InputField).GetText()
	username := af.form.GetFormItem(1).(*tview.InputField).GetText()
	password := af.form.GetFormItem(2).(*tview.InputField).GetText()

	// Extract category from input field
	category := af.form.GetFormItem(3).(*tview.InputField).GetText()
	category = normalizeCategory(category) // Convert "Uncategorized" to empty string

	url := af.form.GetFormItem(4).(*tview.InputField).GetText()
	notes := af.form.GetFormItem(5).(*tview.TextArea).GetText()

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
	service := af.form.GetFormItem(0).(*tview.InputField).GetText()
	if service == "" {
		return fmt.Errorf("service is required")
	}

	// Username is required (minimum validation)
	username := af.form.GetFormItem(1).(*tview.InputField).GetText()
	if username == "" {
		return fmt.Errorf("username is required")
	}

	// Password validation (basic check)
	password := af.form.GetFormItem(2).(*tview.InputField).GetText()
	if password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

// T048, T049: updatePasswordLabel updates the password field label with strength indicator
func (af *AddForm) updatePasswordLabel(field *tview.InputField, password []byte) {
	policy := security.DefaultPasswordPolicy
	strength := policy.Strength(password)

	var label string
	if len(password) == 0 {
		label = "Password"
	} else {
		switch strength {
		case security.PasswordStrengthWeak:
			label = "Password [yellow](Weak)[-]"
		case security.PasswordStrengthMedium:
			label = "Password [orange](Medium)[-]"
		case security.PasswordStrengthStrong:
			label = "Password [green](Strong)[-]"
		}
	}

	field.SetLabel(label)
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

// wrapInFrame wraps the form in a Frame to add keyboard hints footer.
func (af *AddForm) wrapInFrame() {
	theme := styles.GetCurrentTheme()

	frame := tview.NewFrame(af.form)
	frame.SetBorder(true).
		SetTitle(" Add Credential ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(theme.BorderColor)

	// Add keyboard hints as footer text
	hintsText := "Tab: Next  •  Shift+Tab: Prev  •  Ctrl+S: Add  •  Ctrl+H: Toggle pwd  •  Esc: Cancel"
	frame.AddText(hintsText, false, tview.AlignCenter, theme.TextSecondary)

	af.Frame = frame
}


// setupKeyboardShortcuts configures form-level keyboard shortcuts.
// Adds Ctrl+S for quick-save and ensures Tab/Shift+Tab stay within form.
func (af *AddForm) setupKeyboardShortcuts() {
	af.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			// Ctrl+S for quick-save
			af.onAddPressed()
			return nil

		case tcell.KeyCtrlH:
			// Only toggle if Ctrl modifier is actually pressed
			// Backspace sends KeyCtrlH but without ModCtrl set
			if event.Modifiers()&tcell.ModCtrl != 0 {
				af.togglePasswordVisibility()
				return nil
			}
			// Regular backspace - pass through to allow deletion
			return event

		case tcell.KeyEscape:
			// Handle Escape to close form
			af.onCancelPressed()
			return nil

		case tcell.KeyTab, tcell.KeyBacktab:
			// Let form handle Tab/Shift+Tab internally
			// tview's Form will cycle through focusable items and skip disabled ones
			// Return event to allow built-in navigation
			return event
		}
		return event
	})
}

// applyStyles applies theme colors and border styling to the form.
func (af *AddForm) applyStyles() {
	theme := styles.GetCurrentTheme()

	// Apply form-level styling
	styles.ApplyFormStyle(af.form)

	// Style individual input fields
	// Use BackgroundLight for input fields - lighter than form Background for contrast
	for i := 0; i < 6; i++ {
		item := af.form.GetFormItem(i)
		switch field := item.(type) {
		case *tview.InputField:
			field.SetFieldBackgroundColor(theme.BackgroundLight).
				SetFieldTextColor(theme.TextPrimary)
		case *tview.TextArea:
			field.SetTextStyle(tcell.StyleDefault.
				Background(theme.BackgroundLight).
				Foreground(theme.TextPrimary))
		case *tview.DropDown:
			field.SetFieldBackgroundColor(theme.BackgroundLight).
				SetFieldTextColor(theme.TextPrimary)
		}
	}

	// Button alignment
	af.form.SetButtonsAlign(tview.AlignRight)
}

// togglePasswordVisibility switches between masked and plaintext password display.
// Updates both the mask character and the field label to indicate current state.
func (af *AddForm) togglePasswordVisibility() {
	af.passwordVisible = !af.passwordVisible

	passwordField := af.form.GetFormItem(2).(*tview.InputField)

	if af.passwordVisible {
		passwordField.SetMaskCharacter(0) // 0 = plaintext (tview convention)
		passwordField.SetLabel("Password [VISIBLE]")
	} else {
		passwordField.SetMaskCharacter('*')
		// Restore label with current strength if password exists
		text := passwordField.GetText()
		if text != "" {
			af.updatePasswordLabel(passwordField, []byte(text))
		} else {
			passwordField.SetLabel("Password")
		}
	}
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
		form:       form,
		appState:   appState,
		credential: credential,
	}

	ef.buildFormFieldsWithValues()
	ef.applyStyles()
	ef.setupKeyboardShortcuts()
	ef.wrapInFrame()

	return ef
}

// buildFormFieldsWithValues constructs form fields pre-populated with credential data.
func (ef *EditForm) buildFormFieldsWithValues() {
	categories := ef.getCategories()

	// Pre-populate fields with existing credential data
	// Service field is read-only (cannot be changed in edit mode)
	// Use 0 width to make fields fill available space (prevents black rectangles)
	ef.form.AddInputField("Service (UID)", ef.credential.Service, 0, nil, nil)
	serviceField := ef.form.GetFormItem(0).(*tview.InputField)
	serviceField.SetDisabled(true) // Make read-only to prevent confusion

	ef.form.AddInputField("Username", ef.credential.Username, 0, nil, nil)

	// T048, T049: Password field with real-time strength indicator
	// Defer fetching until user focuses field (lazy loading)
	// This prevents blocking UI and avoids incrementing usage stats on form open
	passwordField := tview.NewInputField().
		SetLabel("Password").
		SetFieldWidth(0).
		SetMaskCharacter('*')

	// Attach focus handler to fetch password lazily
	passwordField.SetFocusFunc(func() {
		ef.fetchPasswordIfNeeded(passwordField)
	})

	// T049: Update password label with strength indicator inline
	passwordField.SetChangedFunc(func(text string) {
		ef.updatePasswordLabel(passwordField, []byte(text))
	})

	ef.form.AddFormItem(passwordField)

	// Optional metadata fields - pre-populated from credential
	// ef.AddDropDown("Category", categories, categoryIndex, nil)

	// Replace DropDown with InputField + autocomplete
	// Pre-populate with existing category (or "Uncategorized" if empty)
	initialCategory := ef.credential.Category
	if initialCategory == "" {
		initialCategory = "Uncategorized"
	}

	categoryField := tview.NewInputField().
		SetLabel("Category").
		SetFieldWidth(0).
		SetText(initialCategory)

	categoryField.SetAutocompleteFunc(func(currentText string) []string {
		if currentText == "" {
			return categories
		}
		var matches []string
		lowerText := strings.ToLower(currentText)
		for _, cat := range categories {
			if strings.HasPrefix(strings.ToLower(cat), lowerText) {
				matches = append(matches, cat)
			}
		}
		return matches
	})

	ef.form.AddFormItem(categoryField)

	ef.form.AddInputField("URL", ef.credential.URL, 0, nil, nil)
	ef.form.AddTextArea("Notes", ef.credential.Notes, 0, 5, 0, nil)

	// Action buttons
	ef.form.AddButton("Save", ef.onSavePressed)
	ef.form.AddButton("Cancel", ef.onCancelPressed)
}

// fetchPasswordIfNeeded lazily fetches the password when the field is focused.
// Uses track=false to avoid incrementing usage statistics on form pre-population.
// Caches result to avoid redundant fetches on repeated focus events.
func (ef *EditForm) fetchPasswordIfNeeded(passwordField *tview.InputField) {
	// Only fetch once
	if ef.passwordFetched {
		return
	}

	// Fetch credential without tracking (track=false)
	cred, err := ef.appState.GetFullCredentialWithTracking(ef.credential.Service, false)
	if err != nil {
		// Surface error via AppState error handler without blocking UI
		// Leave password field empty on error
		ef.passwordFetched = true // Mark as attempted to avoid retry loops
		return
	}

	// Set password field and cache original value
	// T020d: Convert []byte password to string for display
	if cred != nil {
		ef.originalPassword = string(cred.Password)
		passwordField.SetText(string(cred.Password))
	}

	ef.passwordFetched = true
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
	// Note: Use original credential.Service as identifier, not form field
	// This prevents ErrCredentialNotFound if user tries to edit service name
	// For service renaming, a dedicated rename flow should be implemented
	service := ef.credential.Service
	username := ef.form.GetFormItem(1).(*tview.InputField).GetText()
	password := ef.form.GetFormItem(2).(*tview.InputField).GetText()

	// Extract category from input field
	category := ef.form.GetFormItem(3).(*tview.InputField).GetText()
	category = normalizeCategory(category) // Convert "Uncategorized" to empty string

	url := ef.form.GetFormItem(4).(*tview.InputField).GetText()
	notes := ef.form.GetFormItem(5).(*tview.TextArea).GetText()

	// Build UpdateCredentialOpts with only non-empty fields
	opts := models.UpdateCredentialOpts{}

	if username != "" {
		opts.Username = &username
	}

	// Only update password if user changed it (not empty AND different from original)
	// This prevents unnecessary updates when user just views the form
	// T020d: Convert string password to []byte for UpdateOpts
	if password != "" && password != ef.originalPassword {
		passwordBytes := []byte(password)
		opts.Password = &passwordBytes
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
	service := ef.form.GetFormItem(0).(*tview.InputField).GetText()
	if service == "" {
		return fmt.Errorf("service is required")
	}

	// Username is required (minimum validation)
	username := ef.form.GetFormItem(1).(*tview.InputField).GetText()
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
	// If credential has empty category, prefer "Uncategorized" in dropdown
	if ef.credential.Category == "" {
		for i, c := range categories {
			if c == "Uncategorized" {
				return i
			}
		}
		return 0
	}
	// Match credential's category against the categories list
	for i, c := range categories {
		if c == ef.credential.Category {
			return i
		}
	}
	// Return 0 (first category) if no match found
	return 0
}

// wrapInFrame wraps the form in a Frame to add keyboard hints footer.
func (ef *EditForm) wrapInFrame() {
	theme := styles.GetCurrentTheme()

	frame := tview.NewFrame(ef.form)
	frame.SetBorder(true).
		SetTitle(" Edit Credential ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(theme.BorderColor)

	// Add keyboard hints as footer text
	hintsText := "Tab: Next  •  Shift+Tab: Prev  •  Ctrl+S: Save  •  Ctrl+H: Toggle pwd  •  Esc: Cancel"
	frame.AddText(hintsText, false, tview.AlignCenter, theme.TextSecondary)

	ef.Frame = frame
}

// setupKeyboardShortcuts configures form-level keyboard shortcuts.
// Adds Ctrl+S for quick-save and ensures Tab/Shift+Tab stay within form.
func (ef *EditForm) setupKeyboardShortcuts() {
	ef.form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlS:
			// Ctrl+S for quick-save
			ef.onSavePressed()
			return nil

		case tcell.KeyCtrlH:
			// Only toggle if Ctrl modifier is actually pressed
			// Backspace sends KeyCtrlH but without ModCtrl set
			if event.Modifiers()&tcell.ModCtrl != 0 {
				ef.togglePasswordVisibility()
				return nil
			}
			// Regular backspace - pass through to allow deletion
			return event

		case tcell.KeyEscape:
			// Handle Escape to close form
			ef.onCancelPressed()
			return nil

		case tcell.KeyTab, tcell.KeyBacktab:
			// Let form handle Tab/Shift+Tab internally
			// tview's Form will cycle through focusable items and skip disabled ones
			// Return event to allow built-in navigation
			return event
		}
		return event
	})
}

// T048, T049: updatePasswordLabel updates the password field label with strength indicator
func (ef *EditForm) updatePasswordLabel(field *tview.InputField, password []byte) {
	policy := security.DefaultPasswordPolicy
	strength := policy.Strength(password)

	var label string
	if len(password) == 0 {
		label = "Password"
	} else {
		switch strength {
		case security.PasswordStrengthWeak:
			label = "Password [yellow](Weak)[-]"
		case security.PasswordStrengthMedium:
			label = "Password [orange](Medium)[-]"
		case security.PasswordStrengthStrong:
			label = "Password [green](Strong)[-]"
		}
	}

	field.SetLabel(label)
}

// applyStyles applies theme colors and border styling to the form.
func (ef *EditForm) applyStyles() {
	theme := styles.GetCurrentTheme()

	// Apply form-level styling
	styles.ApplyFormStyle(ef.form)

	// Style individual input fields
	// Use BackgroundLight for input fields - lighter than form Background for contrast
	for i := 0; i < 6; i++ {
		item := ef.form.GetFormItem(i)
		switch field := item.(type) {
		case *tview.InputField:
			field.SetFieldBackgroundColor(theme.BackgroundLight).
				SetFieldTextColor(theme.TextPrimary)
		case *tview.TextArea:
			field.SetTextStyle(tcell.StyleDefault.
				Background(theme.BackgroundLight).
				Foreground(theme.TextPrimary))
		case *tview.DropDown:
			field.SetFieldBackgroundColor(theme.BackgroundLight).
				SetFieldTextColor(theme.TextPrimary)
		}
	}

	// Button alignment
	ef.form.SetButtonsAlign(tview.AlignRight)
}

// togglePasswordVisibility switches between masked and plaintext password display.
// Updates both the mask character and the field label to indicate current state.
func (ef *EditForm) togglePasswordVisibility() {
	ef.passwordVisible = !ef.passwordVisible

	passwordField := ef.form.GetFormItem(2).(*tview.InputField)

	if ef.passwordVisible {
		passwordField.SetMaskCharacter(0) // 0 = plaintext (tview convention)
		passwordField.SetLabel("Password [VISIBLE]")
	} else {
		passwordField.SetMaskCharacter('*')
		// Restore label with current strength if password exists
		text := passwordField.GetText()
		if text != "" {
			ef.updatePasswordLabel(passwordField, []byte(text))
		} else {
			passwordField.SetLabel("Password")
		}
	}
}

// SetOnSubmit registers a callback to be invoked after successful update.
func (ef *EditForm) SetOnSubmit(callback func()) {
	ef.onSubmit = callback
}

// SetOnCancel registers a callback to be invoked when cancel is pressed.
func (ef *EditForm) SetOnCancel(callback func()) {
	ef.onCancel = callback
}
