package markup

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

var assertContains = testhelpers.AssertContains

type errorWriter struct{}

var errWrite = errors.New("write error")

func (e *errorWriter) Write(_ []byte) (int, error) {
	return 0, errWrite
}

var _ io.Writer = (*errorWriter)(nil)

type writeNThenFailWriter struct {
	remaining int
}

func (w *writeNThenFailWriter) Write(p []byte) (int, error) {
	if w.remaining <= 0 {
		return 0, errWrite
	}

	w.remaining--

	return len(p), nil
}

var _ io.Writer = (*writeNThenFailWriter)(nil)

type errorRenderer struct{}

var errTest = errors.New("test error")

func (e *errorRenderer) Render() (string, error) {
	return "", errTest
}

type testRenderer struct {
	output string
}

func (r *testRenderer) Render() (string, error) {
	return r.output, nil
}

func testEmptyRendererOutput(
	t *testing.T,
	renderer output.Renderer,
	expected []struct{ Substring, Message string },
) {
	t.Helper()

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	for _, e := range expected {
		if !strings.Contains(got, e.Substring) {
			t.Error(e.Message)
		}
	}
}

func testHTMLEmptyExpected() []struct{ Substring, Message string } {
	return []struct{ Substring, Message string }{
		{"<table", "Empty table should still be valid HTML"},
		{"</table>", "Empty table should have closing tag"},
	}
}

func testHTMLEscape(t *testing.T, newRenderer func() interface {
	SetHeaders([]string)
	AddRow([]string)
	Render() (string, error)
}, name string,
) {
	t.Helper()

	t.Run("escapes <brackets>", func(t *testing.T) {
		t.Parallel()

		r := newRenderer()
		r.SetHeaders([]string{"Data"})
		r.AddRow([]string{"<script>alert('xss')</script>"})

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertContains(t, got, "&lt;script&gt;", name+" should escape < and >")
	})

	t.Run("escapes & ampersand", func(t *testing.T) {
		t.Parallel()

		r := newRenderer()
		r.SetHeaders([]string{"Name"})
		r.AddRow([]string{"Tom & Jerry"})

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertContains(t, got, "Tom &amp; Jerry", name+" should escape &")
	})
}
