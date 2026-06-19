package output

// TreeNode represents a node in a tree structure.
type TreeNode struct {
	// ID is the unique identifier for the node.
	ID TreeNodeID
	// Label is the display text for the node.
	Label TreeNodeLabel
	// Children holds the child nodes.
	Children []*TreeNode
	// Metadata holds arbitrary key-value pairs for custom data.
	Metadata map[string]string
	parent   *TreeNode
	depth    int // cached depth (root = 0), set in AddChild
}

// NewTreeNode creates a new TreeNode with the given ID and label.
func NewTreeNode(id, label string) *TreeNode {
	return &TreeNode{
		ID:       NewBrandedID[TreeNodeIDBrand](id),
		Label:    NewBrandedID[TreeNodeLabelBrand](label),
		Children: make([]*TreeNode, 0),
		Metadata: make(map[string]string),
		parent:   nil, // parent is set via AddChild
	}
}

// AddChild adds a child node to this node.
func (n *TreeNode) AddChild(child *TreeNode) {
	child.parent = n
	n.Children = append(n.Children, child)
	child.propagateDepth(n.depth + 1)
}

// propagateDepth sets the depth for this node and all descendants recursively.
// Called from AddChild to handle subtrees that change depth when re-parented.
func (n *TreeNode) propagateDepth(d int) {
	n.depth = d
	for _, c := range n.Children {
		c.propagateDepth(d + 1)
	}
}

// Depth returns the depth of this node in the tree (root = 0).
func (n *TreeNode) Depth() int {
	return n.depth
}

// Parent returns the parent node (nil for root).
func (n *TreeNode) Parent() *TreeNode {
	return n.parent
}

// TreeOutputRenderer defines the interface for tree format renderers.
type TreeOutputRenderer interface {
	Renderer
	// SetRoot sets the root node of the tree.
	SetRoot(node *TreeNode)
}
