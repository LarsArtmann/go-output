package serialization

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*JSONTableRenderer)(nil)
	_ output.TableRenderer = (*JSONTableRenderer)(nil)
)

//nolint:gochecknoinits // Registers JSON TableData marshaler for registry-based dispatch.
func init() {
	output.RegisterTableDataMarshaler(output.FormatJSON, renderJSONTableData)
}

// MarshalJSON encodes v to JSON.
func MarshalJSON(v any) ([]byte, error) {
	return output.MarshalFormat("json", json.Marshal, v)
}

// UnmarshalJSON decodes JSON data into v.
func UnmarshalJSON(data []byte, v any) error {
	return output.UnmarshalFormat("json", json.Unmarshal, data, v)
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

// JSONTableRenderer renders TableData as a JSON array of objects.
type JSONTableRenderer struct {
	output.TableDataBase
}

// NewJSONTableRenderer creates a new JSONTableRenderer.
func NewJSONTableRenderer() *JSONTableRenderer {
	return &JSONTableRenderer{}
}

// Render returns the table data as a JSON string.
func (r *JSONTableRenderer) Render() (string, error) {
	return renderTable(r.Data(), "[]", "json", func(v any) ([]byte, error) {
		return json.MarshalIndent(v, "", "  ")
	})
}

func renderJSONTableData(w io.Writer, data *output.TableData, _ output.RenderOptions) error {
	return renderViaRenderer(w, data, NewJSONTableRenderer(), "json")
}
