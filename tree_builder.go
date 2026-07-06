package output

// TreeBuilder is the CQRS write-side builder for tree data.
// It provides a fluent construction API for assembling a tree,
// then returns the assembled *TreeNode (root) via Build().
// Note: *TreeNode has exported fields and an AddChild method; callers
// SHOULD treat the result as read-only after Build(), but Go does not
// enforce this at the type level.
//
// Usage:
//
//	root := NewTreeBuilder().
//	    SetRoot("build", "Build").
//	    AddChild("build", "compile", "Compile").
//	    AddChild("build", "lint", "Lint").
//	    AddChild("compile", "test", "Test").
//	    Build()
type TreeBuilder struct {
	root  *TreeNode
	nodes map[string]*TreeNode
}

// NewTreeBuilder creates a new TreeBuilder.
func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{
		nodes: make(map[string]*TreeNode),
	}
}

// SetRoot creates and sets the root node.
func (b *TreeBuilder) SetRoot(id, label string) *TreeBuilder {
	b.root = NewTreeNode(id, label)
	b.nodes[id] = b.root

	return b
}

// AddChild adds a child node under the specified parent ID.
// If the parent ID is not found, the child is silently skipped.
func (b *TreeBuilder) AddChild(parentID, id, label string) *TreeBuilder {
	parent, ok := b.nodes[parentID]
	if !ok {
		return b
	}

	child := NewTreeNode(id, label)
	parent.AddChild(child)
	b.nodes[id] = child

	return b
}

// AddChildren adds multiple child nodes under the specified parent ID.
// Each child is specified as an (id, label) pair. If the parent ID is not
// found, all children are silently skipped.
func (b *TreeBuilder) AddChildren(parentID string, children ...[2]string) *TreeBuilder {
	parent, ok := b.nodes[parentID]
	if !ok {
		return b
	}

	for _, c := range children {
		child := NewTreeNode(c[0], c[1])
		parent.AddChild(child)
		b.nodes[c[0]] = child
	}

	return b
}

// Build returns the root node of the assembled tree.
// Returns nil if SetRoot was never called.
func (b *TreeBuilder) Build() *TreeNode {
	return b.root
}
