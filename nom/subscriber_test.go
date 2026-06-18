package nom

import (
	"context"
	"errors"
	"testing"
	"time"
)

type testEvent struct {
	eventType string
	workflow  bool
	activity  bool
	wID       WorkflowID
	wName     WorkflowName
	aID       ActivityID
	aName     ActivityName
	deps      []ActivityID
	duration  time.Duration
	err       error
}

func (e *testEvent) GetEventType() string { return e.eventType }

func (e *testEvent) GetWorkflowID() WorkflowID     { return e.wID }
func (e *testEvent) GetWorkflowName() WorkflowName { return e.wName }
func (e *testEvent) GetActivityID() ActivityID     { return e.aID }
func (e *testEvent) GetActivityName() ActivityName { return e.aName }
func (e *testEvent) GetDuration() time.Duration    { return e.duration }
func (e *testEvent) GetError() error               { return e.err }
func (e *testEvent) GetDependencies() []ActivityID { return e.deps }

// setupWithWorkflow creates a subscriber and fires workflow.started.
func setupWithWorkflow(t *testing.T) (*NOMStyleSubscriber, context.Context) {
	t.Helper()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	ns.OnEvent(ctx, &testEvent{
		eventType: EventWorkflowStarted,
		workflow:  true,
		wID:       WorkflowID("wf-1"),
	})

	return ns, ctx
}

// sendActivityStarted fires an activity.started event with the given ID and name.
func sendActivityStarted(t *testing.T, ns *NOMStyleSubscriber, ctx context.Context, id ActivityID, name ActivityName) {
	t.Helper()

	err := ns.OnEvent(ctx, &testEvent{
		eventType: EventActivityStarted,
		activity:  true,
		aID:       id,
		aName:     name,
	})
	if err != nil {
		t.Fatalf("activity.started OnEvent() error: %v", err)
	}
}

// registerActivity fires an activity.registered event with optional dependencies.
// Use for golden test setup where the same workflow is repeated across frames.
func registerActivity(
	ns *NOMStyleSubscriber,
	ctx context.Context,
	id ActivityID,
	name ActivityName,
	deps ...ActivityID,
) {
	_ = ns.OnEvent(ctx, &testEvent{
		eventType: EventActivityRegistered,
		aID:       id,
		aName:     name,
		deps:      deps,
	})
}

func TestNewNOMStyleSubscriber(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	if ns == nil {
		t.Fatal("NewNOMStyleSubscriber() returned nil")
	}

	if !ns.IsEnabled() {
		t.Error("new subscriber should be enabled")
	}

	if ns.IsWorkflowRunning() {
		t.Error("new subscriber should not be running")
	}
}

func TestNOMStyleSubscriber_Configuration(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()

	ns.SetEnabled(false)

	if ns.IsEnabled() {
		t.Error("should be disabled after SetEnabled(false)")
	}

	ns.SetEnabled(true)

	if !ns.IsEnabled() {
		t.Error("should be enabled after SetEnabled(true)")
	}
}

func TestNOMStyleSubscriber_Reset(t *testing.T) {
	t.Parallel()

	// Reset test uses a custom workflow name to verify it gets cleared.
	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	ns.OnEvent(ctx, &testEvent{
		eventType: EventWorkflowStarted,
		workflow:  true,
		wID:       WorkflowID("wf-1"),
		wName:     WorkflowName("CI"),
	})
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("Build"))

	ns.Reset()

	if ns.IsWorkflowRunning() {
		t.Error("should not be running after Reset()")
	}

	if ns.GetWorkflowID() != "" {
		t.Error("workflow ID should be empty after Reset()")
	}

	if ns.GetWorkflowName() != "" {
		t.Error("workflow name should be empty after Reset()")
	}

	activities := ns.GetActivities()
	if len(activities) != 0 {
		t.Errorf("activities should be empty after Reset(), got %d", len(activities))
	}
}

func TestNOMStyleSubscriber_WorkflowStarted(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	ctx := context.Background()

	err := ns.OnEvent(ctx, &testEvent{
		eventType: EventWorkflowStarted,
		workflow:  true,
		wID:       WorkflowID("wf-1"),
		wName:     WorkflowName("Deploy"),
	})
	if err != nil {
		t.Fatalf("OnEvent() error: %v", err)
	}

	if !ns.IsWorkflowRunning() {
		t.Error("workflow should be running after started event")
	}

	if ns.GetWorkflowID() != "wf-1" {
		t.Errorf("workflow ID = %q, want %q", ns.GetWorkflowID(), "wf-1")
	}

	if ns.GetWorkflowName() != "Deploy" {
		t.Errorf("workflow name = %q, want %q", ns.GetWorkflowName(), "Deploy")
	}

	if ns.GetStartTime().IsZero() {
		t.Error("start time should be set")
	}
}

func TestNOMStyleSubscriber_WorkflowCompleted(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)

	err := ns.OnEvent(ctx, &testEvent{
		eventType: EventWorkflowCompleted,
		workflow:  true,
		wID:       WorkflowID("wf-1"),
	})
	if err != nil {
		t.Fatalf("OnEvent() error: %v", err)
	}

	if ns.IsWorkflowRunning() {
		t.Error("workflow should not be running after completed event")
	}
}

func TestNOMStyleSubscriber_WorkflowFailed(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)

	err := ns.OnEvent(ctx, &testEvent{
		eventType: EventWorkflowFailed,
		err:       errors.New("timeout"),
	})
	if err != nil {
		t.Fatalf("OnEvent() error: %v", err)
	}

	if ns.IsWorkflowRunning() {
		t.Error("workflow should not be running after failed event")
	}
}

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

	if activity.ActivityName != ActivityName("Build") {
		t.Errorf("activity name = %q, want %q", activity.ActivityName, "Build")
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

	if activity.Error == nil {
		t.Error("activity should have error")
	}
}

func TestNOMStyleSubscriber_UnknownEventType(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
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

	ns := NewNOMStyleSubscriber()
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

	ns := NewNOMStyleSubscriber()

	paused := NewActivityDisplayState(ActivityID("p1"), ActivityName("Wait"))
	paused.Status = ActivityStatusPaused
	ns.SetActivityState(paused)

	failed := NewActivityDisplayState(ActivityID("f1"), ActivityName("Build"))
	failed.SetFailed(errors.New("crash"))
	ns.SetActivityState(failed)

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

	ns := NewNOMStyleSubscriber()

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

func TestNOMStyleSubscriber_SyncActivityTimingToTree(t *testing.T) {
	t.Parallel()

	ns, ctx := setupWithWorkflow(t)
	sendActivityStarted(t, ns, ctx, ActivityID("a1"), ActivityName("Build"))

	time.Sleep(10 * time.Millisecond)
	ns.UpdateRunningActivityElapsed()
	ns.SyncActivityTimingToTree()

	tree := ns.GetDependencyTree()

	node := tree.GetNode(ActivityID("a1"))
	if node == nil {
		t.Fatal("tree node should exist")
	}

	if node.CurrentElapsed <= 0 {
		t.Error("tree node CurrentElapsed should be positive after sync")
	}
}

func TestNOMStyleSubscriber_SetActivityState(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()
	activity := NewActivityDisplayState(ActivityID("custom"), ActivityName("Custom"))
	activity.SetRunning()

	ns.SetActivityState(activity)

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

	ns := NewNOMStyleSubscriber()

	tree := ns.GetDependencyTree()
	if tree == nil {
		t.Error("GetDependencyTree() should not return nil")
	}
}

func TestNOMStyleSubscriber_GetTimingCache(t *testing.T) {
	t.Parallel()

	ns := NewNOMStyleSubscriber()

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
