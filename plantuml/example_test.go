package plantuml_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/testhelpers"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewPlantUMLDiagram() {
	diagram := plantuml.NewPlantUMLDiagram()
	diagram.AddNode(graphtest.NewTestNode("client", "Client"))
	diagram.AddNode(graphtest.NewTestNode("server", "Server"))
	diagram.AddEdge(graphtest.NewTestEdge("client", "server", "requests"))

	fmt.Println(testhelpers.MustRender(diagram))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewPlantUMLFromTable() {
	data := output.NewTable([]string{"Name"})
	data.AddRow([]string{"Alpha"})
	data.AddRow([]string{"Beta"})

	fmt.Println(testhelpers.MustRender(plantuml.NewPlantUMLFromTable(data)))
}
