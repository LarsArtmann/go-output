package nom

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestActivityStatus_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status ActivityStatus
		want   string
	}{
		{ActivityStatusPending, "pending"},
		{ActivityStatusRunning, "running"},
		{ActivityStatusCompleted, "completed"},
		{ActivityStatusFailed, "failed"},
		{ActivityStatusPaused, "paused"},
		{ActivityStatus(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			got := tt.status.String()
			testhelpers.AssertEqual(t, "ActivityStatus.String", tt.status, got, tt.want)
		})
	}
}

func TestActivityStatus_GetSymbol(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status ActivityStatus
		want   Symbol
	}{
		{ActivityStatusPending, SymbolPaused},
		{ActivityStatusRunning, SymbolRunning},
		{ActivityStatusCompleted, SymbolCompleted},
		{ActivityStatusFailed, SymbolFailed},
		{ActivityStatusPaused, SymbolPaused},
		{ActivityStatus(99), "?"},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			t.Parallel()

			got := tt.status.GetSymbol()
			testhelpers.AssertEqual(t, "ActivityStatus.GetSymbol", tt.status, got, tt.want)
		})
	}
}

func TestActivityStatus_GetColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status ActivityStatus
	}{
		{ActivityStatusPending},
		{ActivityStatusRunning},
		{ActivityStatusCompleted},
		{ActivityStatusFailed},
		{ActivityStatusPaused},
		{ActivityStatus(99)},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			t.Parallel()

			color := tt.status.GetColor()
			if color == nil {
				t.Errorf("ActivityStatus(%d).GetColor() returned nil", tt.status)
			}
		})
	}
}

func TestStatusStringUnknown(t *testing.T) {
	t.Parallel()

	if StatusStringUnknown != "unknown" {
		t.Errorf("StatusStringUnknown = %q, want %q", StatusStringUnknown, "unknown")
	}
}
