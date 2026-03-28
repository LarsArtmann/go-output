// Package enum provides utilities for type-safe string-based enums.
package enum

import (
	"fmt"
	"slices"
)

// Parse parses a string into an enum value, returning an error if invalid.
func Parse[T comparable](values []T, s string, toString func(T) string) (T, error) {
	for _, v := range values {
		if toString(v) == s {
			return v, nil
		}
	}
	return *new(T), &ParseError[T]{Value: s, Values: values}
}

// Contains checks if a value is in the list of allowed values.
func Contains[T comparable](values []T, v T) bool {
	return slices.Contains(values, v)
}

// AllowedStrings returns the string representations of all enum values.
func AllowedStrings[T any](values []T, toString func(T) string) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = toString(v)
	}
	return result
}

// ParseError represents a parse error for enum types.
type ParseError[T any] struct {
	Value  string
	Values []T
}

func (e *ParseError[T]) Error() string {
	return fmt.Sprintf("invalid value: %q", e.Value)
}
