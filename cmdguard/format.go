package cmdguard

import (
	"fmt"

	"github.com/larsartmann/go-output"
)

// OutputFormatFlag parses output format flags.
type OutputFormatFlag struct {
	value *output.OutputFormat
}

// NewOutputFormatFlag creates a new OutputFormatFlag.
func NewOutputFormatFlag(value *output.OutputFormat) *OutputFormatFlag {
	return &OutputFormatFlag{value: value}
}

// Parse parses the flag value.
func (f *OutputFormatFlag) Parse(s string) error {
	parsed, err := output.ParseOutputFormat(s)
	if err != nil {
		return fmt.Errorf("parse output format: %w", err)
	}
	*f.value = parsed
	return nil
}

// AllowedValues returns the allowed flag values.
func (f *OutputFormatFlag) AllowedValues() []string {
	return f.value.AllowedValues()
}

// Default returns the default value.
func (f *OutputFormatFlag) Default() string {
	return f.value.String()
}
