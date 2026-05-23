package output

import (
	"errors"
	"fmt"
	"strings"

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

// ErrInvalidGraphShape is returned when an invalid graph shape is provided.
var ErrInvalidGraphShape = errors.New("invalid graph shape")

// ParseGraphShape converts a string to GraphShape, returning an error if invalid.
func ParseGraphShape(s string) (GraphShape, error) {
	v, err := enum.Parse(graphShapeValues, s, func(g GraphShape) string { return string(g) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidGraphShape, s)
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
	}
}

// EdgeStyle represents styling attributes for an edge.
type EdgeStyle struct {
	Color     string
	Style     string // solid, dashed, dotted
	ArrowHead string
	ArrowTail string
}

// GraphNodeLabelFunc is a function that formats a cell value with its header into a label.
type GraphNodeLabelFunc func(header, cell string) string

// DefaultGraphNodeLabel returns a label in the format "header: cell".
func DefaultGraphNodeLabel(header, cell string) string {
	return fmt.Sprintf("%s: %s", header, cell)
}

// TreeNodeIDFunc resolves a TreeNode's ID for a specific graph format.
type TreeNodeIDFunc func(*TreeNode) string

// AddTreeNodes recursively adds tree nodes and edges to the provided graph slices.
func AddTreeNodes(
	nodes *[]GraphNode, edges *[]GraphEdge,
	node *TreeNode, parentID string,
	idFunc TreeNodeIDFunc, shape GraphShape,
) {
	nodeID := idFunc(node)
	graphNodeID := NewBrandedID[GraphNodeIDBrand](nodeID)
	graphNodeLabel := NewBrandedID[GraphNodeLabelBrand](node.Label.Get())

	//nolint:exhaustruct // Uses defaults for optional fields
	*nodes = append(*nodes, GraphNode{
		ID:    graphNodeID,
		Label: graphNodeLabel,
		Shape: shape,
	})

	if parentID != "" {
		//nolint:exhaustruct // Uses defaults for optional fields
		*edges = append(*edges, GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand](parentID),
			To:   graphNodeID,
		})
	}

	for _, child := range node.Children {
		AddTreeNodes(nodes, edges, child, nodeID, idFunc, shape)
	}
}

// NodesFromTableData creates GraphNodes from TableData using the provided label function.
func NodesFromTableData(data *TableData, labelFn GraphNodeLabelFunc) []GraphNode {
	if data == nil {
		return nil
	}

	nodes := make([]GraphNode, 0, len(data.Rows))
	for i, row := range data.Rows {
		var labelParts []string

		for j, cell := range row {
			if j < len(data.Headers) {
				labelParts = append(labelParts, labelFn(data.Headers[j], cell))
			} else {
				labelParts = append(labelParts, cell)
			}
		}

		label := strings.Join(labelParts, "\n")
		//nolint:exhaustruct // Uses defaults for optional fields
		nodes = append(nodes, GraphNode{
			ID:    NewBrandedID[GraphNodeIDBrand](fmt.Sprintf("row%d", i)),
			Label: NewBrandedID[GraphNodeLabelBrand](label),
		})
	}

	return nodes
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

// NodesPtr returns a pointer to the graph nodes slice for mutation.
func (m *GraphRendererMixin) NodesPtr() *[]GraphNode {
	return &m.nodes
}

// EdgesPtr returns a pointer to the graph edges slice for mutation.
func (m *GraphRendererMixin) EdgesPtr() *[]GraphEdge {
	return &m.edges
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

// SetNodesFromTableData creates nodes from TableData, applies per-node modifications,
// adds them to the graph, and adds row edges.
func (m *GraphRendererMixin) SetNodesFromTableData(
	data *TableData,
	modifyNode func(i int, n *GraphNode),
) {
	nodes := NodesFromTableData(data, DefaultGraphNodeLabel)
	for i := range nodes {
		modifyNode(i, &nodes[i])
	}

	m.nodes = append(m.nodes, nodes...)
	m.AddRowEdges(data)
}
