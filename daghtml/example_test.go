package daghtml_test

import (
	"fmt"
	"strings"

	"github.com/larsartmann/go-output/daghtml"
)

func ExampleRender() {
	dag := daghtml.DAG{
		Nodes: []daghtml.Node{
			{ID: "fetch", Label: "📥 Fetch Data", Color: "var(--success)", Tooltip: "Fetch | 42ms"},
			{ID: "parse", Label: "🔍 Parse", Color: "var(--accent)", Tooltip: "Parse | 12ms"},
			{ID: "save", Label: "💾 Save", Color: "var(--success)", Tooltip: "Save | 8ms"},
		},
		Edges: []daghtml.Edge{
			{From: "fetch", To: "parse"},
			{From: "parse", To: "save"},
		},
	}

	html, err := daghtml.Render(dag, daghtml.WithTitle("Data Pipeline"))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Contains DOCTYPE: %v\n", strings.HasPrefix(html, "<!DOCTYPE html>"))
	fmt.Printf("Contains graph JS: %v\n", strings.Contains(html, "initDAGGraph"))
	fmt.Printf("Contains CSS theme: %v\n", strings.Contains(html, "--accent"))
	fmt.Printf("Nodes in JSON: %v\n", strings.Contains(html, `"id":"fetch"`))
	// Output:
	// Contains DOCTYPE: true
	// Contains graph JS: true
	// Contains CSS theme: true
	// Nodes in JSON: true
}
