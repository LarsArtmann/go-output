package output

import (
	"strings"
	"testing"
)

func TestTableDataStore_SetHeaders(t *testing.T) {
	t.Parallel()

	var base TableDataStore
	base.SetHeaders([]string{"A", "B"})

	data := base.Data()
	if data == nil {
		t.Fatal("Data() should not be nil after SetHeaders")
	}

	if len(data.Headers) != 2 {
		t.Errorf("Headers = %v, want 2 elements", data.Headers)
	}
}

func TestTableDataStore_AddRow(t *testing.T) {
	t.Parallel()

	var base TableDataStore
	base.AddRow([]string{"1", "2"})

	data := base.Data()
	if len(data.Rows) != 1 {
		t.Fatalf("Rows = %d, want 1", len(data.Rows))
	}

	if data.Rows[0][0] != "1" {
		t.Errorf("Row[0][0] = %q, want %q", data.Rows[0][0], "1")
	}
}

func TestTableDataStore_SetData(t *testing.T) {
	t.Parallel()

	var base TableDataStore

	d := &TableData{Headers: []string{"X"}}
	base.SetData(d)

	if base.Data() != d {
		t.Error("Data() should return the same pointer set via SetData")
	}
}

func TestTableDataStore_SetFooter(t *testing.T) {
	t.Parallel()

	var base TableDataStore
	base.SetFooter([]string{"Total", "10"})

	if !base.HasFooter() {
		t.Error("HasFooter() should return true after SetFooter")
	}

	data := base.Data()
	if data.Footer[0] != "Total" {
		t.Errorf("Footer[0] = %q, want %q", data.Footer[0], "Total")
	}
}

func TestTableDataStore_HasFooterNil(t *testing.T) {
	t.Parallel()

	var base TableDataStore
	if base.HasFooter() {
		t.Error("HasFooter() should return false when data is nil")
	}
}

func TestTableDataStore_DataNil(t *testing.T) {
	t.Parallel()

	var base TableDataStore
	if base.Data() != nil {
		t.Error("Data() should return nil when not initialized")
	}
}

func TestErrorTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want string
	}{
		{"InvalidColorModeError", &InvalidColorModeError{Value: "bad"}, "invalid color mode: bad"},
		{"InvalidShapeError", &InvalidShapeError{Value: "bad"}, "invalid shape: bad"},
		{"InvalidGraphShapeError", &InvalidGraphShapeError{Value: "bad"}, "invalid graph shape: bad"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.err.Error() != tt.want {
				t.Errorf("Error() = %q, want %q", tt.err.Error(), tt.want)
			}
		})
	}
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
