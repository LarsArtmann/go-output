package testhelpers

import (
	"errors"
	"io"
)

// ErrWrite is a sentinel error used by ErrorWriter and WriteNThenFailWriter.
var ErrWrite = errors.New("write error")

// ErrorWriter implements io.Writer, always returning an error.
type ErrorWriter struct{}

func (e *ErrorWriter) Write(_ []byte) (int, error) {
	return 0, ErrWrite
}

var _ io.Writer = (*ErrorWriter)(nil)

// WriteNThenFailWriter implements io.Writer. It succeeds for the first N calls,
// then returns ErrWrite on every subsequent call.
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
