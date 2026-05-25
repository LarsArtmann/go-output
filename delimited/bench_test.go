package delimited

import (
	"bytes"
	"testing"

	"github.com/larsartmann/go-output"
)

var benchHeaders = []string{"Name", "Age", "Email", "City"}

var benchRows = func() [][]string {
	rows := make([][]string, 100)
	for i := range rows {
		rows[i] = []string{"Alice", "30", "alice@example.com", "Berlin"}
	}

	return rows
}()

func benchTableData() *output.TableData {
	data := output.NewTableData(benchHeaders)
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

func BenchmarkMarshalCSVFromTableData(b *testing.B) {
	data := benchTableData()

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalCSVFromTableData(data)
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

func BenchmarkMarshalTSVFromTableData(b *testing.B) {
	data := benchTableData()

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalTSVFromTableData(data)
	}
}
