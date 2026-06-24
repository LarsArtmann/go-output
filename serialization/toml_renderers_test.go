package serialization

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestTOMLGraphRenderer_Empty(t *testing.T) {
	t.Parallel()

	r := NewTOMLGraphRenderer()

	out, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertAllContained(t, out, "nodes")
}

func TestTOMLGraphRenderer_WithNodesAndEdges(t *testing.T) {
	t.Parallel()

	r := NewTOMLGraphRenderer()

	out := renderWithABNodes(t, r)

	testhelpers.AssertAllContained(t, out, "Node A", "Node B")
}

func TestTOMLGraphRenderer_EdgeWithLabel(t *testing.T) {
	t.Parallel()

	r := NewTOMLGraphRenderer()
	r.SetNodes([]output.GraphNode{newTestNode("a", "A")})
	r.SetEdges([]output.GraphEdge{testEdgeAB("connects")})

	out, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertAllContained(t, out, "connects")
}

func TestTOMLTreeRenderer(t *testing.T) {
	t.Parallel()

	t.Run("renders tree", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTreeRenderer()
		root := output.NewTreeNode("root", "Root")
		child := output.NewTreeNode("child1", "Child")
		child.Metadata = map[string]string{"key": "value"}
		root.AddChild(child)
		r.SetRoot(root)

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertAllContained(t, out, "Root", "Child")
	})

	t.Run("nil root returns empty", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTreeRenderer()

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if out != "" {
			t.Errorf("expected empty, got %q", out)
		}
	})

	t.Run("deep tree", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTreeRenderer()
		root := output.NewTreeNode("root", "Root")
		child := output.NewTreeNode("child", "Child")
		grandchild := output.NewTreeNode("gc", "Grandchild")
		child.AddChild(grandchild)
		root.AddChild(child)
		r.SetRoot(root)

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertAllContained(t, out, "Grandchild")
	})
}
