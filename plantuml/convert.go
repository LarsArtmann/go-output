package plantuml

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

//nolint:gochecknoinits // Registers PlantUML TableRenderer for registry-based dispatch.
func init() {
	output.RegisterTableMarshaler(output.FormatPlantUML, renderPlantUMLTable)
}

func renderPlantUMLTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	out, err := NewPlantUMLFromTable(data).Render()
	if err != nil {
		return fmt.Errorf("render PlantUML: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write PlantUML output: %w", err)
	}

	return nil
}

// NewPlantUMLFromTable creates a PlantUML diagram from table data.
func NewPlantUMLFromTable(data *output.Table) *PlantUMLDiagram {
	diagram := NewPlantUMLDiagram()
	if data == nil {
		return diagram
	}

	diagram.SetNodesFromTable(data, func(_ int, _ *output.GraphNode) {})

	return diagram
}

// NewPlantUMLFromTree creates a PlantUML diagram from a tree hierarchy.
func NewPlantUMLFromTree(root *output.TreeNode) *PlantUMLDiagram {
	return output.TreeToRenderer(NewPlantUMLDiagram, (*PlantUMLDiagram).addTreeNodes, root)
}

func (d *PlantUMLDiagram) addTreeNodes(node *output.TreeNode, parentID output.TreeNodeID) {
	output.AddTreeNodes(&d.GraphRendererState, node, parentID.Get(), plantUMLTreeNodeID, "")
}

func plantUMLTreeNodeID(node *output.TreeNode) string {
	if !node.ID.IsZero() {
		return node.ID.Get()
	}

	return sanitizePlantUMLID(node.Label.Get())
}
