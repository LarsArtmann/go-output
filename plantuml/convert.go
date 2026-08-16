package plantuml

import (
	"io"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/escape"
)

//nolint:gochecknoinits // Registers PlantUML TableRenderer for registry-based dispatch.
func init() {
	output.RegisterTableMarshaler(output.FormatPlantUML, renderPlantUMLTable)
}

func renderPlantUMLTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return output.WriteRenderedFrom(
		w,
		NewPlantUMLFromTable(data).Render,
		"PlantUML",
		"render PlantUML",
	)
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
	output.AddTreeNodes(&d.GraphBuilder, node, parentID.Get(), plantUMLTreeNodeID, "")
}

func plantUMLTreeNodeID(node *output.TreeNode) string {
	if !node.ID.IsZero() {
		return escape.PlantUMLID(node.ID.Get())
	}

	return escape.PlantUMLID(node.Label.Get())
}
