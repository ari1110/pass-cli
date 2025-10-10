package tui_test

import (
	"fmt"
	"testing"
	"time"

	"pass-cli/cmd/tui/components"
	"pass-cli/cmd/tui/models"
	"pass-cli/internal/vault"
)

// T063: Performance validation for search filtering and detail rendering

// BenchmarkSearchFiltering_1000Credentials validates search performance meets <100ms requirement
func BenchmarkSearchFiltering_1000Credentials(b *testing.B) {
	// Setup: Create 1000 test credentials
	credentials := make([]vault.CredentialMetadata, 1000)
	for i := 0; i < 1000; i++ {
		credentials[i] = vault.CredentialMetadata{
			Service:  fmt.Sprintf("Service-%d", i),
			Username: fmt.Sprintf("user%d@example.com", i),
			Category: "work",
			URL:      fmt.Sprintf("https://service%d.com", i),
		}
	}

	// Add some matching credentials
	credentials[500].Service = "GitHub"
	credentials[750].Service = "GitLab"

	searchState := &models.SearchState{
		Active: true,
		Query:  "git",
	}

	// Benchmark the filtering operation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matchCount := 0
		for j := range credentials {
			if searchState.MatchesCredential(&credentials[j]) {
				matchCount++
			}
		}
	}
}

// TestSearchFiltering_Performance validates <100ms requirement for 1000 credentials
func TestSearchFiltering_Performance(t *testing.T) {
	// Create 1000 test credentials
	credentials := make([]vault.CredentialMetadata, 1000)
	for i := 0; i < 1000; i++ {
		credentials[i] = vault.CredentialMetadata{
			Service:  fmt.Sprintf("Service-%d", i),
			Username: fmt.Sprintf("user%d@example.com", i),
			Category: "work",
			URL:      fmt.Sprintf("https://service%d.com", i),
		}
	}

	// Add matching credentials
	credentials[500].Service = "GitHub"
	credentials[750].Service = "GitLab"

	searchState := &models.SearchState{
		Active: true,
		Query:  "git",
	}

	// Measure filtering time
	start := time.Now()
	matchCount := 0
	for i := range credentials {
		if searchState.MatchesCredential(&credentials[i]) {
			matchCount++
		}
	}
	elapsed := time.Since(start)

	t.Logf("Search filtering 1000 credentials took %v (found %d matches)", elapsed, matchCount)

	// Validate: Must be under 100ms
	if elapsed > 100*time.Millisecond {
		t.Errorf("Search filtering took %v, exceeds 100ms requirement", elapsed)
	}

	// Sanity check: Should find 2 matches
	if matchCount != 2 {
		t.Errorf("Expected 2 matches, got %d", matchCount)
	}
}

// BenchmarkDetailRendering validates detail panel rendering performance
func BenchmarkDetailRendering(b *testing.B) {
	// Setup: Create credential with usage records
	usageRecords := make(map[string]vault.UsageRecord)
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/home/user/project/file%d.go", i)
		usageRecords[path] = vault.UsageRecord{
			Location:   path,
			Timestamp:  time.Now().Add(-time.Duration(i) * time.Hour),
			GitRepo:    "pass-cli",
			Count:      i + 1,
			LineNumber: 42 + i,
		}
	}

	cred := &vault.Credential{
		Service:     "GitHub",
		Username:    "user@example.com",
		Category:    "work",
		URL:         "https://github.com",
		Password:    "secret123",
		CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:   time.Now().Add(-1 * 24 * time.Hour),
		UsageRecord: usageRecords,
	}

	// Benchmark the formatting operation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = components.FormatUsageLocations(cred)
	}
}

// TestDetailRendering_Performance validates <200ms requirement
func TestDetailRendering_Performance(t *testing.T) {
	// Create credential with typical number of usage records (10 locations)
	usageRecords := make(map[string]vault.UsageRecord)
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/home/user/project/file%d.go", i)
		usageRecords[path] = vault.UsageRecord{
			Location:   path,
			Timestamp:  time.Now().Add(-time.Duration(i) * time.Hour),
			GitRepo:    "pass-cli",
			Count:      i + 1,
			LineNumber: 42 + i,
		}
	}

	cred := &vault.Credential{
		Service:     "GitHub",
		Username:    "user@example.com",
		Category:    "work",
		URL:         "https://github.com",
		Password:    "secret123",
		CreatedAt:   time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:   time.Now().Add(-1 * 24 * time.Hour),
		UsageRecord: usageRecords,
	}

	// Measure rendering time
	start := time.Now()
	result := components.FormatUsageLocations(cred)
	elapsed := time.Since(start)

	t.Logf("Detail rendering with 10 usage locations took %v", elapsed)

	// Validate: Must be under 200ms
	if elapsed > 200*time.Millisecond {
		t.Errorf("Detail rendering took %v, exceeds 200ms requirement", elapsed)
	}

	// Sanity check: Should contain usage locations header
	if result == "" {
		t.Error("Expected non-empty rendering result")
	}
}

// TestSortUsageLocations_Performance validates sorting performance for large datasets
func TestSortUsageLocations_Performance(t *testing.T) {
	// Create 100 usage records (stress test)
	records := make(map[string]vault.UsageRecord)
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/path/to/file/%d.go", i)
		records[path] = vault.UsageRecord{
			Location:   path,
			Timestamp:  time.Now().Add(-time.Duration(i) * time.Minute),
			GitRepo:    "repo",
			Count:      1,
			LineNumber: i,
		}
	}

	// Measure sorting time
	start := time.Now()
	sorted := components.SortUsageLocations(records)
	elapsed := time.Since(start)

	t.Logf("Sorting 100 usage records took %v", elapsed)

	// Validate: Should be very fast (well under 10ms)
	if elapsed > 10*time.Millisecond {
		t.Errorf("Sorting 100 records took %v, expected <10ms", elapsed)
	}

	// Sanity check: Verify sort correctness
	if len(sorted) != 100 {
		t.Errorf("Expected 100 sorted records, got %d", len(sorted))
	}

	// Verify descending order (most recent first)
	for i := 0; i < len(sorted)-1; i++ {
		if sorted[i].Timestamp.Before(sorted[i+1].Timestamp) {
			t.Errorf("Sort order incorrect at index %d", i)
		}
	}
}
