package security

import (
	"testing"
	"time"
)

// T051a: Test rate limiting after 3 failures
func TestValidationRateLimiter_ThreeFailures(t *testing.T) {
	rl := NewValidationRateLimiter()

	// First failure - should pass
	err := rl.CheckAndRecordFailure()
	if err != nil {
		t.Errorf("First failure should not trigger cooldown, got: %v", err)
	}

	// Second failure - should pass
	err = rl.CheckAndRecordFailure()
	if err != nil {
		t.Errorf("Second failure should not trigger cooldown, got: %v", err)
	}

	// Third failure - should trigger cooldown
	err = rl.CheckAndRecordFailure()
	if err == nil {
		t.Error("Third failure should trigger cooldown")
	}
	if err != nil && err.Error() != "too many failed attempts - please wait 5 seconds before trying again" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

// T051a: Test cooldown enforcement
func TestValidationRateLimiter_CooldownEnforcement(t *testing.T) {
	rl := NewValidationRateLimiter()

	// Trigger cooldown
	rl.CheckAndRecordFailure()
	rl.CheckAndRecordFailure()
	rl.CheckAndRecordFailure()

	// Attempt during cooldown - should fail
	err := rl.CheckAndRecordFailure()
	if err == nil {
		t.Error("Attempt during cooldown should fail")
	}

	// Wait for cooldown to expire
	time.Sleep(6 * time.Second)

	// Attempt after cooldown - should pass
	err = rl.CheckAndRecordFailure()
	if err != nil {
		t.Errorf("Attempt after cooldown should pass, got: %v", err)
	}
}

// T051a: Test reset clears state
func TestValidationRateLimiter_Reset(t *testing.T) {
	rl := NewValidationRateLimiter()

	// Record two failures
	rl.CheckAndRecordFailure()
	rl.CheckAndRecordFailure()

	// Reset
	rl.Reset()

	// Next three failures should behave like first-time failures
	err := rl.CheckAndRecordFailure()
	if err != nil {
		t.Errorf("First failure after reset should not trigger cooldown, got: %v", err)
	}

	err = rl.CheckAndRecordFailure()
	if err != nil {
		t.Errorf("Second failure after reset should not trigger cooldown, got: %v", err)
	}

	err = rl.CheckAndRecordFailure()
	if err == nil {
		t.Error("Third failure after reset should trigger cooldown")
	}
}

// T051a: Test auto-reset after 30 seconds of inactivity
func TestValidationRateLimiter_AutoReset(t *testing.T) {
	rl := NewValidationRateLimiter()

	// Record two failures
	rl.CheckAndRecordFailure()
	rl.CheckAndRecordFailure()

	// Manually advance last failure time (simulate 31 seconds passing)
	rl.mu.Lock()
	rl.lastFailure = time.Now().Add(-31 * time.Second)
	rl.mu.Unlock()

	// Next failure should not count previous two (auto-reset)
	err := rl.CheckAndRecordFailure()
	if err != nil {
		t.Errorf("First failure after auto-reset should not trigger cooldown, got: %v", err)
	}

	// Verify counter was reset
	rl.mu.Lock()
	count := rl.failureCount
	rl.mu.Unlock()

	if count != 1 {
		t.Errorf("Failure count after auto-reset should be 1, got: %d", count)
	}
}

// T051a: Test concurrent access safety
func TestValidationRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := NewValidationRateLimiter()

	// Spawn 10 goroutines trying to record failures concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			rl.CheckAndRecordFailure()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// No panic = success (mutex protects concurrent access)
}

// T051a: Test cooldown remaining time message
func TestValidationRateLimiter_RemainingTimeMessage(t *testing.T) {
	rl := NewValidationRateLimiter()

	// Trigger cooldown
	rl.CheckAndRecordFailure()
	rl.CheckAndRecordFailure()
	rl.CheckAndRecordFailure()

	// Wait 1 second
	time.Sleep(1 * time.Second)

	// Check error message includes remaining time
	err := rl.CheckAndRecordFailure()
	if err == nil {
		t.Error("Expected cooldown error")
	}

	// Error should mention "wait" and remaining time (approximately 4s)
	errMsg := err.Error()
	if errMsg != "too many failed attempts - please wait 4s before trying again" &&
		errMsg != "too many failed attempts - please wait 3s before trying again" {
		t.Errorf("Unexpected cooldown message: %v", errMsg)
	}
}
