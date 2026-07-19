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

//nolint:gochecknoinits // Registers DOT format capabilities and TableRenderer.
func init() {
	output.RegisterFormatShapes(output.FormatDOT, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableMarshaler(output.FormatDOT, renderDOTTable)
}

func renderDOTTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	out, err := NewDOTFromTable(data).Render()
	if err != nil {
		return fmt.Errorf("render DOT: %w", err)
	}

	return output.WriteRendered(w, "DOT", out)
}

// DOTRenderer implements the GraphRenderer interface for DOT/Graphviz output.
// defaultNodeSep is the default minimum space between adjacent nodes in the same rank.
const defaultNodeSep = "0.5"

// defaultRankSep is the default minimum space between two consecutive ranks.
const defaultRankSep = "0.5"

type DOTRenderer struct {
	output.GraphBuilder

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
		GraphBuilder: *output.NewGraphBuilder(),
		directed:     directed,
		graphID:      "G",
		rankdir:      RankDirTB,
		splines:      SplineOrtho,
		nodesep:      defaultNodeSep,
		ranksep:      defaultRankSep,
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

// SetDirection sets the graph layout direction from the canonical
// output.Direction enum. This bridges D2 and DOT vocabulary through a single
// type — prefer this over SetRankDir when the direction value originates from
// shared code or user input that is format-agnostic.
func (r *DOTRenderer) SetDirection(d output.Direction) *DOTRenderer {
	r.rankdir = RankDir(d.ToRankDir())
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
	return renderDOTString(r.Nodes(), r.Edges(), dotConfig{
		directed: r.directed,
		graphID:  r.graphID,
		rankdir:  r.rankdir,
		splines:  r.splines,
		nodesep:  r.nodesep,
		ranksep:  r.ranksep,
	}), nil
}

// renderDOTString builds the complete DOT document from nodes, edges, and config.
// This is the single source of truth for DOT formatting — shared by both
// WriteDOT (CQRS path) and DOTRenderer.Render() (legacy path).
func renderDOTString(nodes []output.GraphNode, edges []output.GraphEdge, cfg dotConfig) string {
	var b strings.Builder

	if cfg.directed {
		b.WriteString("digraph ")
	} else {
		b.WriteString("graph ")
	}

	b.WriteString(cfg.graphID)
	b.WriteString(" {\n")

	b.WriteString("  // Graph attributes\n")
	fmt.Fprintf(&b, "  rankdir=%s;\n", cfg.rankdir.String())
	fmt.Fprintf(&b, "  splines=%s;\n", cfg.splines.String())
	fmt.Fprintf(&b, "  nodesep=%s;\n", cfg.nodesep)
	fmt.Fprintf(&b, "  ranksep=%s;\n\n", cfg.ranksep)

	b.WriteString("  // Default node attributes\n")
	b.WriteString("  node [\n")
	b.WriteString("    shape=box\n")
	b.WriteString("    fontname=\"Helvetica\"\n")
	b.WriteString("    fontsize=12\n")
	b.WriteString("  ];\n\n")

	b.WriteString("  // Nodes\n")

	for _, node := range nodes {
		writeDOTNodeStmt(&b, node)
	}

	b.WriteString("\n  // Edges\n")

	for _, edge := range edges {
		writeDOTEdgeStmt(&b, edge, cfg.directed)
	}

	b.WriteString("}\n")

	return b.String()
}

func writeDOTNodeStmt(b *strings.Builder, node output.GraphNode) {
	b.WriteString("  \"")
	b.WriteString(escape.DOT(node.ID.Get()))
	b.WriteString("\" [\n")

	b.WriteString("    label=\"")
	b.WriteString(escape.DOT(node.Label.Get()))
	b.WriteString("\"\n")

	writeDOTAttrStmt(b, "shape", string(node.Shape), node.Shape != "")
	writeDOTAttrStmt(b, "fillcolor", node.Style.Fill, node.Style.Fill != "")
	writeDOTAttrStmt(b, "color", node.Style.Stroke, node.Style.Stroke != "")

	b.WriteString("  ];\n")
}

func writeDOTAttrStmt(b *strings.Builder, attrName, attrValue string, condition bool) {
	if condition {
		b.WriteString("    ")
		b.WriteString(attrName)
		b.WriteString("=\"")
		b.WriteString(escape.DOT(attrValue))
		b.WriteString("\"\n")
	}
}

func writeDOTEdgeStmt(b *strings.Builder, edge output.GraphEdge, directed bool) {
	op := "->"
	if !directed {
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

	if edge.Style.Line != "" {
		attrs = append(attrs, "style="+escape.DOT(edge.Style.Line.String()))
	}

	if len(attrs) > 0 {
		b.WriteString(" [\n    ")
		b.WriteString(strings.Join(attrs, "\n    "))
		b.WriteString("\n  ]")
	}

	b.WriteString(";\n")
}

// NewDOTFromTable converts Table to a DOT graph.
func NewDOTFromTable(data *output.Table) *DOTRenderer {
	renderer := NewDOTRenderer()
	if data == nil {
		return renderer
	}

	renderer.SetNodesFromTable(data, nil)

	return renderer
}

// NewDOTFromTree converts a TreeNode to DOT format.
func NewDOTFromTree(root *output.TreeNode) *DOTRenderer {
	return output.TreeToRenderer(NewDOTRenderer, (*DOTRenderer).addTreeNodes, root)
}

func dotTreeNodeID(node *output.TreeNode) string {
	if !node.ID.IsZero() {
		return node.ID.Get()
	}

	return escape.SlugifyID(node.Label.Get())
}

func (r *DOTRenderer) addTreeNodes(node *output.TreeNode, parentID output.TreeNodeID) {
	output.AddTreeNodes(&r.GraphBuilder, node, parentID.Get(), dotTreeNodeID, "")
}
