package markdown_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/markdown"
	"github.com/larsartmann/go-output/testhelpers"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewMarkdownTable() {
	md := markdown.NewMarkdownTable()
	md.SetHeaders([]string{"Name", "Status", "Duration"})
	md.AddRow([]string{"Build", "✓", "1.2s"})
	md.AddRow([]string{"Test", "✓", "0.5s"})
	md.AddRow([]string{"Deploy", "…", "-"})

	fmt.Println(testhelpers.MustRender(md))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewMarkdownTableFromTable() {
	data := output.NewTable([]string{"Module", "Coverage"})
	data.AddRow([]string{"core", "92%"})
	data.AddRow([]string{"cli", "78%"})

	md := markdown.NewMarkdownTableFromTable(data)
	md.SetColorMode(output.ColorModeNever)

	fmt.Println(testhelpers.MustRender(md))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleMarkdownTable_SetFooter() {
	md := markdown.NewMarkdownTable()
	md.SetHeaders([]string{"Package", "Downloads"})
	md.AddRow([]string{"v1.0.0", "1234"})
	md.AddRow([]string{"v1.1.0", "5678"})
	md.SetFooter([]string{"Total", "6912"})

	fmt.Println(testhelpers.MustRender(md))
}
