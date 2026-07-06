// Package main demonstrates usage of the go-output library.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/examples/shared"
)

// Project represents a sample data structure for demonstration.
type Project struct {
	Name       string
	Health     int
	Complexity int
}

// rendererFunc is a function that renders projects in a specific format.
type rendererFunc func([]Project)

// getRenderers returns a map of format to renderer function.
func getRenderers() map[output.Format]rendererFunc {
	return map[output.Format]rendererFunc{
		output.FormatTable:    renderTable,
		output.FormatJSON:     renderJSON,
		output.FormatMarkdown: renderMarkdown,
		output.FormatCSV:      renderCSV,
		output.FormatTSV:      renderTSV,
		output.FormatXML:      renderXML,
		output.FormatYAML:     renderYAML,
		output.FormatD2:       renderD2,
		output.FormatHTML:     renderHTML,
		output.FormatTree:     renderTree,
		output.FormatMermaid:  renderMermaid,
		output.FormatDOT:      renderDOT,
		output.FormatJSONL:    renderJSONL,
		output.FormatAsciiDoc: renderAsciiDoc,
		output.FormatTOML:     renderTOML,
		output.FormatPlantUML: renderPlantUML,
	}
}

// colorMode is the global color mode for output, parsed from --color flag.
var colorMode output.ColorMode

func main() {
	projects := []Project{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
		{Name: "Gamma", Health: 85, Complexity: 8},
	}

	format := output.FormatTable
	colorMode = output.ColorModeAuto

	for i := 1; i < len(os.Args); i++ {
		switch {
		case os.Args[i] == "--color" && i+1 < len(os.Args):
			i++

			cm, err := output.ParseColorMode(os.Args[i])
			if err != nil {
				shared.HandleError(err)
			}

			colorMode = cm
		default:
			f, err := output.ParseFormat(os.Args[i])
			if err != nil {
				shared.HandleError(err)
			}

			format = f
		}
	}

	renderOutput(format, projects)
}

func renderOutput(format output.Format, projects []Project) {
	renderers := getRenderers()
	if renderer, ok := renderers[format]; ok {
		renderer(projects)

		return
	}

	// Handle unknown format safely - format is validated by ParseFormat
	//nolint:gosec // format is validated enum type
	fmt.Fprintf(os.Stderr, "Unsupported format: %s\n", format)
	fmt.Fprintf(os.Stderr, "Available formats: %v\n", output.FormatTable.AllowedValues())
	os.Exit(1)
}

// projectHeaders defines the common headers for project data.
var projectHeaders = []string{"Name", "Health", "Complexity"}

// projectToRow converts a Project to a formatted row slice.
func projectToRow(p Project) []string {
	return []string{
		p.Name,
		strconv.Itoa(p.Health) + "%",
		strconv.Itoa(p.Complexity) + "/10",
	}
}

// projectsToTable creates Table from projects.
func projectsToTable(projects []Project) *output.Table {
	data := output.NewTable(projectHeaders)
	for _, p := range projects {
		data.AddRow(projectToRow(p))
	}

	data.SetFooter([]string{
		"TOTAL",
		fmt.Sprintf("%d projects", len(projects)),
		"-",
	})

	return data
}
