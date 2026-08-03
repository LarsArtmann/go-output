package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/teatest/v2"
	"github.com/charmbracelet/x/vt"
)

// vtScreenFromBytes feeds raw bytes into a VT emulator and returns the
// reconstructed screen. This lets us assert on what the terminal ACTUALLY
// displays after processing cursor movement, erase, and sync sequences.
//
// This is a pure function (no *testing.T dependency) so it can be safely
// called from any context, including teatest polling callbacks.
func vtScreenFromBytes(raw []byte, width, height int) (string, error) {
	term := vt.NewEmulator(width, height)
	defer func() { _ = term.Close() }()

	if _, err := term.Write(raw); err != nil {
		return "", fmt.Errorf("vt write: %w", err)
	}

	screen := term.String()
	lines := strings.Split(screen, "\n")

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		result = append(result, strings.TrimRight(line, " \x00"))
	}

	return strings.Join(result, "\n"), nil
}

// TestTeatest_VTScreen_ShowsActivityLabels feeds teatest's raw diff output
// through a VT emulator and asserts that activity labels appear on the
// reconstructed screen. This is a deeper assertion than ANSI-stripped text —
// it proves the diff renderer's cursor sequences produce correct screen state.
//
// The wait uses ANSI-strip (the proven stable approach) to know when content
// is ready, capturing the raw bytes during polling. A SINGLE VT reconstruction
// is done afterward — never inside the polling loop, which created dozens of
// VT emulators under -race and caused CI deadlocks.
func TestTeatest_VTScreen_ShowsActivityLabels(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	var capturedOutput []byte

	teatest.WaitFor(t, tm.Output(), func(raw []byte) bool {
		capturedOutput = raw

		return strings.Contains(ansi.Strip(string(raw)), "Build Module") &&
			strings.Contains(ansi.Strip(string(raw)), "Run Tests")
	}, teatest.WithDuration(5*time.Second), teatest.WithCheckInterval(50*time.Millisecond))

	screenSnapshot, err := vtScreenFromBytes(capturedOutput, 100, 30)
	if err != nil {
		t.Fatalf("vt screen reconstruction: %v", err)
	}

	if !strings.Contains(screenSnapshot, "Build Module") {
		t.Errorf("VT screen should contain 'Build Module'\n\nScreen:\n%s", screenSnapshot)
	}

	if !strings.Contains(screenSnapshot, "Run Tests") {
		t.Errorf("VT screen should contain 'Run Tests'\n\nScreen:\n%s", screenSnapshot)
	}
}
