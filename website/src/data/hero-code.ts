export const heroCode = `package main

import (
    "fmt"
    "github.com/larsartmann/go-output"
    "github.com/larsartmann/go-output/graph"
)

func main() {
    b := output.NewGraphBuilder()
    b.AddNode(output.NewGraphNode("compile", "Compile"))
    b.AddNode(output.NewGraphNode("test", "Test"))
    b.AddEdge(output.NewGraphEdge("compile", "test"))

    g := b.Build()

    dot, _     := graph.RenderDOT(g)
    mermaid, _ := graph.RenderMermaid(g)

    fmt.Println(dot)
    fmt.Println(mermaid)
}`;
