package nom

import (
	"fmt"
	"testing"
	"time"
)

func buildBenchmarkTree(nodeCount int) *DependencyTree {
	dt := NewDependencyTree()
	now := time.Now()

	dt.AddActivity(ActivityID("root"), NewActivity("root", "Root"), nil)
	node := dt.GetNode(ActivityID("root"))
	node.Status = ActivityStatusRunning
	node.applyVisualStyle()
	node.StartTime = now

	for i := 1; i < nodeCount; i++ {
		id := ActivityID(fmt.Sprintf("step-%04d", i))
		name := fmt.Sprintf("Step %d", i)

		var deps []ActivityID
		if i%5 == 0 {
			deps = []ActivityID{ActivityID("root")}
		} else {
			deps = []ActivityID{ActivityID(fmt.Sprintf("step-%04d", i-1))}
		}

		dt.AddActivity(id, NewActivity(string(id), name), deps)

		status := ActivityStatusPending

		if i < nodeCount/3 {
			status = ActivityStatusCompleted
		} else if i < nodeCount*2/3 {
			status = ActivityStatusRunning
		}

		n := dt.GetNode(id)
		n.Status = status
		n.applyVisualStyle()
		n.StartTime = now
	}

	_ = dt.Build()

	return dt
}

func BenchmarkDependencyTree_Render(b *testing.B) {
	for _, size := range []int{100, 500} {
		b.Run(fmt.Sprintf("%dNodes", size), func(b *testing.B) {
			dt := buildBenchmarkTree(size)
			maxVisible := size / 2

			b.ResetTimer()

			for b.Loop() {
				if dt.RenderString(maxVisible) == "" {
					b.Fatal("Render() should produce output")
				}
			}
		})
	}
}

func BenchmarkDependencyTree_VisibleNodes(b *testing.B) {
	for _, size := range []int{100, 500} {
		b.Run(fmt.Sprintf("%dNodes", size), func(b *testing.B) {
			dt := buildBenchmarkTree(size)
			maxVisible := size / 2

			b.ResetTimer()

			for b.Loop() {
				if len(dt.VisibleNodes(maxVisible)) == 0 {
					b.Fatal("VisibleNodes() should return nodes")
				}
			}
		})
	}
}

func BenchmarkDependencyTree_RenderWithWidth(b *testing.B) {
	for _, size := range []int{100, 500} {
		b.Run(fmt.Sprintf("%dNodes", size), func(b *testing.B) {
			dt := buildBenchmarkTree(size)
			maxVisible := size / 2

			b.ResetTimer()

			for b.Loop() {
				if dt.RenderWithWidth(maxVisible, 80) == "" {
					b.Fatal("RenderWithWidth() should produce output")
				}
			}
		})
	}
}

func BenchmarkDependencyTree_ChildPriority(b *testing.B) {
	dt := buildBenchmarkTree(500)

	root := dt.GetNode(ActivityID("root"))
	if root == nil {
		b.Fatal("root node not found")
	}

	b.ResetTimer()

	for b.Loop() {
		_ = dt.childPriority(root)
	}
}
