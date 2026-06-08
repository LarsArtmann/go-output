package output

import (
	"github.com/larsartmann/go-output/enum"
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
	// ID is the unique identifier for the node.
	ID GraphNodeID
	// Label is the display text for the node.
	Label GraphNodeLabel
	// Shape defines the visual shape (box, ellipse, diamond, etc.).
	Shape GraphShape
	// Style contains optional visual styling attributes.
	Style GraphStyle
	// Metadata holds arbitrary key-value pairs for custom data.
	Metadata map[string]string
}

// NewGraphNode creates a new GraphNode.
func NewGraphNode(id, label string) *GraphNode {
	return &GraphNode{
		ID:       NewBrandedID[GraphNodeIDBrand](id),
		Label:    NewBrandedID[GraphNodeLabelBrand](label),
		Shape:    ShapeBox,
		Metadata: make(map[string]string),
	}
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

// InvalidGraphShapeError is returned when an invalid graph shape is provided.
type InvalidGraphShapeError struct {
	Value string
}

// Error returns a descriptive error message for the invalid graph shape.
func (e *InvalidGraphShapeError) Error() string {
	return "invalid graph shape: " + e.Value
}

// ParseGraphShape converts a string to GraphShape, returning an error if invalid.
func ParseGraphShape(s string) (GraphShape, error) {
	v, err := enum.Parse(graphShapeValues, s, func(g GraphShape) string { return string(g) })
	if err != nil {
		return "", &InvalidGraphShapeError{Value: s}
	}

	return v, nil
}

// String returns the string representation of the graph shape.
func (s GraphShape) String() string {
	return string(s)
}

// AllowedValues returns all valid graph shape values.
func (s GraphShape) AllowedValues() []string {
	return enum.AllowedValues(graphShapeValues)
}

// IsValid checks if the graph shape is valid.
func (s GraphShape) IsValid() bool {
	return enum.Contains(graphShapeValues, s)
}

// GraphStyle represents styling attributes for a graph node.
type GraphStyle struct {
	// FillColor is the background color (e.g., "#f9f9f9").
	FillColor string
	// StrokeColor is the border color.
	StrokeColor string
	// FontColor is the text color.
	FontColor string
	// FontSize is the text size in points.
	FontSize int
}

// GraphEdge represents an edge between two nodes.
type GraphEdge struct {
	// From is the source node ID.
	From GraphNodeID
	// To is the target node ID.
	To GraphNodeID
	// Label is the optional display text on the edge.
	Label GraphNodeLabel
	// Style contains optional visual styling attributes.
	Style EdgeStyle
}

// NewGraphEdge creates a new GraphEdge.
func NewGraphEdge(from, to string) *GraphEdge {
	return &GraphEdge{
		From:  NewBrandedID[GraphNodeIDBrand](from),
		To:    NewBrandedID[GraphNodeIDBrand](to),
		Label: NewBrandedID[GraphNodeLabelBrand](""),
	}
}

// EdgeStyle represents styling attributes for an edge.
type EdgeStyle struct {
	// Color is the edge line color.
	Color string
	// Style is the line style ("solid", "dashed", "dotted").
	Style string
	// ArrowHead is the arrowhead style at the target end.
	ArrowHead string
	// ArrowTail is the arrowhead style at the source end.
	ArrowTail string
}

// AddNode appends a node to the graph.
func (m *GraphRendererMixin) AddNode(node GraphNode) {
	m.nodes = append(m.nodes, node)
}

// AddEdge appends an edge to the graph.
func (m *GraphRendererMixin) AddEdge(edge GraphEdge) {
	m.edges = append(m.edges, edge)
}

// NodeEdgeAppender is implemented by types that can add nodes and edges.
type NodeEdgeAppender interface {
	AddNode(node GraphNode)
	AddEdge(edge GraphEdge)
}

// AddTreeNodes recursively adds tree nodes and edges to the provided appender.
func AddTreeNodes(
	a NodeEdgeAppender,
	node *TreeNode, parentID string,
	idFunc TreeNodeIDFunc, shape GraphShape,
) {
	nodeID := idFunc(node)
	graphNodeID := NewBrandedID[GraphNodeIDBrand](nodeID)
	graphNodeLabel := NewBrandedID[GraphNodeLabelBrand](node.Label.Get())

	//nolint:exhaustruct // Uses defaults for optional fields
	a.AddNode(GraphNode{
		ID:    graphNodeID,
		Label: graphNodeLabel,
		Shape: shape,
	})

	if parentID != "" {
		//nolint:exhaustruct // Uses defaults for optional fields
		a.AddEdge(GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand](parentID),
			To:   graphNodeID,
		})
	}

	for _, child := range node.Children {
		AddTreeNodes(a, child, nodeID, idFunc, shape)
	}
}

// GraphRendererMixin contains shared fields and methods for graph renderers.
//
// D2 does not use this mixin because it has richer domain-specific types
// (D2Node, D2Edge with classes, SQL tables, shapes, arrow types, etc.)
// that do not map to the simpler GraphNode/GraphEdge model.
type GraphRendererMixin struct {
	nodes []GraphNode
	edges []GraphEdge
}

// NewGraphRendererMixin creates a new GraphRendererMixin with initialized slices.
func NewGraphRendererMixin() GraphRendererMixin {
	return GraphRendererMixin{
		nodes: make([]GraphNode, 0),
		edges: make([]GraphEdge, 0),
	}
}

// SetNodes sets the graph nodes.
func (m *GraphRendererMixin) SetNodes(nodes []GraphNode) {
	m.nodes = nodes
}

// SetEdges sets the graph edges.
func (m *GraphRendererMixin) SetEdges(edges []GraphEdge) {
	m.edges = edges
}

// Nodes returns the graph nodes.
func (m *GraphRendererMixin) Nodes() []GraphNode {
	return m.nodes
}

// Edges returns the graph edges.
func (m *GraphRendererMixin) Edges() []GraphEdge {
	return m.edges
}
