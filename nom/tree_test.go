package nom

import (
	"strings"
	"testing"
	"time"
)

func TestNewDependencyTree(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	if dt == nil {
		t.Fatal("NewDependencyTree() returned nil")
	}
}

func TestDependencyTree_AddActivity(t *testing.T) {
	t.Parallel()

	t.Run("single root activity", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		err := dt.AddActivity(ActivityID("a"), nil)
		if err != nil {
			t.Fatalf("AddActivity() error: %v", err)
		}

		node := dt.GetNode(ActivityID("a"))
		if node == nil {
			t.Fatal("node should exist after AddActivity")
		}

		if node.ID != ActivityID("a") {
			t.Errorf("node.ID = %q, want %q", node.ID, "a")
		}
	})

	t.Run("activity with dependency creates parent-child", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		err := dt.AddActivity(ActivityID("parent"), nil)
		if err != nil {
			t.Fatalf("AddActivity() error: %v", err)
		}

		err = dt.AddActivity(ActivityID("child"), []ActivityID{"parent"})
		if err != nil {
			t.Fatalf("AddActivity() error: %v", err)
		}

		child := dt.GetNode(ActivityID("child"))
		if child == nil {
			t.Fatal("child node should exist")
		}

		assertChildParentID(t, child, "parent")

		parent := dt.GetNode(ActivityID("parent"))
		if len(parent.Children) != 1 || parent.Children[0].ID != "child" {
			t.Error("parent should have child")
		}
	})
}

func TestDependencyTree_Build(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("b"), nil)

	err := dt.Build()
	if err != nil {
		t.Fatalf("Build() error: %v", err)
	}

	roots := dt.GetRootNodes()
	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}

	for _, root := range roots {
		if !root.IsRoot {
			t.Errorf("node %q should be a root", root.ID)
		}

		if root.Depth != 0 {
			t.Errorf("root %q depth = %d, want 0", root.ID, root.Depth)
		}
	}
}

func TestDependencyTree_Build_DepthCalculation(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("b"), []ActivityID{"a"})
	dt.AddActivity(ActivityID("c"), []ActivityID{"b"})

	dt.Build()

	c := dt.GetNode(ActivityID("c"))
	if c.Depth != 2 {
		t.Errorf("c depth = %d, want 2", c.Depth)
	}
}

func TestDependencyTree_FindNodesByStatus(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("b"), nil)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("a"), "A", ActivityStatusRunning, 0)
	snaps.set(ActivityID("b"), "B", ActivityStatusCompleted, 0)

	running := dt.findNodesByStatus(ActivityStatusRunning, snaps.snaps)
	if len(running) != 1 {
		t.Errorf("expected 1 running node, got %d", len(running))
	}

	if running[0].ID != "a" {
		t.Errorf("running node = %q, want %q", running[0].ID, "a")
	}
}

func TestDependencyTree_Clear(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)
	dt.Build()

	dt.Clear()

	node := dt.GetNode(ActivityID("a"))
	if node != nil {
		t.Error("node should be nil after Clear()")
	}
}

func TestDependencyTree_SnapshotRoots(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("b"), nil)

	snapshot := dt.snapshotRoots()
	if len(snapshot) != 2 {
		t.Errorf("expected 2 roots, got %d", len(snapshot))
	}
}

func TestDependencyTree_Render(t *testing.T) {
	t.Parallel()

	t.Run("empty tree shows message", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		got := dt.RenderWithSnapshots(nil, 10, 0)
		if got != msgNoActivitiesToDisplay {
			t.Errorf("Render() on empty tree = %q, want %q", got, msgNoActivitiesToDisplay)
		}
	})

	t.Run("tree with activities renders content", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), nil)

		snaps := newSnapshotBuilder()
		snaps.set(ActivityID("a"), "Activity A", ActivityStatusRunning, 0)

		got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
		if got == "" {
			t.Error("Render should not return empty string for non-empty tree")
		}
	})

	t.Run("max height limits display", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		for i := range 10 {
			id := ActivityID(string(rune('a' + i)))
			dt.AddActivity(id, nil)
		}

		got := dt.RenderWithSnapshots(nil, 3, 0)
		if got == "" {
			t.Error("Render should not return empty")
		}
	})
}

func TestDependencyTree_AddActivity_DedupSecondaryParents(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), nil)
	dt.AddActivity(ActivityID("step1"), []ActivityID{"phase"})

	dt.AddActivity(ActivityID("step2"), []ActivityID{"phase", "step1"})
	dt.AddActivity(ActivityID("step2"), []ActivityID{"phase", "step1"})

	node := dt.GetNode(ActivityID("step2"))
	if len(node.SecondaryParents) != 1 {
		t.Errorf("SecondaryParents = %v, want 1 entry", node.SecondaryParents)
	}
}

func TestDependencyTree_RenderNode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), nil)
	dt.AddActivity(ActivityID("child"), []ActivityID{"root"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root Task", ActivityStatusRunning, 0)
	snaps.set(ActivityID("child"), "Child Task", ActivityStatusPending, 0)

	nodes := dt.VisibleNodesWithSnapshots(snaps.snaps, 10)
	if len(nodes) == 0 {
		t.Fatal("expected at least one display node")
	}

	out := dt.RenderNode(nodes[0], nil, snaps.snaps)
	if out == "" {
		t.Error("RenderNode should produce non-empty output")
	}

	if !strings.Contains(out, "Root") {
		t.Errorf("RenderNode output should contain node name, got: %q", out)
	}
}

func TestDependencyTree_VisibleNodes(t *testing.T) {
	t.Run("returns nodes up to maxHeight", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), nil)
		dt.AddActivity(ActivityID("b"), nil)
		dt.AddActivity(ActivityID("c"), nil)

		visible := dt.VisibleNodesWithSnapshots(nil, 2)
		if len(visible) != 2 {
			t.Errorf("VisibleNodes(2) = %d nodes, want 2", len(visible))
		}
	})

	t.Run("zero or negative maxHeight returns all", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), nil)
		dt.AddActivity(ActivityID("b"), nil)

		visible := dt.VisibleNodesWithSnapshots(nil, 0)
		if len(visible) != 2 {
			t.Errorf("VisibleNodes(0) = %d nodes, want 2", len(visible))
		}
	})

	t.Run("empty tree returns empty", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		visible := dt.VisibleNodesWithSnapshots(nil, 10)
		if len(visible) != 0 {
			t.Errorf("VisibleNodes on empty tree = %d, want 0", len(visible))
		}
	})
}

func TestActivityNode_removeChild(t *testing.T) {
	t.Parallel()

	parent := newActivityNode(ActivityID("parent"), "")
	childA := newActivityNode(ActivityID("a"), "")
	childB := newActivityNode(ActivityID("b"), "")
	parent.Children = []*ActivityNode{childA, childB}

	parent.removeChild(ActivityID("a"))

	if parent.hasChild(ActivityID("a")) {
		t.Error("child 'a' should be removed")
	}

	if !parent.hasChild(ActivityID("b")) {
		t.Error("child 'b' should still be present")
	}

	if len(parent.Children) != 1 {
		t.Errorf("len(Children) = %d, want 1", len(parent.Children))
	}

	parent.removeChild(ActivityID("missing"))

	if len(parent.Children) != 1 {
		t.Errorf("removing missing child changed Children: len = %d, want 1", len(parent.Children))
	}
}

// TestDependencyTree_SnapshotRendering_VerifiesStatusInOutput ensures the
// snapshot-based render path reflects the activity's status correctly. This
// replaces the old TestDependencyTree_DirectStatusMutation test which mutated
// node.Activity fields directly (no longer possible — the node stores only ID).
func TestDependencyTree_SnapshotRendering_VerifiesStatusInOutput(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("a"), "Activity A", ActivityStatusRunning, 5*time.Second)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)

	if !strings.Contains(got, string(SymbolRunning)) {
		t.Errorf("render should contain running symbol, got:\n%s", got)
	}

	if !strings.Contains(got, "Activity A") {
		t.Errorf("render should contain label, got:\n%s", got)
	}
}
