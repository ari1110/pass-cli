package tui

import (
	"fmt"
	"time"
)

// formatRelativeTime formats a time as a relative string (e.g., "2 hr ago")
func formatRelativeTime(t time.Time) string {
	if t.IsZero() {
		return "Never"
	}

	duration := time.Since(t)

	if duration < time.Minute {
		return "Just now"
	}
	if duration < time.Hour {
		mins := int(duration.Minutes())
		return fmt.Sprintf("%d min ago", mins)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d hr ago", hours)
	}
	if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d day ago", days)
	}
	if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / 24 / 7)
		return fmt.Sprintf("%d wk ago", weeks)
	}

	// For older dates, show actual date
	return t.Format("2006-01-02")
}
