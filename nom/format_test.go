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

func TestOperationSymbol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		opType  string
		want    string
		wantLen int
	}{
		{"download", OperationTypeDownload, SymbolDownload, 1},
		{"upload", OperationTypeUpload, SymbolUpload, 1},
		{"unknown", "unknown", "", 0},
		{"empty", "", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := OperationSymbol(tt.opType)
			if got != tt.want {
				t.Errorf("OperationSymbol(%q) = %q, want %q", tt.opType, got, tt.want)
			}
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

	t.Run("paused with elapsed", func(t *testing.T) {
		t.Parallel()

		got := FormatActivityNodeTiming(ActivityStatusPaused, 2*time.Second, 0)
		if got == "" {
			t.Error("expected non-empty timing for paused")
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

func TestGetActivitySummaryString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		running   int
		completed int
		failed    int
		total     int
		wantEmpty bool
	}{
		{"all zero returns empty", 0, 0, 0, 0, true},
		{"only total", 0, 0, 0, 5, false},
		{"with running", 3, 0, 0, 5, false},
		{"with all categories", 1, 2, 3, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetActivitySummaryString(tt.running, tt.completed, tt.failed, tt.total)
			if tt.wantEmpty && got != "" {
				t.Errorf("expected empty, got %q", got)
			}

			if !tt.wantEmpty && got == "" {
				t.Error("expected non-empty")
			}
		})
	}
}

func TestFormatTimingInfo(t *testing.T) {
	t.Parallel()

	t.Run("running activity shows elapsed", func(t *testing.T) {
		t.Parallel()

		ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
		ads.SetRunning()

		got := FormatTimingInfo(ads)
		if got == "" {
			t.Error("expected non-empty timing for running activity")
		}
	})

	t.Run("completed activity shows duration", func(t *testing.T) {
		t.Parallel()

		ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
		ads.SetRunning()
		time.Sleep(10 * time.Millisecond)
		ads.SetCompleted()

		got := FormatTimingInfo(ads)
		if got == "" {
			t.Error("expected non-empty timing for completed activity")
		}
	})

	t.Run("pending with estimated time shows average", func(t *testing.T) {
		t.Parallel()

		ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
		ads.SetEstimatedTime(5 * time.Second)

		got := FormatTimingInfo(ads)
		if got == "" {
			t.Error("expected non-empty timing for pending with estimate")
		}
	})

	t.Run("pending without estimated time returns empty", func(t *testing.T) {
		t.Parallel()

		ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))

		got := FormatTimingInfo(ads)
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})
}
