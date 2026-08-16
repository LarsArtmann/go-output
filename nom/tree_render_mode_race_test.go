package nom

import (
	"testing"
	"time"
)

// TestRenderMode_ConcurrentToggleAndRender guards the data race between
// SetRenderMode (writer, e.g. the TUI update loop toggling display mode) and
// the render entry points (readers, e.g. the tick loop). Before the fix, the
// render dispatch sites read dt.renderMode without holding dt.mu. Run with
// -race to verify.
func TestRenderMode_ConcurrentToggleAndRender(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root1"), nil)
	_ = dt.AddActivity(ActivityID("child1"), []ActivityID{"root1"})
	_ = dt.AddActivity(ActivityID("grandchild1"), []ActivityID{"child1"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root1"), "Root One", ActivityStatusRunning, 2*time.Second)
	snaps.set(ActivityID("child1"), "Child One", ActivityStatusCompleted, 1*time.Second)
	snaps.set(ActivityID("grandchild1"), "Grandchild One", ActivityStatusPending, 0)

	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)

		for i := 0; i < 500; i++ {
			dt.SetRenderMode(RenderModeLayered)
			dt.SetRenderMode(RenderModeTree)
		}
	}()

	entries := dt.VisibleEntriesWithSnapshots(snaps.snaps, 0)
	for i := 0; i < 500; i++ {
		_ = dt.RenderWithSnapshots(snaps.snaps, 0, 80)
		_ = dt.VisibleEntriesWithSnapshots(snaps.snaps, 0)

		for _, entry := range entries {
			_ = dt.RenderVisibleEntry(entry, snaps.snaps, 80)
		}
	}

	<-writerDone
}
