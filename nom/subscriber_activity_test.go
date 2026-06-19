package nom

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNOMStyleSubscriber_ActivityStarted(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)

	err := ns.OnEvent(ctx, &testEvent{
		eventType: EventActivityStarted,
		activity:  true,
		aID:       ActivityID("a1"),
		aName:     ActivityName("Build"),
	})
	if err != nil {
		t.Fatalf("OnEvent() error: %v", err)
	}

	activity := ns.GetActivity(ActivityID("a1"))
	if activity == nil {
		t.Fatal("activity should exist after started event")
	}

	if activity.Status != ActivityStatusRunning {
		t.Errorf("activity status = %v, want Running", activity.Status)
	}

	if activity.Label.Get() != "Build" {
		t.Errorf("activity name = %q, want %q", activity.Label.Get(), "Build")
	}
}

func TestNOMStyleSubscriber_ActivityCompleted(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("Build"))

	err := ns.OnEvent(ctx, &testEvent{
		eventType: EventActivityCompleted,
		activity:  true,
		aID:       ActivityID("a1"),
		aName:     ActivityName("Build"),
		duration:  5 * time.Second,
	})
	if err != nil {
		t.Fatalf("OnEvent() error: %v", err)
	}

	activity := ns.GetActivity(ActivityID("a1"))
	if activity == nil {
		t.Fatal("activity should exist")
	}

	if activity.Status != ActivityStatusCompleted {
		t.Errorf("activity status = %v, want Completed", activity.Status)
	}
}

func TestNOMStyleSubscriber_ActivityFailed(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("Build"))

	testErr := errors.New("build failed")

	err := ns.OnEvent(ctx, &testEvent{
		eventType: EventActivityFailed,
		activity:  true,
		aID:       ActivityID("a1"),
		aName:     ActivityName("Build"),
		err:       testErr,
		duration:  3 * time.Second,
	})
	if err != nil {
		t.Fatalf("OnEvent() error: %v", err)
	}

	activity := ns.GetActivity(ActivityID("a1"))
	if activity == nil {
		t.Fatal("activity should exist")
	}

	if activity.Status != ActivityStatusFailed {
		t.Errorf("activity status = %v, want Failed", activity.Status)
	}

	if activity.Err == nil {
		t.Error("activity should have error")
	}
}

func TestNOMStyleSubscriber_UnknownEventType(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)
	ctx := context.Background()

	err := ns.OnEvent(ctx, &testEvent{eventType: "unknown.event"})
	if err != nil {
		t.Errorf("unknown event type should not return error: %v", err)
	}
}

type minimalEvent struct {
	eventType string
}

func (e *minimalEvent) GetEventType() string { return e.eventType }

func TestNOMStyleSubscriber_EventWithoutAccessor(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)
	ctx := context.Background()

	err := ns.OnEvent(ctx, &minimalEvent{eventType: EventWorkflowStarted})
	if err != nil {
		t.Errorf("event without accessor should not return error: %v", err)
	}
}

func TestNOMStyleSubscriber_GetActivities(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("Build"))
	sendActivityStarted(t, ns, ctx, ActivityID("a2"), ActivityName("Test"))

	activities := ns.GetActivities()
	if len(activities) != 2 {
		t.Errorf("expected 2 activities, got %d", len(activities))
	}
}

func TestNOMStyleSubscriber_GetActivityCounts(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("Build"))
	sendActivityStarted(t, ns, ctx, ActivityID("a2"), ActivityName("Test"))

	ns.OnEvent(ctx, &testEvent{
		eventType: EventActivityCompleted,
		activity:  true,
		aID:       ActivityID("a1"),
		aName:     ActivityName("Build"),
	})

	counts := ns.GetActivityCounts()
	if counts.Running != 1 {
		t.Errorf("running = %d, want 1", counts.Running)
	}

	if counts.Completed != 1 {
		t.Errorf("completed = %d, want 1", counts.Completed)
	}

	if counts.Failed != 0 {
		t.Errorf("failed = %d, want 0", counts.Failed)
	}

	if counts.Pending != 0 {
		t.Errorf("pending = %d, want 0", counts.Pending)
	}
}

func TestNOMStyleSubscriber_GetActivityCounts_PausedAndFailed(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)

	paused := NewActivity("p1", "Wait")
	paused.Status = ActivityStatusPaused
	ns.SetActivityState(ActivityID("p1"), paused)

	failed := NewActivity("f1", "Build")
	failed.SetFailed(errors.New("crash"))
	ns.SetActivityState(ActivityID("f1"), failed)

	counts := ns.GetActivityCounts()
	if counts.Running != 0 {
		t.Errorf("running = %d, want 0", counts.Running)
	}

	if counts.Completed != 0 {
		t.Errorf("completed = %d, want 0", counts.Completed)
	}

	if counts.Failed != 1 {
		t.Errorf("failed = %d, want 1", counts.Failed)
	}

	if counts.Pending != 1 {
		t.Errorf("pending = %d, want 1 (paused counts as pending)", counts.Pending)
	}
}

func TestNOMStyleSubscriber_GetActivity_NotFound(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)

	activity := ns.GetActivity(ActivityID("nonexistent"))
	if activity != nil {
		t.Error("expected nil for nonexistent activity")
	}
}

func TestNOMStyleSubscriber_UpdateRunningActivityElapsed(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("Build"))

	time.Sleep(10 * time.Millisecond)
	ns.UpdateRunningActivityElapsed()

	activity := ns.GetActivity(ActivityID("a1"))
	if activity == nil {
		t.Fatal("activity should exist")
	}

	if activity.CurrentElapsed <= 0 {
		t.Error("CurrentElapsed should be positive after UpdateRunningActivityElapsed()")
	}
}

func TestNOMStyleSubscriber_SetActivityState(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)
	activity := NewActivity("custom", "Custom")
	activity.SetRunning()

	ns.SetActivityState(ActivityID("custom"), activity)

	got := ns.GetActivity(ActivityID("custom"))
	if got == nil {
		t.Fatal("activity should exist after SetActivityState")
	}

	if got.Status != ActivityStatusRunning {
		t.Errorf("status = %v, want Running", got.Status)
	}
}

func TestNOMStyleSubscriber_GetDependencyTree(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)

	tree := ns.GetDependencyTree()
	if tree == nil {
		t.Error("GetDependencyTree() should not return nil")
	}
}

func TestNOMStyleSubscriber_GetTimingCache(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)

	cache := ns.GetTimingCache()
	if cache == nil {
		t.Error("GetTimingCache() should not return nil")
	}
}

func TestNOMStyleSubscriber_ActivitiesDeepCopy(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("Build"))

	activities := ns.GetActivities()
	activities[ActivityID("a1")].Status = ActivityStatusCompleted

	original := ns.GetActivity(ActivityID("a1"))
	if original.Status != ActivityStatusRunning {
		t.Error("modifying returned map should not affect internal state (deep copy failed)")
	}
}
