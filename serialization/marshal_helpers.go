package serialization

import "fmt"

// stringFromBytes wraps the standard "marshal bytes to string with an error
// context" idiom that every tree/graph renderer uses. It removes the
// four-line copy/paste that toml/yaml/json renderers all carried.
//
// Callers pass the raw multi-value result of a Marshal function directly:
//
//	return stringFromBytes("toml", "tree", toml.Marshal(node))
//
// The helper takes both []byte and error positions of the marshal call's
// return signature, matching the function-style marshal APIs in this module.
func stringFromBytes(format, subject string, data []byte, err error) (string, error) {
	if err != nil {
		return "", fmt.Errorf("marshal %s %s: %w", format, subject, err)
	}

	return string(data), nil
}
