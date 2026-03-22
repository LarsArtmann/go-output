package cmdguard

import (
	"github.com/larsartmann/go-output"
)

type OutputFormatFlag struct {
	value *output.OutputFormat
}

func NewOutputFormatFlag(value *output.OutputFormat) *OutputFormatFlag {
	return &OutputFormatFlag{value: value}
}

func (f *OutputFormatFlag) Parse(s string) error {
	parsed, err := output.ParseOutputFormat(s)
	if err != nil {
		return err
	}
	*f.value = parsed
	return nil
}

func (f *OutputFormatFlag) AllowedValues() []string {
	return (*f.value).AllowedValues()
}

func (f *OutputFormatFlag) Default() string {
	return f.value.String()
}
