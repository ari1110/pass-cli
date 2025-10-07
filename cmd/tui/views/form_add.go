package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FormField represents which field is focused
type FormField int

const (
	FieldService FormField = iota
	FieldUsername
	FieldPassword
	FieldNotes
)

// AddFormView manages the add credential form
type AddFormView struct {
	serviceInput  textinput.Model
	usernameInput textinput.Model
	passwordInput textinput.Model
	notesInput    textarea.Model
	focusedField  FormField
	errorMsg      string
	notification  string
	width         int
	height        int
}

// NewAddFormView creates a new add form
func NewAddFormView() *AddFormView {
	// Service input
	serviceInput := textinput.New()
	serviceInput.Placeholder = "e.g., github, aws-prod"
	serviceInput.Focus()
	serviceInput.CharLimit = 100

	// Username input
	usernameInput := textinput.New()
	usernameInput.Placeholder = "username or email"
	usernameInput.CharLimit = 100

	// Password input
	passwordInput := textinput.New()
	passwordInput.Placeholder = "leave empty to generate"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.CharLimit = 100

	// Notes textarea
	notesInput := textarea.New()
	notesInput.Placeholder = "optional notes"
	notesInput.CharLimit = 500

	return &AddFormView{
		serviceInput:  serviceInput,
		usernameInput: usernameInput,
		passwordInput: passwordInput,
		notesInput:    notesInput,
		focusedField:  FieldService,
	}
}

// SetSize updates dimensions
func (v *AddFormView) SetSize(width, height int) {
	v.width = width
	v.height = height

	inputWidth := width - 20
	if inputWidth < 40 {
		inputWidth = 40
	}

	v.serviceInput.Width = inputWidth
	v.usernameInput.Width = inputWidth
	v.passwordInput.Width = inputWidth
	v.notesInput.SetWidth(inputWidth)
	v.notesInput.SetHeight(3)
}

// Update handles messages
func (v *AddFormView) Update(msg tea.Msg) (*AddFormView, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			v.nextField()
			return v, nil
		case "shift+tab", "up":
			v.prevField()
			return v, nil
		}
	}

	// Update focused field
	var cmd tea.Cmd
	switch v.focusedField {
	case FieldService:
		v.serviceInput, cmd = v.serviceInput.Update(msg)
	case FieldUsername:
		v.usernameInput, cmd = v.usernameInput.Update(msg)
	case FieldPassword:
		v.passwordInput, cmd = v.passwordInput.Update(msg)
	case FieldNotes:
		v.notesInput, cmd = v.notesInput.Update(msg)
	}
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

// View renders the form
func (v *AddFormView) View() string {
	var content strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")).
		MarginBottom(1)
	content.WriteString(titleStyle.Render("Add New Credential"))
	content.WriteString("\n\n")

	// Notification
	if v.notification != "" {
		notifStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)
		content.WriteString(notifStyle.Render("✓ " + v.notification))
		content.WriteString("\n\n")
	}

	// Error message
	if v.errorMsg != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).
			Bold(true)
		content.WriteString(errorStyle.Render("✗ " + v.errorMsg))
		content.WriteString("\n\n")
	}

	// Service field
	content.WriteString(v.renderField("Service*", v.serviceInput.View(), FieldService))
	content.WriteString("\n")

	// Username field
	content.WriteString(v.renderField("Username", v.usernameInput.View(), FieldUsername))
	content.WriteString("\n")

	// Password field
	content.WriteString(v.renderField("Password", v.passwordInput.View(), FieldPassword))
	content.WriteString("\n")

	// Notes field
	content.WriteString(v.renderField("Notes", v.notesInput.View(), FieldNotes))
	content.WriteString("\n\n")

	// Help
	help := "tab/↓: next field | shift+tab/↑: prev field | g: generate password | ctrl+s: save | esc: cancel"
	content.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(help))

	return content.String()
}

// renderField renders a form field with focus indicator
func (v *AddFormView) renderField(label, input string, field FormField) string {
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Width(12)

	if v.focusedField == field {
		labelStyle = labelStyle.Foreground(lipgloss.Color("6")) // Cyan when focused
	} else {
		labelStyle = labelStyle.Foreground(lipgloss.Color("240")) // Gray when not focused
	}

	return labelStyle.Render(label) + " " + input
}

// nextField moves to the next field
func (v *AddFormView) nextField() {
	v.blurAll()
	v.focusedField = (v.focusedField + 1) % 4
	v.focusCurrent()
}

// prevField moves to the previous field
func (v *AddFormView) prevField() {
	v.blurAll()
	if v.focusedField == 0 {
		v.focusedField = 3
	} else {
		v.focusedField--
	}
	v.focusCurrent()
}

// focusCurrent focuses the current field
func (v *AddFormView) focusCurrent() {
	switch v.focusedField {
	case FieldService:
		v.serviceInput.Focus()
	case FieldUsername:
		v.usernameInput.Focus()
	case FieldPassword:
		v.passwordInput.Focus()
	case FieldNotes:
		v.notesInput.Focus()
	}
}

// blurAll blurs all fields
func (v *AddFormView) blurAll() {
	v.serviceInput.Blur()
	v.usernameInput.Blur()
	v.passwordInput.Blur()
	v.notesInput.Blur()
}

// GetValues returns the form values
func (v *AddFormView) GetValues() (service, username, password, notes string) {
	return strings.TrimSpace(v.serviceInput.Value()),
		strings.TrimSpace(v.usernameInput.Value()),
		v.passwordInput.Value(),
		strings.TrimSpace(v.notesInput.Value())
}

// SetError sets an error message
func (v *AddFormView) SetError(msg string) {
	v.errorMsg = msg
	v.notification = ""
}

// SetNotification sets a notification message
func (v *AddFormView) SetNotification(msg string) {
	v.notification = msg
	v.errorMsg = ""
}

// SetPassword sets the password value
func (v *AddFormView) SetPassword(password string) {
	v.passwordInput.SetValue(password)
}

// HasChanges returns true if any field has been modified
func (v *AddFormView) HasChanges() bool {
	service, username, password, notes := v.GetValues()
	return service != "" || username != "" || password != "" || notes != ""
}
