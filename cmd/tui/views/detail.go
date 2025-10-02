package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"pass-cli/cmd/tui/styles"
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
	// Status bar handles shortcuts, so just return viewport
	return v.viewport.View()
}

// updateContent updates the viewport content
func (v *DetailView) updateContent() {
	var content strings.Builder

	// Title
	content.WriteString(styles.TitleStyle.Render("Credential Details"))
	content.WriteString("\n\n")

	// Notification
	if v.notification != "" {
		content.WriteString(styles.NotificationStyle.Render("✓ " + v.notification))
		content.WriteString("\n\n")
	}

	// Service
	content.WriteString(styles.LabelStyle.Render("Service: "))
	content.WriteString(styles.ValueStyle.Render(v.credential.Service))
	content.WriteString("\n\n")

	// Username
	content.WriteString(styles.LabelStyle.Render("Username: "))
	username := v.credential.Username
	if username == "" {
		username = "(not set)"
	}
	content.WriteString(styles.ValueStyle.Render(username))
	content.WriteString("\n\n")

	// Password
	content.WriteString(styles.LabelStyle.Render("Password: "))
	password := v.credential.Password
	if v.passwordMasked {
		password = strings.Repeat("*", len(password))
	}
	content.WriteString(styles.PasswordStyle.Render(password))
	content.WriteString("\n\n")

	// Notes
	if v.credential.Notes != "" {
		content.WriteString(styles.LabelStyle.Render("Notes: "))
		content.WriteString(styles.ValueStyle.Render(v.credential.Notes))
		content.WriteString("\n\n")
	}

	// Timestamps
	content.WriteString(styles.LabelStyle.Render("Created: "))
	content.WriteString(styles.ValueStyle.Render(formatTime(v.credential.CreatedAt)))
	content.WriteString("\n")

	content.WriteString(styles.LabelStyle.Render("Updated: "))
	content.WriteString(styles.ValueStyle.Render(formatTime(v.credential.UpdatedAt)))
	content.WriteString("\n\n")

	// Usage records
	if len(v.credential.UsageRecord) > 0 {
		content.WriteString(styles.LabelStyle.Render("Usage Records:"))
		content.WriteString("\n\n")

		// Header
		content.WriteString(styles.TableHeaderStyle.Render(fmt.Sprintf("%-40s %-20s %s\n", "Location", "Last Accessed", "Count")))
		content.WriteString(styles.TableDividerStyle.Render(strings.Repeat("-", 80)))
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

			content.WriteString(styles.TableRowStyle.Render(
				fmt.Sprintf("%-40s %-20s %d\n",
					location+repo,
					formatTime(usage.Timestamp),
					usage.Count,
				),
			))
		}
	} else {
		content.WriteString(styles.LabelStyle.Render("Usage Records: "))
		content.WriteString(styles.ValueStyle.Render("None"))
		content.WriteString("\n")
	}

	v.viewport.SetContent(content.String())
}

// renderHelp renders the help line
func (v *DetailView) renderHelp() string {
	help := "m: toggle password | c: copy password | e: edit | d: delete | esc: back to list | q: quit"
	return styles.HelpStyle.Render(help)
}

// GetCredential returns the current credential
func (v *DetailView) GetCredential() *vault.Credential {
	return v.credential
}

// formatTime formats a timestamp
func formatTime(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}
	return t.Format("2006-01-02 15:04:05")
}
