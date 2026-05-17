package output

import (
	"fmt"

	"github.com/go-faster/yaml"
)

// Compile-time interface checks.
var (
	_ Renderer      = (*YAMLTableRenderer)(nil)
	_ TableRenderer = (*YAMLTableRenderer)(nil)
)

// MarshalYAML encodes v to YAML.
func MarshalYAML(v any) ([]byte, error) {
	return marshal("yaml", yaml.Marshal, v)
}

// UnmarshalYAML decodes YAML data into v.
func UnmarshalYAML(data []byte, v any) error {
	return unmarshal("yaml", yaml.Unmarshal, data, v)
}

// YAMLTableRenderer renders TableData as a YAML sequence of mappings.
// Each row becomes a YAML map with headers as keys.
type YAMLTableRenderer struct {
	tableDataBase
}

// NewYAMLTableRenderer creates a new YAMLTableRenderer.
func NewYAMLTableRenderer() *YAMLTableRenderer {
	return &YAMLTableRenderer{}
}

//nolint:gochecknoglobals // Constant used for empty YAML output.
var emptyYAML = "[]\n"

// Render returns the table data as a YAML string.
func (r *YAMLTableRenderer) Render() (string, error) {
	if r.data == nil || len(r.data.Headers) == 0 {
		return emptyYAML, nil
	}

	rows := r.data.ToMapSlice()

	data, err := yaml.Marshal(rows)
	if err != nil {
		return "", fmt.Errorf("marshal yaml table (%d rows): %w", len(rows), err)
	}

	return string(data), nil
}
