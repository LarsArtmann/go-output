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

	if data.HasFooter() {
		err = r.writeFooter(w, data)
	} else {
		err = r.writeChunk(w, []byte("</tbody>\n"), "close tbody")
	}

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
	if err := r.writeRowBoundary(w, "<tr>\n", rowIndex, "start"); err != nil {
		return err
	}

	for colIndex, cell := range row {
		if _, err := w.Write([]byte("<td>" + escape.HTML(cell) + "</td>\n")); err != nil {
			return fmt.Errorf("write row %d cell %d: %w", rowIndex, colIndex, err)
		}
	}

	return r.writeRowBoundary(w, "</tr>\n", rowIndex, "end")
}

func (r *StreamingHTMLRenderer) writeRowBoundary(w io.Writer, tag string, rowIndex int, phase string) error {
	return r.writeChunk(w, []byte(tag), fmt.Sprintf("write row %d %s", rowIndex, phase))
}

func (r *StreamingHTMLRenderer) writeTableClose(w io.Writer) error {
	return r.writeChunk(w, []byte("</table>\n"), "write table end")
}

func (r *StreamingHTMLRenderer) writeFooter(w io.Writer, data *output.TableData) error {
	if err := r.writeChunk(w, []byte("</tbody>\n"), "close tbody before tfoot"); err != nil {
		return err
	}

	if err := r.writeChunk(w, []byte("<tfoot>\n<tr>\n"), "write tfoot open"); err != nil {
		return err
	}

	for i, cell := range data.Footer {
		if _, err := w.Write([]byte("<td>" + escape.HTML(cell) + "</td>\n")); err != nil {
			return fmt.Errorf("write footer cell %d: %w", i, err)
		}
	}

	if err := r.writeChunk(w, []byte("</tr>\n</tfoot>\n"), "write tfoot close"); err != nil {
		return err
	}

	return nil
}
