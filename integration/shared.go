package integration

import (
	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/table"
)

// TestProject represents a project for integration testing.
type TestProject struct {
	Name       string
	Health     int
	Complexity int
}

// SampleProjects returns sample projects for testing.
func SampleProjects() []TestProject {
	return []TestProject{
		{Name: "Alpha", Health: 90, Complexity: 7},
		{Name: "Beta", Health: 75, Complexity: 5},
	}
}

// sharedAlphaProject returns a test project with default values.
func sharedAlphaProject() []TestProject {
	return []TestProject{
		{Name: "Alpha", Health: 90, Complexity: 7},
	}
}

// sharedTestData contains common test data used across workflow tests.
func sharedTestData() (headers []string, rows [][]string) {
	return []string{"Name", "Value"}, [][]string{
		{"Alpha", "100"},
		{"Beta", "200"},
		{"Gamma", "150"},
	}
}

// renderTableFormat renders projects as a table.
func renderTableFormat(projects []TestProject) string {
	tbl := table.New()
	tbl.SetHeaders("Name", "Health", "Complexity")

	for _, p := range projects {
		tbl.AddRow(p.Name, formatHealth(p.Health), formatComplexity(p.Complexity))
	}

	return tbl.Render()
}

// formatHealth formats health as a percentage string.
func formatHealth(health int) string {
	return output.Sprintf("%d%%", health)
}

// formatComplexity formats complexity as a ratio string.
func formatComplexity(complexity int) string {
	return output.Sprintf("%d/10", complexity)
}
