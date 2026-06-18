package d2

import (
	"errors"
	"fmt"

	"github.com/larsartmann/go-output/enum"
)

// D2Direction constants for diagram layout direction.
// Use output.Direction for cross-format portability; D2Direction stays
// for D2-specific APIs. Convert via Direction.ToD2Direction().
type D2Direction string

// D2 direction constants.
const (
	D2DirDown  D2Direction = ""
	D2DirRight D2Direction = "right"
	D2DirLeft  D2Direction = "left"
	D2DirUp    D2Direction = "up"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var d2DirectionValues = []D2Direction{
	D2DirDown,
	D2DirRight,
	D2DirLeft,
	D2DirUp,
}

// ErrInvalidD2Direction is returned when an invalid D2 direction is provided.
var ErrInvalidD2Direction = errors.New("invalid D2 direction")

// ParseD2Direction converts a string to D2Direction, returning an error if invalid.
func ParseD2Direction(s string) (D2Direction, error) {
	v, err := enum.Parse(d2DirectionValues, s, func(d D2Direction) string { return string(d) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidD2Direction, s)
	}

	return v, nil
}

// IsValid returns true if the direction is a valid D2Direction value.
func (d D2Direction) IsValid() bool {
	return enum.Contains(d2DirectionValues, d)
}

// AllowedValues returns all valid D2 direction values for CLI help text.
func (d D2Direction) AllowedValues() []string {
	return enum.AllowedValues(d2DirectionValues)
}

// String returns the string representation of the direction.
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

//nolint:gochecknoglobals // Global variable used for value iteration.
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

// ErrInvalidD2NodeShape is returned when an invalid D2 node shape is provided.
var ErrInvalidD2NodeShape = errors.New("invalid D2 node shape")

// ParseD2NodeShape converts a string to D2NodeShape, returning an error if invalid.
func ParseD2NodeShape(s string) (D2NodeShape, error) {
	v, err := enum.Parse(d2NodeShapeValues, s, func(ns D2NodeShape) string { return string(ns) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidD2NodeShape, s)
	}

	return v, nil
}

// IsValid returns true if the shape is a valid D2NodeShape value.
func (s D2NodeShape) IsValid() bool {
	return enum.Contains(d2NodeShapeValues, s)
}

// AllowedValues returns all valid D2 node shape values for CLI help text.
func (s D2NodeShape) AllowedValues() []string {
	return enum.AllowedValues(d2NodeShapeValues)
}

// String returns the string representation of the node shape.
func (s D2NodeShape) String() string {
	return string(s)
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

//nolint:gochecknoglobals // Global variable used for value iteration.
var d2ArrowTypeValues = []D2ArrowType{
	D2ArrowNone,
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

// ErrInvalidD2ArrowType is returned when an invalid D2 arrow type is provided.
var ErrInvalidD2ArrowType = errors.New("invalid D2 arrow type")

// ParseD2ArrowType converts a string to D2ArrowType, returning an error if invalid.
func ParseD2ArrowType(s string) (D2ArrowType, error) {
	v, err := enum.Parse(d2ArrowTypeValues, s, func(a D2ArrowType) string { return string(a) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidD2ArrowType, s)
	}

	return v, nil
}

// IsValid returns true if the arrow type is a valid D2ArrowType value.
func (a D2ArrowType) IsValid() bool {
	return a == D2ArrowNone || enum.Contains(d2ArrowTypeValues, a)
}

// AllowedValues returns all valid D2 arrow type values for CLI help text.
func (a D2ArrowType) AllowedValues() []string {
	return enum.AllowedValues(d2ArrowTypeValues)
}

// String returns the string representation of the arrow type.
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

// ErrInvalidD2Constraint is returned when an invalid D2 constraint is provided.
var ErrInvalidD2Constraint = errors.New("invalid D2 constraint")

//nolint:gochecknoglobals // Allowed values for D2Constraint validation.
var allD2Constraints = []D2Constraint{
	D2ConstraintPrimary,
	D2ConstraintForeign,
	D2ConstraintUnique,
}

// AllD2Constraints returns all valid D2Constraint values.
func AllD2Constraints() []D2Constraint {
	return allD2Constraints
}

// ParseD2Constraint converts a string to D2Constraint, returning an error if invalid.
func ParseD2Constraint(s string) (D2Constraint, error) {
	v, err := enum.Parse(allD2Constraints, s, func(d D2Constraint) string { return string(d) })
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidD2Constraint, s)
	}

	return v, nil
}

// String returns the string representation of the constraint.
func (c D2Constraint) String() string { return string(c) }

// AllowedValues returns all valid D2Constraint values for CLI help text.
func (c D2Constraint) AllowedValues() []string {
	return enum.AllowedValues(allD2Constraints)
}

// IsValid returns true if the constraint is a valid D2Constraint value.
func (c D2Constraint) IsValid() bool {
	return enum.Contains(allD2Constraints, c)
}
