package plantuml_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/plantuml"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewPlantUMLDiagram() {
	diagram := plantuml.NewPlantUMLDiagram()
	diagram.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("client"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Client"),
	})
	diagram.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("server"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Server"),
	})
	diagram.AddEdge(output.GraphEdge{
		From:  output.NewBrandedID[output.GraphNodeIDBrand]("client"),
		To:    output.NewBrandedID[output.GraphNodeIDBrand]("server"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("requests"),
	})

	result, err := diagram.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExamplePlantUMLFromTableData() {
	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Alpha"})
	data.AddRow([]string{"Beta"})

	result, err := plantuml.PlantUMLFromTableData(data).Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}
