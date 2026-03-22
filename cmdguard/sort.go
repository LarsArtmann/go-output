package cmdguard

import (
	"github.com/larsartmann/go-output"
)

type SortByFlag struct {
	value *output.SortBy
}

func NewSortByFlag(value *output.SortBy) *SortByFlag {
	return &SortByFlag{value: value}
}

func (f *SortByFlag) Parse(s string) error {
	parsed, err := output.ParseSortBy(s)
	if err != nil {
		return err
	}
	*f.value = parsed
	return nil
}

func (f *SortByFlag) AllowedValues() []string {
	return (*f.value).AllowedValues()
}

func (f *SortByFlag) Default() string {
	return f.value.String()
}
