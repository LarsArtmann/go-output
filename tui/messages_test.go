package tui

import (
	"testing"
)

func TestUpdateType(t *testing.T) {
	t.Parallel()

	if ProgressUpdate != 0 {
		t.Errorf("ProgressUpdate = %d, want 0", ProgressUpdate)
	}

	if MessageUpdate != 1 {
		t.Errorf("MessageUpdate = %d, want 1", MessageUpdate)
	}

	if StepUpdate != 2 {
		t.Errorf("StepUpdate = %d, want 2", StepUpdate)
	}
}

func TestProgressUpdateMsg(t *testing.T) {
	t.Parallel()

	msg := ProgressUpdateMsg{
		Type:     ProgressUpdate,
		Progress: 50.0,
		Message:  "halfway",
		Current:  5,
		Total:    10,
	}

	if msg.Type != ProgressUpdate {
		t.Errorf("Type = %d, want %d", msg.Type, ProgressUpdate)
	}

	if msg.Progress != 50.0 {
		t.Errorf("Progress = %f, want 50.0", msg.Progress)
	}

	if msg.Message != "halfway" {
		t.Errorf("Message = %q, want %q", msg.Message, "halfway")
	}
}
