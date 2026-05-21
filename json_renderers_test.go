package output

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
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
	root := NewTreeNode("root", "Root")
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, `"id": "root"`)

	gentest.AssertOutputContains(t, got, `"label": "Root"`)
}

func TestJSONTreeRenderer_WithChildren(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := NewTreeNode("root", "Root")
	child := NewTreeNode("child", "Child")
	root.AddChild(child)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, `"children"`)

	gentest.AssertOutputContains(t, got, `"id": "child"`)
}

func TestJSONTreeRenderer_WithMetadata(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := NewTreeNode("root", "Root")
	root.Metadata["key"] = "value" //nolint:goconst // test value
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, `"metadata"`)

	gentest.AssertOutputContains(t, got, `"key": "value"`)
}

func TestJSONTreeRenderer_ValidJSON(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
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

	gentest.AssertValidJSON(t, got)
}

func TestJSONTreeRenderer_DeepNesting(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
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

	gentest.AssertOutputContains(t, got, `"id": "grandchild"`)
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

	gentest.AssertOutputContains(t, got, `"id": "A"`)

	gentest.AssertOutputContains(t, got, `"from": "A"`)
}

func TestJSONGraphRenderer_EdgeWithLabel(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	r.SetNodes(testNodesAB())
	r.SetEdges([]GraphEdge{testEdgeAB("connects")})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, `"label": "connects"`)
}

func TestJSONGraphRenderer_NodeWithShape(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	r.SetNodes([]GraphNode{newTestNodeWithShape("A", "Node A", ShapeDiamond)})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, `"shape": "diamond"`)
}

func TestJSONGraphRenderer_NodeWithMetadata(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	node := NewGraphNode("A", "Node A")
	node.Metadata["type"] = "service"
	r.SetNodes([]GraphNode{*node})
	r.SetEdges(nil)

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	gentest.AssertOutputContains(t, got, `"type": "service"`)
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

	gentest.AssertValidJSON(t, got)
}

func TestJSONTreeRendererMustRender(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := NewTreeNode("root", "Root")
	r.SetRoot(root)

	got := MustRender(r)
	gentest.AssertOutputContains(t, got, `"id": "root"`)
}

func TestJSONGraphRendererMustRender(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	r.SetNodes(testNodesAB())

	got := MustRender(r)
	gentest.AssertOutputContains(t, got, `"id": "A"`)
}
