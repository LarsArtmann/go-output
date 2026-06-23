// Package markdown renders TableData as Markdown tables.
//
// It is an optional format renderer: import it to activate Markdown output
// through output.RenderTableData, or use NewMarkdownTable directly.
//
//	import "github.com/larsartmann/go-output/markdown"
//
//	_ = markdown.NewMarkdownTable()
package markdown

import (
	"fmt"
	"strings"

	"github.com/larsartmann/go-output"
)

// Standard ANSI escape sequences for terminal styling. Local to this module
// so markdown stays dependency-free beyond the core output types.
const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiDim   = "\033[2m"
)

// Compile-time interface check.
var _ output.Renderer = (*MarkdownTable)(nil)

// Alignment represents text alignment within a table cell.
type Alignment int

// Column alignment constants.
const (
	AlignLeft   Alignment = alignmentLeft
	AlignRight  Alignment = alignmentRight
	AlignCenter Alignment = alignmentCenter
)

// Column alignment iota values.
const (
	alignmentLeft Alignment = iota
	alignmentRight
	alignmentCenter
)

// MarkdownTable builds Markdown tables.
type MarkdownTable struct {
	headers   []string
	rows      [][]string
	footer    []string
	align     []Alignment
	colorMode output.ColorMode
}

// NewMarkdownTable creates a new MarkdownTable.
func NewMarkdownTable() *MarkdownTable {
	return &MarkdownTable{
		headers:   nil,
		rows:      nil,
		align:     nil,
		colorMode: output.ColorModeAuto, //nolint:exhaustruct // align set via SetHeaders
	}
}

// NewMarkdownTableFromData creates a MarkdownTable populated from TableData.
func NewMarkdownTableFromData(data *output.TableData) *MarkdownTable {
	m := NewMarkdownTable()
	m.SetHeaders(data.Headers)

	for _, row := range data.Rows {
		m.AddRow(row)
	}

	if data.HasFooter() {
		m.SetFooter(data.Footer)
	}

	return m
}

// SetColorMode sets the color mode for terminal output.
func (m *MarkdownTable) SetColorMode(mode output.ColorMode) *MarkdownTable {
	m.colorMode = mode
	return m
}

// SetHeaders sets the table headers.
func (m *MarkdownTable) SetHeaders(headers []string) *MarkdownTable {
	m.headers = headers

	m.align = make([]Alignment, 0, len(headers))
	for range headers {
		m.align = append(m.align, alignmentLeft)
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

// SetFooter sets the footer row for the table.
func (m *MarkdownTable) SetFooter(footer []string) *MarkdownTable {
	m.footer = footer

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

	if len(m.footer) > 0 {
		m.writeSeparator(&b, colWidths)
		m.writeFooterRow(&b, colWidths)
	}

	return b.String(), nil
}

func (m *MarkdownTable) calculateColumnWidths() []int {
	colWidths := make([]int, 0, len(m.headers))
	for _, h := range m.headers {
		colWidths = append(colWidths, len(h))
	}

	for _, row := range m.rows {
		updateMaxWidths(colWidths, row)
	}

	updateMaxWidths(colWidths, m.footer)

	return colWidths
}

func updateMaxWidths(colWidths []int, cells []string) {
	for i, cell := range cells {
		if i < len(colWidths) && len(cell) > colWidths[i] {
			colWidths[i] = len(cell)
		}
	}
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
	case alignmentRight:
		return "", ":"
	case alignmentCenter:
		return ":", ":"
	case alignmentLeft:
		fallthrough
	default:
		return "", ""
	}
}

func (m *MarkdownTable) writeRows(b *strings.Builder, colWidths []int) {
	for _, row := range m.rows {
		m.writeSingleRow(b, row, colWidths)
	}
}

func (m *MarkdownTable) writeFooterRow(b *strings.Builder, colWidths []int) {
	if m.useColor() {
		b.WriteString(ansiBold)
	}

	m.writeSingleRow(b, m.footer, colWidths)

	if m.useColor() {
		b.WriteString(ansiReset)
	}
}

func (m *MarkdownTable) writeSingleRow(b *strings.Builder, row []string, colWidths []int) {
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

func (m *MarkdownTable) writeCell(b *strings.Builder, i int, cell string, colWidths []int) {
	width := colWidths[i]
	alignment := m.getAlignment(i)

	switch alignment {
	case alignmentRight:
		fmt.Fprintf(b, "%*s", width, cell)
	case alignmentCenter:
		leftPad := (width - len(cell)) / 2
		rightPad := width - len(cell) - leftPad
		b.WriteString(strings.Repeat(" ", leftPad))
		b.WriteString(cell)
		b.WriteString(strings.Repeat(" ", rightPad))
	case alignmentLeft:
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

	return alignmentLeft
}

// markdownTableAdapter wraps a MarkdownTable to satisfy the TableRenderer interface.
// It adapts the fluent API (returning *MarkdownTable) to the void-returning TableRenderer methods.
type markdownTableAdapter struct {
	inner *MarkdownTable
}

func (a *markdownTableAdapter) Render() (string, error)     { return a.inner.Render() }
func (a *markdownTableAdapter) SetHeaders(headers []string) { a.inner.SetHeaders(headers) }
func (a *markdownTableAdapter) AddRow(row []string)         { a.inner.AddRow(row) }

// AsTableRenderer returns a TableRenderer that delegates to this MarkdownTable.
// This adapts the fluent API (returning *MarkdownTable) to the TableRenderer interface.
func (m *MarkdownTable) AsTableRenderer() output.TableRenderer {
	return &markdownTableAdapter{inner: m}
}
