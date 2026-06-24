package output

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
	Shape NodeShape
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
		Shape:    NodeShapeBox,
		Metadata: make(map[string]string),
	}
}

// NodeShape represents the shape of a graph node.
type NodeShape string

// NodeShape constants define the available shapes for graph nodes.
const (
	NodeShapeBox           NodeShape = "box"
	NodeShapeEllipse       NodeShape = "ellipse"
	NodeShapeDiamond       NodeShape = "diamond"
	NodeShapeCircle        NodeShape = "circle"
	NodeShapeCylinder      NodeShape = "cylinder"
	NodeShapeHexagon       NodeShape = "hexagon"
	NodeShapeParallelogram NodeShape = "parallelogram"

	// Deprecated: Use NodeShapeBox instead. NodeShapeRect will be removed in v2.
	NodeShapeRect NodeShape = "rect"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var nodeShapeValues = []NodeShape{
	NodeShapeBox,
	NodeShapeEllipse,
	NodeShapeDiamond,
	NodeShapeCircle,
	NodeShapeCylinder,
	NodeShapeHexagon,
	NodeShapeParallelogram,
	NodeShapeRect, //nolint:staticcheck // deprecated but must stay in allowed-values for backward compat
}

// LineStyle represents the visual style of a line (edge).
type LineStyle string

const (
	LineStyleSolid  LineStyle = "solid"
	LineStyleDashed LineStyle = "dashed"
	LineStyleDotted LineStyle = "dotted"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var AllLineStyles = []LineStyle{
	LineStyleSolid,
	LineStyleDashed,
	LineStyleDotted,
}

// ParseLineStyle converts a string to LineStyle, returning an error if invalid.
func ParseLineStyle(s string) (LineStyle, error) {
	v, err := ParseEnum(AllLineStyles, s, func(l LineStyle) string { return string(l) })
	if err != nil {
		return "", &InvalidLineStyleError{Value: s, Allowed: AllLineStyles}
	}

	return v, nil
}

// AllowedValues returns all valid line style values for CLI help text.
func (l LineStyle) AllowedValues() []string {
	return EnumAllowedValues(AllLineStyles)
}

// IsValid returns true if the LineStyle is a recognized value.
func (l LineStyle) IsValid() bool {
	return ContainsEnum(AllLineStyles, l)
}

// String returns the string representation of the LineStyle.
func (l LineStyle) String() string { return string(l) }

// InvalidLineStyleError is returned when an invalid line style is provided.
type InvalidLineStyleError struct {
	Value   string
	Allowed []LineStyle
}

// Error returns a descriptive error message for the invalid line style.
func (e *InvalidLineStyleError) Error() string {
	return "invalid line style: " + e.Value + " (allowed: " + joinStrings(EnumAllowedValues(e.Allowed)) + ")"
}

// InvalidNodeShapeError is returned when an invalid graph shape is provided.
type InvalidNodeShapeError struct {
	Value string
}

// Error returns a descriptive error message for the invalid graph shape.
func (e *InvalidNodeShapeError) Error() string {
	return "invalid graph shape: " + e.Value
}

// ParseNodeShape converts a string to NodeShape, returning an error if invalid.
func ParseNodeShape(s string) (NodeShape, error) {
	v, err := ParseEnum(nodeShapeValues, s, func(g NodeShape) string { return string(g) })
	if err != nil {
		return "", &InvalidNodeShapeError{Value: s}
	}

	return v, nil
}

// String returns the string representation of the graph shape.
func (s NodeShape) String() string {
	return string(s)
}

// AllowedValues returns all valid graph shape values.
func (s NodeShape) AllowedValues() []string {
	return EnumAllowedValues(nodeShapeValues)
}

// IsValid checks if the graph shape is valid.
func (s NodeShape) IsValid() bool {
	return ContainsEnum(nodeShapeValues, s)
}

// GraphStyle represents styling attributes for a graph node.
type GraphStyle struct {
	// Fill is the background color (e.g., "#f9f9f9").
	Fill string
	// Stroke is the border color.
	Stroke string
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
	// Line is the line style (solid, dashed, dotted).
	Line LineStyle
	// ArrowHead is the arrowhead style at the target end.
	ArrowHead string
	// ArrowTail is the arrowhead style at the source end.
	ArrowTail string
}

// AddNode appends a node to the graph.
func (m *GraphRendererState) AddNode(node GraphNode) {
	m.nodes = append(m.nodes, node)
}

// AddEdge appends an edge to the graph.
func (m *GraphRendererState) AddEdge(edge GraphEdge) {
	m.edges = append(m.edges, edge)
}

// DedupEdges removes duplicate edges in-place. Two edges are considered
// duplicates if they share the same From and To node IDs. The first occurrence
// is kept; subsequent duplicates are silently discarded. Edges with different
// labels between the same node pair are also considered duplicates — if you
// need parallel edges, do not call this method.
func (m *GraphRendererState) DedupEdges() {
	if len(m.edges) <= 1 {
		return
	}

	seen := make(map[string]struct{}, len(m.edges))
	deduped := make([]GraphEdge, 0, len(m.edges))

	for _, edge := range m.edges {
		key := edge.From.Get() + "\x00" + edge.To.Get()
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}

		deduped = append(deduped, edge)
	}

	m.edges = deduped
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
	idFunc TreeNodeIDFunc, shape NodeShape,
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

// GraphRendererState contains shared fields and methods for graph renderers.
//
// D2 does not use this mixin because it has richer domain-specific types
// (D2Node, D2Edge with classes, SQL tables, shapes, arrow types, etc.)
// that do not map to the simpler GraphNode/GraphEdge model.
type GraphRendererState struct {
	nodes []GraphNode
	edges []GraphEdge
}

// NewGraphRendererState creates a new GraphRendererState with initialized slices.
func NewGraphRendererState() GraphRendererState {
	return GraphRendererState{
		nodes: make([]GraphNode, 0),
		edges: make([]GraphEdge, 0),
	}
}

// SetNodes sets the graph nodes.
func (m *GraphRendererState) SetNodes(nodes []GraphNode) {
	m.nodes = nodes
}

// SetEdges sets the graph edges.
func (m *GraphRendererState) SetEdges(edges []GraphEdge) {
	m.edges = edges
}

// Nodes returns the graph nodes.
func (m *GraphRendererState) Nodes() []GraphNode {
	return m.nodes
}

// Edges returns the graph edges.
func (m *GraphRendererState) Edges() []GraphEdge {
	return m.edges
}
