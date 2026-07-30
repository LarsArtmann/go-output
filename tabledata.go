package output

import (
	"errors"
	"fmt"
)

var (
	// ErrColumnMismatch is the sentinel for a row or footer whose column
	// count does not match the header count. Match via errors.Is.
	ErrColumnMismatch = errors.New("column count does not match headers")
	// ErrNilRow is the sentinel for a nil row encountered during validation.
	// Match via errors.Is.
	ErrNilRow = errors.New("nil row in data")
)

// Table represents tabular data with headers, rows, and an optional footer.
type Table struct {
	// Headers are the column header labels.
	Headers []string
	// Rows contains the data rows, each a slice of cell values.
	Rows [][]string
	// Footer is an optional totals/summary row rendered after all data rows.
	// Tabular formats (CSV, TSV, Markdown, HTML, XML, AsciiDoc, Table) render it visually.
	// Data formats (JSON, YAML, TOML, JSONL) and non-tabular formats skip it.
	Footer []string
}

// NewTable creates a new Table with the given headers.
func NewTable(headers []string) *Table {
	return &Table{
		Headers: headers,
		Rows:    make([][]string, 0),
	}
}

// AddRow adds a row to the table data.
//
// This method does NOT validate the row's column count. Use AddRowChecked
// to return an error, or call Validate() before rendering to surface
// mismatched rows. Existing callers that rely on silent acceptance are
// preserved; new code should prefer AddRowChecked for fail-fast behavior.
func (d *Table) AddRow(row []string) {
	d.Rows = append(d.Rows, row)
}

// AddRowChecked adds a row to the table data and returns an error if the
// row's column count does not match the header count.
//
// Returns nil if no headers are set, deferring validation to Validate().
// Returns ErrColumnMismatch if the row length differs from len(Headers);
// callers can match via errors.Is(err, output.ErrColumnMismatch).
func (d *Table) AddRowChecked(row []string) error {
	if len(d.Headers) > 0 && len(row) != len(d.Headers) {
		return fmt.Errorf("%w: row has %d columns, expected %d",
			ErrColumnMismatch, len(row), len(d.Headers))
	}

	d.Rows = append(d.Rows, row)

	return nil
}

// RowCount returns the number of data rows.
func (d *Table) RowCount() int {
	if d == nil {
		return 0
	}

	return len(d.Rows)
}

// ColCount returns the number of columns (based on headers).
func (d *Table) ColCount() int {
	if d == nil {
		return 0
	}

	return len(d.Headers)
}

// GetHeaders returns the column headers.
// Satisfies the table.TableProvider interface.
func (d *Table) GetHeaders() []string {
	if d == nil {
		return nil
	}

	return d.Headers
}

// GetRows returns the data rows.
// Satisfies the table.TableProvider interface.
func (d *Table) GetRows() [][]string {
	if d == nil {
		return nil
	}

	return d.Rows
}

// GetFooter returns the footer row, or nil if none is set.
// Satisfies the table.FooterProvider optional interface.
func (d *Table) GetFooter() []string {
	if d == nil {
		return nil
	}

	return d.Footer
}

// HasFooter returns true if a footer row is present.
func (d *Table) HasFooter() bool {
	if d == nil {
		return false
	}

	return len(d.Footer) > 0
}

// SetFooter sets the footer row.
func (d *Table) SetFooter(footer []string) {
	if d == nil {
		return
	}

	d.Footer = footer
}

// Validate checks the Table for consistency.
// Returns ErrNilRow (wrap-matched via errors.Is) if any row is nil.
// Returns ErrColumnMismatch if the footer column count does not match the header count.
func (d *Table) Validate() error {
	if d == nil {
		return nil
	}

	for i, row := range d.Rows {
		if row == nil {
			return fmt.Errorf("%w at index %d", ErrNilRow, i)
		}
	}

	if len(d.Headers) == 0 || !d.HasFooter() {
		return nil
	}

	if len(d.Footer) != len(d.Headers) {
		return fmt.Errorf("%w: footer has %d columns, expected %d",
			ErrColumnMismatch, len(d.Footer), len(d.Headers))
	}

	return nil
}

// ToMapSlice converts Table to a slice of maps (header→cell).
// Returns nil if data is nil or has no headers.
func (d *Table) ToMapSlice() []map[string]string {
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
	From GraphNodeID
	To   GraphNodeID
}

// CreateRowEdges generates edge data connecting consecutive rows.
// Used by graph renderers to create edges between table rows.
func (d *Table) CreateRowEdges() []RowEdge {
	if d == nil || len(d.Rows) < 2 {
		return nil
	}

	edges := make([]RowEdge, 0, len(d.Rows)-1)
	for i := range len(d.Rows) - 1 {
		edges = append(edges, RowEdge{
			From: NewBrandedID[GraphNodeIDBrand](fmt.Sprintf("row%d", i)),
			To:   NewBrandedID[GraphNodeIDBrand](fmt.Sprintf("row%d", i+1)),
		})
	}

	return edges
}

// TableStore provides common table data storage for renderers.
type TableStore struct {
	data *Table
}

// ensureData initializes data if nil.
func (b *TableStore) ensureData() {
	if b.data == nil {
		b.data = &Table{}
	}
}

// SetHeaders sets the column headers.
func (b *TableStore) SetHeaders(headers []string) {
	b.ensureData()
	b.data.Headers = headers
}

// AddRow adds a data row.
func (b *TableStore) AddRow(row []string) {
	b.ensureData()
	b.data.Rows = append(b.data.Rows, row)
}

// SetData sets the table data directly.
func (b *TableStore) SetData(data *Table) {
	b.data = data
}

// Data returns the underlying Table.
func (b *TableStore) Data() *Table {
	return b.data
}

// SetFooter sets the footer row.
func (b *TableStore) SetFooter(footer []string) {
	b.ensureData()
	b.data.Footer = footer
}

// HasFooter returns true if a footer row is present.
func (b *TableStore) HasFooter() bool {
	if b.data == nil {
		return false
	}

	return b.data.HasFooter()
}
