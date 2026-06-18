package nom

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/golden"
)

// setStatusWithElapsed updates activity status, timing, and derived visual
// style directly on the shared Activity pointer. The symbol/color params are
// ignored (derived from status via applyVisualStyle) but kept for call-site
// readability and backward compatibility with golden snapshots.
func setStatusWithElapsed(
	dt *DependencyTree,
	id ActivityID,
	status ActivityStatus, _, _ interface{},
	t time.Time, elapsed time.Duration,
) {
	node := dt.nodes[id]
	if node == nil {
		return
	}

	node.Status = status
	node.applyVisualStyle()
	node.StartTime = t

	if elapsed > 0 {
		node.CurrentElapsed = elapsed
	}
}

// TestDependencyTreeRenderGolden_PhaseSteps renders a tree with a phase node
// and multiple child steps in various states. Uses "phase:" prefix for
// phase-styled rendering (SymbolPhase/ColorPhase).
func TestDependencyTreeRenderGolden_PhaseSteps(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase:build"), NewActivity("phase:build", "Build Phase"), nil)
	dt.AddActivity(ActivityID("compile"), NewActivity("compile", "Compile"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("test"), NewActivity("test", "Run Tests"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("lint"), NewActivity("lint", "Lint Code"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("deploy"), NewActivity("deploy", "Deploy"), []ActivityID{"phase:build"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	setStatusWithElapsed(dt, ActivityID("phase:build"), ActivityStatusRunning, nil, nil, now, 5*time.Second)
	setStatusWithElapsed(dt, ActivityID("compile"), ActivityStatusCompleted, nil, nil, now, 2*time.Second)
	setStatusWithElapsed(dt, ActivityID("test"), ActivityStatusRunning, nil, nil, now, 3*time.Second)
	setStatusWithElapsed(dt, ActivityID("lint"), ActivityStatusPending, nil, nil, time.Time{}, 0)
	setStatusWithElapsed(dt, ActivityID("deploy"), ActivityStatusFailed, nil, nil, now, 1*time.Second)

	got := dt.RenderString(10)
	golden.RequireEqual(t, got)
}

// TestDependencyTreeRenderGolden_SecondaryDeps renders secondary dependency
// labels when an activity depends on multiple parents.
func TestDependencyTreeRenderGolden_SecondaryDeps(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase:main"), NewActivity("phase:main", "Phase"), nil)
	dt.AddActivity(ActivityID("step1"), NewActivity("step1", "Step 1"), []ActivityID{"phase:main"})
	dt.AddActivity(ActivityID("step2"), NewActivity("step2", "Step 2"), []ActivityID{"phase:main", "step1"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	setStatusWithElapsed(dt, ActivityID("phase:main"), ActivityStatusRunning, nil, nil, now, 0)
	setStatusWithElapsed(dt, ActivityID("step1"), ActivityStatusCompleted, nil, nil, now, 0)
	setStatusWithElapsed(dt, ActivityID("step2"), ActivityStatusPending, nil, nil, time.Time{}, 0)

	got := dt.RenderString(10)
	golden.RequireEqual(t, got)
}

// TestDependencyTreeRenderGolden_MixedStates renders activities in all
// possible states: running, completed, failed, pending.
func TestDependencyTreeRenderGolden_MixedStates(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), NewActivity("root", "Workflow"), nil)
	dt.AddActivity(ActivityID("running"), NewActivity("running", "Running Step"), []ActivityID{"root"})
	dt.AddActivity(ActivityID("completed"), NewActivity("completed", "Completed Step"), []ActivityID{"root"})
	dt.AddActivity(ActivityID("failed"), NewActivity("failed", "Failed Step"), []ActivityID{"root"})
	dt.AddActivity(ActivityID("pending"), NewActivity("pending", "Pending Step"), []ActivityID{"root"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	setStatusWithElapsed(dt, ActivityID("root"), ActivityStatusRunning, nil, nil, now, 10*time.Second)
	setStatusWithElapsed(dt, ActivityID("running"), ActivityStatusRunning, nil, nil, now, 5*time.Second)
	setStatusWithElapsed(dt, ActivityID("completed"), ActivityStatusCompleted, nil, nil, now, 2*time.Second)
	setStatusWithElapsed(dt, ActivityID("failed"), ActivityStatusFailed, nil, nil, now, 1*time.Second)
	setStatusWithElapsed(dt, ActivityID("pending"), ActivityStatusPending, nil, nil, time.Time{}, 0)

	got := dt.RenderString(10)
	golden.RequireEqual(t, got)
}

func TestInlineRendererGolden_FirstFrame(t *testing.T) {
	t.Parallel()

	sub := NewNOMStyleSubscriber()

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 20)
	renderer.SetNoColor(true)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	registerActivity(sub, ctx, ActivityID("phase:build"), ActivityName("Build"))
	registerActivity(sub, ctx, ActivityID("compile"), ActivityName("Compile"), "phase:build")
	registerActivity(sub, ctx, ActivityID("test"), ActivityName("Run Tests"), "phase:build")
	registerActivity(sub, ctx, ActivityID("lint"), ActivityName("Lint"), "phase:build")

	renderer.Draw()

	golden.RequireEqual(t, buf.String())
}

func TestInlineRendererGolden_SecondFrame(t *testing.T) {
	t.Parallel()

	sub := NewNOMStyleSubscriber()

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 20)
	renderer.SetNoColor(true)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	registerActivity(sub, ctx, ActivityID("phase:build"), ActivityName("Build"))
	registerActivity(sub, ctx, ActivityID("compile"), ActivityName("Compile"), "phase:build")
	registerActivity(sub, ctx, ActivityID("lint"), ActivityName("Lint"), "phase:build")

	renderer.Draw()
	buf.Reset()

	sendActivityStarted(t, sub, ctx, ActivityID("compile"), ActivityName("Compile"))

	renderer.Draw()

	golden.RequireEqual(t, buf.String())
}
