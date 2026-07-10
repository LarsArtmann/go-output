package serialization

import (
	"bufio"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*JSONLTableRenderer)(nil)
	_ output.TableRenderer = (*JSONLTableRenderer)(nil)
)

//nolint:gochecknoinits // Registers JSONL Table marshaler and format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatJSONL, output.ShapeTable)
	output.RegisterTableMarshaler(output.FormatJSONL, renderJSONLTable)
}

// JSONLWriter writes JSON Lines output — one JSON object per line.
type JSONLWriter struct {
	writer  *bufio.Writer
	encoder *jsontext.Encoder
}

// NewJSONLWriter creates a new JSONLWriter.
func NewJSONLWriter(w io.Writer) *JSONLWriter {
	bufWriter := bufio.NewWriter(w)

	return &JSONLWriter{
		writer:  bufWriter,
		encoder: jsontext.NewEncoder(bufWriter),
	}
}

// Encode writes v as a single JSON line to the underlying writer.
func (j *JSONLWriter) Encode(v any) error {
	if err := json.MarshalEncode(j.encoder, v, json.Deterministic(true)); err != nil {
		return fmt.Errorf("encode jsonl (%T): %w", v, err)
	}

	return nil
}

// Flush flushes the underlying buffered writer.
func (j *JSONLWriter) Flush() error {
	if err := j.writer.Flush(); err != nil {
		return fmt.Errorf("flush jsonl writer: %w", err)
	}

	return nil
}

// JSONLTableRenderer renders Table as JSON Lines.
type JSONLTableRenderer struct {
	output.TableStore
}

// NewJSONLTableRenderer creates a new JSONLTableRenderer.
func NewJSONLTableRenderer() *JSONLTableRenderer {
	return &JSONLTableRenderer{}
}

//nolint:gochecknoglobals // Constant-like value for empty JSONL output.
var emptyJSONL = "\n"

// Render returns the table data as JSON Lines.
func (r *JSONLTableRenderer) Render() (string, error) {
	data := r.Data()
	if data == nil || len(data.Headers) == 0 {
		return emptyJSONL, nil
	}

	return marshalJSONLRows(data.ToMapSlice())
}

// MarshalJSONLFromTable marshals Table as JSON Lines.
func MarshalJSONLFromTable(data *output.Table) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	rows := data.ToMapSlice()
	if len(rows) == 0 {
		return nil, nil
	}

	out, err := marshalJSONLRows(rows)
	if err != nil {
		return nil, err
	}

	return []byte(out), nil
}

func renderJSONLTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return WriteJSONL(w, data)
}

func marshalJSONLRows(rows []map[string]string) (string, error) {
	var buf strings.Builder

	for _, row := range rows {
		b, err := json.Marshal(row, json.Deterministic(true))
		if err != nil {
			return "", fmt.Errorf("marshal jsonl row (%d fields): %w", len(row), err)
		}

		buf.Write(b)
		buf.WriteByte('\n')
	}

	return buf.String(), nil
}
