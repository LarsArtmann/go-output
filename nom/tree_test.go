package nom

import (
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/go-output/testhelpers"
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

		err := dt.AddActivity(ActivityID("a"), NewActivity("a", "Activity A"), nil)
		if err != nil {
			t.Fatalf("AddActivity() error: %v", err)
		}

		node := dt.GetNode(ActivityID("a"))
		if node == nil {
			t.Fatal("node should exist after AddActivity")
		}

		testhelpers.AssertEqual(t, "ActivityName", "", node.Label.Get(), "Activity A")
	})

	t.Run("activity with dependency creates parent-child", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		err := dt.AddActivity(ActivityID("parent"), NewActivity("parent", "Parent"), nil)
		if err != nil {
			t.Fatalf("AddActivity() error: %v", err)
		}

		err = dt.AddActivity(ActivityID("child"), NewActivity("child", "Child"), []ActivityID{"parent"})
		if err != nil {
			t.Fatalf("AddActivity() error: %v", err)
		}

		child := dt.GetNode(ActivityID("child"))
		if child == nil {
			t.Fatal("child node should exist")
		}

		assertChildParentID(t, child, "parent")

		parent := dt.GetNode(ActivityID("parent"))
		if len(parent.Children) != 1 || parent.Children[0].ID.Get() != "child" {
			t.Error("parent should have child")
		}
	})
}

func TestDependencyTree_Build(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
	dt.AddActivity(ActivityID("b"), NewActivity("b", "B"), nil)

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
			t.Errorf("node %q should be a root", root.ID.Get())
		}

		if root.Depth != 0 {
			t.Errorf("root %q depth = %d, want 0", root.ID.Get(), root.Depth)
		}
	}
}

func TestDependencyTree_Build_DepthCalculation(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
	dt.AddActivity(ActivityID("b"), NewActivity("b", "B"), []ActivityID{"a"})
	dt.AddActivity(ActivityID("c"), NewActivity("c", "C"), []ActivityID{"b"})

	dt.Build()

	c := dt.GetNode(ActivityID("c"))
	if c.Depth != 2 {
		t.Errorf("c depth = %d, want 2", c.Depth)
	}
}

func TestDependencyTree_FindNodesByStatus(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
	dt.AddActivity(ActivityID("b"), NewActivity("b", "B"), nil)
	testSetStatus(dt, ActivityID("a"), ActivityStatusRunning, time.Now())
	testSetStatus(dt, ActivityID("b"), ActivityStatusCompleted, time.Now())

	running := dt.findNodesByStatus(ActivityStatusRunning)
	if len(running) != 1 {
		t.Errorf("expected 1 running node, got %d", len(running))
	}

	if running[0].ID.Get() != "a" {
		t.Errorf("running node = %q, want %q", running[0].ID.Get(), "a")
	}
}

func TestDependencyTree_Clear(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
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
	dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
	dt.AddActivity(ActivityID("b"), NewActivity("b", "B"), nil)

	snapshot := dt.snapshotRoots()
	if len(snapshot) != 2 {
		t.Errorf("expected 2 roots, got %d", len(snapshot))
	}
}

func TestDependencyTree_DirectStatusMutation(t *testing.T) {
	t.Parallel()

	t.Run("existing node", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)

		testSetStatus(dt, ActivityID("a"), ActivityStatusRunning, time.Now())

		node := dt.GetNode(ActivityID("a"))
		if node.Status != ActivityStatusRunning {
			t.Errorf("Status = %v, want Running", node.Status)
		}

		if node.Symbol != SymbolRunning {
			t.Errorf("Symbol = %q, want %q", node.Symbol, SymbolRunning)
		}
	})

	t.Run("nonexistent node is no-op", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		testSetStatus(dt, ActivityID("nonexistent"), ActivityStatusRunning, time.Now())
		// No panic, no error — the shared-pointer model makes this a no-op.
	})
}

func TestDependencyTree_Render(t *testing.T) {
	t.Parallel()

	t.Run("empty tree shows message", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		got := dt.RenderString(10)
		if got != msgNoActivitiesToDisplay {
			t.Errorf("Render() on empty tree = %q, want %q", got, msgNoActivitiesToDisplay)
		}
	})

	t.Run("tree with activities renders content", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), NewActivity("a", "Activity A"), nil)
		testSetStatus(dt, ActivityID("a"), ActivityStatusRunning, time.Now())

		got := dt.RenderString(10)
		if got == "" {
			t.Error("Render() should not return empty string for non-empty tree")
		}
	})

	t.Run("max height limits display", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		for i := range 10 {
			id := ActivityID(string(rune('a' + i)))
			dt.AddActivity(id, NewActivity(string(id), string(rune('a'+i))), nil)
		}

		got := dt.RenderString(3)
		if got == "" {
			t.Error("Render() should not return empty")
		}
	})
}

func TestDependencyTree_AddActivity_DedupSecondaryParents(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), NewActivity("phase", "Phase"), nil)
	dt.AddActivity(ActivityID("step1"), NewActivity("step1", "Step1"), []ActivityID{"phase"})

	// Add step2 with same secondary dep twice (simulating re-registration)
	dt.AddActivity(ActivityID("step2"), NewActivity("step2", "Step2"), []ActivityID{"phase", "step1"})
	dt.AddActivity(ActivityID("step2"), NewActivity("step2", "Step2"), []ActivityID{"phase", "step1"})

	node := dt.GetNode(ActivityID("step2"))
	if len(node.SecondaryParents) != 1 {
		t.Errorf("SecondaryParents = %v, want 1 entry", node.SecondaryParents)
	}
}

func TestDependencyTree_RenderNode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), NewActivity("root", "Root Task"), nil)
	dt.AddActivity(ActivityID("child"), NewActivity("child", "Child Task"), []ActivityID{"root"})

	nodes := dt.VisibleNodes(10)
	if len(nodes) == 0 {
		t.Fatal("expected at least one display node")
	}

	out := dt.RenderNode(nodes[0], nil)
	if out == "" {
		t.Error("RenderNode should produce non-empty output")
	}

	if !strings.Contains(out, "Root") {
		t.Errorf("RenderNode output should contain node name, got: %q", out)
	}
}

func TestDependencyTree_VisibleNodes(t *testing.T) {
	t.Parallel()

	t.Run("returns nodes up to maxHeight", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
		dt.AddActivity(ActivityID("b"), NewActivity("b", "B"), nil)
		dt.AddActivity(ActivityID("c"), NewActivity("c", "C"), nil)

		visible := dt.VisibleNodes(2)
		if len(visible) != 2 {
			t.Errorf("VisibleNodes(2) = %d nodes, want 2", len(visible))
		}
	})

	t.Run("zero or negative maxHeight returns all", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
		dt.AddActivity(ActivityID("b"), NewActivity("b", "B"), nil)

		visible := dt.VisibleNodes(0)
		if len(visible) != 2 {
			t.Errorf("VisibleNodes(0) = %d nodes, want 2", len(visible))
		}
	})

	t.Run("empty tree returns empty", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		visible := dt.VisibleNodes(10)
		if len(visible) != 0 {
			t.Errorf("VisibleNodes on empty tree = %d, want 0", len(visible))
		}
	})
}
