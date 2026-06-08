package d2

import (
	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

// SetNodes sets graph nodes from the generic GraphNode type, satisfying GraphRenderer.
func (d *D2Diagram) SetNodes(nodes []output.GraphNode) {
	d.nodes = make([]D2Node, len(nodes))
	for i, n := range nodes {
		d.nodes[i] = graphNodeToD2(n)
	}
}

// SetEdges sets graph edges from the generic GraphEdge type, satisfying GraphRenderer.
func (d *D2Diagram) SetEdges(edges []output.GraphEdge) {
	d.edges = make([]D2Edge, len(edges))
	for i, e := range edges {
		d.edges[i] = graphEdgeToD2(e)
	}
}

func graphNodeToD2(n output.GraphNode) D2Node {
	return D2Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand](n.ID.Get()),
		Label: output.NewBrandedID[output.D2NodeLabelBrand](n.Label.Get()),
		Shape: graphShapeToD2(n.Shape),
		Style: graphStyleToD2(n.Style),
	}
}

func graphEdgeToD2(e output.GraphEdge) D2Edge {
	return D2Edge{
		From:  output.NewBrandedID[output.D2NodeIDBrand](e.From.Get()),
		To:    output.NewBrandedID[output.D2NodeIDBrand](e.To.Get()),
		Label: output.NewBrandedID[output.D2NodeLabelBrand](e.Label.Get()),
	}
}

func graphShapeToD2(s output.GraphShape) D2NodeShape {
	switch s {
	case output.ShapeBox, output.ShapeRect:
		return D2ShapeRectangle
	case output.ShapeEllipse:
		return D2ShapeOval
	case output.ShapeDiamond:
		return D2ShapeDiamond
	case output.ShapeCircle:
		return D2ShapeCircle
	case output.ShapeCylinder:
		return D2ShapeCylinder
	case output.ShapeHexagon:
		return D2ShapeHexagon
	case output.ShapeParallelogram:
		return D2ShapeParallelogram
	default:
		return D2ShapeRectangle
	}
}

func graphStyleToD2(s output.GraphStyle) D2NodeStyle {
	return D2NodeStyle{
		Fill: s.FillColor,
		D2StrokeStyle: D2StrokeStyle{
			Stroke:   s.StrokeColor,
			FontSize: s.FontSize,
		},
	}
}

// D2FromTableData converts TableData to a D2 diagram with per-row nodes connected by edges.
func D2FromTableData(data *output.TableData) *D2Diagram {
	diagram := NewD2Diagram()
	if data == nil {
		return diagram
	}

	nodes := output.NodesFromTableData(data, output.DefaultGraphNodeLabel)

	for _, n := range nodes {
		diagram.AddNode(graphNodeToD2(n))
	}

	for _, edge := range data.CreateRowEdges() {
		diagram.AddEdgeSimple(edge.From, edge.To)
	}

	return diagram
}

// D2FromTree converts a TreeNode hierarchy to a D2 diagram.
func D2FromTree(root *output.TreeNode) *D2Diagram {
	diagram := NewD2Diagram()
	if root == nil {
		return diagram
	}

	diagram.addTreeNodes(root, "")

	return diagram
}

func (d *D2Diagram) addTreeNodes(node *output.TreeNode, parentID string) {
	nodeID := node.ID.Get()
	if nodeID == "" {
		nodeID = escape.SlugifyID(node.Label.Get())
	}

	d.AddNode(D2Node{
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
