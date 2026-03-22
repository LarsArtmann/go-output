package main

import (
	"fmt"
	"os"

	"github.com/larsartmann/go-output"
)

func main() {
	// Define sample data
	type Project struct {
		Name       string
		Health     int
		Complexity int
	}

	projects := []Project{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
		{Name: "Gamma", Health: 85, Complexity: 8},
	}

	// Parse command line format (default to JSON)
	format := output.OutputFormatJSON
	if len(os.Args) > 1 {
		f, err := output.ParseOutputFormat(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid format: %v\n", err)
			os.Exit(1)
		}
		format = f
	}

	// Output in the specified format
	switch format {
	case output.OutputFormatJSON:
		data, err := output.MarshalJSONIndent(projects, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))

	case output.OutputFormatMarkdown:
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

	case output.OutputFormatCSV:
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

	case output.OutputFormatYAML:
		data, err := output.MarshalYAML(projects)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(data))

	case output.OutputFormatD2:
		d2 := output.NewD2Diagram()
		d2.AddTable("projects", []output.D2Column{
			{Name: "name", Type: "string"},
			{Name: "health", Type: "int"},
			{Name: "complexity", Type: "int"},
		})
		fmt.Println(d2.Render())

	default:
		fmt.Fprintf(os.Stderr, "Unsupported format: %s\n", format)
		os.Exit(1)
	}
}
