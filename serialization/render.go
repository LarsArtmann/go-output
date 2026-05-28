package serialization

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

func renderTable(
	data *output.TableData,
	empty string,
	formatName string,
	marshalFunc func(any) ([]byte, error),
) (string, error) {
	if data == nil || len(data.Headers) == 0 {
		return empty, nil
	}

	rows := data.ToMapSlice()

	b, err := marshalFunc(rows)
	if err != nil {
		return "", fmt.Errorf("marshal %s table (%d rows): %w", formatName, len(rows), err)
	}

	return string(b), nil
}

func renderViaRenderer(w io.Writer, data *output.TableData, renderer dataSetter, formatName string) error {
	renderer.SetData(data)

	out, err := renderer.Render()
	if err != nil {
		return fmt.Errorf("render %s: %w", formatName, err)
	}

	_, err = fmt.Fprint(w, out)
	if err != nil {
		return fmt.Errorf("write %s output: %w", formatName, err)
	}

	return nil
}

type dataSetter interface {
	SetData(*output.TableData)
	Render() (string, error)
}
