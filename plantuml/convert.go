package plantuml

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

//nolint:gochecknoinits // Registers PlantUML TableDataRenderer for registry-based dispatch.
func init() {
	output.RegisterTableDataRenderer(output.FormatPlantUML, renderPlantUMLTableData)
}

func renderPlantUMLTableData(w io.Writer, data *output.TableData, _ output.RenderOptions) error {
	out, err := PlantUMLFromTableData(data).Render()
	if err != nil {
		return fmt.Errorf("render PlantUML: %w", err)
	}

	_, err = fmt.Fprintln(w, out)
	if err != nil {
		return fmt.Errorf("write PlantUML output: %w", err)
	}

	return nil
}

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
	output.AddTreeNodes(&d.GraphRendererState, node, parentID.Get(), plantUMLTreeNodeID, "")
}

func plantUMLTreeNodeID(node *output.TreeNode) string {
	if !node.ID.IsZero() {
		return node.ID.Get()
	}

	return sanitizePlantUMLID(node.Label.Get())
}
