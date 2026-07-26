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
	case nom.ActivityStatusPending:
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

// setupLayeredTestModel returns a model wired with a two-level (root + child)
// dependency tree, both activities registered, and layered render mode active.
// Tests that only need a root pass withChild=false.
func setupLayeredTestModel(t *testing.T, withChild bool) (*ProgressModel, *nom.DependencyTree) {
	t.Helper()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.displayMode = DisplayModeNOM

	tree := setupTestTree(model)
	_ = tree.AddActivity(nom.ActivityID("root"), nil)

	addTestActivity(model, "root", "Root", nom.ActivityStatusRunning)

	if withChild {
		_ = tree.AddActivity(nom.ActivityID("child"), []nom.ActivityID{"root"})

		addTestActivity(model, "child", "Child", nom.ActivityStatusCompleted)
	}

	_ = tree.GetRootNodes()
	tree.SetRenderMode(nom.RenderModeLayered)

	return model, tree
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
	model.workflowState = workflowStateRunning

	updatedModel, _ := model.Update(progressUpdateMsg{
		Type:     progressUpdate,
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
	model.workflowState = workflowStateRunning

	updatedModel, _ := model.Update(progressUpdateMsg{
		Type:     progressUpdate,
		Progress: 100.0,
	})

	m := updatedModel.(*ProgressModel)
	if m.workflowState != workflowStateCompleted {
		t.Errorf("workflow state = %v, want Completed", m.workflowState)
	}
}

func TestProgressModel_Update_ProgressUpdateMsg_MessageUpdate(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = workflowStateRunning

	updatedModel, _ := model.Update(progressUpdateMsg{
		Type:    messageUpdate,
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
	model.workflowState = workflowStateCompleted

	updatedModel, _ := model.Update(progressUpdateMsg{
		Type:     progressUpdate,
		Progress: 50.0,
	})

	m := updatedModel.(*ProgressModel)
	if m.currentProgress == 50.0 {
		t.Error("progress should not update when workflow is completed")
	}
}

func TestProgressModel_AcceptedUpdatesStampLastUpdate(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = workflowStateRunning
	before := model.lastUpdate

	model.Update(progressUpdateMsg{Type: messageUpdate, Message: "Building"})

	progressUpdate := model.lastUpdate
	if !progressUpdate.After(before) {
		t.Error("progress update should advance lastUpdate")
	}

	model.Update(stepUpdateMsg{Current: 1, Total: 2, Message: "Compile"})

	if model.lastUpdate.Before(progressUpdate) {
		t.Error("step update should not move lastUpdate backwards")
	}
}

func TestProgressModel_Update_TickMsg_RejectedWhenNotRunning(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = workflowStateIdle

	updatedModel, _ := model.Update(tickMsg(time.Now()))

	m := updatedModel.(*ProgressModel)
	if m.workflowState != workflowStateIdle {
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
