package d2

import "github.com/larsartmann/go-output"

//nolint:gochecknoinits // Registers D2 format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatD2, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
}

// D2NodeID is a branded identifier for D2 diagram nodes.
type D2NodeID = output.D2NodeID

// D2NodeLabel is a branded identifier for D2 diagram node labels.
type D2NodeLabel = output.D2NodeLabel

// StrokeStyle represents shared stroke/font styling for D2 nodes and edges.
type StrokeStyle struct {
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

func (s StrokeStyle) isSet() bool {
	return s.Stroke != "" || s.StrokeWidth > 0 || s.StrokeDash > 0 ||
		s.FontSize > 0 || s.FontColor != ""
}

// NodeStyle represents styling for a D2 node.
type NodeStyle struct {
	StrokeStyle

	// Fill is the background color (e.g., "blue", "#0000ff").
	Fill string
	// Opacity is the transparency level (0.0 to 1.0).
	Opacity float64
	// Shadow enables a drop shadow effect.
	Shadow bool
	// BorderRadius is the corner radius in pixels.
	BorderRadius int
	// TextTransform controls text casing (uppercase, lowercase, capitalize).
	TextTransform TextTransform
}

func (s NodeStyle) isSet() bool {
	return s.StrokeStyle.isSet() || s.Fill != "" ||
		s.Opacity > 0 || s.Shadow || s.BorderRadius > 0 || s.TextTransform != ""
}

// Node represents a node in a D2 diagram.
type Node struct {
	// ID is the unique identifier for the node (used in connections).
	ID D2NodeID
	// Label is the display text. Defaults to ID if empty.
	Label D2NodeLabel
	// Shape defines the visual shape (rectangle, circle, etc.).
	Shape NodeShape
	// Style contains optional visual styling.
	Style NodeStyle
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

func (n Node) hasBlockAttrs() bool {
	return n.hasVisualAttrs() || n.hasLayoutAttrs()
}

func (n Node) hasVisualAttrs() bool {
	hasShape := n.Shape != "" && n.Shape != ShapeRectangle

	return hasShape || n.Style.isSet() || n.Icon != "" || n.Link != "" || n.Tooltip != ""
}

func (n Node) hasLayoutAttrs() bool {
	return n.Class != "" || n.Near != "" || n.hasGrid() || n.hasSize()
}

func (n Node) hasGrid() bool {
	return n.GridRows > 0 || n.GridColumns > 0 || n.GridGap > 0
}

func (n Node) hasSize() bool {
	return n.Width > 0 || n.Height > 0
}

// EdgeStyle represents styling for a D2 edge.
type EdgeStyle struct {
	StrokeStyle

	// Animated enables edge animation.
	Animated bool
}

// Edge represents an edge in a D2 diagram.
type Edge struct {
	// From is the source node ID.
	From D2NodeID
	// To is the target node ID.
	To D2NodeID
	// Label is the optional display text on the edge.
	Label D2NodeLabel
	// Style contains optional visual styling.
	Style EdgeStyle
	// SourceArrow is the arrowhead style at the source end.
	SourceArrow ArrowType
	// TargetArrow is the arrowhead style at the target end.
	TargetArrow ArrowType
}

func (e Edge) hasBlockAttrs() bool {
	hasStyle := e.Style.isSet() || e.Style.Animated
	hasArrows := e.SourceArrow != "" || e.TargetArrow != ""

	return hasStyle || hasArrows
}

// Column represents a column in a D2 SQL table shape.
type Column struct {
	// Name is the column name.
	Name string
	// Type is the column data type (e.g., "varchar", "int").
	Type string
	// Constraint is the column constraint (e.g., primary key, not null).
	Constraint Constraint
}

// Table represents a SQL table shape in D2 diagrams.
type Table struct {
	// Name is the table name.
	Name string
	// Columns lists the table columns.
	Columns []Column
}
