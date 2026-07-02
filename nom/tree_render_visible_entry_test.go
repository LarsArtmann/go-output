package nom

import (
	"strings"
	"testing"
	"time"
)

func TestRenderVisibleEntry_LayeredDispatch(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root"), nil)
	_ = dt.AddActivity(ActivityID("child"), []ActivityID{"root"})

	dt.SetRenderMode(RenderModeLayered)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root", ActivityStatusRunning, 0)
	snaps.set(ActivityID("child"), "Child", ActivityStatusCompleted, 0)

	entries := dt.VisibleEntriesWithSnapshots(snaps.snaps, 0)

	if len(entries) < 4 {
		t.Fatalf("expected entries, got %d", len(entries))
	}

	layerHeader := dt.RenderVisibleEntry(entries[0], snaps.snaps, 0)
	if !strings.Contains(layerHeader, "Layer 0") {
		t.Errorf("expected layer header, got %q", layerHeader)
	}

	nodeRow := dt.RenderVisibleEntry(entries[1], snaps.snaps, 0)
	if !strings.Contains(nodeRow, "Root") {
		t.Errorf("expected node row to contain Root, got %q", nodeRow)
	}

	separator := dt.RenderVisibleEntry(entries[2], snaps.snaps, 0)
	if !strings.Contains(separator, "┼") {
		t.Errorf("expected separator, got %q", separator)
	}
}

func TestVisibleEntry_ContainsNode(t *testing.T) {
	node := &ActivityNode{ID: ActivityID("a")}

	entry := VisibleEntry{Node: node}
	if !entry.ContainsNode(ActivityID("a")) {
		t.Error("expected ContainsNode true for Node")
	}

	if entry.ContainsNode(ActivityID("b")) {
		t.Error("expected ContainsNode false for unrelated ID")
	}

	layered := VisibleEntry{LayerNodes: []*ActivityNode{node}}
	if !layered.ContainsNode(ActivityID("a")) {
		t.Error("expected ContainsNode true for LayerNodes")
	}
}

func TestRenderNode_Layered(t *testing.T) {
	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("a"), nil)
	dt.SetRenderMode(RenderModeLayered)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("a"), "Alpha", ActivityStatusRunning, time.Second)

	node := dt.GetNode(ActivityID("a"))
	got := dt.RenderNode(node, snaps.snaps)

	if !strings.Contains(got, "Alpha") {
		t.Errorf("expected layered RenderNode to contain Alpha, got %q", got)
	}
}
