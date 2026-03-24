package output

import (
	"io"
	"strings"
)

// StreamingRenderer is an interface for renderers that support streaming output.
// This is useful for rendering large datasets without loading everything into memory.
type StreamingRenderer interface {
	Renderer
	// Stream writes the rendered output to an io.Writer in chunks.
	Stream(w io.Writer) error
}

// StreamingHTMLRenderer is a streaming implementation of HTMLRenderer.
// It writes output incrementally to minimize memory usage.
type StreamingHTMLRenderer struct {
	data *TableData
}

// NewStreamingHTMLRenderer creates a new StreamingHTMLRenderer.
func NewStreamingHTMLRenderer() *StreamingHTMLRenderer {
	return &StreamingHTMLRenderer{}
}

// SetHeaders sets the column headers.
func (r *StreamingHTMLRenderer) SetHeaders(headers []string) {
	if r.data == nil {
		r.data = &TableData{}
	}
	r.data.Headers = headers
}

// AddRow adds a data row.
func (r *StreamingHTMLRenderer) AddRow(row []string) {
	if r.data == nil {
		r.data = &TableData{}
	}
	r.data.Rows = append(r.data.Rows, row)
}

// SetData sets the table data directly.
func (r *StreamingHTMLRenderer) SetData(data *TableData) {
	r.data = data
}

// Render returns the HTML table as a string.
func (r *StreamingHTMLRenderer) Render() string {
	var b strings.Builder
	_ = r.Stream(&b)
	return b.String()
}

// Stream writes the HTML table incrementally to an io.Writer.
func (r *StreamingHTMLRenderer) Stream(w io.Writer) error {
	if r.data == nil {
		_, err := w.Write([]byte(`<table class="data-table"></table>`))
		return err
	}

	// Write table opening
	if _, err := w.Write([]byte(`<table class="data-table">
<thead>
<tr>
`)); err != nil {
		return err
	}

	// Write headers
	for _, h := range r.data.Headers {
		if _, err := w.Write([]byte("<th>" + escapeHTML(h) + "</th>\n")); err != nil {
			return err
		}
	}

	// Write table body opening
	if _, err := w.Write([]byte(`</tr>
</thead>
<tbody>
`)); err != nil {
		return err
	}

	// Write rows
	for _, row := range r.data.Rows {
		if _, err := w.Write([]byte("<tr>\n")); err != nil {
			return err
		}
		for _, cell := range row {
			if _, err := w.Write([]byte("<td>" + escapeHTML(cell) + "</td>\n")); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte("</tr>\n")); err != nil {
			return err
		}
	}

	// Write table closing
	_, err := w.Write([]byte(`</tbody>
</table>
`))
	return err
}

// escapeHTML escapes HTML special characters.
func escapeHTML(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '&':
			b.WriteString("&amp;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&#39;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// StreamingRendererFromRenderer wraps a standard Renderer to implement StreamingRenderer.
// The Stream method collects all output and writes it at once.
func StreamingRendererFromRenderer(r Renderer) StreamingRenderer {
	return &adapterRenderer{r: r}
}

type adapterRenderer struct {
	r Renderer
}

func (a *adapterRenderer) Render() string {
	return a.r.Render()
}

func (a *adapterRenderer) Stream(w io.Writer) error {
	_, err := w.Write([]byte(a.r.Render()))
	return err
}
