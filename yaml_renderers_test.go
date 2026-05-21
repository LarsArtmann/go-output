package output

import (
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
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
	root := NewTreeNode("root", "Root")
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "id: root")

	gentest.AssertOutputContains(t, got, "label: Root")
}

func TestYAMLTreeRenderer_WithChildren(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := NewTreeNode("root", "Root")
	child := NewTreeNode("child", "Child")
	root.AddChild(child)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "children")

	gentest.AssertOutputContains(t, got, "id: child")
}

func TestYAMLTreeRenderer_WithMetadata(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := NewTreeNode("root", "Root")
	root.Metadata["key"] = "value" //nolint:goconst // test value
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "metadata")

	gentest.AssertOutputContains(t, got, "key: value")
}

func TestYAMLTreeRenderer_ValidYAML(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := NewTreeNode("root", "Root")
	child1 := NewTreeNode("c1", "Child 1")
	child2 := NewTreeNode("c2", "Child 2")

	root.AddChild(child1)
	root.AddChild(child2)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertValidYAML(t, got)
}

func TestYAMLTreeRenderer_DeepNesting(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := NewTreeNode("root", "Root")
	child := NewTreeNode("child", "Child")
	grandchild := NewTreeNode("grandchild", "Grandchild")
	child.AddChild(grandchild)
	root.AddChild(child)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "id: grandchild")
}

func TestYAMLGraphRenderer_Empty(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "nodes:")

	gentest.AssertOutputContains(t, got, "edges:")
}

func TestYAMLGraphRenderer_WithNodesAndEdges(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	r.SetNodes(testNodesAB())
	r.SetEdges(testEdgesAB())

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "id: A")

	gentest.AssertOutputContains(t, got, "from: A")
}

func TestYAMLGraphRenderer_EdgeWithLabel(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	r.SetNodes(testNodesAB())
	r.SetEdges([]GraphEdge{testEdgeAB("connects")})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "label: connects")
}

func TestYAMLGraphRenderer_NodeWithShape(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	r.SetNodes([]GraphNode{newTestNodeWithShape("A", "Node A", ShapeDiamond)})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "shape: diamond")
}

func TestYAMLGraphRenderer_NodeWithMetadata(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	node := NewGraphNode("A", "Node A")
	node.Metadata["type"] = "service"
	r.SetNodes([]GraphNode{*node})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, "type: service")
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

	gentest.AssertValidYAML(t, got)
}

func TestYAMLTreeRendererMustRender(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := NewTreeNode("root", "Root")
	r.SetRoot(root)

	got := MustRender(r)
	gentest.AssertOutputContains(t, got, "id: root")
}

func TestYAMLGraphRendererMustRender(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	r.SetNodes(testNodesAB())

	got := MustRender(r)
	gentest.AssertOutputContains(t, got, "id: A")
}
