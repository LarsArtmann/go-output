package testhelpers

import (
	"errors"
	"io"
)

var ErrWrite = errors.New("write error")

type ErrorWriter struct{}

func (e *ErrorWriter) Write(_ []byte) (int, error) {
	return 0, ErrWrite
}

var _ io.Writer = (*ErrorWriter)(nil)

type WriteNThenFailWriter struct {
	Remaining int
}

func (w *WriteNThenFailWriter) Write(p []byte) (int, error) {
	if w.Remaining <= 0 {
		return 0, ErrWrite
	}

	w.Remaining--

	return len(p), nil
}

var _ io.Writer = (*WriteNThenFailWriter)(nil)
