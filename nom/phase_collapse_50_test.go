package nom

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestPhaseCollapse_50StepsAcross5Categories(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.collapseCompletedPhases = true

	const (
		nPhases       = 5
		stepsPerPhase = 10
	)

	snaps := newSnapshotBuilder()

	for p := range nPhases {
		phaseID := ActivityID(fmt.Sprintf("phase-%d", p))
		dt.AddActivity(phaseID, nil)
		snaps.setPhase(phaseID, fmt.Sprintf("Phase-%d", p), ActivityStatusCompleted, 0)

		for s := range stepsPerPhase {
			stepID := ActivityID(fmt.Sprintf("step-%d-%d", p, s))
			dt.AddActivity(stepID, []ActivityID{phaseID})
			snaps.set(stepID, fmt.Sprintf("Step-%d-%d", p, s), ActivityStatusCompleted, 100*time.Millisecond)
		}
	}

	// Render with enough height to show collapsed phases.
	got := dt.RenderWithSnapshots(snaps.snaps, 100, 0)

	if got == "" {
		t.Fatal("expected non-empty output")
	}

	for p := range nPhases {
		phaseLabel := fmt.Sprintf("Phase-%d", p)
		if !strings.Contains(got, phaseLabel) {
			t.Errorf("expected phase label %q in output, got:\n%s", phaseLabel, got)
		}

		// Each phase should show 10/10 in the collapsed label.
		expected := fmt.Sprintf("10/%d", stepsPerPhase)
		if !strings.Contains(got, expected) {
			t.Errorf("expected %q in collapsed label for Phase-%d, got:\n%s", expected, p, got)
		}
	}

	// Individual step labels should NOT appear (they're collapsed).
	if strings.Contains(got, "Step-0-0") {
		t.Errorf("expected step labels hidden by collapse, but found Step-0-0 in:\n%s", got)
	}

	// The output should be much shorter than 50 lines (collapsed = ~5 phase lines).
	lineCount := strings.Count(got, "\n")
	if lineCount > 30 {
		t.Errorf("expected collapsed output to be compact (<30 lines), got %d lines", lineCount)
	}
}
