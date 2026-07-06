package table

import (
	"fmt"
	"io"
	"strings"

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

	if _, err := fmt.Fprintln(w, out); err != nil {
		return fmt.Errorf("write table output: %w", err)
	}

	return nil
}

// Render renders a Table as a styled terminal table string.
func Render(data *output.Table, opts ...Option) (string, error) {
	var buf strings.Builder
	if err := Write(&buf, data, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}
