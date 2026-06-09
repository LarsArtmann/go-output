package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/larsartmann/go-output/nom"
)

func newTestModel() *ProgressModel {
	return &ProgressModel{
		messages:       make([]string, 0),
		steps:          make([]ProgressStep, 0),
		startTime:      time.Now(),
		lastUpdate:     time.Now(),
		workflowState:  WorkflowStateIdle,
		displayMode:    DisplayModeUniversal,
		activities:     make(map[nom.ActivityID]*nom.ActivityDisplayState),
		dependencyTree: nom.NewDependencyTree(),
		timingCache:    nom.NewTimingCache(),
		nomSubscriber:  nom.NewNOMStyleSubscriber(),
	}
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
	model.activities[nom.ActivityID("a")] = nom.NewActivityDisplayState(
		nom.ActivityID("a"), nom.ActivityName("A"),
	)
	model.activities[nom.ActivityID("b")] = nom.NewActivityDisplayState(
		nom.ActivityID("b"), nom.ActivityName("B"),
	)
	model.activities[nom.ActivityID("a")].SetRunning()
	model.activities[nom.ActivityID("b")].SetCompleted()

	running, completed, failed, pending := model.getActivityCounts()
	if running != 1 {
		t.Errorf("running = %d, want 1", running)
	}

	if completed != 1 {
		t.Errorf("completed = %d, want 1", completed)
	}

	if failed != 0 {
		t.Errorf("failed = %d, want 0", failed)
	}

	if pending != 0 {
		t.Errorf("pending = %d, want 0", pending)
	}
}
