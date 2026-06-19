package nom

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestFormatDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{"milliseconds", 500 * time.Millisecond, "500ms"},
		{"1 second", 1 * time.Second, "1.0s"},
		{"1.5 seconds", 1500 * time.Millisecond, "1.5s"},
		{"59.95 seconds boundary", 59950 * time.Millisecond, "59.9s"},
		{"1 minute", 1 * time.Minute, "1m"},
		{"1 minute 30 seconds", 90 * time.Second, "1m30s"},
		{"2 minutes", 2 * time.Minute, "2m"},
		{"59 minutes", 59 * time.Minute, "59m"},
		{"59 minutes 59 seconds", 59*time.Minute + 59*time.Second, "59m59s"},
		{"1 hour exactly", 1 * time.Hour, "1h"},
		{"1 hour 30 minutes", 90 * time.Minute, "1h30m"},
		{"1 hour 30 minutes 45 seconds", 1*time.Hour + 30*time.Minute + 45*time.Second, "1h30m"},
		{"2 hours", 2 * time.Hour, "2h"},
		{"2 hours 5 minutes", 2*time.Hour + 5*time.Minute, "2h5m"},
		{"24 hours", 24 * time.Hour, "24h"},
		{"zero", 0, "0ms"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := FormatDuration(tt.duration)
			testhelpers.AssertEqual(t, "FormatDuration", tt.duration, got, tt.want)
		})
	}
}

func TestShouldDisplayTiming(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		duration time.Duration
		want     bool
	}{
		{"below 1s", 500 * time.Millisecond, false},
		{"exactly 1s", 1 * time.Second, true},
		{"above 1s", 5 * time.Second, true},
		{"zero", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ShouldDisplayTiming(tt.duration)
			testhelpers.AssertEqual(t, "ShouldDisplayTiming", tt.duration, got, tt.want)
		})
	}
}

func TestFormatActivityNodeTiming(t *testing.T) {
	t.Parallel()

	t.Run("running with visible timing", func(t *testing.T) {
		t.Parallel()

		got := FormatActivityNodeTiming(ActivityStatusRunning, 5*time.Second, 0)
		if got == "" {
			t.Error("expected non-empty timing for running status")
		}
	})

	t.Run("running with sub-second elapsed returns empty", func(t *testing.T) {
		t.Parallel()

		got := FormatActivityNodeTiming(ActivityStatusRunning, 100*time.Millisecond, 0)
		if got != "" {
			t.Errorf("expected empty for sub-second, got %q", got)
		}
	})

	t.Run("pending with estimated time", func(t *testing.T) {
		t.Parallel()

		got := FormatActivityNodeTiming(ActivityStatusPending, 0, 10*time.Second)
		if got == "" {
			t.Error("expected non-empty timing for pending with estimate")
		}
	})

	t.Run("completed with elapsed", func(t *testing.T) {
		t.Parallel()

		got := FormatActivityNodeTiming(ActivityStatusCompleted, 3*time.Second, 0)
		if got == "" {
			t.Error("expected non-empty timing for completed")
		}
	})

	t.Run("failed with elapsed", func(t *testing.T) {
		t.Parallel()

		got := FormatActivityNodeTiming(ActivityStatusFailed, 2*time.Second, 0)
		if got == "" {
			t.Error("expected non-empty timing for failed")
		}
	})

	t.Run("zero elapsed returns empty", func(t *testing.T) {
		t.Parallel()

		got := FormatActivityNodeTiming(ActivityStatusRunning, 0, 0)
		if got != "" {
			t.Errorf("expected empty for zero elapsed, got %q", got)
		}
	})
}
