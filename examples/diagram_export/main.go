// Package main demonstrates exporting NOM subscriber state as DOT and Mermaid
// diagrams. The subscriber's Store() method projects live activity state as
// GraphNode/GraphEdge slices that any output.GraphRenderer can consume.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/nom"
)

func main() {
	subscriber := nom.NewNOMSubscriber()
	ctx := context.Background()

	// Build a CI pipeline: fetch → compile → {test, lint} → deploy
	fire := func(id, name string, deps ...string) {
		depIDs := make([]nom.ActivityID, 0, len(deps))
		for _, d := range deps {
			depIDs = append(depIDs, nom.ActivityID(d))
		}

		if err := subscriber.OnEvent(ctx, nom.ActivityStarted{
			ID:   nom.ActivityID(id),
			Name: nom.ActivityName(name),
			Deps: depIDs,
		}); err != nil {
			log.Fatalf("ActivityStarted(%q) failed: %v", id, err)
		}
	}

	fire("fetch", "Fetch Dependencies")
	fire("compile", "Compile Sources", "fetch")
	fire("test", "Run Tests", "compile")
	fire("lint", "Lint Code", "compile")
	fire("deploy", "Deploy", "test")

	// Store() returns an ActivityReader that projects Nodes()/Edges()
	reader := subscriber.Store()

	fmt.Printf("Nodes: %d, Edges: %d\n\n", len(reader.Nodes()), len(reader.Edges()))

	// Project the reader's nodes and edges into an immutable Graph, then
	// render via the CQRS pure functions — the canonical v0.30+ API.
	builder := output.NewGraphBuilder()

	for _, node := range reader.Nodes() {
		builder.AddNode(node)
	}

	for _, edge := range reader.Edges() {
		builder.AddEdge(edge)
	}

	g := builder.Build()

	fmt.Println("=== DOT Diagram ===")

	if err := graph.WriteDOT(os.Stdout, g); err != nil {
		log.Fatalf("WriteDOT failed: %v", err)
	}

	fmt.Println()
	fmt.Println("=== Mermaid Diagram ===")

	if err := graph.WriteMermaid(os.Stdout, g); err != nil {
		log.Fatalf("WriteMermaid failed: %v", err)
	}
}
