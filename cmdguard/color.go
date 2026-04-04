// Package cmdguard provides flag parsing utilities for CLI applications.
package cmdguard

import (
	"github.com/larsartmann/go-output"
)

// ColorModeFlag is a flag parser for ColorMode values.
//
// Deprecated: Use NewEnumFlag with output.ColorMode instead.
type ColorModeFlag = EnumFlag[output.ColorMode]

// NewColorModeFlag creates a new ColorModeFlag.
//
// Deprecated: Use NewEnumFlag[output.ColorMode](value, "color mode", output.ParseColorMode) instead.
func NewColorModeFlag(value *output.ColorMode) *ColorModeFlag {
	return NewEnumFlag(value, "color mode", output.ParseColorMode)
}
