package nom

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLayeredRender_BasicGrouping(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root1"), nil)
	_ = dt.AddActivity(ActivityID("root2"), nil)
	_ = dt.AddActivity(ActivityID("child1"), []ActivityID{"root1"})
	_ = dt.AddActivity(ActivityID("child2"), []ActivityID{"root1", "root2"})
	_ = dt.AddActivity(ActivityID("grandchild1"), []ActivityID{"child1"})

	dt.SetRenderMode(RenderModeLayered)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root1"), "Root One", ActivityStatusRunning, 2*time.Second)
	snaps.set(ActivityID("root2"), "Root Two", ActivityStatusPending, 0)
	snaps.set(ActivityID("child1"), "Child One", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("child2"), "Child Two", ActivityStatusRunning, 3*time.Second)
	snaps.set(ActivityID("grandchild1"), "Grandchild One", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)
	lines := strings.Split(got, "\n")

	if len(lines) < 7 {
		t.Fatalf("expected at least 7 lines, got %d:\n%s", len(lines), got)
	}

	wantContains := []string{
		"Layer 0",
		"Root One",
		"Root Two",
		"Layer 1",
		"Child One",
		"Child Two",
		"Layer 2",
		"Grandchild One",
	}

	for _, want := range wantContains {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q:\n%s", want, got)
		}
	}

	layer0 := strings.Index(got, "Layer 0")
	layer1 := strings.Index(got, "Layer 1")
	layer2 := strings.Index(got, "Layer 2")

	if layer0 >= layer1 || layer1 >= layer2 {
		t.Errorf("layers out of order: layer0=%d layer1=%d layer2=%d", layer0, layer1, layer2)
	}
}

func TestLayeredRender_SkipsPlaceholders(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("child"), []ActivityID{"missing-parent"})
	dt.SetRenderMode(RenderModeLayered)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("child"), "Child", ActivityStatusRunning, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)

	if strings.Contains(got, "missing-parent") {
		t.Errorf("placeholder node should not be rendered:\n%s", got)
	}

	if !strings.Contains(got, "Child") {
		t.Errorf("real child should be rendered:\n%s", got)
	}
}

func TestLayeredRender_PrioritySort(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("completed"), nil)
	_ = dt.AddActivity(ActivityID("running"), nil)
	_ = dt.AddActivity(ActivityID("pending"), nil)

	dt.SetRenderMode(RenderModeLayered)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("completed"), "Completed", ActivityStatusCompleted, 0)
	snaps.set(ActivityID("running"), "Running", ActivityStatusRunning, 0)
	snaps.set(ActivityID("pending"), "Pending", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)

	runningIdx := strings.Index(got, "Running")
	pendingIdx := strings.Index(got, "Pending")
	completedIdx := strings.Index(got, "Completed")

	if runningIdx >= pendingIdx || pendingIdx >= completedIdx {
		t.Errorf(
			"priority sort wrong: running=%d pending=%d completed=%d\n%s",
			runningIdx, pendingIdx, completedIdx, got,
		)
	}
}

func TestLayeredRender_HeightPressureCollapse(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root"), nil)
	_ = dt.AddActivity(ActivityID("child1"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("child2"), []ActivityID{"root"})

	dt.SetRenderMode(RenderModeLayered)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root", ActivityStatusRunning, 0)
	snaps.set(ActivityID("child1"), "Child One", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("child2"), "Child Two", ActivityStatusCompleted, 2*time.Second)

	got := dt.RenderWithSnapshots(snaps.snaps, 4, 0)

	if !strings.Contains(got, "Layer 1: 2 done") {
		t.Errorf("expected layer 1 to be collapsed under height pressure:\n%s", got)
	}

	if strings.Contains(got, "Child One") || strings.Contains(got, "Child Two") {
		t.Errorf("collapsed layer should not show individual children:\n%s", got)
	}
}

func TestLayeredRender_EmptyTree(t *testing.T) {
	dt := NewDependencyTree()
	dt.SetRenderMode(RenderModeLayered)

	got := dt.RenderWithSnapshots(nil, 0, 0)

	if got != msgNoActivitiesToDisplay {
		t.Errorf("expected %q, got %q", msgNoActivitiesToDisplay, got)
	}
}

func TestLayeredRender_VisibleEntries(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root"), nil)
	_ = dt.AddActivity(ActivityID("child"), []ActivityID{"root"})

	dt.SetRenderMode(RenderModeLayered)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root", ActivityStatusRunning, 0)
	snaps.set(ActivityID("child"), "Child", ActivityStatusCompleted, 0)

	entries := dt.VisibleEntriesWithSnapshots(snaps.snaps, 0)

	if len(entries) < 4 {
		t.Fatalf("expected at least 4 entries, got %d", len(entries))
	}

	if entries[0].LayerHeader != "Layer 0" {
		t.Errorf("first entry should be layer 0 header, got %q", entries[0].LayerHeader)
	}

	if entries[1].LayerHeader != "" || len(entries[1].LayerNodes) != 1 {
		t.Errorf("second entry should be a node row, got %+v", entries[1])
	}

	if entries[2].LayerHeader == "" || len(entries[2].LayerNodes) != 0 {
		t.Errorf("third entry should be a separator/header, got %+v", entries[2])
	}
}

func TestLayeredRender_WithRenderModeOption(t *testing.T) {
	ns := NewNOMStyleSubscriber(
		WithCachePath(filepath.Join(t.TempDir(), cacheFilename)),
		WithRenderMode(RenderModeLayered),
	)

	dt := ns.DependencyTree()

	if dt.RenderMode() != RenderModeLayered {
		t.Errorf("WithRenderMode did not set layered mode, got %v", dt.RenderMode())
	}
}

func TestLayeredRender_ModeDefault(t *testing.T) {
	ns := NewNOMStyleSubscriber(WithCachePath(filepath.Join(t.TempDir(), cacheFilename)))

	dt := ns.DependencyTree()

	if dt.RenderMode() != RenderModeTree {
		t.Errorf("default render mode should be tree, got %v", dt.RenderMode())
	}
}

func TestLayeredRender_ThemeApplied(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("done"), nil)
	dt.SetRenderMode(RenderModeLayered)
	dt.theme = ThemeHighContrast

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("done"), "Done", ActivityStatusCompleted, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)

	if !strings.Contains(got, "Done") {
		t.Errorf("label not rendered:\n%s", got)
	}

	stripped := stripANSI(got)

	if !strings.Contains(stripped, string(SymbolCompleted)+" Done") {
		t.Errorf("expected symbol+label layout, got:\n%s", stripped)
	}
}
