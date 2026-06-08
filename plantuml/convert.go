package plantuml

import (
	"github.com/larsartmann/go-output"
)

// PlantUMLFromTableData creates a PlantUML diagram from table data.
func PlantUMLFromTableData(data *output.TableData) *PlantUMLDiagram {
	diagram := NewPlantUMLDiagram()
	if data == nil {
		return diagram
	}

	diagram.SetNodesFromTableData(data, func(_ int, _ *output.GraphNode) {})

	return diagram
}

// PlantUMLFromTree creates a PlantUML diagram from a tree hierarchy.
func PlantUMLFromTree(root *output.TreeNode) *PlantUMLDiagram {
	diagram := NewPlantUMLDiagram()
	if root == nil {
		return diagram
	}

	diagram.addTreeNodes(root, output.NewBrandedID[output.TreeNodeIDBrand](""))

	return diagram
}

func (d *PlantUMLDiagram) addTreeNodes(node *output.TreeNode, parentID output.TreeNodeID) {
	output.AddTreeNodes(&d.GraphRendererMixin, node, parentID.Get(), plantUMLTreeNodeID, "")
}

func plantUMLTreeNodeID(node *output.TreeNode) string {
	if !node.ID.IsZero() {
		return node.ID.Get()
	}

	return sanitizePlantUMLID(node.Label.Get())
}
