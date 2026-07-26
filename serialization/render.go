package serialization

import (
	"fmt"
	"io"

	"github.com/larsartmann/go-output"
)

func renderUnknown(
	w io.Writer,
	data any,
	formatName string,
	marshalFunc func(any) ([]byte, error),
) error {
	b, err := marshalFunc(data)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", formatName, err)
	}

	return output.WriteRendered(w, formatName, string(b))
}

func renderTable(
	data *output.Table,
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
