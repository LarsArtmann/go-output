package tui

import (
	"context"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/larsartmann/go-output/nom"
)

func TestProgressModel_MouseClick_SelectsNode(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.displayMode = DisplayModeNOM

	tree := setupTestTree(model)
	_ = tree.AddActivity(nom.ActivityID("step-a"), nil)
	_ = tree.AddActivity(nom.ActivityID("step-b"), []nom.ActivityID{"step-a"})
	_ = tree.GetRootNodes()
	model.visibleNodes = tree.VisibleNodesWithSnapshots(nil, 20)

	// Click on the first tree line (line 0 relative to tree = line 7 absolute with chrome)
	clickY := model.treeStartLine + chromeLinesAboveTree + 0

	m := clickAt(model, clickY)
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
	_ = tree.AddActivity(nom.ActivityID("step-a"), nil)
	_ = tree.GetRootNodes()
	model.visibleNodes = tree.VisibleNodesWithSnapshots(nil, 20)
	model.selectedNode = nom.ActivityID("step-a")

	clickY := model.treeStartLine + chromeLinesAboveTree + 0

	m := clickAt(model, clickY)
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
	_ = tree.AddActivity(nom.ActivityID("step-a"), nil)
	_ = tree.GetRootNodes()

	model.dependencyTree = tree
	model.visibleNodes = tree.VisibleNodesWithSnapshots(nil, 20)
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
		model.workflowState = workflowStateRunning
		model.displayMode = DisplayModeNOM

		ctx := context.Background()
		startWorkflow(t, model, ctx, nom.WorkflowID("wf-1"))
		startActivity(t, model, ctx, nom.ActivityID("test"), nom.ActivityName("Test"))
		_ = model.nomSubscriber.OnEvent(ctx, nom.ActivityFailed{
			ID:   nom.ActivityID("test"),
			Name: nom.ActivityName("Test"),
			Err:  errTestFail,
		})

		model.updateWorkflowCompletionState()

		if model.workflowState != workflowStateErrored {
			t.Errorf("state = %v, want Errored (failed > 0)", model.workflowState)
		}
	})

	t.Run("all completed transitions to Completed", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		model.workflowState = workflowStateRunning
		model.displayMode = DisplayModeNOM

		ctx := context.Background()
		startWorkflow(t, model, ctx, nom.WorkflowID("wf-1"))
		startActivity(t, model, ctx, nom.ActivityID("build"), nom.ActivityName("Build"))
		_ = model.nomSubscriber.OnEvent(ctx, nom.ActivityCompleted{
			ID:       nom.ActivityID("build"),
			Name:     nom.ActivityName("Build"),
			Duration: 5 * time.Second,
		})

		model.updateWorkflowCompletionState()

		if model.workflowState != workflowStateCompleted {
			t.Errorf("state = %v, want Completed (running=0, completed>0)", model.workflowState)
		}
	})

	t.Run("still running stays Running", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		model.workflowState = workflowStateRunning
		model.displayMode = DisplayModeNOM

		ctx := context.Background()
		startWorkflow(t, model, ctx, nom.WorkflowID("wf-1"))
		startActivity(t, model, ctx, nom.ActivityID("active"), nom.ActivityName("Active"))

		model.updateWorkflowCompletionState()

		if model.workflowState != workflowStateRunning {
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
		t.Error("subscriber() should still return a valid subscriber after Stop()")
	}
}
