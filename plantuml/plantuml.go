package plantuml

import (
	"fmt"
	"strings"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

var (
	_ output.Renderer      = (*PlantUMLDiagram)(nil)
	_ output.GraphRenderer = (*PlantUMLDiagram)(nil)
)

//nolint:gochecknoinits // Registers PlantUML format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatPlantUML, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
}

// PlantUMLDiagram renders a PlantUML component/class diagram.
type PlantUMLDiagram struct {
	output.GraphRendererState

	diagramType string
}

// NewPlantUMLDiagram creates a new PlantUMLDiagram.
func NewPlantUMLDiagram() *PlantUMLDiagram {
	return &PlantUMLDiagram{
		GraphRendererState: output.NewGraphRendererState(),
		diagramType:        "component",
	}
}

// AddNode adds a node to the diagram.
func (d *PlantUMLDiagram) AddNode(node output.GraphNode) *PlantUMLDiagram {
	d.GraphRendererState.AddNode(node)
	return d
}

// AddEdge adds an edge to the diagram.
func (d *PlantUMLDiagram) AddEdge(edge output.GraphEdge) *PlantUMLDiagram {
	d.GraphRendererState.AddEdge(edge)
	return d
}

// Render returns the PlantUML diagram as a string.
func (d *PlantUMLDiagram) Render() (string, error) {
	var b strings.Builder

	b.WriteString("@startuml\n")
	b.WriteString("skinparam componentStyle uml2\n")
	b.WriteString("skinparam defaultFontSize 12\n\n")

	for _, node := range d.Nodes() {
		fmt.Fprintf(&b, "[%s] as %s\n", node.Label.Get(), sanitizePlantUMLID(node.ID.Get()))
	}

	b.WriteString("\n")

	for _, edge := range d.Edges() {
		label := ""
		if !edge.Label.IsZero() {
			label = " : " + edge.Label.Get()
		}

		fmt.Fprintf(
			&b, "%s -->%s %s\n",
			sanitizePlantUMLID(edge.From.Get()),
			label,
			sanitizePlantUMLID(edge.To.Get()),
		)
	}

	b.WriteString("\n@enduml")

	return b.String(), nil
}

// sanitizePlantUMLID converts a string to a valid PlantUML identifier.
func sanitizePlantUMLID(s string) string {
	return escape.SlugifyID(s)
}
