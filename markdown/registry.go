package markdown

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

//nolint:gochecknoinits // Self-registers the Markdown TableData renderer, mirroring table/.
func init() {
	output.RegisterTableDataRenderer(output.FormatMarkdown, renderMarkdownTableData)
}

// renderMarkdownTableData renders TableData as a Markdown table to w.
// It honors RenderOptions.Title (document title + row count header) and ColorMode.
func renderMarkdownTableData(w io.Writer, data *output.TableData, opts output.RenderOptions) error {
	if opts.Title != "" {
		_, err := fmt.Fprintf(w, "# %s\n\n", opts.Title)
		if err != nil {
			return fmt.Errorf("write markdown title: %w", err)
		}

		_, err = fmt.Fprintf(w, "%d rows\n\n", data.RowCount())
		if err != nil {
			return fmt.Errorf("write markdown row count: %w", err)
		}
	}

	mdTable := NewMarkdownTableFromData(data)
	mdTable.SetColorMode(opts.ColorMode)

	out, err := mdTable.Render()
	if err != nil {
		return fmt.Errorf("render markdown: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write markdown output: %w", err)
	}

	return nil
}
