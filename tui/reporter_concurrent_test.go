package tui

import (
	"sync"
	"testing"
)

// TestBubbleTeaProgressReporter_ConcurrentAccess verifies that concurrent
// Report* calls from multiple goroutines do not trigger data races.
//
// With the old design (direct model field mutation from caller goroutine),
// this test would fail under -race because the TUI goroutine's Update()
// and the caller goroutine's Report* methods both mutated model fields
// (workflowState, currentProgress, currentMessage, steps) without
// synchronization.
//
// The fix routes ALL model mutations through send() → model.Update(),
// serialized by pr.mu. The reporter owns its own workflowState (protected
// by the same mutex) for decision logic.
func TestBubbleTeaProgressReporter_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	reporter := newTestReporter()

	var wg sync.WaitGroup

	for range 50 {
		wg.Add(4)

		go func() {
			defer wg.Done()

			reporter.ReportProgress(50.0)
		}()

		go func() {
			defer wg.Done()

			reporter.ReportMessage("working")
		}()

		go func() {
			defer wg.Done()

			reporter.ReportStep(1, 2, "step")
		}()

		go func() {
			defer wg.Done()

			reporter.ReportError(errTestFail)
		}()
	}

	wg.Wait()

	// After concurrent reporting, the reporter should be in a terminal state
	// (Errored from ReportError calls, or Completed from ReportProgress(100)).
	// The exact terminal state depends on goroutine scheduling, but it must
	// be either Running, Completed, or Errored — never an invalid state.
	pr := reporter
	pr.mu.RLock()
	state := pr.workflowState
	pr.mu.RUnlock()

	if state == workflowStateIdle {
		t.Error("workflow should have transitioned away from Idle after concurrent reports")
	}
}

// TestBubbleTeaProgressReporter_ConcurrentReportProgress verifies that
// concurrent ReportProgress calls from different goroutines race-free
// transition the reporter through Idle → Running → Completed.
func TestBubbleTeaProgressReporter_ConcurrentReportProgress(t *testing.T) {
	t.Parallel()

	reporter := newTestReporter()

	var wg sync.WaitGroup

	for i := range 100 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			// Every 10th call reports 100% to trigger completion.
			pct := float64(idx % 100)
			if idx%10 == 0 {
				pct = 100.0
			}

			reporter.ReportProgress(pct)
		}(i)
	}

	wg.Wait()

	// At least one goroutine called ReportProgress(100), so the reporter
	// should have transitioned to Completed.
	pr := reporter
	pr.mu.RLock()
	state := pr.workflowState
	pr.mu.RUnlock()

	if state != workflowStateCompleted {
		t.Errorf("workflow state = %v, want Completed after concurrent 100%% reports", state)
	}
}

// TestBubbleTeaProgressReporter_ConcurrentStepUpdates verifies that
// concurrent ReportStep calls from different goroutines race-free
// create and update steps without data races on the steps slice.
func TestBubbleTeaProgressReporter_ConcurrentStepUpdates(t *testing.T) {
	t.Parallel()

	reporter := newTestReporter()

	var wg sync.WaitGroup

	messages := []string{"Compile", "Link", "Test", "Lint", "Package"}

	for i := range 100 {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			msg := messages[idx%len(messages)]
			reporter.ReportStep(uint(idx%5+1), 5, msg)
		}(i)
	}

	wg.Wait()

	// Verify steps slice is in a consistent state (no panic, no corruption).
	// We can't assert exact step count due to concurrent find-or-create logic,
	// but the slice must not be corrupted.
	pr := reporter
	pr.mu.RLock()
	steps := pr.model.steps
	pr.mu.RUnlock()

	if len(steps) == 0 {
		t.Error("expected at least one step after concurrent ReportStep calls")
	}
}
