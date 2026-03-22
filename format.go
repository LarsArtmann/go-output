package output

import "fmt"

type OutputFormat string

const (
	OutputFormatTable    OutputFormat = "table"
	OutputFormatJSON     OutputFormat = "json"
	OutputFormatCSV      OutputFormat = "csv"
	OutputFormatMarkdown OutputFormat = "markdown"
	OutputFormatD2       OutputFormat = "d2"
	OutputFormatYAML     OutputFormat = "yaml"
)

var outputFormatValues = []OutputFormat{
	OutputFormatTable,
	OutputFormatJSON,
	OutputFormatCSV,
	OutputFormatMarkdown,
	OutputFormatD2,
	OutputFormatYAML,
}

func ParseOutputFormat(s string) (OutputFormat, error) {
	for _, v := range outputFormatValues {
		if string(v) == s {
			return v, nil
		}
	}
	return "", fmt.Errorf("invalid output format: %q (allowed: %v)", s, outputFormatValues)
}

func (f OutputFormat) String() string {
	return string(f)
}

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

func (f OutputFormat) IsValid() bool {
	for _, v := range outputFormatValues {
		if f == v {
			return true
		}
	}
	return false
}
