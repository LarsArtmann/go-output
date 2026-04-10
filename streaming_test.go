package output

import (
	"bytes"
	"strings"
	"testing"

)

const emptyTableHTML = `<table class="data-table"></table>`

func TestStreamingRendererInterface(t *testing.T) {
	t.Parallel()

	var _ StreamingRenderer = (*StreamingHTMLRenderer)(nil)
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

	if len(r.data.Headers) != 3 {
		t.Errorf("Headers length = %d, want 3", len(r.data.Headers))
	}

	if r.data.Headers[0] != "Name" {
		t.Errorf("Headers[0] = %q, want %q", r.data.Headers[0], "Name")
	}
}

func TestStreamingHTMLRendererAddRow(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"Alice"})
	r.AddRow([]string{"Bob"})

	if len(r.data.Rows) != 2 {
		t.Errorf("Rows length = %d, want 2", len(r.data.Rows))
	}
}

func TestStreamingHTMLRendererSetData(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	data := NewTableData([]string{"Col1", "Col2"})
	data.AddRow([]string{"a", "b"})

	r.SetData(data)

	if r.data != data {
		t.Error("SetData() did not set data correctly")
	}
}

func TestStreamingHTMLRendererRender(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name", "Value"})
	r.AddRow([]string{"test", "123"})

	got := r.Render()
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

	got := r.Render()

	want := emptyTableHTML
	if got != want {
		t.Errorf("Render() = %q, want %q", got, want)
	}
}

func TestStreamingHTMLRendererStreamEmpty(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	got := streamRenderer(t, r)

	assertEmptyStreamingOutput(t, got)
}

// streamRenderer streams the renderer output and returns the result.
func streamRenderer(t *testing.T, r *StreamingHTMLRenderer) string {
	t.Helper()

	var buf bytes.Buffer

	err := r.Stream(&buf)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	return buf.String()
}

// assertEmptyStreamingOutput verifies that empty streaming output matches expected.
func assertEmptyStreamingOutput(t *testing.T, got string) {
	t.Helper()

	if got != emptyTableHTML {
		t.Errorf("Stream() = %q, want %q", got, emptyTableHTML)
	}
}

func TestStreamingHTMLRendererEscapeHTML(t *testing.T) {
	t.Parallel()
	testHTMLEscapeShared(
		t,
		func() htmlEscapeTestRenderer { return NewStreamingHTMLRenderer() },
		"StreamingHTMLRenderer",
	)
}

func TestStreamingHTMLRendererEscapeAmpersand(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"Name"})
	r.AddRow([]string{"Tom & Jerry"})

	got := r.Render()
	assertContains(t, got, "Tom &amp; Jerry", "Render() did not escape ampersand")
}

func TestStreamingRendererFromRenderer(t *testing.T) {
	t.Parallel()

	original := &testRenderer{output: "test-output"}
	adapter := StreamingRendererFromRenderer(original)

	if adapter.Render() != "test-output" {
		t.Errorf("Render() = %q, want %q", adapter.Render(), "test-output")
	}

	var buf bytes.Buffer

	err := adapter.Stream(&buf)
	if err != nil {
		t.Fatalf("Stream() error = %v", err)
	}

	if buf.String() != "test-output" {
		t.Errorf("Stream() wrote %q, want %q", buf.String(), "test-output")
	}
}

func TestStreamingHTMLRendererMultipleRows(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"ID", "Name", "Score"})

	for i := range 5 {
		r.AddRow([]string{string(rune('0' + i)), "Name" + string(rune('A'+i)), "100"})
	}

	got := r.Render()
	for i := range 5 {
		if !strings.Contains(got, "Name"+string(rune('A'+i))) {
			t.Errorf("Render() missing row %d", i)
		}
	}
}
