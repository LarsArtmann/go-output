package markup

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*AsciiDocTableRenderer)(nil)
	_ output.TableRenderer = (*AsciiDocTableRenderer)(nil)
)

//nolint:gochecknoinits // Registers AsciiDoc TableData marshaler for registry-based dispatch.
func init() {
	output.RegisterTableDataMarshaler(output.FormatAsciiDoc, renderAsciiDocTableData)
}

// AsciiDocTableRenderer renders TableData as an AsciiDoc table.
type AsciiDocTableRenderer struct {
	output.TableDataBase
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

	b.WriteString("|===")

	return b.String(), nil
}

// MarshalAsciiDocFromTableData marshals TableData as an AsciiDoc table.
func MarshalAsciiDocFromTableData(data *output.TableData) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	renderer := NewAsciiDocTableRenderer()
	renderer.SetData(data)

	out, err := renderer.Render()
	if err != nil {
		return nil, fmt.Errorf("render asciidoc table: %w", err)
	}

	return []byte(out), nil
}

func renderAsciiDocTableData(w io.Writer, data *output.TableData, _ output.RenderOptions) error {
	return renderMarshalAndWrite(w, data, "asciidoc", MarshalAsciiDocFromTableData)
}

// escapeAsciiDoc escapes special AsciiDoc characters in cell content.
func escapeAsciiDoc(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}

func writeAsciiDocCells(b *strings.Builder, cells []string) {
	for _, cell := range cells {
		b.WriteString("| ")
		b.WriteString(escapeAsciiDoc(cell))
		b.WriteString(" ")
	}
}
