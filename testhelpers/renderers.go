package testhelpers

import (
	"errors"
	"fmt"
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

// Renderer is the minimal interface shared by every go-output renderer.
// It mirrors the contract documented in the project's AGENTS.md (Render naming
// convention): any type that can be rendered to a string implements this.
type Renderer interface {
	Render() (string, error)
}

// MustRender calls r.Render() and returns the result, panicking on error.
// Intended for Example* functions and other demonstration code where the
// canonical "render or print error, otherwise print result" boilerplate
// is identical across every renderer module.
func MustRender(r Renderer) string {
	out, err := r.Render()
	if err != nil {
		panic(fmt.Sprintf("MustRender: %v", err))
	}

	return out
}
