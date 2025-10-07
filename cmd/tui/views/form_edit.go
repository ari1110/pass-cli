package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pass-cli/internal/vault"
)

// EditFormView manages the edit credential form
type EditFormView struct {
	serviceInput       textinput.Model
	usernameInput      textinput.Model
	passwordInput      textinput.Model
	notesInput         textarea.Model
	focusedField       FormField
	errorMsg           string
	notification       string
	width              int
	height             int
	originalCredential *vault.Credential
	hasUsageRecords    bool
	usageCount         int
}

// NewEditFormView creates a new edit form pre-filled with credential data
func NewEditFormView(cred *vault.Credential) *EditFormView {
	// Service input (read-only)
	serviceInput := textinput.New()
	serviceInput.SetValue(cred.Service)
	serviceInput.CharLimit = 100

	// Username input
	usernameInput := textinput.New()
	usernameInput.SetValue(cred.Username)
	usernameInput.Placeholder = "username or email"
	usernameInput.Focus()
	usernameInput.CharLimit = 100

	// Password input
	passwordInput := textinput.New()
	passwordInput.SetValue(cred.Password)
	passwordInput.Placeholder = "password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.CharLimit = 100

	// Notes textarea
	notesInput := textarea.New()
	notesInput.SetValue(cred.Notes)
	notesInput.Placeholder = "optional notes"
	notesInput.CharLimit = 500

	// Check for usage records
	hasUsage := len(cred.UsageRecord) > 0
	usageCount := len(cred.UsageRecord)

	return &EditFormView{
		serviceInput:       serviceInput,
		usernameInput:      usernameInput,
		passwordInput:      passwordInput,
		notesInput:         notesInput,
		focusedField:       FieldUsername, // Start at username since service is read-only
		originalCredential: cred,
		hasUsageRecords:    hasUsage,
		usageCount:         usageCount,
	}
}

// SetSize updates dimensions
func (v *EditFormView) SetSize(width, height int) {
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
func (v *EditFormView) Update(msg tea.Msg) (*EditFormView, tea.Cmd) {
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

	// Update focused field (skip service as it's read-only)
	var cmd tea.Cmd
	switch v.focusedField {
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
func (v *EditFormView) View() string {
	var content strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")).
		MarginBottom(1)
	content.WriteString(titleStyle.Render("Edit Credential"))
	content.WriteString("\n\n")

	// Usage warning if credential has usage records
	if v.hasUsageRecords {
		warningStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")). // Yellow/orange
			Bold(true).
			MarginBottom(1)

		warningMsg := fmt.Sprintf("⚠ This credential is used in %d location", v.usageCount)
		if v.usageCount > 1 {
			warningMsg += "s"
		}
		warningMsg += ". Update will affect existing usage."

		content.WriteString(warningStyle.Render(warningMsg))
		content.WriteString("\n")

		// Show usage locations
		if v.originalCredential != nil && len(v.originalCredential.UsageRecord) > 0 {
			locationStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("3")).
				MarginLeft(2)

			content.WriteString(locationStyle.Render("Locations:"))
			content.WriteString("\n")

			for location := range v.originalCredential.UsageRecord {
				content.WriteString(locationStyle.Render("  • " + location))
				content.WriteString("\n")
			}
		}
		content.WriteString("\n")
	}

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

	// Service field (read-only)
	readOnlyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Faint(true)
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Width(12).
		Foreground(lipgloss.Color("240"))
	content.WriteString(labelStyle.Render("Service*") + " " + readOnlyStyle.Render(v.serviceInput.Value()) + " (read-only)")
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
func (v *EditFormView) renderField(label, input string, field FormField) string {
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

// nextField moves to the next field (skipping service since it's read-only)
func (v *EditFormView) nextField() {
	v.blurAll()

	// Cycle through editable fields: Username -> Password -> Notes -> Username
	switch v.focusedField {
	case FieldUsername:
		v.focusedField = FieldPassword
	case FieldPassword:
		v.focusedField = FieldNotes
	case FieldNotes:
		v.focusedField = FieldUsername
	default:
		v.focusedField = FieldUsername
	}

	v.focusCurrent()
}

// prevField moves to the previous field (skipping service since it's read-only)
func (v *EditFormView) prevField() {
	v.blurAll()

	// Cycle through editable fields: Username -> Notes -> Password -> Username
	switch v.focusedField {
	case FieldUsername:
		v.focusedField = FieldNotes
	case FieldPassword:
		v.focusedField = FieldUsername
	case FieldNotes:
		v.focusedField = FieldPassword
	default:
		v.focusedField = FieldUsername
	}

	v.focusCurrent()
}

// focusCurrent focuses the current field
func (v *EditFormView) focusCurrent() {
	switch v.focusedField {
	case FieldUsername:
		v.usernameInput.Focus()
	case FieldPassword:
		v.passwordInput.Focus()
	case FieldNotes:
		v.notesInput.Focus()
	}
}

// blurAll blurs all fields
func (v *EditFormView) blurAll() {
	v.usernameInput.Blur()
	v.passwordInput.Blur()
	v.notesInput.Blur()
}

// GetValues returns the form values
func (v *EditFormView) GetValues() (service, username, password, notes string) {
	return v.serviceInput.Value(), // Service is read-only but still return current value
		strings.TrimSpace(v.usernameInput.Value()),
		v.passwordInput.Value(),
		strings.TrimSpace(v.notesInput.Value())
}

// SetError sets an error message
func (v *EditFormView) SetError(msg string) {
	v.errorMsg = msg
	v.notification = ""
}

// SetNotification sets a notification message
func (v *EditFormView) SetNotification(msg string) {
	v.notification = msg
	v.errorMsg = ""
}

// SetPassword sets the password value
func (v *EditFormView) SetPassword(password string) {
	v.passwordInput.SetValue(password)
}

// HasChanges returns true if any editable field has been modified
func (v *EditFormView) HasChanges() bool {
	if v.originalCredential == nil {
		return false
	}

	// Compare current values with original
	username := strings.TrimSpace(v.usernameInput.Value())
	password := v.passwordInput.Value()
	notes := strings.TrimSpace(v.notesInput.Value())

	return username != v.originalCredential.Username ||
		password != v.originalCredential.Password ||
		notes != v.originalCredential.Notes
}
