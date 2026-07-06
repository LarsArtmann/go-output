package serialization

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestCQRS_RenderJSON(t *testing.T) {
	tbl := output.NewTable([]string{"Name", "Age"})
	tbl.AddRow([]string{"Alice", "30"})

	got, err := RenderJSON(tbl)
	if err != nil {
		t.Fatalf("RenderJSON failed: %v", err)
	}

	if !strings.Contains(got, "Alice") || !strings.Contains(got, "Name") {
		t.Errorf("expected JSON to contain Name and Alice, got %q", got)
	}
}

func TestCQRS_WriteJSON(t *testing.T) {
	tbl := output.NewTable([]string{"A"})
	tbl.AddRow([]string{"1"})

	var buf strings.Builder

	if err := WriteJSON(&buf, tbl); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output")
	}
}

func TestCQRS_RenderYAML(t *testing.T) {
	tbl := output.NewTable([]string{"Name", "Age"})
	tbl.AddRow([]string{"Alice", "30"})

	got, err := RenderYAML(tbl)
	if err != nil {
		t.Fatalf("RenderYAML failed: %v", err)
	}

	if !strings.Contains(got, "Alice") {
		t.Errorf("expected YAML to contain Alice, got %q", got)
	}
}

func TestCQRS_RenderTOML(t *testing.T) {
	tbl := output.NewTable([]string{"Name"})
	tbl.AddRow([]string{"Alice"})

	got, err := RenderTOML(tbl)
	if err != nil {
		t.Fatalf("RenderTOML failed: %v", err)
	}

	if !strings.Contains(got, "Alice") {
		t.Errorf("expected TOML to contain Alice, got %q", got)
	}
}

func TestCQRS_RenderJSONL(t *testing.T) {
	tbl := output.NewTable([]string{"Name"})
	tbl.AddRow([]string{"Alice"})
	tbl.AddRow([]string{"Bob"})

	got, err := RenderJSONL(tbl)
	if err != nil {
		t.Fatalf("RenderJSONL failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSONL lines, got %d", len(lines))
	}
}

func TestCQRS_WriteJSON_ErrorWriter(t *testing.T) {
	t.Parallel()

	tbl := output.NewTable([]string{"Name"})
	tbl.AddRow([]string{"Alice"})

	err := WriteJSON(&errorWriter{}, tbl)
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestCQRS_WriteYAML_ErrorWriter(t *testing.T) {
	t.Parallel()

	tbl := output.NewTable([]string{"Name"})
	tbl.AddRow([]string{"Alice"})

	err := WriteYAML(&errorWriter{}, tbl)
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestCQRS_WriteTOML_ErrorWriter(t *testing.T) {
	t.Parallel()

	tbl := output.NewTable([]string{"Name"})
	tbl.AddRow([]string{"Alice"})

	err := WriteTOML(&errorWriter{}, tbl)
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestCQRS_WriteJSONL_ErrorWriter(t *testing.T) {
	t.Parallel()

	tbl := output.NewTable([]string{"Name"})
	tbl.AddRow([]string{"Alice"})

	err := WriteJSONL(&errorWriter{}, tbl)
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}
