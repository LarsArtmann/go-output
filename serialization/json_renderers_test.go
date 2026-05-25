package serialization

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestJSONTreeRenderer_Empty(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != "null" {
		t.Errorf("empty tree should render as null, got %q", got)
	}
}

func TestJSONTreeRenderer_Single(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, `"id": "root"`)
	assertOutputContains(t, got, `"label": "Root"`)
}

func TestJSONTreeRenderer_WithChildren(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	child := output.NewTreeNode("child", "Child")
	root.AddChild(child)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, `"children"`)
	assertOutputContains(t, got, `"id": "child"`)
}

func TestJSONTreeRenderer_WithMetadata(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	root.Metadata["key"] = "value"
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, `"metadata"`)
	assertOutputContains(t, got, `"key": "value"`)
}

func TestJSONTreeRenderer_ValidJSON(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
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

	assertValidJSON(t, got)
}

func TestJSONTreeRenderer_DeepNesting(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
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

	assertOutputContains(t, got, `"id": "grandchild"`)
}

func TestJSONGraphRenderer_Empty(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(got, `"nodes"`) || !strings.Contains(got, `"edges"`) {
		t.Errorf("should contain nodes and edges, got %q", got)
	}
}

func TestJSONGraphRenderer_WithNodesAndEdges(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	r.SetNodes(testNodesAB())
	r.SetEdges(testEdgesAB())

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, `"id": "A"`)
	assertOutputContains(t, got, `"from": "A"`)
}

func TestJSONGraphRenderer_EdgeWithLabel(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	r.SetNodes(testNodesAB())
	r.SetEdges([]output.GraphEdge{testEdgeAB("connects")})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, `"label": "connects"`)
}

func TestJSONGraphRenderer_NodeWithShape(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	r.SetNodes([]output.GraphNode{newTestNodeWithShape("A", "Node A", output.ShapeDiamond)})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, `"shape": "diamond"`)
}

func TestJSONGraphRenderer_NodeWithMetadata(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	node := output.NewGraphNode("A", "Node A")
	node.Metadata["type"] = "service"
	r.SetNodes([]output.GraphNode{*node})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertOutputContains(t, got, `"type": "service"`)
}

func TestJSONGraphRenderer_ValidJSON(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	r.SetNodes(testNodesABC())
	r.SetEdges(testEdgesABC())

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertValidJSON(t, got)
}
