package output

import (
	"encoding/json"
	"io"
)

func MarshalJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalJSONIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func UnmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

type JSONWriter struct {
	Writer io.Writer
}

func NewJSONWriter(w io.Writer) *JSONWriter {
	return &JSONWriter{Writer: w}
}

func (j *JSONWriter) Encode(v any) error {
	encoder := json.NewEncoder(j.Writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}
