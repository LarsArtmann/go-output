package cmdguard

import (
	"fmt"

	"github.com/larsartmann/go-output"
)

// SortByFlag parses sort by flags.
type SortByFlag struct {
	value *output.SortBy
}

// NewSortByFlag creates a new SortByFlag.
func NewSortByFlag(value *output.SortBy) *SortByFlag {
	return &SortByFlag{value: value}
}

// Parse parses the flag value.
func (f *SortByFlag) Parse(s string) error {
	parsed, err := output.ParseSortBy(s)
	if err != nil {
		return fmt.Errorf("parse sort by %q: %w", s, err)
	}
	*f.value = parsed
	return nil
}

// AllowedValues returns the allowed flag values.
func (f *SortByFlag) AllowedValues() []string {
	return f.value.AllowedValues()
}

// Default returns the default value.
func (f *SortByFlag) Default() string {
	return f.value.String()
}
