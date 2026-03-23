package output

import (
	"fmt"
	"slices"
	"strings"
)

// Format represents the available output format options for CLI applications.
type Format string

// Output format constants.
const (
	FormatTable   Format = "table"
	FormatJSON    Format = "json"
	FormatCSV     Format = "csv"
	FormatMarkdown Format = "markdown"
	FormatD2      Format = "d2"
	FormatYAML    Format = "yaml"
	FormatHTML    Format = "html"
	FormatTree    Format = "tree"
	FormatMermaid Format = "mermaid"
	FormatDOT     Format = "dot"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var formatValues = []Format{
	FormatTable,
	FormatJSON,
	FormatCSV,
	FormatMarkdown,
	FormatD2,
	FormatYAML,
	FormatHTML,
	FormatTree,
	FormatMermaid,
	FormatDOT,
}

// ParseFormat converts a string to Format, returning an error if invalid.
func ParseFormat(s string) (Format, error) {
	for _, v := range formatValues {
		if string(v) == s {
			return v, nil
		}
	}
	return "", &InvalidFormatError{Value: s, Allowed: formatValues}
}

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

// AllowedValues returns all valid output format values for CLI help text.
func (f Format) AllowedValues() []string {
	values := make([]string, len(formatValues))
	for i, v := range formatValues {
		values[i] = string(v)
	}
	return values
}

// IsValid returns true if the format is a valid Format value.
func (f Format) IsValid() bool {
	return slices.Contains(formatValues, f)
}

// IsTableFormat returns true if this is a table-based format.
func (f Format) IsTableFormat() bool {
	switch f {
	case FormatTable, FormatJSON, FormatCSV, FormatMarkdown, FormatYAML, FormatD2:
		return true
	default:
		return false
	}
}

// IsTreeFormat returns true if this is a tree-based format.
func (f Format) IsTreeFormat() bool {
	return f == FormatTree || f == FormatHTML
}

// IsGraphFormat returns true if this is a graph/diagram format.
func (f Format) IsGraphFormat() bool {
	switch f {
	case FormatD2, FormatMermaid, FormatDOT:
		return true
	default:
		return false
	}
}

// InvalidFormatError represents an invalid format error.
type InvalidFormatError struct {
	Value   string
	Allowed []Format
}

func (e *InvalidFormatError) Error() string {
	return "invalid format: " + e.Value + " (allowed: " + formatStrings(e.Allowed) + ")"
}

func formatStrings(formats []Format) string {
	var b strings.Builder
	for i, f := range formats {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(string(f))
	}
	return b.String()
}

// Backward compatibility aliases for the renamed type.
type (
	OutputFormat = Format
)

// Backward compatibility aliases for format constants.
const (
	OutputFormatTable   = FormatTable
	OutputFormatJSON    = FormatJSON
	OutputFormatCSV     = FormatCSV
	OutputFormatMarkdown = FormatMarkdown
	OutputFormatD2      = FormatD2
	OutputFormatYAML    = FormatYAML
	OutputFormatHTML    = FormatHTML
	OutputFormatTree    = FormatTree
	OutputFormatMermaid = FormatMermaid
	OutputFormatDOT     = FormatDOT
)

// Backward compatibility function.
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
	for i := 0; i < len(d.Rows)-1; i++ {
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
	ID       string
	Label    string
	Children []*TreeNode
	Metadata map[string]string
}

// NewTreeNode creates a new TreeNode with the given ID and label.
func NewTreeNode(id, label string) *TreeNode {
	return &TreeNode{
		ID:       id,
		Label:    label,
		Children: make([]*TreeNode, 0),
		Metadata: make(map[string]string),
	}
}

// AddChild adds a child node to this node.
func (n *TreeNode) AddChild(child *TreeNode) {
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
	return nil
}

// GraphRenderer defines the interface for graph format renderers.
type GraphRenderer interface {
	Renderer
	// SetNodes sets the graph nodes.
	SetNodes(nodes []GraphNode)
	// SetEdges sets the graph edges.
	SetEdges(edges []GraphEdge)
}

// GraphNode represents a node in a graph.
type GraphNode struct {
	ID       string
	Label    string
	Shape    GraphShape
	Style    GraphStyle
	Metadata map[string]string
}

// NewGraphNode creates a new GraphNode.
func NewGraphNode(id, label string) *GraphNode {
	return &GraphNode{
		ID:       id,
		Label:    label,
		Shape:    ShapeBox,
		Style:    GraphStyle{},
		Metadata: make(map[string]string),
	}
}

// GraphShape represents the shape of a graph node.
type GraphShape string

const (
	ShapeBox       GraphShape = "box"
	ShapeEllipse   GraphShape = "ellipse"
	ShapeDiamond   GraphShape = "diamond"
	ShapeCircle    GraphShape = "circle"
	ShapeCylinder  GraphShape = "cylinder"
	ShapeHexagon   GraphShape = "hexagon"
	ShapeParallelogram GraphShape = "parallelogram"
	ShapeRect      GraphShape = "rect"
)

// GraphStyle represents styling attributes for a graph node.
type GraphStyle struct {
	FillColor   string
	StrokeColor string
	FontColor   string
	FontSize    int
}

// GraphEdge represents an edge between two nodes.
type GraphEdge struct {
	From       string
	To         string
	Label      string
	Style      EdgeStyle
}

// NewGraphEdge creates a new GraphEdge.
func NewGraphEdge(from, to string) *GraphEdge {
	return &GraphEdge{
		From:  from,
		To:    to,
		Style: EdgeStyle{},
	}
}

// EdgeStyle represents styling attributes for an edge.
type EdgeStyle struct {
	Color      string
	Style      string // solid, dashed, dotted
	ArrowHead  string
	ArrowTail  string
}
