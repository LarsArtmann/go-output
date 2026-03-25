package output

import (
	"fmt"
	"strings"
)

// D2Shape represents a table shape in D2 diagrams.
type D2Shape struct {
	Name    string
	Columns []D2Column
}

// D2Column represents a column in a D2 table shape.
type D2Column struct {
	Name string
	Type string
}

// D2Node represents a node in a D2 diagram.
type D2Node struct {
	ID     D2NodeID
	Label  D2NodeLabel
	Shape  D2NodeShape
	Style  D2NodeStyle
	Nested string
}

// D2NodeShape represents the shape of a D2 node.
type D2NodeShape string

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
)

// D2NodeStyle represents styling for a D2 node.
type D2NodeStyle struct {
	Fill        string
	Stroke      string
	StrokeWidth int
	FontSize    int
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

// D2EdgeStyle represents styling for a D2 edge.
type D2EdgeStyle struct {
	Stroke      string
	StrokeWidth int
	Animated    bool
	Dashed      bool
}

// D2ArrowType represents the type of arrow.
type D2ArrowType string

const (
	D2ArrowNone     D2ArrowType = "none"
	D2ArrowPoint    D2ArrowType = "arrow"
	D2ArrowTriangle D2ArrowType = "triangle"
	D2ArrowDiamond  D2ArrowType = "diamond"
	D2ArrowOval     D2ArrowType = "oval"
)

// D2Diagram builds D2 diagram output.
type D2Diagram struct {
	tables []D2Shape
	nodes  []D2Node
	edges  []D2Edge
}

// NewD2Diagram creates a new D2Diagram.
func NewD2Diagram() *D2Diagram {
	return &D2Diagram{
		tables: nil,
		nodes:  make([]D2Node, 0),
		edges:  make([]D2Edge, 0),
	}
}

// AddTable adds a table to the diagram.
func (d *D2Diagram) AddTable(name string, columns []D2Column) *D2Diagram {
	d.tables = append(d.tables, D2Shape{
		Name:    name,
		Columns: columns,
	})
	return d
}

// AddNode adds a node to the diagram.
func (d *D2Diagram) AddNode(node D2Node) *D2Diagram {
	d.nodes = append(d.nodes, node)
	return d
}

// AddNodeSimple adds a simple node with just ID and label.
func (d *D2Diagram) AddNodeSimple(id, label string) *D2Diagram {
	return d.AddNode(
		D2Node{
			ID:    NewBrandedID[D2NodeIDBrand](id),
			Label: NewBrandedID[D2NodeLabelBrand](label),
			Shape: D2ShapeRectangle,
		},
	)
}

// AddNodeWithShape adds a node with a specific shape.
func (d *D2Diagram) AddNodeWithShape(id, label string, shape D2NodeShape) *D2Diagram {
	return d.AddNode(
		D2Node{
			ID:    NewBrandedID[D2NodeIDBrand](id),
			Label: NewBrandedID[D2NodeLabelBrand](label),
			Shape: shape,
		},
	)
}

// AddEdge adds an edge between two nodes.
func (d *D2Diagram) AddEdge(edge D2Edge) *D2Diagram {
	d.edges = append(d.edges, edge)
	return d
}

// AddEdgeSimple adds a simple edge between two nodes.
func (d *D2Diagram) AddEdgeSimple(from, to string) *D2Diagram {
	return d.AddEdge(
		D2Edge{From: NewBrandedID[D2NodeIDBrand](from), To: NewBrandedID[D2NodeIDBrand](to)},
	)
}

// AddLabeledEdge adds an edge with a label.
func (d *D2Diagram) AddLabeledEdge(from, to, label string) *D2Diagram {
	return d.AddEdge(
		D2Edge{
			From:  NewBrandedID[D2NodeIDBrand](from),
			To:    NewBrandedID[D2NodeIDBrand](to),
			Label: NewBrandedID[D2NodeLabelBrand](label),
		},
	)
}

// Render returns the D2 diagram string.
func (d *D2Diagram) Render() string {
	var b strings.Builder

	// Write tables
	for _, table := range d.tables {
		b.WriteString(table.Name)
		b.WriteString(": {\n")

		for _, col := range table.Columns {
			fmt.Fprintf(&b, "  %s: %s\n", col.Name, col.Type)
		}

		b.WriteString("}\n\n")
	}

	// Write nodes
	for _, node := range d.nodes {
		d.renderNode(&b, node)
	}

	// Write edges
	for _, edge := range d.edges {
		d.renderEdge(&b, edge)
	}

	return b.String()
}

func (d *D2Diagram) renderNode(b *strings.Builder, node D2Node) {
	if node.Nested != "" {
		_, _ = fmt.Fprintf(
			b,
			"%s.%s%s %s {\n",
			node.ID.Get(),
			node.Label.Get(),
			d.renderShapeAttr(node.Shape),
			d.renderStyle(node.Style),
		)
		b.WriteString(node.Nested)
		b.WriteString("}\n")
		return
	}

	shapeAttr := d.renderShapeAttr(node.Shape)
	styleStr := d.renderStyle(node.Style)
	_, _ = fmt.Fprintf(b, "%s%s%s %s\n", node.ID.Get(), shapeAttr, node.Label.Get(), styleStr)
}

func (d *D2Diagram) renderShapeAttr(shape D2NodeShape) string {
	if shape == "" || shape == D2ShapeRectangle {
		return " "
	}
	return fmt.Sprintf(":%s ", shape)
}

func (d *D2Diagram) renderStyle(style D2NodeStyle) string {
	var parts []string
	if style.Fill != "" {
		parts = append(parts, "fill:"+style.Fill)
	}
	if style.Stroke != "" {
		parts = append(parts, "stroke:"+style.Stroke)
	}
	if style.StrokeWidth > 0 {
		parts = append(parts, fmt.Sprintf("stroke-width:%d", style.StrokeWidth))
	}
	if style.FontSize > 0 {
		parts = append(parts, fmt.Sprintf("font-size:%d", style.FontSize))
	}
	if len(parts) == 0 {
		return ""
	}
	return "{" + strings.Join(parts, "; ") + "}"
}

func (d *D2Diagram) renderEdge(b *strings.Builder, edge D2Edge) {
	sourceArrow := d.renderArrow(edge.SourceArrow)
	targetArrow := d.renderArrow(edge.TargetArrow)

	if !edge.Label.IsEmpty() {
		_, _ = fmt.Fprintf(
			b,
			"%s %s-> %s: %s %s\n",
			edge.From.Get(),
			sourceArrow,
			edge.To.Get(),
			edge.Label.Get(),
			targetArrow,
		)
	} else {
		_, _ = fmt.Fprintf(
			b,
			"%s %s-> %s %s\n",
			edge.From.Get(),
			sourceArrow,
			edge.To.Get(),
			targetArrow,
		)
	}
}

func (d *D2Diagram) renderArrow(arrow D2ArrowType) string {
	if arrow == "" || arrow == D2ArrowNone {
		return ""
	}
	return fmt.Sprintf("-%s", arrow)
}

// D2FromTableData converts TableData to a D2 diagram.
func D2FromTableData(data *TableData) *D2Diagram {
	diagram := NewD2Diagram()
	if data == nil {
		return diagram
	}

	columns := make([]D2Column, len(data.Headers))
	for i, h := range data.Headers {
		columns[i] = D2Column{Name: h, Type: "string"}
	}
	diagram.AddTable(data.Headers[0], columns)

	return diagram
}
