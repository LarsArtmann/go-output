package d2

import "github.com/larsartmann/go-output"

// D2NodeID is a branded identifier for D2 diagram nodes.
type D2NodeID = output.D2NodeID

// D2NodeLabel is a branded identifier for D2 diagram node labels.
type D2NodeLabel = output.D2NodeLabel

// D2StrokeStyle represents shared stroke/font styling for D2 nodes and edges.
type D2StrokeStyle struct {
	// Stroke is the border color (e.g., "red", "#ff0000").
	Stroke string
	// StrokeWidth is the border thickness in pixels.
	StrokeWidth int
	// StrokeDash is the dash gap length (0 = solid).
	StrokeDash int
	// FontSize is the label text size in points.
	FontSize int
	// FontColor is the label text color.
	FontColor string
}

func (s D2StrokeStyle) isSet() bool {
	return s.Stroke != "" || s.StrokeWidth > 0 || s.StrokeDash > 0 ||
		s.FontSize > 0 || s.FontColor != ""
}

// D2NodeStyle represents styling for a D2 node.
type D2NodeStyle struct {
	D2StrokeStyle

	// Fill is the background color (e.g., "blue", "#0000ff").
	Fill string
	// Opacity is the transparency level (0.0 to 1.0).
	Opacity float64
	// Shadow enables a drop shadow effect.
	Shadow bool
	// BorderRadius is the corner radius in pixels.
	BorderRadius int
	// TextTransform controls text casing ("uppercase", "lowercase", "capitalize").
	TextTransform string
}

func (s D2NodeStyle) isSet() bool {
	return s.D2StrokeStyle.isSet() || s.Fill != "" ||
		s.Opacity > 0 || s.Shadow || s.BorderRadius > 0 || s.TextTransform != ""
}

// D2Node represents a node in a D2 diagram.
type D2Node struct {
	// ID is the unique identifier for the node (used in connections).
	ID D2NodeID
	// Label is the display text. Defaults to ID if empty.
	Label D2NodeLabel
	// Shape defines the visual shape (rectangle, circle, etc.).
	Shape D2NodeShape
	// Style contains optional visual styling.
	Style D2NodeStyle
	// Icon is an optional icon name or URL.
	Icon string
	// Link is an optional URL the node links to.
	Link string
	// Tooltip is an optional hover tooltip text.
	Tooltip string
	// Class is an optional CSS-like class name for shared styling.
	Class string
	// Near positions the node near another element.
	Near string
	// Width sets an explicit width in pixels.
	Width int
	// Height sets an explicit height in pixels.
	Height int
	// GridRows defines the number of rows in a grid layout.
	GridRows int
	// GridColumns defines the number of columns in a grid layout.
	GridColumns int
	// GridGap sets the spacing between grid cells in pixels.
	GridGap int
	// Nested is the D2 source for nested diagram content.
	Nested string
}

func (n D2Node) hasBlockAttrs() bool {
	return n.hasVisualAttrs() || n.hasLayoutAttrs()
}

func (n D2Node) hasVisualAttrs() bool {
	hasShape := n.Shape != "" && n.Shape != D2ShapeRectangle

	return hasShape || n.Style.isSet() || n.Icon != "" || n.Link != "" || n.Tooltip != ""
}

func (n D2Node) hasLayoutAttrs() bool {
	return n.Class != "" || n.Near != "" || n.hasGrid() || n.hasSize()
}

func (n D2Node) hasGrid() bool {
	return n.GridRows > 0 || n.GridColumns > 0 || n.GridGap > 0
}

func (n D2Node) hasSize() bool {
	return n.Width > 0 || n.Height > 0
}

// D2EdgeStyle represents styling for a D2 edge.
type D2EdgeStyle struct {
	D2StrokeStyle

	// Animated enables edge animation.
	Animated bool
}

// D2Edge represents an edge in a D2 diagram.
type D2Edge struct {
	// From is the source node ID.
	From D2NodeID
	// To is the target node ID.
	To D2NodeID
	// Label is the optional display text on the edge.
	Label D2NodeLabel
	// Style contains optional visual styling.
	Style D2EdgeStyle
	// SourceArrow is the arrowhead style at the source end.
	SourceArrow D2ArrowType
	// TargetArrow is the arrowhead style at the target end.
	TargetArrow D2ArrowType
}

func (e D2Edge) hasBlockAttrs() bool {
	s := e.Style
	hasStyle := s.Stroke != "" || s.StrokeWidth > 0 || s.StrokeDash > 0 ||
		s.Animated || s.FontColor != "" || s.FontSize > 0
	hasArrows := e.SourceArrow != "" || e.TargetArrow != ""

	return hasStyle || hasArrows
}

// D2Column represents a column in a D2 SQL table shape.
type D2Column struct {
	// Name is the column name.
	Name string
	// Type is the column data type (e.g., "varchar", "int").
	Type string
	// Constraint is the column constraint (e.g., primary key, not null).
	Constraint D2Constraint
}

// D2Table represents a SQL table shape in D2 diagrams.
type D2Table struct {
	// Name is the table name.
	Name string
	// Columns lists the table columns.
	Columns []D2Column
}
