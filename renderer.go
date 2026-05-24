package output

import "fmt"

// Renderer defines the interface for output format renderers.
type Renderer interface {
	// Render returns the formatted output as a string.
	Render() (string, error)
}

// MustRender calls Render on the provided Renderer and panics if it returns an error.
// Useful for tests and examples where rendering failure is unexpected.
func MustRender(r Renderer) string {
	out, err := r.Render()
	if err != nil {
		panic(fmt.Sprintf("MustRender: %v", err))
	}

	return out
}

// TableRenderer defines the interface for table format renderers.
type TableRenderer interface {
	Renderer
	// SetHeaders sets the column headers.
	SetHeaders(headers []string)
	// AddRow adds a data row.
	AddRow(row []string)
}
