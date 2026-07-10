package output

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
)

// MarshalJSONIndent encodes v to indented JSON.
func MarshalJSONIndent(v any, prefix, indent string) ([]byte, error) {
	data, err := json.Marshal(
		v,
		json.Deterministic(true),
		jsontext.WithIndentPrefix(prefix),
		jsontext.WithIndent(indent),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal json indent (prefix=%q, indent=%q) for %T: %w",
			prefix, indent, v, err,
		)
	}

	return data, nil
}
