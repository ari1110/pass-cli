package components

import (
	"strings"

	"pass-cli/cmd/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// Breadcrumb displays navigation path
type Breadcrumb struct {
	path  []string
	width int
}

// NewBreadcrumb creates a new breadcrumb component
func NewBreadcrumb() *Breadcrumb {
	return &Breadcrumb{
		path:  []string{"Home"},
		width: 80,
	}
}

// SetPath updates the breadcrumb path
func (b *Breadcrumb) SetPath(path []string) {
	if len(path) == 0 {
		b.path = []string{"Home"}
	} else {
		b.path = path
	}
}

// SetSize updates the breadcrumb width
func (b *Breadcrumb) SetSize(width int) {
	b.width = width
}

// AddSegment adds a segment to the path
func (b *Breadcrumb) AddSegment(segment string) {
	b.path = append(b.path, segment)
}

// PopSegment removes the last segment from the path
func (b *Breadcrumb) PopSegment() {
	if len(b.path) > 1 {
		b.path = b.path[:len(b.path)-1]
	}
}

// Reset resets the breadcrumb to Home
func (b *Breadcrumb) Reset() {
	b.path = []string{"Home"}
}

// GetPath returns the current path
func (b *Breadcrumb) GetPath() []string {
	return b.path
}

// View renders the breadcrumb
func (b *Breadcrumb) View() string {
	if len(b.path) == 0 {
		return ""
	}

	separator := styles.BreadcrumbSeparatorStyle.Render(" > ")
	fullPath := b.renderPath(separator)

	// Check if path fits within width
	if lipgloss.Width(fullPath) <= b.width {
		return fullPath
	}

	// Path is too long, truncate middle segments
	return b.renderTruncatedPath(separator)
}

// renderPath renders the full path with separators
func (b *Breadcrumb) renderPath(separator string) string {
	styledSegments := make([]string, len(b.path))
	for i, segment := range b.path {
		if i == len(b.path)-1 {
			// Last segment is bold/highlighted
			styledSegments[i] = styles.BreadcrumbStyle.Render(segment)
		} else {
			styledSegments[i] = styles.ValueStyle.Render(segment)
		}
	}

	return strings.Join(styledSegments, separator)
}

// renderTruncatedPath renders path with middle segments collapsed
func (b *Breadcrumb) renderTruncatedPath(separator string) string {
	if len(b.path) <= 2 {
		// Can't truncate further, just truncate last segment
		return b.renderPath(separator)
	}

	// Always show first and last segments
	first := styles.ValueStyle.Render(b.path[0])
	last := styles.BreadcrumbStyle.Render(b.path[len(b.path)-1])
	ellipsis := styles.SubtleStyle.Render("...")

	truncated := first + separator + ellipsis + separator + last

	// Check if this fits
	if lipgloss.Width(truncated) <= b.width {
		return truncated
	}

	// If still too long, truncate the last segment
	maxLastLen := b.width - lipgloss.Width(first) - lipgloss.Width(separator) -
		lipgloss.Width(ellipsis) - lipgloss.Width(separator) - 3 // 3 for "..."

	if maxLastLen < 3 {
		// Not enough space, just show first segment
		if lipgloss.Width(first) <= b.width {
			return first
		}
		// Even first segment doesn't fit, truncate it
		return styles.ValueStyle.Render(truncateString(b.path[0], b.width))
	}

	lastTruncated := truncateString(b.path[len(b.path)-1], maxLastLen)
	return first + separator + ellipsis + separator + styles.BreadcrumbStyle.Render(lastTruncated)
}

// Helper function to truncate a string
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 4 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
