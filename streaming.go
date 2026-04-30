package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output/escape"
)

// Compile-time interface checks.
var (
	_ Renderer          = (*StreamingHTMLRenderer)(nil)
	_ TableRenderer     = (*StreamingHTMLRenderer)(nil)
	_ StreamingRenderer = (*StreamingHTMLRenderer)(nil)
	_ StreamingRenderer = (*adapterRenderer)(nil)
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

func (r *StreamingHTMLRenderer) writeChunk(w io.Writer, chunk []byte, description string) error {
	_, err := w.Write(chunk)
	if err != nil {
		return fmt.Errorf("%s: %w", description, err)
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeEmptyTable(w io.Writer) error {
	return r.writeChunk(w, []byte(`<table class="data-table"></table>`), "write empty table")
}

func (r *StreamingHTMLRenderer) writeTableOpen(w io.Writer) error {
	return r.writeChunk(w, []byte(`<table class="data-table">
<thead>
<tr>
`), "write table header")
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
	return r.writeChunk(w, []byte(`</tr>
</thead>
<tbody>
`), "write table body")
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
	if err := r.writeChunkWithError(w, []byte("<tr>\n"), rowIndex, "start"); err != nil {
		return err
	}

	for colIndex, cell := range row {
		if _, err := w.Write([]byte("<td>" + escape.HTML(cell) + "</td>\n")); err != nil {
			return fmt.Errorf("write row %d cell %d: %w", rowIndex, colIndex, err)
		}
	}

	return r.writeChunkWithError(w, []byte("</tr>\n"), rowIndex, "end")
}

func (r *StreamingHTMLRenderer) writeChunkWithError(
	w io.Writer,
	chunk []byte,
	rowIndex int,
	location string,
) error {
	_, err := w.Write(chunk)
	if err != nil {
		return fmt.Errorf("write row %d %s: %w", rowIndex, location, err)
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeTableClose(w io.Writer) error {
	return r.writeChunk(w, []byte(`</tbody>
</table>
`), "write table end")
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

func (a *adapterRenderer) Render() (string, error) {
	out, err := a.r.Render()
	if err != nil {
		return "", fmt.Errorf("adapter render: %w", err)
	}

	return out, nil
}

func (a *adapterRenderer) Stream(w io.Writer) error {
	out, err := a.r.Render()
	if err != nil {
		return fmt.Errorf("render for streaming: %w", err)
	}

	_, err = w.Write([]byte(out))
	if err != nil {
		return fmt.Errorf("stream render output: %w", err)
	}

	return nil
}
