package output

import (
	"errors"
	"fmt"
	"slices"
)

// D2Direction constants for diagram layout direction.
type D2Direction string

const (
	D2DirDown  D2Direction = ""
	D2DirRight D2Direction = "right"
	D2DirLeft  D2Direction = "left"
	D2DirUp    D2Direction = "up"
)

var d2DirectionValues = []D2Direction{
	D2DirDown,
	D2DirRight,
	D2DirLeft,
	D2DirUp,
}

var ErrInvalidD2Direction = errors.New("invalid D2 direction")

func ParseD2Direction(s string) (D2Direction, error) {
	if slices.Contains(d2DirectionValues, D2Direction(s)) {
		return D2Direction(s), nil
	}

	return "", fmt.Errorf("%w: %q (allowed: %v)", ErrInvalidD2Direction, s, d2DirectionValues)
}

func (d D2Direction) IsValid() bool {
	return slices.Contains(d2DirectionValues, d)
}

func (d D2Direction) AllowedValues() []string {
	values := make([]string, len(d2DirectionValues))
	for i, v := range d2DirectionValues {
		values[i] = string(v)
	}

	return values
}

func (d D2Direction) String() string {
	return string(d)
}

// D2NodeShape represents the shape of a D2 node.
type D2NodeShape string

// D2NodeShape constants define the available shapes for D2 nodes.
const (
	D2ShapeRectangle     D2NodeShape = "rectangle"
	D2ShapeSquare        D2NodeShape = "square"
	D2ShapeCircle        D2NodeShape = "circle"
	D2ShapeDiamond       D2NodeShape = "diamond"
	D2ShapeHexagon       D2NodeShape = "hexagon"
	D2ShapeCloud         D2NodeShape = "cloud"
	D2ShapeCylinder      D2NodeShape = "cylinder"
	D2ShapePerson        D2NodeShape = "person"
	D2ShapeQueue         D2NodeShape = "queue"
	D2ShapeOval          D2NodeShape = "oval"
	D2ShapeParallelogram D2NodeShape = "parallelogram"
	D2ShapeTriangle      D2NodeShape = "triangle"
	D2ShapeSQLTable      D2NodeShape = "sql_table"
	D2ShapeImage         D2NodeShape = "image"
	D2ShapeCode          D2NodeShape = "code"
	D2ShapeText          D2NodeShape = "text"
	D2ShapeClass         D2NodeShape = "class"
	D2ShapePage          D2NodeShape = "page"
	D2ShapeStep          D2NodeShape = "step"
	D2ShapeStoredData    D2NodeShape = "stored_data"
)

var d2NodeShapeValues = []D2NodeShape{
	D2ShapeRectangle,
	D2ShapeSquare,
	D2ShapeCircle,
	D2ShapeDiamond,
	D2ShapeHexagon,
	D2ShapeCloud,
	D2ShapeCylinder,
	D2ShapePerson,
	D2ShapeQueue,
	D2ShapeOval,
	D2ShapeParallelogram,
	D2ShapeTriangle,
	D2ShapeSQLTable,
	D2ShapeImage,
	D2ShapeCode,
	D2ShapeText,
	D2ShapeClass,
	D2ShapePage,
	D2ShapeStep,
	D2ShapeStoredData,
}

var ErrInvalidD2NodeShape = errors.New("invalid D2 node shape")

func ParseD2NodeShape(s string) (D2NodeShape, error) {
	if slices.Contains(d2NodeShapeValues, D2NodeShape(s)) {
		return D2NodeShape(s), nil
	}

	return "", fmt.Errorf("%w: %q (allowed: %v)", ErrInvalidD2NodeShape, s, d2NodeShapeValues)
}

func (s D2NodeShape) IsValid() bool {
	return slices.Contains(d2NodeShapeValues, s)
}

func (s D2NodeShape) AllowedValues() []string {
	values := make([]string, len(d2NodeShapeValues))
	for i, v := range d2NodeShapeValues {
		values[i] = string(v)
	}

	return values
}

func (s D2NodeShape) String() string {
	return string(s)
}

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

// D2ArrowType represents the type of arrow for D2 edges.
type D2ArrowType string

// D2ArrowType constants define the available arrow shapes for D2 edges.
const (
	D2ArrowNone           D2ArrowType = ""
	D2ArrowArrow          D2ArrowType = "arrow"
	D2ArrowTriangle       D2ArrowType = "triangle"
	D2ArrowDiamond        D2ArrowType = "diamond"
	D2ArrowCircle         D2ArrowType = "circle"
	D2ArrowFilled         D2ArrowType = "filled"
	D2ArrowBox            D2ArrowType = "box"
	D2ArrowCross          D2ArrowType = "cross"
	D2ArrowCFOne          D2ArrowType = "cf-one"
	D2ArrowCFMany         D2ArrowType = "cf-many"
	D2ArrowCFOneRequired  D2ArrowType = "cf-one-required"
	D2ArrowCFManyRequired D2ArrowType = "cf-many-required"
)

// Deprecated: Use D2ArrowArrow instead.
const D2ArrowPoint = D2ArrowArrow

// Deprecated: Use D2ArrowCircle instead.
const D2ArrowOval = D2ArrowCircle

var d2ArrowTypeValues = []D2ArrowType{
	D2ArrowArrow,
	D2ArrowTriangle,
	D2ArrowDiamond,
	D2ArrowCircle,
	D2ArrowFilled,
	D2ArrowBox,
	D2ArrowCross,
	D2ArrowCFOne,
	D2ArrowCFMany,
	D2ArrowCFOneRequired,
	D2ArrowCFManyRequired,
}

var ErrInvalidD2ArrowType = errors.New("invalid D2 arrow type")

func ParseD2ArrowType(s string) (D2ArrowType, error) {
	if slices.Contains(d2ArrowTypeValues, D2ArrowType(s)) {
		return D2ArrowType(s), nil
	}

	return "", fmt.Errorf("%w: %q (allowed: %v)", ErrInvalidD2ArrowType, s, d2ArrowTypeValues)
}

func (a D2ArrowType) IsValid() bool {
	return a == D2ArrowNone || slices.Contains(d2ArrowTypeValues, a)
}

func (a D2ArrowType) AllowedValues() []string {
	values := make([]string, len(d2ArrowTypeValues))
	for i, v := range d2ArrowTypeValues {
		values[i] = string(v)
	}

	return values
}

func (a D2ArrowType) String() string {
	return string(a)
}

// D2Constraint represents a SQL constraint on a table column.
type D2Constraint string

// D2Constraint constants define the available SQL constraints for D2 table columns.
const (
	D2ConstraintPrimary D2Constraint = "primary_key"
	D2ConstraintForeign D2Constraint = "foreign_key"
	D2ConstraintUnique  D2Constraint = "unique"
)

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
