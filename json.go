package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// MarshalJSON encodes v to JSON.
func MarshalJSON(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal json (%T): %w", v, err)
	}
	return data, nil
}

// MarshalJSONIndent encodes v to indented JSON.
func MarshalJSONIndent(v any, prefix, indent string) ([]byte, error) {
	data, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf("marshal json indent (%T, prefix=%q, indent=%q): %w", v, prefix, indent, err)
	}
	return data, nil
}

// UnmarshalJSON decodes JSON data into v.
func UnmarshalJSON(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal json into %T: %w", v, err)
	}
	return nil
}

// JSONWriter writes JSON output to an io.Writer.
type JSONWriter struct {
	Writer io.Writer
}

// NewJSONWriter creates a new JSONWriter.
func NewJSONWriter(w io.Writer) *JSONWriter {
	return &JSONWriter{Writer: w}
}

// Encode writes v as JSON to the underlying writer.
func (j *JSONWriter) Encode(v any) error {
	encoder := json.NewEncoder(j.Writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("encode json (%T): %w", v, err)
	}
	return nil
}
