package table

import (
	"fmt"
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

	out, err := t.Render()
	if err != nil {
		return fmt.Errorf("render table: %w", err)
	}

	return output.WriteRendered(w, "table", out)
}

// Render renders a Table as a styled terminal table string.
func Render(data *output.Table, opts ...Option) (string, error) {
	return output.RenderFromWrite(func(w io.Writer) error {
		return Write(w, data, opts...)
	})
}
