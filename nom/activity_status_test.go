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
		{ActivityStatusPending, SymbolPending},
		{ActivityStatusRunning, SymbolRunning},
		{ActivityStatusCompleted, SymbolCompleted},
		{ActivityStatusFailed, SymbolFailed},
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

func TestParseActivityStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    ActivityStatus
		wantErr bool
	}{
		{"pending", ActivityStatusPending, false},
		{"running", ActivityStatusRunning, false},
		{"completed", ActivityStatusCompleted, false},
		{"failed", ActivityStatusFailed, false},
		{"", ActivityStatusPending, true},
		{"unknown", ActivityStatusPending, true},
		{"PENDING", ActivityStatusPending, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := ParseActivityStatus(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseActivityStatus(%q) expected error, got nil", tt.input)
				}

				return
			}

			if err != nil {
				t.Errorf("ParseActivityStatus(%q) unexpected error: %v", tt.input, err)

				return
			}

			if got != tt.want {
				t.Errorf("ParseActivityStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestActivityStatus_IsValid(t *testing.T) {
	t.Parallel()

	valid := []ActivityStatus{
		ActivityStatusPending,
		ActivityStatusRunning,
		ActivityStatusCompleted,
		ActivityStatusFailed,
	}
	for _, s := range valid {
		if !s.IsValid() {
			t.Errorf("ActivityStatus(%d).IsValid() = false, want true", s)
		}
	}

	if ActivityStatus(99).IsValid() {
		t.Errorf("ActivityStatus(99).IsValid() = true, want false")
	}
}

func TestActivityStatus_AllowedValues(t *testing.T) {
	t.Parallel()

	values := ActivityStatus(0).AllowedValues()

	// The registry may contain custom statuses from other tests; verify the
	// four core values appear in order at the start of the list.
	want := []string{"pending", "running", "completed", "failed"}
	if len(values) < len(want) {
		t.Fatalf("AllowedValues() returned %d values, want at least %d", len(values), len(want))
	}

	for i, v := range want {
		if values[i] != v {
			t.Errorf("AllowedValues()[%d] = %q, want %q", i, values[i], v)
		}
	}
}

func TestInvalidActivityStatusError(t *testing.T) {
	t.Parallel()

	err := &InvalidActivityStatusError{Value: "bogus"}
	msg := err.Error()
	testhelpers.AssertContains(t, msg, "invalid activity status", "Error should describe the problem")
	testhelpers.AssertContains(t, msg, "bogus", "Error should contain the invalid value")
}
