package output

import (
	"strings"

	"github.com/larsartmann/go-output/enum"
)

// Format represents the available output format options for CLI applications.
type Format string

// Output format constants.
const (
	FormatTable    Format = "table"
	FormatJSON     Format = "json"
	FormatCSV      Format = "csv"
	FormatTSV      Format = "tsv"
	FormatMarkdown Format = "markdown"
	FormatXML      Format = "xml"
	FormatD2       Format = "d2"
	FormatYAML     Format = "yaml"
	FormatHTML     Format = "html"
	FormatTree     Format = "tree"
	FormatMermaid  Format = "mermaid"
	FormatDOT      Format = "dot"
)

// AllFormats contains all valid output format values.
// Use this for testing AllowedValues or generating help text.
var AllFormats = []Format{
	FormatTable,
	FormatJSON,
	FormatCSV,
	FormatTSV,
	FormatMarkdown,
	FormatXML,
	FormatD2,
	FormatYAML,
	FormatHTML,
	FormatTree,
	FormatMermaid,
	FormatDOT,
}

// ParseFormat converts a string to Format, returning an error if invalid.
func ParseFormat(s string) (Format, error) {
	v, err := enum.Parse(AllFormats, s, func(f Format) string { return string(f) })
	if err != nil {
		return "", &InvalidFormatError{Value: s, Allowed: AllFormats}
	}

	return v, nil
}

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

// AllowedValues returns all valid output format values for CLI help text.
func (f Format) AllowedValues() []string {
	return enum.AllowedValues(AllFormats)
}

// IsValid returns true if the format is a valid Format value.
func (f Format) IsValid() bool {
	return enum.Contains(AllFormats, f)
}

// InvalidFormatError represents an invalid format error.
type InvalidFormatError struct {
	Value   string
	Allowed []Format
}

// Error returns a descriptive error message including allowed values.
func (e *InvalidFormatError) Error() string {
	if e.Allowed == nil {
		return "invalid format: " + e.Value
	}

	var b strings.Builder
	b.WriteString("invalid format: ")
	b.WriteString(e.Value)
	b.WriteString(" (allowed: ")

	for i, f := range e.Allowed {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(string(f))
	}

	b.WriteString(")")

	return b.String()
}
