package output

import (
	"encoding/json"
	"fmt"
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
