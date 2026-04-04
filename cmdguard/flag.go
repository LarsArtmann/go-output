package cmdguard

import (
	"fmt"
)

// EnumValue is the interface that enum types must implement to work with EnumFlag.
type EnumValue interface {
	~string
	String() string
	AllowedValues() []string
}

// EnumFlag is a generic flag parser for enum types.
type EnumFlag[T EnumValue] struct {
	value     *T
	name      string
	parseFunc func(string) (T, error)
}

// NewEnumFlag creates a new EnumFlag for the given enum type.
func NewEnumFlag[T EnumValue](
	value *T,
	name string,
	parseFunc func(string) (T, error),
) *EnumFlag[T] {
	return &EnumFlag[T]{
		value:     value,
		name:      name,
		parseFunc: parseFunc,
	}
}

// Parse parses the flag value.
func (f *EnumFlag[T]) Parse(s string) error {
	parsed, err := f.parseFunc(s)
	if err != nil {
		return fmt.Errorf("parse %s %q: %w", f.name, s, err)
	}

	*f.value = parsed

	return nil
}

// AllowedValues returns the allowed flag values.
func (f *EnumFlag[T]) AllowedValues() []string {
	return (*f.value).AllowedValues()
}

// Default returns the default value.
func (f *EnumFlag[T]) Default() string {
	return (*f.value).String()
}
