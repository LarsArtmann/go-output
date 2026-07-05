package nom

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/vt"
)

// vtHarness wraps a vt.Emulator as an io.Writer and provides screen assertion
// helpers. The InlineRenderer writes escape sequences to this harness; the
// underlying VT emulator processes them into a 2D screen buffer that we assert
// against — testing what a REAL terminal would display, not just which escape
// codes were emitted.
type vtHarness struct {
	term            *vt.Emulator
	cursorHidden    bool
	syncWasActive   bool // tracks if sync mode was EVER activated (set+reset in one write)
	syncToggleCount int  // number of sync mode transitions
}

// newVTHarness creates a VT emulator of the given size with cursor-visibility
// and sync-output-mode tracking callbacks installed. The harness implements
// io.Writer so it can be passed directly to NewInlineRenderer.
func newVTHarness(t *testing.T, width, height int) *vtHarness {
	t.Helper()

	term := vt.NewEmulator(width, height)

	t.Cleanup(func() { _ = term.Close() })

	h := &vtHarness{term: term}

	term.SetCallbacks(vt.Callbacks{
		CursorVisibility: func(visible bool) {
			h.cursorHidden = !visible
		},
		EnableMode: func(mode ansi.Mode) {
			if mode == ansi.ModeSynchronizedOutput {
				h.syncWasActive = true
				h.syncToggleCount++
			}
		},
		DisableMode: func(mode ansi.Mode) {
			if mode == ansi.ModeSynchronizedOutput {
				h.syncToggleCount++
			}
		},
	})

	return h
}

// Write implements io.Writer — feeds bytes into the VT emulator.
func (h *vtHarness) Write(p []byte) (int, error) {
	n, err := h.term.Write(p)
	if err != nil {
		return n, fmt.Errorf("vt emulator write: %w", err)
	}

	return n, nil
}

// screenString returns the full visible screen as a trimmed multi-line string.
// Trailing spaces and null characters on each line are removed; fully-empty
// trailing lines are dropped.
func (h *vtHarness) screenString() string {
	raw := h.term.String()
	lines := strings.Split(raw, "\n")

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \x00")
		result = append(result, trimmed)
	}

	return strings.Join(result, "\n")
}

// nonEmptyLines returns only screen lines that contain visible content.
func (h *vtHarness) nonEmptyLines() []string {
	raw := h.term.String()
	lines := strings.Split(raw, "\n")

	result := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " \x00")
		if strings.TrimSpace(trimmed) != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// assertScreenContains asserts the rendered screen contains substr.
func (h *vtHarness) assertScreenContains(t *testing.T, substr string) {
	t.Helper()

	screen := h.screenString()
	if !strings.Contains(screen, substr) {
		t.Errorf("screen should contain %q\n\nFull screen:\n%s", substr, screen)
	}
}
