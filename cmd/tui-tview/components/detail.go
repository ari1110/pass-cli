package components

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/rivo/tview"
	"pass-cli/cmd/tui-tview/models"
	"pass-cli/cmd/tui-tview/styles"
	"pass-cli/internal/vault"
)

// DetailView displays full credential information with password masking and copy support.
// Wraps tview.TextView with credential formatting and clipboard integration.
type DetailView struct {
	*tview.TextView

	appState        *models.AppState
	passwordVisible bool // Toggle for password visibility (false = masked)
}

// NewDetailView creates and configures a new DetailView component.
// Password is masked by default for security.
func NewDetailView(appState *models.AppState) *DetailView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWordWrap(true)

	dv := &DetailView{
		TextView:        textView,
		appState:        appState,
		passwordVisible: false, // Default to masked for security
	}

	dv.applyStyles()
	dv.Refresh()

	return dv
}

// Refresh rebuilds the detail view from the currently selected credential.
// Displays formatted credential information or empty state if no selection.
func (dv *DetailView) Refresh() {
	// Debug: Uncomment to trace selection changes
	// fmt.Printf("[DetailView] Refresh called, selected: %v\n", dv.appState.GetSelectedCredential())

	cred := dv.appState.GetSelectedCredential()

	if cred == nil {
		dv.showEmptyState()
		return
	}

	content := dv.formatCredential(cred)
	dv.SetText(content)
	dv.ScrollToBeginning()
}

// formatCredential creates formatted text display for a credential.
// Uses tview color tags for styling and box drawing characters for sections.
func (dv *DetailView) formatCredential(cred *vault.CredentialMetadata) string {
	var b strings.Builder

	// Header with service name
	b.WriteString("[cyan]═══════════════════════════════════[-]\n")
	b.WriteString(fmt.Sprintf("          [yellow]%s[-]\n", cred.Service))
	b.WriteString("[cyan]═══════════════════════════════════[-]\n\n")

	// Main credential fields
	b.WriteString(fmt.Sprintf("[gray]Service:[-]    [white]%s[-]\n", cred.Service))
	b.WriteString(fmt.Sprintf("[gray]Username:[-]   [white]%s[-]\n", cred.Username))

	// Category (if present)
	if cred.Category != "" {
		b.WriteString(fmt.Sprintf("[gray]Category:[-]   [white]%s[-]\n", cred.Category))
	}

	// URL (if present)
	if cred.URL != "" {
		b.WriteString(fmt.Sprintf("[gray]URL:[-]        [white]%s[-]\n", cred.URL))
	}

	// Password field with masking
	dv.formatPasswordField(&b, cred)

	// Notes (if present)
	if cred.Notes != "" {
		b.WriteString("\n[gray]Notes:[-]\n")
		// Indent multi-line notes
		indentedNotes := strings.ReplaceAll(cred.Notes, "\n", "\n  ")
		b.WriteString(fmt.Sprintf("[white]  %s[-]\n", indentedNotes))
	}

	// Metadata section
	b.WriteString("\n[cyan]═══════════════════════════════════[-]\n")
	b.WriteString("[gray]Metadata[-]\n")
	b.WriteString("[cyan]═══════════════════════════════════[-]\n\n")

	b.WriteString(fmt.Sprintf("[gray]Created:[-]     [white]%s[-]\n", cred.CreatedAt.Format("2006-01-02 03:04 PM")))
	b.WriteString(fmt.Sprintf("[gray]Modified:[-]    [white]%s[-]\n", cred.UpdatedAt.Format("2006-01-02 03:04 PM")))

	if !cred.LastAccessed.IsZero() {
		relativeTime := formatRelativeTime(cred.LastAccessed)
		b.WriteString(fmt.Sprintf("[gray]Last Used:[-]   [white]%s[-]\n", relativeTime))
	}

	if cred.UsageCount > 0 {
		b.WriteString(fmt.Sprintf("[gray]Usage Count:[-] [white]%d times[-]\n", cred.UsageCount))
	}

	// Locations (if any)
	if len(cred.Locations) > 0 {
		b.WriteString(fmt.Sprintf("[gray]Locations:[-]  [white]%d unique locations[-]\n", len(cred.Locations)))
	}

	// Keyboard shortcuts hint
	b.WriteString("\n[gray][e] Edit  [d] Delete  [p] Show/Hide  [c] Copy  [Esc] Back[-]")

	return b.String()
}

// formatPasswordField adds the password field with masking and toggle hint.
// Fetches full credential to display password when visible.
func (dv *DetailView) formatPasswordField(b *strings.Builder, cred *vault.CredentialMetadata) {
	password := "********" // Default masked display
	hint := "  [gray](Press 'p' to reveal)[-]"

	if dv.passwordVisible {
		// Fetch full credential to get password
		fullCred, err := dv.appState.GetFullCredential(cred.Service)
		if err == nil && fullCred != nil {
			password = fullCred.Password
			hint = "  [gray](Press 'p' to hide)[-]"
		} else {
			password = "[red]Error loading password[-]"
			hint = ""
		}
	}

	b.WriteString(fmt.Sprintf("[gray]Password:[-]   [white]%s[-]%s\n", password, hint))
}

// showEmptyState displays a message when no credential is selected.
func (dv *DetailView) showEmptyState() {
	content := `
[cyan]═══════════════════════════════════[-]

        [gray]No Credential Selected[-]

    Select a credential from the list
    to view its details.

[cyan]═══════════════════════════════════[-]
`
	dv.SetText(content)
}

// TogglePasswordVisibility toggles the password display state and refreshes.
// Alternates between masked (••••••••) and plaintext display.
func (dv *DetailView) TogglePasswordVisibility() {
	dv.passwordVisible = !dv.passwordVisible
	dv.Refresh()
}

// CopyPasswordToClipboard copies the selected credential's password to clipboard.
// Returns error if no credential selected or clipboard operation fails.
func (dv *DetailView) CopyPasswordToClipboard() error {
	cred := dv.appState.GetSelectedCredential()
	if cred == nil {
		return fmt.Errorf("no credential selected")
	}

	// Fetch full credential to get password
	fullCred, err := dv.appState.GetFullCredential(cred.Service)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	// Copy password to clipboard
	err = clipboard.WriteAll(fullCred.Password)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	return nil
}

// applyStyles applies theme colors and borders to the detail view.
// Uses theme system for consistent styling.
func (dv *DetailView) applyStyles() {
	styles.ApplyBorderedStyle(dv.TextView, "Details", true)
}
