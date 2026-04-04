package output

import (
	"fmt"
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

//nolint:gochecknoglobals // Global variable used for value iteration.
var formatValues = []Format{
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
	v, err := enum.Parse(formatValues, s, func(f Format) string { return string(f) })
	if err != nil {
		return "", &InvalidFormatError{Value: s, Allowed: formatValues}
	}

	return v, nil
}

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

// AllowedValues returns all valid output format values for CLI help text.
func (f Format) AllowedValues() []string {
	return enum.AllowedValues(formatValues)
}

// IsValid returns true if the format is a valid Format value.
func (f Format) IsValid() bool {
	return enum.Contains(formatValues, f)
}

// Format category for classification.
type FormatCategory int

const (
	CategoryTable FormatCategory = iota
	CategoryTree
	CategoryGraph
)

var (
	//nolint:gochecknoglobals // Map for table format classification.
	tableFormats = map[Format]bool{
		FormatTable:    true,
		FormatJSON:     true,
		FormatCSV:      true,
		FormatTSV:      true,
		FormatXML:      true,
		FormatMarkdown: true,
		FormatYAML:     true,
		FormatD2:       true,
	}

	//nolint:gochecknoglobals // Map for tree format classification.
	treeFormats = map[Format]bool{
		FormatTree: true,
		FormatHTML: true,
	}

	//nolint:gochecknoglobals // Map for graph format classification.
	graphFormats = map[Format]bool{
		FormatD2:      true,
		FormatMermaid: true,
		FormatDOT:     true,
	}
)

// IsTableFormat returns true if this is a table-based format.
func (f Format) IsTableFormat() bool {
	return tableFormats[f]
}

// IsTreeFormat returns true if this is a tree-based format.
func (f Format) IsTreeFormat() bool {
	return treeFormats[f]
}

// IsGraphFormat returns true if this is a graph/diagram format.
func (f Format) IsGraphFormat() bool {
	return graphFormats[f]
}

// Category returns the category of the format.
func (f Format) Category() FormatCategory {
	if treeFormats[f] {
		return CategoryTree
	}

	if graphFormats[f] {
		return CategoryGraph
	}

	return CategoryTable
}

// InvalidFormatError represents an invalid format error.
type InvalidFormatError struct {
	Value   string
	Allowed []Format
}

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

// Renderer defines the interface for output format renderers.
type Renderer interface {
	// Render returns the formatted output as a string.
	Render() string
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

// CreateRowEdges generates edge data connecting consecutive rows.
// Used by graph renderers to create edges between table rows.
func (d *TableData) CreateRowEdges() []struct{ From, To string } {
	if d == nil || len(d.Rows) < 2 {
		return nil
	}

	edges := make([]struct{ From, To string }, 0, len(d.Rows)-1)
	for i := range len(d.Rows) - 1 {
		edges = append(edges, struct{ From, To string }{
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
