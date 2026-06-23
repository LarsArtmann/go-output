package output

import (
	"slices"
	"strings"
)

// Shape represents a data shape that a format can render.
type Shape string

// Data shape constants for format capability classification.
const (
	ShapeTable Shape = "table" // Tabular data with headers and rows
	ShapeTree  Shape = "tree"  // Hierarchical data with parent-child nodes
	ShapeGraph Shape = "graph" // Network data with nodes and edges
)

// AllShapes contains all valid data shape values.
//
//nolint:gochecknoglobals // Global variable used for value iteration.
var AllShapes = []Shape{
	ShapeTable,
	ShapeTree,
	ShapeGraph,
}

// String returns the string representation of the data shape.
func (s Shape) String() string {
	return string(s)
}

// IsValid returns true if the shape is a valid Shape value.
func (s Shape) IsValid() bool {
	return ContainsEnum(AllShapes, s)
}

// AllowedValues returns all valid data shape values for CLI help text.
func (s Shape) AllowedValues() []string {
	return EnumAllowedValues(AllShapes)
}

// InvalidShapeError is returned when an invalid data shape is provided.
type InvalidShapeError struct {
	Value   string
	Allowed []Shape
}

// Error returns a descriptive error message for the invalid shape.
func (e *InvalidShapeError) Error() string {
	return "invalid shape: " + e.Value + " (allowed: " + shapesToString(e.Allowed) + ")"
}

// ParseShape converts a string to Shape, returning an error if invalid.
func ParseShape(s string) (Shape, error) {
	v, err := ParseEnum(AllShapes, s, func(sh Shape) string { return string(sh) })
	if err != nil {
		return "", &InvalidShapeError{Value: s, Allowed: AllShapes}
	}

	return v, nil
}

func shapesToString(shapes []Shape) string {
	parts := make([]string, len(shapes))
	for i, s := range shapes {
		parts[i] = string(s)
	}

	return strings.Join(parts, ", ")
}

//nolint:gochecknoglobals // Registry for format capabilities, populated by sub-module init().
var formatCapabilities = newFormatRegistry[[]Shape]()

// RegisterFormatShapes registers the data shapes a format supports.
// Sub-modules call this from their init() to declare capabilities.
func RegisterFormatShapes(format Format, shapes ...Shape) {
	formatCapabilities.register(format, shapes)
}

//nolint:gochecknoinits // Registers all format capabilities (sub-modules may override via their own init).
func init() {
	RegisterFormatShapes(FormatTable, ShapeTable)
	RegisterFormatShapes(FormatJSON, ShapeTable, ShapeTree, ShapeGraph)
	RegisterFormatShapes(FormatCSV, ShapeTable)
	RegisterFormatShapes(FormatTSV, ShapeTable)
	RegisterFormatShapes(FormatXML, ShapeTable)
	RegisterFormatShapes(FormatMarkdown, ShapeTable)
	RegisterFormatShapes(FormatD2, ShapeTable, ShapeTree, ShapeGraph)
	RegisterFormatShapes(FormatYAML, ShapeTable, ShapeTree, ShapeGraph)
	RegisterFormatShapes(FormatHTML, ShapeTable, ShapeTree)
	RegisterFormatShapes(FormatTree, ShapeTree)
	RegisterFormatShapes(FormatMermaid, ShapeTable, ShapeTree, ShapeGraph)
	RegisterFormatShapes(FormatDOT, ShapeTable, ShapeTree, ShapeGraph)
	RegisterFormatShapes(FormatJSONL, ShapeTable)
	RegisterFormatShapes(FormatAsciiDoc, ShapeTable)
	RegisterFormatShapes(FormatTOML, ShapeTable, ShapeTree, ShapeGraph)
	RegisterFormatShapes(FormatPlantUML, ShapeTable, ShapeTree, ShapeGraph)
}

func getFormatShapes(format Format) ([]Shape, bool) {
	return formatCapabilities.get(format)
}

// Supports returns true if the format can render the given data shape.
func (f Format) Supports(s Shape) bool {
	shapes, ok := getFormatShapes(f)
	if !ok {
		return false
	}

	return slices.Contains(shapes, s)
}

// Shapes returns all data shapes this format supports.
func (f Format) Shapes() []Shape {
	shapes, _ := getFormatShapes(f)

	return shapes
}

// FormatsForShape returns all formats that support the given data shape.
func FormatsForShape(s Shape) []Format {
	var result []Format

	for _, f := range AllFormats {
		if f.Supports(s) {
			result = append(result, f)
		}
	}

	return result
}
