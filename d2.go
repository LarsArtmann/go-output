package output

import (
	"strings"
)

type D2Shape struct {
	Name    string
	Columns []D2Column
}

type D2Column struct {
	Name string
	Type string
}

type D2Diagram struct {
	tables []D2Shape
}

func NewD2Diagram() *D2Diagram {
	return &D2Diagram{}
}

func (d *D2Diagram) AddTable(name string, columns []D2Column) *D2Diagram {
	d.tables = append(d.tables, D2Shape{
		Name:    name,
		Columns: columns,
	})
	return d
}

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
