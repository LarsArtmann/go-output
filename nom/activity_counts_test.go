package nom

import (
	"testing"
)

func TestActivityCounts_Total(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		counts ActivityCounts
		want   int
	}{
		{"zero", ActivityCounts{}, 0},
		{"only running", ActivityCounts{Running: 5}, 5},
		{"mixed", ActivityCounts{Running: 2, Completed: 3, Failed: 1, Pending: 4}, 10},
		{"all completed", ActivityCounts{Completed: 100}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.counts.Total()
			if got != tt.want {
				t.Errorf("Total() = %d, want %d", got, tt.want)
			}
		})
	}
}
