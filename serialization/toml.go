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

//nolint:gochecknoinits // Registers TOML TableData and AnyData marshalers plus format capabilities.
func init() {
	output.RegisterFormatShapes(output.FormatTOML, output.ShapeTable, output.ShapeTree, output.ShapeGraph)
	output.RegisterTableDataMarshaler(output.FormatTOML, renderTOMLTableData)
	output.RegisterAnyDataMarshaler(output.FormatTOML, renderTOMLAnyData)
}

func renderTOMLAnyData(w io.Writer, data any, _ output.RenderOptions) error {
	b, err := toml.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal TOML: %w", err)
	}

	_, err = fmt.Fprintln(w, string(b))
	if err != nil {
		return fmt.Errorf("write TOML output: %w", err)
	}

	return nil
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

// TOMLTableRenderer renders TableData as a TOML array of tables.
type TOMLTableRenderer struct {
	output.TableDataStore
}

// NewTOMLTableRenderer creates a new TOMLTableRenderer.
func NewTOMLTableRenderer() *TOMLTableRenderer {
	return &TOMLTableRenderer{}
}

//nolint:gochecknoglobals // Constant-like value for empty TOML output.
var emptyTOML = "[]\n"

// Render returns the table data as a TOML string.
func (r *TOMLTableRenderer) Render() (string, error) {
	return renderTable(r.Data(), emptyTOML, "toml", toml.Marshal)
}

func renderTOMLTableData(w io.Writer, data *output.TableData, _ output.RenderOptions) error {
	return renderViaRenderer(w, data, NewTOMLTableRenderer(), "toml")
}

// MarshalTOMLFromTableData marshals TableData as TOML.
func MarshalTOMLFromTableData(data *output.TableData) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	renderer := NewTOMLTableRenderer()
	renderer.SetData(data)

	out, err := renderer.Render()
	if err != nil {
		return nil, fmt.Errorf("render toml: %w", err)
	}

	return []byte(out), nil
}
