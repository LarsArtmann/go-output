package delimited_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/delimited"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewCSVWriter() {
	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})
	data.AddRow([]string{"Bob", "25"})

	result, err := delimited.MarshalCSVFromTableData(data)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Print(string(result))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewTSVWriter() {
	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	result, err := delimited.MarshalTSVFromTableData(data)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Print(string(result))
}
