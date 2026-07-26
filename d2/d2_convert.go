package d2

import (
	"io"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

//nolint:gochecknoinits // Registers D2 TableRenderer for registry-based dispatch.
func init() {
	output.RegisterTableMarshaler(output.FormatD2, renderTable)
}

func renderTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return output.WriteRenderedFrom(w, NewD2FromTable(data).Render, "D2", "render D2")
}

// SetNodes sets graph nodes from the generic GraphNode type, satisfying GraphRenderer.
func (d *Diagram) SetNodes(nodes []output.GraphNode) {
	d.nodes = make([]Node, 0, len(nodes))
	for _, n := range nodes {
		d.nodes = append(d.nodes, graphNodeToD2(n))
	}
}

// SetEdges sets graph edges from the generic GraphEdge type, satisfying GraphRenderer.
func (d *Diagram) SetEdges(edges []output.GraphEdge) {
	d.edges = make([]Edge, 0, len(edges))
	for _, e := range edges {
		d.edges = append(d.edges, graphEdgeToD2(e))
	}
}

func graphNodeToD2(n output.GraphNode) Node {
	return Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand](n.ID.Get()),
		Label: output.NewBrandedID[output.D2NodeLabelBrand](n.Label.Get()),
		Shape: nodeShapeToD2(n.Shape),
		Style: graphStyleToD2(n.Style),
	}
}

func graphEdgeToD2(e output.GraphEdge) Edge {
	return Edge{
		From:  output.NewBrandedID[output.D2NodeIDBrand](e.From.Get()),
		To:    output.NewBrandedID[output.D2NodeIDBrand](e.To.Get()),
		Label: output.NewBrandedID[output.D2NodeLabelBrand](e.Label.Get()),
		Style: EdgeStyle{
			StrokeStyle: StrokeStyle{
				Stroke: e.Style.Color,
			},
		},
	}
}

func nodeShapeToD2(s output.NodeShape) NodeShape {
	switch s {
	case output.NodeShapeBox:
		return ShapeRectangle
	case output.NodeShapeEllipse:
		return ShapeOval
	case output.NodeShapeDiamond:
		return ShapeDiamond
	case output.NodeShapeCircle:
		return ShapeCircle
	case output.NodeShapeCylinder:
		return ShapeCylinder
	case output.NodeShapeHexagon:
		return ShapeHexagon
	case output.NodeShapeParallelogram:
		return ShapeParallelogram
	default:
		return ShapeRectangle
	}
}

func graphStyleToD2(s output.NodeStyle) NodeStyle {
	return NodeStyle{
		Fill: s.Fill,
		StrokeStyle: StrokeStyle{
			Stroke:    s.Stroke,
			FontSize:  s.FontSize,
			FontColor: s.FontColor,
		},
	}
}

// NewD2FromTable converts Table to a D2 diagram with per-row nodes connected by edges.
func NewD2FromTable(data *output.Table) *Diagram {
	diagram := NewDiagram()
	if data == nil {
		return diagram
	}

	nodes := output.NodesFromTable(data, output.DefaultGraphNodeLabel)

	for _, n := range nodes {
		diagram.AddNode(graphNodeToD2(n))
	}

	for _, edge := range data.CreateRowEdges() {
		diagram.AddEdgeSimple(edge.From.Get(), edge.To.Get())
	}

	return diagram
}

// NewD2FromTree converts a TreeNode hierarchy to a D2 diagram.
func NewD2FromTree(root *output.TreeNode) *Diagram {
	return output.TreeToRenderer(NewDiagram, (*Diagram).addTreeNodes, root)
}

func (d *Diagram) addTreeNodes(node *output.TreeNode, parentID string) {
	nodeID := node.ID.Get()
	if nodeID == "" {
		nodeID = escape.SlugifyID(node.Label.Get())
	}

	d.AddNode(Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand](nodeID),
		Label: output.NewBrandedID[output.D2NodeLabelBrand](node.Label.Get()),
	})

	if parentID != "" {
		d.AddEdgeSimple(parentID, nodeID)
	}

	for _, child := range node.Children {
		d.addTreeNodes(child, nodeID)
	}
}
