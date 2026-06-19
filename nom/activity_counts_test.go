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

func TestActivityCounts_Summary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		counts ActivityCounts
		want   string
	}{
		{"zero", ActivityCounts{}, ""},
		{"only running", ActivityCounts{Running: 3}, "⏵3"},
		{"only completed", ActivityCounts{Completed: 2}, "✔2"},
		{"running + completed", ActivityCounts{Running: 1, Completed: 2}, "⏵1 ✔2"},
		{"all four", ActivityCounts{Running: 1, Completed: 2, Failed: 3, Pending: 4}, "⏵1 ✔2 ⚠3 ⏸4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.counts.Summary()
			if got != tt.want {
				t.Errorf("Summary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestActivityCounts_CompletionPercent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		counts ActivityCounts
		want   int
	}{
		{"zero", ActivityCounts{}, 0},
		{"all pending", ActivityCounts{Pending: 10}, 0},
		{"half done", ActivityCounts{Completed: 5, Pending: 5}, 50},
		{"all completed", ActivityCounts{Completed: 10}, 100},
		{"all failed", ActivityCounts{Failed: 10}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tt.counts.CompletionPercent()
			if got != tt.want {
				t.Errorf("CompletionPercent() = %d, want %d", got, tt.want)
			}
		})
	}
}
