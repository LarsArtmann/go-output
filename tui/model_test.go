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
	activity := nom.NewActivityDisplayState(nom.ActivityID(id), nom.ActivityName(name))
	activity.Status = status
	switch status {
	case nom.ActivityStatusRunning:
		activity.SetRunning()
	case nom.ActivityStatusCompleted:
		activity.SetCompleted()
	case nom.ActivityStatusFailed:
		activity.SetFailed(errors.New("test failure"))
	}

	model.nomSubscriber.SetActivityState(activity)
}

func setupTestTree(model *ProgressModel) *nom.DependencyTree {
	tree := nom.NewDependencyTree()
	model.dependencyTree = tree
	model.treeStartLine = 2

	return tree
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

func TestProgressModel_MouseClick_SelectsNode(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.displayMode = DisplayModeNOM

	tree := setupTestTree(model)
	_ = tree.AddActivity(nom.ActivityID("step-a"), "Step A", nil)
	_ = tree.AddActivity(nom.ActivityID("step-b"), "Step B", []nom.ActivityID{"step-a"})
	_ = tree.GetRootNodes()
	model.visibleNodes = tree.VisibleNodes(20)

	// Click on the first tree line (line 0 relative to tree = line 7 absolute with chrome)
	clickY := model.treeStartLine + chromeLinesAboveTree + 0
	updatedModel, _ := model.Update(tea.MouseClickMsg{
		X: 5, Y: clickY, Button: tea.MouseLeft,
	})

	m := updatedModel.(*ProgressModel)
	if m.selectedNode != nom.ActivityID("step-a") {
		t.Errorf("selectedNode = %q, want %q", m.selectedNode, "step-a")
	}
}

func TestProgressModel_MouseClick_ToggleOffNode(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.displayMode = DisplayModeNOM

	tree := setupTestTree(model)
	_ = tree.AddActivity(nom.ActivityID("step-a"), "Step A", nil)
	_ = tree.GetRootNodes()
	model.visibleNodes = tree.VisibleNodes(20)
	model.selectedNode = nom.ActivityID("step-a")

	clickY := model.treeStartLine + chromeLinesAboveTree + 0
	updatedModel, _ := model.Update(tea.MouseClickMsg{
		X: 5, Y: clickY, Button: tea.MouseLeft,
	})

	m := updatedModel.(*ProgressModel)
	if m.selectedNode != "" {
		t.Errorf("second click should deselect, got %q", m.selectedNode)
	}
}

func TestProgressModel_MouseClick_IgnoresRightClick(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.displayMode = DisplayModeNOM

	tree := nom.NewDependencyTree()
	_ = tree.AddActivity(nom.ActivityID("step-a"), "Step A", nil)
	_ = tree.GetRootNodes()

	model.dependencyTree = tree
	model.visibleNodes = tree.VisibleNodes(20)
	model.treeStartLine = 2

	updatedModel, _ := model.Update(tea.MouseClickMsg{
		X: 5, Y: 7, Button: tea.MouseRight,
	})

	m := updatedModel.(*ProgressModel)
	if m.selectedNode != "" {
		t.Error("right click should not select a node")
	}
}

func TestProgressModel_UpdateWorkflowCompletionState(t *testing.T) {
	t.Parallel()

	t.Run("failed activities transition to Errored", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		model.workflowState = WorkflowStateRunning
		model.displayMode = DisplayModeNOM

		ctx := context.Background()
		startWorkflow(t, model, ctx, nom.WorkflowID("wf-1"))
		startActivity(t, model, ctx, nom.ActivityID("test"), nom.ActivityName("Test"))
		_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
			eventType: nom.EventActivityFailed,
			aID:       nom.ActivityID("test"),
			aName:     nom.ActivityName("Test"),
			err:       errTestFail,
		})

		model.updateWorkflowCompletionState()

		if model.workflowState != WorkflowStateErrored {
			t.Errorf("state = %v, want Errored (failed > 0)", model.workflowState)
		}
	})

	t.Run("all completed transitions to Completed", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		model.workflowState = WorkflowStateRunning
		model.displayMode = DisplayModeNOM

		ctx := context.Background()
		startWorkflow(t, model, ctx, nom.WorkflowID("wf-1"))
		startActivity(t, model, ctx, nom.ActivityID("build"), nom.ActivityName("Build"))
		_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
			eventType: nom.EventActivityCompleted,
			aID:       nom.ActivityID("build"),
			aName:     nom.ActivityName("Build"),
			duration:  5 * time.Second,
		})

		model.updateWorkflowCompletionState()

		if model.workflowState != WorkflowStateCompleted {
			t.Errorf("state = %v, want Completed (running=0, completed>0)", model.workflowState)
		}
	})

	t.Run("still running stays Running", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		model.workflowState = WorkflowStateRunning
		model.displayMode = DisplayModeNOM

		ctx := context.Background()
		startWorkflow(t, model, ctx, nom.WorkflowID("wf-1"))
		startActivity(t, model, ctx, nom.ActivityID("active"), nom.ActivityName("Active"))

		model.updateWorkflowCompletionState()

		if model.workflowState != WorkflowStateRunning {
			t.Errorf("state = %v, want Running (still has running activities)", model.workflowState)
		}
	})
}

func TestBubbleTeaProgressReporter_Stop_NoProgramIsSafe(t *testing.T) {
	t.Parallel()

	pr := newTestReporter()

	// Stop() should be safe to call when no program was started (pr.program == nil)
	pr.Stop()

	// Verify the reporter is still usable
	if pr.Subscriber() == nil {
		t.Error("Subscriber() should still return a valid subscriber after Stop()")
	}
}
