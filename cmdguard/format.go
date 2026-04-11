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
// Deprecated: Use NewEnumFlag[output.Format] with enumFlagParams instead.
func NewOutputFormatFlag(value *output.Format) *OutputFormatFlag {
	return NewEnumFlag(enumFlagParams[output.Format]{
		Value:     value,
		Name:      "output format",
		ParseFunc: output.ParseFormat,
	})
}
