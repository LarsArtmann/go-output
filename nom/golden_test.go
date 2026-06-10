package nom

import (
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/golden"
)

// TestDependencyTreeRenderGolden_PhaseSteps renders a tree with a phase node
// and multiple child steps in various states.
func TestDependencyTreeRenderGolden_PhaseSteps(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase1"), "Build Phase", nil)
	dt.AddActivity(ActivityID("compile"), "Compile", []ActivityID{"phase1"})
	dt.AddActivity(ActivityID("test"), "Run Tests", []ActivityID{"phase1"})
	dt.AddActivity(ActivityID("lint"), "Lint Code", []ActivityID{"phase1"})
	dt.AddActivity(ActivityID("deploy"), "Deploy", []ActivityID{"phase1"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	dt.UpdateActivityStatus(ActivityID("phase1"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)
	dt.nodes[ActivityID("phase1")].CurrentElapsed = 5 * time.Second

	dt.UpdateActivityStatus(ActivityID("compile"), ActivityStatusCompleted, SymbolCompleted, ColorCompleted, now, 0)
	dt.nodes[ActivityID("compile")].CurrentElapsed = 2 * time.Second

	dt.UpdateActivityStatus(ActivityID("test"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)
	dt.nodes[ActivityID("test")].CurrentElapsed = 3 * time.Second

	dt.UpdateActivityStatus(ActivityID("lint"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)

	dt.UpdateActivityStatus(ActivityID("deploy"), ActivityStatusFailed, SymbolFailed, ColorFailed, now, 0)
	dt.nodes[ActivityID("deploy")].CurrentElapsed = 1 * time.Second

	got := dt.Render(10)
	golden.RequireEqual(t, got)
}

// TestDependencyTreeRenderGolden_SecondaryDeps renders secondary dependency
// labels when an activity depends on multiple parents.
func TestDependencyTreeRenderGolden_SecondaryDeps(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), "Phase", nil)
	dt.AddActivity(ActivityID("step1"), "Step 1", []ActivityID{"phase"})
	dt.AddActivity(ActivityID("step2"), "Step 2", []ActivityID{"phase", "step1"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	dt.UpdateActivityStatus(ActivityID("phase"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)
	dt.UpdateActivityStatus(ActivityID("step1"), ActivityStatusCompleted, SymbolCompleted, ColorCompleted, now, 0)
	dt.UpdateActivityStatus(ActivityID("step2"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)

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

	dt.UpdateActivityStatus(ActivityID("root"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)
	dt.nodes[ActivityID("root")].CurrentElapsed = 10 * time.Second

	dt.UpdateActivityStatus(ActivityID("running"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)
	dt.nodes[ActivityID("running")].CurrentElapsed = 5 * time.Second

	dt.UpdateActivityStatus(ActivityID("completed"), ActivityStatusCompleted, SymbolCompleted, ColorCompleted, now, 0)
	dt.nodes[ActivityID("completed")].CurrentElapsed = 2 * time.Second

	dt.UpdateActivityStatus(ActivityID("failed"), ActivityStatusFailed, SymbolFailed, ColorFailed, now, 0)
	dt.nodes[ActivityID("failed")].CurrentElapsed = 1 * time.Second

	dt.UpdateActivityStatus(ActivityID("pending"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)

	got := dt.Render(10)
	golden.RequireEqual(t, got)
}
