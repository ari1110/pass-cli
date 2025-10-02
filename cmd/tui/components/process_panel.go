package components

import (
	"fmt"
	"strings"
	"time"

	"pass-cli/cmd/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ProcessStatus represents the status of a process
type ProcessStatus string

const (
	ProcessPending ProcessStatus = "pending"
	ProcessRunning ProcessStatus = "running"
	ProcessSuccess ProcessStatus = "success"
	ProcessFailed  ProcessStatus = "failed"
)

// Process represents an operation being tracked
type Process struct {
	ID          string
	Description string
	Status      ProcessStatus
	Error       string
	Timestamp   time.Time
}

// ProcessPanel displays async operation feedback
type ProcessPanel struct {
	processes  []Process
	maxDisplay int
	width      int
	height     int
}

// NewProcessPanel creates a new process panel
func NewProcessPanel() *ProcessPanel {
	return &ProcessPanel{
		processes:  []Process{},
		maxDisplay: 5, // Show maximum 5 processes
	}
}

// SetSize updates the panel dimensions
func (p *ProcessPanel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// AddProcess adds a new process to track
func (p *ProcessPanel) AddProcess(id, description string) {
	process := Process{
		ID:          id,
		Description: description,
		Status:      ProcessRunning,
		Timestamp:   time.Now(),
	}

	p.processes = append([]Process{process}, p.processes...)

	// Trim to max display count
	if len(p.processes) > p.maxDisplay {
		p.processes = p.processes[:p.maxDisplay]
	}
}

// UpdateProcess updates the status of a process
func (p *ProcessPanel) UpdateProcess(id string, status ProcessStatus, errorMsg string) {
	for i := range p.processes {
		if p.processes[i].ID == id {
			p.processes[i].Status = status
			p.processes[i].Error = errorMsg
			return
		}
	}
}

// RemoveProcess removes a process by ID
func (p *ProcessPanel) RemoveProcess(id string) {
	for i, proc := range p.processes {
		if proc.ID == id {
			p.processes = append(p.processes[:i], p.processes[i+1:]...)
			return
		}
	}
}

// HasActiveProcesses returns true if there are any running or pending processes
func (p *ProcessPanel) HasActiveProcesses() bool {
	for _, proc := range p.processes {
		if proc.Status == ProcessRunning || proc.Status == ProcessPending {
			return true
		}
	}
	return false
}

// GetProcesses returns all processes (for inspection)
func (p *ProcessPanel) GetProcesses() []Process {
	return p.processes
}

// View renders the process panel
func (p *ProcessPanel) View() string {
	if len(p.processes) == 0 {
		return ""
	}

	title := styles.PanelTitleStyle.Render("⚙️  Processes")

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	// Render processes (most recent first)
	for _, proc := range p.processes {
		line := p.renderProcess(proc)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	return content
}

// renderProcess renders a single process line
func (p *ProcessPanel) renderProcess(proc Process) string {
	icon := styles.GetStatusIcon(string(proc.Status))

	var statusStyle lipgloss.Style
	switch proc.Status {
	case ProcessSuccess:
		statusStyle = styles.SuccessStyle
	case ProcessFailed:
		statusStyle = styles.ErrorStyle
	case ProcessRunning, ProcessPending:
		statusStyle = styles.ValueStyle
	default:
		statusStyle = styles.SubtleStyle
	}

	// Build the process line
	description := truncateToWidth(proc.Description, p.width-20) // Reserve space for icon and timestamp

	var line string
	if proc.Status == ProcessFailed && proc.Error != "" {
		line = fmt.Sprintf("%s %s - %s",
			statusStyle.Render(icon),
			statusStyle.Render(description),
			styles.ErrorStyle.Render(proc.Error),
		)
	} else {
		line = fmt.Sprintf("%s %s",
			statusStyle.Render(icon),
			statusStyle.Render(description),
		)
	}

	return "  " + line
}

// Helper functions

func truncateToWidth(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth < 4 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}
