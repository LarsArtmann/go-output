package plantuml

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// Option configures PlantUML rendering.
type Option func(*Config)

// Config holds PlantUML rendering configuration.
type Config struct {
	diagramType string
}

// WithDiagramType sets the PlantUML diagram type (default "component").
func WithDiagramType(t string) Option {
	return func(c *Config) { c.diagramType = t }
}

// Write writes a Graph as PlantUML to the provided writer.
func Write(w io.Writer, g output.Graph, opts ...Option) error {
	cfg := Config{diagramType: "component"}
	for _, opt := range opts {
		opt(&cfg)
	}

	d := NewPlantUMLDiagram()
	d.diagramType = cfg.diagramType
	d.SetNodes(g.Nodes())
	d.SetEdges(g.Edges())

	out, err := d.Render()
	if err != nil {
		return fmt.Errorf("plantuml: %w", err)
	}

	return output.WriteRenderedRaw(w, "plantuml", out)
}

// Render renders a Graph as a PlantUML string.
func Render(g output.Graph, opts ...Option) (string, error) {
	var buf strings.Builder
	if err := Write(&buf, g, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}
