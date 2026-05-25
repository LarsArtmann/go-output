package output

import (
	"fmt"
	"strings"
)

// Compile-time interface check.
var _ Renderer = (*MarkdownTable)(nil)

// Alignment represents text alignment within a table cell.
type Alignment int

// Column alignment constants.
const (
	AlignLeft   Alignment = AlignmentLeft
	AlignRight  Alignment = AlignmentRight
	AlignCenter Alignment = AlignmentCenter
)

// Column alignment iota values (unexported).
const (
	AlignmentLeft Alignment = iota
	AlignmentRight
	AlignmentCenter
)

// MarkdownTable builds Markdown tables.
type MarkdownTable struct {
	headers   []string
	rows      [][]string
	align     []Alignment
	colorMode ColorMode
}

// NewMarkdownTable creates a new MarkdownTable.
func NewMarkdownTable() *MarkdownTable {
	return &MarkdownTable{
		headers:   nil,
		rows:      nil,
		align:     nil,
		colorMode: ColorModeAuto, //nolint:exhaustruct // align set via SetHeaders
	}
}

// NewMarkdownTableFromData creates a MarkdownTable populated from TableData.
func NewMarkdownTableFromData(data *TableData) *MarkdownTable {
	m := NewMarkdownTable()
	m.SetHeaders(data.Headers)

	for _, row := range data.Rows {
		m.AddRow(row)
	}

	return m
}

// SetColorMode sets the color mode for terminal output.
func (m *MarkdownTable) SetColorMode(mode ColorMode) *MarkdownTable {
	m.colorMode = mode
	return m
}

// SetHeaders sets the table headers.
func (m *MarkdownTable) SetHeaders(headers []string) *MarkdownTable {
	m.headers = headers

	m.align = make([]Alignment, len(headers))
	for i := range m.align {
		m.align[i] = AlignmentLeft
	}

	return m
}

// SetAlign sets column alignment.
func (m *MarkdownTable) SetAlign(col int, alignment Alignment) *MarkdownTable {
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
func (m *MarkdownTable) Render() (string, error) {
	if len(m.headers) == 0 {
		return "", nil
	}

	colWidths := m.calculateColumnWidths()

	var b strings.Builder

	m.writeHeader(&b, colWidths)
	m.writeSeparator(&b, colWidths)
	m.writeRows(&b, colWidths)

	return b.String(), nil
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

func (m *MarkdownTable) useColor() bool {
	return m.colorMode.ShouldColor()
}

func (m *MarkdownTable) writeHeader(b *strings.Builder, colWidths []int) {
	b.WriteString("|")

	for i, header := range m.headers {
		b.WriteString(" ")

		if m.useColor() {
			b.WriteString(ansiBold)
		}

		b.WriteString(header)

		if m.useColor() {
			b.WriteString(ansiReset)
		}

		b.WriteString(strings.Repeat(" ", colWidths[i]-len(header)+1))
		b.WriteString("|")
	}

	b.WriteString("\n")
}

func (m *MarkdownTable) writeSeparator(b *strings.Builder, colWidths []int) {
	if m.useColor() {
		b.WriteString(ansiDim)
	}

	b.WriteString("|")

	for i, width := range colWidths {
		prefix, suffix := m.getAlignmentMarkers(i)
		b.WriteString(prefix)
		b.WriteString(strings.Repeat("-", width+1))
		b.WriteString(suffix)
		b.WriteString("|")
	}

	b.WriteString("\n")

	if m.useColor() {
		b.WriteString(ansiReset)
	}
}

func (m *MarkdownTable) getAlignmentMarkers(col int) (prefix, suffix string) {
	switch m.getAlignment(col) {
	case AlignmentRight:
		return "", ":"
	case AlignmentCenter:
		return ":", ":"
	case AlignmentLeft:
		fallthrough
	default:
		return "", ""
	}
}

func (m *MarkdownTable) writeRows(b *strings.Builder, colWidths []int) {
	for _, row := range m.rows {
		b.WriteString("|")

		for i := range m.headers {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}

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
	case AlignmentRight:
		fmt.Fprintf(b, "%*s", width, cell)
	case AlignmentCenter:
		leftPad := (width - len(cell)) / 2
		rightPad := width - len(cell) - leftPad
		b.WriteString(strings.Repeat(" ", leftPad))
		b.WriteString(cell)
		b.WriteString(strings.Repeat(" ", rightPad))
	case AlignmentLeft:
		fallthrough
	default:
		b.WriteString(cell)
		b.WriteString(strings.Repeat(" ", width-len(cell)))
	}
}

func (m *MarkdownTable) getAlignment(col int) Alignment {
	if col >= 0 && col < len(m.align) {
		return m.align[col]
	}

	return AlignmentLeft
}
