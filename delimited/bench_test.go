package delimited

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
)

func BenchmarkCSVWriter(b *testing.B) {
	headers := []string{"Name", "Age", "Email", "City"}
	rows := make([][]string, 100)
	for i := range rows {
		rows[i] = []string{"Alice", "30", "alice@example.com", "Berlin"}
	}

	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer
		w := NewCSVWriter(&buf)
		_ = w.WriteHeader(headers)
		for _, row := range rows {
			_ = w.WriteRow(row)
		}
		w.Flush()
	}
}

func BenchmarkMarshalCSVFromTableData(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age", "Email", "City"})
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com", "Berlin"})
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalCSVFromTableData(data)
	}
}

func BenchmarkTSVWriter(b *testing.B) {
	headers := []string{"Name", "Age", "Email", "City"}
	rows := make([][]string, 100)
	for i := range rows {
		rows[i] = []string{"Alice", "30", "alice@example.com", "Berlin"}
	}

	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer
		w := NewTSVWriter(&buf)
		_ = w.WriteHeader(headers)
		for _, row := range rows {
			_ = w.WriteRow(row)
		}
		w.Flush()
	}
}

func BenchmarkMarshalTSVFromTableData(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age", "Email", "City"})
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com", "Berlin"})
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalTSVFromTableData(data)
	}
}
