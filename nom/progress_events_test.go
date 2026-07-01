package nom

import (
	"context"
	"strings"
	"testing"
	"time"
)

// TestActivityProgressEvent verifies that an ActivityProgress event sets the
// progress message on the activity, and that it renders as a dim sub-line.
func TestActivityProgressEvent(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "go-mod-tidy"})
	_ = ns.OnEvent(ctx, ActivityProgress{
		ID:      "step1",
		Name:    "go-mod-tidy",
		Message: "Tidying module [2/26]: modules/gitignore",
	})

	activity := ns.GetActivity("step1")
	if activity == nil {
		t.Fatal("activity not found")
	}

	if activity.Progress != "Tidying module [2/26]: modules/gitignore" {
		t.Errorf("Progress = %q, want %q", activity.Progress, "Tidying module [2/26]: modules/gitignore")
	}

	// Verify it appears in the snapshot.
	snapshots := ns.SnapshotActivities()
	snap := snapshots["step1"]

	if snap.Progress != "Tidying module [2/26]: modules/gitignore" {
		t.Errorf("snapshot Progress = %q, want %q", snap.Progress, "Tidying module [2/26]: modules/gitignore")
	}
}

// TestActivityProgressClearOnComplete verifies that completing an activity
// clears the progress message (terminal state has no sub-step).
func TestActivityProgressClearOnComplete(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "go-mod-tidy"})
	_ = ns.OnEvent(ctx, ActivityProgress{
		ID:      "step1",
		Name:    "go-mod-tidy",
		Message: "Tidying module [2/26]",
	})
	_ = ns.OnEvent(ctx, ActivityCompleted{
		ID:       "step1",
		Name:     "go-mod-tidy",
		Duration: 3 * time.Second,
	})

	snapshots := ns.SnapshotActivities()
	snap := snapshots["step1"]

	if snap.Progress != "" {
		t.Errorf("Progress after complete = %q, want empty", snap.Progress)
	}
}

// TestActivityProgressEmptyMessage verifies that an empty message clears
// prior progress.
func TestActivityProgressEmptyMessage(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "go-mod-tidy"})
	_ = ns.OnEvent(ctx, ActivityProgress{
		ID:      "step1",
		Name:    "go-mod-tidy",
		Message: "Working on it",
	})
	_ = ns.OnEvent(ctx, ActivityProgress{
		ID:      "step1",
		Name:    "go-mod-tidy",
		Message: "",
	})

	snapshots := ns.SnapshotActivities()
	snap := snapshots["step1"]

	if snap.Progress != "" {
		t.Errorf("Progress after empty = %q, want empty", snap.Progress)
	}
}

// TestActivityProgressSetDirect verifies the SetActivityProgress direct
// accessor works (non-event path).
func TestActivityProgressSetDirect(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "go-mod-tidy"})

	ns.SetActivityProgress("step1", "Tidying [3/26]")

	snapshots := ns.SnapshotActivities()
	snap := snapshots["step1"]

	if snap.Progress != "Tidying [3/26]" {
		t.Errorf("Progress = %q, want %q", snap.Progress, "Tidying [3/26]")
	}
}

// TestActivityRetryingEvent verifies that an ActivityRetrying event
// transitions a failed activity back to running and increments the retry count.
func TestActivityRetryingEvent(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "golangci-lint"})
	_ = ns.OnEvent(ctx, ActivityFailed{
		ID:       "step1",
		Name:     "golangci-lint",
		Err:      errRetrying,
		Duration: 2 * time.Second,
	})

	// Verify it's failed.
	counts := ns.GetActivityCounts()
	if counts.Failed != 1 {
		t.Fatalf("Failed count before retry = %d, want 1", counts.Failed)
	}

	_ = ns.OnEvent(ctx, ActivityRetrying{
		ID:      "step1",
		Name:    "golangci-lint",
		Attempt: 1,
	})

	// After retry: should be running again, retry count = 1.
	counts = ns.GetActivityCounts()
	if counts.Failed != 0 {
		t.Errorf("Failed count after retry = %d, want 0", counts.Failed)
	}

	if counts.Running != 1 {
		t.Errorf("Running count after retry = %d, want 1", counts.Running)
	}

	snapshots := ns.SnapshotActivities()
	snap := snapshots["step1"]

	if snap.RetryCount != 1 {
		t.Errorf("RetryCount = %d, want 1", snap.RetryCount)
	}

	if snap.Status != ActivityStatusRunning {
		t.Errorf("Status = %v, want %v", snap.Status, ActivityStatusRunning)
	}
}

// TestActivityRetryingMultipleAttempts verifies retry count increments
// correctly across multiple retries.
func TestActivityRetryingMultipleAttempts(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "flaky-test"})

	for attempt := 1; attempt <= 3; attempt++ {
		_ = ns.OnEvent(ctx, ActivityFailed{
			ID:       "step1",
			Name:     "flaky-test",
			Duration: 1 * time.Second,
		})
		_ = ns.OnEvent(ctx, ActivityRetrying{
			ID:      "step1",
			Name:    "flaky-test",
			Attempt: attempt,
		})
	}

	snapshots := ns.SnapshotActivities()
	snap := snapshots["step1"]

	if snap.RetryCount != 3 {
		t.Errorf("RetryCount after 3 retries = %d, want 3", snap.RetryCount)
	}
}

// TestSetEstimatedTimeDirect verifies the external estimate injection API.
func TestSetEstimatedTimeDirect(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityRegistered{ID: "step1", Name: "go-build"})

	// Inject estimate from external store (e.g. SQLite).
	ns.SetEstimatedTime("step1", 42*time.Second)

	snapshots := ns.SnapshotActivities()
	snap := snapshots["step1"]

	if snap.EstimatedTime != 42*time.Second {
		t.Errorf("EstimatedTime = %v, want 42s", snap.EstimatedTime)
	}
}

// TestEstimatedRemainingInSummary verifies that the summary bar shows
// "~Xm left" when the estimatedRemaining callback is set.
func TestEstimatedRemainingInSummary(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "build"})

	renderer := NewInlineRenderer(ns, &strings.Builder{}, 0)
	renderer.SetStartTime(time.Now().Add(-10 * time.Second))
	renderer.SetEstimatedRemainingFunc(func() time.Duration {
		return 2 * time.Minute
	})

	summary := renderer.renderSummary(time.Now().Add(-10 * time.Second))

	if !strings.Contains(summary, "~2m left") {
		t.Errorf("summary missing '~2m left':\n%s", summary)
	}
}

// TestEstimatedRemainingZero verifies that ~Xm left is NOT shown when
// the callback returns 0.
func TestEstimatedRemainingZero(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "build"})

	renderer := NewInlineRenderer(ns, &strings.Builder{}, 0)
	renderer.SetStartTime(time.Now().Add(-10 * time.Second))
	renderer.SetEstimatedRemainingFunc(func() time.Duration {
		return 0
	})

	summary := renderer.renderSummary(time.Now().Add(-10 * time.Second))

	if strings.Contains(summary, "left") {
		t.Errorf("summary should not contain 'left' when remaining=0:\n%s", summary)
	}
}

// TestEstimatedRemainingNil verifies that ~Xm left is NOT shown when
// the callback is not set (nil).
func TestEstimatedRemainingNil(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "build"})

	renderer := NewInlineRenderer(ns, &strings.Builder{}, 0)
	renderer.SetStartTime(time.Now().Add(-10 * time.Second))

	summary := renderer.renderSummary(time.Now().Add(-10 * time.Second))

	if strings.Contains(summary, "left") {
		t.Errorf("summary should not contain 'left' when callback is nil:\n%s", summary)
	}
}

// TestProgressRendersInTree verifies that the progress message appears
// in the rendered tree output.
func TestProgressRendersInTree(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "go-mod-tidy"})
	_ = ns.OnEvent(ctx, ActivityProgress{
		ID:      "step1",
		Name:    "go-mod-tidy",
		Message: "Tidying [2/26]",
	})

	frame, ok := ns.RenderSnapshot(0, 0)
	if !ok {
		t.Fatal("RenderSnapshot returned false")
	}

	if !strings.Contains(frame, "Tidying [2/26]") {
		t.Errorf("tree output missing progress message:\n%s", frame)
	}
}

// TestRetryRendersInTree verifies that the retry suffix appears in the
// rendered tree output.
func TestRetryRendersInTree(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "flaky-test"})
	_ = ns.OnEvent(ctx, ActivityFailed{ID: "step1", Name: "flaky-test"})
	_ = ns.OnEvent(ctx, ActivityRetrying{ID: "step1", Name: "flaky-test", Attempt: 1})

	frame, ok := ns.RenderSnapshot(0, 0)
	if !ok {
		t.Fatal("RenderSnapshot returned false")
	}

	if !strings.Contains(frame, string(SymbolRetrying)) {
		t.Errorf("tree output missing retry symbol %q:\n%s", SymbolRetrying, frame)
	}
}

// sentinel error for testing.
var errRetrying = &TransientError{}

type TransientError struct{}

func (*TransientError) Error() string { return "transient failure" }

// TestRetryReasonEvent verifies that the Reason field on ActivityRetrying is
// stored on the activity and rendered as "⟳N (reason)".
func TestRetryReasonEvent(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "flaky-test"})
	_ = ns.OnEvent(ctx, ActivityFailed{ID: "step1", Name: "flaky-test"})
	_ = ns.OnEvent(ctx, ActivityRetrying{
		ID:      "step1",
		Name:    "flaky-test",
		Attempt: 1,
		Reason:  "timeout",
	})

	snapshots := ns.SnapshotActivities()
	snap := snapshots["step1"]

	if snap.RetryReason != "timeout" {
		t.Errorf("RetryReason = %q, want %q", snap.RetryReason, "timeout")
	}

	frame, ok := ns.RenderSnapshot(0, 0)
	if !ok {
		t.Fatal("RenderSnapshot returned false")
	}

	if !strings.Contains(frame, "(timeout)") {
		t.Errorf("tree output missing retry reason:\n%s", frame)
	}
}

// TestRetryReasonEmpty verifies that no "(reason)" suffix renders when the
// Reason field is empty (backward-compatible with the pre-reason behavior).
func TestRetryReasonEmpty(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityStarted{ID: "step1", Name: "flaky-test"})
	_ = ns.OnEvent(ctx, ActivityFailed{ID: "step1", Name: "flaky-test"})
	_ = ns.OnEvent(ctx, ActivityRetrying{ID: "step1", Name: "flaky-test", Attempt: 1})

	frame, ok := ns.RenderSnapshot(0, 0)
	if !ok {
		t.Fatal("RenderSnapshot returned false")
	}

	if strings.Contains(frame, "()") {
		t.Errorf("tree output should not contain empty reason parens:\n%s", frame)
	}

	if !strings.Contains(frame, "⟳1") {
		t.Errorf("tree output should still contain retry count:\n%s", frame)
	}
}

// TestEstimatedTotalRemaining verifies that the subscriber sums the remaining
// estimates of pending activities (full estimate) and excludes completed work.
func TestEstimatedTotalRemaining(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityRegistered{ID: "a", Name: "a"})
	_ = ns.OnEvent(ctx, ActivityRegistered{ID: "b", Name: "b"})

	ns.SetEstimatedTime("a", 10*time.Second)
	ns.SetEstimatedTime("b", 5*time.Second)

	remaining := ns.EstimatedTotalRemaining()
	if remaining != 15*time.Second {
		t.Errorf("EstimatedTotalRemaining = %v, want 15s", remaining)
	}

	// Completing one removes its contribution.
	_ = ns.OnEvent(ctx, ActivityCompleted{ID: "a", Name: "a", Duration: 3 * time.Second})

	remaining = ns.EstimatedTotalRemaining()
	if remaining != 5*time.Second {
		t.Errorf("EstimatedTotalRemaining after complete = %v, want 5s", remaining)
	}
}

// TestEstimatedTotalRemainingZero verifies that 0 is returned when no
// unfinished activity has an estimate.
func TestEstimatedTotalRemainingZero(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: "wf", Name: "test"})
	_ = ns.OnEvent(ctx, ActivityRegistered{ID: "a", Name: "a"})

	if got := ns.EstimatedTotalRemaining(); got != 0 {
		t.Errorf("EstimatedTotalRemaining with no estimates = %v, want 0", got)
	}
}
