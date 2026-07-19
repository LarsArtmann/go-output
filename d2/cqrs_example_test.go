package d2_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderGraph() {
	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("design", "Design")).
		AddNode(*output.NewGraphNode("implement", "Implement")).
		AddEdge(*output.NewGraphEdge("design", "implement")).
		Build()

	out, err := d2.RenderGraph(g)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(out)
}
