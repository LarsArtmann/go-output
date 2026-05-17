package output

import "fmt"

// TableData represents tabular data with headers and rows.
type TableData struct {
	Headers []string
	Rows    [][]string
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

// tableDataBase provides common table data storage for renderers.
type tableDataBase struct {
	data *TableData
}

// ensureData initializes data if nil.
func (b *tableDataBase) ensureData() {
	if b.data == nil {
		b.data = &TableData{}
	}
}

// SetHeaders sets the column headers.
func (b *tableDataBase) SetHeaders(headers []string) {
	b.ensureData()
	b.data.Headers = headers
}

// AddRow adds a data row.
func (b *tableDataBase) AddRow(row []string) {
	b.ensureData()
	b.data.Rows = append(b.data.Rows, row)
}

// SetData sets the table data directly.
func (b *tableDataBase) SetData(data *TableData) {
	b.data = data
}
