package markdown

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

//nolint:gochecknoinits // Self-registers the Markdown Table renderer, mirroring table/.
func init() {
	output.RegisterTableMarshaler(output.FormatMarkdown, renderMarkdownTable)
}

// renderMarkdownTable renders Table as a Markdown table to w.
// It honors RenderOptions.Title (document title + row count header) and ColorMode.
func renderMarkdownTable(w io.Writer, data *output.Table, opts output.RenderOptions) error {
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

	mdTable := NewMarkdownTableFromTable(data)
	mdTable.SetColorMode(opts.ColorMode)

	out, err := mdTable.Render()
	if err != nil {
		return fmt.Errorf("render markdown: %w", err)
	}

	return output.WriteRendered(w, "markdown", out)
}
