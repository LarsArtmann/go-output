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
	result := make([]string, 0, len(values))
	for _, v := range values {
		result = append(result, toString(v))
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
	result := make([]string, 0, len(values))
	for _, v := range values {
		result = append(result, v.String())
	}

	return result
}

// ParseError represents a parse error for enum types.
//
// This is the internal error returned by ParseEnum. Domain-specific Parse
// functions (ParseShape, ParseFormat, etc.) discard it and return their own
// typed errors (InvalidShapeError, InvalidFormatError, etc.) which carry
// richer type information. Consumers should match against the domain-specific
// types, not ParseError.
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

// EnumErrorMessage formats the canonical "invalid <kind>: <value> (allowed: a, b, c)"
// message used by every module's InvalidXxxError.Error() implementation. The
// empty-allowed branch mirrors the "len(Allowed) == 0" guard the typed errors
// all carry, so callers can collapse their 4-line Error() body to a one-liner.
//
// This helper exists because the same body was duplicated in 13 sites across
// root (color/line-style/node-shape), d2/ (direction/node-shape/arrow-type/
// constraint/text-transform), graph/ (rank-dir/spline-style), and nom/
// (activity-status/activity-kind). Each module keeps its own typed error
// struct (root cannot import sub-modules per the Core Invariant); only the
// message-formatting predicate is shared.
func EnumErrorMessage(kind, value string, allowed []string) string {
	if len(allowed) == 0 {
		return "invalid " + kind + ": " + value
	}

	return "invalid " + kind + ": " + value + " (allowed: " + joinStrings(allowed) + ")"
}
