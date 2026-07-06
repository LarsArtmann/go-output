package output

// TableBuilder is the CQRS write-side builder for tabular data.
// It provides a fluent construction API, then freezes the result
// into a *Table via Build().
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

// SetFooter sets the footer row.
func (b *TableBuilder) SetFooter(footer ...string) *TableBuilder {
	b.footer = footer
	return b
}

// Build freezes the builder state into a *Table snapshot.
// The returned Table has copied slices — further mutations to the
// builder do not affect previously built Tables.
func (b *TableBuilder) Build() *Table {
	headers := append([]string(nil), b.headers...)

	rows := make([][]string, 0, len(b.rows))
	for _, row := range b.rows {
		rows = append(rows, append([]string(nil), row...))
	}

	var footer []string
	if len(b.footer) > 0 {
		footer = append([]string(nil), b.footer...)
	}

	return &Table{
		Headers: headers,
		Rows:    rows,
		Footer:  footer,
	}
}
