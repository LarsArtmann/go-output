package table

import (
	"fmt"
	"strconv"
	"testing"
)

func BenchmarkTableRender_100Rows5Cols(b *testing.B) {
	tbl := New()
	tbl.SetHeaders("Name", "Status", "Duration", "Count", "Size")

	for i := range 100 {
		tbl.AddRow(
			fmt.Sprintf("Row %d", i),
			"✓",
			fmt.Sprintf("%.1fs", float64(i)*0.1),
			strconv.Itoa(i*10),
			fmt.Sprintf("%dKB", i*5),
		)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_, _ = tbl.Render()
	}
}
