package output

import (
	"strings"
)

// SetNodes sets graph nodes from the generic GraphNode type, satisfying GraphRenderer.
func (d *D2Diagram) SetNodes(nodes []GraphNode) {
	d.nodes = make([]D2Node, len(nodes))
	for i, n := range nodes {
		d.nodes[i] = graphNodeToD2(n)
	}
}

// SetEdges sets graph edges from the generic GraphEdge type, satisfying GraphRenderer.
func (d *D2Diagram) SetEdges(edges []GraphEdge) {
	d.edges = make([]D2Edge, len(edges))
	for i, e := range edges {
		d.edges[i] = graphEdgeToD2(e)
	}
}

func graphNodeToD2(n GraphNode) D2Node {
	return D2Node{
		ID:    NewBrandedID[D2NodeIDBrand](n.ID.Get()),
		Label: NewBrandedID[D2NodeLabelBrand](n.Label.Get()),
		Shape: graphShapeToD2(n.Shape),
		Style: graphStyleToD2(n.Style),
	}
}

func graphEdgeToD2(e GraphEdge) D2Edge {
	return D2Edge{
		From:  NewBrandedID[D2NodeIDBrand](e.From.Get()),
		To:    NewBrandedID[D2NodeIDBrand](e.To.Get()),
		Label: NewBrandedID[D2NodeLabelBrand](e.Label.Get()),
	}
}

func graphShapeToD2(s GraphShape) D2NodeShape {
	switch s {
	case ShapeBox, ShapeRect:
		return D2ShapeRectangle
	case ShapeEllipse:
		return D2ShapeOval
	case ShapeDiamond:
		return D2ShapeDiamond
	case ShapeCircle:
		return D2ShapeCircle
	case ShapeCylinder:
		return D2ShapeCylinder
	case ShapeHexagon:
		return D2ShapeHexagon
	case ShapeParallelogram:
		return D2ShapeParallelogram
	default:
		return D2ShapeRectangle
	}
}

func graphStyleToD2(s GraphStyle) D2NodeStyle {
	return D2NodeStyle{
		Fill:     s.FillColor,
		Stroke:   s.StrokeColor,
		FontSize: s.FontSize,
	}
}

// D2FromTableData converts TableData to a D2 diagram with per-row nodes connected by edges.
func D2FromTableData(data *TableData) *D2Diagram {
	diagram := NewD2Diagram()
	if data == nil {
		return diagram
	}

	nodes := NodesFromTableData(data, DefaultGraphNodeLabel)

	for _, n := range nodes {
		diagram.AddNode(graphNodeToD2(n))
	}

	for _, edge := range data.CreateRowEdges() {
		diagram.AddEdgeSimple(edge.From, edge.To)
	}

	return diagram
}

// D2FromTree converts a TreeNode hierarchy to a D2 diagram.
func D2FromTree(root *TreeNode) *D2Diagram {
	diagram := NewD2Diagram()
	if root == nil {
		return diagram
	}

	diagram.addTreeNodes(root, "")

	return diagram
}

func (d *D2Diagram) addTreeNodes(node *TreeNode, parentID string) {
	nodeID := node.ID.Get()
	if nodeID == "" {
		nodeID = strings.ReplaceAll(node.Label.Get(), " ", "_")
	}

	d.AddNode(D2Node{
		ID:    NewBrandedID[D2NodeIDBrand](nodeID),
		Label: NewBrandedID[D2NodeLabelBrand](node.Label.Get()),
	})

	if parentID != "" {
		d.AddEdgeSimple(parentID, nodeID)
	}

	for _, child := range node.Children {
		d.addTreeNodes(child, nodeID)
	}
}
