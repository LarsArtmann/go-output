package output

import (
	"fmt"
	"strings"
)

// MarkdownTable builds Markdown tables.
type MarkdownTable struct {
	headers []string
	rows    [][]string
	align   []int
}

// NewMarkdownTable creates a new MarkdownTable.
func NewMarkdownTable() *MarkdownTable {
	return &MarkdownTable{
		headers: nil,
		rows:    nil,
		align:   nil,
	}
}

// SetHeaders sets the table headers.
func (m *MarkdownTable) SetHeaders(headers []string) *MarkdownTable {
	m.headers = headers
	m.align = make([]int, len(headers))
	for i := range m.align {
		m.align[i] = 0
	}
	return m
}

// SetAlign sets column alignment: 0=left, 1=right, 2=center.
func (m *MarkdownTable) SetAlign(col, alignment int) *MarkdownTable {
	if col >= 0 && col < len(m.align) {
		m.align[col] = alignment
	}
	return m
}

// AddRow adds a row to the table.
func (m *MarkdownTable) AddRow(row []string) *MarkdownTable {
	m.rows = append(m.rows, row)
	return m
}

// Render returns the Markdown table string.
func (m *MarkdownTable) Render() string {
	if len(m.headers) == 0 {
		return ""
	}

	colWidths := m.calculateColumnWidths()
	var b strings.Builder

	m.writeHeader(&b, colWidths)
	m.writeSeparator(&b, colWidths)
	m.writeRows(&b, colWidths)

	return b.String()
}

func (m *MarkdownTable) calculateColumnWidths() []int {
	colWidths := make([]int, len(m.headers))
	for i, h := range m.headers {
		colWidths[i] = len(h)
	}
	for _, row := range m.rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}
	return colWidths
}

func (m *MarkdownTable) writeHeader(b *strings.Builder, colWidths []int) {
	b.WriteString("|")
	for i, header := range m.headers {
		b.WriteString(" ")
		b.WriteString(header)
		b.WriteString(strings.Repeat(" ", colWidths[i]-len(header)+1))
		b.WriteString("|")
	}
	b.WriteString("\n")
}

func (m *MarkdownTable) writeSeparator(b *strings.Builder, colWidths []int) {
	b.WriteString("|")
	for _, width := range colWidths {
		b.WriteString("-")
		b.WriteString(strings.Repeat("-", width+1))
		b.WriteString("|")
	}
	b.WriteString("\n")
}

func (m *MarkdownTable) writeRows(b *strings.Builder, colWidths []int) {
	for _, row := range m.rows {
		b.WriteString("|")
		for i, cell := range row {
			b.WriteString(" ")
			m.writeCell(b, i, cell, colWidths)
			b.WriteString(" |")
		}
		b.WriteString("\n")
	}
}

func (m *MarkdownTable) writeCell(b *strings.Builder, i int, cell string, colWidths []int) {
	width := colWidths[i]
	alignment := m.getAlignment(i)

	switch alignment {
	case AlignRight:
		fmt.Fprintf(b, "%*s", width, cell)
	case AlignCenter:
		fmt.Fprintf(b, "%-*s", width, cell)
	default:
		b.WriteString(cell)
		b.WriteString(strings.Repeat(" ", width-len(cell)))
	}
}

func (m *MarkdownTable) getAlignment(col int) int {
	if col >= 0 && col < len(m.align) {
		return m.align[col]
	}
	return AlignLeft
}

// Column alignment constants.
const (
	AlignLeft   = 0
	AlignRight  = 1
	AlignCenter = 2
)
