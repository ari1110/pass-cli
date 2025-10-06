package components

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"pass-cli/cmd/tui-tview/models"
	"pass-cli/internal/vault"
)

// CredentialTable wraps tview.Table to display credentials in tabular format.
// Supports filtering by category and selection handling.
type CredentialTable struct {
	*tview.Table

	appState      *models.AppState
	filteredCreds []vault.CredentialMetadata
	selectedIndex int
}

// NewCredentialTable creates and configures a new CredentialTable component.
// Creates Table with fixed header row and builds initial credential list.
func NewCredentialTable(appState *models.AppState) *CredentialTable {
	table := tview.NewTable()

	ct := &CredentialTable{
		Table:         table,
		appState:      appState,
		filteredCreds: make([]vault.CredentialMetadata, 0),
		selectedIndex: 0,
	}

	// Configure table
	ct.SetSelectable(true, false) // Select rows, not columns
	ct.SetFixed(1, 0)             // Fix header row

	// Apply styling
	ct.applyStyles()

	// Build header
	ct.buildHeader()

	// Setup selection handler
	ct.SetSelectedFunc(ct.onSelect)

	// Initial population
	ct.Refresh()

	return ct
}

// buildHeader creates the fixed header row with column titles.
// Header row is not selectable and uses accent color.
func (ct *CredentialTable) buildHeader() {
	headers := []string{"Service", "Username", "Last Used"}

	for col, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tcell.NewRGBColor(139, 233, 253)). // Cyan accent
			SetAlign(tview.AlignLeft).
			SetSelectable(false).
			SetExpansion(1)
		ct.SetCell(0, col, cell)
	}
}

// Refresh rebuilds the table from filtered credentials.
// Gets credentials from AppState, filters by selected category, and repopulates rows.
func (ct *CredentialTable) Refresh() {
	// Get credentials and filter by category (thread-safe read)
	allCreds := ct.appState.GetCredentials()
	category := ct.appState.GetSelectedCategory()
	ct.filteredCreds = ct.filterByCategory(allCreds, category)

	// Clear table (keep header at row 0)
	for row := ct.GetRowCount() - 1; row > 0; row-- {
		ct.RemoveRow(row)
	}

	// Populate rows
	ct.populateRows(ct.filteredCreds)

	// Update title with count
	ct.SetTitle(fmt.Sprintf(" Credentials (%d) ", len(ct.filteredCreds)))

	// Restore selection if possible
	if len(ct.filteredCreds) > 0 {
		// Select first row if no selection or out of bounds
		if ct.selectedIndex >= len(ct.filteredCreds) {
			ct.selectedIndex = 0
		}
		ct.Select(ct.selectedIndex+1, 0) // +1 to account for header row
	}
}

// populateRows adds credential rows to the table.
// Stores credential reference in cell metadata for selection handling.
func (ct *CredentialTable) populateRows(credentials []vault.CredentialMetadata) {
	for i, cred := range credentials {
		row := i + 1 // +1 to skip header row

		// Service column
		serviceCell := tview.NewTableCell(cred.Service).
			SetTextColor(tcell.NewRGBColor(248, 248, 242)). // White text
			SetAlign(tview.AlignLeft).
			SetReference(cred) // Store credential in cell for selection

		// Username column
		usernameCell := tview.NewTableCell(cred.Username).
			SetTextColor(tcell.NewRGBColor(189, 147, 249)). // Purple text
			SetAlign(tview.AlignLeft)

		// Last used column
		lastUsed := "Never"
		if !cred.LastAccessed.IsZero() {
			lastUsed = formatRelativeTime(cred.LastAccessed)
		}
		lastUsedCell := tview.NewTableCell(lastUsed).
			SetTextColor(tcell.NewRGBColor(98, 114, 164)). // Gray text
			SetAlign(tview.AlignLeft)

		ct.SetCell(row, 0, serviceCell)
		ct.SetCell(row, 1, usernameCell)
		ct.SetCell(row, 2, lastUsedCell)
	}
}

// onSelect handles row selection by updating AppState with selected credential.
// Called when user presses Enter on a table row.
func (ct *CredentialTable) onSelect(row, column int) {
	if row == 0 {
		return // Header row, ignore
	}

	// Get credential from cell reference
	cell := ct.GetCell(row, 0)
	if cell != nil {
		if cred, ok := cell.GetReference().(vault.CredentialMetadata); ok {
			ct.selectedIndex = row - 1 // Store index without header offset
			ct.appState.SetSelectedCredential(&cred)
		}
	}
}

// filterByCategory filters credentials by selected category.
// Empty category returns all credentials.
func (ct *CredentialTable) filterByCategory(creds []vault.CredentialMetadata, category string) []vault.CredentialMetadata {
	if category == "" {
		return creds // Show all
	}

	filtered := make([]vault.CredentialMetadata, 0)
	for _, cred := range creds {
		if cred.Service == category {
			filtered = append(filtered, cred)
		}
	}
	return filtered
}

// applyStyles applies borders, colors, and title to the table.
// Uses rounded borders with cyan accent color and dark background.
func (ct *CredentialTable) applyStyles() {
	ct.SetBorder(true).
		SetTitle(" Credentials ").
		SetTitleAlign(tview.AlignLeft).
		SetBorderColor(tcell.NewRGBColor(139, 233, 253)). // Cyan border
		SetBackgroundColor(tcell.NewRGBColor(40, 42, 54))  // Dark background
}

// formatRelativeTime formats a timestamp as a relative time string.
// Examples: "2m ago", "5h ago", "3d ago"
func formatRelativeTime(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "Just now"
	case duration < time.Hour:
		return fmt.Sprintf("%dm ago", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(duration.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(duration.Hours()/24))
	}
}
