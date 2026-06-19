package tree

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

//nolint:gochecknoinits // Self-registers the Tree TableData renderer, mirroring markdown/ and table/.
func init() {
	output.RegisterTableDataRenderer(output.FormatTree, renderTreeTableData)
}

// renderTreeTableData renders TableData as an ASCII tree to w.
func renderTreeTableData(w io.Writer, data *output.TableData, opts output.RenderOptions) error {
	renderer := TreeRendererFromTableData(data)
	renderer.SetColorMode(opts.ColorMode)

	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render tree: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write tree output: %w", err)
	}

	return nil
}
