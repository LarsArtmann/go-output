package output

import (
	"errors"
	"fmt"
)

var errColumnMismatch = errors.New("footer column count does not match headers")

// TableData represents tabular data with headers, rows, and an optional footer.
type TableData struct {
	// Headers are the column header labels.
	Headers []string
	// Rows contains the data rows, each a slice of cell values.
	Rows [][]string
	// Footer is an optional totals/summary row rendered after all data rows.
	// Tabular formats (CSV, TSV, Markdown, HTML, XML, AsciiDoc, Table) render it visually.
	// Data formats (JSON, YAML, TOML, JSONL) and non-tabular formats skip it.
	Footer []string
}

// NewTableData creates a new TableData with the given headers.
func NewTableData(headers []string) *TableData {
	return &TableData{
		Headers: headers,
		Rows:    make([][]string, 0),
	}
}

// AddRow adds a row to the table data.
func (d *TableData) AddRow(row []string) {
	d.Rows = append(d.Rows, row)
}

// RowCount returns the number of data rows.
func (d *TableData) RowCount() int {
	return len(d.Rows)
}

// ColCount returns the number of columns (based on headers).
func (d *TableData) ColCount() int {
	return len(d.Headers)
}

// GetHeaders returns the column headers.
// Satisfies the table.TableDataProvider interface.
func (d *TableData) GetHeaders() []string {
	return d.Headers
}

// GetRows returns the data rows.
// Satisfies the table.TableDataProvider interface.
func (d *TableData) GetRows() [][]string {
	return d.Rows
}

// GetFooter returns the footer row, or nil if none is set.
// Satisfies the table.FooterProvider optional interface.
func (d *TableData) GetFooter() []string {
	return d.Footer
}

// HasFooter returns true if a footer row is present.
func (d *TableData) HasFooter() bool {
	return len(d.Footer) > 0
}

// SetFooter sets the footer row.
func (d *TableData) SetFooter(footer []string) {
	d.Footer = footer
}

// Validate checks the TableData for consistency.
// Returns an error if the footer column count does not match the header count.
func (d *TableData) Validate() error {
	if d == nil || len(d.Headers) == 0 || !d.HasFooter() {
		return nil
	}

	if len(d.Footer) != len(d.Headers) {
		return fmt.Errorf("%w: footer has %d columns, expected %d",
			errColumnMismatch, len(d.Footer), len(d.Headers))
	}

	return nil
}

// ToMapSlice converts TableData to a slice of maps (header→cell).
// Returns nil if data is nil or has no headers.
func (d *TableData) ToMapSlice() []map[string]string {
	if d == nil || len(d.Headers) == 0 {
		return nil
	}

	result := make([]map[string]string, 0, len(d.Rows))

	for _, row := range d.Rows {
		m := make(map[string]string, len(d.Headers))

		for i, header := range d.Headers {
			if i < len(row) {
				m[header] = row[i]
			}
		}

		result = append(result, m)
	}

	return result
}

// RowEdge represents a directed edge between two row identifiers.
type RowEdge struct {
	From string
	To   string
}

// CreateRowEdges generates edge data connecting consecutive rows.
// Used by graph renderers to create edges between table rows.
func (d *TableData) CreateRowEdges() []RowEdge {
	if d == nil || len(d.Rows) < 2 {
		return nil
	}

	edges := make([]RowEdge, 0, len(d.Rows)-1)
	for i := range len(d.Rows) - 1 {
		edges = append(edges, RowEdge{
			From: fmt.Sprintf("row%d", i),
			To:   fmt.Sprintf("row%d", i+1),
		})
	}

	return edges
}

// TableDataBase provides common table data storage for renderers.
type TableDataBase struct {
	data *TableData
}

// ensureData initializes data if nil.
func (b *TableDataBase) ensureData() {
	if b.data == nil {
		b.data = &TableData{}
	}
}

// SetHeaders sets the column headers.
func (b *TableDataBase) SetHeaders(headers []string) {
	b.ensureData()
	b.data.Headers = headers
}

// AddRow adds a data row.
func (b *TableDataBase) AddRow(row []string) {
	b.ensureData()
	b.data.Rows = append(b.data.Rows, row)
}

// SetData sets the table data directly.
func (b *TableDataBase) SetData(data *TableData) {
	b.data = data
}

// Data returns the underlying TableData.
func (b *TableDataBase) Data() *TableData {
	return b.data
}

// SetFooter sets the footer row.
func (b *TableDataBase) SetFooter(footer []string) {
	b.ensureData()
	b.data.Footer = footer
}

// HasFooter returns true if a footer row is present.
func (b *TableDataBase) HasFooter() bool {
	if b.data == nil {
		return false
	}

	return b.data.HasFooter()
}
