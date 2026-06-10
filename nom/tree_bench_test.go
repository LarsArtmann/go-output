package nom

import (
	"fmt"
	"testing"
	"time"
)

func buildBenchmarkTree(nodeCount int) *DependencyTree {
	dt := NewDependencyTree()
	now := time.Now()

	dt.AddActivity(ActivityID("root"), "Root", nil)
	dt.UpdateActivityStatus(ActivityID("root"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)

	for i := 1; i < nodeCount; i++ {
		id := ActivityID(fmt.Sprintf("step-%04d", i))
		name := fmt.Sprintf("Step %d", i)

		var deps []ActivityID
		if i%5 == 0 {
			deps = []ActivityID{ActivityID("root")}
		} else {
			deps = []ActivityID{ActivityID(fmt.Sprintf("step-%04d", i-1))}
		}

		dt.AddActivity(id, name, deps)

		status := ActivityStatusPending
		symbol := SymbolPaused
		color := ColorPaused
		if i < nodeCount/3 {
			status = ActivityStatusCompleted
			symbol = SymbolCompleted
			color = ColorCompleted
		} else if i < nodeCount*2/3 {
			status = ActivityStatusRunning
			symbol = SymbolRunning
			color = ColorRunning
		}

		dt.UpdateActivityStatus(id, status, symbol, color, now, 0)
	}

	dt.EnsureBuild()

	return dt
}

func BenchmarkDependencyTree_Render_100Nodes(b *testing.B) {
	dt := buildBenchmarkTree(100)

	b.ResetTimer()

	for b.Loop() {
		result := dt.Render(50)
		if result == "" {
			b.Fatal("Render() should produce output")
		}
	}
}

func BenchmarkDependencyTree_Render_500Nodes(b *testing.B) {
	dt := buildBenchmarkTree(500)

	b.ResetTimer()

	for b.Loop() {
		result := dt.Render(100)
		if result == "" {
			b.Fatal("Render() should produce output")
		}
	}
}

func BenchmarkDependencyTree_VisibleNodes_100Nodes(b *testing.B) {
	dt := buildBenchmarkTree(100)

	b.ResetTimer()

	for b.Loop() {
		nodes := dt.VisibleNodes(50)
		if len(nodes) == 0 {
			b.Fatal("VisibleNodes() should return nodes")
		}
	}
}

func BenchmarkDependencyTree_VisibleNodes_500Nodes(b *testing.B) {
	dt := buildBenchmarkTree(500)

	b.ResetTimer()

	for b.Loop() {
		nodes := dt.VisibleNodes(100)
		if len(nodes) == 0 {
			b.Fatal("VisibleNodes() should return nodes")
		}
	}
}
