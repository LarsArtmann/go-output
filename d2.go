package output

import (
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

// D2Diagram builds D2 diagram output.
type D2Diagram struct {
	tables []D2Shape
}

// NewD2Diagram creates a new D2Diagram.
func NewD2Diagram() *D2Diagram {
	return &D2Diagram{
		tables: nil,
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

// Render returns the D2 diagram string.
func (d *D2Diagram) Render() string {
	var b strings.Builder

	for _, table := range d.tables {
		b.WriteString(table.Name)
		b.WriteString(": {\n")

		for _, col := range table.Columns {
			b.WriteString("  ")
			b.WriteString(col.Name)
			b.WriteString(": ")
			b.WriteString(col.Type)
			b.WriteString("\n")
		}

		b.WriteString("}\n\n")
	}

	return b.String()
}
