package output_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewTableBuilder() {
	tbl := output.NewTableBuilder().
		SetHeaders("Name", "Status", "Duration").
		AddRow("Compile", "done", "2.1s").
		AddRow("Test", "done", "5.3s").
		SetFooter("Total", "", "7.4s").
		Build()

	fmt.Printf("%d rows, %d cols\n", tbl.RowCount(), tbl.ColCount())
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewTreeBuilder() {
	root := output.NewTreeBuilder().
		SetRoot("build", "Build").
		AddChild("build", "compile", "Compile").
		AddChild("build", "test", "Test").
		AddChild("compile", "lint", "Lint").
		Build()

	fmt.Println(root.Label.Get())
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewGraphBuilder() {
	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("a", "Alpha")).
		AddNode(*output.NewGraphNode("b", "Beta")).
		AddEdge(*output.NewGraphEdge("a", "b")).
		Build()

	fmt.Printf("%d nodes, %d edges\n", len(g.Nodes()), len(g.Edges()))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleTableToGraph() {
	tbl := output.NewTableBuilder().
		SetHeaders("Step").
		AddRow("Fetch").
		AddRow("Build").
		AddRow("Deploy").
		Build()

	g := output.TableToGraph(tbl)

	fmt.Printf("%d nodes\n", len(g.Nodes()))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleGraphToTree() {
	g := output.NewGraphBuilder().
		AddNode(*output.NewGraphNode("root", "Root")).
		AddNode(*output.NewGraphNode("child", "Child")).
		AddEdge(*output.NewGraphEdge("root", "child")).
		Build()
	root := output.GraphToTree(g)

	fmt.Println(root.Label.Get())
}
