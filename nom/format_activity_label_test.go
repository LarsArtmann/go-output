package nom

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
)

func TestFormatActivityLabel_SlowStepEscalation(t *testing.T) {
	gomega.RegisterTestingT(t)

	baseSnap := func(elapsed time.Duration) ActivitySnapshot {
		return ActivitySnapshot{
			Label:          "goimports",
			Symbol:         SymbolCompleted,
			Color:          Colors.Completed,
			Status:         ActivityStatusCompleted,
			CurrentElapsed: elapsed,
		}
	}

	t.Run("9.9s: no escalation color", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		snap := baseSnap(9*time.Second + 900*time.Millisecond)

		display, c := formatActivityLabel(snap)

		gomega.Expect(display).To(gomega.ContainSubstring("goimports"))
		gomega.Expect(c).To(gomega.Equal(Colors.Completed))
	})

	t.Run("10.0s: yellow escalation threshold", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		snap := baseSnap(10 * time.Second)

		display, _ := formatActivityLabel(snap)

		// At >= 10s, the timing info should have ANSI styling (faint yellow).
		gomega.Expect(display).To(gomega.ContainSubstring("goimports"))
		gomega.Expect(display).To(gomega.ContainSubstring("\x1b["))
	})

	t.Run("29.9s: still yellow (below red threshold)", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		snap := baseSnap(29*time.Second + 900*time.Millisecond)

		display, _ := formatActivityLabel(snap)

		gomega.Expect(display).To(gomega.ContainSubstring("goimports"))
		gomega.Expect(display).To(gomega.ContainSubstring("\x1b["))
	})

	t.Run("30.0s: red escalation threshold", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		snap := baseSnap(30 * time.Second)

		display, _ := formatActivityLabel(snap)

		gomega.Expect(display).To(gomega.ContainSubstring("goimports"))
		gomega.Expect(display).To(gomega.ContainSubstring("\x1b["))
	})

	t.Run("60s: still red", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		snap := baseSnap(60 * time.Second)

		display, _ := formatActivityLabel(snap)

		gomega.Expect(display).To(gomega.ContainSubstring("goimports"))
		gomega.Expect(display).To(gomega.ContainSubstring("\x1b["))
	})

	t.Run("0s: no timing info shown for instant steps", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		snap := baseSnap(0)

		display, c := formatActivityLabel(snap)

		gomega.Expect(display).To(gomega.ContainSubstring("goimports"))
		gomega.Expect(c).To(gomega.Equal(Colors.Completed))
	})
}
