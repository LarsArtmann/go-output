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

//nolint:gochecknoinits // Registers TOML TableData marshaler for registry-based dispatch.
func init() {
	output.RegisterTableDataMarshaler(output.FormatTOML, renderTOMLTableData)
}

// MarshalTOML encodes v to TOML.
func MarshalTOML(v any) ([]byte, error) {
	return output.MarshalFormat("toml", toml.Marshal, v)
}

// UnmarshalTOML decodes TOML data into v.
func UnmarshalTOML(data []byte, v any) error {
	return output.UnmarshalFormat("toml", toml.Unmarshal, data, v)
}

// TOMLTableRenderer renders TableData as a TOML array of tables.
type TOMLTableRenderer struct {
	output.TableDataBase
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
	renderer := NewTOMLTableRenderer()
	renderer.SetData(data)

	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render toml: %w", err)
	}

	_, err = fmt.Fprint(w, out)
	if err != nil {
		return fmt.Errorf("write toml output: %w", err)
	}

	return nil
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
