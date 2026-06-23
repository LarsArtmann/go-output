package tui

import (
	"testing"
)

func TestUpdateType(t *testing.T) {
	t.Parallel()

	if progressUpdate != 0 {
		t.Errorf("progressUpdate = %d, want 0", progressUpdate)
	}

	if messageUpdate != 1 {
		t.Errorf("messageUpdate = %d, want 1", messageUpdate)
	}
}

func TestProgressUpdateMsg(t *testing.T) {
	t.Parallel()

	msg := progressUpdateMsg{
		Type:     progressUpdate,
		Progress: 50.0,
		Message:  "halfway",
	}

	if msg.Type != progressUpdate {
		t.Errorf("Type = %d, want %d", msg.Type, progressUpdate)
	}

	if msg.Progress != 50.0 {
		t.Errorf("Progress = %f, want 50.0", msg.Progress)
	}

	if msg.Message != "halfway" {
		t.Errorf("Message = %q, want %q", msg.Message, "halfway")
	}
}
