package output

import (
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
	output := renderer.Render()

	assertContains(t, output, "Node with Meta", "Output should contain node label")
	assertContains(t, output, "key: value", "Output should contain metadata")
}

func TestTreeRendererFromTableData(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"Name", "Age", "City"})
	data.AddRow([]string{"Alice", "30", "NYC"})
	data.AddRow([]string{"Bob", "25", "LA"})

	renderer := TreeRendererFromTableData(data)
	output := renderer.Render()

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

	testTreeNodeDepth(t, root, child, grandchild)
}
