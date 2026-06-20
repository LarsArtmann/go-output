package nom

import (
	"fmt"
	"testing"
)

// BenchmarkGetActivityCounts measures the cost of GetActivityCounts across
// increasing activity counts. With the incremental cache this is O(1) — the
// time must stay flat as N grows. (The previous O(n) scan-every-activity
// implementation scaled linearly.)
func BenchmarkGetActivityCounts(b *testing.B) {
	for _, count := range []int{10, 100, 1000, 10000} {
		b.Run(fmt.Sprintf("%dActivities", count), func(b *testing.B) {
			sub := buildBenchSubscriber(b, count)

			b.ResetTimer()

			for b.Loop() {
				_ = sub.GetActivityCounts()
			}
		})
	}
}
