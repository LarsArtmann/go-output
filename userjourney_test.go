package output_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/internal/testutils"
	"github.com/larsartmann/go-output/sort" //nolint:staticcheck // intentionally testing deprecated package
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
	data := output.NewTableData([]string{"Name", "Health"})
	data.AddRow([]string{"Alpha", "90%"})

	// When: I render it as JSON
	jsonBytes, err := output.MarshalJSONIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("MarshalJSONIndent() error = %v", err)
	}

	// Then: I get valid JSON with the data
	jsonStr := string(jsonBytes)
	testutils.AssertContains(t, jsonStr, "Alpha", "JSON should contain project name")
	testutils.AssertContains(t, jsonStr, "90%", "JSON should contain health value")
}

func TestRenderDataAsCSV(t *testing.T) {
	t.Parallel()

	// When: I render it as CSV
	var buf strings.Builder

	w := output.NewCSVWriter(&buf)
	_ = w.WriteHeader([]string{"Name", "Health"})
	_ = w.WriteRow([]string{"Alpha", "90%"})
	w.Flush()

	// Then: I get valid CSV
	csvStr := buf.String()
	testutils.AssertContains(t, csvStr, "Name", "CSV should contain header")
	testutils.AssertContains(t, csvStr, "Alpha", "CSV should contain data")
}

func TestRenderDataAsMarkdown(t *testing.T) {
	t.Parallel()

	// When: I render it as Markdown
	mdStr := testutils.RenderSampleMarkdownTable()

	// Then: I get valid Markdown table
	testutils.AssertContains(t, mdStr, "| Name", "Markdown should contain header row")
	testutils.AssertContains(t, mdStr, "| Alpha", "Markdown should contain data row")
	testutils.AssertContains(t, mdStr, "|----", "Markdown should contain separator")
}

func TestRenderDataAsYAML(t *testing.T) {
	t.Parallel()

	// Given: User has project data
	data := output.NewTableData([]string{"Name", "Health"})
	data.AddRow([]string{"Alpha", "90%"})

	// When: I render it as YAML
	yamlBytes, err := output.MarshalYAML(data)
	if err != nil {
		t.Fatalf("MarshalYAML() error = %v", err)
	}

	// Then: I get valid YAML
	yamlStr := string(yamlBytes)
	testutils.AssertContains(t, yamlStr, "Name", "YAML should contain field name")
	testutils.AssertContains(t, yamlStr, "Alpha", "YAML should contain data")
}

// User Journey: CLI Developer wants to handle edge cases gracefully

func TestHandleEdgeCases(t *testing.T) {
	testutils.RunEmptyDataRendersJSONWithoutPanic(t)

	t.Run("empty markdown table returns empty string", func(t *testing.T) {
		t.Parallel()

		// Given: User creates table without headers
		md := output.NewMarkdownTable()

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

// User Journey: CLI Developer wants consistent sorting behavior

func TestSortingBehavior(t *testing.T) {
	type Project struct {
		Name string
	}

	makeProjects := func(names ...string) []Project {
		projects := make([]Project, len(names))
		for i, name := range names {
			projects[i] = Project{Name: name}
		}

		return projects
	}

	t.Run("can sort by name", func(t *testing.T) {
		t.Parallel()

		type testCase struct {
			name     string
			data     []Project
			desc     bool
			expected string
		}

		cases := []testCase{
			{
				name:     "ascending",
				data:     makeProjects("zebra", "apple", "banana"),
				desc:     false,
				expected: "apple",
			},
			{
				name:     "descending",
				data:     makeProjects("apple", "zebra", "banana"),
				desc:     true,
				expected: "zebra",
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				sorted := sort.New(tc.data, output.SortByName, tc.desc)
				sorted.WithLessFunc(sort.ByField(func(p Project) string { return p.Name }))
				sorted.Sort()

				if tc.data[0].Name != tc.expected {
					t.Errorf("Expected first item to be %q, got %s", tc.expected, tc.data[0].Name)
				}
			})
		}
	})

	t.Run("invalid sort field sorts without panic", func(t *testing.T) {
		t.Parallel()

		// Given: Data with invalid sort field
		data := []Project{{Name: "test"}}

		// When: I try to sort by invalid field - should not panic
		sorted := sort.New(data, output.SortBy("invalid"), false)
		sorted.Sort()

		// Then: No panic occurred (no LessFunc = no-op)
	})
}
