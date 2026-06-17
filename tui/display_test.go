package tui

import (
	"testing"
)

func TestDisplayMode_Values(t *testing.T) {
	t.Parallel()

	if DisplayModeNOM == DisplayModeUniversal {
		t.Error("DisplayModeNOM and DisplayModeUniversal should be distinct values")
	}

	if DisplayModeUniversal != DisplayMode(0) {
		t.Errorf("DisplayModeUniversal = %d, want 0", DisplayModeUniversal)
	}

	if DisplayModeNOM != DisplayMode(1) {
		t.Errorf("DisplayModeNOM = %d, want 1", DisplayModeNOM)
	}
}
