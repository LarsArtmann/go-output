package output

import (
	"strings"
	"testing"

	"github.com/go-faster/yaml"
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

	if !strings.Contains(got, "id: root") {
		t.Errorf("should contain node id, got %q", got)
	}

	if !strings.Contains(got, "label: Root") {
		t.Errorf("should contain node label, got %q", got)
	}
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

	if !strings.Contains(got, "children") {
		t.Errorf("should contain children, got %q", got)
	}

	if !strings.Contains(got, "id: child") {
		t.Errorf("should contain child id, got %q", got)
	}
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

	if !strings.Contains(got, "metadata") {
		t.Errorf("should contain metadata, got %q", got)
	}

	if !strings.Contains(got, "key: value") {
		t.Errorf("should contain metadata key/value, got %q", got)
	}
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

	var result map[string]any

	if unmarshalErr := yaml.Unmarshal([]byte(got), &result); unmarshalErr != nil {
		t.Errorf("output should be valid YAML: %v, got %q", unmarshalErr, got)
	}
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

	if !strings.Contains(got, "id: grandchild") {
		t.Errorf("should contain deeply nested node, got %q", got)
	}
}

func TestYAMLGraphRenderer_Empty(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(got, "nodes:") {
		t.Errorf("should contain nodes key, got %q", got)
	}

	if !strings.Contains(got, "edges:") {
		t.Errorf("should contain edges key, got %q", got)
	}
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

	if !strings.Contains(got, "id: A") {
		t.Errorf("should contain node A, got %q", got)
	}

	if !strings.Contains(got, "from: A") {
		t.Errorf("should contain edge from A, got %q", got)
	}
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

	if !strings.Contains(got, "label: connects") {
		t.Errorf("should contain edge label, got %q", got)
	}
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

	if !strings.Contains(got, "shape: diamond") {
		t.Errorf("should contain node shape, got %q", got)
	}
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

	if !strings.Contains(got, "type: service") {
		t.Errorf("should contain node metadata, got %q", got)
	}
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

	var result map[string]any

	if unmarshalErr := yaml.Unmarshal([]byte(got), &result); unmarshalErr != nil {
		t.Errorf("output should be valid YAML: %v, got %q", unmarshalErr, got)
	}
}

func TestYAMLTreeRendererMustRender(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := NewTreeNode("root", "Root")
	r.SetRoot(root)

	got := MustRender(r)
	if !strings.Contains(got, "id: root") {
		t.Errorf("MustRender should work, got %q", got)
	}
}

func TestYAMLGraphRendererMustRender(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	r.SetNodes(testNodesAB())

	got := MustRender(r)
	if !strings.Contains(got, "id: A") {
		t.Errorf("MustRender should work, got %q", got)
	}
}
