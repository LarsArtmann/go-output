package markup

import (
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

var assertContains = testhelpers.AssertContains

type errorWriter = testhelpers.ErrorWriter

var errTest = errors.New("test error")

type writeNThenFailWriter = testhelpers.WriteNThenFailWriter

type errorRenderer struct{}

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
