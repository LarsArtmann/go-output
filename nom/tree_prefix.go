package nom
import (
	"strings"
)
// buildTreePrefix builds the tree structure prefix (├──, └──, etc.)
func (dt *DependencyTree) buildTreePrefix(node *TreeNode, displayNodes []*TreeNode) string {
	if node.IsRoot {
		return ""
	}
	// Build prefix based on depth
	var prefixBuilder strings.Builder
	// Add indentation for each level
	for range node.Depth - 1 {
		prefixBuilder.WriteString("│   ")
	}
	// Determine if this is the last child of its parent
	if dt.isLastChild(node, displayNodes) {
		prefixBuilder.WriteString("└── ")
	} else {
		prefixBuilder.WriteString("├── ")
	}
	return prefixBuilder.String()
}
// isLastChild determines if this node is the last child of its parent.
func (dt *DependencyTree) isLastChild(node *TreeNode, displayNodes []*TreeNode) bool {
	if node.Parent == nil {
		return true
	}
	// Find the last sibling that's displayed
	var lastDisplayedSibling *TreeNode
	for _, sibling := range node.Parent.Children {
		for _, displayed := range displayNodes {
			if sibling.ActivityID == displayed.ActivityID {
				lastDisplayedSibling = sibling
			}
		}
	}
	return lastDisplayedSibling == node
}
