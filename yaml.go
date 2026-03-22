package output

import (
	"go.yaml.in/yaml/v4"
)

func MarshalYAML(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func UnmarshalYAML(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
