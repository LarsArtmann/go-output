package markup

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

func renderMarshalAndWrite(
	w io.Writer,
	data *output.Table,
	formatName string,
	marshalFunc func(*output.Table) ([]byte, error),
) error {
	b, err := marshalFunc(data)
	if err != nil {
		return fmt.Errorf("marshal %s table data: %w", formatName, err)
	}

	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write %s output: %w", formatName, err)
	}

	return nil
}

func writeRowTag(w io.Writer, indent, tag string, isClose bool) error {
	var content string
	if isClose {
		content = indent + "</" + tag + ">\n"
	} else {
		content = indent + "<" + tag + ">\n"
	}

	if _, err := io.WriteString(w, content); err != nil {
		return fmt.Errorf(
			"write %s row tag %q at indent %q content=%q: %w",
			map[bool]string{true: "close", false: "open"}[isClose],
			tag,
			indent,
			content,
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
		return fmt.Errorf("open row %q: %w", rowTag, err)
	}

	for _, cell := range row {
		if _, err := io.WriteString(w, indent+indent+"<"+cellTag+">"); err != nil {
			return fmt.Errorf("write cell tag open %q at indent %q: %w", cellTag, indent, err)
		}

		if _, err := io.WriteString(w, escapeFn(cell)); err != nil {
			return fmt.Errorf("write cell content in row %q: %w", rowTag, err)
		}

		if _, err := io.WriteString(w, "</"+cellTag+">\n"); err != nil {
			return fmt.Errorf("write cell tag close %q at indent %q: %w", cellTag, indent, err)
		}
	}

	if err := writeRowTag(w, indent, rowTag, true); err != nil {
		return fmt.Errorf("close row %q: %w", rowTag, err)
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
			return fmt.Errorf("write column tag open at indent %q: %w", indent, err)
		}

		if _, err := io.WriteString(w, escapeFn(col)); err != nil {
			return fmt.Errorf("write column content %q: %w", col, err)
		}

		if _, err := io.WriteString(w, "</column>\n"); err != nil {
			return fmt.Errorf("write column tag close at indent %q: %w", indent, err)
		}
	}

	return nil
}
