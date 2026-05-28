package serialization

import (
	"io"

	"github.com/go-faster/yaml"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*YAMLTableRenderer)(nil)
	_ output.TableRenderer = (*YAMLTableRenderer)(nil)
)

//nolint:gochecknoinits // Registers YAML TableData marshaler for registry-based dispatch.
func init() {
	output.RegisterTableDataMarshaler(output.FormatYAML, renderYAMLTableData)
}

// MarshalYAML encodes v to YAML.
func MarshalYAML(v any) ([]byte, error) {
	return output.MarshalFormat("yaml", yaml.Marshal, v)
}

// UnmarshalYAML decodes YAML data into v.
func UnmarshalYAML(data []byte, v any) error {
	return output.UnmarshalFormat("yaml", yaml.Unmarshal, data, v)
}

// YAMLTableRenderer renders TableData as a YAML sequence of mappings.
type YAMLTableRenderer struct {
	output.TableDataBase
}

// NewYAMLTableRenderer creates a new YAMLTableRenderer.
func NewYAMLTableRenderer() *YAMLTableRenderer {
	return &YAMLTableRenderer{}
}

//nolint:gochecknoglobals // Constant-like value for empty YAML output.
var emptyYAML = "[]\n"

// Render returns the table data as a YAML string.
func (r *YAMLTableRenderer) Render() (string, error) {
	return renderTable(r.Data(), emptyYAML, "yaml", yaml.Marshal)
}

func renderYAMLTableData(w io.Writer, data *output.TableData, _ output.RenderOptions) error {
	return renderViaRenderer(w, data, NewYAMLTableRenderer(), "yaml")
}
