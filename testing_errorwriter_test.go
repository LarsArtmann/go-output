package output

import (
	"errors"
	"io"
)

type errorWriter struct{}

var errWrite = errors.New("write error")

func (e *errorWriter) Write(_ []byte) (int, error) {
	return 0, errWrite
}

var _ io.Writer = (*errorWriter)(nil)

type writeNThenFailWriter struct {
	remaining int
}

func (w *writeNThenFailWriter) Write(p []byte) (int, error) {
	if w.remaining <= 0 {
		return 0, errWrite
	}

	w.remaining--

	return len(p), nil
}

var _ io.Writer = (*writeNThenFailWriter)(nil)
