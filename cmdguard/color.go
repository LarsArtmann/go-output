package cmdguard

import (
	"github.com/larsartmann/go-output"
)

type ColorModeFlag struct {
	value *output.ColorMode
}

func NewColorModeFlag(value *output.ColorMode) *ColorModeFlag {
	return &ColorModeFlag{value: value}
}

func (f *ColorModeFlag) Parse(s string) error {
	parsed, err := output.ParseColorMode(s)
	if err != nil {
		return err
	}
	*f.value = parsed
	return nil
}

func (f *ColorModeFlag) AllowedValues() []string {
	return f.value.AllowedValues()
}

func (f *ColorModeFlag) Default() string {
	return f.value.String()
}
