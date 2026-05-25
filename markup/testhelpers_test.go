package markup

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

var assertContains = testhelpers.AssertContains

type errorWriter = testhelpers.ErrorWriter

type writeNThenFailWriter = testhelpers.WriteNThenFailWriter

type errorRenderer = testhelpers.ErrorRenderer

type htmlEscapeTestRenderer interface {
	SetHeaders([]string)
	AddRow([]string)
	Render() (string, error)
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

func testHTMLEscape(t *testing.T, newRenderer func() htmlEscapeTestRenderer, name string) {
	t.Helper()

	tests := []struct {
		subtest string
		header  string
		input   string
		want    string
		msg     string
	}{
		{
			subtest: "escapes <brackets>",
			header:  "Data",
			input:   "<script>alert('xss')</script>",
			want:    "&lt;script&gt;",
			msg:     "should escape < and >",
		},
		{
			subtest: "escapes & ampersand",
			header:  "Name",
			input:   "Tom & Jerry",
			want:    "Tom &amp; Jerry",
			msg:     "should escape &",
		},
	}

	for _, tt := range tests {
		t.Run(tt.subtest, func(t *testing.T) {
			t.Parallel()

			r := newRenderer()
			r.SetHeaders([]string{tt.header})
			r.AddRow([]string{tt.input})

			got, err := r.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			assertContains(t, got, tt.want, name+" "+tt.msg)
		})
	}
}
