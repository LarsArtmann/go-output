package serialization_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/testhelpers"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewJSONTableRenderer() {
	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})
	data.AddRow([]string{"Bob", "25"})

	renderer := serialization.NewJSONTableRenderer()
	renderer.SetData(data)

	result, err := renderer.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleMarshalJSONLFromTableData() {
	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})
	data.AddRow([]string{"Bob", "25"})

	result, err := serialization.MarshalJSONLFromTableData(data)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Print(string(result))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewYAMLTableRenderer() {
	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	renderer := serialization.NewYAMLTableRenderer()
	renderer.SetData(data)

	result, err := renderer.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewTOMLTableRenderer() {
	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	renderer := serialization.NewTOMLTableRenderer()
	renderer.SetData(data)

	result, err := renderer.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewJSONTreeRenderer() {
	root := output.NewTreeNode("root", "Root")
	root.AddChild(output.NewTreeNode("child1", "Child 1"))
	root.AddChild(output.NewTreeNode("child2", "Child 2"))

	renderer := serialization.NewJSONTreeRenderer()
	renderer.SetRoot(root)

	result, err := renderer.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewJSONGraphRenderer() {
	renderer := serialization.NewJSONGraphRenderer()
	renderer.SetNodes([]output.GraphNode{
		{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("a"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Node A"),
		},
	})
	renderer.SetEdges([]output.GraphEdge{
		graphtest.NewTestEdge("a", "b", ""),
	})

	fmt.Println(testhelpers.MustRender(renderer))
}
