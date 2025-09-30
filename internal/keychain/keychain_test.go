package keychain

import (
	"testing"

	"github.com/zalando/go-keyring"
)

func TestNew(t *testing.T) {
	ks := New()
	if ks == nil {
		t.Fatal("New() returned nil")
	}

	// Availability depends on the test environment
	// Just verify the field is set (true or false)
	t.Logf("Keychain available: %v", ks.IsAvailable())
}

func TestStoreAndRetrieve(t *testing.T) {
	ks := New()

	if !ks.IsAvailable() {
		t.Skip("Keychain not available in test environment")
	}

	// Clean up before test
	_ = ks.Delete()

	testPassword := "test-master-password-12345"

	// Test Store
	err := ks.Store(testPassword)
	if err != nil {
		t.Fatalf("Store() failed: %v", err)
	}

	// Test Retrieve
	retrieved, err := ks.Retrieve()
	if err != nil {
		t.Fatalf("Retrieve() failed: %v", err)
	}

	if retrieved != testPassword {
		t.Errorf("Retrieved password = %q, want %q", retrieved, testPassword)
	}

	// Clean up after test
	_ = ks.Delete()
}

func TestRetrieveNonExistent(t *testing.T) {
	ks := New()

	if !ks.IsAvailable() {
		t.Skip("Keychain not available in test environment")
	}

	// Ensure password doesn't exist
	_ = ks.Delete()

	// Try to retrieve non-existent password
	_, err := ks.Retrieve()
	if err == nil {
		t.Fatal("Retrieve() should fail for non-existent password")
	}

	if err != ErrPasswordNotFound {
		t.Errorf("Retrieve() error = %v, want %v", err, ErrPasswordNotFound)
	}
}

func TestDelete(t *testing.T) {
	ks := New()

	if !ks.IsAvailable() {
		t.Skip("Keychain not available in test environment")
	}

	// Store a password first
	testPassword := "test-password-to-delete"
	err := ks.Store(testPassword)
	if err != nil {
		t.Fatalf("Store() failed: %v", err)
	}

	// Delete it
	err = ks.Delete()
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Verify it's gone
	_, err = ks.Retrieve()
	if err != ErrPasswordNotFound {
		t.Errorf("After Delete(), Retrieve() error = %v, want %v", err, ErrPasswordNotFound)
	}
}

func TestDeleteNonExistent(t *testing.T) {
	ks := New()

	if !ks.IsAvailable() {
		t.Skip("Keychain not available in test environment")
	}

	// Ensure password doesn't exist
	_ = ks.Delete()

	// Delete should not error for non-existent password
	err := ks.Delete()
	if err != nil {
		t.Errorf("Delete() on non-existent password failed: %v", err)
	}
}

func TestClear(t *testing.T) {
	ks := New()

	if !ks.IsAvailable() {
		t.Skip("Keychain not available in test environment")
	}

	// Store a password
	testPassword := "test-password-to-clear"
	err := ks.Store(testPassword)
	if err != nil {
		t.Fatalf("Store() failed: %v", err)
	}

	// Clear it
	err = ks.Clear()
	if err != nil {
		t.Fatalf("Clear() failed: %v", err)
	}

	// Verify it's gone
	_, err = ks.Retrieve()
	if err != ErrPasswordNotFound {
		t.Errorf("After Clear(), Retrieve() error = %v, want %v", err, ErrPasswordNotFound)
	}
}

func TestUnavailableKeychain(t *testing.T) {
	// Create a service with forced unavailability
	ks := &KeychainService{available: false}

	// Test Store
	err := ks.Store("test")
	if err != ErrKeychainUnavailable {
		t.Errorf("Store() with unavailable keychain error = %v, want %v", err, ErrKeychainUnavailable)
	}

	// Test Retrieve
	_, err = ks.Retrieve()
	if err != ErrKeychainUnavailable {
		t.Errorf("Retrieve() with unavailable keychain error = %v, want %v", err, ErrKeychainUnavailable)
	}

	// Test Delete
	err = ks.Delete()
	if err != ErrKeychainUnavailable {
		t.Errorf("Delete() with unavailable keychain error = %v, want %v", err, ErrKeychainUnavailable)
	}

	// Test Clear
	err = ks.Clear()
	if err != ErrKeychainUnavailable {
		t.Errorf("Clear() with unavailable keychain error = %v, want %v", err, ErrKeychainUnavailable)
	}
}

func TestStoreEmptyPassword(t *testing.T) {
	ks := New()

	if !ks.IsAvailable() {
		t.Skip("Keychain not available in test environment")
	}

	// Clean up before test
	_ = ks.Delete()

	// Store empty password (should be allowed)
	err := ks.Store("")
	if err != nil {
		t.Fatalf("Store() with empty password failed: %v", err)
	}

	// Retrieve it
	retrieved, err := ks.Retrieve()
	if err != nil {
		t.Fatalf("Retrieve() failed: %v", err)
	}

	if retrieved != "" {
		t.Errorf("Retrieved password = %q, want empty string", retrieved)
	}

	// Clean up
	_ = ks.Delete()
}

func TestMultipleStoreOverwrites(t *testing.T) {
	ks := New()

	if !ks.IsAvailable() {
		t.Skip("Keychain not available in test environment")
	}

	// Clean up before test
	_ = ks.Delete()

	// Store first password
	password1 := "first-password"
	err := ks.Store(password1)
	if err != nil {
		t.Fatalf("First Store() failed: %v", err)
	}

	// Store second password (should overwrite)
	password2 := "second-password"
	err = ks.Store(password2)
	if err != nil {
		t.Fatalf("Second Store() failed: %v", err)
	}

	// Retrieve should get the second password
	retrieved, err := ks.Retrieve()
	if err != nil {
		t.Fatalf("Retrieve() failed: %v", err)
	}

	if retrieved != password2 {
		t.Errorf("Retrieved password = %q, want %q", retrieved, password2)
	}

	// Clean up
	_ = ks.Delete()
}

// TestCheckAvailability verifies the availability check works
func TestCheckAvailability(t *testing.T) {
	ks := New()

	// The availability check should have run during New()
	available := ks.IsAvailable()

	// Try a manual operation to verify consistency
	err := keyring.Set(ServiceName, "test-check", "test")
	if err != nil {
		if available {
			t.Error("IsAvailable() returned true but keyring.Set() failed")
		}
	} else {
		if !available {
			t.Error("IsAvailable() returned false but keyring.Set() succeeded")
		}
		_ = keyring.Delete(ServiceName, "test-check")
	}
}