package d2

import (
	"io"

	"github.com/larsartmann/go-output"
)

// Option configures D2 rendering.
type Option func(*Config)

// Config holds D2 rendering configuration.
type Config struct {
	direction string
	layout    string
	title     string
}

// WithDirection sets the D2 diagram direction.
func WithDirection(dir Direction) Option {
	return func(c *Config) { c.direction = string(dir) }
}

// WithLayout sets the D2 layout engine.
func WithLayout(engine string) Option {
	return func(c *Config) { c.layout = engine }
}

// WithTitle sets the D2 diagram title.
func WithTitle(title string) Option {
	return func(c *Config) { c.title = title }
}

// Write writes a D2 diagram to the provided writer. Options are applied to
// a shallow copy — the caller's diagram is never mutated, so rendering the
// same diagram with different options (or no options) cannot leak settings
// across calls.
func Write(w io.Writer, diagram *Diagram, opts ...Option) error {
	cfg := Config{}
	for _, opt := range opts {
		opt(&cfg)
	}

	rendered := *diagram

	if cfg.direction != "" {
		rendered.SetDirection(Direction(cfg.direction))
	}

	if cfg.layout != "" {
		rendered.SetLayout(cfg.layout)
	}

	if cfg.title != "" {
		rendered.SetTitle(cfg.title)
	}

	return output.WriteRenderedRawFrom(w, rendered.Render, "d2", "render d2")
}

// Render renders a D2 diagram as a string.
func Render(diagram *Diagram, opts ...Option) (string, error) {
	return output.RenderFromWrite(func(w io.Writer) error {
		return Write(w, diagram, opts...)
	})
}

// WriteGraph writes a generic Graph as a D2 diagram.
func WriteGraph(w io.Writer, g output.Graph, opts ...Option) error {
	diagram := NewDiagram()
	diagram.SetNodes(g.Nodes())
	diagram.SetEdges(g.Edges())

	return Write(w, diagram, opts...)
}

// RenderGraph renders a generic Graph as a D2 string.
func RenderGraph(g output.Graph, opts ...Option) (string, error) {
	return output.RenderFromWrite(func(w io.Writer) error {
		return WriteGraph(w, g, opts...)
	})
}
