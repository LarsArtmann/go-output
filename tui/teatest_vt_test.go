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

// pollTeatestOutput reads a teatest TestModel's output with bounded reads and a
// hard deadline, returning the accumulated bytes. teatest.WaitFor uses io.ReadAll
// which blocks indefinitely when the program writes continuously (the tick loop
// keeps appending to the output buffer, so io.ReadAll never sees EOF). Bounded
// reads drain available bytes and return immediately on the empty-buffer EOF,
// keeping the deadline check live.
func pollTeatestOutput(t *testing.T, tm *teatest.TestModel, cond func([]byte) bool, timeout time.Duration) []byte {
	t.Helper()
	out := tm.Output()
	deadline := time.Now().Add(timeout)
	var accumulated []byte
	for time.Now().Before(deadline) {
		chunk := make([]byte, 8192) //nolint:mnd
		n, _ := out.Read(chunk)     // empty buffer returns (0, io.EOF) — non-blocking
		if n > 0 {
			accumulated = append(accumulated, chunk[:n]...)
		}
		if cond(accumulated) {
			return accumulated
		}
		time.Sleep(50 * time.Millisecond) //nolint:mnd
	}
	t.Fatalf("teatest output condition not met after %s\nLast output:\n%s", timeout, string(accumulated))
	return accumulated
}

// TestTeatest_VTScreen_ShowsActivityLabels feeds teatest's raw diff output
// through a VT emulator and asserts that activity labels appear on the
// reconstructed screen. This is a deeper assertion than ANSI-stripped text —
// it proves the diff renderer's cursor sequences produce correct screen state.
//
// Uses bounded polling (pollTeatestOutput) instead of teatest.WaitFor to avoid
// the io.ReadAll deadlock under -race when the program writes continuously.
func TestTeatest_VTScreen_ShowsActivityLabels(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	capturedOutput := pollTeatestOutput(t, tm, func(raw []byte) bool {
		return strings.Contains(ansi.Strip(string(raw)), "Build Module") &&
			strings.Contains(ansi.Strip(string(raw)), "Run Tests")
	}, 5*time.Second)

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
