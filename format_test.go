package output

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/larsartmann/go-output/internal/gentest"
)

// errTest is a static test error for MustRender panic testing.
var errTest = errors.New("test error")

func TestFormatSupports(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format Format
		shape  Shape
		want   bool
	}{
		{FormatJSON, ShapeTable, true},
		{FormatJSON, ShapeTree, true},
		{FormatJSON, ShapeGraph, true},
		{FormatCSV, ShapeTable, true},
		{FormatCSV, ShapeTree, false},
		{FormatCSV, ShapeGraph, false},
		{FormatD2, ShapeTable, true},
		{FormatD2, ShapeTree, false},
		{FormatD2, ShapeGraph, true},
		{FormatTree, ShapeTree, true},
		{FormatTree, ShapeTable, false},
		{FormatTable, ShapeTable, true},
		{FormatTable, ShapeTree, false},
		{Format("invalid"), ShapeTable, false},
	}
	for _, tt := range tests {
		t.Run(string(tt.format)+"_"+string(tt.shape), func(t *testing.T) {
			t.Parallel()

			got := tt.format.Supports(tt.shape)
			if got != tt.want {
				t.Errorf("Format(%q).Supports(%v) = %v, want %v", tt.format, tt.shape, got, tt.want)
			}
		})
	}
}

func TestFormatShapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format Format
		want   []Shape
	}{
		{FormatJSON, []Shape{ShapeTable, ShapeTree, ShapeGraph}},
		{FormatCSV, []Shape{ShapeTable}},
		{FormatD2, []Shape{ShapeTable, ShapeGraph}},
		{FormatTree, []Shape{ShapeTree}},
		{FormatTable, []Shape{ShapeTable}},
		{FormatYAML, []Shape{ShapeTable, ShapeTree, ShapeGraph}},
	}
	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			t.Parallel()

			got := tt.format.Shapes()
			if len(got) != len(tt.want) {
				t.Errorf("Format(%q).Shapes() = %v, want %v", tt.format, got, tt.want)

				return
			}

			for i, s := range got {
				if s != tt.want[i] {
					t.Errorf("Format(%q).Shapes()[%d] = %v, want %v", tt.format, i, s, tt.want[i])
				}
			}
		})
	}
}

func TestFormatsForShape(t *testing.T) {
	t.Parallel()

	assertContainsAll(t, "graph", FormatsForShape(ShapeGraph),
		FormatJSON, FormatYAML, FormatD2, FormatMermaid, FormatDOT)
	assertContainsAll(t, "tree", FormatsForShape(ShapeTree),
		FormatJSON, FormatYAML, FormatHTML, FormatTree)
	assertContainsAll(t, "table", FormatsForShape(ShapeTable),
		FormatTable, FormatJSON, FormatCSV, FormatD2, FormatMermaid, FormatDOT)
}

func assertContainsAll(t *testing.T, label string, formats []Format, required ...Format) {
	t.Helper()

	for _, f := range required {
		if !slices.Contains(formats, f) {
			t.Errorf("FormatsForShape(%s) missing %v", label, f)
		}
	}
}

func TestInvalidFormatError(t *testing.T) {
	t.Parallel()

	err := &InvalidFormatError{
		Value:   "invalid",
		Allowed: nil,
	}

	got := err.Error()

	wantContains := []string{"invalid format", "invalid"}
	for _, want := range wantContains {
		if !strings.Contains(got, want) {
			t.Errorf("Error() = %q, should contain %q", got, want)
		}
	}
}

func TestInvalidFormatErrorWithAllowed(t *testing.T) {
	t.Parallel()

	err := &InvalidFormatError{
		Value:   "bogus",
		Allowed: []Format{FormatTable, FormatJSON},
	}

	got := err.Error()

	gentest.AssertOutputContains(t, got, "bogus")

	gentest.AssertOutputContains(t, got, "table")
}

func TestMustRender(t *testing.T) {
	t.Parallel()

	renderer := NewJSONTableRenderer()
	renderer.SetHeaders([]string{"Name"})
	renderer.AddRow([]string{"test"})

	got := MustRender(renderer)
	gentest.AssertOutputContains(t, got, `"Name"`)
}

func TestMustRenderPanics(t *testing.T) {
	t.Parallel()

	defer func() {
		r := recover()
		if r == nil {
			t.Error("MustRender should panic on error")
		}
	}()

	// Use a renderer that always errors
	_ = MustRender(&errorRenderer{})
}

type errorRenderer struct{}

func (e *errorRenderer) Render() (string, error) {
	return "", errTest
}
