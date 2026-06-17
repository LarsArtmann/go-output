package graph

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

// Compile-time interface checks.
var (
	_ output.Renderer      = (*DOTRenderer)(nil)
	_ output.GraphRenderer = (*DOTRenderer)(nil)
)

//nolint:gochecknoinits // Registers DOT format capabilities and TableDataMarshaler.
func init() {
	output.RegisterFormatShapes(output.FormatDOT, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableDataMarshaler(output.FormatDOT, renderDOTTableData)
}

func renderDOTTableData(w io.Writer, data *output.TableData, _ output.RenderOptions) error {
	out, err := DOTFromTableData(data).Render()
	if err != nil {
		return fmt.Errorf("render DOT: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write DOT output: %w", err)
	}

	return nil
}

// DOTRenderer implements the GraphRenderer interface for DOT/Graphviz output.
type DOTRenderer struct {
	output.GraphRendererState

	directed bool
	graphID  string
	rankdir  RankDir
	splines  SplineStyle
	nodesep  string
	ranksep  string
}

// newDOTRenderer creates a new DOTRenderer with the specified direction.
func newDOTRenderer(directed bool) *DOTRenderer {
	return &DOTRenderer{
		GraphRendererState: output.NewGraphRendererState(),
		directed:           directed,
		graphID:            "G",
		rankdir:            RankDirTB,
		splines:            SplineOrtho,
		nodesep:            "0.5",
		ranksep:            "0.5",
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

// SetRankDir sets the graph layout direction (TB, LR, BT, RL).
func (r *DOTRenderer) SetRankDir(direction RankDir) *DOTRenderer {
	r.rankdir = direction
	return r
}

// SetSplines sets the edge routing style (ortho, spline, polyline, line, curved, none).
func (r *DOTRenderer) SetSplines(style SplineStyle) *DOTRenderer {
	r.splines = style
	return r
}

// isValidNumericSep reports whether s is a valid DOT nodesep/ranksep value:
// a non-negative number (int or float) without any characters that could
// inject attributes or statements.
func isValidNumericSep(s string) bool {
	if s == "" {
		return false
	}

	hasDigit := false

	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			hasDigit = true
		case c == '.':
			// allowed — decimal point
		default:
			return false
		}
	}

	return hasDigit
}

// SetNodeSep sets the minimum space between two adjacent nodes in the same rank.
// The value must be numeric (e.g., "0.5", "2") to prevent DOT injection.
func (r *DOTRenderer) SetNodeSep(sep string) *DOTRenderer {
	if isValidNumericSep(sep) {
		r.nodesep = sep
	}

	return r
}

// SetRankSep sets the minimum space between two consecutive ranks.
// The value must be numeric (e.g., "0.5", "2") to prevent DOT injection.
func (r *DOTRenderer) SetRankSep(sep string) *DOTRenderer {
	if isValidNumericSep(sep) {
		r.ranksep = sep
	}

	return r
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
	fmt.Fprintf(&b, "  rankdir=%s;\n", r.rankdir.String())
	fmt.Fprintf(&b, "  splines=%s;\n", r.splines.String())
	fmt.Fprintf(&b, "  nodesep=%s;\n", r.nodesep)
	fmt.Fprintf(&b, "  ranksep=%s;\n\n", r.ranksep)

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
		attrs = append(attrs, "color="+escape.DOT(edge.Style.Color))
	}

	if edge.Style.Style != "" {
		attrs = append(attrs, "style="+escape.DOT(edge.Style.Style))
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

	renderer.SetNodesFromTableData(data, nil)

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

	return escape.SlugifyID(node.Label.Get())
}

func (r *DOTRenderer) addTreeNodes(node *output.TreeNode, parentID output.TreeNodeID) {
	output.AddTreeNodes(&r.GraphRendererState, node, parentID.Get(), dotTreeNodeID, "")
}
