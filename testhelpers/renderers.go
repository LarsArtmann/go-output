package testhelpers

import (
	"errors"
)

var ErrTest = errors.New("test error")

type ErrorRenderer struct{}

func (e *ErrorRenderer) Render() (string, error) {
	return "", ErrTest
}

type FixedRenderer struct {
	Output string
}

func (r *FixedRenderer) Render() (string, error) {
	return r.Output, nil
}
