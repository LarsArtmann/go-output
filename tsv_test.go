package output

import (
	"strings"
	"testing"
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
	if !strings.Contains(result, "Name") {
		t.Error("TSV should contain header")
	}
	if !strings.Contains(result, "Alpha") {
		t.Error("TSV should contain data")
	}
	if !strings.Contains(result, "\t") {
		t.Error("TSV should use tabs")
	}
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
	if !strings.Contains(tsv, "Alpha") {
		t.Error("TSV should contain Alpha")
	}
	if !strings.Contains(tsv, "\t") {
		t.Error("TSV should use tabs")
	}
}

func BenchmarkTSVWriter(b *testing.B) {
	var buf strings.Builder
	headers := make([]string, 10)
	for i := range headers {
		headers[i] = "Header"
	}
	rows := make([][]string, 100)
	for i := range rows {
		row := make([]string, 10)
		for j := range row {
			row[j] = "Cell"
		}
		rows[i] = row
	}

	b.ResetTimer()
	for b.Loop() {
		buf.Reset()
		w := NewTSVWriter(&buf)
		_ = w.WriteHeader(headers)
		for _, row := range rows {
			_ = w.WriteRow(row)
		}
		w.Flush()
	}
}
