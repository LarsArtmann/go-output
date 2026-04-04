package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output/internal/escape"
)

// StreamingRenderer is an interface for renderers that support streaming output.
// This is useful for rendering large datasets without loading everything into memory.
//
// Note: Not all implementations provide true streaming. The adapter returned by
// StreamingRendererFromRenderer collects output before writing. Only
// StreamingHTMLRenderer provides genuine streaming behavior.
type StreamingRenderer interface {
	Renderer
	// Stream writes the rendered output to an io.Writer in chunks.
	Stream(w io.Writer) error
}

// StreamingHTMLRenderer is a streaming implementation of HTMLRenderer.
// It writes output incrementally to minimize memory usage.
type StreamingHTMLRenderer struct {
	tableDataBase
}

// NewStreamingHTMLRenderer creates a new StreamingHTMLRenderer.
func NewStreamingHTMLRenderer() *StreamingHTMLRenderer {
	return &StreamingHTMLRenderer{}
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
		return r.writeEmptyTable(w)
	}

	err := r.writeTableOpen(w)
	if err != nil {
		return err
	}

	err = r.writeHeaders(w)
	if err != nil {
		return err
	}

	err = r.writeTableBodyOpen(w)
	if err != nil {
		return err
	}

	err = r.writeRows(w)
	if err != nil {
		return err
	}

	return r.writeTableClose(w)
}

func (r *StreamingHTMLRenderer) writeEmptyTable(w io.Writer) error {
	_, err := w.Write([]byte(`<table class="data-table"></table>`))
	if err != nil {
		return fmt.Errorf("write empty table: %w", err)
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeTableOpen(w io.Writer) error {
	_, err := w.Write([]byte(`<table class="data-table">
<thead>
<tr>
`))
	if err != nil {
		return fmt.Errorf("write table header: %w", err)
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeHeaders(w io.Writer) error {
	for _, h := range r.data.Headers {
		if _, err := w.Write([]byte("<th>" + escape.HTML(h) + "</th>\n")); err != nil {
			return fmt.Errorf("write header cell: %w", err)
		}
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeTableBodyOpen(w io.Writer) error {
	_, err := w.Write([]byte(`</tr>
</thead>
<tbody>
`))
	if err != nil {
		return fmt.Errorf("write table body: %w", err)
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeRows(w io.Writer) error {
	for i, row := range r.data.Rows {
		err := r.writeRow(w, row, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeRow(w io.Writer, row []string, rowIndex int) error {
	if _, err := w.Write([]byte("<tr>\n")); err != nil {
		return fmt.Errorf("write row %d start: %w", rowIndex, err)
	}

	for colIndex, cell := range row {
		if _, err := w.Write([]byte("<td>" + escape.HTML(cell) + "</td>\n")); err != nil {
			return fmt.Errorf("write row %d cell %d: %w", rowIndex, colIndex, err)
		}
	}

	if _, err := w.Write([]byte("</tr>\n")); err != nil {
		return fmt.Errorf("write row %d end: %w", rowIndex, err)
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeTableClose(w io.Writer) error {
	_, err := w.Write([]byte(`</tbody>
</table>
`))
	if err != nil {
		return fmt.Errorf("write table end: %w", err)
	}

	return nil
}

// StreamingRendererFromRenderer wraps a standard Renderer to implement StreamingRenderer.
// Note: This adapter does not provide true streaming behavior - it collects all output
// via Render() and then writes it at once. It exists to satisfy the StreamingRenderer
// interface for renderers that don't have native streaming support.
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
	if err != nil {
		return fmt.Errorf("stream render output: %w", err)
	}

	return nil
}
