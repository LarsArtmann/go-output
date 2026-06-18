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

type diagramEvent struct {
	eventType string
	aID       nom.ActivityID
	aName     nom.ActivityName
	deps      []nom.ActivityID
}

func (e *diagramEvent) GetEventType() string              { return e.eventType }
func (e *diagramEvent) GetActivityID() nom.ActivityID     { return e.aID }
func (e *diagramEvent) GetActivityName() nom.ActivityName { return e.aName }
func (e *diagramEvent) GetDependencies() []nom.ActivityID { return e.deps }

func main() {
	subscriber := nom.NewNOMStyleSubscriber()
	ctx := context.Background()

	// Build a CI pipeline: fetch → compile → {test, lint} → deploy
	fire := func(eventType, id, name string, deps ...string) {
		depIDs := make([]nom.ActivityID, len(deps))
		for i, d := range deps {
			depIDs[i] = nom.ActivityID(d)
		}

		_ = subscriber.OnEvent(ctx, &diagramEvent{
			eventType: eventType,
			aID:       nom.ActivityID(id),
			aName:     nom.ActivityName(name),
			deps:      depIDs,
		})
	}

	fire(nom.EventActivityStarted, "fetch", "Fetch Dependencies")
	fire(nom.EventActivityStarted, "compile", "Compile Sources", "fetch")
	fire(nom.EventActivityStarted, "test", "Run Tests", "compile")
	fire(nom.EventActivityStarted, "lint", "Lint Code", "compile")
	fire(nom.EventActivityStarted, "deploy", "Deploy", "test")

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
