package layout

import (
	"testing"

	"github.com/rivo/tview"
)

// TestNewPageManager verifies PageManager initialization
func TestNewPageManager(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	if pm == nil {
		t.Fatal("NewPageManager returned nil")
	}

	if pm.Pages == nil {
		t.Error("Pages not initialized")
	}

	if pm.app != app {
		t.Error("Application reference not set")
	}

	if pm.modalStack == nil {
		t.Error("Modal stack not initialized")
	}

	if len(pm.modalStack) != 0 {
		t.Errorf("Expected empty modal stack, got length %d", len(pm.modalStack))
	}
}

// TestShowPage verifies adding a non-modal page
func TestShowPage(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	mainLayout := tview.NewFlex()
	pm.ShowPage("main", mainLayout)

	// Verify page was added
	if pm.GetPageCount() != 1 {
		t.Errorf("Expected 1 page, got %d", pm.GetPageCount())
	}

	// Modal stack should still be empty
	if len(pm.modalStack) != 0 {
		t.Errorf("Expected empty modal stack, got length %d", len(pm.modalStack))
	}
}

// TestShowModal verifies modal display and stack management
func TestShowModal(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	// Add main page first
	pm.ShowPage("main", tview.NewFlex())

	// Show a modal
	modal := tview.NewTextView()
	pm.ShowModal("test-modal", modal, 40, 10)

	// Verify modal was added to stack
	if len(pm.modalStack) != 1 {
		t.Fatalf("Expected modal stack length 1, got %d", len(pm.modalStack))
	}

	if pm.modalStack[0] != "test-modal" {
		t.Errorf("Expected modal name 'test-modal', got '%s'", pm.modalStack[0])
	}

	// Verify HasModals returns true
	if !pm.HasModals() {
		t.Error("HasModals should return true when modal is shown")
	}
}

// TestShowMultipleModals verifies modal stacking
func TestShowMultipleModals(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())

	// Show multiple modals
	pm.ShowModal("modal1", tview.NewTextView(), 40, 10)
	pm.ShowModal("modal2", tview.NewTextView(), 40, 10)
	pm.ShowModal("modal3", tview.NewTextView(), 40, 10)

	// Verify stack order
	if len(pm.modalStack) != 3 {
		t.Fatalf("Expected modal stack length 3, got %d", len(pm.modalStack))
	}

	expectedOrder := []string{"modal1", "modal2", "modal3"}
	for i, expected := range expectedOrder {
		if pm.modalStack[i] != expected {
			t.Errorf("Stack position %d: expected '%s', got '%s'", i, expected, pm.modalStack[i])
		}
	}
}

// TestCloseModal verifies modal removal
func TestCloseModal(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())
	pm.ShowModal("modal1", tview.NewTextView(), 40, 10)
	pm.ShowModal("modal2", tview.NewTextView(), 40, 10)

	// Close middle modal
	pm.CloseModal("modal1")

	// Verify stack
	if len(pm.modalStack) != 1 {
		t.Fatalf("Expected modal stack length 1, got %d", len(pm.modalStack))
	}

	if pm.modalStack[0] != "modal2" {
		t.Errorf("Expected remaining modal 'modal2', got '%s'", pm.modalStack[0])
	}
}

// TestCloseTopModal verifies closing most recent modal
func TestCloseTopModal(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())
	pm.ShowModal("modal1", tview.NewTextView(), 40, 10)
	pm.ShowModal("modal2", tview.NewTextView(), 40, 10)
	pm.ShowModal("modal3", tview.NewTextView(), 40, 10)

	// Close top modal
	pm.CloseTopModal()

	// Verify modal3 was removed
	if len(pm.modalStack) != 2 {
		t.Fatalf("Expected modal stack length 2, got %d", len(pm.modalStack))
	}

	if pm.modalStack[1] != "modal2" {
		t.Errorf("Expected top modal 'modal2', got '%s'", pm.modalStack[1])
	}
}

// TestCloseTopModalWhenEmpty verifies safe no-op behavior
func TestCloseTopModalWhenEmpty(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())

	// This should not panic
	pm.CloseTopModal()

	if len(pm.modalStack) != 0 {
		t.Errorf("Expected empty modal stack, got length %d", len(pm.modalStack))
	}
}

// TestCloseNonExistentModal verifies safe no-op behavior
func TestCloseNonExistentModal(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())
	pm.ShowModal("modal1", tview.NewTextView(), 40, 10)

	// This should not panic
	pm.CloseModal("non-existent")

	// Stack should be unchanged
	if len(pm.modalStack) != 1 {
		t.Errorf("Expected modal stack length 1, got %d", len(pm.modalStack))
	}
}

// TestHasModals verifies modal detection
func TestHasModals(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())

	// Initially no modals
	if pm.HasModals() {
		t.Error("HasModals should return false when no modals are shown")
	}

	// Show modal
	pm.ShowModal("test", tview.NewTextView(), 40, 10)

	if !pm.HasModals() {
		t.Error("HasModals should return true when modal is shown")
	}

	// Close modal
	pm.CloseTopModal()

	if pm.HasModals() {
		t.Error("HasModals should return false after closing all modals")
	}
}

// TestShowForm verifies form modal helper
func TestShowForm(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())

	form := tview.NewForm()
	pm.ShowForm(form, "Test Form")

	// Verify modal was added
	if len(pm.modalStack) != 1 {
		t.Fatalf("Expected modal stack length 1, got %d", len(pm.modalStack))
	}

	if pm.modalStack[0] != "form" {
		t.Errorf("Expected modal name 'form', got '%s'", pm.modalStack[0])
	}
}

// TestShowConfirmDialog verifies confirmation dialog creation
func TestShowConfirmDialog(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())

	// Note: Callbacks are tested for creation, not invocation
	// (simulating button press in tview is complex)
	pm.ShowConfirmDialog(
		"Test Title",
		"Test message",
		func() { /* onYes callback */ },
		func() { /* onNo callback */ },
	)

	// Verify modal was added
	if len(pm.modalStack) != 1 {
		t.Fatalf("Expected modal stack length 1, got %d", len(pm.modalStack))
	}

	if pm.modalStack[0] != "confirm" {
		t.Errorf("Expected modal name 'confirm', got '%s'", pm.modalStack[0])
	}

	// Note: Testing callback invocation would require simulating button press,
	// which is complex with tview. The important part is verifying the modal
	// is created and added to the stack correctly.
}

// TestSwitchToPage verifies page switching
func TestSwitchToPage(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("page1", tview.NewFlex())
	pm.ShowPage("page2", tview.NewFlex())

	// Switch to page2
	pm.SwitchToPage("page2")

	// Verify no modals were added
	if len(pm.modalStack) != 0 {
		t.Errorf("Expected empty modal stack, got length %d", len(pm.modalStack))
	}
}

// TestCloseModalStackManagement verifies proper stack cleanup
func TestCloseModalStackManagement(t *testing.T) {
	app := tview.NewApplication()
	pm := NewPageManager(app)

	pm.ShowPage("main", tview.NewFlex())

	// Show and close multiple modals
	pm.ShowModal("modal1", tview.NewTextView(), 40, 10)
	pm.ShowModal("modal2", tview.NewTextView(), 40, 10)
	pm.ShowModal("modal3", tview.NewTextView(), 40, 10)

	// Close in different order
	pm.CloseModal("modal2") // Close middle
	pm.CloseTopModal()      // Close top (modal3)

	// Should have only modal1 left
	if len(pm.modalStack) != 1 {
		t.Fatalf("Expected modal stack length 1, got %d", len(pm.modalStack))
	}

	if pm.modalStack[0] != "modal1" {
		t.Errorf("Expected remaining modal 'modal1', got '%s'", pm.modalStack[0])
	}

	// Close last modal
	pm.CloseTopModal()

	// Stack should be empty
	if len(pm.modalStack) != 0 {
		t.Errorf("Expected empty modal stack, got length %d", len(pm.modalStack))
	}

	if pm.HasModals() {
		t.Error("HasModals should return false after closing all modals")
	}
}

// =============================================================================
// User Story 1 Tests: Terminal Size Warning Display
// =============================================================================

// Note: PageManager tests using ShowSizeWarning/HideSizeWarning call app.Draw()
// which can block in test environment. The core functionality is tested via
// LayoutManager tests (TestHandleResize_BelowMinimum, TestHandleResize_StartupCheck)
// which use mocks. These tests verify basic integration without blocking.

// TestSizeWarningStateTracking verifies IsSizeWarningActive
// returns the correct boolean state.
func TestSizeWarningStateTracking(t *testing.T) {
	// Create PageManager with real app (but don't call methods that trigger Draw)
	app := tview.NewApplication()
	pm := NewPageManager(app)

	// Initially not active
	if pm.IsSizeWarningActive() {
		t.Error("IsSizeWarningActive should return false initially")
	}

	// Manually set state to simulate warning being shown
	pm.sizeWarningActive = true

	// Should be active
	if !pm.IsSizeWarningActive() {
		t.Error("IsSizeWarningActive should return true when state flag is true")
	}

	// Manually clear state
	pm.sizeWarningActive = false

	// Should be inactive again
	if pm.IsSizeWarningActive() {
		t.Error("IsSizeWarningActive should return false when state flag is false")
	}
}
