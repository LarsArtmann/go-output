package output

import (
	"testing"
)

func buildBenchmarkTree() *TreeNode {
	root := NewTreeNode("root", "Root")

	for i := range 100 {
		child := NewTreeNode("child", "Child")

		for j := range 10 {
			_ = j

			child.AddChild(NewTreeNode("leaf", "Leaf"))
		}

		root.AddChild(child)

		_ = i
	}

	return root
}

func BenchmarkASCIITreeRenderer(b *testing.B) {
	renderer := NewASCIITreeRenderer()
	renderer.SetRoot(buildBenchmarkTree())

	b.ResetTimer()

	for b.Loop() {
		_, _ = renderer.Render()
	}
}

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

func BenchmarkASCIITreeColored(b *testing.B) {
	renderer := NewASCIITreeRenderer()
	renderer.SetRoot(buildBenchmarkTree())
	renderer.SetColorMode(ColorModeAlways)

	b.ResetTimer()

	for b.Loop() {
		_, _ = renderer.Render()
	}
}
