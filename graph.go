package output

import (
	"fmt"
	"slices"
)

// GraphRenderer defines the interface for graph format renderers.
type GraphRenderer interface {
	Renderer
	// SetNodes sets the graph nodes.
	SetNodes(nodes []GraphNode)
	// SetEdges sets the graph edges.
	SetEdges(edges []GraphEdge)
}

// GraphNode represents a node in a graph.
type GraphNode struct {
	ID       GraphNodeID
	Label    GraphNodeLabel
	Shape    GraphShape
	Style    GraphStyle
	Metadata map[string]string
}

// NewGraphNode creates a new GraphNode.
func NewGraphNode(id, label string) *GraphNode {
	return &GraphNode{
		ID:       NewBrandedID[GraphNodeIDBrand](id),
		Label:    NewBrandedID[GraphNodeLabelBrand](label),
		Shape:    ShapeBox,
		Style:    GraphStyle{FillColor: "", StrokeColor: "", FontColor: "", FontSize: 0},
		Metadata: make(map[string]string),
	}
}

// GetStyle returns the node's style.
func (n *GraphNode) GetStyle() GraphStyle {
	return n.Style
}

// GraphShape represents the shape of a graph node.
type GraphShape string

// GraphShape constants define the available shapes for graph nodes.
const (
	ShapeBox           GraphShape = "box"
	ShapeEllipse       GraphShape = "ellipse"
	ShapeDiamond       GraphShape = "diamond"
	ShapeCircle        GraphShape = "circle"
	ShapeCylinder      GraphShape = "cylinder"
	ShapeHexagon       GraphShape = "hexagon"
	ShapeParallelogram GraphShape = "parallelogram"
	ShapeRect          GraphShape = "rect"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var graphShapeValues = []GraphShape{
	ShapeBox,
	ShapeEllipse,
	ShapeDiamond,
	ShapeCircle,
	ShapeCylinder,
	ShapeHexagon,
	ShapeParallelogram,
	ShapeRect,
}

// ParseGraphShape parses a graph shape string.
func ParseGraphShape(s string) (GraphShape, error) {
	if slices.Contains(graphShapeValues, GraphShape(s)) {
		return GraphShape(s), nil
	}

	return "", fmt.Errorf("invalid graph shape: %q (allowed: %v)", s, graphShapeValues)
}

func (s GraphShape) String() string {
	return string(s)
}

// AllowedValues returns all valid graph shape values.
func (s GraphShape) AllowedValues() []string {
	values := make([]string, len(graphShapeValues))
	for i, v := range graphShapeValues {
		values[i] = string(v)
	}

	return values
}

// IsValid checks if the graph shape is valid.
func (s GraphShape) IsValid() bool {
	return slices.Contains(graphShapeValues, s)
}

// GraphStyle represents styling attributes for a graph node.
type GraphStyle struct {
	FillColor   string
	StrokeColor string
	FontColor   string
	FontSize    int
}

// GraphEdge represents an edge between two nodes.
type GraphEdge struct {
	From  GraphNodeID
	To    GraphNodeID
	Label GraphNodeLabel
	Style EdgeStyle
}

// NewGraphEdge creates a new GraphEdge.
func NewGraphEdge(from, to string) *GraphEdge {
	return &GraphEdge{
		From:  NewBrandedID[GraphNodeIDBrand](from),
		To:    NewBrandedID[GraphNodeIDBrand](to),
		Label: NewBrandedID[GraphNodeLabelBrand](""),
		Style: EdgeStyle{Color: "", Style: "", ArrowHead: "", ArrowTail: ""},
	}
}

// EdgeStyle represents styling attributes for an edge.
type EdgeStyle struct {
	Color     string
	Style     string // solid, dashed, dotted
	ArrowHead string
	ArrowTail string
}
