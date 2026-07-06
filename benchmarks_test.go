package output

import (
	"testing"
)

func BenchmarkTableCreateRowEdges(b *testing.B) {
	data := NewTable([]string{"A", "B", "C", "D", "E"})
	for range 1000 {
		data.AddRow([]string{"1", "2", "3", "4", "5"})
	}

	b.ResetTimer()

	for b.Loop() {
		data.CreateRowEdges()
	}
}
