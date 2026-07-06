package output

import (
	"fmt"
	"io"
)

// StreamingRenderer is an interface for renderers that support streaming output.
// This is useful for rendering large datasets without loading everything into memory.
//
// Important: Not all implementations provide true streaming. The adapter returned by
// RendererAsWriter collects output before writing. Only
// StreamingHTMLRenderer (in the markup sub-module) provides genuine streaming behavior.
type StreamingRenderer interface {
	Renderer
	// Stream writes the rendered output to an io.Writer in chunks.
	Stream(w io.Writer) error
}

// RendererAsWriter wraps a standard Renderer to implement StreamingRenderer.
// The name is honest: it renders the full output via Render() then writes it
// to the writer in one shot — NOT true streaming. For genuine streaming, use
// StreamingHTMLRenderer (markup sub-module).
func RendererAsWriter(r Renderer) StreamingRenderer {
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
