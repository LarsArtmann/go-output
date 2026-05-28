package delimited

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestTSVWriterHeaderAndRow(t *testing.T) {
	t.Parallel()

	var buf strings.Builder

	w := NewTSVWriter(&buf)

	err := w.WriteHeader([]string{"Name", "Value"})
	if err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	err = w.WriteRow([]string{"Alpha", "100"})
	if err != nil {
		t.Fatalf("WriteRow() error = %v", err)
	}

	w.Flush()

	result := buf.String()
	assertContains(t, result, "Name", "TSV should contain header")
	assertContains(t, result, "Alpha", "TSV should contain data")
	assertContains(t, result, "\t", "TSV should use tabs")
}

func TestTSVWriterMultipleRows(t *testing.T) {
	t.Parallel()

	var buf strings.Builder

	w := NewTSVWriter(&buf)

	_ = w.WriteHeader([]string{"A", "B"})
	_ = w.WriteRow([]string{"1", "2"})
	_ = w.WriteRow([]string{"3", "4"})
	w.Flush()

	result := buf.String()

	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestMarshalTSV(t *testing.T) {
	t.Parallel()

	data := [][]string{
		{"Name", "Value"},
		{"Alpha", "100"},
		{"Beta", "200"},
	}

	result, err := MarshalTSV(data)
	if err != nil {
		t.Fatalf("MarshalTSV() error = %v", err)
	}

	tsv := string(result)
	assertContains(t, tsv, "Alpha", "TSV should contain Alpha")
	assertContains(t, tsv, "\t", "TSV should use tabs")
}

func TestMarshalTSVSingleRow(t *testing.T) {
	t.Parallel()

	result, err := MarshalTSV([]string{"A", "B", "C"})
	if err != nil {
		t.Fatalf("MarshalTSV() error = %v", err)
	}

	assertContains(t, string(result), "A\tB\tC", "should marshal single row")
}

func TestMarshalTSVUnsupportedType(t *testing.T) {
	t.Parallel()

	_, err := MarshalTSV(42)
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}

	if !errors.Is(err, ErrUnsupportedType) {
		t.Errorf("error = %v, want ErrUnsupportedType", err)
	}
}

func TestTSVWriterRowError(t *testing.T) {
	t.Parallel()

	w := NewTSVWriter(&errorWriter{})
	_ = w.WriteRow([]string{"test"})
	w.Flush()

	err := w.Error()
	if err == nil {
		t.Fatal("expected error from errorWriter after flush")
	}
}

func TestTSVWriterHeaderError(t *testing.T) {
	t.Parallel()

	w := NewTSVWriter(&errorWriter{})
	_ = w.WriteHeader([]string{"Name"})
	w.Flush()

	err := w.Error()
	if err == nil {
		t.Fatal("expected error from errorWriter after flush")
	}
}

func TestTSVWriterRowsError(t *testing.T) {
	t.Parallel()

	w := NewTSVWriter(&errorWriter{})

	err := w.WriteRows([][]string{{"a"}, {"b"}})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}
}

func TestTSVWriterError(t *testing.T) {
	t.Parallel()

	w := NewTSVWriter(&errorWriter{})
	_ = w.WriteRow([]string{"test"})
	w.Flush()

	err := w.Error()
	if err == nil {
		t.Error("Error() should return error after failed write")
	}
}

func TestMarshalTSVFromTableData(t *testing.T) {
	t.Parallel()

	t.Run("with headers and rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name", "Value"})
		data.AddRow([]string{"Alpha", "100"})
		data.AddRow([]string{"Beta", "200"})

		b, err := MarshalTSVFromTableData(data)
		if err != nil {
			t.Fatalf("MarshalTSVFromTableData() error = %v", err)
		}

		result := string(b)
		assertContains(t, result, "Name", "TSV should contain header")
		assertContains(t, result, "Alpha", "TSV should contain data")
		assertContains(t, result, "\t", "TSV should use tabs")
	})

	t.Run("nil data returns nil", func(t *testing.T) {
		t.Parallel()

		b, err := MarshalTSVFromTableData(nil)
		if err != nil {
			t.Fatalf("MarshalTSVFromTableData(nil) error = %v", err)
		}

		if b != nil {
			t.Errorf("MarshalTSVFromTableData(nil) = %q, want nil", b)
		}
	})

	t.Run("empty data", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData(nil)

		b, err := MarshalTSVFromTableData(data)
		if err != nil {
			t.Fatalf("MarshalTSVFromTableData() error = %v", err)
		}

		if len(b) != 0 {
			t.Errorf("MarshalTSVFromTableData(empty) = %q, want empty", b)
		}
	})

	t.Run("headers only no rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})

		b, err := MarshalTSVFromTableData(data)
		if err != nil {
			t.Fatalf("MarshalTSVFromTableData() error = %v", err)
		}

		result := string(b)
		assertContains(t, result, "Name", "TSV should contain header even with no rows")
	})

	t.Run("with footer row", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name", "Count"})
		data.AddRow([]string{"Alice", "10"})
		data.Footer = []string{"Total", "10"}

		b, err := MarshalTSVFromTableData(data)
		if err != nil {
			t.Fatalf("MarshalTSVFromTableData() error = %v", err)
		}

		result := string(b)
		assertContains(t, result, "Total", "TSV should contain footer row")

		lines := strings.Split(strings.TrimSpace(result), "\n")
		if len(lines) != 3 {
			t.Errorf("expected 3 lines (header + row + footer), got %d", len(lines))
		}
	})
}

func TestTSVRenderTableData(t *testing.T) {
	t.Parallel()

	testRenderTableData(t, output.FormatTSV, "TSV")
}

func TestTSVWriter_WriteFooter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	w := NewTSVWriter(&buf)
	_ = w.WriteHeader([]string{"Name", "Count"})
	_ = w.WriteRow([]string{"Alice", "10"})
	_ = w.WriteFooter([]string{"Total", "10"})
	w.Flush()

	result := buf.String()
	assertContains(t, result, "Total", "TSV WriteFooter should contain footer text")
}
