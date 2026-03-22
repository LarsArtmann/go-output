// Package output provides consistent output formatting for CLI applications.
package output

import (
	"fmt"
	"slices"
)

// OutputFormat represents the available output format options for CLI applications.
type OutputFormat string

// Output format constants.
const (
	// OutputFormatTable renders data as terminal tables with lipgloss styling.
	OutputFormatTable OutputFormat = "table"
	// OutputFormatJSON renders data as formatted JSON.
	OutputFormatJSON OutputFormat = "json"
	// OutputFormatCSV renders data as CSV with headers.
	OutputFormatCSV OutputFormat = "csv"
	// OutputFormatMarkdown renders data as Markdown tables.
	OutputFormatMarkdown OutputFormat = "markdown"
	// OutputFormatD2 renders data as D2 diagram shapes.
	OutputFormatD2 OutputFormat = "d2"
	// OutputFormatYAML renders data as YAML.
	OutputFormatYAML OutputFormat = "yaml"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var outputFormatValues = []OutputFormat{
	OutputFormatTable,
	OutputFormatJSON,
	OutputFormatCSV,
	OutputFormatMarkdown,
	OutputFormatD2,
	OutputFormatYAML,
}

// ParseOutputFormat converts a string to OutputFormat, returning an error if invalid.
func ParseOutputFormat(s string) (OutputFormat, error) {
	for _, v := range outputFormatValues {
		if string(v) == s {
			return v, nil
		}
	}
	return "", fmt.Errorf("invalid output format: %q (allowed: %v)", s, outputFormatValues)
}

// String returns the string representation of the format.
func (f OutputFormat) String() string {
	return string(f)
}

// AllowedValues returns all valid output format values for CLI help text.
func (f OutputFormat) AllowedValues() []string {
	return []string{
		string(OutputFormatTable),
		string(OutputFormatJSON),
		string(OutputFormatCSV),
		string(OutputFormatMarkdown),
		string(OutputFormatD2),
		string(OutputFormatYAML),
	}
}

// IsValid returns true if the format is a valid OutputFormat value.
func (f OutputFormat) IsValid() bool {
	return slices.Contains(outputFormatValues, f)
}
