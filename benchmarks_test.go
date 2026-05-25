package output

import (
	"testing"
)

func BenchmarkASCIITreeRenderer(b *testing.B) {
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

	renderer := NewASCIITreeRenderer()
	renderer.SetRoot(root)

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

func BenchmarkMarkdownTable(b *testing.B) {
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

type BenchmarkData struct {
	ID        int      `json:"id"         yaml:"id"`
	Name      string   `json:"name"       yaml:"name"`
	Items     []string `json:"items"      yaml:"items"`
	Count     int      `json:"count"      yaml:"count"`
	Active    bool     `json:"active"     yaml:"active"`
	CreatedAt string   `json:"created_at" yaml:"created_at"`
	UpdatedAt string   `json:"updated_at" yaml:"updated_at"`
}

type BenchmarkYAMLStruct = BenchmarkData

func NewBenchmarkData() BenchmarkData {
	return BenchmarkData{
		ID:        12345,
		Name:      "Test Project Alpha",
		Items:     []string{"item1", "item2", "item3", "item4", "item5"},
		Count:     100,
		Active:    true,
		CreatedAt: "2026-03-22T10:00:00Z",
		UpdatedAt: "2026-03-22T12:00:00Z",
	}
}
