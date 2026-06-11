package markup

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

const emptyTableHTML = `<table class="data-table"></table>`

func TestStreamingRendererInterface(t *testing.T) {
	t.Parallel()

	var _ output.StreamingRenderer = (*StreamingHTMLRenderer)(nil)
}

func TestNewStreamingHTMLRenderer(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	if r == nil {
		t.Fatal("NewStreamingHTMLRenderer() returned nil")
	}
}

func TestStreamingHTMLRendererSetHeaders(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name", "Age", "City"})

	data := r.Data()
	if len(data.Headers) != 3 {
		t.Errorf("Headers length = %d, want 3", len(data.Headers))
	}

	if data.Headers[0] != "Name" {
		t.Errorf("Headers[0] = %q, want %q", data.Headers[0], "Name")
	}
}

func TestStreamingHTMLRendererAddRow(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"Alice"})
	r.AddRow([]string{"Bob"})

	data := r.Data()
	if len(data.Rows) != 2 {
		t.Errorf("Rows length = %d, want 2", len(data.Rows))
	}
}

func TestStreamingHTMLRendererSetData(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	data := output.NewTableData([]string{"Col1", "Col2"})
	data.AddRow([]string{"a", "b"})

	r.SetData(data)

	if r.Data() != data {
		t.Error("SetData() did not set data correctly")
	}
}

func TestStreamingHTMLRendererRender(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name", "Value"})
	r.AddRow([]string{"test", "123"})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "<th>Name</th>", "Render() missing header Name")
	assertContains(t, got, "<td>test</td>", "Render() missing cell test")
}

func TestStreamingHTMLRendererStream(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"A", "B"})
	r.AddRow([]string{"1", "2"})

	got := streamRenderer(t, r)
	assertContains(t, got, "<th>A</th>", "Stream() missing header A")
	assertContains(t, got, "<td>1</td>", "Stream() missing cell 1")
}

func TestStreamingHTMLRendererRenderEmpty(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != emptyTableHTML {
		t.Errorf("Render() = %q, want %q", got, emptyTableHTML)
	}
}

func TestStreamingHTMLRendererStreamEmpty(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()

	var buf bytes.Buffer

	err := r.Stream(&buf)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	got := buf.String()
	if got != emptyTableHTML {
		t.Errorf("Stream() = %q, want %q", got, emptyTableHTML)
	}
}

func streamRenderer(t *testing.T, r *StreamingHTMLRenderer) string {
	t.Helper()

	var buf bytes.Buffer

	err := r.Stream(&buf)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	return buf.String()
}

func TestStreamingHTMLRendererEscapeHTML(t *testing.T) {
	t.Parallel()

	testHTMLEscape(t, func() htmlEscapeTestRenderer {
		return NewStreamingHTMLRenderer()
	}, "StreamingHTMLRenderer")
}

func TestStreamingHTMLRendererEscapeAmpersand(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"Tom & Jerry"})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "Tom &amp; Jerry", "Render() did not escape ampersand")
}

func TestStreamingHTMLRendererMultipleRows(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"ID", "Name", "Score"})

	for i := range 5 {
		r.AddRow([]string{string(rune('0' + i)), "Name" + string(rune('A'+i)), "100"})
	}

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	for i := range 5 {
		if !strings.Contains(got, "Name"+string(rune('A'+i))) {
			t.Errorf("Render() missing row %d", i)
		}
	}
}

func TestStreamingHTMLRendererStreamMidWriteError(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name      string
		remaining int
	}{
		{"header cell write", 2},
		{"row cell write", 4},
		{"row end chunk write", 5},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := NewStreamingHTMLRenderer()
			r.SetHeaders([]string{"A"})
			r.AddRow([]string{"1"})

			err := r.Stream(&writeNThenFailWriter{Remaining: tt.remaining})
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestStreamingHTMLRendererStreamError(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"A"})
	r.AddRow([]string{"1"})

	err := r.Stream(&errorWriter{})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestStreamingHTMLRendererWithFooter(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name", "Count"})
	r.AddRow([]string{"Alice", "10"})
	r.SetFooter([]string{"Total", "10"})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "<tfoot>", "should contain <tfoot>")
	assertContains(t, got, `<td class="footer-cell">Total</td>`, "should contain footer cell with class")
	assertContains(t, got, "</tfoot>", "should contain </tfoot>")

	if strings.Contains(got, "<tfoot>") && !strings.Contains(got, "</tbody>") {
		t.Error("<tfoot> should come after </tbody>")
	}
}

func TestStreamingHTMLRendererNoFooter(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"Alice"})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, "<tfoot>") {
		t.Error("should not contain <tfoot> when no footer")
	}

	assertContains(t, got, "</tbody>", "should still close tbody")
}

func TestStreamHTMLTable_NilData_WriteError(t *testing.T) {
	t.Parallel()

	err := streamHTMLTable(&errorWriter{}, nil)
	if err == nil {
		t.Fatal("expected error for nil data with error writer")
	}

	if !strings.Contains(err.Error(), "write empty table") {
		t.Errorf("error should mention write empty table, got: %v", err)
	}
}
