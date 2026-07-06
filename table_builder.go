package output

import "slices"

// TableBuilder is the CQRS write-side builder for tabular data.
// It provides a fluent construction API, then copies the result
// into a *Table via Build(). The returned Table is a defensive copy —
// further mutations to the builder do not affect previously built Tables.
// Note: *Table has exported fields; callers SHOULD treat it as read-only
// after Build(), but Go does not enforce this at the type level.
//
// Usage:
//
//	t := NewTableBuilder().
//	    SetHeaders("Name", "Status").
//	    AddRow("Compile", "done").
//	    AddRow("Test", "done").
//	    SetFooter("Total", "2 tasks").
//	    Build()
type TableBuilder struct {
	headers []string
	rows    [][]string
	footer  []string
}

// NewTableBuilder creates a new TableBuilder.
func NewTableBuilder() *TableBuilder {
	return &TableBuilder{
		rows: make([][]string, 0),
	}
}

// SetHeaders sets the column headers.
func (b *TableBuilder) SetHeaders(headers ...string) *TableBuilder {
	b.headers = headers
	return b
}

// AddRow appends a data row.
func (b *TableBuilder) AddRow(row ...string) *TableBuilder {
	b.rows = append(b.rows, row)
	return b
}

// AddRows appends multiple data rows.
func (b *TableBuilder) AddRows(rows ...[]string) *TableBuilder {
	b.rows = append(b.rows, rows...)
	return b
}

// SetFooter sets the footer row.
func (b *TableBuilder) SetFooter(footer ...string) *TableBuilder {
	b.footer = footer
	return b
}

// Build returns a defensive copy of the builder's data as a *Table.
// The returned slices are copied — further mutations to the builder do
// not affect previously built Tables. Callers SHOULD treat the result
// as read-only, though Go does not enforce immutability on *Table.
func (b *TableBuilder) Build() *Table {
	headers := slices.Clone(b.headers)

	rows := make([][]string, 0, len(b.rows))
	for _, row := range b.rows {
		rows = append(rows, slices.Clone(row))
	}

	var footer []string
	if len(b.footer) > 0 {
		footer = slices.Clone(b.footer)
	}

	return &Table{
		Headers: headers,
		Rows:    rows,
		Footer:  footer,
	}
}
