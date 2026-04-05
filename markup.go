package output

import (
	"strings"

	"github.com/larsartmann/go-output/internal/escape"
)

func writeMarkupRow(
	b *strings.Builder,
	row []string,
	rowTag, cellTag, indent string,
	escapeFn func(string) string,
) {
	b.WriteString(indent + "<" + rowTag + ">\n")

	for _, cell := range row {
		b.WriteString(indent + indent + "<" + cellTag + ">")
		b.WriteString(escapeFn(cell))
		b.WriteString("</" + cellTag + ">\n")
	}

	b.WriteString(indent + "</" + rowTag + ">\n")
}

func writeMarkupColumns(
	b *strings.Builder,
	cols []string,
	indent string,
	escapeFn func(string) string,
) {
	for _, col := range cols {
		b.WriteString(indent + "<column>")
		b.WriteString(escapeFn(col))
		b.WriteString("</column>\n")
	}
}

func htmlEscape(s string) string {
	return escape.HTML(s)
}

func xmlEscape(s string) string {
	return escape.XML(s)
}
