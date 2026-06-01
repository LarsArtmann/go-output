package serialization

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestRenderTableData_NilAndError(t *testing.T) {
	t.Parallel()

	formats := []struct {
		name   string
		format output.Format
	}{
		{"JSONL", output.FormatJSONL},
		{"YAML", output.FormatYAML},
		{"TOML", output.FormatTOML},
	}

	t.Run("nil data produces empty output", func(t *testing.T) {
		t.Parallel()

		for _, f := range formats {
			t.Run(f.name, func(t *testing.T) {
				t.Parallel()

				var buf bytes.Buffer

				err := output.RenderTableData(nil, f.format, output.RenderOptions{Writer: &buf})
				if err != nil {
					t.Fatalf("RenderTableData %s nil: %v", f.name, err)
				}

				if buf.Len() > 0 {
					t.Errorf("expected empty output for nil data, got %q", buf.String())
				}
			})
		}
	})

	t.Run("writer error propagates", func(t *testing.T) {
		t.Parallel()

		for _, f := range formats {
			t.Run(f.name, func(t *testing.T) {
				t.Parallel()

				data := output.NewTableData([]string{"Name"})
				data.AddRow([]string{"Alice"})

				err := output.RenderTableData(data, f.format, output.RenderOptions{Writer: &errorWriter{}})
				if err == nil {
					t.Fatal("expected error from errorWriter")
				}
			})
		}
	})
}

func TestRenderJSONLTableData(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})
	data.AddRow([]string{"Bob", "25"})

	var buf bytes.Buffer

	err := output.RenderTableData(data, output.FormatJSONL, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData jsonl: %v", err)
	}

	out := buf.String()
	assertContains(t, out, `"Name":"Alice"`, "jsonl output")
	assertContains(t, out, `"Name":"Bob"`, "jsonl output")

	lines := strings.Count(out, "\n")
	if lines != 2 {
		t.Errorf("expected 2 JSONL lines, got %d", lines)
	}
}

func TestRenderYAMLTableData(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	var buf bytes.Buffer

	err := output.RenderTableData(data, output.FormatYAML, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData yaml: %v", err)
	}

	out := buf.String()
	assertContains(t, out, "Name: Alice", "yaml output")
	assertContains(t, out, "Age: \"30\"", "yaml output")
}

func TestRenderTOMLTableData(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	var buf bytes.Buffer

	err := output.RenderTableData(data, output.FormatTOML, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData toml: %v", err)
	}

	out := buf.String()
	assertContains(t, out, "Name", "toml output")
	assertContains(t, out, "Alice", "toml output")
}
