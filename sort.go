// Package output provides consistent output formatting for CLI applications.
package output

import (
	"fmt"

	"github.com/larsartmann/go-output/enum"
)

// SortBy represents the available sort field options for CLI applications.
type SortBy string

// Sort field constants.
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

// ParseSortBy converts a string to SortBy, returning an error if invalid.
func ParseSortBy(s string) (SortBy, error) {
	v, err := enum.Parse(sortByValues, s, func(s SortBy) string { return string(s) })
	if err != nil {
		return "", fmt.Errorf("invalid sort by: %q", s)
	}

	return v, nil
}

// String returns the string representation of the sort field.
func (s SortBy) String() string {
	return string(s)
}

// AllowedValues returns all valid sort field values for CLI help text.
func (s SortBy) AllowedValues() []string {
	return enum.AllowedValues(sortByValues)
}

// IsValid returns true if the sort field is a valid SortBy value.
func (s SortBy) IsValid() bool {
	return enum.Contains(sortByValues, s)
}
