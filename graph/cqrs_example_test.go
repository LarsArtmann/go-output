package graph_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/graph"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderDOT() {
	b := output.NewGraphBuilder()
	b.AddNode(*output.NewGraphNode("compile", "Compile"))
	b.AddNode(*output.NewGraphNode("test", "Test"))
	b.AddEdge(*output.NewGraphEdge("compile", "test"))

	g := b.Build()

	dot, err := graph.RenderDOT(g)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(dot)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderMermaid() {
	b := output.NewGraphBuilder()
	b.AddNode(*output.NewGraphNode("build", "Build"))
	b.AddNode(*output.NewGraphNode("deploy", "Deploy"))
	b.AddEdge(*output.NewGraphEdge("build", "deploy"))

	g := b.Build()

	mermaid, err := graph.RenderMermaid(g)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(mermaid)
}
