package d2_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderGraph() {
	b := output.NewGraphBuilder()
	b.AddNode(*output.NewGraphNode("design", "Design"))
	b.AddNode(*output.NewGraphNode("implement", "Implement"))
	b.AddEdge(*output.NewGraphEdge("design", "implement"))

	g := b.Build()

	out, err := d2.RenderGraph(g)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(out)
}
