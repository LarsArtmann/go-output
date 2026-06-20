package nom

import (
	"fmt"
	"testing"
	"time"
)

func buildBenchmarkTree(nodeCount int) (*DependencyTree, map[ActivityID]ActivitySnapshot) {
	dt := NewDependencyTree()
	snaps := make(map[ActivityID]ActivitySnapshot)
	now := time.Now()
	_ = now

	dt.AddActivity(ActivityID("root"), nil)

	snaps[ActivityID("root")] = ActivitySnapshot{
		Label: "Root", Status: ActivityStatusRunning,
		Symbol: ActivityStatusRunning.GetSymbol(), Color: ActivityStatusRunning.GetColor(),
	}

	for i := 1; i < nodeCount; i++ {
		id := ActivityID(fmt.Sprintf("step-%04d", i))
		name := fmt.Sprintf("Step %d", i)

		var deps []ActivityID
		if i%5 == 0 {
			deps = []ActivityID{ActivityID("root")}
		} else {
			deps = []ActivityID{ActivityID(fmt.Sprintf("step-%04d", i-1))}
		}

		dt.AddActivity(id, deps)

		status := ActivityStatusPending
		if i < nodeCount/3 {
			status = ActivityStatusCompleted
		} else if i < nodeCount*2/3 {
			status = ActivityStatusRunning
		}

		snaps[id] = ActivitySnapshot{
			Label: name, Status: status,
			Symbol: status.GetSymbol(), Color: status.GetColor(),
		}
	}

	_ = dt.Build()

	return dt, snaps
}

func BenchmarkDependencyTree_Render(b *testing.B) {
	for _, size := range []int{100, 500} {
		b.Run(fmt.Sprintf("%dNodes", size), func(b *testing.B) {
			dt, snaps := buildBenchmarkTree(size)
			maxVisible := size / 2

			b.ResetTimer()

			for b.Loop() {
				if dt.RenderWithSnapshots(snaps, maxVisible, 0) == "" {
					b.Fatal("Render should produce output")
				}
			}
		})
	}
}

func BenchmarkDependencyTree_VisibleNodes(b *testing.B) {
	for _, size := range []int{100, 500} {
		b.Run(fmt.Sprintf("%dNodes", size), func(b *testing.B) {
			dt, snaps := buildBenchmarkTree(size)
			maxVisible := size / 2

			b.ResetTimer()

			for b.Loop() {
				if len(dt.VisibleNodesWithSnapshots(snaps, maxVisible)) == 0 {
					b.Fatal("VisibleNodes should return nodes")
				}
			}
		})
	}
}

func BenchmarkDependencyTree_RenderWithWidth(b *testing.B) {
	for _, size := range []int{100, 500} {
		b.Run(fmt.Sprintf("%dNodes", size), func(b *testing.B) {
			dt, snaps := buildBenchmarkTree(size)
			maxVisible := size / 2

			b.ResetTimer()

			for b.Loop() {
				if dt.RenderWithSnapshots(snaps, maxVisible, 80) == "" {
					b.Fatal("Render should produce output")
				}
			}
		})
	}
}

func BenchmarkDependencyTree_ChildPriority(b *testing.B) {
	dt, snaps := buildBenchmarkTree(500)

	root := dt.GetNode(ActivityID("root"))
	if root == nil {
		b.Fatal("root node not found")
	}

	b.ResetTimer()

	for b.Loop() {
		_ = dt.childPriority(root, snaps)
	}
}
