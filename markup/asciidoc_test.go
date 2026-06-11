package markup

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func assertAllContained(t *testing.T, haystack string, needles ...string) {
	t.Helper()

	for _, n := range needles {
		if !strings.Contains(haystack, n) {
			t.Errorf("should contain %q", n)
		}
	}
}

func TestMarshalAsciiDocFromTableData(t *testing.T) {
	t.Parallel()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		b, err := MarshalAsciiDocFromTableData(nil)
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

		b, err := MarshalAsciiDocFromTableData(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		result := string(b)
		if !strings.Contains(result, "|===") {
			t.Error("AsciiDoc output should contain table delimiter")
		}

		if !strings.Contains(result, "Alice") {
			t.Error("AsciiDoc output should contain 'Alice'")
		}
	})

	t.Run("escapes pipe characters", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"a|b"})

		b, err := MarshalAsciiDocFromTableData(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		result := string(b)
		if !strings.Contains(result, `\|`) {
			t.Error("pipe in cell should be escaped")
		}
	})

	t.Run("with footer row", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name", "Count"})
		data.AddRow([]string{"Alice", "10"})
		data.Footer = []string{"Total", "10"}

		b, err := MarshalAsciiDocFromTableData(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		result := string(b)
		if !strings.Contains(result, "Total") {
			t.Error("AsciiDoc output should contain footer text 'Total'")
		}
	})
}

func TestAsciiDocTableRenderer(t *testing.T) {
	t.Parallel()

	t.Run("renders table", func(t *testing.T) {
		t.Parallel()

		r := NewAsciiDocTableRenderer()
		r.SetHeaders([]string{"Name", "Value"})
		r.AddRow([]string{"Alpha", "100"})
		r.AddRow([]string{"Beta", "200"})

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertAllContained(t, out, "|===", "| Name", "| Alpha")
	})

	t.Run("nil data returns empty", func(t *testing.T) {
		t.Parallel()

		r := NewAsciiDocTableRenderer()

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if out != "" {
			t.Errorf("expected empty, got %q", out)
		}
	})

	t.Run("empty headers returns empty", func(t *testing.T) {
		t.Parallel()

		r := NewAsciiDocTableRenderer()
		r.SetHeaders(nil)

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if out != "" {
			t.Errorf("expected empty, got %q", out)
		}
	})
}

func TestEscapeAsciiDoc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"a|b", `a\|b`},
		{"a|b|c", `a\|b\|c`},
		{"no pipes", "no pipes"},
		{"*bold*", `\*bold\*`},
		{"_italic_", `\_italic\_`},
		{"`code`", `\` + "`" + `code\` + "`"},
		{"~sub~", `\~sub\~`},
		{"^super^", `\^super\^`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got := escapeAsciiDoc(tt.input)
			if got != tt.want {
				t.Errorf("escapeAsciiDoc(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
