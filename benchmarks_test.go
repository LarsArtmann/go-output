package output

import (
	"testing"
)

func BenchmarkTableDataCreateRowEdges(b *testing.B) {
	data := NewTableData([]string{"A", "B", "C", "D", "E"})
	for range 1000 {
		data.AddRow([]string{"1", "2", "3", "4", "5"})
	}

	b.ResetTimer()

	for b.Loop() {
		data.CreateRowEdges()
	}
}
