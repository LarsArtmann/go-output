package nom

import (
	"testing"
	"time"
)

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
			snap := ActivitySnapshot{
				Label:          "Compile Source Files",
				Status:         tc.status,
				Symbol:         tc.status.GetSymbol(),
				Color:          tc.status.GetColor(),
				CurrentElapsed: 5 * time.Second,
			}

			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				_, _ = formatActivityLabel(snap)
			}
		})
	}
}
