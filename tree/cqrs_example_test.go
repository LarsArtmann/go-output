package tree_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/tree"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderASCII() {
	root := output.NewTreeBuilder().
		SetRoot("build", "Build").
		AddChild("build", "compile", "Compile").
		AddChild("build", "lint", "Lint").
		Build()

	ascii, err := tree.RenderASCII(root)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(ascii)
}
