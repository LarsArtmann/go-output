package graph_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/graph"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderDOT() {
	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("compile", "Compile")).
		AddNode(*output.NewGraphNode("test", "Test")).
		AddEdge(*output.NewGraphEdge("compile", "test")).
		Build()

	dot, err := graph.RenderDOT(g)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(dot)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderMermaid() {
	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("build", "Build")).
		AddNode(*output.NewGraphNode("deploy", "Deploy")).
		AddEdge(*output.NewGraphEdge("build", "deploy")).
		Build()

	mermaid, err := graph.RenderMermaid(g)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(mermaid)
}
