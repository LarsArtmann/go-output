package output

import (
	"strings"
	"testing"
)

func TestASCIITreeRenderer(t *testing.T) {
	t.Parallel()
	renderer := NewASCIITreeRenderer()

	// Test empty tree
	if got := renderer.Render(); got != "" {
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
	output := renderer.Render()

	// Verify output contains expected elements
	if !strings.Contains(output, "Root") {
		t.Error("Output should contain 'Root'")
	}
	if !strings.Contains(output, "Child 1") {
		t.Error("Output should contain 'Child 1'")
	}
	if !strings.Contains(output, "Child 2") {
		t.Error("Output should contain 'Child 2'")
	}
	if !strings.Contains(output, "Leaf 1") {
		t.Error("Output should contain 'Leaf 1'")
	}
}

func TestASCIITreeRendererWithMetadata(t *testing.T) {
	t.Parallel()
	renderer := NewASCIITreeRenderer()

	node := NewTreeNode("node", "Node with Meta")
	node.Metadata["key"] = "value"
	node.Metadata["count"] = "42"

	renderer.SetRoot(node)
	output := renderer.Render()

	if !strings.Contains(output, "Node with Meta") {
		t.Error("Output should contain node label")
	}
	if !strings.Contains(output, "key: value") {
		t.Error("Output should contain metadata")
	}
}

func TestTreeRendererFromTableData(t *testing.T) {
	t.Parallel()
	data := NewTableData([]string{"Name", "Age", "City"})
	data.AddRow([]string{"Alice", "30", "NYC"})
	data.AddRow([]string{"Bob", "25", "LA"})

	renderer := TreeRendererFromTableData(data)
	output := renderer.Render()

	if !strings.Contains(output, "Data") {
		t.Error("Output should contain 'Data' as root")
	}
	if !strings.Contains(output, "Headers") {
		t.Error("Output should contain 'Headers' section")
	}
	if !strings.Contains(output, "Rows") {
		t.Error("Output should contain 'Rows' section")
	}
	if !strings.Contains(output, "Alice") {
		t.Error("Output should contain 'Alice' row")
	}
}

func TestNewTreeNode(t *testing.T) {
	t.Parallel()
	node := NewTreeNode("id1", "Label 1")

	if node.ID != "id1" {
		t.Errorf("Expected ID 'id1', got '%s'", node.ID)
	}
	if node.Label != "Label 1" {
		t.Errorf("Expected Label 'Label 1', got '%s'", node.Label)
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

	if root.Depth() != 0 {
		t.Errorf("Root depth should be 0, got %d", root.Depth())
	}
	if child.Depth() != 0 {
		t.Errorf("Child depth should be 0 (no parent pointer), got %d", child.Depth())
	}
}
