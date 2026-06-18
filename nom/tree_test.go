package nom

import (
	"errors"
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

		err := dt.AddActivity(ActivityID("a"), "Activity A", nil)
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

		err := dt.AddActivity(ActivityID("parent"), "Parent", nil)
		if err != nil {
			t.Fatalf("AddActivity() error: %v", err)
		}

		err = dt.AddActivity(ActivityID("child"), "Child", []ActivityID{"parent"})
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
	dt.AddActivity(ActivityID("a"), "A", nil)
	dt.AddActivity(ActivityID("b"), "B", nil)

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
	dt.AddActivity(ActivityID("a"), "A", nil)
	dt.AddActivity(ActivityID("b"), "B", []ActivityID{"a"})
	dt.AddActivity(ActivityID("c"), "C", []ActivityID{"b"})

	dt.Build()

	c := dt.GetNode(ActivityID("c"))
	if c.Depth != 2 {
		t.Errorf("c depth = %d, want 2", c.Depth)
	}
}

func TestDependencyTree_GetDisplayActivities(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), "A", nil)

	activities := dt.getDisplayActivities()
	if activities == nil {
		t.Error("GetDisplayActivities() should not return nil")
	}
}

func TestDependencyTree_FindNodesByStatus(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), "A", nil)
	dt.AddActivity(ActivityID("b"), "B", nil)
	dt.UpdateActivityStatus(ActivityID("a"), ActivityStatusRunning, SymbolRunning, ColorRunning, time.Now(), 0)
	dt.UpdateActivityStatus(ActivityID("b"), ActivityStatusCompleted, SymbolCompleted, ColorCompleted, time.Now(), 0)

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
	dt.AddActivity(ActivityID("a"), "A", nil)
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
	dt.AddActivity(ActivityID("a"), "A", nil)
	dt.AddActivity(ActivityID("b"), "B", nil)

	snapshot := dt.snapshotRoots()
	if len(snapshot) != 2 {
		t.Errorf("expected 2 roots, got %d", len(snapshot))
	}
}

func TestDependencyTree_UpdateActivityStatus(t *testing.T) {
	t.Parallel()

	t.Run("existing node", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), "A", nil)

		err := dt.UpdateActivityStatus(
			ActivityID("a"),
			ActivityStatusRunning,
			SymbolRunning,
			ColorRunning,
			time.Now(),
			5*time.Second,
		)
		if err != nil {
			t.Fatalf("UpdateActivityStatus() error: %v", err)
		}

		node := dt.GetNode(ActivityID("a"))
		if node.Status != ActivityStatusRunning {
			t.Errorf("Status = %v, want Running", node.Status)
		}

		if node.Symbol != SymbolRunning {
			t.Errorf("Symbol = %q, want %q", node.Symbol, SymbolRunning)
		}
	})

	t.Run("nonexistent node returns error", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		err := dt.UpdateActivityStatus(
			ActivityID("nonexistent"),
			ActivityStatusRunning,
			SymbolRunning,
			ColorRunning,
			time.Now(),
			0,
		)
		if err == nil {
			t.Error("expected error for nonexistent activity")
		}

		if !errors.Is(err, ErrActivityNotFound) {
			t.Errorf("error should wrap ErrActivityNotFound, got: %v", err)
		}
	})
}

func TestDependencyTree_Render(t *testing.T) {
	t.Parallel()

	t.Run("empty tree shows message", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		got := dt.Render(10)
		if got != msgNoActivitiesToDisplay {
			t.Errorf("Render() on empty tree = %q, want %q", got, msgNoActivitiesToDisplay)
		}
	})

	t.Run("tree with activities renders content", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), "Activity A", nil)
		dt.UpdateActivityStatus(ActivityID("a"), ActivityStatusRunning, SymbolRunning, ColorRunning, time.Now(), 0)

		got := dt.Render(10)
		if got == "" {
			t.Error("Render() should not return empty string for non-empty tree")
		}
	})

	t.Run("max height limits display", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()

		for i := range 10 {
			id := ActivityID(string(rune('a' + i)))
			dt.AddActivity(id, string(rune('a'+i)), nil)
		}

		got := dt.Render(3)
		if got == "" {
			t.Error("Render() should not return empty")
		}
	})
}

func TestDependencyTree_AddActivity_DedupSecondaryParents(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), "Phase", nil)
	dt.AddActivity(ActivityID("step1"), "Step1", []ActivityID{"phase"})

	// Add step2 with same secondary dep twice (simulating re-registration)
	dt.AddActivity(ActivityID("step2"), "Step2", []ActivityID{"phase", "step1"})
	dt.AddActivity(ActivityID("step2"), "Step2", []ActivityID{"phase", "step1"})

	node := dt.GetNode(ActivityID("step2"))
	if len(node.SecondaryParents) != 1 {
		t.Errorf("SecondaryParents = %v, want 1 entry", node.SecondaryParents)
	}
}

func TestDependencyTree_RenderNode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), "Root Task", nil)
	dt.AddActivity(ActivityID("child"), "Child Task", []ActivityID{"root"})

	nodes := dt.VisibleNodes(10)
	if len(nodes) == 0 {
		t.Fatal("expected at least one display node")
	}

	out := dt.RenderNode(nodes[0], nil)
	if out == "" {
		t.Error("RenderNode should produce non-empty output")
	}

	if !contains(out, "Root") {
		t.Errorf("RenderNode output should contain node name, got: %q", out)
	}
}

func TestDependencyTree_VisibleNodes(t *testing.T) {
	t.Parallel()

	t.Run("returns nodes up to maxHeight", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), "A", nil)
		dt.AddActivity(ActivityID("b"), "B", nil)
		dt.AddActivity(ActivityID("c"), "C", nil)

		visible := dt.VisibleNodes(2)
		if len(visible) != 2 {
			t.Errorf("VisibleNodes(2) = %d nodes, want 2", len(visible))
		}
	})

	t.Run("zero or negative maxHeight returns all", func(t *testing.T) {
		t.Parallel()

		dt := NewDependencyTree()
		dt.AddActivity(ActivityID("a"), "A", nil)
		dt.AddActivity(ActivityID("b"), "B", nil)

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

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
