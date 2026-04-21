package output

// OutputFormat is a backward compatibility alias for Format.
//
// Deprecated: Use Format directly instead of OutputFormat.
//
//revive:disable-next-line:exported
type OutputFormat = Format

// Backward compatibility aliases for format constants.
const (
	OutputFormatTable    = FormatTable
	OutputFormatJSON     = FormatJSON
	OutputFormatCSV      = FormatCSV
	OutputFormatTSV      = FormatTSV
	OutputFormatXML      = FormatXML
	OutputFormatMarkdown = FormatMarkdown
	OutputFormatD2       = FormatD2
	OutputFormatYAML     = FormatYAML
	OutputFormatHTML     = FormatHTML
	OutputFormatTree     = FormatTree
	OutputFormatMermaid  = FormatMermaid
	OutputFormatDOT      = FormatDOT
)

// ParseOutputFormat is a backward compatibility function that calls ParseFormat.
func ParseOutputFormat(s string) (Format, error) {
	return ParseFormat(s)
}
