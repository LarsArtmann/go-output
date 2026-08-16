package integration

import (
	"cmp"
	"slices"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/delimited"
	"github.com/larsartmann/go-output/markdown"
	"github.com/larsartmann/go-output/testhelpers"
	_ "github.com/larsartmann/go-output/tree" // activate FormatTree registry
)

// User Journey: CLI Developer wants to add output formatting to their tool

func TestCLIDeveloperJourney(t *testing.T) {
	t.Run("can parse format from command line input", func(t *testing.T) {
		t.Parallel()

		// As a CLI developer, I want to parse user input directly
		format, err := output.ParseFormat("json")
		if err != nil {
			t.Fatalf("ParseFormat() error = %v", err)
		}

		if format != output.FormatJSON {
			t.Errorf("ParseFormat() = %v, want %v", format, output.FormatJSON)
		}
	})

	t.Run("receives helpful error for invalid format", func(t *testing.T) {
		t.Parallel()

		// As a CLI developer, I want to give users helpful feedback
		_, err := output.ParseFormat("invalid")
		if err == nil {
			t.Fatal("ParseFormat() expected error for invalid input")
		}

		// Error should guide user to valid options
		errMsg := err.Error()
		if !strings.Contains(errMsg, "allowed") {
			t.Errorf("Error should mention allowed values, got: %s", errMsg)
		}
	})

	t.Run("can list all available formats for help text", func(t *testing.T) {
		t.Parallel()

		// As a CLI developer, I want to show users all available options
		allowed := output.FormatJSON.AllowedValues()
		if len(allowed) < 8 {
			t.Errorf("Expected at least 8 formats, got %d", len(allowed))
		}

		// Should contain common formats
		expected := []string{
			"table",
			"json",
			"csv",
			"tsv",
			"markdown",
			"yaml",
			"html",
			"tree",
			"d2",
			"mermaid",
			"dot",
		}
		for _, exp := range expected {
			if !slices.Contains(allowed, exp) {
				t.Errorf("Expected format %q in allowed values", exp)
			}
		}
	})
}

// User Journey: CLI Developer wants to render data in multiple formats

func TestRenderDataAsJSON(t *testing.T) {
	t.Parallel()

	// Given: User has project data
	data := output.NewTable([]string{"Name", "Health"})
	data.AddRow([]string{"Alpha", "90%"})

	// When: I render it as JSON
	jsonBytes, err := output.MarshalJSONIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("MarshalJSONIndent() error = %v", err)
	}

	// Then: I get valid JSON with the data
	jsonStr := string(jsonBytes)
	testhelpers.AssertContains(t, jsonStr, "Alpha", "JSON should contain project name")
	testhelpers.AssertContains(t, jsonStr, "90%", "JSON should contain health value")
}

func TestRenderDataAsCSV(t *testing.T) {
	t.Parallel()

	// When: I render it as CSV
	var buf strings.Builder

	w := delimited.NewCSVWriter(&buf)
	if err := w.WriteHeader([]string{"Name", "Health"}); err != nil {
		t.Fatalf("WriteHeader failed: %v", err)
	}

	if err := w.WriteRow([]string{"Alpha", "90%"}); err != nil {
		t.Fatalf("WriteRow failed: %v", err)
	}

	w.Flush()

	// Then: I get valid CSV
	csvStr := buf.String()
	testhelpers.AssertContains(t, csvStr, "Name", "CSV should contain header")
	testhelpers.AssertContains(t, csvStr, "Alpha", "CSV should contain data")
}

func TestRenderDataAsMarkdown(t *testing.T) {
	t.Parallel()

	// When: I render it as Markdown
	md := markdown.NewMarkdownTable()
	md.SetHeaders([]string{"Name", "Health"})
	md.AddRow([]string{"Alpha", "90%"})

	// Then: I get valid Markdown table
	testhelpers.RenderAssert(t, md, "| Name", "| Alpha", "|----")
}

// User Journey: CLI Developer wants to handle edge cases gracefully

func TestHandleEdgeCases(t *testing.T) {
	t.Run("empty data renders JSON without panic", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{})
		if _, err := output.MarshalJSONIndent(data, "", ""); err != nil {
			t.Errorf("MarshalJSONIndent on empty data should not error: %v", err)
		}
	})

	t.Run("empty markdown table returns empty string", func(t *testing.T) {
		t.Parallel()

		// Given: User creates table without headers
		md := markdown.NewMarkdownTable()

		// When: I render it
		result, err := md.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		// Then: I get empty string (not panic)
		if result != "" {
			t.Errorf("Empty MarkdownTable should return empty string, got: %q", result)
		}
	})

	t.Run("invalid format gives clear error", func(t *testing.T) {
		t.Parallel()

		// Given: User provides invalid format string
		// When: I parse it
		_, err := output.ParseFormat("invalid_format")

		// Then: I get clear error
		if err == nil {
			t.Error("Expected error for invalid format")
		}
	})
}

// User Journey: CLI Developer wants sorted output — sort data first, then
// render; the renderer must preserve the caller's row order.

func TestSortingBehavior(t *testing.T) {
	type Project struct {
		Name string
	}

	makeProjects := func(names ...string) []Project {
		projects := make([]Project, 0, len(names))
		for _, name := range names {
			projects = append(projects, Project{Name: name})
		}

		return projects
	}

	// renderOrder renders the projects as a markdown table and returns the
	// line positions of the given names, proving the output preserves input order.
	renderOrder := func(t *testing.T, projects []Project) map[string]int {
		t.Helper()

		data := output.NewTable([]string{"Name"})
		for _, p := range projects {
			data.AddRow([]string{p.Name})
		}

		md := markdown.NewMarkdownTableFromTable(data)

		out, err := md.Render()
		if err != nil {
			t.Fatalf("markdown render failed: %v", err)
		}

		lines := strings.Split(out, "\n")

		positions := make(map[string]int, len(projects))
		for i, line := range lines {
			for _, p := range projects {
				if strings.Contains(line, "| "+p.Name+" ") {
					positions[p.Name] = i
				}
			}
		}

		return positions
	}

	t.Run("rendered output preserves sorted order", func(t *testing.T) {
		t.Parallel()

		type testCase struct {
			name     string
			data     []Project
			desc     bool
			first    string
			second   string
			third    string
		}

		cases := []testCase{
			{
				name:   "ascending",
				data:   makeProjects("zebra", "apple", "banana"),
				desc:   false,
				first:  "apple",
				second: "banana",
				third:  "zebra",
			},
			{
				name:   "descending",
				data:   makeProjects("apple", "zebra", "banana"),
				desc:   true,
				first:  "zebra",
				second: "banana",
				third:  "apple",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				slices.SortStableFunc(tc.data, func(a, b Project) int {
					if tc.desc {
						return cmp.Compare(b.Name, a.Name)
					}

					return cmp.Compare(a.Name, b.Name)
				})

				positions := renderOrder(t, tc.data)

				if positions[tc.first] >= positions[tc.second] || positions[tc.second] >= positions[tc.third] {
					t.Errorf(
						"rendered order = %v, want %s < %s < %s",
						positions, tc.first, tc.second, tc.third,
					)
				}
			})
		}
	})

	t.Run("sorting with no comparator leaves render unchanged", func(t *testing.T) {
		t.Parallel()

		data := []Project{{Name: "solo"}}
		slices.SortStableFunc(data, nil)

		positions := renderOrder(t, data)
		if positions["solo"] == 0 {
			t.Error("expected solo row to appear in rendered output")
		}
	})
}
