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
	return renderUnknown(w, data, "YAML", yaml.Marshal)
}

// MarshalYAML encodes v to YAML. Panics from the underlying encoder (e.g.
// for channel or function values, which go-faster/yaml panics on rather
// than erroring) are converted to errors so callers never crash on
// unmarshalable input.
func MarshalYAML(v any) (out []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			//nolint:err113 // r is a recovered panic value (any), not an error — %w cannot wrap it.
			err = fmt.Errorf("marshal yaml %T: %v", v, r)
		}
	}()

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
	return WriteYAML(w, data)
}
