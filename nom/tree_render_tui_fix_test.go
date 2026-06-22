package nom

import (
	"fmt"
	"strings"
	"testing"
)

// buildElisionTree creates a root with 6 children: the first 4 completed, the
// last 2 running. Under a small maxHeight this triggers elideCompletedUnderPressure,
// which appends a synthetic collapse-marker entry (node == nil) for the 4
// completed children. This is the exact shape that panicked the TUI render path
// (VisibleNodesWithSnapshots + RenderNode) before the fix.
func buildElisionTree(t *testing.T) (*DependencyTree, map[ActivityID]ActivitySnapshot) {
	t.Helper()

	tree := NewDependencyTree()
	tree.AddActivity(ActivityID("root"), nil)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root", ActivityStatusRunning, 0)

	for i := range 6 {
		id := ActivityID(fmt.Sprintf("c%d", i))
		tree.AddActivity(id, []ActivityID{"root"})

		status := ActivityStatusCompleted
		if i >= 4 { // last two stay running
			status = ActivityStatusRunning
		}

		snaps.set(id, fmt.Sprintf("Child %d", i), status, 0)
	}

	if err := tree.Build(); err != nil {
		t.Fatalf("build: %v", err)
	}

	return tree, snaps.snaps
}

// TestVisibleNodesWithSnapshots_NeverReturnsNilUnderHeightPressure reproduces the
// TUI panic: under height pressure a collapse-marker entry carries a nil node,
// and the old VisibleNodesWithSnapshots copied it verbatim into the result slice,
// so RenderNode dereferenced node.ID and crashed.
func TestVisibleNodesWithSnapshots_NeverReturnsNilUnderHeightPressure(t *testing.T) {
	t.Parallel()

	tree, snaps := buildElisionTree(t)

	// maxHeight=4: root + 2 running + 1 collapse-marker line for the 4 completed.
	nodes := tree.VisibleNodesWithSnapshots(snaps, 4)

	for i, node := range nodes {
		if node == nil {
			t.Fatalf("VisibleNodesWithSnapshots returned a nil node at index %d "+
				"(collapse marker leaked into the node slice)", i)
		}

		// RenderNode must not panic on any returned node.
		_ = tree.RenderNode(node, snaps)
	}
}

// TestVisibleEntriesWithSnapshots_IncludesCollapseMarker verifies the marker-aware
// API exposes the synthetic collapse line explicitly (Node == nil,
// CollapsedCompleted > 0) instead of smuggling a nil into a []*ActivityNode.
func TestVisibleEntriesWithSnapshots_IncludesCollapseMarker(t *testing.T) {
	t.Parallel()

	tree, snaps := buildElisionTree(t)

	entries := tree.VisibleEntriesWithSnapshots(snaps, 4)

	var markers, realNodes int

	for _, entry := range entries {
		rendered := tree.RenderVisibleEntry(entry, snaps, 0)

		if entry.CollapsedCompleted > 0 {
			if entry.Node != nil {
				t.Errorf("marker entry should have a nil Node, got %v", entry.Node.ID)
			}

			markers++

			if !strings.Contains(rendered, "⋯") || !strings.Contains(rendered, "4 completed") {
				t.Errorf("marker line should render '⋯ 4 completed', got %q", rendered)
			}

			continue
		}

		if entry.Node == nil {
			t.Fatal("non-marker entry should have a non-nil Node")
		}

		realNodes++

		if rendered == "" {
			t.Errorf("real node %v rendered an empty line", entry.Node.ID)
		}
	}

	if markers != 1 {
		t.Errorf("expected exactly 1 collapse marker, got %d (entries=%v)", markers, entries)
	}

	// root + 2 running children = 3 real lines.
	if realNodes != 3 {
		t.Errorf("expected 3 real node lines, got %d (entries=%v)", realNodes, entries)
	}
}
