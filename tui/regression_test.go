package tui

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output/nom"
)

// TestRenderProgressBar_NarrowWidth does not panic when terminal width is very small.
// Regression for a bug where strings.Repeat received a negative count.
func TestRenderProgressBar_NarrowWidth(t *testing.T) {
	t.Parallel()

	widths := []int{1, 5, 10, 20, 29}

	for _, w := range widths {
		t.Run("narrow_width", func(t *testing.T) {
			t.Parallel()

			model := newTestModel()
			model.width = w
			model.currentProgress = 50.0

			// Must not panic.
			_ = model.renderProgressBar()
		})
	}
}

// TestHandleStepUpdate_TotalZeroNotCompleted ensures that a step with Total=0
// is NOT immediately marked as completed. Regression for an unsigned-comparison
// flaw where Current >= Total was always true when Total=0.
func TestHandleStepUpdate_TotalZeroNotCompleted(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = workflowStateRunning

	_, _ = model.handleStepUpdate(stepUpdateMsg{
		Message: "Starting...",
		Current: 0,
		Total:   0,
	})

	if len(model.steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(model.steps))
	}

	if model.steps[0].CompletedAt != nil {
		t.Error("step with Total=0 should not be marked completed")
	}
}

// TestHandleStepUpdate_MatchesByMessageNotActive ensures that a step update
// targets the step with the matching message, not any arbitrary active step.
// Regression for a bug where || m.steps[i].isActive() matched the wrong step.
func TestHandleStepUpdate_MatchesByMessageNotActive(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.workflowState = workflowStateRunning

	// Add two steps with different messages.
	model.steps = []progressStep{
		{Message: "Step A", Current: 0, Total: 10},
		{Message: "Step B", Current: 0, Total: 10},
	}

	// Update Step B specifically.
	_, _ = model.handleStepUpdate(stepUpdateMsg{
		Message: "Step B",
		Current: 5,
		Total:   10,
	})

	if model.steps[0].Current != 0 {
		t.Errorf("Step A should be unchanged, got Current=%d", model.steps[0].Current)
	}

	if model.steps[1].Current != 5 {
		t.Errorf("Step B should have Current=5, got %d", model.steps[1].Current)
	}
}

// TestStateSummary_NoDoubleSuffix ensures that elapsed-time templates don't
// produce double unit suffixes (e.g. "1.5ss"). Regression for a bug where
// {time}s was appended after a string that already contained the unit.
func TestStateSummary_NoDoubleSuffix(t *testing.T) {
	t.Parallel()

	for _, state := range []workflowState{workflowStateIdle, workflowStateCompleted, workflowStateErrored} {
		t.Run(state.String(), func(t *testing.T) {
			t.Parallel()

			stateSummary, _ := getStateStyle(state, nom.ThemeDefault)

			// The template must NOT have "{time}s" — formatElapsedTime already
			// returns unit-suffixed strings like "1.5s", "500ms", "2m30s".
			if strings.Contains(stateSummary, "{time}s") {
				t.Errorf("template for %s has double-suffix bug: %q", state, stateSummary)
			}
		})
	}
}
