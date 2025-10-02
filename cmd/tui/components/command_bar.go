package components

import (
	"strings"

	"pass-cli/cmd/tui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Command represents a parsed command
type Command struct {
	Name string
	Args []string
}

// CommandBar represents a vim-style command input bar
type CommandBar struct {
	input      textinput.Model
	history    []string
	historyIdx int
	error      string
	width      int
	height     int
}

// NewCommandBar creates a new command bar
func NewCommandBar() *CommandBar {
	ti := textinput.New()
	ti.Placeholder = "Enter command (e.g., :add, :search, :help, :quit)"
	ti.CharLimit = 200
	ti.Width = 50

	return &CommandBar{
		input:      ti,
		history:    []string{},
		historyIdx: -1,
	}
}

// SetSize updates the command bar dimensions
func (c *CommandBar) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.input.Width = width - 2 // Account for prompt (:) and padding
}

// Focus puts focus on the command bar input
func (c *CommandBar) Focus() tea.Cmd {
	c.input.SetValue(":")
	c.input.CursorEnd()
	c.error = ""
	return c.input.Focus()
}

// Blur removes focus from the command bar
func (c *CommandBar) Blur() {
	c.input.Blur()
	c.input.SetValue("")
	c.error = ""
	c.historyIdx = -1
}

// Update handles tea messages
func (c *CommandBar) Update(msg tea.Msg) (*CommandBar, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			// Navigate history backwards
			if len(c.history) > 0 {
				if c.historyIdx < len(c.history)-1 {
					c.historyIdx++
					c.input.SetValue(":" + c.history[len(c.history)-1-c.historyIdx])
					c.input.CursorEnd()
				}
			}
			return c, nil

		case "down":
			// Navigate history forwards
			if c.historyIdx > 0 {
				c.historyIdx--
				c.input.SetValue(":" + c.history[len(c.history)-1-c.historyIdx])
				c.input.CursorEnd()
			} else if c.historyIdx == 0 {
				c.historyIdx = -1
				c.input.SetValue(":")
				c.input.CursorEnd()
			}
			return c, nil
		}
	}

	c.input, cmd = c.input.Update(msg)
	return c, cmd
}

// GetCommand parses and returns the current command, adding it to history
func (c *CommandBar) GetCommand() (*Command, error) {
	input := c.input.Value()

	// Remove leading :
	if !strings.HasPrefix(input, ":") {
		c.SetError("Commands must start with :")
		return nil, nil
	}
	input = strings.TrimPrefix(input, ":")
	input = strings.TrimSpace(input)

	if input == "" {
		c.SetError("Empty command")
		return nil, nil
	}

	// Add to history
	c.history = append(c.history, input)
	c.historyIdx = -1

	// Parse command
	parts := strings.Fields(input)
	cmd := &Command{
		Name: parts[0],
	}

	if len(parts) > 1 {
		cmd.Args = parts[1:]
	} else {
		cmd.Args = []string{}
	}

	return cmd, nil
}

// SetError sets an error message to display
func (c *CommandBar) SetError(err string) {
	c.error = err
}

// ClearError clears the error message
func (c *CommandBar) ClearError() {
	c.error = ""
}

// View renders the command bar
func (c *CommandBar) View() string {
	var lines []string

	// Show error if present
	if c.error != "" {
		errorLine := styles.ErrorStyle.Render("Error: " + c.error)
		lines = append(lines, errorLine)
	}

	// Render input
	inputView := c.input.View()
	lines = append(lines, inputView)

	// Add help text
	helpText := styles.SubtleStyle.Render("Commands: :add [service] | :search [query] | :category [name] | :help | :quit  •  Esc to cancel")
	lines = append(lines, helpText)

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return content
}

// GetHistory returns the command history
func (c *CommandBar) GetHistory() []string {
	return c.history
}
