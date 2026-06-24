package serialization

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func TestYAMLTreeRenderer_Empty(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != "null\n" {
		t.Errorf("empty tree should render as null, got %q", got)
	}
}

func TestYAMLTreeRenderer_Single(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, "id: root")
	assertOutputContains(t, got, "label: Root")
}

func TestYAMLTreeRenderer_WithChildren(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	child := output.NewTreeNode("child", "Child")
	root.AddChild(child)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, "children")
	assertOutputContains(t, got, "id: child")
}

func TestYAMLTreeRenderer_WithMetadata(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	root.Metadata["key"] = "value"
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, "metadata")
	assertOutputContains(t, got, "key: value")
}

func TestYAMLTreeRenderer_ValidYAML(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	child1 := output.NewTreeNode("c1", "Child 1")
	child2 := output.NewTreeNode("c2", "Child 2")

	root.AddChild(child1)
	root.AddChild(child2)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertValidYAML(t, got)
}

func TestYAMLTreeRenderer_DeepNesting(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	child := output.NewTreeNode("child", "Child")
	grandchild := output.NewTreeNode("grandchild", "Grandchild")
	child.AddChild(grandchild)
	root.AddChild(child)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, "id: grandchild")
}

func TestYAMLGraphRenderer_Empty(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, "nodes:")
	assertOutputContains(t, got, "edges:")
}

func TestYAMLGraphRenderer_WithNodesAndEdges(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()

	got := renderWithABNodes(t, r)

	assertOutputContains(t, got, "id: A")
	assertOutputContains(t, got, "from: A")
}

func TestYAMLGraphRenderer_EdgeWithLabel(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	r.SetNodes(testNodesAB())
	r.SetEdges([]output.GraphEdge{testEdgeAB("connects")})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, "label: connects")
}

func TestYAMLGraphRenderer_NodeWithShape(t *testing.T) {
	t.Parallel()

	testGraphRendererNodeWithShape(t, NewYAMLGraphRenderer(), "shape: diamond")
}

func TestYAMLGraphRenderer_NodeWithMetadata(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	node := output.NewGraphNode("A", "Node A")
	node.Metadata["type"] = "service"
	r.SetNodes([]output.GraphNode{*node})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, "type: service")
}

func TestYAMLGraphRenderer_ValidYAML(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	r.SetNodes(testNodesABC())
	r.SetEdges(testEdgesABC())

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertValidYAML(t, got)
}
