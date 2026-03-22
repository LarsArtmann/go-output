package output

import (
	"fmt"
	"strings"
)

type MarkdownTable struct {
	headers []string
	rows    [][]string
	align   []int
}

func NewMarkdownTable() *MarkdownTable {
	return &MarkdownTable{}
}

func (m *MarkdownTable) SetHeaders(headers []string) *MarkdownTable {
	m.headers = headers
	m.align = make([]int, len(headers))
	for i := range m.align {
		m.align[i] = 0
	}
	return m
}

func (m *MarkdownTable) SetAlign(col int, alignment int) *MarkdownTable {
	if col >= 0 && col < len(m.align) {
		m.align[col] = alignment
	}
	return m
}

func (m *MarkdownTable) AddRow(row []string) *MarkdownTable {
	m.rows = append(m.rows, row)
	return m
}

func (m *MarkdownTable) Render() (string, error) {
	if len(m.headers) == 0 {
		return "", nil
	}

	var b strings.Builder

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

	b.WriteString("|")
	for i, header := range m.headers {
		b.WriteString(" ")
		b.WriteString(header)
		b.WriteString(strings.Repeat(" ", colWidths[i]-len(header)+1))
		b.WriteString("|")
	}
	b.WriteString("\n")

	b.WriteString("|")
	for _, width := range colWidths {
		b.WriteString("-")
		b.WriteString(strings.Repeat("-", width+1))
		b.WriteString("|")
	}
	b.WriteString("\n")

	for _, row := range m.rows {
		b.WriteString("|")
		for i, cell := range row {
			b.WriteString(" ")
			width := colWidths[i]
			if i < len(m.align) {
				switch m.align[i] {
				case 1:
					b.WriteString(fmt.Sprintf("%*s", width, cell))
				case 2:
					b.WriteString(fmt.Sprintf("%-*s", width, cell))
				default:
					b.WriteString(cell)
					b.WriteString(strings.Repeat(" ", width-len(cell)))
				}
			} else {
				b.WriteString(cell)
				b.WriteString(strings.Repeat(" ", width-len(cell)))
			}
			b.WriteString(" |")
		}
		b.WriteString("\n")
	}

	return b.String(), nil
}

const (
	AlignLeft  = 0
	AlignRight = 1
	AlignCenter = 2
)
