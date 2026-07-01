package nom

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestActivityKind_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind ActivityKind
		want string
	}{
		{ActivityKindTask, "task"},
		{ActivityKindPhase, "phase"},
		{ActivityKind(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			got := tt.kind.String()
			testhelpers.AssertEqual(t, "ActivityKind.String", tt.kind, got, tt.want)
		})
	}
}

func TestActivityKind_IsPhase(t *testing.T) {
	t.Parallel()

	if !ActivityKindPhase.IsPhase() {
		t.Error("ActivityKindPhase.IsPhase() = false, want true")
	}

	if ActivityKindTask.IsPhase() {
		t.Error("ActivityKindTask.IsPhase() = true, want false")
	}
}

func TestParseActivityKind(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input   string
		want    ActivityKind
		wantErr bool
	}{
		{"task", ActivityKindTask, false},
		{"phase", ActivityKindPhase, false},
		{"", ActivityKindTask, true},
		{"unknown", ActivityKindTask, true},
		{"TASK", ActivityKindTask, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := ParseActivityKind(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseActivityKind(%q) expected error, got nil", tt.input)
				}

				return
			}

			if err != nil {
				t.Errorf("ParseActivityKind(%q) unexpected error: %v", tt.input, err)

				return
			}

			if got != tt.want {
				t.Errorf("ParseActivityKind(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestActivityKind_IsValid(t *testing.T) {
	t.Parallel()

	valid := []ActivityKind{ActivityKindTask, ActivityKindPhase}
	for _, k := range valid {
		if !k.IsValid() {
			t.Errorf("ActivityKind(%d).IsValid() = false, want true", k)
		}
	}

	if ActivityKind(99).IsValid() {
		t.Errorf("ActivityKind(99).IsValid() = true, want false")
	}
}

func TestInvalidActivityKindError(t *testing.T) {
	t.Parallel()

	err := &InvalidActivityKindError{Value: "bogus"}
	msg := err.Error()
	testhelpers.AssertContains(t, msg, "invalid activity kind", "Error should describe the problem")
	testhelpers.AssertContains(t, msg, "bogus", "Error should contain the invalid value")
}
