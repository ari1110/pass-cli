package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pass-cli/internal/vault"
)

// ConfirmType represents the type of confirmation dialog
type ConfirmType int

const (
	ConfirmSimple ConfirmType = iota // Simple y/n confirmation
	ConfirmTyped                      // Requires typing service name
	ConfirmDelete                     // Delete confirmation with credential info
)

// ConfirmView displays a confirmation dialog
type ConfirmView struct {
	confirmType    ConfirmType
	message        string
	credential     *vault.Credential // For delete confirmations
	textInput      textinput.Model   // For typed confirmations
	errorMsg       string
	width          int
	height         int
}

// NewConfirmView creates a new simple confirmation dialog
func NewConfirmView(message string) *ConfirmView {
	return &ConfirmView{
		confirmType: ConfirmSimple,
		message:     message,
	}
}

// NewDeleteConfirmView creates a delete confirmation dialog
func NewDeleteConfirmView(cred *vault.Credential) *ConfirmView {
	hasUsage := len(cred.UsageRecord) > 0
	confirmType := ConfirmSimple

	var input textinput.Model
	if hasUsage {
		// High-risk deletion requires typed confirmation
		confirmType = ConfirmTyped
		input = textinput.New()
		input.Placeholder = "type service name to confirm"
		input.Focus()
		input.CharLimit = 100
	}

	return &ConfirmView{
		confirmType: confirmType,
		credential:  cred,
		textInput:   input,
	}
}

// SetSize updates dimensions
func (v *ConfirmView) SetSize(width, height int) {
	v.width = width
	v.height = height

	if v.confirmType == ConfirmTyped {
		inputWidth := 50
		if inputWidth > width-20 {
			inputWidth = width - 20
		}
		v.textInput.Width = inputWidth
	}
}

// Update handles messages
func (v *ConfirmView) Update(msg tea.Msg) (*ConfirmView, tea.Cmd) {
	if v.confirmType == ConfirmTyped {
		var cmd tea.Cmd
		v.textInput, cmd = v.textInput.Update(msg)
		return v, cmd
	}
	return v, nil
}

// View renders the confirmation dialog
func (v *ConfirmView) View() string {
	var content strings.Builder

	if v.credential != nil {
		// Delete confirmation
		content.WriteString(v.renderDeleteConfirmation())
	} else {
		// Simple confirmation (discard changes, etc.)
		content.WriteString(v.renderSimpleConfirmation())
	}

	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("3")). // Yellow for warning
		Padding(1, 2).
		Width(70).
		Align(lipgloss.Left)

	dialog := dialogStyle.Render(content.String())

	return lipgloss.Place(
		v.width,
		v.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}

// renderSimpleConfirmation renders a simple y/n confirmation
func (v *ConfirmView) renderSimpleConfirmation() string {
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	return messageStyle.Render(v.message) + "\n\n" +
		promptStyle.Render("y: yes, discard changes | n: no, go back")
}

// renderDeleteConfirmation renders a delete confirmation with credential details
func (v *ConfirmView) renderDeleteConfirmation() string {
	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")). // Red for danger
		Bold(true)
	content.WriteString(titleStyle.Render(fmt.Sprintf("Delete '%s'?", v.credential.Service)))
	content.WriteString("\n\n")

	// Username for verification
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	content.WriteString(labelStyle.Render("Username: "))
	username := v.credential.Username
	if username == "" {
		username = "(not set)"
	}
	content.WriteString(valueStyle.Render(username))
	content.WriteString("\n\n")

	// Usage warning if applicable
	if len(v.credential.UsageRecord) > 0 {
		warningStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")). // Red
			Bold(true)

		locationCount := len(v.credential.UsageRecord)
		warningMsg := fmt.Sprintf("⚠ Used in %d location", locationCount)
		if locationCount > 1 {
			warningMsg += "s"
		}
		warningMsg += ":"

		content.WriteString(warningStyle.Render(warningMsg))
		content.WriteString("\n")

		// List usage locations
		locationStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))

		for location := range v.credential.UsageRecord {
			// Truncate long paths
			displayLocation := location
			if len(displayLocation) > 60 {
				displayLocation = "..." + displayLocation[len(displayLocation)-57:]
			}
			content.WriteString(locationStyle.Render("  • " + displayLocation))
			content.WriteString("\n")
		}
		content.WriteString("\n")
	}

	// Prompt based on confirmation type
	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	if v.confirmType == ConfirmTyped {
		// Typed confirmation required
		content.WriteString(promptStyle.Render("Type service name to confirm:"))
		content.WriteString("\n")
		content.WriteString(v.textInput.View())
		content.WriteString("\n\n")

		// Error message if any
		if v.errorMsg != "" {
			errorStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Bold(true)
			content.WriteString(errorStyle.Render("✗ " + v.errorMsg))
			content.WriteString("\n\n")
		}

		content.WriteString(promptStyle.Render("enter: confirm | esc: cancel"))
	} else {
		// Simple y/n confirmation
		content.WriteString(promptStyle.Render("y: yes, delete | n: no, cancel"))
	}

	return content.String()
}

// GetTypedValue returns the typed input value
func (v *ConfirmView) GetTypedValue() string {
	if v.confirmType == ConfirmTyped {
		return strings.TrimSpace(v.textInput.Value())
	}
	return ""
}

// SetError sets an error message
func (v *ConfirmView) SetError(msg string) {
	v.errorMsg = msg
}

// IsTypedConfirmation returns true if this is a typed confirmation
func (v *ConfirmView) IsTypedConfirmation() bool {
	return v.confirmType == ConfirmTyped
}

// GetCredential returns the credential being confirmed for deletion
func (v *ConfirmView) GetCredential() *vault.Credential {
	return v.credential
}
