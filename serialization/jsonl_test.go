package serialization

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestMarshalJSONLFromTableData(t *testing.T) {
	t.Parallel()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		b, err := MarshalJSONLFromTableData(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(b) != 0 {
			t.Errorf("expected empty output, got %q", string(b))
		}
	})

	t.Run("single row", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name", "Age"})
		data.AddRow([]string{"Alice", "30"})

		b, err := MarshalJSONLFromTableData(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
		if len(lines) != 1 {
			t.Fatalf("expected 1 line, got %d", len(lines))
		}

		if !json.Valid([]byte(lines[0])) {
			t.Errorf("invalid JSON line: %q", lines[0])
		}
	})

	t.Run("multiple rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name", "Score"})
		data.AddRow([]string{"Alpha", "90"})
		data.AddRow([]string{"Beta", "75"})

		b, err := MarshalJSONLFromTableData(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d", len(lines))
		}

		for _, line := range lines {
			if !json.Valid([]byte(line)) {
				t.Errorf("invalid JSON line: %q", line)
			}
		}
	})
}

func TestJSONLTableRenderer(t *testing.T) {
	t.Parallel()

	t.Run("renders table data as JSON lines", func(t *testing.T) {
		t.Parallel()

		r := NewJSONLTableRenderer()
		r.SetHeaders([]string{"Name", "Value"})
		r.AddRow([]string{"Alpha", "100"})
		r.AddRow([]string{"Beta", "200"})

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d", len(lines))
		}

		for _, line := range lines {
			if !json.Valid([]byte(line)) {
				t.Errorf("invalid JSON line: %q", line)
			}
		}
	})

	t.Run("nil data returns newline", func(t *testing.T) {
		t.Parallel()

		r := NewJSONLTableRenderer()

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if out != "\n" {
			t.Errorf("expected newline, got %q", out)
		}
	})

	t.Run("empty headers returns newline", func(t *testing.T) {
		t.Parallel()

		r := NewJSONLTableRenderer()
		r.SetHeaders(nil)

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if out != "\n" {
			t.Errorf("expected newline, got %q", out)
		}
	})
}

func TestJSONLWriter(t *testing.T) {
	t.Parallel()

	t.Run("writes JSON lines", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		w := NewJSONLWriter(&buf)

		err := w.Encode(map[string]string{"name": "Alpha"})
		if err != nil {
			t.Fatalf("Encode() error = %v", err)
		}

		err = w.Encode(map[string]string{"name": "Beta"})
		if err != nil {
			t.Fatalf("Encode() error = %v", err)
		}

		err = w.Flush()
		if err != nil {
			t.Fatalf("Flush() error = %v", err)
		}

		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		if len(lines) != 2 {
			t.Fatalf("expected 2 lines, got %d", len(lines))
		}

		for _, line := range lines {
			if !json.Valid([]byte(line)) {
				t.Errorf("invalid JSON line: %q", line)
			}
		}
	})
}
