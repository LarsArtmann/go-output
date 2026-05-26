// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/testhelpers"
)

type TestProject struct {
	Name       string
	Health     int
	Complexity int
}

func SampleProjects() []TestProject {
	return []TestProject{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
	}
}

// SampleProject returns a single sample project for testing.
func SampleProject() []TestProject {
	return []TestProject{{Name: "Alpha", Health: 90, Complexity: 7}}
}

// sharedTestData contains common test data used across workflow tests.
func sharedTestData() (headers []string, rows [][]string) {
	return []string{"Name", "Value"}, [][]string{
		{"Alpha", "100"},
		{"Beta", "200"},
		{"Gamma", "150"},
	}
}

func TestAllFormatsRender(t *testing.T) {
	t.Parallel()

	projects := SampleProjects()

	formats := []output.Format{
		output.FormatTable,
		output.FormatJSON,
		output.FormatMarkdown,
		output.FormatCSV,
		output.FormatTSV,
		output.FormatXML,
		output.FormatYAML,
		output.FormatHTML,
		output.FormatTree,
		output.FormatD2,
		output.FormatMermaid,
		output.FormatDOT,
		output.FormatJSONL,
		output.FormatAsciiDoc,
		output.FormatTOML,
		output.FormatPlantUML,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			result := renderProject(format, projects)
			if result == "" {
				t.Errorf("Format %s returned empty output", format)
			}
		})
	}
}

func TestStreamingRenderer(t *testing.T) {
	r := markup.NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"Alpha"})

	var buf bytes.Buffer

	err := r.Stream(&buf)
	if err != nil {
		t.Fatalf("Stream failed: %v", err)
	}

	result := buf.String()
	testhelpers.AssertContains(t, result, "<table", "Streaming HTML should contain table tag")
}

func TestTableDataRowEdges(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Row0"})
	data.AddRow([]string{"Row1"})
	data.AddRow([]string{"Row2"})

	edges := data.CreateRowEdges()
	if len(edges) != 2 {
		t.Errorf("Expected 2 edges for 3 rows, got %d", len(edges))
	}
}

func TestTreeNodeDepth(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Root")
	child := output.NewTreeNode("child", "Child")
	grandchild := output.NewTreeNode("grandchild", "Grandchild")

	root.AddChild(child)
	child.AddChild(grandchild)

	if root.Depth() != 0 {
		t.Errorf("Root depth should be 0, got %d", root.Depth())
	}

	if child.Depth() != 1 {
		t.Errorf("Child depth should be 1, got %d", child.Depth())
	}

	if grandchild.Depth() != 2 {
		t.Errorf("Grandchild depth should be 2, got %d", grandchild.Depth())
	}
}

func TestColorModeRenderTableData(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name", "Value"})
	data.AddRow([]string{"Alpha", "100"})

	assertColorMode(t, data, output.FormatMarkdown, "Markdown", output.ColorModeAlways, true, "Alpha")
	assertColorMode(t, data, output.FormatTree, "Tree", output.ColorModeAlways, true, "Alpha")
	assertColorMode(t, data, output.FormatMarkdown, "Markdown", output.ColorModeNever, false, "")
	assertColorMode(t, data, output.FormatTree, "Tree", output.ColorModeNever, false, "")
}

func assertColorMode(
	t *testing.T,
	data *output.TableData,
	format output.Format,
	name string,
	mode output.ColorMode,
	wantANSI bool,
	wantContains string,
) {
	t.Helper()

	t.Run(fmt.Sprintf("%s with color %s", name, mode), func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		err := output.RenderTableData(data, format, output.RenderOptions{
			Writer:    &buf,
			ColorMode: mode,
		})
		if err != nil {
			t.Fatalf("RenderTableData %s: %v", name, err)
		}

		if wantContains != "" {
			testhelpers.AssertContains(t, buf.String(), wantContains, name+" should contain data")
		}

		hasANSI := bytes.Contains(buf.Bytes(), []byte("\033["))
		if wantANSI && !hasANSI {
			t.Errorf("%s output with %s should contain ANSI escape codes", name, mode)
		}

		if !wantANSI && hasANSI {
			t.Errorf("%s output with %s should not contain ANSI escape codes", name, mode)
		}
	})
}
