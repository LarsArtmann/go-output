package output

import (
	"fmt"
	"io"
)

// StreamingRenderer is an interface for renderers that support streaming output.
// This is useful for rendering large datasets without loading everything into memory.
//
// Important: Not all implementations provide true streaming. The adapter returned by
// StreamingRendererFromRenderer collects output before writing. Only
// StreamingHTMLRenderer (in the markup sub-module) provides genuine streaming behavior.
type StreamingRenderer interface {
	Renderer
	// Stream writes the rendered output to an io.Writer in chunks.
	Stream(w io.Writer) error
}

// StreamingRendererFromRenderer wraps a standard Renderer to implement StreamingRenderer.
// Important: This adapter does not provide true streaming behavior - it collects all output
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
