package output

// NewTableWithRow creates a Table pre-populated with the given headers and a
// single data row. It collapses the "NewTable + AddRow" idiom into one call,
// so callers don't repeat two lines for the most common table construction
// pattern:
//
//	t := output.NewTableWithRow([]string{"Name"}, "Alice")
//
// For tables that need additional rows, append them with AddRow. For more
// complex construction (multiple headers, footer, etc.), use NewTableBuilder.
func NewTableWithRow(headers []string, row ...string) *Table {
	return NewTableBuilder().
		SetHeaders(headers...).
		AddRow(row...).
		Build()
}
