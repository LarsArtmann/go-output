package plantuml

import (
	"fmt"
	"strings"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*PlantUMLDiagram)(nil)
	_ output.GraphRenderer = (*PlantUMLDiagram)(nil)
)

// PlantUMLDiagram renders a PlantUML component/class diagram.
type PlantUMLDiagram struct {
	output.GraphRendererMixin

	diagramType string
}

// NewPlantUMLDiagram creates a new PlantUMLDiagram.
func NewPlantUMLDiagram() *PlantUMLDiagram {
	return &PlantUMLDiagram{
		GraphRendererMixin: output.NewGraphRendererMixin(),
		diagramType:        "component",
	}
}

// AddNode adds a node to the diagram.
func (d *PlantUMLDiagram) AddNode(node output.GraphNode) *PlantUMLDiagram {
	*d.NodesPtr() = append(*d.NodesPtr(), node)
	return d
}

// AddEdge adds an edge to the diagram.
func (d *PlantUMLDiagram) AddEdge(edge output.GraphEdge) *PlantUMLDiagram {
	*d.EdgesPtr() = append(*d.EdgesPtr(), edge)
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
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")

	return s
}
