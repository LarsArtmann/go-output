package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/larsartmann/go-output/nom"
)

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
	model.dependencyTree.AddActivity(nom.ActivityID("a"), "Activity A", nil)
	model.dependencyTree.UpdateActivityStatus(
		nom.ActivityID("a"),
		nom.ActivityStatusRunning,
		nom.SymbolRunning,
		nom.ColorRunning,
		time.Now(),
		0,
	)

	output := model.renderNOMStyle()
	if output == "" {
		t.Error("renderNOMStyle() should not be empty")
	}
}

func TestProgressModel_RenderNOMSummaryBar(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.startTime = time.Now()
	model.activities[nom.ActivityID("a")] = nom.NewActivityDisplayState(
		nom.ActivityID("a"), nom.ActivityName("A"),
	)
	model.activities[nom.ActivityID("a")].SetRunning()

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
		model.dependencyTree.AddActivity(nom.ActivityID("a"), "Activity A", nil)
		model.dependencyTree.UpdateActivityStatus(
			nom.ActivityID("a"),
			nom.ActivityStatusRunning,
			nom.SymbolRunning,
			nom.ColorRunning,
			time.Now(),
			0,
		)

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
