package tree_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/tree"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewASCIITreeRenderer() {
	root := &output.TreeNode{
		Label: output.NewBrandedID[output.TreeNodeLabelBrand]("project"),
		Children: []*output.TreeNode{
			{Label: output.NewBrandedID[output.TreeNodeLabelBrand]("src")},
			{Label: output.NewBrandedID[output.TreeNodeLabelBrand]("docs")},
			{Label: output.NewBrandedID[output.TreeNodeLabelBrand]("tests")},
		},
	}

	r := tree.NewASCIITreeRenderer()
	r.SetRoot(root)

	result, err := r.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleASCIITreeRenderer_SetColorMode() {
	root := &output.TreeNode{
		Label: output.NewBrandedID[output.TreeNodeLabelBrand]("root"),
		Children: []*output.TreeNode{
			{Label: output.NewBrandedID[output.TreeNodeLabelBrand]("child1")},
			{Label: output.NewBrandedID[output.TreeNodeLabelBrand]("child2")},
		},
	}

	r := tree.NewASCIITreeRenderer()
	r.SetColorMode(output.ColorModeNever)
	r.SetRoot(root)

	result, err := r.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}
