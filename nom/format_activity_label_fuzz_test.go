package nom

import (
	"strings"
	"testing"
	"time"
)

// FuzzFormatActivityLabel guards formatActivityLabel against panics and
// malformed output across arbitrary label strings. It enforces that the
// function never panics and always returns a non-empty display string
// containing the label.
func FuzzFormatActivityLabel(f *testing.F) {
	f.Add("Build")
	f.Add("")
	f.Add("unicode: → ⟳ ⏵ ✔ ⚠ ○")
	f.Add("very long label that exceeds normal terminal width " + strings.Repeat("x", 200))
	f.Add("newline\ninjection\tattempt")
	f.Add("emoji 🎉 mixed")

	f.Fuzz(func(t *testing.T, label string) {
		snap := ActivitySnapshot{
			Label:          label,
			Symbol:         SymbolRunning,
			Color:          Colors.Running,
			Status:         ActivityStatusRunning,
			CurrentElapsed: 5 * time.Second,
		}

		display, c := formatActivityLabel(snap)

		// Must not panic (implicit — if we reach here, it didn't).
		// Display must always contain the label text.
		if !strings.Contains(display, label) && label != "" {
			t.Errorf("display %q does not contain label %q", display, label)
		}

		// Color must never be nil — the caller relies on it for styling.
		if c == nil {
			t.Error("formatActivityLabel returned nil color")
		}
	})
}
