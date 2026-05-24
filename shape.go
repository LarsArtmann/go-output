package output

import (
	"slices"

	"github.com/larsartmann/go-output/enum"
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
	return enum.Contains(AllShapes, s)
}

// AllowedValues returns all valid data shape values for CLI help text.
func (s Shape) AllowedValues() []string {
	return enum.AllowedValues(AllShapes)
}

// InvalidShapeError is returned when an invalid data shape is provided.
type InvalidShapeError struct {
	Value string
}

func (e *InvalidShapeError) Error() string {
	return "invalid shape: " + e.Value
}

// ParseShape converts a string to Shape, returning an error if invalid.
func ParseShape(s string) (Shape, error) {
	v, err := enum.Parse(AllShapes, s, func(sh Shape) string { return string(sh) })
	if err != nil {
		return "", &InvalidShapeError{Value: s}
	}

	return v, nil
}

// formatCapabilities maps each format to the data shapes it supports.
//
//nolint:gochecknoglobals // Capability matrix is the single source of truth.
var formatCapabilities = map[Format][]Shape{
	FormatTable:    {ShapeTable},
	FormatJSON:     {ShapeTable, ShapeTree, ShapeGraph},
	FormatCSV:      {ShapeTable},
	FormatTSV:      {ShapeTable},
	FormatXML:      {ShapeTable},
	FormatMarkdown: {ShapeTable},
	FormatD2:       {ShapeTable, ShapeGraph},
	FormatYAML:     {ShapeTable, ShapeTree, ShapeGraph},
	FormatHTML:     {ShapeTable, ShapeTree},
	FormatTree:     {ShapeTree},
	FormatMermaid:  {ShapeTable, ShapeGraph},
	FormatDOT:      {ShapeTable, ShapeGraph},
}

// Supports returns true if the format can render the given data shape.
func (f Format) Supports(s Shape) bool {
	shapes, ok := formatCapabilities[f]
	if !ok {
		return false
	}

	return slices.Contains(shapes, s)
}

// Shapes returns all data shapes this format supports.
func (f Format) Shapes() []Shape {
	return formatCapabilities[f]
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
