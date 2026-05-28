package output

import (
	"errors"
	"strings"
	"testing"
)

func TestTableDataBase_SetHeaders(t *testing.T) {
	t.Parallel()

	var base TableDataBase
	base.SetHeaders([]string{"A", "B"})

	data := base.Data()
	if data == nil {
		t.Fatal("Data() should not be nil after SetHeaders")
	}

	if len(data.Headers) != 2 {
		t.Errorf("Headers = %v, want 2 elements", data.Headers)
	}
}

func TestTableDataBase_AddRow(t *testing.T) {
	t.Parallel()

	var base TableDataBase
	base.AddRow([]string{"1", "2"})

	data := base.Data()
	if len(data.Rows) != 1 {
		t.Fatalf("Rows = %d, want 1", len(data.Rows))
	}

	if data.Rows[0][0] != "1" {
		t.Errorf("Row[0][0] = %q, want %q", data.Rows[0][0], "1")
	}
}

func TestTableDataBase_SetData(t *testing.T) {
	t.Parallel()

	var base TableDataBase
	d := &TableData{Headers: []string{"X"}}
	base.SetData(d)

	if base.Data() != d {
		t.Error("Data() should return the same pointer set via SetData")
	}
}

func TestTableDataBase_SetFooter(t *testing.T) {
	t.Parallel()

	var base TableDataBase
	base.SetFooter([]string{"Total", "10"})

	if !base.HasFooter() {
		t.Error("HasFooter() should return true after SetFooter")
	}

	data := base.Data()
	if data.Footer[0] != "Total" {
		t.Errorf("Footer[0] = %q, want %q", data.Footer[0], "Total")
	}
}

func TestTableDataBase_HasFooterNil(t *testing.T) {
	t.Parallel()

	var base TableDataBase
	if base.HasFooter() {
		t.Error("HasFooter() should return false when data is nil")
	}
}

func TestTableDataBase_DataNil(t *testing.T) {
	t.Parallel()

	var base TableDataBase
	if base.Data() != nil {
		t.Error("Data() should return nil when not initialized")
	}
}

func TestErrorTypes(t *testing.T) {
	t.Parallel()

	t.Run("InvalidColorModeError", func(t *testing.T) {
		t.Parallel()

		err := &InvalidColorModeError{Value: "bad"}
		if err.Error() != "invalid color mode: bad" {
			t.Errorf("Error() = %q, want %q", err.Error(), "invalid color mode: bad")
		}
	})

	t.Run("InvalidShapeError", func(t *testing.T) {
		t.Parallel()

		err := &InvalidShapeError{Value: "bad"}
		if err.Error() != "invalid shape: bad" {
			t.Errorf("Error() = %q, want %q", err.Error(), "invalid shape: bad")
		}
	})

	t.Run("InvalidGraphShapeError", func(t *testing.T) {
		t.Parallel()

		err := &InvalidGraphShapeError{Value: "bad"}
		if err.Error() != "invalid graph shape: bad" {
			t.Errorf("Error() = %q, want %q", err.Error(), "invalid graph shape: bad")
		}
	})
}

func TestFormatIsValid(t *testing.T) {
	t.Parallel()

	if !FormatCSV.IsValid() {
		t.Error("FormatCSV.IsValid() should return true")
	}

	if Format("nonexistent").IsValid() {
		t.Error("invalid format IsValid() should return false")
	}
}

func TestNewGraphNode(t *testing.T) {
	t.Parallel()

	node := NewGraphNode("id1", "Label 1")
	if node.ID.Get() != "id1" {
		t.Errorf("ID = %q, want %q", node.ID.Get(), "id1")
	}

	if node.Label.Get() != "Label 1" {
		t.Errorf("Label = %q, want %q", node.Label.Get(), "Label 1")
	}

	if node.Shape != ShapeBox {
		t.Errorf("Shape = %v, want %v", node.Shape, ShapeBox)
	}

	if node.Metadata == nil {
		t.Error("Metadata should be initialized")
	}
}

// fixedRenderer is a test helper that returns a fixed string.
type fixedRenderer struct {
	output string
}

func (f *fixedRenderer) Render() (string, error) { return f.output, nil }

func TestStreamingRendererFromRenderer(t *testing.T) {
	t.Parallel()

	sr := StreamingRendererFromRenderer(&fixedRenderer{output: "hello"})

	got, err := sr.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != "hello" {
		t.Errorf("Render() = %q, want %q", got, "hello")
	}

	var buf strings.Builder

	err = sr.Stream(&buf)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	if buf.String() != "hello" {
		t.Errorf("Stream() = %q, want %q", buf.String(), "hello")
	}
}

func TestUnmarshalFormat(t *testing.T) {
	t.Parallel()

	var result string

	err := UnmarshalFormat("json", func(data []byte, v any) error {
		return errors.New("test error") //nolint:goerr113 // intentional test error
	}, []byte(`"hello"`), &result)
	if err == nil {
		t.Fatal("UnmarshalFormat should propagate unmarshal error")
	}
}
