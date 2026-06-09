package tui

import (
	"testing"
)

func TestDisplayMode_Values(t *testing.T) {
	t.Parallel()

	if DisplayModeNOM != "nom" {
		t.Errorf("DisplayModeNOM = %q, want %q", DisplayModeNOM, "nom")
	}

	if DisplayModeUniversal != "universal" {
		t.Errorf("DisplayModeUniversal = %q, want %q", DisplayModeUniversal, "universal")
	}
}
