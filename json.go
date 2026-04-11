package output

import (
	"encoding/json"
	"fmt"
	"io"
)

// MarshalJSON encodes v to JSON.
func MarshalJSON(v any) ([]byte, error) {
	return marshal("json", json.Marshal, v)
}

// UnmarshalJSON decodes JSON data into v.
func UnmarshalJSON(data []byte, v any) error {
	return unmarshal("json", json.Unmarshal, data, v)
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

	err := encoder.Encode(v)
	if err != nil {
		return fmt.Errorf("encode json (%T): %w", v, err)
	}

	return nil
}
