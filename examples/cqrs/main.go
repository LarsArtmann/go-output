// Command cqrs demonstrates the Build → Freeze → Render CQRS pipeline
// for all three data shapes: Table, Tree, and Graph.
package main

import (
	"fmt"
	"os"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/delimited"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/tree"
)

func main() {
	// ── Table: Build → Freeze → Render in multiple formats ──
	tbl := output.NewTableBuilder().
		SetHeaders("Step", "Status", "Duration").
		AddRow("Fetch", "done", "0.3s").
		AddRow("Build", "done", "2.1s").
		AddRow("Deploy", "running", "—").
		SetFooter("Total", "", "2.4s").
		Build()

	fmt.Println("=== CSV ===")

	if err := delimited.WriteCSV(os.Stdout, tbl); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println("=== JSON ===")

	if err := serialization.WriteJSON(os.Stdout, tbl); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// ── Tree: Build → Freeze → Render ──
	root := output.NewTreeBuilder().
		SetRoot("ci", "CI Pipeline").
		AddChild("ci", "build", "Build").
		AddChild("ci", "test", "Test").
		AddChild("build", "compile", "Compile").
		AddChild("build", "lint", "Lint").
		Build()

	fmt.Println("=== ASCII Tree ===")

	if err := tree.WriteASCII(os.Stdout, root); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println("=== Markdown Tree ===")

	if err := tree.WriteMarkdown(os.Stdout, root); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// ── Graph: Build → Freeze → Render ──
	g := output.TableToGraph(tbl)

	fmt.Println("=== DOT ===")

	if err := graph.WriteDOT(os.Stdout, g); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
