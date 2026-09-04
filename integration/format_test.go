// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestFormatParseRoundtrip(t *testing.T) {
	t.Parallel()

	formats := output.FormatTable.AllowedValues()
	for _, formatStr := range formats {
		t.Run(formatStr, func(t *testing.T) {
			t.Parallel()

			format, err := output.ParseFormat(formatStr)
			if err != nil {
				t.Errorf("ParseFormat(%q) failed: %v", formatStr, err)

				return
			}

			if !format.IsValid() {
				t.Errorf("Format %q should be valid after parsing", formatStr)
			}

			if string(format) != formatStr {
				t.Errorf("Format.String() = %q, want %q", format, formatStr)
			}
		})
	}
}

func TestInvalidFormatError(t *testing.T) {
	t.Parallel()

	err := &output.InvalidFormatError{
		Value:   "invalid",
		Allowed: nil,
	}
	result := err.Error()

	testhelpers.AssertContains(
		t,
		result,
		"invalid format",
		"Error message should contain 'invalid format'",
	)
	testhelpers.AssertContains(
		t,
		result,
		"invalid",
		"Error message should contain the invalid value",
	)
}

func TestFormatCategories(t *testing.T) {
	t.Parallel()

	// This test pins the LOAD-BEARING edges of the shape matrix: positive
	// anchors for the shape-agnostic formats, negative boundaries a format
	// must not claim, and one structural invariant. The full matrix is
	// declared in root's shape.go plus each sub-module's init() — restating
	// all of it here would just duplicate the source.

	// Positive anchors: serialization formats are shape-agnostic.
	for _, f := range []output.Format{output.FormatJSON, output.FormatYAML, output.FormatTOML} {
		for _, shape := range []output.Shape{output.ShapeTable, output.ShapeTree, output.ShapeGraph} {
			if !f.Supports(shape) {
				t.Errorf("Format %s is shape-agnostic and must support %s", f, shape)
			}
		}
	}

	// Negative boundaries: formats that must NOT claim a shape.
	negative := []struct {
		format output.Format
		shape  output.Shape
	}{
		{output.FormatCSV, output.ShapeTree},
		{output.FormatCSV, output.ShapeGraph},
		{output.FormatTSV, output.ShapeGraph},
		{output.FormatMarkdown, output.ShapeGraph},
		{output.FormatJSONL, output.ShapeTree},
		{output.FormatXML, output.ShapeTree},
		{output.FormatAsciiDoc, output.ShapeTree},
		{output.FormatTable, output.ShapeTree},
		{output.FormatTree, output.ShapeTable},
		{output.FormatTree, output.ShapeGraph},
		{output.FormatHTML, output.ShapeGraph},
	}
	for _, tc := range negative {
		if tc.format.Supports(tc.shape) {
			t.Errorf("Format %s should NOT support %s", tc.format, tc.shape)
		}
	}

	// Structural invariant: every registered format supports at least one
	// shape — catches a sub-module whose init() registration broke.
	for _, f := range output.AllFormats {
		if !f.Supports(output.ShapeTable) && !f.Supports(output.ShapeTree) && !f.Supports(output.ShapeGraph) {
			t.Errorf("Format %s supports no shapes — its init() registration is broken", f)
		}
	}
}
