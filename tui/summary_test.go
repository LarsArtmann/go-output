package tui

import (
	"testing"
	"time"
)

func TestFormatElapsedTime(t *testing.T) {
	t.Parallel()

	got := formatElapsedTime(5 * time.Second)
	if got == "" {
		t.Error("formatElapsedTime() should not return empty")
	}
}

func TestBuildUniversalSummary(t *testing.T) {
	t.Parallel()

	summary := buildUniversalSummary(2, 5, 30*time.Second, 75.0)
	if summary == "" {
		t.Error("buildUniversalSummary() should not return empty")
	}
}

func TestBuildActivityCountsSummary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		running   int
		completed int
		failed    int
		pending   int
		wantEmpty bool
	}{
		{"all zero returns empty", 0, 0, 0, 0, true},
		{"with counts", 1, 2, 1, 3, false},
		{"only running", 3, 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildActivityCountsSummary(tt.running, tt.completed, tt.failed, tt.pending)
			if tt.wantEmpty && got != "" {
				t.Errorf("expected empty, got %q", got)
			}
			if !tt.wantEmpty && got == "" {
				t.Error("expected non-empty summary")
			}
		})
	}
}

func TestBuildNOMSummary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		ok, fail, skip, activity int
		duration                 time.Duration
	}{
		{"with activity counts", 1, 2, 0, 1, 10 * time.Second},
		{"with zero counts still shows timing", 0, 0, 0, 0, 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildNOMSummary(tt.ok, tt.fail, tt.skip, tt.activity, tt.duration)
			if got == "" {
				t.Error("expected non-empty summary")
			}
		})
	}
}

func TestGetStateStyle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		state WorkflowState
	}{
		{"idle", WorkflowStateIdle},
		{"running", WorkflowStateRunning},
		{"completed", WorkflowStateCompleted},
		{"errored", WorkflowStateErrored},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			msg, color := getStateStyle(tt.state)
			_ = msg
			_ = color
		})
	}
}

func TestApplyStateSummary(t *testing.T) {
	t.Parallel()

	t.Run("with state summary", func(t *testing.T) {
		t.Parallel()

		summary, style := applyStateSummary("test", WorkflowStateCompleted, 5, 10*time.Second)
		if summary == "" {
			t.Error("expected non-empty summary")
		}

		_ = style
	})

	t.Run("running state returns original summary", func(t *testing.T) {
		t.Parallel()

		summary, _ := applyStateSummary("original", WorkflowStateRunning, 0, 0)
		if summary != "original" {
			t.Errorf("summary = %q, want %q", summary, "original")
		}
	})
}

func TestCreateSummaryStyle(t *testing.T) {
	t.Parallel()

	style := createSummaryStyle()
	_ = style
}
