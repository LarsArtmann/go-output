package output

import (
	"strconv"
	"strings"
)

// TreeOutputRenderer defines the interface for tree format renderers.
type TreeOutputRenderer interface {
	Renderer
	// SetRoot sets the root node of the tree.
	SetRoot(node *TreeNode)
}

// TreeNode represents a node in a tree structure.
type TreeNode struct {
	ID       TreeNodeID
	Label    TreeNodeLabel
	Children []*TreeNode
	Metadata map[string]string
	parent   *TreeNode
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
}

// Depth returns the depth of this node in the tree (root = 0).
func (n *TreeNode) Depth() int {
	depth := 0

	current := n
	for current.Parent() != nil {
		depth++
		current = current.Parent()
	}

	return depth
}

// Parent returns the parent node (nil for root).
func (n *TreeNode) Parent() *TreeNode {
	return n.parent
}

// ASCIITreeRenderer implements the TreeOutputRenderer interface for ASCII tree output.
type ASCIITreeRenderer struct {
	root *TreeNode
}

// NewASCIITreeRenderer creates a new ASCIITreeRenderer.
func NewASCIITreeRenderer() *ASCIITreeRenderer {
	return &ASCIITreeRenderer{} //nolint:exhaustruct // root and builder are initialized lazily
}

// Compile-time interface checks.
var (
	_ Renderer           = (*ASCIITreeRenderer)(nil)
	_ TreeOutputRenderer = (*ASCIITreeRenderer)(nil)
)

// SetRoot sets the root node of the tree.
func (r *ASCIITreeRenderer) SetRoot(node *TreeNode) {
	r.root = node
}

// Render returns the tree as a string.
func (r *ASCIITreeRenderer) Render() (string, error) {
	if r.root == nil {
		return "", nil
	}

	var b strings.Builder
	r.renderNode(&b, r.root, "", true)

	return b.String(), nil
}

func (r *ASCIITreeRenderer) renderNode(
	b *strings.Builder,
	node *TreeNode,
	prefix string,
	isLast bool,
) {
	// Determine connector characters
	var connector string
	if isLast {
		connector = "└── "
	} else {
		connector = "├── "
	}

	// Write this node
	b.WriteString(prefix)
	b.WriteString(connector)
	b.WriteString(node.Label.Get())

	// Add metadata summary if present
	if len(node.Metadata) > 0 {
		metaParts := make([]string, 0, len(node.Metadata))
		for k, v := range node.Metadata {
			metaParts = append(metaParts, k+": "+v)
		}

		b.WriteString(" (")
		b.WriteString(strings.Join(metaParts, ", "))
		b.WriteString(")")
	}

	b.WriteString("\n")

	// Prepare prefix for children
	var childPrefix string
	if isLast {
		childPrefix = prefix + "    "
	} else {
		childPrefix = prefix + "│   "
	}

	// Render children
	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		r.renderNode(b, child, childPrefix, isLastChild)
	}
}

// TreeRendererFromTableData converts TableData to a tree using the first column as hierarchy.
func TreeRendererFromTableData(data *TableData) *ASCIITreeRenderer {
	renderer := NewASCIITreeRenderer()
	if data == nil || len(data.Rows) == 0 {
		return renderer
	}

	// Build a simple tree from the data
	root := NewTreeNode("root", "Data")

	if len(data.Headers) > 0 {
		headerNode := NewTreeNode("headers", "Headers")
		for _, h := range data.Headers {
			headerNode.AddChild(NewTreeNode(h, h))
		}

		root.AddChild(headerNode)
	}

	rowsNode := NewTreeNode("rows", "Rows")
	for i, row := range data.Rows {
		rowNode := NewTreeNode("row-"+strconv.Itoa(i), "Row "+strconv.Itoa(i+1))

		for j, cell := range row {
			var headerName string
			if j < len(data.Headers) {
				headerName = data.Headers[j]
			} else {
				headerName = "Col " + strconv.Itoa(j)
			}

			rowNode.AddChild(NewTreeNode(headerName, cell))
		}

		rowsNode.AddChild(rowNode)
	}

	root.AddChild(rowsNode)
	renderer.SetRoot(root)

	return renderer
}
