package serialization

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*JSONTableRenderer)(nil)
	_ output.TableRenderer = (*JSONTableRenderer)(nil)
)

//nolint:gochecknoinits // Registers JSON Table and Unknown marshalers plus format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatJSON, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableMarshaler(output.FormatJSON, renderJSONTable)
	output.RegisterUnknownRenderer(output.FormatJSON, renderJSONUnknown)
}

func renderJSONUnknown(w io.Writer, data any, _ output.RenderOptions) error {
	b, err := json.Marshal(data, json.Deterministic(true), jsontext.WithIndentPrefix(""), jsontext.WithIndent("  "))
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	_, err = fmt.Fprintln(w, string(b))
	if err != nil {
		return fmt.Errorf("write JSON output: %w", err)
	}

	return nil
}

// MarshalJSON encodes v to JSON.
func MarshalJSON(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal json %T: %w", v, err)
	}

	return data, nil
}

// UnmarshalJSON decodes JSON data into v.
func UnmarshalJSON(data []byte, v any) error {
	err := json.Unmarshal(data, v)
	if err != nil {
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
	encoder := jsontext.NewEncoder(j.Writer, jsontext.WithIndentPrefix(""), jsontext.WithIndent("  "))

	err := json.MarshalEncode(encoder, v)
	if err != nil {
		return fmt.Errorf("encode json (%T): %w", v, err)
	}

	return nil
}

// JSONTableRenderer renders Table as a JSON array of objects.
type JSONTableRenderer struct {
	output.TableStore
}

// NewJSONTableRenderer creates a new JSONTableRenderer.
func NewJSONTableRenderer() *JSONTableRenderer {
	return &JSONTableRenderer{}
}

// Render returns the table data as a JSON string.
func (r *JSONTableRenderer) Render() (string, error) {
	return renderTable(r.Data(), "[]", "json", func(v any) ([]byte, error) {
		return json.Marshal(v, json.Deterministic(true), jsontext.WithIndentPrefix(""), jsontext.WithIndent("  "))
	})
}

func renderJSONTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return WriteJSON(w, data)
}
