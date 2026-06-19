package nom

import (
	"testing"
	"time"
)

// BenchmarkFormatActivityLabel measures the per-node label formatting cost
// (symbol + label + conditional timing) that runs once per visible node during a
// render. Complements the integrated BenchmarkDependencyTree_Render* benchmarks
// by isolating this hot path across every activity status.
func BenchmarkFormatActivityLabel(b *testing.B) {
	statuses := []struct {
		name   string
		status ActivityStatus
	}{
		{"running", ActivityStatusRunning},
		{"completed", ActivityStatusCompleted},
		{"failed", ActivityStatusFailed},
		{"pending", ActivityStatusPending},
	}

	for _, tc := range statuses {
		b.Run(tc.name, func(b *testing.B) {
			node := newActivityNode(ActivityID("compile"), "Compile Source Files")
			node.Status = tc.status
			node.CurrentElapsed = 5 * time.Second
			node.applyVisualStyle()

			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				_, _ = formatActivityLabel(node)
			}
		})
	}
}
