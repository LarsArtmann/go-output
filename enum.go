package output

import (
	"fmt"
	"slices"
	"strings"
)

// ParseEnum parses a string into an enum value, returning an error if invalid.
func ParseEnum[T comparable](values []T, s string, toString func(T) string) (T, error) {
	for _, v := range values {
		if toString(v) == s {
			return v, nil
		}
	}

	return *new(T), &ParseError{Value: s, Values: EnumAllowedStrings(values, toString)}
}

// ContainsEnum checks if a value is in the list of allowed values.
func ContainsEnum[T comparable](values []T, v T) bool {
	return slices.Contains(values, v)
}

// EnumAllowedStrings returns the string representations of all enum values.
func EnumAllowedStrings[T any](values []T, toString func(T) string) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = toString(v)
	}

	return result
}

// StringEnum is an interface for string-based enum types.
type StringEnum interface {
	String() string
}

// EnumAllowedValues returns the string representations of all enum values
// for types that implement the StringEnum interface.
func EnumAllowedValues[T StringEnum](values []T) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.String()
	}

	return result
}

// ParseError represents a parse error for enum types.
type ParseError struct {
	Value  string
	Values []string
}

// Error returns a descriptive error message including the invalid value and allowed values.
func (e *ParseError) Error() string {
	return fmt.Sprintf("invalid value: %q (allowed: %s)", e.Value, joinStrings(e.Values))
}

func joinStrings(ss []string) string {
	return strings.Join(ss, ", ")
}
