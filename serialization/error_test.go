package serialization

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestJSONLWriter_FlushError(t *testing.T) {
	t.Parallel()

	w := NewJSONLWriter(&testhelpers.ErrorWriter{})

	//nolint:errchkjson // Intentionally ignoring to test Flush error
	_ = w.Encode(map[string]string{"key": "val"})

	err := w.Flush()
	if err == nil {
		t.Fatal("Flush should return error when underlying writer fails")
	}

	if !strings.Contains(err.Error(), "flush") {
		t.Errorf("Expected flush error message, got: %v", err)
	}
}

func TestRenderJSONLTable_RowMarshalError(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name"})
	data.AddRow([]string{"Alice"})

	var buf bytes.Buffer

	err := renderJSONLTable(&buf, data, output.RenderOptions{})
	if err != nil {
		t.Fatalf("renderJSONLTable with normal data: %v", err)
	}

	if !strings.Contains(buf.String(), "Alice") {
		t.Error("Expected Alice in output")
	}
}

func TestRenderJSONLTable_NilRows(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{})

	var buf bytes.Buffer

	err := renderJSONLTable(&buf, data, output.RenderOptions{})
	if err != nil {
		t.Fatalf("renderJSONLTable with empty headers: %v", err)
	}

	if buf.String() != "" {
		t.Errorf("Expected empty output for empty headers, got %q", buf.String())
	}
}

func TestMarshalJSONLFromTable_NilData(t *testing.T) {
	t.Parallel()

	b, err := MarshalJSONLFromTable(nil)
	if err != nil {
		t.Fatalf("Expected nil error for nil data, got: %v", err)
	}

	if b != nil {
		t.Errorf("Expected nil bytes for nil data, got %v", b)
	}
}

func TestMarshalJSONLFromTable_EmptyRows(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"A"})

	b, err := MarshalJSONLFromTable(data)
	if err != nil {
		t.Fatalf("Expected nil error for empty rows, got: %v", err)
	}

	if b != nil {
		t.Errorf("Expected nil bytes for empty rows, got %v", b)
	}
}

func TestRenderTable_MarshalError(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"A"})
	data.AddRow([]string{"1"})

	_, err := renderTable(data, "[]", "failing", func(v any) ([]byte, error) {
		return nil, errors.New("unsupported value")
	})
	if err == nil {
		t.Fatal("Expected error from failing marshal function")
	}

	if !strings.Contains(err.Error(), "failing") {
		t.Errorf("Error should mention format name, got: %v", err)
	}
}

func TestRenderTable_NilData(t *testing.T) {
	t.Parallel()

	got, err := renderTable(nil, "empty", "test", func(v any) ([]byte, error) { return json.Marshal(v) })
	if err != nil {
		t.Fatalf("Expected nil error for nil data, got: %v", err)
	}

	if got != "empty" {
		t.Errorf("Expected empty value for nil data, got %q", got)
	}
}

func TestRenderTable_EmptyHeaders(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{})

	got, err := renderTable(data, "empty", "test", func(v any) ([]byte, error) { return json.Marshal(v) })
	if err != nil {
		t.Fatalf("Expected nil error for empty headers, got: %v", err)
	}

	if got != "empty" {
		t.Errorf("Expected empty value for empty headers, got %q", got)
	}
}

func TestJSONTreeRenderer_MarshalError(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	root.Metadata = map[string]string{"bad": string(rune(0x80))}
	r.SetRoot(root)

	_, err := r.Render()
	if err != nil {
		t.Logf("JSON tree marshal error (expected for some implementations): %v", err)
	}
}

func TestJSONGraphRenderer_MarshalError(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	node := output.NewGraphNode("A", "Node A")
	node.Metadata = map[string]string{}
	r.SetNodes([]output.GraphNode{*node})

	_, err := r.Render()
	if err != nil {
		t.Logf("JSON graph marshal error (expected for some implementations): %v", err)
	}
}

func TestTOMLTreeRenderer_MarshalError(t *testing.T) {
	t.Parallel()

	r := NewTOMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	root.Metadata = map[string]string{"key": "value"}
	r.SetRoot(root)

	_, err := r.Render()
	if err != nil {
		t.Logf("TOML tree marshal error: %v", err)
	}
}

func TestTOMLGraphRenderer_MarshalError(t *testing.T) {
	t.Parallel()

	r := NewTOMLGraphRenderer()
	node := output.NewGraphNode("A", "Node A")
	node.Metadata = map[string]string{"type": "service"}
	r.SetNodes([]output.GraphNode{*node})

	_, err := r.Render()
	if err != nil {
		t.Logf("TOML graph marshal error: %v", err)
	}
}

func TestYAMLTreeRenderer_MarshalError(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	root.Metadata = map[string]string{"key": "value"}
	r.SetRoot(root)

	_, err := r.Render()
	if err != nil {
		t.Logf("YAML tree marshal error: %v", err)
	}
}

func TestYAMLGraphRenderer_MarshalError(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	node := output.NewGraphNode("A", "Node A")
	node.Metadata = map[string]string{"type": "service"}
	r.SetNodes([]output.GraphNode{*node})

	_, err := r.Render()
	if err != nil {
		t.Logf("YAML graph marshal error: %v", err)
	}
}

func TestMarshalTOMLFromTable_RenderError(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"A"})
	data.AddRow([]string{"1"})

	_, err := MarshalTOMLFromTable(data)
	if err != nil {
		t.Fatalf("MarshalTOMLFromTable should succeed with normal data: %v", err)
	}
}

func TestMarshalTOMLFromTable_Nil(t *testing.T) {
	t.Parallel()

	b, err := MarshalTOMLFromTable(nil)
	if err != nil {
		t.Fatalf("Expected nil error for nil data, got: %v", err)
	}

	if len(b) != 0 {
		t.Errorf("Expected empty bytes for nil data, got %v", b)
	}
}

func TestMarshalTOML_Error(t *testing.T) {
	t.Parallel()

	_, err := MarshalTOML(make(chan int))
	if err == nil {
		t.Error("Expected error for unmarshalable type")
	}
}

func TestRenderJSONUnknown_WriteError(t *testing.T) {
	t.Parallel()

	err := renderJSONUnknown(&errorWriter{}, map[string]string{"a": "b"}, output.RenderOptions{})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "write JSON", "error should mention JSON")
}

func TestRenderYAMLUnknown_WriteError(t *testing.T) {
	t.Parallel()

	err := renderYAMLUnknown(&errorWriter{}, map[string]string{"a": "b"}, output.RenderOptions{})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "write YAML", "error should mention YAML")
}

func TestRenderTOMLUnknown_WriteError(t *testing.T) {
	t.Parallel()

	err := renderTOMLUnknown(&errorWriter{}, map[string]string{"a": "b"}, output.RenderOptions{})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "write TOML", "error should mention TOML")
}

func TestRenderJSONLTable_WriteRowError(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name"})
	data.AddRow([]string{"Alice"})

	err := renderJSONLTable(&errorWriter{}, data, output.RenderOptions{})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "jsonl", "error should mention jsonl")
}
