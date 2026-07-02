package markdown_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/markdown"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewMarkdownTable() {
	md := markdown.NewMarkdownTable()
	md.SetHeaders([]string{"Name", "Status", "Duration"})
	md.AddRow([]string{"Build", "✓", "1.2s"})
	md.AddRow([]string{"Test", "✓", "0.5s"})
	md.AddRow([]string{"Deploy", "…", "-"})

	result, err := md.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewMarkdownTableFromData() {
	data := output.NewTableData([]string{"Module", "Coverage"})
	data.AddRow([]string{"core", "92%"})
	data.AddRow([]string{"cli", "78%"})

	md := markdown.NewMarkdownTableFromData(data)
	md.SetColorMode(output.ColorModeNever)

	result, err := md.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleMarkdownTable_SetFooter() {
	md := markdown.NewMarkdownTable()
	md.SetHeaders([]string{"Package", "Downloads"})
	md.AddRow([]string{"v1.0.0", "1234"})
	md.AddRow([]string{"v1.1.0", "5678"})
	md.SetFooter([]string{"Total", "6912"})

	result, err := md.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}
