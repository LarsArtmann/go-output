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

func BenchmarkMarkdownTableColored(b *testing.B) {
	md := NewMarkdownTable()

	const (
		headerCell = "Header"
		dataCell   = "Cell"
	)

	headers := make([]string, 10)
	for i := range headers {
		headers[i] = headerCell
	}

	md.SetHeaders(headers)
	md.SetColorMode(ColorModeAlways)

	rows := make([][]string, 100)
	for i := range rows {
		row := make([]string, 10)
		for j := range row {
			row[j] = dataCell
		}

		rows[i] = row
		md.AddRow(row)
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = md.Render()
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
