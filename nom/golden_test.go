package nom

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/golden"
)

func TestDependencyTreeRenderGolden_PhaseSteps(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("build"), nil)
	dt.AddActivity(ActivityID("compile"), []ActivityID{"build"})
	dt.AddActivity(ActivityID("test"), []ActivityID{"build"})
	dt.AddActivity(ActivityID("lint"), []ActivityID{"build"})
	dt.AddActivity(ActivityID("deploy"), []ActivityID{"build"})

	snaps := newSnapshotBuilder()
	snaps.setPhase(ActivityID("build"), "Build Phase", ActivityStatusRunning, 5*time.Second)
	snaps.set(ActivityID("compile"), "Compile", ActivityStatusCompleted, 2*time.Second)
	snaps.set(ActivityID("test"), "Run Tests", ActivityStatusRunning, 3*time.Second)
	snaps.set(ActivityID("lint"), "Lint Code", ActivityStatusPending, 0)
	snaps.set(ActivityID("deploy"), "Deploy", ActivityStatusFailed, 1*time.Second)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
	golden.RequireEqual(t, got)
}

func TestDependencyTreeRenderGolden_SecondaryDeps(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("main"), nil)
	dt.AddActivity(ActivityID("step1"), []ActivityID{"main"})
	dt.AddActivity(ActivityID("step2"), []ActivityID{"main", "step1"})

	snaps := newSnapshotBuilder()
	snaps.setPhase(ActivityID("main"), "Phase", ActivityStatusRunning, 0)
	snaps.set(ActivityID("step1"), "Step 1", ActivityStatusCompleted, 0)
	snaps.set(ActivityID("step2"), "Step 2", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
	golden.RequireEqual(t, got)
}

func TestDependencyTreeRenderGolden_MixedStates(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), nil)
	dt.AddActivity(ActivityID("running"), []ActivityID{"root"})
	dt.AddActivity(ActivityID("completed"), []ActivityID{"root"})
	dt.AddActivity(ActivityID("failed"), []ActivityID{"root"})
	dt.AddActivity(ActivityID("pending"), []ActivityID{"root"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Workflow", ActivityStatusRunning, 10*time.Second)
	snaps.set(ActivityID("running"), "Running Step", ActivityStatusRunning, 5*time.Second)
	snaps.set(ActivityID("completed"), "Completed Step", ActivityStatusCompleted, 2*time.Second)
	snaps.set(ActivityID("failed"), "Failed Step", ActivityStatusFailed, 1*time.Second)
	snaps.set(ActivityID("pending"), "Pending Step", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
	golden.RequireEqual(t, got)
}

func TestInlineRendererGolden_FirstFrame(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := newInlineTestRenderer(sub, &buf, 20)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	registerPhase(sub, ctx, ActivityID("build"), ActivityName("Build"))
	registerActivity(sub, ctx, ActivityID("compile"), ActivityName("Compile"), "build")
	registerActivity(sub, ctx, ActivityID("test"), ActivityName("Run Tests"), "build")
	registerActivity(sub, ctx, ActivityID("lint"), ActivityName("Lint"), "build")

	renderer.Draw()

	golden.RequireEqual(t, buf.String())
}

func TestInlineRendererGolden_SecondFrame(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := newInlineTestRenderer(sub, &buf, 20)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	registerPhase(sub, ctx, ActivityID("build"), ActivityName("Build"))
	registerActivity(sub, ctx, ActivityID("compile"), ActivityName("Compile"), "build")
	registerActivity(sub, ctx, ActivityID("lint"), ActivityName("Lint"), "build")

	renderer.Draw()
	buf.Reset()

	sendActivityStarted(t, sub, ctx, ActivityID("compile"), ActivityName("Compile"))

	renderer.Draw()

	golden.RequireEqual(t, buf.String())
}

func TestDependencyTreeRenderGolden_LayeredMode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.SetRenderMode(RenderModeLayered)

	_ = dt.AddActivity(ActivityID("setup"), nil)
	_ = dt.AddActivity(ActivityID("compile"), []ActivityID{"setup"})
	_ = dt.AddActivity(ActivityID("lint"), []ActivityID{"setup"})
	_ = dt.AddActivity(ActivityID("test"), []ActivityID{"compile"})
	_ = dt.AddActivity(ActivityID("deploy"), []ActivityID{"test", "lint"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("setup"), "Setup", ActivityStatusCompleted, 2*time.Second)
	snaps.set(ActivityID("compile"), "Compile", ActivityStatusCompleted, 5*time.Second)
	snaps.set(ActivityID("lint"), "Lint", ActivityStatusRunning, 3*time.Second)
	snaps.set(ActivityID("test"), "Test", ActivityStatusRunning, 1*time.Second)
	snaps.set(ActivityID("deploy"), "Deploy", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 20, 0)
	golden.RequireEqual(t, got)
}

func TestDependencyTreeRenderGolden_LayeredHeightPressure(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.SetRenderMode(RenderModeLayered)

	_ = dt.AddActivity(ActivityID("root"), nil)
	_ = dt.AddActivity(ActivityID("a"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("b"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("c"), []ActivityID{"a"})
	_ = dt.AddActivity(ActivityID("d"), []ActivityID{"b"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("a"), "A", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("b"), "B", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("c"), "C", ActivityStatusRunning, 0)
	snaps.set(ActivityID("d"), "D", ActivityStatusRunning, 0)

	// maxHeight=2 forces height-pressure collapse of completed layers.
	got := dt.RenderWithSnapshots(snaps.snaps, 2, 0)
	golden.RequireEqual(t, got)
}
