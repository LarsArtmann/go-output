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
	dt.AddActivity(ActivityID("phase:build"), nil)
	dt.AddActivity(ActivityID("compile"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("test"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("lint"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("deploy"), []ActivityID{"phase:build"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("phase:build"), "Build Phase", ActivityStatusRunning, 5*time.Second)
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
	dt.AddActivity(ActivityID("phase:main"), nil)
	dt.AddActivity(ActivityID("step1"), []ActivityID{"phase:main"})
	dt.AddActivity(ActivityID("step2"), []ActivityID{"phase:main", "step1"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("phase:main"), "Phase", ActivityStatusRunning, 0)
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

	sub := newTestSubscriber(t)

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
