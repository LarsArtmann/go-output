package nom

import (
	"testing"
)

func TestActivityID(t *testing.T) {
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

	t.Run("valid input", func(t *testing.T) {
		t.Parallel()

		got, err := ParseActivityID("build")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != ActivityID("build") {
			t.Errorf("ParseActivityID(%q) = %q, want %q", "build", got, "build")
		}
	})

	t.Run("empty input returns error", func(t *testing.T) {
		t.Parallel()

		_, err := ParseActivityID("")
		if err == nil {
			t.Error("expected error for empty input")
		}
	})
}

func TestWorkflowID(t *testing.T) {
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

	t.Run("valid input", func(t *testing.T) {
		t.Parallel()

		got, err := ParseWorkflowID("deploy")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != WorkflowID("deploy") {
			t.Errorf("ParseWorkflowID(%q) = %q, want %q", "deploy", got, "deploy")
		}
	})

	t.Run("empty input returns error", func(t *testing.T) {
		t.Parallel()

		_, err := ParseWorkflowID("")
		if err == nil {
			t.Error("expected error for empty input")
		}
	})
}

func TestActivityName(t *testing.T) {
	t.Parallel()

	got := NewActivityName("compile")
	if got.String() != "compile" {
		t.Errorf("ActivityName.String() = %q, want %q", got.String(), "compile")
	}
}

func TestWorkflowName(t *testing.T) {
	t.Parallel()

	got := NewWorkflowName("CI Pipeline")
	if got.String() != "CI Pipeline" {
		t.Errorf("WorkflowName.String() = %q, want %q", got.String(), "CI Pipeline")
	}
}
