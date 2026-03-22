package output

import "fmt"

type SortBy string

const (
	SortByName       SortBy = "name"
	SortByImportance SortBy = "importance"
	SortByCreatedAt  SortBy = "created_at"
	SortByUpdatedAt  SortBy = "updated_at"
	SortByHealth     SortBy = "health"
	SortByComplexity SortBy = "complexity"
)

var sortByValues = []SortBy{
	SortByName,
	SortByImportance,
	SortByCreatedAt,
	SortByUpdatedAt,
	SortByHealth,
	SortByComplexity,
}

func ParseSortBy(s string) (SortBy, error) {
	for _, v := range sortByValues {
		if string(v) == s {
			return v, nil
		}
	}
	return "", fmt.Errorf("invalid sort by: %q (allowed: %v)", s, sortByValues)
}

func (s SortBy) String() string {
	return string(s)
}

func (s SortBy) AllowedValues() []string {
	values := make([]string, len(sortByValues))
	for i, v := range sortByValues {
		values[i] = string(v)
	}
	return values
}

func (s SortBy) IsValid() bool {
	for _, v := range sortByValues {
		if s == v {
			return true
		}
	}
	return false
}
