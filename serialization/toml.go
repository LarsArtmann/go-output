package serialization

import (
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"

	"github.com/larsartmann/go-output"
)

var (
	_ output.Renderer      = (*TOMLTableRenderer)(nil)
	_ output.TableRenderer = (*TOMLTableRenderer)(nil)
)

//nolint:gochecknoinits // Registers TOML Table and Unknown marshalers plus format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatTOML, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableMarshaler(output.FormatTOML, renderTOMLTable)
	output.RegisterUnknownRenderer(output.FormatTOML, renderTOMLUnknown)
}

func renderTOMLUnknown(w io.Writer, data any, _ output.RenderOptions) error {
	return renderUnknown(w, data, "TOML", toml.Marshal)
}

// MarshalTOML encodes v to TOML.
func MarshalTOML(v any) ([]byte, error) {
	data, err := toml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal toml %T: %w", v, err)
	}

	return data, nil
}

// UnmarshalTOML decodes TOML data into v.
func UnmarshalTOML(data []byte, v any) error {
	err := toml.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("unmarshal toml into %T: %w", v, err)
	}

	return nil
}

// TOMLTableRenderer renders Table as a TOML array of tables.
type TOMLTableRenderer struct {
	output.TableStore
}

// NewTOMLTableRenderer creates a new TOMLTableRenderer.
func NewTOMLTableRenderer() *TOMLTableRenderer {
	return &TOMLTableRenderer{}
}

//nolint:gochecknoglobals // Constant-like value for empty TOML output.
var emptyTOML = "\n"

// tomlTableKey is the key used for the array-of-tables wrapper.
// TOML cannot encode a bare top-level array, so rows are nested under this key.
const tomlTableKey = "row"

// Render returns the table data as a TOML string using array-of-tables syntax.
func (r *TOMLTableRenderer) Render() (string, error) {
	data := r.Data()
	if data == nil || len(data.Headers) == 0 {
		return emptyTOML, nil
	}

	rows := data.ToMapSlice()

	wrapped := map[string]any{tomlTableKey: rows}

	b, err := toml.Marshal(wrapped)
	if err != nil {
		return "", fmt.Errorf("marshal toml table (%d rows): %w", len(rows), err)
	}

	return string(b), nil
}

func renderTOMLTable(w io.Writer, data *output.Table, _ output.RenderOptions) error {
	return WriteTOML(w, data)
}

// MarshalTOMLFromTable marshals Table as TOML.
func MarshalTOMLFromTable(data *output.Table) ([]byte, error) {
	return output.MarshalViaRenderer(data, "toml", NewTOMLTableRenderer)
}
