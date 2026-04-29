package output

import (
	"slices"
	"sync"
	"testing"
)

// testRegistryFormat avoids collisions with real format constants in parallel tests.
const testRegistryFormat Format = "__test_registry__"

const testOutput = "test-output"

func TestRegister(t *testing.T) {
	Unregister(testRegistryFormat)

	err := Register(testRegistryFormat, testRendererFunc("test"))
	if err != nil {
		t.Fatalf("Register(%v) error = %v", testRegistryFormat, err)
	}

	if !IsRegistered(testRegistryFormat) {
		t.Errorf("IsRegistered(%v) = false, want true", testRegistryFormat)
	}
}

// testRendererFunc creates a test renderer for the given output string.
func testRendererFunc(output string) func() Renderer {
	return func() Renderer { return &testRenderer{output: output} }
}

func TestRegisterDuplicate(t *testing.T) {
	format := Format("__test_dup__")
	Unregister(format)

	_ = Register(format, testRendererFunc("1"))

	err := Register(format, testRendererFunc("2"))
	if err == nil {
		t.Error("Register() expected error for duplicate registration")
	}
}

func TestUnregister(t *testing.T) {
	format := Format("__test_unreg__")
	Unregister(format)
	_ = Register(format, testRendererFunc("yaml"))

	Unregister(format)

	if IsRegistered(format) {
		t.Errorf("IsRegistered(%v) = true, want false after Unregister", format)
	}
}

func TestCreate(t *testing.T) {
	t.Run("existing format", func(t *testing.T) {
		format := Format("__test_create__")
		Unregister(format)

		err := Register(format, testRendererFunc(testOutput))
		if err != nil {
			t.Fatalf("Register() error = %v", err)
		}

		r, err := Create(format)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if got != testOutput {
			t.Errorf("Create().Render() = %q, want %q", got, testOutput)
		}
	})

	t.Run("unregistered format", func(t *testing.T) {
		format := Format("__test_missing__")
		Unregister(format)

		r, err := Create(format)
		if err == nil {
			t.Error("Create(unregistered) expected error, got nil")
		}

		if r != nil {
			t.Errorf("Create(unregistered) = %v, want nil", r)
		}
	})
}

func TestRegisteredFormats(t *testing.T) {
	formatA := Format("__test_fmt_a__")
	formatB := Format("__test_fmt_b__")

	Unregister(formatA)
	Unregister(formatB)

	err := Register(formatA, func() Renderer { return &testRenderer{output: ""} })
	if err != nil {
		t.Fatalf("Register(A) error = %v", err)
	}

	err = Register(formatB, func() Renderer { return &testRenderer{output: ""} })
	if err != nil {
		t.Fatalf("Register(B) error = %v", err)
	}

	formats := RegisteredFormats()
	if len(formats) < 2 {
		t.Fatalf("RegisteredFormats() returned %d formats, want at least 2", len(formats))
	}

	if !slices.Contains(formats, formatA) {
		t.Error("RegisteredFormats() does not contain format A")
	}
}

func TestIsRegistered(t *testing.T) {
	t.Run("registered", func(t *testing.T) {
		format := Format("__test_isreg__")
		Unregister(format)

		err := Register(format, testRendererFunc(""))
		if err != nil {
			t.Fatalf("Register() error = %v", err)
		}

		if !IsRegistered(format) {
			t.Errorf("IsRegistered(%v) = false, want true", format)
		}
	})

	t.Run("unregistered", func(t *testing.T) {
		format := Format("__test_unreg2__")
		Unregister(format)

		if IsRegistered(format) {
			t.Errorf("IsRegistered(%v) = true, want false", format)
		}
	})
}

func TestRegistryConcurrency(t *testing.T) {
	format := Format("__test_concurrent__")
	Unregister(format)

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 100 {
				_ = Register(format, testRendererFunc("x"))
				Unregister(format)
			}
		})
	}

	for range 10 {
		wg.Go(func() {
			for range 100 {
				IsRegistered(format)
				RegisteredFormats()
			}
		})
	}

	wg.Wait()
}

type testRenderer struct {
	output string
}

func (r *testRenderer) Render() (string, error) {
	return r.output, nil
}
