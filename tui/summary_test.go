package tui

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output/nom"
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
		counts    nom.ActivityCounts
		wantEmpty bool
	}{
		{"all zero returns empty", nom.ActivityCounts{}, true},
		{"with counts", nom.ActivityCounts{Running: 1, Completed: 2, Failed: 1, Pending: 3}, false},
		{"only running", nom.ActivityCounts{Running: 3}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildActivityCountsSummary(tt.counts)
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
		name     string
		counts   nom.ActivityCounts
		duration time.Duration
	}{
		{"with activity counts", nom.ActivityCounts{Running: 1, Completed: 2, Failed: 0, Pending: 1}, 10 * time.Second},
		{"with zero counts still shows timing", nom.ActivityCounts{}, 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := buildNOMSummary(tt.counts, tt.duration)
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
