package cmdguard

import (
	"github.com/larsartmann/go-output"
)

// OutputFormatFlag is a flag parser for Format values.
//
// Deprecated: Use NewEnumFlag with output.Format instead.
type OutputFormatFlag = EnumFlag[output.Format]

// NewOutputFormatFlag creates a new OutputFormatFlag.
//
// Deprecated: Use NewEnumFlag[output.Format](value, "output format", output.ParseFormat) instead.
func NewOutputFormatFlag(value *output.Format) *OutputFormatFlag {
	return NewEnumFlag(value, "output format", output.ParseFormat)
}
