package nom

import (
	"path/filepath"
	"strconv"
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

	if got != MsgNoActivities {
		t.Errorf("expected %q, got %q", MsgNoActivities, got)
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
	ns := NewNOMSubscriber(
		WithCachePath(filepath.Join(t.TempDir(), cacheFilename)),
		WithRenderMode(RenderModeLayered),
	)

	dt := ns.DependencyTree()

	if dt.RenderMode() != RenderModeLayered {
		t.Errorf("WithRenderMode did not set layered mode, got %v", dt.RenderMode())
	}
}

func TestLayeredRender_ModeDefault(t *testing.T) {
	ns := NewNOMSubscriber(WithCachePath(filepath.Join(t.TempDir(), cacheFilename)))

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

func TestLayeredRender_HideFutureLayers(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("setup"), nil)
	_ = dt.AddActivity(ActivityID("compile"), []ActivityID{"setup"})
	_ = dt.AddActivity(ActivityID("test"), []ActivityID{"compile"})
	_ = dt.AddActivity(ActivityID("deploy"), []ActivityID{"test"})

	dt.SetRenderMode(RenderModeLayered)
	dt.hideFutureLayers = true

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("setup"), "Setup", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("compile"), "Compile", ActivityStatusRunning, 0)
	snaps.set(ActivityID("test"), "Test", ActivityStatusPending, 0)
	snaps.set(ActivityID("deploy"), "Deploy", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)
	stripped := stripANSI(got)

	// Active layers (0 and 1) should show full labels.
	if !strings.Contains(stripped, "Setup") {
		t.Errorf("expected Setup in active layer:\n%s", stripped)
	}

	if !strings.Contains(stripped, "Compile") {
		t.Errorf("expected Compile in active layer:\n%s", stripped)
	}

	// Future layers (all pending) should be collapsed.
	if strings.Contains(stripped, "Test") {
		t.Errorf("Test should be hidden in future layer:\n%s", stripped)
	}

	if strings.Contains(stripped, "Deploy") {
		t.Errorf("Deploy should be hidden in future layer:\n%s", stripped)
	}

	// Should show pending summary lines.
	if !strings.Contains(stripped, "pending") {
		t.Errorf("expected pending summary in collapsed layers:\n%s", stripped)
	}
}

func TestLayeredRender_FutureLayersDisabled(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("setup"), nil)
	_ = dt.AddActivity(ActivityID("test"), []ActivityID{"setup"})

	dt.SetRenderMode(RenderModeLayered)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("setup"), "Setup", ActivityStatusRunning, 0)
	snaps.set(ActivityID("test"), "Test", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)
	stripped := stripANSI(got)

	// Without hideFutureLayers, pending nodes should be visible.
	if !strings.Contains(stripped, "Test") {
		t.Errorf("Test should be visible when hideFutureLayers is off:\n%s", stripped)
	}
}

func TestLayeredRender_CategoryCollapse(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root"), nil)
	_ = dt.AddActivity(ActivityID("compile-a"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("compile-b"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("compile-c"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("deploy"), []ActivityID{"root"})

	dt.SetRenderMode(RenderModeLayered)
	dt.showCategory = true
	dt.collapseCategories = true

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root", ActivityStatusCompleted, 1*time.Second)
	snaps.setCategory(ActivityID("compile-a"), "Compile A", ActivityStatusRunning, ActivityCategory("build"))
	snaps.setCategory(ActivityID("compile-b"), "Compile B", ActivityStatusRunning, ActivityCategory("build"))
	snaps.setCategory(ActivityID("compile-c"), "Compile C", ActivityStatusPending, ActivityCategory("build"))
	snaps.setCategory(ActivityID("deploy"), "Deploy", ActivityStatusPending, ActivityCategory("deploy"))

	got := dt.RenderWithSnapshots(snaps.snaps, 0, 0)
	stripped := stripANSI(got)

	// The 3 build tasks should be collapsed into a summary line.
	if !strings.Contains(stripped, "3 build tasks") {
		t.Errorf("expected collapsed category summary '3 build tasks':\n%s", stripped)
	}

	// Deploy (single node in its category) should still be shown individually.
	if !strings.Contains(stripped, "Deploy") {
		t.Errorf("expected Deploy as individual node:\n%s", stripped)
	}

	// Individual compile labels should NOT appear (collapsed).
	if strings.Contains(stripped, "Compile A") {
		t.Errorf("Compile A should be collapsed:\n%s", stripped)
	}
}

func TestVisibleEntryKind_Classification(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root1"), nil)
	_ = dt.AddActivity(ActivityID("child1"), []ActivityID{"root1"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root1"), "Root One", ActivityStatusRunning, 1*time.Second)
	snaps.set(ActivityID("child1"), "Child One", ActivityStatusPending, 0)

	if err := dt.Build(); err != nil {
		t.Fatalf("Build: %v", err)
	}

	entries := dt.collectLayeredEntries(snaps.snaps, 50)

	var sawHeader, sawSeparator, sawRow bool
	for _, entry := range entries {
		switch entry.Kind() {
		case KindLayerHeader:
			sawHeader = true
			if entry.LayerHeader == "" {
				t.Error("KindLayerHeader entry must carry header text")
			}
		case KindSeparator:
			sawSeparator = true
			if entry.LayerHeader == "" {
				t.Error("KindSeparator entry must carry the separator line")
			}
		case KindLayerRow:
			sawRow = true
			if len(entry.LayerNodes) == 0 {
				t.Error("KindLayerRow entry must carry nodes")
			}
		case KindEmpty, KindNode, KindCollapsed, KindPhase:
			t.Errorf("unexpected kind %d in layered collection", entry.Kind())
		}
	}

	if !sawHeader || !sawSeparator || !sawRow {
		t.Errorf("missing kinds: header=%v separator=%v row=%v", sawHeader, sawSeparator, sawRow)
	}
}

func TestVisibleEntryKind_ZeroValueAndPayloadKinds(t *testing.T) {
	t.Parallel()

	if got := (VisibleEntry{}).Kind(); got != KindEmpty {
		t.Errorf("zero entry: want KindEmpty, got %d", got)
	}

	node := &ActivityNode{ID: ActivityID("a")}
	if got := (VisibleEntry{Node: node}).Kind(); got != KindNode {
		t.Errorf("node entry: want KindNode, got %d", got)
	}

	if got := (VisibleEntry{CollapsedCompleted: 3}).Kind(); got != KindCollapsed {
		t.Errorf("collapsed entry: want KindCollapsed, got %d", got)
	}

	if got := (VisibleEntry{PhaseCounts: &PhaseCounts{}}).Kind(); got != KindPhase {
		t.Errorf("phase entry: want KindPhase, got %d", got)
	}
}

func TestLayeredSeparator_AlignsWithHeaderPipe(t *testing.T) {
	t.Parallel()

	for _, maxDepth := range []int{0, 3, 9, 10, 12, 99, 100, 1234} {
		header := (&DependencyTree{}).renderLayeredHeader("Layer "+strconv.Itoa(maxDepth), 0)
		separator := layeredSeparator(maxDepth)

		pipeCol := strings.Index(header, "│")
		crossCol := strings.Index(separator, "┼")

		if pipeCol < 0 || crossCol < 0 {
			t.Fatalf("maxDepth=%d: missing markers header=%q separator=%q", maxDepth, header, separator)
		}

		if pipeCol != crossCol {
			t.Errorf("maxDepth=%d: ┼ at column %d but header │ at column %d\nheader=%q\nseparator=%q",
				maxDepth, crossCol, pipeCol, header, separator)
		}
	}
}
