package table

import (
	"fmt"
	"testing"
)

func BenchmarkTableRender_100Rows5Cols(b *testing.B) {
	tbl := New()
	tbl.SetHeaders("Name", "Status", "Duration", "Count", "Size")
	for i := 0; i < 100; i++ {
		tbl.AddRow(
			fmt.Sprintf("Row %d", i),
			"✓",
			fmt.Sprintf("%.1fs", float64(i)*0.1),
			fmt.Sprintf("%d", i*10),
			fmt.Sprintf("%dKB", i*5),
		)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = tbl.Render()
	}
}
