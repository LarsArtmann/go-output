package output

import (
	"strconv"
	"strings"
)

// ASCIITreeRenderer implements the TreeOutputRenderer interface for ASCII tree output.
type ASCIITreeRenderer struct {
	root *TreeNode
}

// NewASCIITreeRenderer creates a new ASCIITreeRenderer.
func NewASCIITreeRenderer() *ASCIITreeRenderer {
	return &ASCIITreeRenderer{} //nolint:exhaustruct // root and builder are initialized lazily
}

// Ensure ASCIITreeRenderer implements TreeOutputRenderer.
var _ TreeOutputRenderer = (*ASCIITreeRenderer)(nil)

// SetRoot sets the root node of the tree.
func (r *ASCIITreeRenderer) SetRoot(node *TreeNode) {
	r.root = node
}

// Render returns the tree as a string.
func (r *ASCIITreeRenderer) Render() string {
	if r.root == nil {
		return ""
	}

	var b strings.Builder
	r.renderNode(&b, r.root, "", true)

	return b.String()
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
