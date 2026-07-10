package serialization

import (
	"bytes"
	"encoding/json/jsontext"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func assertValidJSONLines(t *testing.T, input string, wantLines int) {
	t.Helper()

	lines := strings.Split(strings.TrimRight(input, "\n"), "\n")
	if len(lines) != wantLines {
		t.Fatalf("expected %d lines, got %d", wantLines, len(lines))
	}

	for _, line := range lines {
		if !jsontext.Value([]byte(line)).IsValid() {
			t.Errorf("invalid JSON line: %q", line)
		}
	}
}

func TestMarshalJSONLFromTable(t *testing.T) {
	t.Parallel()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		b, err := MarshalJSONLFromTable(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(b) != 0 {
			t.Errorf("expected empty output, got %q", string(b))
		}
	})

	t.Run("single row", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Name", "Age"})
		data.AddRow([]string{"Alice", "30"})

		b, err := MarshalJSONLFromTable(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertValidJSONLines(t, string(b), 1)
	})

	t.Run("multiple rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Name", "Score"})
		data.AddRow([]string{"Alpha", "90"})
		data.AddRow([]string{"Beta", "75"})

		b, err := MarshalJSONLFromTable(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertValidJSONLines(t, string(b), 2)
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

		assertValidJSONLines(t, out, 2)
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

		assertValidJSONLines(t, buf.String(), 2)
	})
}

func TestJSONLWriter_EncodeError(t *testing.T) {
	t.Parallel()

	w := NewJSONLWriter(&errorWriter{})

	err := w.Encode(map[string]any{"key": make(chan int)})
	if err == nil {
		t.Fatal("Encode should return error for unmarshalable type")
	}
}
