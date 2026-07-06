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
	output.GraphBuilder

	diagramType string
}

// NewPlantUMLDiagram creates a new PlantUMLDiagram.
func NewPlantUMLDiagram() *PlantUMLDiagram {
	return &PlantUMLDiagram{
		GraphBuilder: output.NewGraphBuilder(),
		diagramType:  "component",
	}
}

// AddNode adds a node to the diagram.
func (d *PlantUMLDiagram) AddNode(node output.GraphNode) *PlantUMLDiagram {
	d.GraphBuilder.AddNode(node)
	return d
}

// AddEdge adds an edge to the diagram.
func (d *PlantUMLDiagram) AddEdge(edge output.GraphEdge) *PlantUMLDiagram {
	d.GraphBuilder.AddEdge(edge)
	return d
}

// Render returns the PlantUML diagram as a string.
func (d *PlantUMLDiagram) Render() (string, error) {
	var b strings.Builder

	b.WriteString("@startuml\n")
	b.WriteString("skinparam componentStyle uml2\n")
	b.WriteString("skinparam defaultFontSize 12\n\n")

	for _, node := range d.Nodes() {
		colorSpec := plantumlColorSpec(node.Style)

		label := escape.PlantUML(node.Label.Get())
		if colorSpec != "" {
			fmt.Fprintf(&b, "[%s] as %s %s\n", label, sanitizePlantUMLID(node.ID.Get()), colorSpec)
		} else {
			fmt.Fprintf(&b, "[%s] as %s\n", label, sanitizePlantUMLID(node.ID.Get()))
		}
	}

	b.WriteString("\n")

	for _, edge := range d.Edges() {
		label := ""
		if !edge.Label.IsZero() {
			label = " : " + escape.PlantUML(edge.Label.Get())
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

// plantumlColorSpec converts a NodeStyle into a PlantUML color specification
// string for per-element styling. Returns empty string when no colors are set.
//
// PlantUML syntax: the spec starts with '#' and joins attributes with ';'.
// Example: #e8a838;line:#4a4030;text:#14110d.
func plantumlColorSpec(s output.NodeStyle) string {
	var parts []string

	if s.Fill != "" {
		parts = append(parts, plantumlColorValue(s.Fill))
	}

	if s.Stroke != "" {
		parts = append(parts, "line:"+plantumlColorValue(s.Stroke))
	}

	if s.FontColor != "" {
		parts = append(parts, "text:"+plantumlColorValue(s.FontColor))
	}

	result := strings.Join(parts, ";")

	if result != "" && !strings.HasPrefix(result, "#") {
		result = "#" + result
	}

	return result
}

// plantumlColorValue escapes a color value for use in a PlantUML color spec.
// Uses escape.PlantUML for general escaping (newline, backslash, quote, ])
// and additionally replaces ';' — the PlantUML attribute separator — to
// prevent attribute injection.
func plantumlColorValue(s string) string {
	return strings.ReplaceAll(escape.PlantUML(s), ";", "_")
}
