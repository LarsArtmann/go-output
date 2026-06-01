package plantuml_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewPlantUMLDiagram() {
	diagram := plantuml.NewPlantUMLDiagram()
	diagram.AddNode(graphtest.NewTestNode("client", "Client"))
	diagram.AddNode(graphtest.NewTestNode("server", "Server"))
	diagram.AddEdge(graphtest.NewTestEdge("client", "server", "requests"))

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
