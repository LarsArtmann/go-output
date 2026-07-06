package markup_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/testhelpers"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewHTMLRenderer() {
	html := markup.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Age"})
	html.AddRow([]string{"Alice", "30"})
	html.AddRow([]string{"Bob", "25"})

	fmt.Println(testhelpers.MustRender(html))
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
	data := output.NewTable([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	renderer := markup.NewAsciiDocTableRenderer()
	renderer.SetData(data)

	fmt.Println(testhelpers.MustRender(renderer))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleMarshalXMLFromTable() {
	data := output.NewTable([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	result, err := markup.MarshalXMLFromTable(data)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(string(result))
}
