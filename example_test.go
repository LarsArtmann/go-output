package output_test

import (
	"fmt"
	"os"

	"github.com/larsartmann/go-output"
)

//nolint:testableexamples
func ExampleFormat_IsValid() {
	fmt.Println(output.FormatCSV.IsValid())
	fmt.Println(output.Format("unknown").IsValid())
	// Output:
	// true
	// false
}

//nolint:testableexamples
func ExampleParseFormat() {
	f, err := output.ParseFormat("json")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println(f)
	// Output: json
}

//nolint:testableexamples
func ExampleColorMode() {
	fmt.Println(output.ColorModeAuto.ShouldColor())
	fmt.Println(output.ColorModeNever.ShouldColor())
	// Output:
	// false
	// false
}

//nolint:testableexamples
func ExampleShape() {
	fmt.Println(output.FormatCSV.Supports(output.ShapeTable))
	fmt.Println(output.FormatCSV.Supports(output.ShapeTree))
	// Output:
	// true
	// false
}

//nolint:testableexamples
func ExampleTable_Validate() {
	data := output.NewTable([]string{"Name", "Count"})
	data.AddRow([]string{"Alice", "10"})
	data.SetFooter([]string{"Total", "10"})

	fmt.Println(data.Validate())
	// Output: <nil>
}
