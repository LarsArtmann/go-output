package output

import (
	"testing"
)

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
