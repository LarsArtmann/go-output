package nom

import (
	"image/color"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/golden"
)

// setStatusWithElapsed updates activity status and sets elapsed time if non-zero.
func setStatusWithElapsed(
	dt *DependencyTree,
	id ActivityID,
	status ActivityStatus, symbol string, c color.Color,
	t time.Time, elapsed time.Duration,
) {
	dt.UpdateActivityStatus(id, status, symbol, c, t, 0)
	if elapsed > 0 {
		dt.nodes[id].CurrentElapsed = elapsed
	}
}

// TestDependencyTreeRenderGolden_PhaseSteps renders a tree with a phase node
// and multiple child steps in various states. Uses "phase:" prefix for
// phase-styled rendering (SymbolPhase/ColorPhase).
func TestDependencyTreeRenderGolden_PhaseSteps(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase:build"), "Build Phase", nil)
	dt.AddActivity(ActivityID("compile"), "Compile", []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("test"), "Run Tests", []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("lint"), "Lint Code", []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("deploy"), "Deploy", []ActivityID{"phase:build"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	setStatusWithElapsed(dt, ActivityID("phase:build"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 5*time.Second)
	setStatusWithElapsed(dt, ActivityID("compile"), ActivityStatusCompleted, SymbolCompleted, ColorCompleted, now, 2*time.Second)
	setStatusWithElapsed(dt, ActivityID("test"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 3*time.Second)
	setStatusWithElapsed(dt, ActivityID("lint"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)
	setStatusWithElapsed(dt, ActivityID("deploy"), ActivityStatusFailed, SymbolFailed, ColorFailed, now, 1*time.Second)

	got := dt.Render(10)
	golden.RequireEqual(t, got)
}

// TestDependencyTreeRenderGolden_SecondaryDeps renders secondary dependency
// labels when an activity depends on multiple parents.
func TestDependencyTreeRenderGolden_SecondaryDeps(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase:main"), "Phase", nil)
	dt.AddActivity(ActivityID("step1"), "Step 1", []ActivityID{"phase:main"})
	dt.AddActivity(ActivityID("step2"), "Step 2", []ActivityID{"phase:main", "step1"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	setStatusWithElapsed(dt, ActivityID("phase:main"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)
	setStatusWithElapsed(dt, ActivityID("step1"), ActivityStatusCompleted, SymbolCompleted, ColorCompleted, now, 0)
	setStatusWithElapsed(dt, ActivityID("step2"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)

	got := dt.Render(10)
	golden.RequireEqual(t, got)
}

// TestDependencyTreeRenderGolden_MixedStates renders activities in all
// possible states: running, completed, failed, pending.
func TestDependencyTreeRenderGolden_MixedStates(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), "Workflow", nil)
	dt.AddActivity(ActivityID("running"), "Running Step", []ActivityID{"root"})
	dt.AddActivity(ActivityID("completed"), "Completed Step", []ActivityID{"root"})
	dt.AddActivity(ActivityID("failed"), "Failed Step", []ActivityID{"root"})
	dt.AddActivity(ActivityID("pending"), "Pending Step", []ActivityID{"root"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	setStatusWithElapsed(dt, ActivityID("root"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 10*time.Second)
	setStatusWithElapsed(dt, ActivityID("running"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 5*time.Second)
	setStatusWithElapsed(dt, ActivityID("completed"), ActivityStatusCompleted, SymbolCompleted, ColorCompleted, now, 2*time.Second)
	setStatusWithElapsed(dt, ActivityID("failed"), ActivityStatusFailed, SymbolFailed, ColorFailed, now, 1*time.Second)
	setStatusWithElapsed(dt, ActivityID("pending"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)

	got := dt.Render(10)
	golden.RequireEqual(t, got)
}
