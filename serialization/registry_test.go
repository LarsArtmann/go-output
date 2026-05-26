package serialization

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

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

func TestRenderJSONLTableData_NilData(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := output.RenderTableData(nil, output.FormatJSONL, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData jsonl nil: %v", err)
	}

	if buf.String() != "" {
		t.Errorf("expected empty output for nil data, got %q", buf.String())
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

func TestRenderYAMLTableData_NilData(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := output.RenderTableData(nil, output.FormatYAML, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData yaml nil: %v", err)
	}

	if buf.String() != "" {
		t.Errorf("expected empty output for nil data, got %q", buf.String())
	}
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

func TestRenderTOMLTableData_NilData(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := output.RenderTableData(nil, output.FormatTOML, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData toml nil: %v", err)
	}

	if buf.String() != "" {
		t.Errorf("expected empty output for nil data, got %q", buf.String())
	}
}

func TestRenderJSONLTableData_WriterError(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Alice"})

	err := output.RenderTableData(data, output.FormatJSONL, output.RenderOptions{Writer: &errorWriter{}})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestRenderYAMLTableData_WriterError(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Alice"})

	err := output.RenderTableData(data, output.FormatYAML, output.RenderOptions{Writer: &errorWriter{}})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestRenderTOMLTableData_WriterError(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Alice"})

	err := output.RenderTableData(data, output.FormatTOML, output.RenderOptions{Writer: &errorWriter{}})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}
