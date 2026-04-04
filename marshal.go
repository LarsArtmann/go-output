package output

import "fmt"

// unmarshal decodes data into v using the provided unmarshal function.
func unmarshal(format string, unmarshalFn func([]byte, any) error, data []byte, v any) error {
	err := unmarshalFn(data, v)
	if err != nil {
		return fmt.Errorf("unmarshal %s into %T: %w", format, v, err)
	}

	return nil
}

// marshal encodes v using the provided marshal function.
func marshal(format string, marshalFn func(any) ([]byte, error), v any) ([]byte, error) {
	data, err := marshalFn(v)
	if err != nil {
		return nil, fmt.Errorf("marshal %s (%T): %w", format, v, err)
	}

	return data, nil
}

// marshalIndent encodes v with indentation using the provided marshal function.
func marshalIndent(
	format string,
	marshalFn func(any, string, string) ([]byte, error),
	v any,
	prefix, indent string,
) ([]byte, error) {
	data, err := marshalFn(v, prefix, indent)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal %s indent (prefix=%q, indent=%q) for %T: %w",
			format,
			prefix,
			indent,
			v,
			err,
		)
	}

	return data, nil
}
