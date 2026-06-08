package serialization

import (
	"fmt"
	"io"

	"github.com/go-faster/yaml"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*YAMLTableRenderer)(nil)
	_ output.TableRenderer = (*YAMLTableRenderer)(nil)
)

//nolint:gochecknoinits // Registers YAML TableData marshaler and format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatYAML, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableDataMarshaler(output.FormatYAML, renderYAMLTableData)
}

// MarshalYAML encodes v to YAML.
func MarshalYAML(v any) ([]byte, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal yaml %T: %w", v, err)
	}

	return data, nil
}

// UnmarshalYAML decodes YAML data into v.
func UnmarshalYAML(data []byte, v any) error {
	err := yaml.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("unmarshal yaml into %T: %w", v, err)
	}

	return nil
}

// YAMLTableRenderer renders TableData as a YAML sequence of mappings.
type YAMLTableRenderer struct {
	output.TableDataStore
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
