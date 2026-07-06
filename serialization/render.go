package serialization

import (
	"fmt"

	"github.com/larsartmann/go-output"
)

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
