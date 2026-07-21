package plantuml

import (
	"io"

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

	return output.WriteRenderedRawFrom(w, d.Render, "plantuml", "plantuml")
}

// Render renders a Graph as a PlantUML string.
func Render(g output.Graph, opts ...Option) (string, error) {
	return output.RenderFromWrite(func(w io.Writer) error {
		return Write(w, g, opts...)
	})
}
