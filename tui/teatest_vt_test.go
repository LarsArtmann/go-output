package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/charmbracelet/x/vt"
)

// vtScreenFromBytes feeds raw bytes into a VT emulator and returns the
// reconstructed screen. This lets us assert on what the terminal ACTUALLY
// displays after processing cursor movement, erase, and sync sequences.
func vtScreenFromBytes(t *testing.T, raw []byte, width, height int) string {
	t.Helper()

	term := vt.NewEmulator(width, height)

	t.Cleanup(func() { _ = term.Close() })

	if _, err := term.Write(raw); err != nil {
		t.Fatalf("vt write: %v", err)
	}

	screen := term.String()
	lines := strings.Split(screen, "\n")

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		result = append(result, strings.TrimRight(line, " \x00"))
	}

	return strings.Join(result, "\n")
}

// TestTeatest_VTScreen_ShowsActivityLabels feeds teatest's raw diff output
// through a VT emulator and asserts that activity labels appear on the
// reconstructed screen. This is a deeper assertion than ANSI-stripped text —
// it proves the diff renderer's cursor sequences produce correct screen state.
func TestTeatest_VTScreen_ShowsActivityLabels(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	var screenSnapshot string

	teatest.WaitFor(t, tm.Output(), func(raw []byte) bool {
		screenSnapshot = vtScreenFromBytes(t, raw, 100, 30)

		return strings.Contains(screenSnapshot, "Build Module") &&
			strings.Contains(screenSnapshot, "Run Tests")
	}, teatest.WithDuration(5*time.Second), teatest.WithCheckInterval(50*time.Millisecond))

	if !strings.Contains(screenSnapshot, "Build Module") {
		t.Errorf("VT screen should contain 'Build Module'\n\nScreen:\n%s", screenSnapshot)
	}

	if !strings.Contains(screenSnapshot, "Run Tests") {
		t.Errorf("VT screen should contain 'Run Tests'\n\nScreen:\n%s", screenSnapshot)
	}
}
