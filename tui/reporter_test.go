package tui

import (
	"testing"
	"time"
)

func TestNewBubbleTeaProgressReporter(t *testing.T) {
	t.Parallel()

	reporter := NewBubbleTeaProgressReporter()
	if reporter == nil {
		t.Fatal("NewBubbleTeaProgressReporter() returned nil")
	}

	if reporter.started {
		t.Error("reporter should not be started initially")
	}
}

func TestBubbleTeaProgressReporter_TransitionWorkflowState(t *testing.T) {
	t.Parallel()

	reporter := NewBubbleTeaProgressReporter()

	t.Run("idle to running is valid", func(t *testing.T) {
		t.Parallel()

		ok := reporter.transitionWorkflowState(workflowStateRunning)
		if !ok {
			t.Error("transition from idle to running should succeed")
		}
	})

	t.Run("running to completed is valid", func(t *testing.T) {
		t.Parallel()

		r := NewBubbleTeaProgressReporter()
		r.transitionWorkflowState(workflowStateRunning)

		ok := r.transitionWorkflowState(workflowStateCompleted)
		if !ok {
			t.Error("transition from running to completed should succeed")
		}
	})

	t.Run("idle to completed is invalid", func(t *testing.T) {
		t.Parallel()

		r := NewBubbleTeaProgressReporter()

		ok := r.transitionWorkflowState(workflowStateCompleted)
		if ok {
			t.Error("transition from idle to completed should fail")
		}
	})
}

func TestBubbleTeaProgressReporter_IsWorkflowActive(t *testing.T) {
	t.Parallel()

	t.Run("idle is active", func(t *testing.T) {
		t.Parallel()

		reporter := NewBubbleTeaProgressReporter()
		if !reporter.isWorkflowActive() {
			t.Error("idle state should be active")
		}
	})

	t.Run("completed is not active", func(t *testing.T) {
		t.Parallel()

		reporter := NewBubbleTeaProgressReporter()
		reporter.transitionWorkflowState(workflowStateRunning)
		reporter.transitionWorkflowState(workflowStateCompleted)

		if reporter.isWorkflowActive() {
			t.Error("completed state should not be active")
		}
	})
}

func TestBubbleTeaProgressReporter_ReportProgress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		progress  float64
		wantState workflowState
	}{
		{"sets running on first report", 50.0, workflowStateRunning},
		{"completes at 100%", 100.0, workflowStateCompleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reporter := newTestReporter()
			reporter.ReportProgress(tt.progress)

			if reporter.model.workflowState != tt.wantState {
				t.Errorf("workflow state = %v, want %v", reporter.model.workflowState, tt.wantState)
			}
		})
	}
}

func TestBubbleTeaProgressReporter_ReportMessage(t *testing.T) {
	t.Parallel()

	reporter := newTestReporter()
	reporter.ReportMessage("Building project")

	if reporter.model.currentMessage != "Building project" {
		t.Errorf("message = %q, want %q", reporter.model.currentMessage, "Building project")
	}

	if reporter.model.workflowState != workflowStateRunning {
		t.Errorf("workflow state = %v, want Running", reporter.model.workflowState)
	}
}

func TestBubbleTeaProgressReporter_ReportStep(t *testing.T) {
	t.Parallel()

	t.Run("creates active step", func(t *testing.T) {
		t.Parallel()

		reporter := newTestReporter()
		reporter.ReportStep(1, 5, "Compile")

		assertSingleStep(t, reporter)

		if reporter.model.steps[0].Message != "Compile" {
			t.Errorf("step message = %q, want %q", reporter.model.steps[0].Message, "Compile")
		}

		assertStepCurrent(t, reporter, 1)

		if !reporter.model.steps[0].isActive() {
			t.Error("step should be active (1 < 5)")
		}
	})

	t.Run("completes when current >= total", func(t *testing.T) {
		t.Parallel()

		reporter := newTestReporter()
		reporter.ReportStep(1, 5, "Compile")
		reporter.ReportStep(5, 5, "Compile")

		if reporter.model.steps[0].CompletedAt == nil {
			t.Error("step should be completed when updated to current >= total")
		}

		if reporter.model.steps[0].isActive() {
			t.Error("step should not be active when current >= total")
		}
	})

	t.Run("updates existing step", func(t *testing.T) {
		t.Parallel()

		reporter := newTestReporter()
		reporter.ReportStep(1, 5, "Compile")
		reporter.ReportStep(3, 5, "Compile")

		if len(reporter.model.steps) != 1 {
			t.Errorf("steps count = %d, want 1 (should update, not add)", len(reporter.model.steps))
		}

		assertStepCurrent(t, reporter, 3)
	})
}

func TestBubbleTeaProgressReporter_ReportError(t *testing.T) {
	t.Parallel()

	reporter := newTestReporter()
	reporter.ReportProgress(50.0)
	reporter.ReportError(errDiskFull)

	if reporter.model.workflowState != workflowStateErrored {
		t.Errorf("workflow state = %v, want Errored", reporter.model.workflowState)
	}
}

func TestBubbleTeaProgressReporter_ReportError_CannotTransitionFromIdle(t *testing.T) {
	t.Parallel()

	reporter := newTestReporter()
	reporter.ReportError(errTestFail)

	if reporter.model.workflowState != workflowStateIdle {
		t.Errorf("cannot transition from idle to errored, state = %v", reporter.model.workflowState)
	}
}

func TestBubbleTeaProgressReporter_Send_NilProgram(t *testing.T) {
	t.Parallel()

	reporter := newTestReporter()
	reporter.send(progressUpdateMsg{
		Type:     progressUpdate,
		Progress: 50.0,
	})
}

func TestProgressModel_Init(t *testing.T) {
	t.Parallel()

	model := &ProgressModel{
		steps:         make([]progressStep, 0),
		startTime:     time.Now(),
		workflowState: workflowStateIdle,
		displayMode:   DisplayModeUniversal,
	}

	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return a non-nil command")
	}
}
