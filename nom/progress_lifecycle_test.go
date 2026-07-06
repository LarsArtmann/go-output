package nom

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestMultiSubscriber_ProgressAndRetryFanout verifies that ActivityProgress and
// ActivityRetrying events propagate to every NOMSubscriber behind a
// MultiSubscriber. Each subscriber independently tracks the progress message and
// retry count.
func TestMultiSubscriber_ProgressAndRetryFanout(t *testing.T) {
	t.Parallel()

	sub1 := NewNOMSubscriber()
	sub2 := NewNOMSubscriber()
	multi := NewMultiSubscriber(sub1, sub2)

	ctx := context.Background()

	// Order matters: retry clears progress (SetRunning), so send progress last.
	events := []Event{
		WorkflowStarted{ID: "wf", Name: "test"},
		ActivityStarted{ID: "step", Name: "build"},
		ActivityFailed{ID: "step", Name: "build"},
		ActivityRetrying{ID: "step", Name: "build", Attempt: 1, Reason: "flaky"},
		ActivityProgress{ID: "step", Name: "build", Message: "Working [1/3]"},
	}

	for _, e := range events {
		if err := multi.OnEvent(ctx, e); err != nil {
			t.Fatalf("OnEvent(%T) error: %v", e, err)
		}
	}

	for _, sub := range []*NOMSubscriber{sub1, sub2} {
		snap := sub.SnapshotActivities()["step"]
		if snap.Progress != "Working [1/3]" {
			t.Errorf("subscriber progress = %q, want %q", snap.Progress, "Working [1/3]")
		}

		if snap.RetryCount != 1 {
			t.Errorf("subscriber retry count = %d, want 1", snap.RetryCount)
		}

		if snap.RetryReason != "flaky" {
			t.Errorf("subscriber retry reason = %q, want %q", snap.RetryReason, "flaky")
		}
	}
}

// TestResetClearsProgressAndRetry verifies that Reset() wipes all per-activity
// state including the new Progress, RetryCount, and RetryReason fields.
func TestResetClearsProgressAndRetry(t *testing.T) {
	t.Parallel()

	ns := NewNOMSubscriber()
	ctx := context.Background()

	// Retry clears progress, so send progress after the retry to populate both.
	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step", Name: "build"})
	_ = ns.OnEvent(ctx, ActivityFailed{ID: "step", Name: "build"})
	_ = ns.OnEvent(ctx, ActivityRetrying{ID: "step", Name: "build", Attempt: 2, Reason: "timeout"})
	_ = ns.OnEvent(ctx, ActivityProgress{ID: "step", Name: "build", Message: "Working"})

	// Preconditions: state is populated.
	snap := ns.SnapshotActivities()["step"]
	if snap.RetryCount != 2 || snap.Progress != "Working" {
		t.Fatalf("precondition: retry=%d progress=%q", snap.RetryCount, snap.Progress)
	}

	ns.Reset()

	// After reset: no activities, counts zeroed.
	if got := ns.GetActivityCounts(); got.Total() != 0 {
		t.Errorf("counts after reset = %+v, want empty", got)
	}

	if ns.GetActivity("step") != nil {
		t.Error("GetActivity should return nil after Reset")
	}
}

// TestProgressClearedOnRetry verifies that when an activity is retried, the
// progress message from the prior attempt is cleared (SetRunning clears Progress).
// A retry starts fresh with no stale sub-step message.
func TestProgressClearedOnRetry(t *testing.T) {
	t.Parallel()

	ns := NewNOMSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step", Name: "build"})
	_ = ns.OnEvent(ctx, ActivityProgress{ID: "step", Name: "build", Message: "Attempt 1 progress"})

	// Verify progress is set.
	snap := ns.SnapshotActivities()["step"]
	if snap.Progress != "Attempt 1 progress" {
		t.Fatalf("precondition progress = %q", snap.Progress)
	}

	// Fail and retry.
	_ = ns.OnEvent(ctx, ActivityFailed{ID: "step", Name: "build"})
	_ = ns.OnEvent(ctx, ActivityRetrying{ID: "step", Name: "build", Attempt: 1})

	// Progress should be cleared by the retry's SetRunning.
	snap = ns.SnapshotActivities()["step"]
	if snap.Progress != "" {
		t.Errorf("progress after retry = %q, want empty (retry starts fresh)", snap.Progress)
	}

	if snap.RetryCount != 1 {
		t.Errorf("retry count = %d, want 1", snap.RetryCount)
	}
}

// TestConcurrentProgressAndRetry exercises the event handlers under concurrent
// load to surface data races between progress updates and retry transitions.
// Run with -race to verify lock correctness.
func TestConcurrentProgressAndRetry(t *testing.T) {
	t.Parallel()

	ns := NewNOMSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step", Name: "build"})

	const workers = 8

	var wg sync.WaitGroup
	wg.Add(workers)

	for range workers {
		go func() {
			defer wg.Done()

			for i := range 100 {
				switch i % 3 {
				case 0:
					_ = ns.OnEvent(ctx, ActivityProgress{
						ID: "step", Name: "build", Message: "working",
					})
				case 1:
					_ = ns.OnEvent(ctx, ActivityFailed{ID: "step", Name: "build"})
				case 2:
					_ = ns.OnEvent(ctx, ActivityRetrying{ID: "step", Name: "build", Attempt: 1})
				}

				// Read state concurrently to catch snapshot races.
				_ = ns.SnapshotActivities()
				_ = ns.GetActivityCounts()
			}
		}()
	}

	wg.Wait()

	// After all goroutines finish, the subscriber must be in a consistent state:
	// counts must match a brute-force recount of activities.
	snapshots := ns.SnapshotActivities()
	counts := ns.GetActivityCounts()

	var expected ActivityCounts

	for _, s := range snapshots {
		switch s.Status {
		case ActivityStatusRunning:
			expected.Running++
		case ActivityStatusPending:
			expected.Pending++
		case ActivityStatusCompleted:
			expected.Completed++
		case ActivityStatusFailed:
			expected.Failed++
		}
	}

	if counts != expected {
		t.Errorf("counts cache drifted: cache=%+v expected=%+v", counts, expected)
	}
}

// TestEstimatedTotalRemainingRunningElapsed verifies that running activities
// contribute only their remaining estimate (estimated - elapsed), not the full
// estimate. This guards the elapsed-subtraction logic in EstimatedTotalRemaining.
func TestEstimatedTotalRemainingRunningElapsed(t *testing.T) {
	t.Parallel()

	ns := NewNOMSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "run", Name: "run"})

	ns.SetEstimatedTime("run", 1*time.Hour)

	// A running activity's remaining estimate must be less than its full estimate
	// (elapsed > 0) but non-negative.
	remaining := ns.EstimatedTotalRemaining()
	if remaining <= 0 || remaining >= 1*time.Hour {
		t.Errorf("running remaining = %v, want (0, 1h)", remaining)
	}
}
