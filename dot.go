package output

import (
	"fmt"
	"strings"
)

// GraphRendererMixin contains shared fields and methods for graph renderers.
type GraphRendererMixin struct {
	nodes []GraphNode
	edges []GraphEdge
}

// SetNodes sets the graph nodes.
func (m *GraphRendererMixin) SetNodes(nodes []GraphNode) {
	m.nodes = nodes
}

// SetEdges sets the graph edges.
func (m *GraphRendererMixin) SetEdges(edges []GraphEdge) {
	m.edges = edges
}

// AddRowEdges adds edges from data.CreateRowEdges() to the graph.
func (m *GraphRendererMixin) AddRowEdges(data *TableData) {
	for _, edge := range data.CreateRowEdges() {
		//nolint:exhaustruct // Uses defaults for optional fields
		m.edges = append(m.edges, GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand](edge.From),
			To:   NewBrandedID[GraphNodeIDBrand](edge.To),
		})
	}
}

// DOTRenderer implements the GraphRenderer interface for DOT/Graphviz output.
type DOTRenderer struct {
	GraphRendererMixin

	directed bool
	graphID  string
}

// newDOTRenderer creates a new DOTRenderer with the specified direction.
func newDOTRenderer(directed bool) *DOTRenderer {
	return &DOTRenderer{
		GraphRendererMixin: GraphRendererMixin{
			nodes: make([]GraphNode, 0),
			edges: make([]GraphEdge, 0),
		},
		directed: directed,
		graphID:  "G",
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
func (r *DOTRenderer) Render() string {
	var b strings.Builder

	// Write header
	if r.directed {
		b.WriteString("digraph ")
	} else {
		b.WriteString("graph ")
	}

	b.WriteString(r.graphID)
	b.WriteString(" {\n")

	// Graph attributes
	b.WriteString("  // Graph attributes\n")
	b.WriteString("  rankdir=TB;\n")
	b.WriteString("  splines=ortho;\n")
	b.WriteString("  nodesep=0.5;\n")
	b.WriteString("  ranksep=0.5;\n\n")

	// Default node attributes
	b.WriteString("  // Default node attributes\n")
	b.WriteString("  node [\n")
	b.WriteString("    shape=box\n")
	b.WriteString("    fontname=\"Helvetica\"\n")
	b.WriteString("    fontsize=12\n")
	b.WriteString("  ];\n\n")

	// Write nodes
	b.WriteString("  // Nodes\n")

	for _, node := range r.nodes {
		r.writeNode(&b, node)
	}

	// Write edges
	b.WriteString("\n  // Edges\n")

	for _, edge := range r.edges {
		r.writeEdge(&b, edge)
	}

	b.WriteString("}\n")

	return b.String()
}

func (r *DOTRenderer) writeNode(b *strings.Builder, node GraphNode) {
	b.WriteString("  \"")
	b.WriteString(r.escapeDOT(node.ID.Get()))
	b.WriteString("\" [\n")

	b.WriteString("    label=\"")
	b.WriteString(r.escapeDOT(node.Label.Get()))
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

func (r *DOTRenderer) writeEdge(b *strings.Builder, edge GraphEdge) {
	op := "->"
	if !r.directed {
		op = "--"
	}

	b.WriteString("  \"")
	b.WriteString(r.escapeDOT(edge.From.Get()))
	b.WriteString("\" ")
	b.WriteString(op)
	b.WriteString(" \"")
	b.WriteString(r.escapeDOT(edge.To.Get()))
	b.WriteString("\"")

	attrs := make([]string, 0)

	if !edge.Label.IsEmpty() {
		attrs = append(attrs, fmt.Sprintf("label=\"%s\"", r.escapeDOT(edge.Label.Get())))
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

func (r *DOTRenderer) escapeDOT(s string) string {
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")

	return s
}

// DOTFromTableData converts TableData to a DOT graph.
func DOTFromTableData(data *TableData) *DOTRenderer {
	renderer := NewDOTRenderer()
	if data == nil {
		return renderer
	}

	// Create nodes for each row
	for i, row := range data.Rows {
		var labelParts []string

		for j, cell := range row {
			if j < len(data.Headers) {
				labelParts = append(labelParts, fmt.Sprintf("%s: %s", data.Headers[j], cell))
			} else {
				labelParts = append(labelParts, cell)
			}
		}

		label := strings.Join(labelParts, "\\n")
		//nolint:exhaustruct // Uses defaults for optional fields
		renderer.nodes = append(renderer.nodes, GraphNode{
			ID:    NewBrandedID[GraphNodeIDBrand](fmt.Sprintf("row%d", i)),
			Label: NewBrandedID[GraphNodeLabelBrand](label),
		})
	}

	renderer.AddRowEdges(data)

	return renderer
}

// DOTFromTree converts a TreeNode to DOT format.
func DOTFromTree(root *TreeNode) *DOTRenderer {
	renderer := NewDOTRenderer()
	if root == nil {
		return renderer
	}

	renderer.addTreeNodes(root, NewBrandedID[TreeNodeIDBrand](""))

	return renderer
}

func (r *DOTRenderer) addTreeNodes(node *TreeNode, parentID TreeNodeID) {
	nodeID := node.ID
	if nodeID.IsEmpty() {
		nodeID = NewBrandedID[TreeNodeIDBrand](strings.ReplaceAll(node.Label.Get(), " ", "_"))
	}

	graphNodeID := NewBrandedID[GraphNodeIDBrand](nodeID.Get())
	graphNodeLabel := NewBrandedID[GraphNodeLabelBrand](node.Label.Get())

	//nolint:exhaustruct // Uses defaults for optional fields
	r.nodes = append(r.nodes, GraphNode{
		ID:    graphNodeID,
		Label: graphNodeLabel,
	})

	if !parentID.IsEmpty() {
		//nolint:exhaustruct // Uses defaults for optional fields
		r.edges = append(r.edges, GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand](parentID.Get()),
			To:   graphNodeID,
		})
	}

	for _, child := range node.Children {
		r.addTreeNodes(child, nodeID)
	}
}
