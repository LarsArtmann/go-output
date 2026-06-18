package tui

import (
	"context"
	"errors"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/larsartmann/go-output/nom"
)

func newTestModel() *ProgressModel {
	return NewProgressModel()
}

// newTestReporter creates a reporter that never starts the real Bubble Tea
// program. Setting started=true causes ensureStarted() to be a no-op, so
// reporter methods (ReportMessage, ReportProgress, etc.) mutate model fields
// directly without spawning a goroutine that races with test assertions.
func newTestReporter() *BubbleTeaProgressReporter {
	r := NewBubbleTeaProgressReporter()
	r.started = true

	return r
}

func addTestActivity(model *ProgressModel, id, name string, status nom.ActivityStatus) {
	activity := nom.NewActivity(id, name)

	activity.Status = status
	switch status {
	case nom.ActivityStatusRunning:
		activity.SetRunning()
	case nom.ActivityStatusCompleted:
		activity.SetCompleted()
	case nom.ActivityStatusFailed:
		activity.SetFailed(errors.New("test failure"))
	case nom.ActivityStatusPending, nom.ActivityStatusPaused:
		// No transition side-effects; Status is already set above.
	}

	model.nomSubscriber.SetActivityState(nom.ActivityID(id), activity)
}

func setupTestTree(model *ProgressModel) *nom.DependencyTree {
	tree := nom.NewDependencyTree()
	model.dependencyTree = tree
	model.treeStartLine = 2

	return tree
}

// clickAt simulates a mouse click at the given Y coordinate and returns the updated model.
func clickAt(model *ProgressModel, clickY int) *ProgressModel {
	updatedModel, _ := model.Update(tea.MouseClickMsg{
		X: 5, Y: clickY, Button: tea.MouseLeft,
	})

	return updatedModel.(*ProgressModel)
}

func TestProgressModel_Update_WindowSize(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	m := updatedModel.(*ProgressModel)
	if m.width != 80 {
		t.Errorf("width = %d, want 80", m.width)
	}

	if m.height != 24 {
		t.Errorf("height = %d, want 24", m.height)
	}
}

func TestProgressModel_Update_KeyPressQuit(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	updatedModel, cmd := model.Update(tea.KeyPressMsg{})
	_ = updatedModel
	_ = cmd
}

func TestProgressModel_Update_ProgressUpdateMsg(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = WorkflowStateRunning

	updatedModel, _ := model.Update(ProgressUpdateMsg{
		Type:     ProgressUpdate,
		Progress: 75.0,
	})

	m := updatedModel.(*ProgressModel)
	if m.currentProgress != 75.0 {
		t.Errorf("progress = %f, want 75.0", m.currentProgress)
	}
}

func TestProgressModel_Update_ProgressUpdateMsg_Completion(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = WorkflowStateRunning

	updatedModel, _ := model.Update(ProgressUpdateMsg{
		Type:     ProgressUpdate,
		Progress: 100.0,
	})

	m := updatedModel.(*ProgressModel)
	if m.workflowState != WorkflowStateCompleted {
		t.Errorf("workflow state = %v, want Completed", m.workflowState)
	}
}

func TestProgressModel_Update_ProgressUpdateMsg_MessageUpdate(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = WorkflowStateRunning

	updatedModel, _ := model.Update(ProgressUpdateMsg{
		Type:    MessageUpdate,
		Message: "Deploying",
	})

	m := updatedModel.(*ProgressModel)
	if m.currentMessage != "Deploying" {
		t.Errorf("message = %q, want %q", m.currentMessage, "Deploying")
	}
}

func TestProgressModel_Update_ProgressUpdateMsg_RejectedWhenCompleted(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = WorkflowStateCompleted

	updatedModel, _ := model.Update(ProgressUpdateMsg{
		Type:     ProgressUpdate,
		Progress: 50.0,
	})

	m := updatedModel.(*ProgressModel)
	if m.currentProgress == 50.0 {
		t.Error("progress should not update when workflow is completed")
	}
}

func TestProgressModel_Update_TickMsg_RejectedWhenNotRunning(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = WorkflowStateIdle

	updatedModel, _ := model.Update(TickMsg(time.Now()))

	m := updatedModel.(*ProgressModel)
	if m.workflowState != WorkflowStateIdle {
		t.Error("tick should not change state when idle")
	}
}

func TestProgressModel_GetActivityCounts(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	addTestActivity(model, "a", "A", nom.ActivityStatusRunning)
	addTestActivity(model, "b", "B", nom.ActivityStatusCompleted)

	counts := model.getActivityCounts()
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

func TestProgressModel_CtrlC_CancelsContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	model := newTestModel()
	model.cancelFunc = cancel

	_, cmd := model.Update(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})

	if ctx.Err() == nil {
		t.Error("ctrl+c should cancel the context")
	}

	if cmd == nil {
		t.Error("ctrl+c should also quit the TUI")
	}
}

func TestProgressModel_QuitKey_DoesNotCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	model := newTestModel()
	model.cancelFunc = cancel

	_, _ = model.Update(tea.KeyPressMsg{Code: 'q'})

	if ctx.Err() != nil {
		t.Error("q should not cancel the context")
	}
}
