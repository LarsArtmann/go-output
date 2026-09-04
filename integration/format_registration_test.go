// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"slices"
	"testing"

	"github.com/larsartmann/go-output"
)

// wantFormats is the explicit, complete list of the library's output formats.
// Asserting against this literal (instead of a magic count) makes a registry
// failure name exactly which format is missing or unexpected.
var wantFormats = []output.Format{
	output.FormatTable,
	output.FormatJSON,
	output.FormatCSV,
	output.FormatTSV,
	output.FormatXML,
	output.FormatMarkdown,
	output.FormatD2,
	output.FormatYAML,
	output.FormatHTML,
	output.FormatTree,
	output.FormatMermaid,
	output.FormatDOT,
	output.FormatJSONL,
	output.FormatAsciiDoc,
	output.FormatTOML,
	output.FormatPlantUML,
}

// TestFormatRegistration verifies that importing all sub-modules activates every
// format's shape capabilities via their init() functions. Without the import,
// Format constants exist but Supports() returns false (no shapes registered).
// This test proves the registry-dispatch-via-init pattern works end-to-end.
func TestFormatRegistration(t *testing.T) {
	t.Parallel()

	got := map[output.Format]bool{}
	for _, f := range output.AllFormats {
		got[f] = true
	}

	var missing, unexpected []output.Format

	for _, f := range wantFormats {
		if !got[f] {
			missing = append(missing, f)
		}
	}

	for f := range got {
		if !slices.Contains(wantFormats, f) {
			unexpected = append(unexpected, f)
		}
	}

	if len(missing) > 0 || len(unexpected) > 0 {
		t.Fatalf("format registry drift — missing: %v, unexpected: %v", missing, unexpected)
	}

	shapes := []output.Shape{
		output.ShapeTable,
		output.ShapeTree,
		output.ShapeGraph,
	}

	for _, format := range output.AllFormats {
		t.Run(string(format), func(t *testing.T) {
			t.Parallel()

			registered := slices.ContainsFunc(shapes, func(shape output.Shape) bool {
				return format.Supports(shape)
			})

			if !registered {
				t.Errorf(
					"Format %s has no shapes registered — its sub-module init() may not have run",
					format,
				)
			}
		})
	}

	// At least 14 formats should support ShapeTable (all except tree-only formats).
	tableFormats := output.FormatsForShape(output.ShapeTable)
	if len(tableFormats) < 14 {
		t.Errorf("table format count = %d, want >= 14: %v", len(tableFormats), tableFormats)
	}

	// Verify tree and graph formats are also registered.
	treeFormats := output.FormatsForShape(output.ShapeTree)
	if len(treeFormats) < 2 {
		t.Errorf("tree format count = %d, want >= 2: %v", len(treeFormats), treeFormats)
	}

	graphFormats := output.FormatsForShape(output.ShapeGraph)
	if len(graphFormats) < 4 {
		t.Errorf("graph format count = %d, want >= 4: %v", len(graphFormats), graphFormats)
	}
}
