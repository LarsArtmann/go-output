package output

import (
	"testing"
)

func BenchmarkTableToGraph_100Rows(b *testing.B) {
	tbl := NewTable([]string{"Name", "Status", "Duration"})

	for i := range 100 {
		tbl.AddRow([]string{"task-" + string(rune(i)), "done", "1.0s"})
	}

	for b.Loop() {
		_ = TableToGraph(tbl)
	}
}

func BenchmarkGraphBuilder_Build_Freeze(b *testing.B) {
	gb := NewGraphBuilder()

	for i := range 100 {
		gb.AddNode(*NewGraphNode("n"+string(rune(i)), "Node"))
	}

	for i := range 99 {
		gb.AddEdge(*NewGraphEdge("n"+string(rune(i)), "n"+string(rune(i+1))))
	}

	for b.Loop() {
		_ = gb.Build()
	}
}

func BenchmarkTableBuilder_Build(b *testing.B) {
	bb := NewTableBuilder().
		SetHeaders("A", "B", "C")

	for i := range 100 {
		bb.AddRow(string(rune(i))+"a", string(rune(i))+"b", string(rune(i))+"c")
	}

	for b.Loop() {
		_ = bb.Build()
	}
}
