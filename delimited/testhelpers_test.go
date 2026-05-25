package delimited

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type errorWriter struct{}

var errWrite = errors.New("write error")

func (e *errorWriter) Write(_ []byte) (int, error) {
	return 0, errWrite
}

var _ io.Writer = (*errorWriter)(nil)

func assertContains(t *testing.T, s, substr, msg string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("%s: %q does not contain %q", msg, s, substr)
	}
}
