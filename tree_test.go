package output

import (
	"strings"
	"testing"
)

func TestASCIITreeRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()

	// Test empty tree
	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != "" {
		t.Errorf("Empty tree should return empty string, got: %s", got)
	}

	// Build a test tree
	root := NewTreeNode("root", "Root")
	child1 := NewTreeNode("child1", "Child 1")
	child2 := NewTreeNode("child2", "Child 2")
	child1.AddChild(NewTreeNode("leaf1", "Leaf 1"))
	child1.AddChild(NewTreeNode("leaf2", "Leaf 2"))
	root.AddChild(child1)
	root.AddChild(child2)

	renderer.SetRoot(root)

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Verify output contains expected elements
	assertContains(t, output, "Root", "Output should contain 'Root'")
	assertContains(t, output, "Child 1", "Output should contain 'Child 1'")
	assertContains(t, output, "Child 2", "Output should contain 'Child 2'")
	assertContains(t, output, "Leaf 1", "Output should contain 'Leaf 1'")
}

func TestASCIITreeRendererWithMetadata(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()

	node := NewTreeNode("node", "Node with Meta")
	node.Metadata["key"] = "value"
	node.Metadata["count"] = "42"

	renderer.SetRoot(node)

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "Node with Meta", "Output should contain node label")
	assertContains(t, output, "key: value", "Output should contain metadata")
}

func TestTreeRendererFromTableData(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"Name", "Age", "City"})
	data.AddRow([]string{"Alice", "30", "NYC"})
	data.AddRow([]string{"Bob", "25", "LA"})

	renderer := TreeRendererFromTableData(data)

	output, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, output, "Data", "Output should contain 'Data' as root")
	assertContains(t, output, "Headers", "Output should contain 'Headers' section")
	assertContains(t, output, "Rows", "Output should contain 'Rows' section")
	assertContains(t, output, "Alice", "Output should contain 'Alice' row")
}

func TestNewTreeNode(t *testing.T) {
	t.Parallel()

	node := NewTreeNode("id1", "Label 1")

	if node.ID.Get() != "id1" {
		t.Errorf("Expected ID 'id1', got '%s'", node.ID.Get())
	}

	if node.Label.Get() != "Label 1" {
		t.Errorf("Expected Label 'Label 1', got '%s'", node.Label.Get())
	}

	if len(node.Children) != 0 {
		t.Error("New node should have no children")
	}

	if node.Metadata == nil {
		t.Error("New node should have non-nil metadata map")
	}
}

func TestTreeNodeAddChild(t *testing.T) {
	t.Parallel()

	parent := NewTreeNode("parent", "Parent")
	child := NewTreeNode("child", "Child")

	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Errorf("Parent should have 1 child, got %d", len(parent.Children))
	}

	if parent.Children[0] != child {
		t.Error("Parent's first child should be the added child")
	}
}

func TestTreeNodeDepth(t *testing.T) {
	t.Parallel()

	root := NewTreeNode("root", "Root")
	child := NewTreeNode("child", "Child")
	grandchild := NewTreeNode("grandchild", "Grandchild")

	child.AddChild(grandchild)
	root.AddChild(child)

	assertTreeNodeDepth(t, root, child, grandchild)
}

func TestTreeColorModeNever(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()
	renderer.SetColorMode(ColorModeNever)

	root := NewTreeNode("root", "Root")
	root.Metadata["key"] = "value"
	renderer.SetRoot(root)

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, "\x1b[") {
		t.Errorf("ColorModeNever should produce no ANSI codes, got: %q", got)
	}

	assertContains(t, got, "Root", "should contain label even without colors")
	assertContains(t, got, "key: value", "should contain metadata even without colors")
}

func TestTreeColorModeAlways(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()
	renderer.SetColorMode(ColorModeAlways)

	root := NewTreeNode("root", "Root")
	child := NewTreeNode("child", "Child")
	root.AddChild(child)
	renderer.SetRoot(root)

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("ColorModeAlways should produce ANSI codes, got: %q", got)
	}

	assertContains(t, got, "Root", "should contain label")
	assertContains(t, got, "Child", "should contain child label")
}

func TestTreeColorModeDefault(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()
	if renderer.colorMode != ColorModeAuto {
		t.Errorf("default ColorMode = %v, want %v", renderer.colorMode, ColorModeAuto)
	}
}

func TestTreeColoredMetadata(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()
	renderer.SetColorMode(ColorModeAlways)

	node := NewTreeNode("node", "Node")
	node.Metadata["count"] = "42"
	node.Metadata["status"] = "active"
	renderer.SetRoot(node)

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "count: 42", "should contain metadata")
	assertContains(t, got, "status: active", "should contain metadata")
}
