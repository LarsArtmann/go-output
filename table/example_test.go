package table_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/table"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNew() {
	tbl := table.New()
	tbl.SetHeaders("Name", "Age", "City")
	tbl.AddRow("Alice", "30", "NYC")
	tbl.AddRow("Bob", "25", "LA")

	result, err := tbl.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleFromTableData() {
	data := output.NewTableData([]string{"Name", "Score"})
	data.AddRow([]string{"Alice", "95"})
	data.AddRow([]string{"Bob", "87"})

	tbl := table.FromTableData(data, table.WithColorMode(output.ColorModeNever))

	result, err := tbl.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}
