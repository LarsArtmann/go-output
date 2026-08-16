package nom

import (
	"context"
	"errors"
	"testing"

	"github.com/larsartmann/go-output"
)

// recount brute-force recomputes ActivityCounts by scanning the subscriber's
// internal activities map. Used as the ground truth to verify the incremental
// count cache stays in sync across every state transition.
func recount(ns *NOMSubscriber) ActivityCounts {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	var c ActivityCounts

	for _, a := range ns.activities {
		switch a.Status {
		case ActivityStatusRunning:
			c.Running++
		case ActivityStatusCompleted:
			c.Completed++
		case ActivityStatusFailed:
			c.Failed++
		case ActivityStatusPending:
			c.Pending++
		default:
			c.Other++
		}
	}

	return c
}

func assertCountsMatch(t *testing.T, ns *NOMSubscriber, msg string) {
	t.Helper()

	got := ns.GetActivityCounts()
	want := recount(ns)

	if got != want {
		t.Errorf("%s: cached counts = %+v, want %+v (brute-force recount)", msg, got, want)
	}
}

// TestActivityCountsCache_LifecycleConsistency verifies the incremental count
// cache matches a brute-force recount after every event in a realistic mixed
// workflow. This is the core invariant: the cache can never drift from the
// ground truth, regardless of transition order.
func TestActivityCountsCache_LifecycleConsistency(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)
	ctx := context.Background()
	_ = ns.OnEvent(ctx, WorkflowStarted{ID: WorkflowID("wf"), Name: WorkflowName("build")})

	assertCountsMatch(t, ns, "after workflow started")

	// Register 3 pending activities.
	for _, id := range []ActivityID{"a1", "a2", "a3"} {
		_ = ns.OnEvent(ctx, ActivityRegistered{ID: id, Name: ActivityName("step")})
	}

	assertCountsMatch(t, ns, "after 3 registered (pending)")

	if c := ns.GetActivityCounts(); c.Pending != 3 {
		t.Fatalf("pending = %d, want 3", c.Pending)
	}

	// Start all 3: Pending → Running.
	for _, id := range []ActivityID{"a1", "a2", "a3"} {
		_ = ns.OnEvent(ctx, ActivityStarted{ID: id, Name: ActivityName("step")})
	}

	assertCountsMatch(t, ns, "after 3 started")

	if c := ns.GetActivityCounts(); c.Running != 3 || c.Pending != 0 {
		t.Fatalf("running=%d pending=%d, want running=3 pending=0", c.Running, c.Pending)
	}

	// Complete a1, fail a2, leave a3 running.
	_ = ns.OnEvent(ctx, ActivityCompleted{ID: ActivityID("a1"), Name: ActivityName("step")})
	assertCountsMatch(t, ns, "after a1 completed")

	_ = ns.OnEvent(ctx, ActivityFailed{
		ID:   ActivityID("a2"),
		Name: ActivityName("step"),
		Err:  errors.New("boom"),
	})
	assertCountsMatch(t, ns, "after a2 failed")

	want := ActivityCounts{Running: 1, Completed: 1, Failed: 1}
	if got := ns.GetActivityCounts(); got != want {
		t.Errorf("final counts = %+v, want %+v", got, want)
	}
}

// TestActivityCountsCache_IdempotentTransitions verifies that re-firing the
// same transition (Running → Running) does not corrupt the cache. applyDelta
// must be a no-op when from == to.
func TestActivityCountsCache_IdempotentTransitions(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)
	ctx := context.Background()
	_ = ns.OnEvent(ctx, WorkflowStarted{ID: WorkflowID("wf")})

	id := ActivityID("a1")
	_ = ns.OnEvent(ctx, ActivityStarted{ID: id, Name: ActivityName("step")})
	before := ns.GetActivityCounts()

	// Re-fire ActivityStarted on an already-running activity.
	_ = ns.OnEvent(ctx, ActivityStarted{ID: id, Name: ActivityName("step")})
	after := ns.GetActivityCounts()

	if after != before {
		t.Errorf("re-firing Running→Running changed counts: before=%+v after=%+v", before, after)
	}

	// Re-fire ActivityCompleted on an already-completed activity.
	_ = ns.OnEvent(ctx, ActivityCompleted{ID: id, Name: ActivityName("step")})
	completed := ns.GetActivityCounts()
	_ = ns.OnEvent(ctx, ActivityCompleted{ID: id, Name: ActivityName("step")})
	again := ns.GetActivityCounts()

	if again != completed {
		t.Errorf("re-firing Completed→Completed changed counts: first=%+v second=%+v", completed, again)
	}
}

// TestActivityCountsCache_Reset verifies that Reset() zeroes the count cache
// along with the activities map. A stale cache after Reset would be a silent
// correctness bug.
func TestActivityCountsCache_Reset(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("build"))

	if c := ns.GetActivityCounts(); c.Running != 1 {
		t.Fatalf("running = %d, want 1 before reset", c.Running)
	}

	ns.Reset()

	if c := ns.GetActivityCounts(); c != (ActivityCounts{}) {
		t.Errorf("counts after Reset = %+v, want zero value", c)
	}
}

// TestActivityCountsCache_SkipRegistration verifies the common path where
// ActivityStarted fires without a prior ActivityRegistered — the activity is
// created and counted as Pending first, then immediately transitions to Running
// in the same handler call.
func TestActivityCountsCache_SkipRegistration(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)
	ctx := context.Background()
	_ = ns.OnEvent(ctx, WorkflowStarted{ID: WorkflowID("wf")})

	// Direct start with no prior registration: getOrCreateActivity creates it
	// as Pending (+1 Pending), then applyDelta flips Pending→Running.
	_ = ns.OnEvent(ctx, ActivityStarted{ID: ActivityID("x"), Name: ActivityName("x")})
	assertCountsMatch(t, ns, "after direct start (skip registration)")

	if c := ns.GetActivityCounts(); c.Running != 1 || c.Pending != 0 {
		t.Errorf("direct start: running=%d pending=%d, want running=1 pending=0", c.Running, c.Pending)
	}
}

// TestActivityCountsCache_SetActivityState verifies the test helper maintains
// the cache correctly on both insertion and replacement.
func TestActivityCountsCache_SetActivityState(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)

	// Insert a running activity.
	a := NewActivity("a1", "Build")
	a.SetRunning()
	ns.SetActivityState(ActivityID("a1"), a)
	assertCountsMatch(t, ns, "after first SetActivityState")

	if c := ns.GetActivityCounts(); c.Running != 1 {
		t.Fatalf("running = %d, want 1", c.Running)
	}

	// Replace with a completed activity: Running → Completed.
	b := NewActivity("a1", "Build")
	b.SetCompleted()
	ns.SetActivityState(ActivityID("a1"), b)
	assertCountsMatch(t, ns, "after replace Running→Completed")

	if c := ns.GetActivityCounts(); c.Running != 0 || c.Completed != 1 {
		t.Errorf("after replace: running=%d completed=%d, want 0/1", c.Running, c.Completed)
	}
}

// TestActivityCountsCache_CustomStatusNotLost verifies that activities in
// registered custom statuses land in the Other bucket instead of vanishing
// from counts, totals, and percentages. The open status registry must never
// make activities invisible to the count cache.
func TestActivityCountsCache_CustomStatusNotLost(t *testing.T) {
	t.Parallel()

	custom := RegisterStatus(
		"counts-test-skipped",
		SymbolOther,
		Colors.Pending,
		3,
		output.NodeShapeBox,
		output.NodeStyle{}, //nolint:exhaustruct // Count cache test ignores diagram style
	)

	ns := newTestSubscriber(t)

	a := NewActivity("a1", "Build")
	a.Status = custom
	ns.SetActivityState(ActivityID("a1"), a)
	assertCountsMatch(t, ns, "after custom-status insert")

	counts := ns.GetActivityCounts()
	if counts.Other != 1 || counts.Total() != 1 {
		t.Errorf("custom status lost: counts=%+v, want Other=1 Total=1", counts)
	}

	// Transition custom → completed: Other decremented, Completed incremented.
	b := NewActivity("a1", "Build")
	b.SetCompleted()
	ns.SetActivityState(ActivityID("a1"), b)
	assertCountsMatch(t, ns, "after custom→completed replace")

	if c := ns.GetActivityCounts(); c.Other != 0 || c.Completed != 1 {
		t.Errorf("after custom→completed: other=%d completed=%d, want 0/1", c.Other, c.Completed)
	}
}

// TestSetActivityStateNilGuard verifies a nil activity is ignored instead of
// panicking on activity.Status.
func TestSetActivityStateNilGuard(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)

	ns.SetActivityState(ActivityID("a1"), nil)

	if c := ns.GetActivityCounts(); c != (ActivityCounts{}) {
		t.Errorf("counts after nil SetActivityState = %+v, want zero value", c)
	}
}
