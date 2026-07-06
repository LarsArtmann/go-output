package serialization

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// WriteJSON writes a Table as JSON to the provided writer.
func WriteJSON(w io.Writer, data *output.Table) error {
	return renderJSONTable(w, data, output.RenderOptions{})
}

// RenderJSON renders a Table as a JSON string.
func RenderJSON(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteJSON(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteYAML writes a Table as YAML to the provided writer.
func WriteYAML(w io.Writer, data *output.Table) error {
	return renderYAMLTable(w, data, output.RenderOptions{})
}

// RenderYAML renders a Table as a YAML string.
func RenderYAML(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteYAML(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteTOML writes a Table as TOML to the provided writer.
func WriteTOML(w io.Writer, data *output.Table) error {
	b, err := MarshalTOMLFromTable(data)
	if err != nil {
		return fmt.Errorf("marshal toml: %w", err)
	}

	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write toml output: %w", err)
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

// WriteJSONL writes a Table as JSONL to the provided writer.
func WriteJSONL(w io.Writer, data *output.Table) error {
	b, err := MarshalJSONLFromTable(data)
	if err != nil {
		return fmt.Errorf("marshal jsonl: %w", err)
	}

	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write jsonl output: %w", err)
	}

	return nil
}

// RenderJSONL renders a Table as a JSONL string.
func RenderJSONL(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteJSONL(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
