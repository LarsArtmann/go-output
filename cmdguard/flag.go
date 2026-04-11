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

// enumFlagParams holds the parameters for creating an EnumFlag.
type enumFlagParams[T EnumValue] struct {
	Value     *T
	Name      string
	ParseFunc func(string) (T, error)
}

// EnumFlag is a generic flag parser for enum types.
type EnumFlag[T EnumValue] struct {
	enumFlagParams[T]
}

// NewEnumFlag creates a new EnumFlag for the given enum type.
func NewEnumFlag[T EnumValue](params enumFlagParams[T]) *EnumFlag[T] {
	return &EnumFlag[T]{
		enumFlagParams: params,
	}
}

// Parse parses the flag value.
func (f *EnumFlag[T]) Parse(s string) error {
	parsed, err := f.ParseFunc(s)
	if err != nil {
		return fmt.Errorf("parse %s %q: %w", f.Name, s, err)
	}

	*f.Value = parsed

	return nil
}

// AllowedValues returns the allowed flag values.
func (f *EnumFlag[T]) AllowedValues() []string {
	return (*f.Value).AllowedValues()
}

// Default returns the default value.
func (f *EnumFlag[T]) Default() string {
	return (*f.Value).String()
}
