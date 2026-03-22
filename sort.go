// Package output provides consistent output formatting for CLI applications.
package output

import (
	"fmt"
	"slices"
)

// SortBy specifies the field to sort by.
type SortBy string

// Sort field options.
const (
	SortByName       SortBy = "name"
	SortByImportance SortBy = "importance"
	SortByCreatedAt  SortBy = "created_at"
	SortByUpdatedAt  SortBy = "updated_at"
	SortByHealth     SortBy = "health"
	SortByComplexity SortBy = "complexity"
)

//nolint:gochecknoglobals // Global variable used for value iteration.
var sortByValues = []SortBy{
	SortByName,
	SortByImportance,
	SortByCreatedAt,
	SortByUpdatedAt,
	SortByHealth,
	SortByComplexity,
}

// ParseSortBy parses a sort field string.
func ParseSortBy(s string) (SortBy, error) {
	if slices.Contains(sortByValues, SortBy(s)) {
		return SortBy(s), nil
	}
	return "", fmt.Errorf("invalid sort by: %q (allowed: %v)", s, sortByValues)
}

func (s SortBy) String() string {
	return string(s)
}

// AllowedValues returns all valid sort field values.
func (s SortBy) AllowedValues() []string {
	values := make([]string, len(sortByValues))
	for i, v := range sortByValues {
		values[i] = string(v)
	}
	return values
}

// IsValid checks if the sort field is valid.
func (s SortBy) IsValid() bool {
	return slices.Contains(sortByValues, s)
}
