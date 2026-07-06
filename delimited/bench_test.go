package delimited

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
)

var benchHeaders = []string{"Name", "Age", "Email", "City"}

var benchRows = func() [][]string {
	rows := make([][]string, 0, 100)
	for range 100 {
		rows = append(rows, []string{"Alice", "30", "alice@example.com", "Berlin"})
	}

	return rows
}()

func benchTable() *output.Table {
	data := output.NewTable(benchHeaders)
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com", "Berlin"})
	}

	return data
}

func BenchmarkCSVWriter(b *testing.B) {
	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer

		w := NewCSVWriter(&buf)

		_ = w.WriteHeader(benchHeaders)

		for _, row := range benchRows {
			_ = w.WriteRow(row)
		}

		w.Flush()
	}
}

func BenchmarkMarshalCSVFromTable(b *testing.B) {
	data := benchTable()

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalCSVFromTable(data)
	}
}

func BenchmarkMarshalCSVFromTableWithFooter(b *testing.B) {
	data := benchTable()
	data.SetFooter([]string{"Total", "100", "", ""})

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalCSVFromTable(data)
	}
}

func BenchmarkTSVWriter(b *testing.B) {
	b.ResetTimer()

	for b.Loop() {
		var buf bytes.Buffer

		w := NewTSVWriter(&buf)

		_ = w.WriteHeader(benchHeaders)

		for _, row := range benchRows {
			_ = w.WriteRow(row)
		}

		w.Flush()
	}
}

func BenchmarkMarshalTSVFromTable(b *testing.B) {
	data := benchTable()

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalTSVFromTable(data)
	}
}
