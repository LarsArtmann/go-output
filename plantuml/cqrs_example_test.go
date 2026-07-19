package plantuml_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/plantuml"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRender() {
	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("client", "Client")).
		AddNode(*output.NewGraphNode("server", "Server")).
		AddEdge(*output.NewGraphEdge("client", "server")).
		Build()

	out, err := plantuml.Render(g)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(out)
}
