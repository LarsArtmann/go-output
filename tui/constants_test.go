package tui

import (
	"testing"
)

func TestConstants(t *testing.T) {
	t.Parallel()

	if timingFormatWithIcon == "" {
		t.Error("timingFormatWithIcon should not be empty")
	}

	if SeparatorLineEquals == "" {
		t.Error("SeparatorLineEquals should not be empty")
	}

	if MsgNoActivitiesToDisplay == "" {
		t.Error("MsgNoActivitiesToDisplay should not be empty")
	}
}
