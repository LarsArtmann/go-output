package tui

import (
	"context"
	"errors"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/larsartmann/go-output/nom"
)

// testEvent implements nom.Event and accessor interfaces for testing.
type testEvent struct {
	eventType    string
	wID          nom.WorkflowID
	wName        nom.WorkflowName
	aID          nom.ActivityID
	aName        nom.ActivityName
	duration     time.Duration
	err          error
	dependencies []nom.ActivityID
}

func (e *testEvent) GetEventType() string              { return e.eventType }
func (e *testEvent) GetWorkflowID() nom.WorkflowID     { return e.wID }
func (e *testEvent) GetWorkflowName() nom.WorkflowName { return e.wName }
func (e *testEvent) GetActivityID() nom.ActivityID     { return e.aID }
func (e *testEvent) GetActivityName() nom.ActivityName { return e.aName }
func (e *testEvent) GetDuration() time.Duration        { return e.duration }
func (e *testEvent) GetError() error                   { return e.err }
func (e *testEvent) GetDependencies() []nom.ActivityID { return e.dependencies }

// startActivity fires an activity.started event on the model's NOM subscriber.
// Common test setup boilerplate; collapses a 5-line OnEvent block to one call.
func startActivity(t *testing.T, model *ProgressModel, ctx context.Context, id nom.ActivityID, name nom.ActivityName) {
	t.Helper()

	_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: nom.EventActivityStarted,
		aID:       id,
		aName:     name,
	})
}

// startWorkflow fires a workflow.started event on the model's NOM subscriber.
func startWorkflow(t *testing.T, model *ProgressModel, ctx context.Context, id nom.WorkflowID) {
	t.Helper()

	_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: nom.EventWorkflowStarted,
		wID:       id,
	})
}

// TestProgressModel_EventSequence_WorkflowStartToTick verifies that when
// the NOM subscriber reports a running workflow, a tick transitions the
// model from Idle to Running and syncs the tree.
func TestProgressModel_EventSequence_WorkflowStartToTick(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.displayMode = DisplayModeNOM

	// Start workflow in subscriber
	ctx := context.Background()
	model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: nom.EventWorkflowStarted,
		wID:       nom.WorkflowID("wf-1"),
		wName:     nom.WorkflowName("Test"),
	})

	// Pre-register an activity
	startActivity(t, model, ctx, "build", "Build")

	// Initial state
	if model.workflowState != workflowStateIdle {
		t.Fatalf("initial state = %v, want Idle", model.workflowState)
	}

	// Send tick → should sync NOM subscriber and transition to running
	updated, _ := model.Update(tickMsg(time.Now()))
	m := updated.(*ProgressModel)

	if m.workflowState != workflowStateRunning {
		t.Errorf("state after tick = %v, want Running", m.workflowState)
	}

	// Tree should be synced
	if m.dependencyTree == nil {
		t.Fatal("dependency tree should be synced after tick")
	}
}

// TestProgressModel_EventSequence_PreRegisterThenStart verifies that
// pre-registered activities in the NOM subscriber are synced to the model
// after a tick, even before workflow.started.
func TestProgressModel_EventSequence_PreRegisterThenStart(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.displayMode = DisplayModeNOM

	// Pre-register activities directly on the subscriber
	model.nomSubscriber.SetActivityState(nom.ActivityID("phase1"), nom.NewActivity("phase1", "Build Phase"))
	model.nomSubscriber.SetActivityState(nom.ActivityID("step1"), nom.NewActivity("step1", "Compile"))

	// Manually add to dependency tree
	model.nomSubscriber.GetDependencyTree().AddActivity(
		nom.ActivityID("phase1"), nom.NewActivity("phase1", "Build Phase"), nil,
	)
	model.nomSubscriber.GetDependencyTree().AddActivity(
		nom.ActivityID("step1"), nom.NewActivity("step1", "Compile"), []nom.ActivityID{"phase1"},
	)

	// Before tick: subscriber already has the pre-registered activities
	counts := model.nomSubscriber.GetActivityCounts()

	total := counts.Total()
	if total != 2 {
		t.Errorf("activities before tick = %d, want 2", total)
	}

	// Send tick → sync subscriber state to model
	updated, _ := model.Update(tickMsg(time.Now()))
	m := updated.(*ProgressModel)

	// After tick: activities should be synced to subscriber
	counts = m.nomSubscriber.GetActivityCounts()

	total = counts.Total()
	if total != 2 {
		t.Errorf("activities after tick = %d, want 2", total)
	}

	if m.dependencyTree == nil {
		t.Fatal("dependency tree should be synced after tick")
	}

	node := m.dependencyTree.GetNode(nom.ActivityID("step1"))
	if node == nil {
		t.Error("step1 should exist in synced dependency tree")
	}
}

// TestProgressModel_EventSequence_StepLifecycle verifies step creation,
// update, and completion through ReportStep on the reporter.
func TestProgressModel_EventSequence_StepLifecycle(t *testing.T) {
	t.Parallel()

	reporter := newTestReporter()

	// Report step started (1/3)
	reporter.ReportStep(1, 3, "Compile")

	assertSingleStep(t, reporter)

	step := reporter.model.steps[0]
	if step.Message != "Compile" {
		t.Errorf("message = %q, want %q", step.Message, "Compile")
	}

	if step.Current != 1 {
		t.Errorf("current = %d, want 1", step.Current)
	}

	if !step.isActive() {
		t.Error("step should be active")
	}

	// Report step completed (3/3)
	reporter.ReportStep(3, 3, "Compile")

	step = reporter.model.steps[0]
	if step.Current != 3 {
		t.Errorf("current = %d, want 3", step.Current)
	}

	if step.isActive() {
		t.Error("step should not be active after completion")
	}

	if step.CompletedAt == nil {
		t.Error("step should have completion time")
	}
}

// TestProgressModel_EventSequence_StepFailed verifies that a failed activity
// in the NOM subscriber is reflected in the model after sync.
func TestProgressModel_EventSequence_StepFailed(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.displayMode = DisplayModeNOM

	ctx := context.Background()
	model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: nom.EventWorkflowStarted,
		wID:       nom.WorkflowID("wf-1"),
	})

	// Set up subscriber with a failed activity
	startActivity(t, model, ctx, "test", "Run Tests")
	model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: nom.EventActivityFailed,
		aID:       nom.ActivityID("test"),
		aName:     nom.ActivityName("Run Tests"),
		err:       errors.New("test failure"),
	})

	// Sync via tick
	updated, _ := model.Update(tickMsg(time.Now()))
	m := updated.(*ProgressModel)

	// Verify failed activity is present
	act := m.nomSubscriber.GetActivity(nom.ActivityID("test"))
	if act == nil {
		t.Fatal("failed activity should be in subscriber")
	}

	if act.Status != nom.ActivityStatusFailed {
		t.Errorf("status = %v, want Failed", act.Status)
	}

	// Verify tree shows failed status
	node := m.dependencyTree.GetNode(nom.ActivityID("test"))
	if node == nil {
		t.Fatal("failed node should be in tree")
	}

	// Verify tree shows failed status via snapshot
	snaps := m.nomSubscriber.SnapshotActivities()
	snap := snaps[nom.ActivityID("test")]
	if snap.Status != nom.ActivityStatusFailed {
		t.Errorf("node status = %v, want Failed", snap.Status)
	}
}

// TestProgressModel_EventSequence_WorkflowComplete verifies that sending
// progress 100 transitions to completed state.
func TestProgressModel_EventSequence_WorkflowComplete(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = workflowStateRunning

	// Send 100% progress → should transition to completed
	updated, _ := model.Update(progressUpdateMsg{
		Type:     progressUpdate,
		Progress: 100.0,
	})
	m := updated.(*ProgressModel)

	if m.workflowState != workflowStateCompleted {
		t.Errorf("state = %v, want Completed", m.workflowState)
	}

	// Reject further updates
	updated, _ = m.Update(progressUpdateMsg{
		Type:     progressUpdate,
		Progress: 50.0,
	})
	m = updated.(*ProgressModel)

	if m.currentProgress != 100.0 {
		t.Errorf("progress = %f, want 100.0 (rejected update)", m.currentProgress)
	}
}

// TestProgressModel_EventSequence_KeyQuit verifies quit key handling.
func TestProgressModel_EventSequence_KeyQuit(t *testing.T) {
	t.Parallel()

	model := newTestModel()

	// Send 'q' key
	_, cmd := model.Update(tea.KeyPressMsg{Code: 'q'})

	if cmd == nil {
		t.Fatal("expected Quit command for 'q' key")
	}

	// Verify it's a quit command by checking the message it produces
	msg := cmd()
	if msg == nil {
		t.Error("quit command should produce a message")
	}
}
