package d2

import (
	"errors"
	"fmt"

	output "github.com/larsartmann/go-output"
)

// Direction constants for diagram layout direction.
// Use output.Direction for cross-format portability; Direction stays
// for D2-specific APIs. Convert via Direction.ToDirection().
type Direction string

// D2 direction constants.
const (
	DirDown  Direction = ""
	DirRight Direction = "right"
	DirLeft  Direction = "left"
	DirUp    Direction = "up"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var directionValues = []Direction{
	DirDown,
	DirRight,
	DirLeft,
	DirUp,
}

// ErrInvalidDirection is returned when an invalid D2 direction is provided.
var ErrInvalidDirection = errors.New("invalid D2 direction")

// ParseDirection converts a string to Direction, returning an error if invalid.
func ParseDirection(s string) (Direction, error) {
	v, err := output.ParseEnum(directionValues, s, func(d Direction) string { return string(d) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidDirection, s)
	}

	return v, nil
}

// IsValid returns true if the direction is a valid Direction value.
func (d Direction) IsValid() bool {
	return output.ContainsEnum(directionValues, d)
}

// AllowedValues returns all valid D2 direction values for CLI help text.
func (d Direction) AllowedValues() []string {
	return output.EnumAllowedValues(directionValues)
}

// String returns the string representation of the direction.
func (d Direction) String() string {
	return string(d)
}

// NodeShape represents the shape of a D2 node.
type NodeShape string

// NodeShape constants define the available shapes for D2 nodes.
const (
	ShapeRectangle     NodeShape = "rectangle"
	ShapeSquare        NodeShape = "square"
	ShapeCircle        NodeShape = "circle"
	ShapeDiamond       NodeShape = "diamond"
	ShapeHexagon       NodeShape = "hexagon"
	ShapeCloud         NodeShape = "cloud"
	ShapeCylinder      NodeShape = "cylinder"
	ShapePerson        NodeShape = "person"
	ShapeQueue         NodeShape = "queue"
	ShapeOval          NodeShape = "oval"
	ShapeParallelogram NodeShape = "parallelogram"
	ShapeTriangle      NodeShape = "triangle"
	ShapeSQLTable      NodeShape = "sql_table"
	ShapeImage         NodeShape = "image"
	ShapeCode          NodeShape = "code"
	ShapeText          NodeShape = "text"
	ShapeClass         NodeShape = "class"
	ShapePage          NodeShape = "page"
	ShapeStep          NodeShape = "step"
	ShapeStoredData    NodeShape = "stored_data"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var nodeShapeValues = []NodeShape{
	ShapeRectangle,
	ShapeSquare,
	ShapeCircle,
	ShapeDiamond,
	ShapeHexagon,
	ShapeCloud,
	ShapeCylinder,
	ShapePerson,
	ShapeQueue,
	ShapeOval,
	ShapeParallelogram,
	ShapeTriangle,
	ShapeSQLTable,
	ShapeImage,
	ShapeCode,
	ShapeText,
	ShapeClass,
	ShapePage,
	ShapeStep,
	ShapeStoredData,
}

// ErrInvalidNodeShape is returned when an invalid D2 node shape is provided.
var ErrInvalidNodeShape = errors.New("invalid D2 node shape")

// ParseNodeShape converts a string to NodeShape, returning an error if invalid.
func ParseNodeShape(s string) (NodeShape, error) {
	v, err := output.ParseEnum(nodeShapeValues, s, func(ns NodeShape) string { return string(ns) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidNodeShape, s)
	}

	return v, nil
}

// IsValid returns true if the shape is a valid NodeShape value.
func (s NodeShape) IsValid() bool {
	return output.ContainsEnum(nodeShapeValues, s)
}

// AllowedValues returns all valid D2 node shape values for CLI help text.
func (s NodeShape) AllowedValues() []string {
	return output.EnumAllowedValues(nodeShapeValues)
}

// String returns the string representation of the node shape.
func (s NodeShape) String() string {
	return string(s)
}

// ArrowType represents the type of arrow for D2 edges.
type ArrowType string

// ArrowType constants define the available arrow shapes for D2 edges.
const (
	ArrowNone           ArrowType = ""
	ArrowArrow          ArrowType = "arrow"
	ArrowTriangle       ArrowType = "triangle"
	ArrowDiamond        ArrowType = "diamond"
	ArrowCircle         ArrowType = "circle"
	ArrowFilled         ArrowType = "filled"
	ArrowBox            ArrowType = "box"
	ArrowCross          ArrowType = "cross"
	ArrowCFOne          ArrowType = "cf-one"
	ArrowCFMany         ArrowType = "cf-many"
	ArrowCFOneRequired  ArrowType = "cf-one-required"
	ArrowCFManyRequired ArrowType = "cf-many-required"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var arrowTypeValues = []ArrowType{
	ArrowNone,
	ArrowArrow,
	ArrowTriangle,
	ArrowDiamond,
	ArrowCircle,
	ArrowFilled,
	ArrowBox,
	ArrowCross,
	ArrowCFOne,
	ArrowCFMany,
	ArrowCFOneRequired,
	ArrowCFManyRequired,
}

// ErrInvalidArrowType is returned when an invalid D2 arrow type is provided.
var ErrInvalidArrowType = errors.New("invalid D2 arrow type")

// ParseArrowType converts a string to ArrowType, returning an error if invalid.
func ParseArrowType(s string) (ArrowType, error) {
	v, err := output.ParseEnum(arrowTypeValues, s, func(a ArrowType) string { return string(a) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidArrowType, s)
	}

	return v, nil
}

// IsValid returns true if the arrow type is a valid ArrowType value.
func (a ArrowType) IsValid() bool {
	return a == ArrowNone || output.ContainsEnum(arrowTypeValues, a)
}

// AllowedValues returns all valid D2 arrow type values for CLI help text.
func (a ArrowType) AllowedValues() []string {
	return output.EnumAllowedValues(arrowTypeValues)
}

// String returns the string representation of the arrow type.
func (a ArrowType) String() string {
	return string(a)
}

// Constraint represents a SQL constraint on a table column.
type Constraint string

// Constraint constants define the available SQL constraints for D2 table columns.
const (
	ConstraintPrimary Constraint = "primary_key"
	ConstraintForeign Constraint = "foreign_key"
	ConstraintUnique  Constraint = "unique"
)

// ErrInvalidConstraint is returned when an invalid D2 constraint is provided.
var ErrInvalidConstraint = errors.New("invalid D2 constraint")

//nolint:gochecknoglobals // Allowed values for Constraint validation.
var allConstraints = []Constraint{
	ConstraintPrimary,
	ConstraintForeign,
	ConstraintUnique,
}

// AllConstraints returns all valid Constraint values.
func AllConstraints() []Constraint {
	return allConstraints
}

// ParseConstraint converts a string to Constraint, returning an error if invalid.
func ParseConstraint(s string) (Constraint, error) {
	v, err := output.ParseEnum(allConstraints, s, func(d Constraint) string { return string(d) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidConstraint, s)
	}

	return v, nil
}

// String returns the string representation of the constraint.
func (c Constraint) String() string { return string(c) }

// AllowedValues returns all valid Constraint values for CLI help text.
func (c Constraint) AllowedValues() []string {
	return output.EnumAllowedValues(allConstraints)
}

// IsValid returns true if the constraint is a valid Constraint value.
func (c Constraint) IsValid() bool {
	return output.ContainsEnum(allConstraints, c)
}

// TextTransform represents the text-casing transform for a D2 node's label.
type TextTransform string

// TextTransform constants.
const (
	TextTransformNone       TextTransform = ""
	TextTransformUpper      TextTransform = "uppercase"
	TextTransformLower      TextTransform = "lowercase"
	TextTransformCapitalize TextTransform = "capitalize"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var textTransformValues = []TextTransform{
	TextTransformUpper,
	TextTransformLower,
	TextTransformCapitalize,
}

// ErrInvalidTextTransform is returned when an invalid text transform is provided.
var ErrInvalidTextTransform = errors.New("invalid D2 text transform")

// ParseTextTransform converts a string to TextTransform, returning an error if invalid.
func ParseTextTransform(s string) (TextTransform, error) {
	if s == "" {
		return TextTransformNone, nil
	}

	v, err := output.ParseEnum(textTransformValues, s, func(t TextTransform) string { return string(t) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidTextTransform, s)
	}

	return v, nil
}

// IsValid returns true if the text transform is a valid TextTransform value.
func (t TextTransform) IsValid() bool {
	return t == TextTransformNone || output.ContainsEnum(textTransformValues, t)
}

// AllowedValues returns all valid D2 text transform values for CLI help text.
func (t TextTransform) AllowedValues() []string {
	return output.EnumAllowedValues(textTransformValues)
}

// String returns the string representation of the text transform.
func (t TextTransform) String() string {
	return string(t)
}
