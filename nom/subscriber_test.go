package nom

import (
	"context"
	"errors"
	"testing"
	"time"
)

// setupWithWorkflow creates a subscriber and fires workflow.started.
func setupWithWorkflow(t *testing.T) (*NOMStyleSubscriber, context.Context) {
	t.Helper()

	ns := newTestSubscriber(t)
	ctx := context.Background()

	_ = ns.OnEvent(ctx, WorkflowStarted{ID: WorkflowID("wf-1")})

	return ns, ctx
}

// sendActivityStarted fires an activity.started event with the given ID and name.
func sendActivityStarted(t *testing.T, ns *NOMStyleSubscriber, ctx context.Context, id ActivityID, name ActivityName) {
	t.Helper()

	if err := ns.OnEvent(ctx, ActivityStarted{ID: id, Name: name}); err != nil {
		t.Fatalf("activity.started OnEvent() error: %v", err)
	}
}

// sendActivityCompleted fires an activity.completed event with the given ID, name, and duration.
func sendActivityCompleted(
	t *testing.T,
	ns *NOMStyleSubscriber,
	ctx context.Context,
	id ActivityID,
	name ActivityName,
	duration time.Duration,
) {
	t.Helper()

	if err := ns.OnEvent(ctx, ActivityCompleted{ID: id, Name: name, Duration: duration}); err != nil {
		t.Fatalf("activity.completed OnEvent() error: %v", err)
	}
}

// sendWorkflowStarted fires a workflow.started event with the given ID and name.
// Returns the error so callers can choose to assert or ignore it.
func sendWorkflowStarted(ns *NOMStyleSubscriber, ctx context.Context, id WorkflowID, name WorkflowName) error {
	return ns.OnEvent(ctx, WorkflowStarted{ID: id, Name: name})
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
	_ = ns.OnEvent(ctx, ActivityRegistered{ID: id, Name: name, Deps: deps})
}

// registerPhase fires an activity.registered event with Kind=Phase, so the
// node renders as a phase grouping (SymbolPhase/Colors.Phase).
func registerPhase(
	ns *NOMStyleSubscriber,
	ctx context.Context,
	id ActivityID,
	name ActivityName,
	deps ...ActivityID,
) {
	_ = ns.OnEvent(ctx, ActivityRegistered{ID: id, Name: name, Kind: ActivityKindPhase, Deps: deps})
}

func TestNewNOMStyleSubscriber(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)
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

	ns := newTestSubscriber(t)

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
	ns := newTestSubscriber(t)
	ctx := context.Background()

	_ = sendWorkflowStarted(ns, ctx, WorkflowID("wf-1"), WorkflowName("CI"))
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

	ns := newTestSubscriber(t)
	ctx := context.Background()

	err := sendWorkflowStarted(ns, ctx, WorkflowID("wf-1"), WorkflowName("Deploy"))
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

	err := ns.OnEvent(ctx, WorkflowCompleted{ID: WorkflowID("wf-1")})
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

	err := ns.OnEvent(ctx, WorkflowFailed{ID: WorkflowID("wf-1"), Err: errors.New("timeout")})
	if err != nil {
		t.Fatalf("OnEvent() error: %v", err)
	}

	if ns.IsWorkflowRunning() {
		t.Error("workflow should not be running after failed event")
	}
}
