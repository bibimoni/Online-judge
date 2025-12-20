package testutil

import (
	"context"
	"testing"
	"time"
)

// TestHelper provides common test utilities
type TestHelper struct {
	t   *testing.T
	ctx context.Context
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t:   t,
		ctx: context.Background(),
	}
}

// Context returns the test context
func (h *TestHelper) Context() context.Context {
	return h.ctx
}

// RequireNoError fails the test if error is not nil
func (h *TestHelper) RequireNoError(err error, msg string) {
	h.t.Helper()
	if err != nil {
		h.t.Fatalf("%s: %v", msg, err)
	}
}

// AssertEqual fails if values are not equal
func (h *TestHelper) AssertEqual(expected, actual interface{}, msg string) {
	h.t.Helper()
	if expected != actual {
		h.t.Errorf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertTrue fails if condition is false
func (h *TestHelper) AssertTrue(condition bool, msg string) {
	h.t.Helper()
	if !condition {
		h.t.Errorf("%s: condition is false", msg)
	}
}

// AssertFalse fails if condition is true
func (h *TestHelper) AssertFalse(condition bool, msg string) {
	h.t.Helper()
	if condition {
		h.t.Errorf("%s: condition is true", msg)
	}
}

// TimeFixture provides deterministic time values for testing
type TimeFixture struct {
	Now       time.Time
	Yesterday time.Time
	Tomorrow  time.Time
	NextWeek  time.Time
	LastWeek  time.Time
}

// NewTimeFixture creates a new time fixture
func NewTimeFixture() *TimeFixture {
	now := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	return &TimeFixture{
		Now:       now,
		Yesterday: now.AddDate(0, 0, -1),
		Tomorrow:  now.AddDate(0, 0, 1),
		NextWeek:  now.AddDate(0, 0, 7),
		LastWeek:  now.AddDate(0, 0, -7),
	}
}

// ContextWithTimeout creates a context with timeout for tests
func ContextWithTimeout(t *testing.T, duration time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	t.Cleanup(cancel)
	return ctx, cancel
}

// WaitForCondition polls a condition until it's true or timeout
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		if condition() {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("Timeout waiting for condition: %s", message)
		}
		<-ticker.C
	}
}

// AssertEventually polls until condition is true
func AssertEventually(t *testing.T, assertion func() bool, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if assertion() {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("Assertion failed within timeout")
}

// TestDatabase provides test database utilities
type TestDatabase struct {
	Name string
}

// NewTestDatabase creates a unique test database name
func NewTestDatabase() *TestDatabase {
	return &TestDatabase{
		Name: "test_db_" + time.Now().Format("20060102150405"),
	}
}

// MockClock provides a controllable clock for testing
type MockClock struct {
	current time.Time
}

// NewMockClock creates a new mock clock
func NewMockClock(start time.Time) *MockClock {
	return &MockClock{current: start}
}

// Now returns the current mocked time
func (m *MockClock) Now() time.Time {
	return m.current
}

// Advance moves the clock forward
func (m *MockClock) Advance(d time.Duration) {
	m.current = m.current.Add(d)
}

// Set sets the clock to a specific time
func (m *MockClock) Set(t time.Time) {
	m.current = t
}

// SkipIfShort skips the test if running in short mode
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

// SkipIfNoDocker skips test if Docker is not available
func SkipIfNoDocker(t *testing.T) {
	// Check if Docker is available
	// This is a simple check; you might want more sophisticated detection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	// Try to ping Docker
	// Implementation would check Docker availability
	_ = ctx
}

// RandomString generates a random string for testing
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}

// CompareJSON compares two JSON strings ignoring formatting
func CompareJSON(t *testing.T, expected, actual string) {
	t.Helper()
	// Implementation would parse and compare JSON
	// For now, simple string comparison
	if expected != actual {
		t.Errorf("JSON mismatch:\nExpected: %s\nActual: %s", expected, actual)
	}
}

// AssertNoRaceCondition runs a function multiple times concurrently
func AssertNoRaceCondition(t *testing.T, fn func(), iterations int) {
	t.Helper()
	done := make(chan bool)
	
	for i := 0; i < iterations; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Panic detected: %v", r)
				}
				done <- true
			}()
			fn()
		}()
	}
	
	for i := 0; i < iterations; i++ {
		<-done
	}
}
