// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"slices"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

// TestFormatRegistration verifies that importing all sub-modules activates every
// format's shape capabilities via their init() functions. Without the import,
// Format constants exist but Supports() returns false (no shapes registered).
// This test proves the registry-dispatch-via-init pattern works end-to-end.
func TestFormatRegistration(t *testing.T) {
	t.Parallel()

	if len(output.AllFormats) != 16 {
		t.Fatalf("expected 16 formats, got %d", len(output.AllFormats))
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
	testhelpers.AssertEqual(t, "table format count", nil, len(tableFormats) >= 14, true)

	// Verify tree and graph formats are also registered.
	treeFormats := output.FormatsForShape(output.ShapeTree)
	testhelpers.AssertEqual(t, "tree format count", nil, len(treeFormats) >= 2, true)

	graphFormats := output.FormatsForShape(output.ShapeGraph)
	testhelpers.AssertEqual(t, "graph format count", nil, len(graphFormats) >= 4, true)
}
