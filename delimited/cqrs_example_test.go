package delimited_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/delimited"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleRenderCSV() {
	tbl := output.NewTableBuilder().
		SetHeaders("Name", "Status").
		AddRow("Compile", "done").
		AddRow("Test", "done").
		Build()

	csv, err := delimited.RenderCSV(tbl)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(csv)
}
