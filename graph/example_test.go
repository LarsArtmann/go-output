package graph_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/graph"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleDOTFromTableData() {
	data := &output.TableData{
		Headers: []string{"Name", "Role"},
		Rows: [][]string{
			{"Alice", "Engineer"},
			{"Bob", "Designer"},
		},
	}

	renderer := graph.DOTFromTableData(data)

	result, err := renderer.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleMermaidFlowchartRenderer() {
	data := &output.TableData{
		Headers: []string{"Step", "Action"},
		Rows: [][]string{
			{"1", "Start"},
			{"2", "Process"},
			{"3", "End"},
		},
	}

	renderer := graph.MermaidFlowchartRenderer(data)

	result, err := renderer.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleDOTFromTree() {
	root := &output.TreeNode{
		Label: output.NewBrandedID[output.TreeNodeLabelBrand]("Root"),
		Children: []*output.TreeNode{
			{
				Label: output.NewBrandedID[output.TreeNodeLabelBrand]("Child 1"),
			},
			{
				Label: output.NewBrandedID[output.TreeNodeLabelBrand]("Child 2"),
			},
		},
	}

	renderer := graph.DOTFromTree(root)

	result, err := renderer.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}
