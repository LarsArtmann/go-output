package output

import (
	"encoding/json"
	"fmt"

	id "github.com/larsartmann/go-branded-id"
)

// MarshalJSONIndent encodes v to indented JSON.
func MarshalJSONIndent(v any, prefix, indent string) ([]byte, error) {
	data, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal json indent (prefix=%q, indent=%q) for %T: %w",
			prefix, indent, v, err,
		)
	}

	return data, nil
}

// UnmarshalFormat decodes data into v using the provided unmarshal function.
func UnmarshalFormat(format string, unmarshalFn func([]byte, any) error, data []byte, v any) error {
	err := unmarshalFn(data, v)
	if err != nil {
		return fmt.Errorf("unmarshal %s into %T: %w", format, v, err)
	}

	return nil
}

// MarshalFormat encodes v using the provided marshal function.
func MarshalFormat(format string, marshalFn func(any) ([]byte, error), v any) ([]byte, error) {
	data, err := marshalFn(v)
	if err != nil {
		return nil, fmt.Errorf("marshal %s %T: %w", format, v, err)
	}

	return data, nil
}

// BrandedValue returns the string value of a branded ID, or empty string if zero.
func BrandedValue[Brand any](label id.ID[Brand, string]) string {
	if label.IsZero() {
		return ""
	}

	return label.Get()
}
