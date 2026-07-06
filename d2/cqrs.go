package d2

import (
	"fmt"
	"io"
	"strings"

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

// Write writes a D2 diagram to the provided writer.
func Write(w io.Writer, diagram *Diagram, opts ...Option) error {
	cfg := Config{}
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.direction != "" {
		diagram.SetDirection(Direction(cfg.direction))
	}

	if cfg.layout != "" {
		diagram.SetLayout(cfg.layout)
	}

	if cfg.title != "" {
		diagram.SetTitle(cfg.title)
	}

	out, err := diagram.Render()
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, out)
	if err != nil {
		return fmt.Errorf("write d2 output: %w", err)
	}

	return nil
}

// Render renders a D2 diagram as a string.
func Render(diagram *Diagram, opts ...Option) (string, error) {
	var buf strings.Builder
	if err := Write(&buf, diagram, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
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
	var buf strings.Builder
	if err := WriteGraph(&buf, g, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}
