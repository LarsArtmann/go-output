package d2

import "github.com/larsartmann/go-output"

// D2NodeID is a branded identifier for D2 diagram nodes.
type D2NodeID = output.D2NodeID

// D2NodeLabel is a branded identifier for D2 diagram node labels.
type D2NodeLabel = output.D2NodeLabel

// D2NodeStyle represents styling for a D2 node.
type D2NodeStyle struct {
	Fill          string
	Stroke        string
	StrokeWidth   int
	StrokeDash    int
	FontSize      int
	FontColor     string
	Opacity       float64
	Shadow        bool
	BorderRadius  int
	TextTransform string
}

func (s D2NodeStyle) isSet() bool {
	return s.Fill != "" || s.Stroke != "" || s.StrokeWidth > 0 ||
		s.StrokeDash > 0 || s.FontSize > 0 || s.FontColor != "" ||
		s.Opacity > 0 || s.Shadow || s.BorderRadius > 0 ||
		s.TextTransform != ""
}

// D2Node represents a node in a D2 diagram.
type D2Node struct {
	ID          D2NodeID
	Label       D2NodeLabel
	Shape       D2NodeShape
	Style       D2NodeStyle
	Icon        string
	Link        string
	Tooltip     string
	Class       string
	Near        string
	Width       int
	Height      int
	GridRows    int
	GridColumns int
	GridGap     int
	Nested      string
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
	Stroke      string
	StrokeWidth int
	StrokeDash  int
	Animated    bool
	FontColor   string
	FontSize    int
}

// D2Edge represents an edge in a D2 diagram.
type D2Edge struct {
	From        D2NodeID
	To          D2NodeID
	Label       D2NodeLabel
	Style       D2EdgeStyle
	SourceArrow D2ArrowType
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
	Name       string
	Type       string
	Constraint D2Constraint
}

// D2Table represents a SQL table shape in D2 diagrams.
type D2Table struct {
	Name    string
	Columns []D2Column
}
