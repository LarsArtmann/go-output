package tree

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

//nolint:gochecknoinits // Self-registers the Tree Table renderer, mirroring markdown/ and table/.
func init() {
	output.RegisterTableMarshaler(output.FormatTree, renderTreeTable)
}

// renderTreeTable renders Table as an ASCII tree to w.
func renderTreeTable(w io.Writer, data *output.Table, opts output.RenderOptions) error {
	renderer := TreeRendererFromTable(data)
	renderer.SetColorMode(opts.ColorMode)

	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render tree: %w", err)
	}

	return output.WriteRendered(w, "tree", out)
}
