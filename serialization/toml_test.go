package serialization

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestMarshalTOML(t *testing.T) {
	t.Parallel()

	t.Run("simple struct", func(t *testing.T) {
		t.Parallel()

		type item struct {
			Name  string `toml:"name"`
			Value int    `toml:"value"`
		}

		b, err := MarshalTOML(item{Name: "Alpha", Value: 42})
		if err != nil {
			t.Fatalf("MarshalTOML() error = %v", err)
		}

		result := string(b)
		if !strings.Contains(result, "Alpha") {
			t.Error("TOML output should contain 'Alpha'")
		}

		if !strings.Contains(result, "42") {
			t.Error("TOML output should contain '42'")
		}
	})
}

func TestUnmarshalTOML(t *testing.T) {
	t.Parallel()

	t.Run("simple struct", func(t *testing.T) {
		t.Parallel()

		type item struct {
			Name  string `toml:"name"`
			Value int    `toml:"value"`
		}

		input := `name = "Alpha"
value = 42
`

		var result item

		err := UnmarshalTOML([]byte(input), &result)
		if err != nil {
			t.Fatalf("UnmarshalTOML() error = %v", err)
		}

		if result.Name != "Alpha" {
			t.Errorf("expected Name=Alpha, got %q", result.Name)
		}

		if result.Value != 42 {
			t.Errorf("expected Value=42, got %d", result.Value)
		}
	})

	t.Run("invalid TOML", func(t *testing.T) {
		t.Parallel()

		var result any

		err := UnmarshalTOML([]byte("not = valid = toml"), &result)
		if err == nil {
			t.Error("expected error for invalid TOML")
		}
	})
}

func TestTOMLTableRenderer(t *testing.T) {
	t.Parallel()

	t.Run("renders table data", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTableRenderer()
		r.SetHeaders([]string{"Name", "Value"})
		r.AddRow([]string{"Alpha", "100"})
		r.AddRow([]string{"Beta", "200"})

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertAllContained(t, out, "Alpha", "Beta")
	})

	t.Run("nil data returns empty marker", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTableRenderer()

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if out != "\n" {
			t.Errorf("expected '\\n', got %q", out)
		}
	})

	t.Run("with rows uses array-of-tables syntax", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTableRenderer()
		r.SetHeaders([]string{"Name", "Value"})
		r.AddRow([]string{"Alpha", "100"})
		r.AddRow([]string{"Beta", "200"})

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertAllContained(t, out, "Alpha", "Beta")

		if !strings.Contains(out, "[[row]]") {
			t.Errorf("TOML output should use [[row]] array-of-tables syntax, got:\n%s", out)
		}
	})
}

func TestMarshalTOMLFromTableData(t *testing.T) {
	t.Parallel()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		b, err := MarshalTOMLFromTableData(nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(b) != 0 {
			t.Errorf("expected empty output, got %q", string(b))
		}
	})

	t.Run("with rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name", "Score"})
		data.AddRow([]string{"Alpha", "90"})
		data.AddRow([]string{"Beta", "75"})

		b, err := MarshalTOMLFromTableData(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		result := string(b)
		if !strings.Contains(result, "Alpha") {
			t.Error("TOML output should contain 'Alpha'")
		}

		if !strings.Contains(result, "[[row]]") {
			t.Error("TOML output should use [[row]] array-of-tables syntax")
		}
	})
}
