package d2

import (
	"fmt"
	"sort"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

// Compile-time interface checks.
var (
	_ output.Renderer      = (*Diagram)(nil)
	_ output.GraphRenderer = (*Diagram)(nil)
)

// d2NeedsQuoting reports whether s contains characters that are special in
// D2 syntax and would break parsing if left unquoted. D2 treats # as a
// comment character, and {}, [], (), :, ;, |, ", \, and whitespace as
// structural or syntactic characters. Strings containing any of these must
// be wrapped in double quotes to produce valid D2.
func d2NeedsQuoting(s string) bool {
	if s == "" {
		return true
	}

	for _, r := range s {
		switch r {
		case ' ', '\t', '\n', '\r',
			'#', ':', ';', '|',
			'{', '}', '[', ']', '(', ')',
			'"', '\\':
			return true
		}
	}

	return false
}

// d2Quote wraps s in double quotes (after D2 escaping) when s contains
// characters that require quoting. Simple identifiers (alphanumeric,
// underscores, hyphens) are returned unquoted for readability and backward
// compatibility.
func d2Quote(s string) string {
	if d2NeedsQuoting(s) {
		return `"` + escape.D2(s) + `"`
	}

	return escape.D2(s)
}

// Diagram builds D2 diagram output with full support for nodes, edges,
// SQL table shapes, styling, nesting, icons, links, tooltips, classes, and layout configuration.
type Diagram struct {
	direction Direction
	layout    string
	title     string
	classes   map[string]NodeStyle
	tables    []Table
	nodes     []Node
	edges     []Edge
}

// NewDiagram creates a new Diagram.
func NewDiagram() *Diagram {
	return &Diagram{
		classes: make(map[string]NodeStyle),
		nodes:   make([]Node, 0),
		edges:   make([]Edge, 0),
	}
}

// SetDirection sets the layout direction for the diagram.
func (d *Diagram) SetDirection(dir Direction) *Diagram {
	d.direction = dir
	return d
}

// SetLayout sets the layout engine (e.g., "elk", "dagre").
func (d *Diagram) SetLayout(engine string) *Diagram {
	d.layout = engine
	return d
}

// SetTitle sets the diagram title.
func (d *Diagram) SetTitle(title string) *Diagram {
	d.title = title
	return d
}

// AddClass adds a reusable style class that can be referenced by nodes.
func (d *Diagram) AddClass(name string, style NodeStyle) *Diagram {
	d.classes[name] = style
	return d
}

// AddTable adds a SQL table shape to the diagram.
func (d *Diagram) AddTable(name string, columns []Column) *Diagram {
	d.tables = append(d.tables, Table{Name: name, Columns: columns})
	return d
}

// AddNode adds a node to the diagram.
func (d *Diagram) AddNode(node Node) *Diagram {
	d.nodes = append(d.nodes, node)
	return d
}

// AddNodeSimple adds a simple node with just ID and label.
func (d *Diagram) AddNodeSimple(id, label string) *Diagram {
	return d.AddNode(Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand](id),
		Label: output.NewBrandedID[output.D2NodeLabelBrand](label),
	})
}

// AddNodeWithShape adds a node with a specific shape.
func (d *Diagram) AddNodeWithShape(id, label string, shape NodeShape) *Diagram {
	return d.AddNode(Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand](id),
		Label: output.NewBrandedID[output.D2NodeLabelBrand](label),
		Shape: shape,
	})
}

// AddEdge adds an edge between two nodes.
func (d *Diagram) AddEdge(edge Edge) *Diagram {
	d.edges = append(d.edges, edge)
	return d
}

// AddEdgeSimple adds a simple edge between two nodes.
func (d *Diagram) AddEdgeSimple(from, to string) *Diagram {
	return d.AddEdge( //nolint:exhaustruct // Simple edge uses defaults for optional fields
		Edge{
			From: output.NewBrandedID[output.D2NodeIDBrand](from),
			To:   output.NewBrandedID[output.D2NodeIDBrand](to),
		},
	)
}

// AddLabeledEdge adds an edge with a label.
func (d *Diagram) AddLabeledEdge(from, to, label string) *Diagram {
	return d.AddEdge( //nolint:exhaustruct // Labeled edge uses defaults for optional fields
		Edge{
			From:  output.NewBrandedID[output.D2NodeIDBrand](from),
			To:    output.NewBrandedID[output.D2NodeIDBrand](to),
			Label: output.NewBrandedID[output.D2NodeLabelBrand](label),
		},
	)
}

// Render returns the D2 diagram as a valid D2 language string.
func (d *Diagram) Render() (string, error) {
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

	return strings.TrimRight(b.String(), "\n"), nil
}

func (d *Diagram) writeConfig(b *strings.Builder) {
	hasConfig := d.direction != "" || d.layout != "" || d.title != ""
	if !hasConfig {
		return
	}

	if d.direction != "" && d.direction != DirDown && d.direction.IsValid() {
		fmt.Fprintf(b, "direction: %s\n", d.direction)
	}

	if d.title != "" {
		fmt.Fprintf(b, "title: {\n  label: %s\n}\n", d2Quote(d.title))
	}

	if d.layout != "" {
		fmt.Fprintf(b, "layout: %s\n", d2Quote(d.layout))
	}

	b.WriteString("\n")
}

func (d *Diagram) writeClasses(b *strings.Builder) {
	b.WriteString("classes: {\n")

	names := make([]string, 0, len(d.classes))
	for name := range d.classes {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		b.WriteString("  ")
		b.WriteString(d2Quote(name))
		b.WriteString(": {\n")
		d.writeStyleAttrs(b, d.classes[name], "    ")
		b.WriteString("  }\n")
	}

	b.WriteString("}\n\n")
}

func (d *Diagram) writeTable(b *strings.Builder, table Table) {
	fmt.Fprintf(b, "%s: {\n  shape: sql_table\n", d2Quote(table.Name))

	for _, col := range table.Columns {
		d.writeColumn(b, col)
	}

	b.WriteString("}\n\n")
}

func (*Diagram) writeColumn(b *strings.Builder, col Column) {
	if col.Constraint != "" {
		fmt.Fprintf(b, "  %s: %s {constraint: %s}\n",
			d2Quote(col.Name), d2Quote(col.Type), d2Quote(string(col.Constraint)))
	} else {
		fmt.Fprintf(b, "  %s: %s\n", d2Quote(col.Name), d2Quote(col.Type))
	}
}

func (d *Diagram) writeNode(b *strings.Builder, node Node) {
	if node.Nested != "" {
		d.writeNestedNode(b, node)
		return
	}

	if node.hasBlockAttrs() {
		fmt.Fprintf(b, "%s: %s {\n", d2Quote(node.ID.Get()), d2Quote(node.Label.Get()))
		d.writeNodeAttrs(b, node)
		b.WriteString("}\n")
	} else {
		fmt.Fprintf(b, "%s: %s\n", d2Quote(node.ID.Get()), d2Quote(node.Label.Get()))
	}
}

func (d *Diagram) writeNestedNode(b *strings.Builder, node Node) {
	fmt.Fprintf(b, "%s: %s {\n", d2Quote(node.ID.Get()), d2Quote(node.Label.Get()))
	d.writeNodeAttrs(b, node)
	b.WriteString(node.Nested)
	b.WriteString("}\n")
}

func (d *Diagram) writeNodeAttrs(b *strings.Builder, node Node) {
	if node.Shape != "" && node.Shape != ShapeRectangle && node.Shape.IsValid() {
		fmt.Fprintf(b, "  shape: %s\n", node.Shape)
	}

	d.writeStyleAttrs(b, node.Style, "  ")
	d.writeNodeSize(b, node)
	d.writeNodeLayout(b, node)
	d.writeNodeRefs(b, node)
}

func (*Diagram) writeNodeSize(b *strings.Builder, node Node) {
	if node.Width > 0 {
		fmt.Fprintf(b, "  width: %d\n", node.Width)
	}

	if node.Height > 0 {
		fmt.Fprintf(b, "  height: %d\n", node.Height)
	}
}

func (*Diagram) writeNodeLayout(b *strings.Builder, node Node) {
	if node.Near != "" {
		fmt.Fprintf(b, "  near: %s\n", d2Quote(node.Near))
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

func (*Diagram) writeNodeRefs(b *strings.Builder, node Node) {
	if node.Class != "" {
		fmt.Fprintf(b, "  class: %s\n", d2Quote(node.Class))
	}

	if node.Icon != "" {
		fmt.Fprintf(b, "  icon: %s\n", d2Quote(node.Icon))
	}

	if node.Link != "" {
		fmt.Fprintf(b, "  link: %s\n", d2Quote(node.Link))
	}

	if node.Tooltip != "" {
		fmt.Fprintf(b, "  tooltip: %s\n", d2Quote(node.Tooltip))
	}
}
