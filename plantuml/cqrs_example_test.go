package plantuml_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/plantuml"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRender() {
	b := output.NewGraphBuilder()
	b.AddNode(*output.NewGraphNode("client", "Client"))
	b.AddNode(*output.NewGraphNode("server", "Server"))
	b.AddEdge(*output.NewGraphEdge("client", "server"))

	g := b.Build()

	out, err := plantuml.Render(g)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(out)
}
