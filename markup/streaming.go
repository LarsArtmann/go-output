package markup

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
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
	output.TableDataBase
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
	data := r.Data()
	if data == nil {
		return r.writeEmptyTable(w)
	}

	err := r.writeTableOpen(w)
	if err != nil {
		return err
	}

	err = r.writeHeaders(w, data)
	if err != nil {
		return err
	}

	err = r.writeTableBodyOpen(w)
	if err != nil {
		return err
	}

	err = r.writeRows(w, data)
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

func (r *StreamingHTMLRenderer) writeHeaders(w io.Writer, data *output.TableData) error {
	for _, h := range data.Headers {
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

func (r *StreamingHTMLRenderer) writeRows(w io.Writer, data *output.TableData) error {
	for i, row := range data.Rows {
		err := r.writeRow(w, row, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *StreamingHTMLRenderer) writeRow(w io.Writer, row []string, rowIndex int) error {
	if err := r.writeChunk(
		w, []byte("<tr>\n"),
		fmt.Sprintf("write row %d start", rowIndex),
	); err != nil {
		return err
	}

	for colIndex, cell := range row {
		if _, err := w.Write([]byte("<td>" + escape.HTML(cell) + "</td>\n")); err != nil {
			return fmt.Errorf("write row %d cell %d: %w", rowIndex, colIndex, err)
		}
	}

	return r.writeChunk(
		w, []byte("</tr>\n"),
		fmt.Sprintf("write row %d end", rowIndex),
	)
}

func (r *StreamingHTMLRenderer) writeTableClose(w io.Writer) error {
	return r.writeChunk(w, []byte(`</tbody>
</table>
`), "write table end")
}
