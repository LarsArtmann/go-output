package serialization

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"io"
	"strings"

	"github.com/go-faster/yaml"
	"github.com/pelletier/go-toml/v2"

	"github.com/larsartmann/go-output"
)

// WriteJSON writes a Table as JSON directly to the provided writer using
// json.NewEncoder — no intermediate string allocation.
func WriteJSON(w io.Writer, data *output.Table) error {
	if data == nil || len(data.Headers) == 0 {
		return writeEmptyArrayPayload(w, "json")
	}

	encoder := jsontext.NewEncoder(w, jsontext.WithIndentPrefix(""), jsontext.WithIndent("  "))

	rows := data.ToMapSlice()

	if err := json.MarshalEncode(encoder, rows, json.Deterministic(true)); err != nil {
		return fmt.Errorf("encode json table (%d rows): %w", len(rows), err)
	}

	return nil
}

// RenderJSON renders a Table as a JSON string.
func RenderJSON(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteJSON(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteYAML writes a Table as YAML directly to the provided writer using
// yaml.NewEncoder — no intermediate string allocation.
func WriteYAML(w io.Writer, data *output.Table) error {
	if data == nil || len(data.Headers) == 0 {
		return writeEmptyArrayPayload(w, "yaml")
	}

	encoder := yaml.NewEncoder(w)

	rows := data.ToMapSlice()

	if err := encoder.Encode(rows); err != nil {
		return fmt.Errorf("encode yaml table (%d rows): %w", len(rows), err)
	}

	return nil
}

// RenderYAML renders a Table as a YAML string.
func RenderYAML(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteYAML(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteTOML writes a Table as TOML directly to the provided writer using
// toml.NewEncoder — no intermediate string allocation.
// Rows are nested under the key "row" because TOML cannot encode a bare
// top-level array.
func WriteTOML(w io.Writer, data *output.Table) error {
	if data == nil || len(data.Headers) == 0 {
		return nil
	}

	rows := data.ToMapSlice()
	wrapped := map[string]any{tomlTableKey: rows}

	if err := toml.NewEncoder(w).Encode(wrapped); err != nil {
		return fmt.Errorf("encode toml table (%d rows): %w", len(rows), err)
	}

	return nil
}

// RenderTOML renders a Table as a TOML string.
func RenderTOML(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteTOML(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteJSONL writes a Table as JSON Lines directly to the provided writer.
// Each row is encoded as a separate JSON object on its own line via
// NewJSONLWriter — true row-level streaming.
func WriteJSONL(w io.Writer, data *output.Table) error {
	if data == nil || len(data.Headers) == 0 {
		return nil
	}

	jw := NewJSONLWriter(w)

	for _, row := range data.ToMapSlice() {
		if err := jw.Encode(row); err != nil {
			return fmt.Errorf("encode jsonl row: %w", err)
		}
	}

	return jw.Flush()
}

// RenderJSONL renders a Table as a JSONL string.
func RenderJSONL(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteJSONL(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
