package graph

import (
	"fmt"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

// Compile-time interface checks.
var (
	_ output.Renderer      = (*DOTRenderer)(nil)
	_ output.GraphRenderer = (*DOTRenderer)(nil)
)

// DOTRenderer implements the GraphRenderer interface for DOT/Graphviz output.
type DOTRenderer struct {
	output.GraphRendererMixin

	directed bool
	graphID  string
}

// newDOTRenderer creates a new DOTRenderer with the specified direction.
func newDOTRenderer(directed bool) *DOTRenderer {
	return &DOTRenderer{
		GraphRendererMixin: output.NewGraphRendererMixin(),
		directed:           directed,
		graphID:            "G",
	}
}

// NewDOTRenderer creates a new DOTRenderer for directed graphs.
func NewDOTRenderer() *DOTRenderer {
	return newDOTRenderer(true)
}

// NewUndirectedDOTRenderer creates a new DOTRenderer for undirected graphs.
func NewUndirectedDOTRenderer() *DOTRenderer {
	return newDOTRenderer(false)
}

// SetGraphID sets the graph ID.
func (r *DOTRenderer) SetGraphID(id string) {
	r.graphID = id
}

// Render returns the DOT graph as a string.
func (r *DOTRenderer) Render() (string, error) {
	var b strings.Builder

	if r.directed {
		b.WriteString("digraph ")
	} else {
		b.WriteString("graph ")
	}

	b.WriteString(r.graphID)
	b.WriteString(" {\n")

	b.WriteString("  // Graph attributes\n")
	b.WriteString("  rankdir=TB;\n")
	b.WriteString("  splines=ortho;\n")
	b.WriteString("  nodesep=0.5;\n")
	b.WriteString("  ranksep=0.5;\n\n")

	b.WriteString("  // Default node attributes\n")
	b.WriteString("  node [\n")
	b.WriteString("    shape=box\n")
	b.WriteString("    fontname=\"Helvetica\"\n")
	b.WriteString("    fontsize=12\n")
	b.WriteString("  ];\n\n")

	b.WriteString("  // Nodes\n")

	for _, node := range r.Nodes() {
		r.writeNode(&b, node)
	}

	b.WriteString("\n  // Edges\n")

	for _, edge := range r.Edges() {
		r.writeEdge(&b, edge)
	}

	b.WriteString("}\n")

	return b.String(), nil
}

func (r *DOTRenderer) writeNode(b *strings.Builder, node output.GraphNode) {
	b.WriteString("  \"")
	b.WriteString(escape.DOT(node.ID.Get()))
	b.WriteString("\" [\n")

	b.WriteString("    label=\"")
	b.WriteString(escape.DOT(node.Label.Get()))
	b.WriteString("\"\n")

	r.writeNodeAttr(b, "shape", string(node.Shape), node.Shape != "")
	r.writeNodeAttr(b, "fillcolor", node.Style.FillColor, node.Style.FillColor != "")
	r.writeNodeAttr(b, "color", node.Style.StrokeColor, node.Style.StrokeColor != "")

	b.WriteString("  ];\n")
}

func (r *DOTRenderer) writeNodeAttr(
	b *strings.Builder,
	attrName, attrValue string,
	condition bool,
) {
	if condition {
		b.WriteString("    ")
		b.WriteString(attrName)
		b.WriteString("=")
		b.WriteString(attrValue)
		b.WriteString("\n")
	}
}

func (r *DOTRenderer) writeEdge(b *strings.Builder, edge output.GraphEdge) {
	op := "->"
	if !r.directed {
		op = "--"
	}

	b.WriteString("  \"")
	b.WriteString(escape.DOT(edge.From.Get()))
	b.WriteString("\" ")
	b.WriteString(op)
	b.WriteString(" \"")
	b.WriteString(escape.DOT(edge.To.Get()))
	b.WriteString("\"")

	attrs := make([]string, 0)

	if !edge.Label.IsZero() {
		attrs = append(attrs, fmt.Sprintf("label=\"%s\"", escape.DOT(edge.Label.Get())))
	}

	if edge.Style.Color != "" {
		attrs = append(attrs, "color="+edge.Style.Color)
	}

	if edge.Style.Style != "" {
		attrs = append(attrs, "style="+edge.Style.Style)
	}

	if len(attrs) > 0 {
		b.WriteString(" [\n    ")
		b.WriteString(strings.Join(attrs, "\n    "))
		b.WriteString("\n  ]")
	}

	b.WriteString(";\n")
}

// DOTFromTableData converts TableData to a DOT graph.
func DOTFromTableData(data *output.TableData) *DOTRenderer {
	renderer := NewDOTRenderer()
	if data == nil {
		return renderer
	}

	renderer.SetNodesFromTableData(data, func(_ int, n *output.GraphNode) {
		n.Label = output.NewBrandedID[output.GraphNodeLabelBrand](escape.DOT(n.Label.Get()))
	})

	return renderer
}

// DOTFromTree converts a TreeNode to DOT format.
func DOTFromTree(root *output.TreeNode) *DOTRenderer {
	renderer := NewDOTRenderer()
	if root == nil {
		return renderer
	}

	renderer.addTreeNodes(root, output.NewBrandedID[output.TreeNodeIDBrand](""))

	return renderer
}

func dotTreeNodeID(node *output.TreeNode) string {
	if !node.ID.IsZero() {
		return node.ID.Get()
	}

	return strings.ReplaceAll(node.Label.Get(), " ", "_")
}

func (r *DOTRenderer) addTreeNodes(node *output.TreeNode, parentID output.TreeNodeID) {
	output.AddTreeNodes(&r.GraphRendererMixin, node, parentID.Get(), dotTreeNodeID, "")
}
