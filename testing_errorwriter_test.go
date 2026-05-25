package output

import (
	"errors"

	"github.com/larsartmann/go-output/testhelpers"
)

type errorWriter = testhelpers.ErrorWriter

var errWrite = errors.New("write error")

type writeNThenFailWriter = testhelpers.WriteNThenFailWriter
