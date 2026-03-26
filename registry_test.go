package output

import (
	"slices"
	"sync"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Parallel()
	Unregister(FormatJSON)

	err := Register(FormatJSON, func() Renderer {
		return &testRenderer{output: "json"}
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if !IsRegistered(FormatJSON) {
		t.Error("IsRegistered(FormatJSON) = false, want true")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	t.Parallel()
	Unregister(FormatCSV)

	_ = Register(FormatCSV, func() Renderer {
		return &testRenderer{output: "csv1"}
	})

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
	_ = Register(FormatYAML, func() Renderer {
		return &testRenderer{output: "yaml"}
	})

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
		Unregister(FormatHTML)
		err := Register(FormatHTML, func() Renderer { return &testRenderer{output: ""} })
		if err != nil {
			t.Fatalf("Register(FormatHTML) error = %v", err)
		}
		if !IsRegistered(FormatHTML) {
			t.Error("IsRegistered(FormatHTML) = false, want true")
		}
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
	Unregister(FormatCSV)
	Unregister(FormatJSON)

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 100 {
				_ = Register(FormatCSV, func() Renderer { return &testRenderer{output: "csv"} })
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
