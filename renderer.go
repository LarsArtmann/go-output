package output

// Renderer defines the interface for output format renderers.
type Renderer interface {
	// Render returns the formatted output as a string.
	Render() (string, error)
}

// TableRenderer defines the interface for table format renderers.
type TableRenderer interface {
	Renderer
	// SetHeaders sets the column headers.
	SetHeaders(headers []string)
	// AddRow adds a data row.
	AddRow(row []string)
}
