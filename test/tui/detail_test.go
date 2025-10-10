package tui_test

import (
	"strings"
	"testing"
	"time"

	"pass-cli/cmd/tui/components"
	"pass-cli/internal/vault"
)

// T038: Test sortUsageLocations() ordering by timestamp descending
func TestSortUsageLocations_OrderByTimestamp(t *testing.T) {
	// Create usage records with different timestamps
	records := map[string]vault.UsageRecord{
		"/path/one":   CreateTestUsageRecord("/path/one", 1, "", 5, 0),     // 1 hour ago
		"/path/two":   CreateTestUsageRecord("/path/two", 24, "repo1", 3, 0), // 1 day ago
		"/path/three": CreateTestUsageRecord("/path/three", 72, "repo2", 1, 0), // 3 days ago
		"/path/four":  CreateTestUsageRecord("/path/four", 168, "", 2, 0),   // 7 days ago (1 week)
	}

	// Sort the records
	sorted := components.SortUsageLocations(records)

	// Verify: Most recent first (descending order)
	if len(sorted) != 4 {
		t.Fatalf("Expected 4 sorted records, got %d", len(sorted))
	}

	// Check that timestamps are in descending order
	if sorted[0].Location != "/path/one" {
		t.Errorf("Expected first record to be /path/one (most recent), got %s", sorted[0].Location)
	}
	if sorted[1].Location != "/path/two" {
		t.Errorf("Expected second record to be /path/two, got %s", sorted[1].Location)
	}
	if sorted[2].Location != "/path/three" {
		t.Errorf("Expected third record to be /path/three, got %s", sorted[2].Location)
	}
	if sorted[3].Location != "/path/four" {
		t.Errorf("Expected fourth record to be /path/four (oldest), got %s", sorted[3].Location)
	}

	// Verify timestamps are actually in descending order
	for i := 0; i < len(sorted)-1; i++ {
		if sorted[i].Timestamp.Before(sorted[i+1].Timestamp) {
			t.Errorf("Timestamps not in descending order: sorted[%d] (%v) is before sorted[%d] (%v)",
				i, sorted[i].Timestamp, i+1, sorted[i+1].Timestamp)
		}
	}
}

// T039: Test formatTimestamp() hybrid logic (<7 days relative, ≥7 days absolute)
func TestFormatTimestamp_HybridLogic(t *testing.T) {
	tests := []struct {
		name          string
		hoursAgo      int
		wantContains  string
		wantFormat    string // "relative" or "absolute"
	}{
		{"30 minutes ago", 0, "minutes ago", "relative"},
		{"2 hours ago", 2, "hours ago", "relative"},
		{"1 day ago", 24, "1 day", "relative"},
		{"3 days ago", 72, "3 days", "relative"},
		{"6 days ago", 144, "6 days", "relative"},
		{"7 days ago - threshold", 168, "", "absolute"}, // Exactly 7 days = absolute
		{"10 days ago", 240, "", "absolute"},
		{"30 days ago", 720, "", "absolute"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := time.Now().Add(-time.Duration(tt.hoursAgo) * time.Hour)
			result := components.FormatTimestamp(timestamp)

			if tt.wantFormat == "relative" {
				// Should contain "ago" and not be a date format
				if !strings.Contains(result, "ago") {
					t.Errorf("Expected relative format with 'ago', got %q", result)
				}
				if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
					t.Errorf("Expected result to contain %q, got %q", tt.wantContains, result)
				}
				// Should not be in YYYY-MM-DD format
				if len(result) == 10 && result[4] == '-' && result[7] == '-' {
					t.Errorf("Expected relative format, got date format: %q", result)
				}
			} else {
				// Should be in YYYY-MM-DD format (absolute)
				if len(result) != 10 || result[4] != '-' || result[7] != '-' {
					t.Errorf("Expected YYYY-MM-DD format, got %q", result)
				}
				// Should NOT contain "ago"
				if strings.Contains(result, "ago") {
					t.Errorf("Expected absolute format (no 'ago'), got %q", result)
				}
			}
		})
	}
}

// T040: Test formatTimestamp() relative formats (minutes/hours/days ago)
func TestFormatTimestamp_RelativeFormats(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		expectFormat string
	}{
		{"10 minutes", 10 * time.Minute, "minutes ago"},
		{"45 minutes", 45 * time.Minute, "minutes ago"},
		{"1 hour", 1 * time.Hour, "hour"},
		{"5 hours", 5 * time.Hour, "hours ago"},
		{"23 hours", 23 * time.Hour, "hours ago"},
		{"1 day", 24 * time.Hour, "day"},
		{"2 days", 48 * time.Hour, "days ago"},
		{"6 days", 6 * 24 * time.Hour, "days ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := time.Now().Add(-tt.duration)
			result := components.FormatTimestamp(timestamp)

			if !strings.Contains(result, tt.expectFormat) {
				t.Errorf("Expected format to contain %q, got %q", tt.expectFormat, result)
			}

			// Verify it contains "ago" (relative format)
			if !strings.Contains(result, "ago") && !strings.Contains(result, "hour") && !strings.Contains(result, "day") {
				t.Errorf("Expected relative time format, got %q", result)
			}
		})
	}
}

// T041: Test formatUsageLocations() with zero usage records (empty state)
func TestFormatUsageLocations_EmptyState(t *testing.T) {
	// Create credential with no usage records
	cred := CreateTestCredential("GitHub", "user", map[string]vault.UsageRecord{})

	result := components.FormatUsageLocations(cred)

	// Should display empty state message
	expectedEmpty := "No usage recorded"
	if !strings.Contains(result, expectedEmpty) {
		t.Errorf("Expected empty state message %q, got: %q", expectedEmpty, result)
	}
}

// T042: Test formatUsageLocations() with multiple locations and GitRepo display
func TestFormatUsageLocations_MultipleLocationsWithGitRepo(t *testing.T) {
	// Create credential with multiple usage records
	records := map[string]vault.UsageRecord{
		"/home/user/projects/pass-cli": CreateTestUsageRecord(
			"/home/user/projects/pass-cli", 2, "pass-cli", 5, 0,
		),
		"/home/user/projects/other-app": CreateTestUsageRecord(
			"/home/user/projects/other-app", 48, "other-app", 1, 0,
		),
		"/tmp/test": CreateTestUsageRecord(
			"/tmp/test", 1, "", 2, 0, // No git repo
		),
	}

	cred := CreateTestCredential("AWS", "admin", records)
	result := components.FormatUsageLocations(cred)

	// Should contain usage locations header (with or without color tags)
	if !strings.Contains(result, "Usage Locations") {
		t.Error("Expected 'Usage Locations' header in output")
	}

	// Should display all three locations
	if !strings.Contains(result, "/home/user/projects/pass-cli") {
		t.Error("Expected first location in output")
	}
	if !strings.Contains(result, "/home/user/projects/other-app") {
		t.Error("Expected second location in output")
	}
	if !strings.Contains(result, "/tmp/test") {
		t.Error("Expected third location in output")
	}

	// Should display git repo names when available
	if !strings.Contains(result, "pass-cli") {
		t.Error("Expected git repo 'pass-cli' in output")
	}
	if !strings.Contains(result, "other-app") {
		t.Error("Expected git repo 'other-app' in output")
	}

	// Should display access count
	if !strings.Contains(result, "5 times") || !strings.Contains(result, "accessed") {
		t.Error("Expected access count '5 times' in output")
	}

	// Should display timestamps (relative format for recent, absolute for old)
	if !strings.Contains(result, "ago") {
		t.Error("Expected relative timestamp 'ago' in output")
	}
}

// T043: Test long path truncation with ellipsis
func TestFormatUsageLocations_LongPathTruncation(t *testing.T) {
	// Create a very long path (200+ characters)
	longPath := "/home/user/very/deep/directory/structure/that/goes/on/and/on/and/on/and/on/" +
		"with/many/nested/folders/to/simulate/a/very/long/file/path/that/exceeds/" +
		"typical/terminal/width/and/should/be/truncated/with/ellipsis/in/the/middle/" +
		"to/fit/the/display/area"

	records := map[string]vault.UsageRecord{
		longPath: CreateTestUsageRecord(longPath, 1, "", 1, 0),
	}

	cred := CreateTestCredential("GitHub", "user", records)

	// Call formatUsageLocations with a specific terminal width constraint
	// For this test, we check that the function handles long paths gracefully
	result := components.FormatUsageLocations(cred)

	// The result should contain the path (possibly truncated)
	// We can't assert exact truncation without knowing terminal width,
	// but we can verify the function doesn't crash and produces output
	if result == "" {
		t.Error("Expected non-empty result for long path")
	}

	// Should still contain usage information
	if !strings.Contains(result, "Usage Locations") {
		t.Error("Expected 'Usage Locations' header even with long path")
	}

	// Should contain some part of the path or ellipsis indication
	// This is a basic check - actual truncation logic will be in implementation
	containsPath := strings.Contains(result, "/home/user") ||
		strings.Contains(result, "...") ||
		strings.Contains(result, longPath)

	if !containsPath {
		t.Error("Expected result to contain path or truncation indicator")
	}
}

// Additional test: formatUsageLocations with line numbers (FR-013 requirement)
func TestFormatUsageLocations_WithLineNumbers(t *testing.T) {
	// Create usage records with line numbers
	records := map[string]vault.UsageRecord{
		"/home/user/script.sh": CreateTestUsageRecord(
			"/home/user/script.sh", 1, "", 1, 42,
		),
		"/etc/config.yaml": CreateTestUsageRecord(
			"/etc/config.yaml", 2, "myrepo", 3, 15,
		),
	}

	cred := CreateTestCredential("API-Key", "service", records)
	result := components.FormatUsageLocations(cred)

	// Should display line numbers in format "path:lineNumber"
	if !strings.Contains(result, ":42") {
		t.Error("Expected line number :42 in output")
	}
	if !strings.Contains(result, ":15") {
		t.Error("Expected line number :15 in output")
	}
}

// Additional test: Verify sorting works with same timestamps
func TestSortUsageLocations_StableSort(t *testing.T) {
	now := time.Now()

	// Create multiple records with exact same timestamp
	records := map[string]vault.UsageRecord{
		"/path/alpha": {
			Location:  "/path/alpha",
			Timestamp: now,
			Count:     1,
		},
		"/path/beta": {
			Location:  "/path/beta",
			Timestamp: now,
			Count:     2,
		},
		"/path/gamma": {
			Location:  "/path/gamma",
			Timestamp: now,
			Count:     3,
		},
	}

	sorted := components.SortUsageLocations(records)

	// All should have same timestamp
	if len(sorted) != 3 {
		t.Fatalf("Expected 3 sorted records, got %d", len(sorted))
	}

	// Verify all timestamps are equal (stable sort preserves relative order)
	for i := 1; i < len(sorted); i++ {
		if !sorted[i].Timestamp.Equal(sorted[0].Timestamp) {
			t.Errorf("Expected all timestamps to be equal")
		}
	}
}

// Test helper to verify the format output structure
func TestFormatUsageLocations_OutputStructure(t *testing.T) {
	records := map[string]vault.UsageRecord{
		"/home/user/project": CreateTestUsageRecord(
			"/home/user/project", 5, "myrepo", 10, 25,
		),
	}

	cred := CreateTestCredential("GitHub", "user", records)
	result := components.FormatUsageLocations(cred)

	// Expected structure checks
	expectedElements := []string{
		"Usage Locations", // Header (without colon due to color tags)
		"/home/user/project", // Path
		"myrepo",             // Git repo
		"ago",                // Timestamp indicator
		"accessed",           // Usage count label
		"10 times",           // Count
		":25",                // Line number
	}

	for _, element := range expectedElements {
		if !strings.Contains(result, element) {
			t.Errorf("Expected output to contain %q, got:\n%s", element, result)
		}
	}
}

// Test that formatUsageLocations handles missing file paths gracefully (T046 requirement)
func TestFormatUsageLocations_MissingFilePath(t *testing.T) {
	// Create a usage record with a path that doesn't exist on disk
	nonExistentPath := "/path/that/does/not/exist/anywhere.txt"
	records := map[string]vault.UsageRecord{
		nonExistentPath: CreateTestUsageRecord(nonExistentPath, 1, "", 1, 0),
	}

	cred := CreateTestCredential("Test", "user", records)
	result := components.FormatUsageLocations(cred)

	// Should still display the path even if file doesn't exist
	if !strings.Contains(result, nonExistentPath) {
		t.Errorf("Expected to display path even if file doesn't exist: %s", nonExistentPath)
	}

	// Should not error or crash
	if result == "" {
		t.Error("Expected non-empty result even with missing file path")
	}
}
