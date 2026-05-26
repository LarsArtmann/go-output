package markup_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/markup"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewHTMLRenderer() {
	html := markup.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Age"})
	html.AddRow([]string{"Alice", "30"})
	html.AddRow([]string{"Bob", "25"})

	result, err := html.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewHTMLRenderer_fullPage() {
	html := markup.NewHTMLRenderer()
	html.SetHeaders([]string{"Name"})
	html.AddRow([]string{"Alice"})

	fmt.Println(html.RenderFullHTML("My Report"))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewAsciiDocTableRenderer() {
	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	renderer := markup.NewAsciiDocTableRenderer()
	renderer.SetData(data)

	result, err := renderer.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleMarshalXMLFromTableData() {
	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	result, err := markup.MarshalXMLFromTableData(data)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(string(result))
}
