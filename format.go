package output

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
	FormatJSONL    Format = "jsonl"
	FormatAsciiDoc Format = "asciidoc"
	FormatTOML     Format = "toml"
	FormatPlantUML Format = "plantuml"
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
	FormatJSONL,
	FormatAsciiDoc,
	FormatTOML,
	FormatPlantUML,
}

// ParseFormat converts a string to Format, returning an error if invalid.
func ParseFormat(s string) (Format, error) {
	v, err := ParseEnum(AllFormats, s, func(f Format) string { return string(f) })
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
	return EnumAllowedValues(AllFormats)
}

// IsValid returns true if the format is a valid Format value.
func (f Format) IsValid() bool {
	return ContainsEnum(AllFormats, f)
}

// InvalidFormatError represents an invalid format error.
type InvalidFormatError struct {
	Value   string
	Allowed []Format
}

// Error returns a descriptive error message including allowed values.
func (e *InvalidFormatError) Error() string {
	return EnumErrorMessage("format", e.Value, EnumAllowedValues(e.Allowed))
}
