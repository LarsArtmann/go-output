package nom

import (
	"testing"
	"time"
)

func TestSetEnabledAndIsEnabled(t *testing.T) {
	t.Parallel()

	t.Run("default state is enabled", func(t *testing.T) {
		t.Parallel()

		ns := newTestSubscriber(t)
		if !ns.IsEnabled() {
			t.Error("default subscriber should be enabled")
		}
	})

	t.Run("disable sets IsEnabled to false", func(t *testing.T) {
		t.Parallel()

		ns := newTestSubscriber(t)
		ns.SetEnabled(false)

		if ns.IsEnabled() {
			t.Error("IsEnabled() should be false after SetEnabled(false)")
		}
	})

	t.Run("re-enable sets IsEnabled to true", func(t *testing.T) {
		t.Parallel()

		ns := newTestSubscriber(t)
		ns.SetEnabled(false)
		ns.SetEnabled(true)

		if !ns.IsEnabled() {
			t.Error("IsEnabled() should be true after SetEnabled(true)")
		}
	})
}

func TestReset_ClearsAllState(t *testing.T) {
	t.Parallel()

	ns := newTestSubscriber(t)

	// Populate state.
	ns.activities = map[ActivityID]*Activity{
		NewActivityID("a1"): NewActivity("a1", "Activity 1"),
	}
	ns.workflowID = NewWorkflowID("wf-1")
	ns.workflowName = NewWorkflowName("Workflow 1")
	ns.startTime = time.Now()
	ns.isRunning = true
	ns.dependencyTree.AddActivity(NewActivityID("a1"), nil)

	ns.Reset()

	if len(ns.activities) != 0 {
		t.Errorf("activities should be empty after Reset, got %d", len(ns.activities))
	}

	if ns.workflowID != "" {
		t.Errorf("workflowID = %q, want empty", ns.workflowID)
	}

	if ns.workflowName != WorkflowName("") {
		t.Errorf("workflowName = %q, want empty", ns.workflowName)
	}

	if !ns.startTime.IsZero() {
		t.Errorf("startTime = %v, want zero", ns.startTime)
	}

	if ns.isRunning {
		t.Error("isRunning should be false after Reset")
	}

	if ns.dependencyTree == nil {
		t.Error("dependencyTree should not be nil after Reset (only cleared)")
	}
}
