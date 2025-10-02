package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"pass-cli/internal/vault"
)

// DetailView displays credential details
type DetailView struct {
	credential     *vault.Credential
	passwordMasked bool
	viewport       viewport.Model
	notification   string
	width          int
	height         int
}

// NewDetailView creates a new detail view
func NewDetailView(credential *vault.Credential) *DetailView {
	vp := viewport.New(0, 0)

	return &DetailView{
		credential:     credential,
		passwordMasked: true,
		viewport:       vp,
	}
}

// SetSize updates the dimensions
func (v *DetailView) SetSize(width, height int) {
	v.width = width
	v.height = height

	// Reserve space for help line
	contentHeight := height - 2
	if contentHeight < 5 {
		contentHeight = 5
	}
	v.viewport.Width = width
	v.viewport.Height = contentHeight

	// Update viewport content
	v.updateContent()
}

// Update handles messages
func (v *DetailView) Update(msg tea.Msg) (*DetailView, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "m":
			v.passwordMasked = !v.passwordMasked
			v.updateContent()
			return v, nil
		case "c":
			if err := clipboard.WriteAll(v.credential.Password); err != nil {
				v.notification = "Failed to copy to clipboard"
			} else {
				v.notification = "Password copied to clipboard"
			}
			v.updateContent()
			// Clear notification after 2 seconds
			return v, tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
				v.notification = ""
				v.updateContent()
				return nil
			})
		}
	}

	v.viewport, cmd = v.viewport.Update(msg)
	return v, cmd
}

// View renders the detail view
func (v *DetailView) View() string {
	help := v.renderHelp()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		v.viewport.View(),
		help,
	)
}

// updateContent updates the viewport content
func (v *DetailView) updateContent() {
	var content strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")). // Cyan
		MarginBottom(1)
	content.WriteString(titleStyle.Render("Credential Details"))
	content.WriteString("\n\n")

	// Notification
	if v.notification != "" {
		notifStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")). // Green
			Bold(true)
		content.WriteString(notifStyle.Render("✓ " + v.notification))
		content.WriteString("\n\n")
	}

	// Service
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("240"))
	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	content.WriteString(labelStyle.Render("Service: "))
	content.WriteString(valueStyle.Render(v.credential.Service))
	content.WriteString("\n\n")

	// Username
	content.WriteString(labelStyle.Render("Username: "))
	username := v.credential.Username
	if username == "" {
		username = "(not set)"
	}
	content.WriteString(valueStyle.Render(username))
	content.WriteString("\n\n")

	// Password
	content.WriteString(labelStyle.Render("Password: "))
	password := v.credential.Password
	if v.passwordMasked {
		password = strings.Repeat("*", len(password))
	}
	passwordStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Background(lipgloss.Color("235"))
	content.WriteString(passwordStyle.Render(" " + password + " "))
	content.WriteString("\n\n")

	// Notes
	if v.credential.Notes != "" {
		content.WriteString(labelStyle.Render("Notes: "))
		content.WriteString(valueStyle.Render(v.credential.Notes))
		content.WriteString("\n\n")
	}

	// Timestamps
	content.WriteString(labelStyle.Render("Created: "))
	content.WriteString(valueStyle.Render(formatTime(v.credential.CreatedAt)))
	content.WriteString("\n")

	content.WriteString(labelStyle.Render("Updated: "))
	content.WriteString(valueStyle.Render(formatTime(v.credential.UpdatedAt)))
	content.WriteString("\n\n")

	// Usage records
	if len(v.credential.UsageRecord) > 0 {
		content.WriteString(labelStyle.Render("Usage Records:"))
		content.WriteString("\n\n")

		tableStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

		// Header
		content.WriteString(tableStyle.Render(fmt.Sprintf("%-40s %-20s %s\n", "Location", "Last Accessed", "Count")))
		content.WriteString(tableStyle.Render(strings.Repeat("-", 80)))
		content.WriteString("\n")

		// Rows
		for _, usage := range v.credential.UsageRecord {
			location := usage.Location
			if len(location) > 38 {
				location = "..." + location[len(location)-35:]
			}

			repo := ""
			if usage.GitRepo != "" {
				repo = fmt.Sprintf(" (%s)", usage.GitRepo)
			}

			content.WriteString(valueStyle.Render(
				fmt.Sprintf("%-40s %-20s %d\n",
					location+repo,
					formatTime(usage.Timestamp),
					usage.Count,
				),
			))
		}
	} else {
		content.WriteString(labelStyle.Render("Usage Records: "))
		content.WriteString(valueStyle.Render("None"))
		content.WriteString("\n")
	}

	v.viewport.SetContent(content.String())
}

// renderHelp renders the help line
func (v *DetailView) renderHelp() string {
	help := "m: toggle password | c: copy password | esc: back to list | q: quit"
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(help)
}

// formatTime formats a timestamp
func formatTime(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	return t.Format("2006-01-02 15:04:05")
}
