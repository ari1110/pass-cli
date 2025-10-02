package components

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"pass-cli/cmd/tui/styles"
	"pass-cli/internal/vault"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MetadataPanel displays detailed information about a selected credential
type MetadataPanel struct {
	credential     *vault.Credential
	passwordMasked bool
	viewport       viewport.Model
	width          int
	height         int
	focused        bool
}

// NewMetadataPanel creates a new metadata panel
func NewMetadataPanel() *MetadataPanel {
	vp := viewport.New(25, 10)

	return &MetadataPanel{
		credential:     nil,
		passwordMasked: true, // Default to masked
		viewport:       vp,
		focused:        false,
	}
}

// SetSize updates the panel dimensions
func (m *MetadataPanel) SetSize(width, height int) {
	m.width = width
	m.height = height

	// Reserve space for title (1), borders/padding (2)
	contentHeight := height - 3
	if contentHeight < 5 {
		contentHeight = 5
	}

	m.viewport.Width = width - 2 // Account for padding
	m.viewport.Height = contentHeight

	// Update content if credential is set
	if m.credential != nil {
		m.updateViewportContent()
	}
}

// SetFocus sets the focus state of the panel
func (m *MetadataPanel) SetFocus(focused bool) {
	m.focused = focused
}

// SetCredential updates the displayed credential
func (m *MetadataPanel) SetCredential(cred *vault.Credential) {
	m.credential = cred
	m.passwordMasked = true // Reset to masked when switching credentials
	m.updateViewportContent()
}

// TogglePasswordMask toggles the password visibility
func (m *MetadataPanel) TogglePasswordMask() {
	m.passwordMasked = !m.passwordMasked
	m.updateViewportContent()
}

// Update handles tea messages
func (m *MetadataPanel) Update(msg tea.Msg) (*MetadataPanel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Only handle keys when focused
		if !m.focused {
			return m, nil
		}

		switch msg.String() {
		case "m":
			// Toggle password mask
			m.TogglePasswordMask()
			return m, nil
		}
	}

	// Update viewport for scrolling
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// updateViewportContent refreshes the viewport with current credential data
func (m *MetadataPanel) updateViewportContent() {
	content := m.renderCredentialDetails()
	m.viewport.SetContent(content)
}

// renderCredentialDetails renders the credential information
func (m *MetadataPanel) renderCredentialDetails() string {
	if m.credential == nil {
		return styles.SubtleStyle.Render("No credential selected\n\nSelect a credential from the sidebar to view details.")
	}

	var sections []string

	// Service (Header)
	serviceHeader := styles.TitleStyle.Render(m.credential.Service)
	sections = append(sections, serviceHeader)

	// Username
	sections = append(sections, "")
	sections = append(sections, styles.LabelStyle.Render("Username:"))
	if m.credential.Username != "" {
		sections = append(sections, styles.ValueStyle.Render("  "+m.credential.Username))
	} else {
		sections = append(sections, styles.SubtleStyle.Render("  (not set)"))
	}

	// Password
	sections = append(sections, "")
	sections = append(sections, styles.LabelStyle.Render("Password:"))
	if m.credential.Password != "" {
		if m.passwordMasked {
			sections = append(sections, styles.PasswordStyle.Render("  "+maskPassword(m.credential.Password)))
			sections = append(sections, styles.SubtleStyle.Render("  [m] to reveal"))
		} else {
			sections = append(sections, styles.ValueStyle.Render("  "+m.credential.Password))
			sections = append(sections, styles.SubtleStyle.Render("  [m] to mask"))
		}
	} else {
		sections = append(sections, styles.SubtleStyle.Render("  (not set)"))
	}

	// Notes
	if m.credential.Notes != "" {
		sections = append(sections, "")
		sections = append(sections, styles.LabelStyle.Render("Notes:"))
		// Word wrap notes
		wrappedNotes := wordWrap(m.credential.Notes, m.viewport.Width-2)
		for _, line := range strings.Split(wrappedNotes, "\n") {
			sections = append(sections, styles.ValueStyle.Render("  "+line))
		}
	}

	// Timestamps
	sections = append(sections, "")
	sections = append(sections, styles.LabelStyle.Render("Created:"))
	sections = append(sections, styles.ValueStyle.Render("  "+formatTime(m.credential.CreatedAt)))

	sections = append(sections, "")
	sections = append(sections, styles.LabelStyle.Render("Updated:"))
	sections = append(sections, styles.ValueStyle.Render("  "+formatTime(m.credential.UpdatedAt)))

	// Usage Records
	if len(m.credential.UsageRecord) > 0 {
		sections = append(sections, "")
		sections = append(sections, styles.LabelStyle.Render("Usage Records:"))
		sections = append(sections, "")

		// Convert map to slice and sort by last accessed
		records := make([]vault.UsageRecord, 0, len(m.credential.UsageRecord))
		for _, record := range m.credential.UsageRecord {
			records = append(records, record)
		}
		sort.Slice(records, func(i, j int) bool {
			return records[i].Timestamp.After(records[j].Timestamp)
		})

		// Render table
		sections = append(sections, styles.TableHeaderStyle.Render(fmt.Sprintf("  %-30s  %8s  %s", "Location", "Count", "Last Accessed")))
		sections = append(sections, styles.TableDividerStyle.Render("  "+strings.Repeat("─", m.viewport.Width-4)))

		for _, record := range records {
			location := truncateMiddle(record.Location, 30)
			lastAccessed := formatRelativeTime(record.Timestamp)
			row := fmt.Sprintf("  %-30s  %8d  %s", location, record.Count, lastAccessed)
			sections = append(sections, styles.TableRowStyle.Render(row))
		}
	} else {
		sections = append(sections, "")
		sections = append(sections, styles.SubtleStyle.Render("No usage records yet"))
	}

	return strings.Join(sections, "\n")
}

// View renders the metadata panel
func (m *MetadataPanel) View() string {
	titleStyle := styles.InactivePanelTitleStyle
	if m.focused {
		titleStyle = styles.PanelTitleStyle
	}

	title := titleStyle.Render("ℹ️  Details")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		m.viewport.View(),
	)

	return content
}

// Helper functions

func maskPassword(password string) string {
	if len(password) == 0 {
		return ""
	}
	return strings.Repeat("•", min(len(password), 12))
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "(not set)"
	}

	now := time.Now()
	duration := now.Sub(t)

	// Absolute format for very old dates
	if duration > 365*24*time.Hour {
		return t.Format("Jan 2, 2006")
	}

	return t.Format("Jan 2, 2006 15:04")
}

func formatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}

	now := time.Now()
	duration := now.Sub(t)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	} else if duration < 365*24*time.Hour {
		months := int(duration.Hours() / 24 / 30)
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	years := int(duration.Hours() / 24 / 365)
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

func truncateMiddle(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	if maxLen < 7 {
		return s[:maxLen]
	}

	// Show first and last parts with ... in middle
	sideLen := (maxLen - 3) / 2
	return s[:sideLen] + "..." + s[len(s)-sideLen:]
}

func wordWrap(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return strings.Join(lines, "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
