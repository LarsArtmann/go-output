package testhelpers

import (
	"testing"
)

func TestErrorRenderer(t *testing.T) {
	t.Parallel()

	r := &ErrorRenderer{}
	got, err := r.Render()

	if got != "" {
		t.Errorf("ErrorRenderer.Render() got = %q, want empty", got)
	}

	if err != ErrTest {
		t.Errorf("ErrorRenderer.Render() err = %v, want ErrTest", err)
	}
}

func TestFixedRenderer(t *testing.T) {
	t.Parallel()

	r := &FixedRenderer{Output: "hello"}
	got, err := r.Render()

	if err != nil {
		t.Errorf("FixedRenderer.Render() err = %v, want nil", err)
	}

	if got != "hello" {
		t.Errorf("FixedRenderer.Render() = %q, want %q", got, "hello")
	}
}
