// Package cmdguard provides flag parsing utilities for CLI applications.
package cmdguard

import (
	"fmt"

	"github.com/larsartmann/go-output"
)

// ColorModeFlag parses color mode flags.
type ColorModeFlag struct {
	value *output.ColorMode
}

// NewColorModeFlag creates a new ColorModeFlag.
func NewColorModeFlag(value *output.ColorMode) *ColorModeFlag {
	return &ColorModeFlag{value: value}
}

// Parse parses the flag value.
func (f *ColorModeFlag) Parse(s string) error {
	parsed, err := output.ParseColorMode(s)
	if err != nil {
		return fmt.Errorf("parse color mode %q: %w", s, err)
	}
	*f.value = parsed
	return nil
}

// AllowedValues returns the allowed flag values.
func (f *ColorModeFlag) AllowedValues() []string {
	return f.value.AllowedValues()
}

// Default returns the default value.
func (f *ColorModeFlag) Default() string {
	return f.value.String()
}
