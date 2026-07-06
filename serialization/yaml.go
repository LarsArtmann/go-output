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

//nolint:gochecknoinits // Registers YAML Table and Unknown marshalers plus format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatYAML, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableMarshaler(output.FormatYAML, renderYAMLTable)
	output.RegisterUnknownRenderer(output.FormatYAML, renderYAMLUnknown)
}

func renderYAMLUnknown(w io.Writer, data any, _ output.RenderOptions) error {
	b, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal YAML: %w", err)
	}

	_, err = fmt.Fprintln(w, string(b))
	if err != nil {
		return fmt.Errorf("write YAML output: %w", err)
	}

	return nil
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

// YAMLTableRenderer renders Table as a YAML sequence of mappings.
type YAMLTableRenderer struct {
	output.TableStore
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

func renderYAMLTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return renderViaRenderer(w, data, NewYAMLTableRenderer(), "yaml")
}
