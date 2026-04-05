// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"slices"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
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

	if !strings.Contains(result, "invalid format") {
		t.Error("Error message should contain 'invalid format'")
	}

	if !strings.Contains(result, "invalid") {
		t.Error("Error message should contain the invalid value")
	}
}

func TestFormatCategories(t *testing.T) {
	t.Parallel()

	tableFormats := []output.Format{
		output.FormatTable,
		output.FormatJSON,
		output.FormatCSV,
		output.FormatMarkdown,
		output.FormatYAML,
	}

	for _, f := range tableFormats {
		if !f.IsTableFormat() {
			t.Errorf("Format %s should be a table format", f)
		}
	}

	treeFormats := []output.Format{
		output.FormatTree,
		output.FormatHTML,
	}

	for _, f := range treeFormats {
		if !f.IsTreeFormat() {
			t.Errorf("Format %s should be a tree format", f)
		}
	}

	graphFormats := []output.Format{
		output.FormatD2,
		output.FormatMermaid,
		output.FormatDOT,
	}

	for _, f := range graphFormats {
		if !f.IsGraphFormat() {
			t.Errorf("Format %s should be a graph format", f)
		}
	}
}

func TestFormatRegistry(t *testing.T) {
	t.Parallel()

	customFormat := output.Format("custom")

	err := output.Register(customFormat, func() output.Renderer {
		return nil
	})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer output.Unregister(customFormat)

	if !output.IsRegistered(customFormat) {
		t.Error("Custom format should be registered")
	}

	formats := output.RegisteredFormats()

	found := slices.Contains(formats, customFormat)
	if !found {
		t.Error("Custom format should be in registered formats list")
	}
}
