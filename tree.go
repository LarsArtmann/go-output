package output

import (
	"strings"
)

// ASCIITreeRenderer implements the TreeOutputRenderer interface for ASCII tree output.
type ASCIITreeRenderer struct {
	root    *TreeNode
	builder strings.Builder
}

// NewASCIITreeRenderer creates a new ASCIITreeRenderer.
func NewASCIITreeRenderer() *ASCIITreeRenderer {
	return &ASCIITreeRenderer{}
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
	r.builder.Reset()
	r.renderNode(r.root, "", true)
	return r.builder.String()
}

func (r *ASCIITreeRenderer) renderNode(node *TreeNode, prefix string, isLast bool) {
	// Determine connector characters
	var connector string
	if isLast {
		connector = "└── "
	} else {
		connector = "├── "
	}

	// Write this node
	r.builder.WriteString(prefix)
	r.builder.WriteString(connector)
	r.builder.WriteString(node.Label)

	// Add metadata summary if present
	if len(node.Metadata) > 0 {
		var metaParts []string
		for k, v := range node.Metadata {
			metaParts = append(metaParts, k+": "+v)
		}
		r.builder.WriteString(" (")
		r.builder.WriteString(strings.Join(metaParts, ", "))
		r.builder.WriteString(")")
	}

	r.builder.WriteString("\n")

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
		r.renderNode(child, childPrefix, isLastChild)
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
		rowNode := NewTreeNode("row-"+string(rune('0'+i)), "Row "+itoa(i+1))
		for j, cell := range row {
			var headerName string
			if j < len(data.Headers) {
				headerName = data.Headers[j]
			} else {
				headerName = "Col " + itoa(j)
			}
			rowNode.AddChild(NewTreeNode(headerName, cell))
		}
		rowsNode.AddChild(rowNode)
	}

	root.AddChild(rowsNode)
	renderer.SetRoot(root)
	return renderer
}

// itoa converts an int to string without importing strconv.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}
