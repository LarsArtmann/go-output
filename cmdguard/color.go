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
// Deprecated: Use NewEnumFlag with enumFlagParams[output.ColorMode] instead.
func NewColorModeFlag(value *output.ColorMode) *ColorModeFlag {
	return NewEnumFlag(enumFlagParams[output.ColorMode]{
		Value:     value,
		Name:      "color mode",
		ParseFunc: output.ParseColorMode,
	})
}
