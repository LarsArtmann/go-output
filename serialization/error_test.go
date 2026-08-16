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

	data := output.NewTableWithRow([]string{"A"}, "1")

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

// The *Renderer_MarshalError tests below were originally written to expect
// marshal failures, but their inputs (TreeNode/GraphNode with string-only
// metadata) marshal fine — and the old bodies only t.Logf'd the error, so
// they could never fail. They now assert the opposite invariant: the tree
// and graph renderers succeed and carry the metadata through.

func TestJSONTreeRenderer_RendersMetadata(t *testing.T) {
	t.Parallel()

	r := NewJSONTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	root.Metadata = map[string]string{"bad": string(rune(0x80))}
	r.SetRoot(root)

	out, err := r.Render()
	if err != nil {
		t.Fatalf("JSON tree render should succeed: %v", err)
	}

	// Invalid UTF-8 is replaced, not rejected — the payload must still carry the key.
	if !strings.Contains(string(out), "bad") {
		t.Errorf("JSON tree output should contain metadata key %q, got: %s", "bad", out)
	}
}

func TestJSONGraphRenderer_RendersMetadata(t *testing.T) {
	t.Parallel()

	r := NewJSONGraphRenderer()
	node := output.NewGraphNode("A", "Node A")
	node.Metadata = map[string]string{}
	r.SetNodes([]output.GraphNode{*node})

	out, err := r.Render()
	if err != nil {
		t.Fatalf("JSON graph render should succeed: %v", err)
	}

	if !strings.Contains(string(out), "Node A") {
		t.Errorf("JSON graph output should contain node label, got: %s", out)
	}
}

func TestTOMLTreeRenderer_RendersMetadata(t *testing.T) {
	t.Parallel()

	r := NewTOMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	root.Metadata = map[string]string{"key": "value"}
	r.SetRoot(root)

	out, err := r.Render()
	if err != nil {
		t.Fatalf("TOML tree render should succeed: %v", err)
	}

	if !strings.Contains(string(out), "key") || !strings.Contains(string(out), "value") {
		t.Errorf("TOML tree output should contain metadata, got: %s", out)
	}
}

func TestTOMLGraphRenderer_RendersMetadata(t *testing.T) {
	t.Parallel()

	r := NewTOMLGraphRenderer()
	node := output.NewGraphNode("A", "Node A")
	node.Metadata = map[string]string{"type": "service"}
	r.SetNodes([]output.GraphNode{*node})

	out, err := r.Render()
	if err != nil {
		t.Fatalf("TOML graph render should succeed: %v", err)
	}

	if !strings.Contains(string(out), "service") {
		t.Errorf("TOML graph output should contain metadata value, got: %s", out)
	}
}

func TestYAMLTreeRenderer_RendersMetadata(t *testing.T) {
	t.Parallel()

	r := NewYAMLTreeRenderer()
	root := output.NewTreeNode("root", "Root")
	root.Metadata = map[string]string{"key": "value"}
	r.SetRoot(root)

	out, err := r.Render()
	if err != nil {
		t.Fatalf("YAML tree render should succeed: %v", err)
	}

	if !strings.Contains(string(out), "key: value") {
		t.Errorf("YAML tree output should contain metadata, got: %s", out)
	}
}

func TestYAMLGraphRenderer_RendersMetadata(t *testing.T) {
	t.Parallel()

	r := NewYAMLGraphRenderer()
	node := output.NewGraphNode("A", "Node A")
	node.Metadata = map[string]string{"type": "service"}
	r.SetNodes([]output.GraphNode{*node})

	out, err := r.Render()
	if err != nil {
		t.Fatalf("YAML graph render should succeed: %v", err)
	}

	if !strings.Contains(string(out), "type: service") {
		t.Errorf("YAML graph output should contain metadata, got: %s", out)
	}
}

func TestMarshalTOMLFromTable_RenderError(t *testing.T) {
	t.Parallel()

	data := output.NewTableWithRow([]string{"A"}, "1")

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

	data := output.NewTableWithRow([]string{"Name"}, "Alice")

	err := renderJSONLTable(&errorWriter{}, data, output.RenderOptions{})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "jsonl", "error should mention jsonl")
}
