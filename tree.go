package output

import (
	"strconv"
	"strings"
)

// ANSI escape codes for terminal coloring.
const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiCyan    = "\033[36m"
	ansiBlue    = "\033[34m"
	ansiGreen   = "\033[32m"
	ansiMagenta = "\033[35m"
)

var depthColors = []string{ansiGreen, ansiBlue, ansiMagenta, ansiCyan}

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
	root      *TreeNode
	colorMode ColorMode
}

// NewASCIITreeRenderer creates a new ASCIITreeRenderer.
func NewASCIITreeRenderer() *ASCIITreeRenderer {
	return &ASCIITreeRenderer{colorMode: ColorModeAuto} //nolint:exhaustruct // root is initialized lazily
}

// SetColorMode sets the color mode for the tree renderer.
func (r *ASCIITreeRenderer) SetColorMode(mode ColorMode) {
	r.colorMode = mode
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
	r.renderNode(&b, r.root, "", true, 0)

	return b.String(), nil
}

func (r *ASCIITreeRenderer) useColor() bool {
	return r.colorMode.ShouldColor()
}

func (r *ASCIITreeRenderer) colorForDepth(depth int) string {
	return depthColors[depth%len(depthColors)]
}

func (r *ASCIITreeRenderer) renderNode(
	b *strings.Builder,
	node *TreeNode,
	prefix string,
	isLast bool,
	depth int,
) {
	var connector string
	if isLast {
		connector = "└── "
	} else {
		connector = "├── "
	}

	b.WriteString(prefix)

	if r.useColor() {
		b.WriteString(ansiDim)
	}

	b.WriteString(connector)

	if r.useColor() {
		b.WriteString(ansiReset)
		b.WriteString(r.colorForDepth(depth))
		b.WriteString(ansiBold)
	}

	b.WriteString(node.Label.Get())

	if r.useColor() {
		b.WriteString(ansiReset)
	}

	if len(node.Metadata) > 0 {
		metaParts := make([]string, 0, len(node.Metadata))
		for k, v := range node.Metadata {
			metaParts = append(metaParts, k+": "+v)
		}

		if r.useColor() {
			b.WriteString(" ")
			b.WriteString(ansiDim)
			b.WriteString(ansiCyan)
		} else {
			b.WriteString(" ")
		}

		b.WriteString("(")
		b.WriteString(strings.Join(metaParts, ", "))
		b.WriteString(")")

		if r.useColor() {
			b.WriteString(ansiReset)
		}
	}

	b.WriteString("\n")

	var childPrefix string
	if isLast {
		childPrefix = prefix + "    "
	} else {
		childPrefix = prefix + "│   "
	}

	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		r.renderNode(b, child, childPrefix, isLastChild, depth+1)
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
