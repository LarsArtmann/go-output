package output

import (
	"encoding/json"
	"strings"
	"testing"
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

	if !strings.Contains(got, `"id": "root"`) {
		t.Errorf("should contain node id, got %q", got)
	}

	if !strings.Contains(got, `"label": "Root"`) {
		t.Errorf("should contain node label, got %q", got)
	}
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

	if !strings.Contains(got, `"children"`) {
		t.Errorf("should contain children array, got %q", got)
	}

	if !strings.Contains(got, `"id": "child"`) {
		t.Errorf("should contain child id, got %q", got)
	}
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

	if !strings.Contains(got, `"metadata"`) {
		t.Errorf("should contain metadata, got %q", got)
	}

	if !strings.Contains(got, `"key": "value"`) {
		t.Errorf("should contain metadata key/value, got %q", got)
	}
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

	var result map[string]any
	if unmarshalErr := json.Unmarshal([]byte(got), &result); unmarshalErr != nil {
		t.Errorf("output should be valid JSON: %v, got %q", unmarshalErr, got)
	}
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

	if !strings.Contains(got, `"id": "grandchild"`) {
		t.Errorf("should contain deeply nested node, got %q", got)
	}
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

	if !strings.Contains(got, `"id": "A"`) {
		t.Errorf("should contain node A, got %q", got)
	}

	if !strings.Contains(got, `"from": "A"`) {
		t.Errorf("should contain edge from A, got %q", got)
	}
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

	if !strings.Contains(got, `"label": "connects"`) {
		t.Errorf("should contain edge label, got %q", got)
	}
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

	if !strings.Contains(got, `"shape": "diamond"`) {
		t.Errorf("should contain node shape, got %q", got)
	}
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

	if !strings.Contains(got, `"type": "service"`) {
		t.Errorf("should contain node metadata, got %q", got)
	}
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

	var result map[string]any
	if unmarshalErr := json.Unmarshal([]byte(got), &result); unmarshalErr != nil {
		t.Errorf("output should be valid JSON: %v, got %q", unmarshalErr, got)
	}
}

func TestJSONTreeRendererMustRender(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := NewTreeNode("root", "Root")
	r.SetRoot(root)

	got := MustRender(r)
	if !strings.Contains(got, `"id": "root"`) {
		t.Errorf("MustRender should work, got %q", got)
	}
}

func TestJSONGraphRendererMustRender(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	r.SetNodes(testNodesAB())

	got := MustRender(r)
	if !strings.Contains(got, `"id": "A"`) {
		t.Errorf("MustRender should work, got %q", got)
	}
}
