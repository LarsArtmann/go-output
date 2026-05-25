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

	testHTMLEscape(t, func() interface {
		SetHeaders([]string)
		AddRow([]string)
		Render() (string, error)
	} {
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

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"A"})
	r.AddRow([]string{"1"})

	err := r.Stream(&writeNThenFailWriter{remaining: 2})
	if err == nil {
		t.Fatal("expected error from mid-write failure")
	}
}

func TestStreamingHTMLRendererWriteHeaderCellError(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"A", "B"})

	err := r.Stream(&writeNThenFailWriter{remaining: 2})
	if err == nil {
		t.Fatal("expected error during header cell write")
	}
}

func TestStreamingHTMLRendererWriteRowCellError(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"A"})
	r.AddRow([]string{"1"})

	err := r.Stream(&writeNThenFailWriter{remaining: 4})
	if err == nil {
		t.Fatal("expected error during row cell write")
	}
}

func TestStreamingHTMLRendererWriteRowEndChunkError(t *testing.T) {
	t.Parallel()

	r := NewStreamingHTMLRenderer()
	r.SetHeaders([]string{"A"})
	r.AddRow([]string{"1"})

	err := r.Stream(&writeNThenFailWriter{remaining: 5})
	if err == nil {
		t.Fatal("expected error during row end chunk write")
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
