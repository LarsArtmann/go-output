package output

import (
	"fmt"
	"strings"

	"github.com/larsartmann/go-output/internal/escape"
)

// D2Diagram builds D2 diagram output with full support for nodes, edges,
// SQL table shapes, styling, nesting, icons, links, tooltips, classes, and layout configuration.
type D2Diagram struct {
	direction D2Direction
	layout    string
	title     string
	classes   map[string]D2NodeStyle
	tables    []D2Table
	nodes     []D2Node
	edges     []D2Edge
}

// NewD2Diagram creates a new D2Diagram.
func NewD2Diagram() *D2Diagram {
	return &D2Diagram{
		classes: make(map[string]D2NodeStyle),
		nodes:   make([]D2Node, 0),
		edges:   make([]D2Edge, 0),
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

// AddClass adds a reusable style class that can be referenced by nodes.
func (d *D2Diagram) AddClass(name string, style D2NodeStyle) *D2Diagram {
	d.classes[name] = style
	return d
}

// AddTable adds a SQL table shape to the diagram.
func (d *D2Diagram) AddTable(name string, columns []D2Column) *D2Diagram {
	d.tables = append(d.tables, D2Table{Name: name, Columns: columns})
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

	if len(d.classes) > 0 {
		d.writeClasses(&b)
	}

	for _, table := range d.tables {
		d.writeTable(&b, table)
	}

	for _, node := range d.nodes {
		d.writeNode(&b, node)
	}

	for _, edge := range d.edges {
		d.writeEdge(&b, edge)
	}

	return strings.TrimRight(b.String(), "\n")
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
		fmt.Fprintf(b, "title: {\n  label: %s\n}\n", escape.D2(d.title))
	}

	if d.layout != "" {
		fmt.Fprintf(b, "layout: %s\n", d.layout)
	}

	b.WriteString("\n")
}

func (d *D2Diagram) writeClasses(b *strings.Builder) {
	b.WriteString("classes: {\n")

	for name, style := range d.classes {
		b.WriteString("  " + escape.D2(name) + ": {\n")
		d.writeStyleAttrs(b, style, "    ")
		b.WriteString("  }\n")
	}

	b.WriteString("}\n\n")
}

func (d *D2Diagram) writeTable(b *strings.Builder, table D2Table) {
	fmt.Fprintf(b, "%s: {\n  shape: sql_table\n", escape.D2(table.Name))

	for _, col := range table.Columns {
		d.writeColumn(b, col)
	}

	b.WriteString("}\n\n")
}

func (*D2Diagram) writeColumn(b *strings.Builder, col D2Column) {
	if col.Constraint != "" {
		fmt.Fprintf(b, "  %s: %s {constraint: %s}\n",
			escape.D2(col.Name), escape.D2(col.Type), string(col.Constraint))
	} else {
		fmt.Fprintf(b, "  %s: %s\n", escape.D2(col.Name), escape.D2(col.Type))
	}
}

func (d *D2Diagram) writeNode(b *strings.Builder, node D2Node) {
	if node.Nested != "" {
		d.writeNestedNode(b, node)
		return
	}

	if node.hasBlockAttrs() {
		fmt.Fprintf(b, "%s: %s {\n", escape.D2(node.ID.Get()), escape.D2(node.Label.Get()))
		d.writeNodeAttrs(b, node)
		b.WriteString("}\n")
	} else {
		fmt.Fprintf(b, "%s: %s\n", escape.D2(node.ID.Get()), escape.D2(node.Label.Get()))
	}
}

func (d *D2Diagram) writeNestedNode(b *strings.Builder, node D2Node) {
	fmt.Fprintf(b, "%s: %s {\n", escape.D2(node.ID.Get()), escape.D2(node.Label.Get()))
	d.writeNodeAttrs(b, node)
	b.WriteString(node.Nested)
	b.WriteString("}\n")
}

func (d *D2Diagram) writeNodeAttrs(b *strings.Builder, node D2Node) {
	if node.Shape != "" && node.Shape != D2ShapeRectangle {
		fmt.Fprintf(b, "  shape: %s\n", node.Shape)
	}

	d.writeStyleAttrs(b, node.Style, "  ")
	d.writeNodeSize(b, node)
	d.writeNodeLayout(b, node)
	d.writeNodeRefs(b, node)
}

func (*D2Diagram) writeNodeSize(b *strings.Builder, node D2Node) {
	if node.Width > 0 {
		fmt.Fprintf(b, "  width: %d\n", node.Width)
	}

	if node.Height > 0 {
		fmt.Fprintf(b, "  height: %d\n", node.Height)
	}
}

func (*D2Diagram) writeNodeLayout(b *strings.Builder, node D2Node) {
	if node.Near != "" {
		fmt.Fprintf(b, "  near: %s\n", escape.D2(node.Near))
	}

	if node.GridRows > 0 {
		fmt.Fprintf(b, "  grid-rows: %d\n", node.GridRows)
	}

	if node.GridColumns > 0 {
		fmt.Fprintf(b, "  grid-columns: %d\n", node.GridColumns)
	}

	if node.GridGap > 0 {
		fmt.Fprintf(b, "  grid-gap: %d\n", node.GridGap)
	}
}

func (*D2Diagram) writeNodeRefs(b *strings.Builder, node D2Node) {
	if node.Class != "" {
		fmt.Fprintf(b, "  class: %s\n", escape.D2(node.Class))
	}

	if node.Icon != "" {
		fmt.Fprintf(b, "  icon: %s\n", node.Icon)
	}

	if node.Link != "" {
		fmt.Fprintf(b, "  link: %s\n", node.Link)
	}

	if node.Tooltip != "" {
		fmt.Fprintf(b, "  tooltip: %s\n", escape.D2(node.Tooltip))
	}
}
