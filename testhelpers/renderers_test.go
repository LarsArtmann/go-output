package testhelpers

import (
	"errors"
	"testing"
)

func TestErrorRenderer(t *testing.T) {
	t.Parallel()

	r := &ErrorRenderer{}
	got, err := r.Render()

	if got != "" {
		t.Errorf("ErrorRenderer.Render() got = %q, want empty", got)
	}

	if !errors.Is(err, ErrTest) {
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

func TestMustRender(t *testing.T) {
	t.Parallel()

	t.Run("returns string on success", func(t *testing.T) {
		t.Parallel()

		if got := MustRender(&FixedRenderer{Output: "ok"}); got != "ok" {
			t.Errorf("MustRender() = %q, want %q", got, "ok")
		}
	})

	t.Run("panics on error", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("MustRender with ErrorRenderer should panic")
			}
		}()

		MustRender(&ErrorRenderer{})
	})
}
