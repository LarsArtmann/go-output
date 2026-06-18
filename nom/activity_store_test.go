package nom

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func TestNewActivityStore(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	if s.Size() != 0 {
		t.Errorf("Size = %d, want 0", s.Size())
	}
}

func TestActivityStore_UpsertGet(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	a := NewActivity("build", "Build")
	s.Upsert(a)

	got, ok := s.Get(a.ID)
	if !ok {
		t.Fatal("activity not found after Upsert")
	}

	if got.Label.Get() != "Build" {
		t.Errorf("Label = %q, want %q", got.Label.Get(), "Build")
	}
}

func TestActivityStore_GetNotFound(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	id := output.NewBrandedID[output.GraphNodeIDBrand]("missing")

	_, ok := s.Get(id)
	if ok {
		t.Error("should not find missing activity")
	}
}

func TestActivityStore_Nodes(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	s.Upsert(NewActivity("a", "Alpha"))
	s.Upsert(NewActivity("b", "Beta"))

	nodes := s.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("len(Nodes) = %d, want 2", len(nodes))
	}

	for _, n := range nodes {
		if n.ID.Get() != "a" && n.ID.Get() != "b" {
			t.Errorf("unexpected ID %q", n.ID.Get())
		}
	}
}

func TestActivityStore_Edges(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	from := output.NewBrandedID[output.GraphNodeIDBrand]("parent")
	to := output.NewBrandedID[output.GraphNodeIDBrand]("child")
	s.AddEdge(from, to)

	edges := s.Edges()
	if len(edges) != 1 {
		t.Fatalf("len(Edges) = %d, want 1", len(edges))
	}

	if edges[0].From.Get() != "parent" || edges[0].To.Get() != "child" {
		t.Errorf("Edge = %s → %s, want parent → child", edges[0].From.Get(), edges[0].To.Get())
	}
}

func TestActivityStore_Roots(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	s.Upsert(NewActivity("root", "Root"))
	s.Upsert(NewActivity("child", "Child"))

	from := output.NewBrandedID[output.GraphNodeIDBrand]("root")
	to := output.NewBrandedID[output.GraphNodeIDBrand]("child")
	s.AddEdge(from, to)

	roots := s.Roots()
	if len(roots) != 1 {
		t.Fatalf("len(Roots) = %d, want 1", len(roots))
	}

	if roots[0].Get() != "root" {
		t.Errorf("Root = %q, want %q", roots[0].Get(), "root")
	}
}

func TestActivityStore_Counts(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()

	r1 := NewActivity("r1", "Running1")
	r1.SetRunning()

	r2 := NewActivity("r2", "Running2")
	r2.SetRunning()
	s.Upsert(r1)
	s.Upsert(r2)

	c := NewActivity("c1", "Completed1")
	c.SetCompleted()
	s.Upsert(c)

	f := NewActivity("f1", "Failed1")
	f.SetFailed(errTestFailure)
	s.Upsert(f)

	p := NewActivity("p1", "Pending1")
	s.Upsert(p)

	running, completed, failed, pending := s.Counts()
	if running != 2 {
		t.Errorf("running = %d, want 2", running)
	}

	if completed != 1 {
		t.Errorf("completed = %d, want 1", completed)
	}

	if failed != 1 {
		t.Errorf("failed = %d, want 1", failed)
	}

	if pending != 1 {
		t.Errorf("pending = %d, want 1", pending)
	}
}

func TestActivityStore_Clear(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	s.Upsert(NewActivity("a", "A"))
	s.AddEdge(
		output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		output.NewBrandedID[output.GraphNodeIDBrand]("b"),
	)

	s.Clear()

	if s.Size() != 0 {
		t.Errorf("Size after Clear = %d, want 0", s.Size())
	}

	if len(s.Edges()) != 0 {
		t.Errorf("Edges after Clear = %d, want 0", len(s.Edges()))
	}
}

func TestActivityStore_All(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	s.Upsert(NewActivity("x", "X"))
	s.Upsert(NewActivity("y", "Y"))

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("len(All) = %d, want 2", len(all))
	}
}

func TestActivityStore_NodesProjectionIsCopy(t *testing.T) {
	t.Parallel()

	s := NewActivityStore()
	s.Upsert(NewActivity("a", "Alpha"))

	nodes := s.Nodes()
	nodes[0].Label = output.NewBrandedID[output.GraphNodeLabelBrand]("MUTATED")

	// Original store should be unaffected
	got, _ := s.Get(nodes[0].ID)
	if got.Label.Get() != "Alpha" {
		t.Errorf("mutation leaked into store: %q", got.Label.Get())
	}
}
