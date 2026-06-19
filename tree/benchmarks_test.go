package tree

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func buildBenchmarkTree() *output.TreeNode {
	root := output.NewTreeNode("root", "Root")

	for i := range 100 {
		child := output.NewTreeNode("child", "Child")

		for j := range 10 {
			_ = j

			child.AddChild(output.NewTreeNode("leaf", "Leaf"))
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

func BenchmarkASCIITreeColored(b *testing.B) {
	renderer := NewASCIITreeRenderer()
	renderer.SetRoot(buildBenchmarkTree())
	renderer.SetColorMode(output.ColorModeAlways)

	b.ResetTimer()

	for b.Loop() {
		_, _ = renderer.Render()
	}
}
