// Package main demonstrates exporting NOM subscriber state as DOT and Mermaid
// diagrams. The subscriber's Store() method projects live activity state as
// GraphNode/GraphEdge slices that any output.GraphRenderer can consume.
package main

import (
	"context"
	"fmt"

	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/nom"
)

func main() {
	subscriber := nom.NewNOMStyleSubscriber()
	ctx := context.Background()

	// Build a CI pipeline: fetch → compile → {test, lint} → deploy
	fire := func(id, name string, deps ...string) {
		depIDs := make([]nom.ActivityID, 0, len(deps))
		for _, d := range deps {
			depIDs = append(depIDs, nom.ActivityID(d))
		}

		_ = subscriber.OnEvent(ctx, nom.ActivityStarted{
			ID:   nom.ActivityID(id),
			Name: nom.ActivityName(name),
			Deps: depIDs,
		})
	}

	fire("fetch", "Fetch Dependencies")
	fire("compile", "Compile Sources", "fetch")
	fire("test", "Run Tests", "compile")
	fire("lint", "Lint Code", "compile")
	fire("deploy", "Deploy", "test")

	// Store() returns an ActivityReader that projects Nodes()/Edges()
	reader := subscriber.Store()

	fmt.Printf("Nodes: %d, Edges: %d\n\n", len(reader.Nodes()), len(reader.Edges()))

	// Export as DOT
	dot := graph.NewDOTRenderer()
	dot.SetNodes(reader.Nodes())
	dot.SetEdges(reader.Edges())

	dotOut, err := dot.Render()
	if err != nil {
		panic(err)
	}

	fmt.Println("=== DOT Diagram ===")
	fmt.Println(dotOut)

	// Export as Mermaid
	mermaid := graph.NewMermaidRenderer()
	mermaid.SetNodes(reader.Nodes())
	mermaid.SetEdges(reader.Edges())

	mermaidOut, err := mermaid.Render()
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Mermaid Diagram ===")
	fmt.Println(mermaidOut)
}
