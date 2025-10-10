package tui_test

import (
	"testing"
	"time"

	"pass-cli/internal/vault"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// MockLayoutManager provides a test double for LayoutManager
type MockLayoutManager struct {
	SidebarOverride     *bool
	DetailPanelOverride *bool
	Width               int
	Height              int
	RebuildCalled       bool
}

// NewMockLayoutManager creates a new mock layout manager
func NewMockLayoutManager() *MockLayoutManager {
	return &MockLayoutManager{
		Width:  120,
		Height: 40,
	}
}

// MockRebuildLayout simulates layout rebuild
func (m *MockLayoutManager) MockRebuildLayout() {
	m.RebuildCalled = true
}

// CreateTestCredential creates a test credential with usage records
func CreateTestCredential(service, username string, usageRecords map[string]vault.UsageRecord) *vault.Credential {
	now := time.Now()
	return &vault.Credential{
		Service:     service,
		Username:    username,
		Password:    "test-password",
		Category:    "test-category",
		URL:         "https://example.com",
		Notes:       "test notes",
		CreatedAt:   now,
		UpdatedAt:   now,
		UsageRecord: usageRecords,
	}
}

// CreateTestUsageRecord creates a test usage record
func CreateTestUsageRecord(location string, hoursAgo int, gitRepo string, count int, lineNumber int) vault.UsageRecord {
	return vault.UsageRecord{
		Location:   location,
		Timestamp:  time.Now().Add(-time.Duration(hoursAgo) * time.Hour),
		GitRepo:    gitRepo,
		Count:      count,
		LineNumber: lineNumber,
	}
}

// CreateTestCredentialMetadata creates test credential metadata
func CreateTestCredentialMetadata(service, username, category, url string) *vault.CredentialMetadata {
	return &vault.CredentialMetadata{
		Service:  service,
		Username: username,
		Category: category,
		URL:      url,
	}
}

// SimulateApp creates a minimal tview.Application for testing
func SimulateApp(t *testing.T) *tview.Application {
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("Failed to initialize simulation screen: %v", err)
	}
	app.SetScreen(screen)
	return app
}

// AssertTrue fails the test if condition is false
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Errorf("Assertion failed: %s", message)
	}
}

// AssertFalse fails the test if condition is true
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Errorf("Assertion failed: %s", message)
	}
}

// AssertEqual fails the test if expected != actual
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Errorf("Assertion failed: %s\nExpected: %v\nActual: %v", message, expected, actual)
	}
}

// AssertNil fails the test if value is not nil
func AssertNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value != nil {
		t.Errorf("Assertion failed: %s\nExpected nil, got: %v", message, value)
	}
}

// AssertNotNil fails the test if value is nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Errorf("Assertion failed: %s\nExpected non-nil value", message)
	}
}
