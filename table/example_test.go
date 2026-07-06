package table_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/table"
	"github.com/larsartmann/go-output/testhelpers"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNew() {
	tbl := table.New()
	tbl.SetHeaders("Name", "Age", "City")
	tbl.AddRow("Alice", "30", "NYC")
	tbl.AddRow("Bob", "25", "LA")

	fmt.Println(testhelpers.MustRender(tbl))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleFromTable() {
	data := output.NewTable([]string{"Name", "Score"})
	data.AddRow([]string{"Alice", "95"})
	data.AddRow([]string{"Bob", "87"})

	tbl := table.FromTable(data, table.WithColorMode(output.ColorModeNever))

	fmt.Println(testhelpers.MustRender(tbl))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleTable_SetFooter() {
	tbl := table.New(table.WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Name", "Score")
	tbl.AddRow("Alice", "95")
	tbl.AddRow("Bob", "87")
	tbl.SetFooter("Total", "182")

	fmt.Println(testhelpers.MustRender(tbl))
}
