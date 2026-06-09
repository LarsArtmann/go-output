package markup

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// Compile-time interface checks.
var (
	_ output.Renderer          = (*StreamingHTMLRenderer)(nil)
	_ output.TableRenderer     = (*StreamingHTMLRenderer)(nil)
	_ output.StreamingRenderer = (*StreamingHTMLRenderer)(nil)
)

// StreamingHTMLRenderer is a streaming implementation of HTMLRenderer.
// It writes output incrementally to minimize memory usage.
type StreamingHTMLRenderer struct {
	output.TableDataStore
}

// NewStreamingHTMLRenderer creates a new StreamingHTMLRenderer.
func NewStreamingHTMLRenderer() *StreamingHTMLRenderer {
	return &StreamingHTMLRenderer{}
}

// Render returns the HTML table as a string.
func (r *StreamingHTMLRenderer) Render() (string, error) {
	var b strings.Builder

	err := r.Stream(&b)
	if err != nil {
		return "", fmt.Errorf("stream to string: %w", err)
	}

	return b.String(), nil
}

// Stream writes the HTML table incrementally to an io.Writer.
func (r *StreamingHTMLRenderer) Stream(w io.Writer) error {
	return streamHTMLTable(w, r.Data())
}

// htmlTableData holds the data for HTML table template rendering.
type htmlTableData struct {
	Headers []string
	Rows    [][]string
	Footer  []string
}

//nolint:gochecknoglobals // Parsed template for HTML table.
var htmlTableTemplate = template.Must(template.New("htmlTable").Parse(
	`<table class="data-table">
{{- if .Headers}}<thead>
<tr>
{{range .Headers}}<th>{{.}}</th>
{{end}}</tr>
</thead>
{{- end}}<tbody>
{{range $i, $row := .Rows}}<tr>
{{range $row}}<td>{{.}}</td>
{{end}}</tr>
{{end}}{{if .Footer}}</tbody>
<tfoot>
<tr>
{{range .Footer}}<td class="footer-cell">{{.}}</td>
{{end}}</tr>
</tfoot>
{{else}}</tbody>
{{end}}</table>`,
))

// streamHTMLTable writes a complete HTML table to w using html/template for auto-escaping.
func streamHTMLTable(w io.Writer, data *output.TableData) error {
	if data == nil {
		if _, err := w.Write([]byte(`<table class="data-table"></table>`)); err != nil {
			return fmt.Errorf("write empty table: %w", err)
		}

		return nil
	}

	err := htmlTableTemplate.Execute(w, htmlTableData{
		Headers: data.Headers,
		Rows:    data.Rows,
		Footer:  data.Footer,
	})
	if err != nil {
		return fmt.Errorf("execute html table template: %w", err)
	}

	return nil
}
