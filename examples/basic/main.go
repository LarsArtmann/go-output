// Package main demonstrates usage of the go-output library.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/table"
)

// Project represents a sample data structure for demonstration.
type Project struct {
	Name       string
	Health     int
	Complexity int
}

func main() {
	projects := []Project{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
		{Name: "Gamma", Health: 85, Complexity: 8},
	}

	// Parse command line format (default to table)
	format := output.OutputFormatTable
	if len(os.Args) > 1 {
		f, err := output.ParseOutputFormat(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid format: %v\n", err)
			os.Exit(1)
		}
		format = f
	}

	// Output in the specified format
	renderOutput(format, projects)
}

func renderOutput(format output.Format, projects []Project) {
	switch format {
	case output.OutputFormatTable:
		renderTable(projects)
	case output.OutputFormatJSON:
		renderJSON(projects)
	case output.OutputFormatMarkdown:
		renderMarkdown(projects)
	case output.OutputFormatCSV:
		renderCSV(projects)
	case output.OutputFormatYAML:
		renderYAML(projects)
	case output.OutputFormatD2:
		renderD2()
	case output.OutputFormatHTML:
		renderHTML(projects)
	case output.OutputFormatTree:
		renderTree(projects)
	case output.OutputFormatMermaid:
		renderMermaid(projects)
	case output.OutputFormatDOT:
		renderDOT(projects)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported format: %q\n", format) //nolint:gosec // G705: format is validated OutputFormat
		fmt.Fprintf(os.Stderr, "Available formats: %v\n", output.FormatTable.AllowedValues())
		os.Exit(1)
	}
}

func renderTable(projects []Project) {
	tbl := table.New()
	tbl.SetHeaders("Name", "Health", "Complexity")
	for _, p := range projects {
		tbl.AddRow(p.Name, strconv.Itoa(p.Health)+"%", strconv.Itoa(p.Complexity)+"/10")
	}
	fmt.Println(tbl.Render())
}

func renderJSON(projects []Project) {
	data, err := output.MarshalJSONIndent(projects, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func renderMarkdown(projects []Project) {
	md := output.NewMarkdownTable()
	md.SetHeaders([]string{"Name", "Health", "Complexity"})
	for _, p := range projects {
		md.AddRow(
			[]string{p.Name, fmt.Sprintf("%d%%", p.Health), fmt.Sprintf("%d/10", p.Complexity)},
		)
	}
	out, err := md.Render()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(out)
}

func renderCSV(projects []Project) {
	w := output.NewCSVWriter(os.Stdout)
	if err := w.WriteHeader([]string{"Name", "Health", "Complexity"}); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing header: %v\n", err)
		os.Exit(1)
	}
	for _, p := range projects {
		if err := w.WriteRow(
			[]string{p.Name, strconv.Itoa(p.Health), strconv.Itoa(p.Complexity)},
		); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing row: %v\n", err)
			os.Exit(1)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Fprintf(os.Stderr, "Error flushing: %v\n", err)
		os.Exit(1)
	}
}

func renderYAML(projects []Project) {
	data, err := output.MarshalYAML(projects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func renderD2() {
	d2 := output.NewD2Diagram()
	d2.AddTable("projects", []output.D2Column{
		{Name: "name", Type: "string"},
		{Name: "health", Type: "int"},
		{Name: "complexity", Type: "int"},
	})
	fmt.Println(d2.Render())
}

func renderHTML(projects []Project) {
	html := output.NewHTMLRenderer()
	html.SetHeaders([]string{"Name", "Health", "Complexity"})
	for _, p := range projects {
		html.AddRow([]string{
			p.Name,
			strconv.Itoa(p.Health) + "%",
			strconv.Itoa(p.Complexity) + "/10",
		})
	}
	fmt.Println(html.RenderFullHTML("Project Health Report"))
}

func renderTree(projects []Project) {
	tree := output.NewASCIITreeRenderer()
	root := output.NewTreeNode("root", "Projects")
	for _, p := range projects {
		projNode := output.NewTreeNode("proj-"+p.Name, p.Name)
		projNode.Metadata["health"] = strconv.Itoa(p.Health) + "%"
		projNode.Metadata["complexity"] = strconv.Itoa(p.Complexity)
		root.AddChild(projNode)
	}
	tree.SetRoot(root)
	fmt.Println(tree.Render())
}

func renderMermaid(projects []Project) {
	data := output.NewTableData([]string{"Name", "Health", "Complexity"})
	for _, p := range projects {
		data.AddRow([]string{
			p.Name,
			strconv.Itoa(p.Health) + "%",
			strconv.Itoa(p.Complexity) + "/10",
		})
	}
	mermaid := output.MermaidFlowchartRenderer(data)
	fmt.Println(mermaid.Render())
}

func renderDOT(projects []Project) {
	data := output.NewTableData([]string{"Name", "Health", "Complexity"})
	for _, p := range projects {
		data.AddRow([]string{
			p.Name,
			strconv.Itoa(p.Health) + "%",
			strconv.Itoa(p.Complexity) + "/10",
		})
	}
	dot := output.DOTFromTableData(data)
	fmt.Println(dot.Render())
}
