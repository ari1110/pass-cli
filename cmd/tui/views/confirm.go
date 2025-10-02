package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmView displays a confirmation dialog
type ConfirmView struct {
	message string
	width   int
	height  int
}

// NewConfirmView creates a new confirmation dialog
func NewConfirmView(message string) *ConfirmView {
	return &ConfirmView{
		message: message,
	}
}

// SetSize updates dimensions
func (v *ConfirmView) SetSize(width, height int) {
	v.width = width
	v.height = height
}

// Update handles messages
func (v *ConfirmView) Update(msg tea.Msg) (*ConfirmView, tea.Cmd) {
	return v, nil
}

// View renders the confirmation dialog
func (v *ConfirmView) View() string {
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("3")). // Yellow for warning
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Align(lipgloss.Center)

	promptStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		MarginTop(1)

	content := messageStyle.Render(v.message) + "\n\n" +
		promptStyle.Render("y: yes, discard changes | n: no, go back")

	dialog := dialogStyle.Render(content)

	// Center the dialog
	verticalPadding := (v.height - lipgloss.Height(dialog)) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	return lipgloss.Place(
		v.width,
		v.height,
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)
}
