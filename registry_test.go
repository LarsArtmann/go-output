package output

import (
	"slices"
	"sync"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Parallel()
	testRegisterAndVerify(t, FormatJSON, "json")
}

func testRegisterAndVerify(t *testing.T, format OutputFormat, output string) {
	t.Helper()
	Unregister(format)

	err := Register(format, func() Renderer {
		return &testRenderer{output: output}
	})
	if err != nil {
		t.Fatalf("Register(%v) error = %v", format, err)
	}

	if !IsRegistered(format) {
		t.Errorf("IsRegistered(%v) = false, want true", format)
	}
}

func registerFormatForTest(format OutputFormat, output string) {
	_ = Register(format, func() Renderer {
		return &testRenderer{output: output}
	})
}

func TestRegisterDuplicate(t *testing.T) {
	t.Parallel()
	Unregister(FormatCSV)

	registerFormatForTest(FormatCSV, "csv1")

	err := Register(FormatCSV, func() Renderer {
		return &testRenderer{output: "csv2"}
	})
	if err == nil {
		t.Error("Register() expected error for duplicate registration")
	}
}

func TestUnregister(t *testing.T) {
	t.Parallel()
	Unregister(FormatYAML)
	registerFormatForTest(FormatYAML, "yaml")

	Unregister(FormatYAML)

	if IsRegistered(FormatYAML) {
		t.Error("IsRegistered(FormatYAML) = true, want false after Unregister")
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	t.Run("existing format", func(t *testing.T) {
		t.Parallel()
		Unregister(FormatD2)

		err := Register(FormatD2, func() Renderer {
			return &testRenderer{output: "d2-output"}
		})
		if err != nil {
			t.Fatalf("Register() error = %v", err)
		}

		r, err := Create(FormatD2)
		if err != nil {
			t.Fatalf("Create(FormatD2) error = %v", err)
		}

		if r.Render() != "d2-output" {
			t.Errorf("Create(FormatD2).Render() = %q, want %q", r.Render(), "d2-output")
		}
	})

	t.Run("unregistered format", func(t *testing.T) {
		t.Parallel()
		Unregister(FormatMermaid)

		r, err := Create(FormatMermaid)
		if err == nil {
			t.Error("Create(unregistered) expected error, got nil")
		}

		if r != nil {
			t.Errorf("Create(unregistered) = %v, want nil", r)
		}
	})
}

func TestRegisteredFormats(t *testing.T) {
	t.Parallel()
	Unregister(FormatTable)
	Unregister(FormatJSON)

	err := Register(FormatTable, func() Renderer { return &testRenderer{output: ""} })
	if err != nil {
		t.Fatalf("Register(FormatTable) error = %v", err)
	}

	err = Register(FormatJSON, func() Renderer { return &testRenderer{output: ""} })
	if err != nil {
		t.Fatalf("Register(FormatJSON) error = %v", err)
	}

	formats := RegisteredFormats()
	if len(formats) < 2 {
		t.Fatalf("RegisteredFormats() returned %d formats, want at least 2", len(formats))
	}

	found := slices.Contains(formats, FormatTable)
	if !found {
		t.Error("RegisteredFormats() does not contain FormatTable")
	}
}

func TestIsRegistered(t *testing.T) {
	t.Parallel()
	t.Run("registered", func(t *testing.T) {
		t.Parallel()
		testRegisterAndVerify(t, FormatHTML, "")
	})

	t.Run("unregistered", func(t *testing.T) {
		t.Parallel()
		Unregister(FormatDOT)

		if IsRegistered(FormatDOT) {
			t.Error("IsRegistered(FormatDOT) = true, want false")
		}
	})
}

func TestRegistryConcurrency(t *testing.T) {
	t.Parallel()
	Unregister(FormatCSV)
	Unregister(FormatJSON)

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 100 {
				registerFormatForTest(FormatCSV, "csv")
				Unregister(FormatCSV)
			}
		})
	}

	for range 10 {
		wg.Go(func() {
			for range 100 {
				IsRegistered(FormatCSV)
				RegisteredFormats()
			}
		})
	}

	wg.Wait()
}

type testRenderer struct {
	output string
}

func (r *testRenderer) Render() string {
	return r.output
}
