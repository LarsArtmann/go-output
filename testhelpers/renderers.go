package testhelpers

import (
	"errors"
)

// ErrTest is a sentinel error used by ErrorRenderer.
var ErrTest = errors.New("test error")

// ErrorRenderer implements output.Renderer, always returning an error.
type ErrorRenderer struct{}

func (e *ErrorRenderer) Render() (string, error) {
	return "", ErrTest
}

// FixedRenderer implements output.Renderer, returning a fixed string.
type FixedRenderer struct {
	Output string
}

func (r *FixedRenderer) Render() (string, error) {
	return r.Output, nil
}
