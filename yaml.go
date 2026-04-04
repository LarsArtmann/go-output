package output

import (
	"github.com/go-faster/yaml"
)

// MarshalYAML encodes v to YAML.
func MarshalYAML(v any) ([]byte, error) {
	return marshal("yaml", yaml.Marshal, v)
}

// UnmarshalYAML decodes YAML data into v.
func UnmarshalYAML(data []byte, v any) error {
	return unmarshal("yaml", yaml.Unmarshal, data, v)
}
