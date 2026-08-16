// Package tree renders Table as an ASCII tree.
//
// It is an optional format renderer: import it to activate Tree output
// through output.RenderTable, or use NewASCIITreeRenderer directly.
//
//	import "github.com/larsartmann/go-output/tree"
//
//	r := tree.NewASCIITreeRenderer()
//	r.SetRoot(rootNode)
//	out, _ := r.Render()
package tree

import (
	"sort"
	"strconv"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

// Local ANSI color constants for depth-based tree coloring. The Reset and Bold
// sequences come from the shared escape package; Cyan, Blue, Green, and Magenta
// are tree-specific and remain local since tree is the only consumer.
const (
	ansiCyan    = "\033[36m"
	ansiBlue    = "\033[34m"
	ansiGreen   = "\033[32m"
	ansiMagenta = "\033[35m"
)

//nolint:gochecknoglobals // Color cycle lookup table for depth-based tree coloring.
var depthColors = []string{ansiGreen, ansiBlue, ansiMagenta, ansiCyan}

// ASCIITreeRenderer implements the output.TreeRenderer interface for ASCII tree output.
type ASCIITreeRenderer struct {
	root      *output.TreeNode
	colorMode output.ColorMode
}

// NewASCIITreeRenderer creates a new ASCIITreeRenderer.
func NewASCIITreeRenderer() *ASCIITreeRenderer {
	return &ASCIITreeRenderer{colorMode: output.ColorModeAuto} //nolint:exhaustruct // root is initialized lazily
}

// SetColorMode sets the color mode for the tree renderer.
func (r *ASCIITreeRenderer) SetColorMode(mode output.ColorMode) {
	r.colorMode = mode
}

// Compile-time interface checks.
var (
	_ output.Renderer     = (*ASCIITreeRenderer)(nil)
	_ output.TreeRenderer = (*ASCIITreeRenderer)(nil)
)

// SetRoot sets the root node of the tree.
func (r *ASCIITreeRenderer) SetRoot(node *output.TreeNode) {
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
	node *output.TreeNode,
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
		b.WriteString(escape.ANSIDim)
	}

	b.WriteString(connector)

	if r.useColor() {
		b.WriteString(escape.ANSIReturn)
		b.WriteString(r.colorForDepth(depth))
		b.WriteString(escape.ANSIBold)
	}

	b.WriteString(node.Label.Get())

	if r.useColor() {
		b.WriteString(escape.ANSIReturn)
	}

	if len(node.Metadata) > 0 {
		metaKeys := make([]string, 0, len(node.Metadata))
		for k := range node.Metadata {
			metaKeys = append(metaKeys, k)
		}

		sort.Strings(metaKeys)

		metaParts := make([]string, 0, len(node.Metadata))
		for _, k := range metaKeys {
			metaParts = append(metaParts, k+": "+node.Metadata[k])
		}

		if r.useColor() {
			b.WriteString(" ")
			b.WriteString(escape.ANSIDim)
			b.WriteString(ansiCyan)
		} else {
			b.WriteString(" ")
		}

		b.WriteString("(")
		b.WriteString(strings.Join(metaParts, ", "))
		b.WriteString(")")

		if r.useColor() {
			b.WriteString(escape.ANSIReturn)
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

// TreeRendererFromTable converts Table to a tree using the first column as hierarchy.
func TreeRendererFromTable(data *output.Table) *ASCIITreeRenderer {
	renderer := NewASCIITreeRenderer()
	if data == nil || len(data.Rows) == 0 {
		return renderer
	}

	// Build a simple tree from the data
	root := output.NewTreeNode("root", "Data")

	if len(data.Headers) > 0 {
		headerNode := output.NewTreeNode("headers", "Headers")
		for _, h := range data.Headers {
			headerNode.AddChild(output.NewTreeNode(h, h))
		}

		root.AddChild(headerNode)
	}

	rowsNode := output.NewTreeNode("rows", "Rows")
	for i, row := range data.Rows {
		rowNode := output.NewTreeNode("row-"+strconv.Itoa(i), "Row "+strconv.Itoa(i+1))

		for j, cell := range row {
			var headerName string
			if j < len(data.Headers) {
				headerName = data.Headers[j]
			} else {
				headerName = "Col " + strconv.Itoa(j)
			}

			rowNode.AddChild(output.NewTreeNode(headerName, cell))
		}

		rowsNode.AddChild(rowNode)
	}

	root.AddChild(rowsNode)
	renderer.SetRoot(root)

	return renderer
}
