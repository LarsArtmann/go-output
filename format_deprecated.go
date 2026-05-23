package output

// OutputFormat is a backward compatibility alias for Format.
//
// Deprecated: Use Format directly instead of OutputFormat.
//
//revive:disable-next-line:exported
type OutputFormat = Format

// Backward compatibility aliases for format constants.
const (
	OutputFormatTable    = FormatTable
	OutputFormatJSON     = FormatJSON
	OutputFormatCSV      = FormatCSV
	OutputFormatTSV      = FormatTSV
	OutputFormatXML      = FormatXML
	OutputFormatMarkdown = FormatMarkdown
	OutputFormatD2       = FormatD2
	OutputFormatYAML     = FormatYAML
	OutputFormatHTML     = FormatHTML
	OutputFormatTree     = FormatTree
	OutputFormatMermaid  = FormatMermaid
	OutputFormatDOT      = FormatDOT
)

// ParseOutputFormat is a backward compatibility function that calls ParseFormat.
func ParseOutputFormat(s string) (Format, error) {
	return ParseFormat(s)
}

// FormatCategory represents a category for format classification.
//
// Deprecated: Use Shape instead. FormatCategory will be removed in a future version.
type FormatCategory int

// Format category constants for classifying output formats.
//
// Deprecated: Use ShapeTable, ShapeTree, ShapeGraph instead.
const (
	CategoryTable FormatCategory = iota
	CategoryTree
	CategoryGraph
)

// String returns the string representation of the format category.
//
// Deprecated: Use Shape.String() instead.
func (c FormatCategory) String() string {
	switch c {
	case CategoryTable:
		return "table" //nolint:goconst // Deprecated: do not modify
	case CategoryTree:
		return "tree" //nolint:goconst // Deprecated: do not modify
	case CategoryGraph:
		return "graph" //nolint:goconst // Deprecated: do not modify
	default:
		return "unknown"
	}
}

// IsTableFormat returns true if this is a table-based format.
//
// Deprecated: Use f.Supports(ShapeTable) instead.
func (f Format) IsTableFormat() bool {
	return f.Supports(ShapeTable)
}

// IsTreeFormat returns true if this is a tree-based format.
//
// Deprecated: Use f.Supports(ShapeTree) instead.
func (f Format) IsTreeFormat() bool {
	return f.Supports(ShapeTree)
}

// IsGraphFormat returns true if this is a graph/diagram format.
//
// Deprecated: Use f.Supports(ShapeGraph) instead.
func (f Format) IsGraphFormat() bool {
	return f.Supports(ShapeGraph)
}

// Category returns the category of the format.
//
// Deprecated: Use f.Shapes() instead. Category returns the primary shape.
func (f Format) Category() FormatCategory {
	if f.Supports(ShapeGraph) {
		return CategoryGraph
	}

	if f.Supports(ShapeTree) {
		return CategoryTree
	}

	return CategoryTable
}
