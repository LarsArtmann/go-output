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

	return *new(T), &ParseError{Value: s, Values: AllowedStrings(values, toString)}
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

// Enum is an interface for string-based enum types.
type Enum interface {
	String() string
}

// AllowedValues returns the string representations of all enum values
// for types that implement the Enum interface.
func AllowedValues[T Enum](values []T) []string {
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

func (e *ParseError) Error() string {
	return fmt.Sprintf("invalid value: %q (allowed: %s)", e.Value, joinStrings(e.Values))
}

func joinStrings(ss []string) string {
	if len(ss) == 0 {
		return ""
	}

	if len(ss) == 1 {
		return ss[0]
	}

	return fmt.Sprintf("%s, %s", joinStrings(ss[:len(ss)-1]), ss[len(ss)-1])
}
