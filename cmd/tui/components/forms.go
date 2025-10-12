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
	"pass-cli/internal/vault"
)

// normalizeCategory converts the "Uncategorized" UI label to empty string for storage.
// Prevents the UI label from leaking into credential data.
func normalizeCategory(c string) string {
	if c == "Uncategorized" {
		return ""
	}
	return c
}

// AddForm provides a modal form for adding new credentials.
// Embeds tview.Form and manages validation and submission.
type AddForm struct {
	*tview.Form

	appState *models.AppState

	passwordVisible bool // Track password visibility state for toggle

	onSubmit func()
	onCancel func()
}

// EditForm provides a modal form for editing existing credentials.
// Embeds tview.Form and pre-populates fields from credential.
type EditForm struct {
	*tview.Form

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
		Form:     form,
		appState: appState,
	}

	af.buildFormFields()
	af.applyStyles()
	af.setupKeyboardShortcuts()

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
	af.AddInputField("Service (UID)", "", 0, nil, nil)
	af.AddInputField("Username", "", 0, nil, nil)
	af.AddPasswordField("Password", "", 0, '*', nil)

	// Optional metadata fields - default to "Uncategorized"

	// Original dropdown approach (commented out for autocomplete field)
	// af.AddDropDown("Category", categories, uncategorizedIndex, nil)

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

  af.AddFormItem(categoryField)
	af.AddInputField("URL", "", 0, nil, nil)
	af.AddTextArea("Notes", "", 0, 5, 0, nil)

	// Action buttons
	af.AddButton("Add", af.onAddPressed)
	af.AddButton("Cancel", af.onCancelPressed)

	// Add keyboard hints at the bottom
	af.addKeyboardHints()
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

	// Extract category from input field (index 3)
  category := af.GetFormItem(3).(*tview.InputField).GetText()
  category = normalizeCategory(category) // Convert "Uncategorized" to empty string
	
	//Original dropdown approach (commented out for autocomplete field) 
	// categoryDropdown := af.GetFormItem(3).(*tview.DropDown)
	// _, category := categoryDropdown.GetCurrentOption()
	// category = normalizeCategory(category) // Convert "Uncategorized" to empty string

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

// addKeyboardHints adds a read-only TextView displaying keyboard shortcuts.
// Helps users discover available commands without needing external documentation.
func (af *AddForm) addKeyboardHints() {
	theme := styles.GetCurrentTheme()

	hintsText := "  Tab: Next field  •  Shift+Tab: Previous  •  Ctrl+S: Add  •  Ctrl+H: Toggle password  •  Esc: Cancel"

	hints := tview.NewTextView()
	hints.SetText(hintsText)
	hints.SetTextAlign(tview.AlignCenter)
	hints.SetTextColor(theme.TextSecondary) // Muted color for hints
	hints.SetBackgroundColor(theme.Background)

	af.AddFormItem(hints)
}

// setupKeyboardShortcuts configures form-level keyboard shortcuts.
// Adds Ctrl+S for quick-save and ensures Tab/Shift+Tab stay within form.
func (af *AddForm) setupKeyboardShortcuts() {
	af.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

		case tcell.KeyTab:
			// Let form handle Tab internally (don't let it escape to app)
			// Return event to allow tview's built-in form navigation
			return event

		case tcell.KeyBacktab: // Shift+Tab
			// Let form handle Shift+Tab internally
			return event
		}
		return event
	})
}

// applyStyles applies theme colors and border styling to the form.
func (af *AddForm) applyStyles() {
	theme := styles.GetCurrentTheme()

	// Apply form-level styling
	styles.ApplyFormStyle(af.Form)

	// Style individual input fields (indices 0-5)
	// Use BackgroundLight for input fields - lighter than form Background for contrast
	for i := 0; i < 6; i++ {
		item := af.GetFormItem(i)
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

	// Set border and title
	af.SetBorder(true).
		SetTitle(" Add Credential ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(theme.BorderColor)

	// Button alignment
	af.SetButtonsAlign(tview.AlignRight)
}

// togglePasswordVisibility switches between masked and plaintext password display.
// Updates both the mask character and the field label to indicate current state.
func (af *AddForm) togglePasswordVisibility() {
	af.passwordVisible = !af.passwordVisible

	passwordField := af.GetFormItem(2).(*tview.InputField)

	if af.passwordVisible {
		passwordField.SetMaskCharacter(0) // 0 = plaintext (tview convention)
		passwordField.SetLabel("Password [VISIBLE]")
	} else {
		passwordField.SetMaskCharacter('*')
		passwordField.SetLabel("Password")
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
		Form:       form,
		appState:   appState,
		credential: credential,
	}

	ef.buildFormFieldsWithValues()
	ef.applyStyles()
	ef.setupKeyboardShortcuts()

	return ef
}

// buildFormFieldsWithValues constructs form fields pre-populated with credential data.
func (ef *EditForm) buildFormFieldsWithValues() {
	categories := ef.getCategories()

	// Pre-populate fields with existing credential data
	// Service field is read-only (cannot be changed in edit mode)
	// Use 0 width to make fields fill available space (prevents black rectangles)
	ef.AddInputField("Service (UID)", ef.credential.Service, 0, nil, nil)
	serviceField := ef.GetFormItem(0).(*tview.InputField)
	serviceField.SetDisabled(true) // Make read-only to prevent confusion

	ef.AddInputField("Username", ef.credential.Username, 0, nil, nil)

	// Password field - defer fetching until user focuses field (lazy loading)
	// This prevents blocking UI and avoids incrementing usage stats on form open
	passwordField := tview.NewInputField().
		SetLabel("Password").
		SetFieldWidth(0).
		SetMaskCharacter('*')

	// Attach focus handler to fetch password lazily
	passwordField.SetFocusFunc(func() {
		ef.fetchPasswordIfNeeded(passwordField)
	})

	ef.AddFormItem(passwordField)

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

	ef.AddFormItem(categoryField)

	ef.AddInputField("URL", ef.credential.URL, 0, nil, nil)
	ef.AddTextArea("Notes", ef.credential.Notes, 0, 5, 0, nil)

	// Action buttons
	ef.AddButton("Save", ef.onSavePressed)
	ef.AddButton("Cancel", ef.onCancelPressed)

	// Add keyboard hints at the bottom
	ef.addKeyboardHints()
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
	username := ef.GetFormItem(1).(*tview.InputField).GetText()
	password := ef.GetFormItem(2).(*tview.InputField).GetText()

	// Extract category from input field (index 3)
	category := ef.GetFormItem(3).(*tview.InputField).GetText()
	category = normalizeCategory(category) // Convert "Uncategorized" to empty string

	url := ef.GetFormItem(4).(*tview.InputField).GetText()
	notes := ef.GetFormItem(5).(*tview.TextArea).GetText()

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

// addKeyboardHints adds a read-only TextView displaying keyboard shortcuts.
// Helps users discover available commands without needing external documentation.
func (ef *EditForm) addKeyboardHints() {
	theme := styles.GetCurrentTheme()

	hintsText := "  Tab: Next field  •  Shift+Tab: Previous  •  Ctrl+S: Save  •  Ctrl+H: Toggle password  •  Esc: Cancel"

	hints := tview.NewTextView()
	hints.SetText(hintsText)
	hints.SetTextAlign(tview.AlignCenter)
	hints.SetTextColor(theme.TextSecondary) // Muted color for hints
	hints.SetBackgroundColor(theme.Background)

	ef.AddFormItem(hints)
}

// setupKeyboardShortcuts configures form-level keyboard shortcuts.
// Adds Ctrl+S for quick-save and ensures Tab/Shift+Tab stay within form.
func (ef *EditForm) setupKeyboardShortcuts() {
	ef.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

		case tcell.KeyTab:
			// Let form handle Tab internally (don't let it escape to app)
			// Return event to allow tview's built-in form navigation
			return event

		case tcell.KeyBacktab: // Shift+Tab
			// Let form handle Shift+Tab internally
			return event
		}
		return event
	})
}

// applyStyles applies theme colors and border styling to the form.
func (ef *EditForm) applyStyles() {
	theme := styles.GetCurrentTheme()

	// Apply form-level styling
	styles.ApplyFormStyle(ef.Form)

	// Style individual input fields (indices 0-5)
	// Use BackgroundLight for input fields - lighter than form Background for contrast
	for i := 0; i < 6; i++ {
		item := ef.GetFormItem(i)
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

	// Set border and title
	ef.SetBorder(true).
		SetTitle(" Edit Credential ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(theme.BorderColor)

	// Button alignment
	ef.SetButtonsAlign(tview.AlignRight)
}

// togglePasswordVisibility switches between masked and plaintext password display.
// Updates both the mask character and the field label to indicate current state.
func (ef *EditForm) togglePasswordVisibility() {
	ef.passwordVisible = !ef.passwordVisible

	passwordField := ef.GetFormItem(2).(*tview.InputField)

	if ef.passwordVisible {
		passwordField.SetMaskCharacter(0) // 0 = plaintext (tview convention)
		passwordField.SetLabel("Password [VISIBLE]")
	} else {
		passwordField.SetMaskCharacter('*')
		passwordField.SetLabel("Password")
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
