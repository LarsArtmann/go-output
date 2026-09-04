package markup

import (
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*AsciiDocTableRenderer)(nil)
	_ output.TableRenderer = (*AsciiDocTableRenderer)(nil)
)

//nolint:gochecknoinits // Registers AsciiDoc Table marshaler and format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatAsciiDoc, output.ShapeTable)
	output.RegisterTableMarshaler(output.FormatAsciiDoc, renderAsciiDocTable)
}

// AsciiDocTableRenderer renders Table as an AsciiDoc table.
type AsciiDocTableRenderer struct {
	output.TableStore
}

// NewAsciiDocTableRenderer creates a new AsciiDocTableRenderer.
func NewAsciiDocTableRenderer() *AsciiDocTableRenderer {
	return &AsciiDocTableRenderer{}
}

// Render returns the table data as an AsciiDoc table string.
func (r *AsciiDocTableRenderer) Render() (string, error) {
	data := r.Data()
	if data == nil || len(data.Headers) == 0 {
		return "", nil
	}

	var b strings.Builder

	b.WriteString("|===\n")

	for _, h := range data.Headers {
		b.WriteString("| ")
		b.WriteString(escapeAsciiDoc(h))
		b.WriteString(" ")
	}

	b.WriteString("\n\n")

	for _, row := range data.Rows {
		writeAsciiDocCells(&b, row)
		b.WriteString("\n\n")
	}

	if data.HasFooter() {
		writeAsciiDocCells(&b, data.Footer)
		b.WriteString("\n\n")
	}

	b.WriteString("|===\n")

	return b.String(), nil
}

// MarshalAsciiDocFromTable marshals Table as an AsciiDoc table.
func MarshalAsciiDocFromTable(data *output.Table) ([]byte, error) {
	return output.MarshalViaRenderer(data, "asciidoc table", NewAsciiDocTableRenderer)
}

func renderAsciiDocTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return renderMarshalAndWrite(w, data, "asciidoc", MarshalAsciiDocFromTable)
}

// asciidocReplacer escapes AsciiDoc special characters: cell delimiter and inline formatting.
//
//nolint:gochecknoglobals // Reusable strings.Replacer, safe to share.
var asciidocReplacer = strings.NewReplacer(
	"|", "\\|",
	"*", "\\*",
	"_", "\\_",
	"`", "\\`",
	"~", "\\~",
	"^", "\\^",
)

// escapeAsciiDoc escapes special AsciiDoc characters in cell content.
func escapeAsciiDoc(s string) string {
	return asciidocReplacer.Replace(s)
}

func writeAsciiDocCells(b *strings.Builder, cells []string) {
	for _, cell := range cells {
		b.WriteString("| ")
		b.WriteString(escapeAsciiDoc(cell))
		b.WriteString(" ")
	}
}
