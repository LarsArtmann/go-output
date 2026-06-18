package nom

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestActivityID_Methods(t *testing.T) {
	t.Parallel()

	t.Run("String returns underlying value", func(t *testing.T) {
		t.Parallel()

		id := ActivityID("build")
		if id.String() != "build" {
			t.Errorf("ActivityID.String() = %q, want %q", id.String(), "build")
		}
	})

	t.Run("IsZero returns true for empty", func(t *testing.T) {
		t.Parallel()

		if !ActivityID("").IsZero() {
			t.Error("empty ActivityID should be zero")
		}
	})

	t.Run("IsZero returns false for non-empty", func(t *testing.T) {
		t.Parallel()

		if ActivityID("x").IsZero() {
			t.Error("non-empty ActivityID should not be zero")
		}
	})
}

func TestNewActivityID(t *testing.T) {
	t.Parallel()

	got := NewActivityID("test")
	if got != ActivityID("test") {
		t.Errorf("NewActivityID(%q) = %q, want %q", "test", got, "test")
	}
}

func TestParseActivityID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    ActivityID
		wantErr bool
	}{
		{name: "valid input", input: "build", want: ActivityID("build"), wantErr: false},
		{name: "empty input returns error", input: "", want: ActivityID(""), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseActivityID(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			testhelpers.AssertEqual(t, "ParseActivityID", tt.input, got, tt.want)
		})
	}
}

func TestWorkflowID_Methods(t *testing.T) {
	t.Parallel()

	t.Run("String returns underlying value", func(t *testing.T) {
		t.Parallel()

		id := WorkflowID("ci-pipeline")
		if id.String() != "ci-pipeline" {
			t.Errorf("WorkflowID.String() = %q, want %q", id.String(), "ci-pipeline")
		}
	})

	t.Run("IsZero returns true for empty", func(t *testing.T) {
		t.Parallel()

		if !WorkflowID("").IsZero() {
			t.Error("empty WorkflowID should be zero")
		}
	})
}

func TestNewWorkflowID(t *testing.T) {
	t.Parallel()

	got := NewWorkflowID("wf")
	if got != WorkflowID("wf") {
		t.Errorf("NewWorkflowID(%q) = %q, want %q", "wf", got, "wf")
	}
}

func TestParseWorkflowID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    WorkflowID
		wantErr bool
	}{
		{name: "valid input", input: "deploy", want: WorkflowID("deploy"), wantErr: false},
		{name: "empty input returns error", input: "", want: WorkflowID(""), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseWorkflowID(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			testhelpers.AssertEqual(t, "ParseWorkflowID", tt.input, got, tt.want)
		})
	}
}

func TestActivityName_String(t *testing.T) {
	t.Parallel()

	got := NewActivityName("compile")
	if got.String() != "compile" {
		t.Errorf("ActivityName.String() = %q, want %q", got.String(), "compile")
	}
}

func TestWorkflowName_String(t *testing.T) {
	t.Parallel()

	got := NewWorkflowName("CI Pipeline")
	if got.String() != "CI Pipeline" {
		t.Errorf("WorkflowName.String() = %q, want %q", got.String(), "CI Pipeline")
	}
}
