package table

import (
	"io"

	"github.com/larsartmann/go-output"
)

// Write writes a Table as a styled terminal table to the provided writer.
// Uses lipgloss for rendering with optional color mode control.
func Write(w io.Writer, data *output.Table, opts ...Option) error {
	if data == nil {
		return nil
	}

	t := FromTable(data, opts...)

	return output.WriteRenderedFrom(w, t.Render, "table", "render table")
}

// Render renders a Table as a styled terminal table string.
func Render(data *output.Table, opts ...Option) (string, error) {
	return output.RenderFromWrite(func(w io.Writer) error {
		return Write(w, data, opts...)
	})
}
