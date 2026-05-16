package output

import (
	"fmt"
	"slices"
	"strings"

	"github.com/larsartmann/go-output/enum"
)

// Format represents the available output format options for CLI applications.
type Format string

// Output format constants.
const (
	FormatTable    Format = "table"
	FormatJSON     Format = "json"
	FormatCSV      Format = "csv"
	FormatTSV      Format = "tsv"
	FormatMarkdown Format = "markdown"
	FormatXML      Format = "xml"
	FormatD2       Format = "d2"
	FormatYAML     Format = "yaml"
	FormatHTML     Format = "html"
	FormatTree     Format = "tree"
	FormatMermaid  Format = "mermaid"
	FormatDOT      Format = "dot"
)

// AllFormats contains all valid output format values.
// Use this for testing AllowedValues or generating help text.
var AllFormats = []Format{
	FormatTable,
	FormatJSON,
	FormatCSV,
	FormatTSV,
	FormatMarkdown,
	FormatXML,
	FormatD2,
	FormatYAML,
	FormatHTML,
	FormatTree,
	FormatMermaid,
	FormatDOT,
}

// ParseFormat converts a string to Format, returning an error if invalid.
func ParseFormat(s string) (Format, error) {
	v, err := enum.Parse(AllFormats, s, func(f Format) string { return string(f) })
	if err != nil {
		return "", &InvalidFormatError{Value: s, Allowed: AllFormats}
	}

	return v, nil
}

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

// AllowedValues returns all valid output format values for CLI help text.
func (f Format) AllowedValues() []string {
	return enum.AllowedValues(AllFormats)
}

// IsValid returns true if the format is a valid Format value.
func (f Format) IsValid() bool {
	return enum.Contains(AllFormats, f)
}

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
		return "table"
	case CategoryTree:
		return "tree"
	case CategoryGraph:
		return "graph"
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

// InvalidFormatError represents an invalid format error.
type InvalidFormatError struct {
	Value   string
	Allowed []Format
}

// Error returns a descriptive error message including allowed values.
func (e *InvalidFormatError) Error() string {
	if e.Allowed == nil {
		return "invalid format: " + e.Value
	}

	var b strings.Builder
	b.WriteString("invalid format: ")
	b.WriteString(e.Value)
	b.WriteString(" (allowed: ")

	for i, f := range e.Allowed {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(string(f))
	}

	b.WriteString(")")

	return b.String()
}

// Renderer defines the interface for output format renderers.
type Renderer interface {
	// Render returns the formatted output as a string.
	Render() (string, error)
}

// MustRender calls Render on the provided Renderer and panics if it returns an error.
// Useful for tests and examples where rendering failure is unexpected.
func MustRender(r Renderer) string {
	out, err := r.Render()
	if err != nil {
		panic(fmt.Sprintf("MustRender: %v", err))
	}

	return out
}

// TableRenderer defines the interface for table format renderers.
type TableRenderer interface {
	Renderer
	// SetHeaders sets the column headers.
	SetHeaders(headers []string)
	// AddRow adds a data row.
	AddRow(row []string)
}

// TableData represents tabular data with headers and rows.
type TableData struct {
	Headers []string
	Rows    [][]string
}

// NewTableData creates a new TableData with the given headers.
func NewTableData(headers []string) *TableData {
	return &TableData{
		Headers: headers,
		Rows:    make([][]string, 0),
	}
}

// AddRow adds a row to the table data.
func (d *TableData) AddRow(row []string) {
	d.Rows = append(d.Rows, row)
}

// RowCount returns the number of data rows.
func (d *TableData) RowCount() int {
	return len(d.Rows)
}

// ColCount returns the number of columns (based on headers).
func (d *TableData) ColCount() int {
	return len(d.Headers)
}

// GetHeaders returns the column headers.
// Satisfies the table.TableDataProvider interface.
func (d *TableData) GetHeaders() []string {
	return d.Headers
}

// GetRows returns the data rows.
// Satisfies the table.TableDataProvider interface.
func (d *TableData) GetRows() [][]string {
	return d.Rows
}

// RowEdge represents a directed edge between two row identifiers.
type RowEdge struct {
	From string
	To   string
}

// CreateRowEdges generates edge data connecting consecutive rows.
// Used by graph renderers to create edges between table rows.
func (d *TableData) CreateRowEdges() []RowEdge {
	if d == nil || len(d.Rows) < 2 {
		return nil
	}

	edges := make([]RowEdge, 0, len(d.Rows)-1)
	for i := range len(d.Rows) - 1 {
		edges = append(edges, RowEdge{
			From: fmt.Sprintf("row%d", i),
			To:   fmt.Sprintf("row%d", i+1),
		})
	}

	return edges
}

// TreeOutputRenderer defines the interface for tree format renderers.
type TreeOutputRenderer interface {
	Renderer
	// SetRoot sets the root node of the tree.
	SetRoot(node *TreeNode)
}

// TreeNode represents a node in a tree structure.
type TreeNode struct {
	ID       TreeNodeID
	Label    TreeNodeLabel
	Children []*TreeNode
	Metadata map[string]string
	parent   *TreeNode
}

// NewTreeNode creates a new TreeNode with the given ID and label.
func NewTreeNode(id, label string) *TreeNode {
	return &TreeNode{
		ID:       NewBrandedID[TreeNodeIDBrand](id),
		Label:    NewBrandedID[TreeNodeLabelBrand](label),
		Children: make([]*TreeNode, 0),
		Metadata: make(map[string]string),
		parent:   nil, // parent is set via AddChild
	}
}

// AddChild adds a child node to this node.
func (n *TreeNode) AddChild(child *TreeNode) {
	child.parent = n
	n.Children = append(n.Children, child)
}

// Depth returns the depth of this node in the tree (root = 0).
func (n *TreeNode) Depth() int {
	depth := 0

	current := n
	for current.Parent() != nil {
		depth++
		current = current.Parent()
	}

	return depth
}

// Parent returns the parent node (nil for root).
func (n *TreeNode) Parent() *TreeNode {
	return n.parent
}
