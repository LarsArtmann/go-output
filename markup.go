package output

import (
	"fmt"
	"io"
)

// writeRowTag writes a row opening or closing tag.
func writeRowTag(w io.Writer, indent, tag string, isClose bool) error {
	var content string
	if isClose {
		content = indent + "</" + tag + ">\n"
	} else {
		content = indent + "<" + tag + ">\n"
	}

	if _, err := io.WriteString(w, content); err != nil {
		return fmt.Errorf(
			"write row tag %s: %w",
			map[bool]string{true: "close", false: "open"}[isClose],
			err,
		)
	}

	return nil
}

func writeMarkupRow(
	w io.Writer,
	row []string,
	rowTag, cellTag, indent string,
	escapeFn func(string) string,
) error {
	if err := writeRowTag(w, indent, rowTag, false); err != nil {
		return err
	}

	for _, cell := range row {
		if _, err := io.WriteString(w, indent+indent+"<"+cellTag+">"); err != nil {
			return fmt.Errorf("write cell tag open: %w", err)
		}

		if _, err := io.WriteString(w, escapeFn(cell)); err != nil {
			return fmt.Errorf("write cell content: %w", err)
		}

		if _, err := io.WriteString(w, "</"+cellTag+">\n"); err != nil {
			return fmt.Errorf("write cell tag close: %w", err)
		}
	}

	if err := writeRowTag(w, indent, rowTag, true); err != nil {
		return err
	}

	return nil
}

func writeMarkupColumns(
	w io.Writer,
	cols []string,
	indent string,
	escapeFn func(string) string,
) error {
	for _, col := range cols {
		if _, err := io.WriteString(w, indent+"<column>"); err != nil {
			return fmt.Errorf("write column tag open: %w", err)
		}

		if _, err := io.WriteString(w, escapeFn(col)); err != nil {
			return fmt.Errorf("write column content: %w", err)
		}

		if _, err := io.WriteString(w, "</column>\n"); err != nil {
			return fmt.Errorf("write column tag close: %w", err)
		}
	}

	return nil
}
