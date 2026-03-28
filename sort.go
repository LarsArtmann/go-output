// Package output provides consistent output formatting for CLI applications.
package output

import (
	"fmt"

	"github.com/larsartmann/go-output/enum"
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
	v, err := enum.Parse(sortByValues, s, func(s SortBy) string { return string(s) })
	if err != nil {
		return "", fmt.Errorf("invalid sort by: %q", s)
	}
	return v, nil
}

func (s SortBy) String() string {
	return string(s)
}

// AllowedValues returns all valid sort field values.
func (s SortBy) AllowedValues() []string {
	return enum.AllowedStrings(sortByValues, func(s SortBy) string { return string(s) })
}

// IsValid checks if the sort field is valid.
func (s SortBy) IsValid() bool {
	return enum.Contains(sortByValues, s)
}
