package output

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

// MarshalYAML encodes v to YAML.
func MarshalYAML(v any) ([]byte, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal yaml: %w", err)
	}
	return data, nil
}

// UnmarshalYAML decodes YAML data into v.
func UnmarshalYAML(data []byte, v any) error {
	if err := yaml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}
	return nil
}
