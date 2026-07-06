package tui

import (
	"testing"
)

func TestConstants(t *testing.T) {
	t.Parallel()

	if timingFormat == "" {
		t.Error("timingFormat should not be empty")
	}

	if MsgNoActivities == "" {
		t.Error("MsgNoActivities should not be empty")
	}
}
