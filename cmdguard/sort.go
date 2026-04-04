package cmdguard

import (
	"github.com/larsartmann/go-output"
)

// SortByFlag is a flag parser for SortBy values.
//
// Deprecated: Use NewEnumFlag with output.SortBy instead.
type SortByFlag = EnumFlag[output.SortBy]

// NewSortByFlag creates a new SortByFlag.
//
// Deprecated: Use NewEnumFlag[output.SortBy](value, "sort by", output.ParseSortBy) instead.
func NewSortByFlag(value *output.SortBy) *SortByFlag {
	return NewEnumFlag(value, "sort by", output.ParseSortBy)
}
