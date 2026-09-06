package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
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

	// The emulator answers terminal query sequences (DA1, DECRQM, DSR) via a
	// synchronous io.Pipe that blocks until read. Program output captured from
	// Bubble Tea contains those queries, so the response side must be drained
	// or Write deadlocks (the 2-minute CI hang in TestTeatest_VTScreen).
	go func() {
		buf := make([]byte, 4096)
		for {
			if _, err := term.Read(buf); err != nil {
				return
			}
		}
	}()

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

// TestVTScreen_HandlesQuerySequences regression-guards the response-pipe
// deadlock: captured program output routinely contains terminal query
// sequences (DA1, DECRQM) whose emulator answers block until the response
// side is read. Without the drain goroutine this test hangs to the package
// timeout, exactly as CI did for 100 consecutive runs.
func TestVTScreen_HandlesQuerySequences(t *testing.T) {
	done := make(chan string, 1)
	go func() {
		screen, err := vtScreenFromBytes([]byte("\x1b[c\x1b[?2026$pBuild Module"), 50, 10)
		if err != nil {
			t.Errorf("vt screen reconstruction: %v", err)
		}
		done <- screen
	}()

	select {
	case screen := <-done:
		if !strings.Contains(screen, "Build Module") {
			t.Errorf("VT screen should contain 'Build Module'\n\nScreen:\n%s", screen)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("vtScreenFromBytes deadlocked on query sequences — response pipe not drained")
	}
}
