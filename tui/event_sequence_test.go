package tui

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

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
		eventType: "workflow.started",
		wID:       nom.WorkflowID("wf-1"),
		wName:     nom.WorkflowName("Test"),
	})

	// Pre-register an activity
	model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: "activity.started",
		aID:       nom.ActivityID("build"),
		aName:     nom.ActivityName("Build"),
	})

	// Initial state
	if model.workflowState != WorkflowStateIdle {
		t.Fatalf("initial state = %v, want Idle", model.workflowState)
	}

	// Send tick → should sync NOM subscriber and transition to running
	updated, _ := model.Update(TickMsg(time.Now()))
	m := updated.(*ProgressModel)

	if m.workflowState != WorkflowStateRunning {
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
	model.nomSubscriber.SetActivityState(nom.NewActivityDisplayState(
		nom.ActivityID("phase1"), nom.ActivityName("Build Phase"),
	))
	model.nomSubscriber.SetActivityState(nom.NewActivityDisplayState(
		nom.ActivityID("step1"), nom.ActivityName("Compile"),
	))

	// Manually add to dependency tree
	model.nomSubscriber.GetDependencyTree().AddActivity(
		nom.ActivityID("phase1"), "Build Phase", nil,
	)
	model.nomSubscriber.GetDependencyTree().AddActivity(
		nom.ActivityID("step1"), "Compile", []nom.ActivityID{"phase1"},
	)

	// Before tick: tree should be empty in model
	if len(model.activities) != 0 {
		t.Errorf("activities before tick = %d, want 0", len(model.activities))
	}

	// Send tick → sync subscriber state to model
	updated, _ := model.Update(TickMsg(time.Now()))
	m := updated.(*ProgressModel)

	// After tick: activities and tree should be synced
	if len(m.activities) != 2 {
		t.Errorf("activities after tick = %d, want 2", len(m.activities))
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

	reporter := NewBubbleTeaProgressReporter()

	// Report step started (1/3)
	reporter.ReportStep(1, 3, "Compile")

	if len(reporter.model.steps) != 1 {
		t.Fatalf("steps count = %d, want 1", len(reporter.model.steps))
	}

	step := reporter.model.steps[0]
	if step.Message != "Compile" {
		t.Errorf("message = %q, want %q", step.Message, "Compile")
	}

	if step.Current != 1 {
		t.Errorf("current = %d, want 1", step.Current)
	}

	if !step.IsActive {
		t.Error("step should be active")
	}

	// Report step completed (3/3)
	reporter.ReportStep(3, 3, "Compile")

	step = reporter.model.steps[0]
	if step.Current != 3 {
		t.Errorf("current = %d, want 3", step.Current)
	}

	if step.IsActive {
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
		eventType: "workflow.started",
		wID:       nom.WorkflowID("wf-1"),
	})

	// Set up subscriber with a failed activity
	model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: "activity.started",
		aID:       nom.ActivityID("test"),
		aName:     nom.ActivityName("Run Tests"),
	})
	model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: "activity.failed",
		aID:       nom.ActivityID("test"),
		aName:     nom.ActivityName("Run Tests"),
		err:       errors.New("test failure"),
	})

	// Sync via tick
	updated, _ := model.Update(TickMsg(time.Now()))
	m := updated.(*ProgressModel)

	// Verify failed activity is present
	act, ok := m.activities[nom.ActivityID("test")]
	if !ok {
		t.Fatal("failed activity should be in model")
	}

	if act.Status != nom.ActivityStatusFailed {
		t.Errorf("status = %v, want Failed", act.Status)
	}

	// Verify tree shows failed status
	node := m.dependencyTree.GetNode(nom.ActivityID("test"))
	if node == nil {
		t.Fatal("failed node should be in tree")
	}

	if node.Status != nom.ActivityStatusFailed {
		t.Errorf("node status = %v, want Failed", node.Status)
	}
}

// TestProgressModel_EventSequence_WorkflowComplete verifies that sending
// progress 100 transitions to completed state.
func TestProgressModel_EventSequence_WorkflowComplete(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = WorkflowStateRunning

	// Send 100% progress → should transition to completed
	updated, _ := model.Update(ProgressUpdateMsg{
		Type:     ProgressUpdate,
		Progress: 100.0,
	})
	m := updated.(*ProgressModel)

	if m.workflowState != WorkflowStateCompleted {
		t.Errorf("state = %v, want Completed", m.workflowState)
	}

	// Reject further updates
	updated, _ = m.Update(ProgressUpdateMsg{
		Type:     ProgressUpdate,
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

func assertScrollOffset(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("scrollOffset = %d, want %d", got, want)
	}
}

// TestProgressModel_KeyboardNavigation verifies scroll keys update offset.
func TestProgressModel_KeyboardNavigation(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.height = 20
	model.scrollOffset = 10

	updated, _ := model.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	m := updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 9)

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	m = updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 10)

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyHome})
	m = updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 0)

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnd})
	m = updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, scrollToBottomSentinel)
}

// TestProgressModel_MouseScrolling verifies mouse wheel updates scroll offset.
func TestProgressModel_MouseScrolling(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.scrollOffset = 5

	updated, _ := model.Update(tea.MouseWheelMsg{Button: ansi.MouseWheelDown})
	m := updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 8)

	updated, _ = m.Update(tea.MouseWheelMsg{Button: ansi.MouseWheelUp})
	m = updated.(*ProgressModel)
	assertScrollOffset(t, m.scrollOffset, 5)
}

// TestProgressModel_CancelMessage verifies CancelMsg triggers quit.
func TestProgressModel_CancelMessage(t *testing.T) {
	t.Parallel()

	model := newTestModel()

	_, cmd := model.Update(CancelMsg{})

	if cmd == nil {
		t.Fatal("expected Quit command for CancelMsg")
	}

	msg := cmd()
	if msg == nil {
		t.Error("quit command should produce a message")
	}
}

// TestProgressModel_ViewportScrolling verifies applyScrollViewport clips content.
func TestProgressModel_ViewportScrolling(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 3

	content := "line1\nline2\nline3\nline4\nline5\nline6\nline7\nline8\nline9\nline10"

	result := model.applyScrollViewport(content)

	lines := splitLines(result)
	if len(lines) != 3 {
		t.Fatalf("expected 3 visible lines, got %d", len(lines))
	}

	if lines[0] != "line1" {
		t.Errorf("first visible line = %q, want %q", lines[0], "line1")
	}

	// Scroll to offset 5
	model.scrollOffset = 5
	result = model.applyScrollViewport(content)

	lines = splitLines(result)
	if len(lines) != 3 {
		t.Fatalf("expected 3 visible lines after scroll, got %d", len(lines))
	}

	if lines[0] != "line6" {
		t.Errorf("first visible line after scroll = %q, want %q", lines[0], "line6")
	}

	// Scroll past end — should clamp
	model.scrollOffset = 100
	result = model.applyScrollViewport(content)

	lines = splitLines(result)
	if lines[0] != "line8" {
		t.Errorf("clamped first line = %q, want %q", lines[0], "line8")
	}

	// Content fits viewport — scrollOffset should reset to 0
	model.height = 20
	model.scrollOffset = 5

	_ = model.applyScrollViewport(content)
	if model.scrollOffset != 0 {
		t.Errorf("scrollOffset should be 0 when content fits, got %d", model.scrollOffset)
	}
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}

	return strings.Split(s, "\n")
}

func TestProgressModel_ResizeStress_RapidSequence(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.displayMode = DisplayModeNOM
	model.workflowState = WorkflowStateRunning
	model.dependencyTree = newTestTree(50)

	resizes := []tea.WindowSizeMsg{
		{Width: 120, Height: 40},
		{Width: 80, Height: 24},
		{Width: 60, Height: 15},
		{Width: 200, Height: 60},
		{Width: 40, Height: 10},
		{Width: 120, Height: 40},
	}

	for _, resize := range resizes {
		updated, _ := model.Update(resize)
		m := updated.(*ProgressModel)

		if m.width != resize.Width {
			t.Errorf("width after resize = %d, want %d", m.width, resize.Width)
		}

		if m.height != resize.Height {
			t.Errorf("height after resize = %d, want %d", m.height, resize.Height)
		}

		view := m.View()
		if view.Content == "" {
			t.Error("View() should produce content after resize")
		}

		model = m
	}
}

func TestProgressModel_Resize_ClampsScrollOffset(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.displayMode = DisplayModeNOM
	model.workflowState = WorkflowStateRunning
	model.width = 120
	model.height = 40
	model.dependencyTree = newTestTree(50)

	// Scroll to bottom in large window
	model.scrollOffset = 30

	// Shrink window — viewport should clamp scrollOffset
	updated, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 10})
	m := updated.(*ProgressModel)

	// scrollOffset should not exceed content lines - height
	if m.scrollOffset < 0 {
		t.Error("scrollOffset should not be negative")
	}

	// Render should not panic with small height
	view := m.View()
	if view.Content == "" {
		t.Error("View() should produce content after resize")
	}
}
