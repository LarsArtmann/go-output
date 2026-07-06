package markup

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestRenderTable_NilAndError(t *testing.T) {
	t.Parallel()

	formats := []struct {
		name   string
		format output.Format
	}{
		{"XML", output.FormatXML},
		{"HTML", output.FormatHTML},
		{"AsciiDoc", output.FormatAsciiDoc},
	}

	t.Run("nil data produces empty output", func(t *testing.T) {
		t.Parallel()

		for _, f := range formats {
			t.Run(f.name, func(t *testing.T) {
				t.Parallel()

				assertNilDataEmptyOutput(t, f.format, f.name)
			})
		}
	})

	t.Run("writer error propagates", func(t *testing.T) {
		t.Parallel()

		for _, f := range formats {
			t.Run(f.name, func(t *testing.T) {
				t.Parallel()

				data := output.NewTable([]string{"Name"})
				data.AddRow([]string{"Alice"})

				err := output.RenderTable(data, f.format, output.RenderOptions{Writer: &testhelpers.ErrorWriter{}})
				if err == nil {
					t.Fatal("expected error from ErrorWriter")
				}
			})
		}
	})
}

func assertNilDataEmptyOutput(t *testing.T, format output.Format, name string) {
	t.Helper()

	var buf bytes.Buffer

	err := output.RenderTable(nil, format, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable %s nil: %v", name, err)
	}

	if buf.String() != "" {
		t.Errorf("expected empty output for nil data, got %q", buf.String())
	}
}

func TestRenderXMLTable(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	var buf bytes.Buffer

	err := output.RenderTable(data, output.FormatXML, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable xml: %v", err)
	}

	out := buf.String()
	testhelpers.AssertContains(t, out, "<table>", "xml output")
	testhelpers.AssertContains(t, out, "Alice", "xml output")
}

func TestRenderHTMLTable(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	var buf bytes.Buffer

	err := output.RenderTable(data, output.FormatHTML, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable html: %v", err)
	}

	out := buf.String()
	testhelpers.AssertContains(t, out, "<table", "html output")
	testhelpers.AssertContains(t, out, "Alice", "html output")
}

func TestRenderHTMLTable_WithTitle(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name"})
	data.AddRow([]string{"Alice"})

	var buf bytes.Buffer

	err := output.RenderTable(data, output.FormatHTML, output.RenderOptions{
		Writer: &buf,
		Title:  "Test Report",
	})
	if err != nil {
		t.Fatalf("RenderTable html with title: %v", err)
	}

	out := buf.String()
	testhelpers.AssertContains(t, out, "Test Report", "html output with title")
	testhelpers.AssertContains(t, out, "<table", "html output")
}

func TestRenderAsciiDocTable(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	var buf bytes.Buffer

	err := output.RenderTable(data, output.FormatAsciiDoc, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable asciidoc: %v", err)
	}

	out := buf.String()
	testhelpers.AssertContains(t, out, "|===", "asciidoc output")
	testhelpers.AssertContains(t, out, "Alice", "asciidoc output")
}
