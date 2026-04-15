package output

import (
	"fmt"
	"strings"
)

// D2Direction constants for diagram layout direction.
type D2Direction string

const (
	D2DirDown  D2Direction = ""
	D2DirRight D2Direction = "right"
	D2DirLeft  D2Direction = "left"
	D2DirUp    D2Direction = "up"
)

// D2Shape represents a SQL table shape in D2 diagrams.
type D2Shape struct {
	Name    string
	Columns []D2Column
}

// D2Column represents a column in a D2 table shape.
type D2Column struct {
	Name string
	Type string
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
)

// D2NodeStyle represents styling for a D2 node.
type D2NodeStyle struct {
	Fill        string
	Stroke      string
	StrokeWidth int
	FontSize    int
	Opacity     float64
	Shadow      bool
}

// D2Node represents a node in a D2 diagram.
type D2Node struct {
	ID      D2NodeID
	Label   D2NodeLabel
	Shape   D2NodeShape
	Style   D2NodeStyle
	Icon    string
	Link    string
	Tooltip string
	Nested  string
}

// D2EdgeStyle represents styling for a D2 edge.
type D2EdgeStyle struct {
	Stroke      string
	StrokeWidth int
	Animated    bool
	Dashed      bool
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

// D2ArrowType represents the type of arrow for D2 edges.
type D2ArrowType string

// D2ArrowType constants define the available arrow shapes for D2 edges.
const (
	D2ArrowNone     D2ArrowType = ""
	D2ArrowArrow    D2ArrowType = "arrow"
	D2ArrowTriangle D2ArrowType = "triangle"
	D2ArrowDiamond  D2ArrowType = "diamond"
	D2ArrowCircle   D2ArrowType = "circle"
	D2ArrowFilled   D2ArrowType = "filled"
)

// Deprecated: Use D2ArrowArrow instead. D2 uses "arrow" as the standard arrowhead.
const D2ArrowPoint = D2ArrowArrow

// Deprecated: Use D2ArrowCircle instead. D2 uses "circle" as the standard round arrowhead.
const D2ArrowOval = D2ArrowCircle

// D2Diagram builds D2 diagram output with full support for nodes, edges,
// SQL table shapes, styling, nesting, icons, links, tooltips, and layout configuration.
type D2Diagram struct {
	direction D2Direction
	layout    string
	title     string
	tables    []D2Shape
	nodes     []D2Node
	edges     []D2Edge
}

// NewD2Diagram creates a new D2Diagram.
func NewD2Diagram() *D2Diagram {
	return &D2Diagram{
		nodes: make([]D2Node, 0),
		edges: make([]D2Edge, 0),
	}
}

// SetDirection sets the layout direction for the diagram.
func (d *D2Diagram) SetDirection(dir D2Direction) *D2Diagram {
	d.direction = dir
	return d
}

// SetLayout sets the layout engine (e.g., "elk", "dagre").
func (d *D2Diagram) SetLayout(engine string) *D2Diagram {
	d.layout = engine
	return d
}

// SetTitle sets the diagram title.
func (d *D2Diagram) SetTitle(title string) *D2Diagram {
	d.title = title
	return d
}

// AddTable adds a SQL table shape to the diagram.
func (d *D2Diagram) AddTable(name string, columns []D2Column) *D2Diagram {
	d.tables = append(d.tables, D2Shape{Name: name, Columns: columns})
	return d
}

// AddNode adds a node to the diagram.
func (d *D2Diagram) AddNode(node D2Node) *D2Diagram {
	d.nodes = append(d.nodes, node)
	return d
}

// AddNodeSimple adds a simple node with just ID and label.
func (d *D2Diagram) AddNodeSimple(id, label string) *D2Diagram {
	return d.AddNode(D2Node{
		ID:    NewBrandedID[D2NodeIDBrand](id),
		Label: NewBrandedID[D2NodeLabelBrand](label),
	})
}

// AddNodeWithShape adds a node with a specific shape.
func (d *D2Diagram) AddNodeWithShape(id, label string, shape D2NodeShape) *D2Diagram {
	return d.AddNode(D2Node{
		ID:    NewBrandedID[D2NodeIDBrand](id),
		Label: NewBrandedID[D2NodeLabelBrand](label),
		Shape: shape,
	})
}

// AddEdge adds an edge between two nodes.
func (d *D2Diagram) AddEdge(edge D2Edge) *D2Diagram {
	d.edges = append(d.edges, edge)
	return d
}

// AddEdgeSimple adds a simple edge between two nodes.
func (d *D2Diagram) AddEdgeSimple(from, to string) *D2Diagram {
	return d.AddEdge( //nolint:exhaustruct // Simple edge uses defaults for optional fields
		D2Edge{From: NewBrandedID[D2NodeIDBrand](from), To: NewBrandedID[D2NodeIDBrand](to)})
}

// AddLabeledEdge adds an edge with a label.
func (d *D2Diagram) AddLabeledEdge(from, to, label string) *D2Diagram {
	return d.AddEdge( //nolint:exhaustruct // Labeled edge uses defaults for optional fields
		D2Edge{
			From:  NewBrandedID[D2NodeIDBrand](from),
			To:    NewBrandedID[D2NodeIDBrand](to),
			Label: NewBrandedID[D2NodeLabelBrand](label),
		})
}

// Render returns the D2 diagram as a valid D2 language string.
func (d *D2Diagram) Render() string {
	var b strings.Builder

	d.writeConfig(&b)

	for _, table := range d.tables {
		d.writeTable(&b, table)
	}

	for _, node := range d.nodes {
		d.writeNode(&b, node)
	}

	for _, edge := range d.edges {
		d.writeEdge(&b, edge)
	}

	return b.String()
}

func (d *D2Diagram) writeConfig(b *strings.Builder) {
	hasConfig := d.direction != "" || d.layout != "" || d.title != ""
	if !hasConfig {
		return
	}

	if d.direction != "" && d.direction != D2DirDown {
		fmt.Fprintf(b, "direction: %s\n", d.direction)
	}

	if d.title != "" {
		fmt.Fprintf(b, "title: {\n  label: %s\n}\n", escapeD2(d.title))
	}

	if d.layout != "" {
		fmt.Fprintf(b, "layout: %s\n", d.layout)
	}

	b.WriteString("\n")
}

func (d *D2Diagram) writeTable(b *strings.Builder, table D2Shape) {
	fmt.Fprintf(b, "%s: {\n  shape: sql_table\n", escapeD2(table.Name))

	for _, col := range table.Columns {
		fmt.Fprintf(b, "  %s: %s\n", escapeD2(col.Name), escapeD2(col.Type))
	}

	b.WriteString("}\n\n")
}

func (d *D2Diagram) writeNode(b *strings.Builder, node D2Node) {
	if node.Nested != "" {
		d.writeNestedNode(b, node)
		return
	}

	if node.hasBlockAttrs() {
		fmt.Fprintf(b, "%s: %s {\n", escapeD2(node.ID.Get()), escapeD2(node.Label.Get()))
		d.writeNodeShape(b, node.Shape)
		d.writeNodeBlockAttrs(b, node)
		b.WriteString("}\n")
	} else {
		fmt.Fprintf(b, "%s: %s\n", escapeD2(node.ID.Get()), escapeD2(node.Label.Get()))
	}
}

func (d *D2Diagram) writeNestedNode(b *strings.Builder, node D2Node) {
	fmt.Fprintf(b, "%s: %s {\n", escapeD2(node.ID.Get()), escapeD2(node.Label.Get()))
	d.writeNodeShape(b, node.Shape)
	d.writeNodeBlockAttrs(b, node)
	b.WriteString(node.Nested)
	b.WriteString("}\n")
}

func (s D2NodeStyle) isSet() bool {
	return s.Fill != "" || s.Stroke != "" || s.StrokeWidth > 0 ||
		s.FontSize > 0 || s.Opacity > 0 || s.Shadow
}

func (n D2Node) hasBlockAttrs() bool {
	hasShape := n.Shape != "" && n.Shape != D2ShapeRectangle

	return hasShape || n.Style.isSet() || n.Icon != "" || n.Link != "" || n.Tooltip != ""
}

func (*D2Diagram) writeNodeShape(b *strings.Builder, shape D2NodeShape) {
	if shape != "" && shape != D2ShapeRectangle {
		fmt.Fprintf(b, "  shape: %s\n", shape)
	}
}

func (*D2Diagram) writeNodeBlockAttrs(b *strings.Builder, node D2Node) {
	s := node.Style
	if s.Fill != "" {
		fmt.Fprintf(b, "  style.fill: %s\n", s.Fill)
	}

	if s.Stroke != "" {
		fmt.Fprintf(b, "  style.stroke: %s\n", s.Stroke)
	}

	if s.StrokeWidth > 0 {
		fmt.Fprintf(b, "  style.stroke-width: %d\n", s.StrokeWidth)
	}

	if s.FontSize > 0 {
		fmt.Fprintf(b, "  style.font-size: %d\n", s.FontSize)
	}

	if s.Opacity > 0 {
		fmt.Fprintf(b, "  style.opacity: %g\n", s.Opacity)
	}

	if s.Shadow {
		b.WriteString("  style.shadow: true\n")
	}

	if node.Icon != "" {
		fmt.Fprintf(b, "  icon: %s\n", node.Icon)
	}

	if node.Link != "" {
		fmt.Fprintf(b, "  link: %s\n", node.Link)
	}

	if node.Tooltip != "" {
		fmt.Fprintf(b, "  tooltip: %s\n", escapeD2(node.Tooltip))
	}
}

func (d *D2Diagram) writeEdge(b *strings.Builder, edge D2Edge) {
	from := escapeD2(edge.From.Get())
	to := escapeD2(edge.To.Get())

	if !edge.hasBlockAttrs() {
		if !edge.Label.IsEmpty() {
			fmt.Fprintf(b, "%s -> %s: %s\n", from, to, escapeD2(edge.Label.Get()))
		} else {
			fmt.Fprintf(b, "%s -> %s\n", from, to)
		}

		return
	}

	if !edge.Label.IsEmpty() {
		fmt.Fprintf(b, "%s -> %s: %s {\n", from, to, escapeD2(edge.Label.Get()))
	} else {
		fmt.Fprintf(b, "%s -> %s: {\n", from, to)
	}

	d.writeEdgeBlockAttrs(b, edge)
	b.WriteString("}\n")
}

func (e D2Edge) hasBlockAttrs() bool {
	s := e.Style
	hasStyle := s.Stroke != "" || s.StrokeWidth > 0 || s.Animated || s.Dashed
	hasArrows := e.SourceArrow != "" || e.TargetArrow != ""

	return hasStyle || hasArrows
}

func (*D2Diagram) writeEdgeBlockAttrs(b *strings.Builder, edge D2Edge) {
	s := edge.Style
	if s.Stroke != "" {
		fmt.Fprintf(b, "  style.stroke: %s\n", s.Stroke)
	}

	if s.StrokeWidth > 0 {
		fmt.Fprintf(b, "  style.stroke-width: %d\n", s.StrokeWidth)
	}

	if s.Animated {
		b.WriteString("  style.animated: true\n")
	}

	if s.Dashed {
		b.WriteString("  style.stroke-dash: 5\n")
	}

	if edge.SourceArrow != "" {
		fmt.Fprintf(b, "  source-arrowhead.shape: %s\n", edge.SourceArrow)
	}

	if edge.TargetArrow != "" {
		fmt.Fprintf(b, "  target-arrowhead.shape: %s\n", edge.TargetArrow)
	}
}

// escapeD2 escapes special characters for safe inclusion in D2 output.
func escapeD2(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\t", `\t`)

	return s
}
