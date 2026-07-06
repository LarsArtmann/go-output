package serialization_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/serialization"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderJSON() {
	tbl := output.NewTableBuilder().
		SetHeaders("Name", "Status").
		AddRow("Compile", "done").
		AddRow("Test", "done").
		Build()

	json, err := serialization.RenderJSON(tbl)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(json)
}
