package output

import (
	"strings"
	"testing"
)

func TestTableStore_SetHeaders(t *testing.T) {
	t.Parallel()

	var base TableStore
	base.SetHeaders([]string{"A", "B"})

	data := base.Data()
	if data == nil {
		t.Fatal("Data() should not be nil after SetHeaders")
	}

	if len(data.Headers) != 2 {
		t.Errorf("Headers = %v, want 2 elements", data.Headers)
	}
}

func TestTableStore_AddRow(t *testing.T) {
	t.Parallel()

	var base TableStore
	base.AddRow([]string{"1", "2"})

	data := base.Data()
	if len(data.Rows) != 1 {
		t.Fatalf("Rows = %d, want 1", len(data.Rows))
	}

	if data.Rows[0][0] != "1" {
		t.Errorf("Row[0][0] = %q, want %q", data.Rows[0][0], "1")
	}
}

func TestTableStore_SetData(t *testing.T) {
	t.Parallel()

	var base TableStore

	d := &Table{Headers: []string{"X"}}
	base.SetData(d)

	if base.Data() != d {
		t.Error("Data() should return the same pointer set via SetData")
	}
}

func TestTableStore_SetFooter(t *testing.T) {
	t.Parallel()

	var base TableStore
	base.SetFooter([]string{"Total", "10"})

	if !base.HasFooter() {
		t.Error("HasFooter() should return true after SetFooter")
	}

	data := base.Data()
	if data.Footer[0] != "Total" {
		t.Errorf("Footer[0] = %q, want %q", data.Footer[0], "Total")
	}
}

func TestTableStore_HasFooterNil(t *testing.T) {
	t.Parallel()

	var base TableStore
	if base.HasFooter() {
		t.Error("HasFooter() should return false when data is nil")
	}
}

func TestTableStore_DataNil(t *testing.T) {
	t.Parallel()

	var base TableStore
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
		{
			"InvalidShapeError",
			&InvalidShapeError{Value: "bad", Allowed: AllShapes},
			"invalid shape: bad (allowed: table, tree, graph)",
		},
		{"InvalidNodeShapeError", &InvalidNodeShapeError{Value: "bad"}, "invalid node shape: bad"},
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

	assertLabel(t, "Label", node.Label.Get(), "Label 1")

	if node.Shape != NodeShapeBox {
		t.Errorf("Shape = %v, want %v", node.Shape, NodeShapeBox)
	}

	if node.Metadata == nil {
		t.Error("Metadata should be initialized")
	}
}

// fixedRenderer is a test helper that returns a fixed string.
type fixedRenderer struct {
	output string
}

// assertLabel checks that got equals want, failing with a descriptive error.
func assertLabel(t *testing.T, name, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}

func (f *fixedRenderer) Render() (string, error) { return f.output, nil }

func TestRendererAsWriter(t *testing.T) {
	t.Parallel()

	sr := RendererAsWriter(&fixedRenderer{output: "hello"})

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
