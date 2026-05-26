package serialization

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func BenchmarkMarshalJSONFromTableData(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age", "Email"})
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com"})
	}

	b.ResetTimer()

	for b.Loop() {
		renderer := NewJSONTableRenderer()
		renderer.SetData(data)
		_, _ = renderer.Render()
	}
}

func BenchmarkMarshalYAMLFromTableData(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age", "Email"})
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com"})
	}

	b.ResetTimer()

	for b.Loop() {
		renderer := NewYAMLTableRenderer()
		renderer.SetData(data)
		_, _ = renderer.Render()
	}
}

func BenchmarkMarshalTOMLFromTableData(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age", "Email"})
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com"})
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalTOMLFromTableData(data)
	}
}

func BenchmarkMarshalJSONLFromTableData(b *testing.B) {
	data := output.NewTableData([]string{"Name", "Age", "Email"})
	for range 100 {
		data.AddRow([]string{"Alice", "30", "alice@example.com"})
	}

	b.ResetTimer()

	for b.Loop() {
		_, _ = MarshalJSONLFromTableData(data)
	}
}

func BenchmarkJSONTreeRenderer(b *testing.B) {
	root := output.NewTreeNode("root", "Root")
	for range 50 {
		child := output.NewTreeNode("child", "Child")
		for range 5 {
			child.AddChild(output.NewTreeNode("leaf", "Leaf"))
		}

		root.AddChild(child)
	}

	b.ResetTimer()

	for b.Loop() {
		renderer := NewJSONTreeRenderer()
		renderer.SetRoot(root)
		_, _ = renderer.Render()
	}
}
