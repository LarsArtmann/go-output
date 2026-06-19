package tui

import (
	"testing"
)

func TestConstants(t *testing.T) {
	t.Parallel()

	if timingFormatWithIcon == "" {
		t.Error("timingFormatWithIcon should not be empty")
	}

	if msgNoActivitiesToDisplay == "" {
		t.Error("msgNoActivitiesToDisplay should not be empty")
	}
}
