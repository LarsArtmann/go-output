package serialization

import (
	"io"
	"testing"

	"github.com/larsartmann/go-output"
)

func benchTable(rows int) *output.Table {
	data := output.NewTable([]string{"Name", "Status", "Duration", "Timestamp", "ID"})
	for i := range rows {
		data.AddRow([]string{
			"task-" + string(rune('a'+i%26)),
			"completed",
			"1.2s",
			"2026-07-06T12:00:00Z",
			"0001",
		})
	}

	return data
}

func BenchmarkCQRS_WriteJSON_100Rows(b *testing.B) {
	data := benchTable(100)

	b.ResetTimer()

	for range b.N {
		_ = WriteJSON(io.Discard, data)
	}
}

func BenchmarkCQRS_WriteYAML_100Rows(b *testing.B) {
	data := benchTable(100)

	b.ResetTimer()

	for range b.N {
		_ = WriteYAML(io.Discard, data)
	}
}

func BenchmarkCQRS_WriteTOML_100Rows(b *testing.B) {
	data := benchTable(100)

	b.ResetTimer()

	for range b.N {
		_ = WriteTOML(io.Discard, data)
	}
}
