package tui

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/larsartmann/go-output/nom"
)

func addRunningActivity(model *ProgressModel, id, name string) {
	model.dependencyTree.AddActivity(nom.ActivityID(id), name, nil)
	model.dependencyTree.UpdateActivityStatus(
		nom.ActivityID(id),
		nom.ActivityStatusRunning,
		nom.SymbolRunning,
		nom.ColorRunning,
		time.Now(),
		0,
	)
}

func TestProgressModel_View_ZeroWidth(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	view := model.View()

	content := view.Content
	if content != "Loading..." {
		t.Errorf("View() with zero width = %q, want %q", content, "Loading...")
	}
}

func TestProgressModel_View_UniversalMode(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.workflowState = WorkflowStateRunning
	model.currentProgress = 50.0
	model.currentMessage = "Building"

	view := model.View()

	content := view.Content
	if content == "" {
		t.Error("View() should not be empty")
	}

	if !view.AltScreen {
		t.Error("View() should use alt screen")
	}
}

func TestProgressModel_View_NOMMode(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.workflowState = WorkflowStateRunning
	model.displayMode = DisplayModeNOM

	view := model.View()

	content := view.Content
	if content == "" {
		t.Error("View() should not be empty in NOM mode")
	}
}

func TestProgressModel_RenderSteps(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	now := time.Now()
	model.steps = []ProgressStep{
		{Message: "Step 1", Current: 5, Total: 5, CompletedAt: &now, IsActive: false},
		{Message: "Step 2", Current: 3, Total: 5, StartTime: now, IsActive: true},
	}

	output := model.renderSteps()
	if output == "" {
		t.Error("renderSteps() should not be empty")
	}
}

func TestProgressModel_RenderProgressBar(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.currentProgress = 50.0

	output := model.renderProgressBar()
	if output == "" {
		t.Error("renderProgressBar() should not be empty")
	}
}

func TestProgressModel_RenderSummaryBar(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.startTime = time.Now()

	output := model.renderSummaryBar()
	if output == "" {
		t.Error("renderSummaryBar() should not be empty")
	}
}

func TestProgressModel_RenderSummaryBar_Completed(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.startTime = time.Now()
	model.workflowState = WorkflowStateCompleted
	now := time.Now()
	model.steps = []ProgressStep{
		{Message: "Step 1", CompletedAt: &now},
	}

	output := model.renderSummaryBar()
	if output == "" {
		t.Error("renderSummaryBar() should not be empty when completed")
	}
}

func TestProgressModel_RenderNOMStyle(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.workflowState = WorkflowStateRunning
	model.displayMode = DisplayModeNOM
	model.dependencyTree = nom.NewDependencyTree()
	addRunningActivity(model, "a", "Activity A")

	output := model.renderNOMStyle()
	if output == "" {
		t.Error("renderNOMStyle() should not be empty")
	}
}

func TestProgressModel_RenderNOMSummaryBar(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.startTime = time.Now()
	addTestActivity(model, "a", "A", func(a *nom.ActivityDisplayState) { a.SetRunning() })

	output := model.renderNOMSummaryBar()
	if output == "" {
		t.Error("renderNOMSummaryBar() should not be empty")
	}
}

func TestProgressModel_RenderDependencyTree(t *testing.T) {
	t.Parallel()

	t.Run("nil tree returns empty", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		model.dependencyTree = nil

		got := model.renderDependencyTree()
		if got != "" {
			t.Errorf("expected empty for nil tree, got %q", got)
		}
	})

	t.Run("tree with activities renders", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		model.height = 24
		addRunningActivity(model, "a", "Activity A")

		got := model.renderDependencyTree()
		if got == "" {
			t.Error("renderDependencyTree() should not be empty for non-empty tree")
		}
	})
}

func TestProgressModel_HelpOverlay(t *testing.T) {
	t.Parallel()

	t.Run("toggle on and off", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		if model.showHelp {
			t.Error("showHelp should start false")
		}

		updatedModel, _ := model.Update(tea.KeyPressMsg{Code: '?'})
		m := updatedModel.(*ProgressModel)
		if !m.showHelp {
			t.Error("showHelp should be true after pressing ?")
		}

		updatedModel2, _ := m.Update(tea.KeyPressMsg{Code: '?'})
		m2 := updatedModel2.(*ProgressModel)
		if m2.showHelp {
			t.Error("showHelp should be false after pressing ? again")
		}
	})

	t.Run("render includes help content", func(t *testing.T) {
		t.Parallel()

		model := newTestModel()
		model.width = 80
		model.height = 24
		model.workflowState = WorkflowStateRunning
		model.showHelp = true

		view := model.View()
		if view.Content == "" {
			t.Error("View() should not be empty with help overlay")
		}
	})
}

func TestProgressModel_SelectedNodeHighlight(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.displayMode = DisplayModeNOM
	model.workflowState = WorkflowStateRunning

	tree := nom.NewDependencyTree()
	_ = tree.AddActivity(nom.ActivityID("step-a"), "Step A", nil)
	_ = tree.AddActivity(nom.ActivityID("step-b"), "Step B", []nom.ActivityID{"step-a"})
	tree.EnsureBuild()

	model.dependencyTree = tree
	model.selectedNode = nom.ActivityID("step-a")

	view := model.View()
	content := view.Content
	if content == "" {
		t.Fatal("View() should not be empty")
	}

	for _, want := range []string{"Step A", "Step B"} {
		if !strings.Contains(content, want) {
			t.Errorf("rendered content should contain %q", want)
		}
	}
}

func TestProgressModel_GetActivityCounts_AllStatuses(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	addTestActivity(model, "running", "R", func(a *nom.ActivityDisplayState) { a.SetRunning() })
	addTestActivity(model, "completed", "C", func(a *nom.ActivityDisplayState) { a.SetCompleted() })
	addTestActivity(model, "failed", "F", func(a *nom.ActivityDisplayState) { a.SetFailed(errors.New("fail")) })

	paused := nom.NewActivityDisplayState(nom.ActivityID("paused"), nom.ActivityName("P"))
	paused.Status = nom.ActivityStatusPaused
	model.activities[nom.ActivityID("paused")] = paused

	running, completed, failed, pending := model.getActivityCounts()
	if running != 1 {
		t.Errorf("running = %d, want 1", running)
	}

	if completed != 1 {
		t.Errorf("completed = %d, want 1", completed)
	}

	if failed != 1 {
		t.Errorf("failed = %d, want 1", failed)
	}

	if pending != 1 {
		t.Errorf("pending = %d, want 1 (paused counts as pending)", pending)
	}
}

func TestProgressModel_RenderHelpOverlay_DefaultDimensions(t *testing.T) {
	t.Parallel()

	model := newTestModel()

	got := model.renderHelpOverlay()
	if got == "" {
		t.Error("renderHelpOverlay() should not be empty with default dimensions")
	}
}

func TestProgressModel_RenderUniversalWorkflowProgress_NoStepsNoProgress(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.width = 80
	model.height = 24
	model.workflowState = WorkflowStateRunning

	got := model.renderUniversalWorkflowProgress()
	if got == "" {
		t.Error("should render even with no steps and no progress")
	}
}
