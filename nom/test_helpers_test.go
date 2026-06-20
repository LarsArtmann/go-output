package nom

import (
	"path/filepath"
	"testing"
	"time"
)

func newTestSubscriber(t *testing.T) *NOMStyleSubscriber {
	t.Helper()
	ns := NewNOMStyleSubscriber(WithCachePath(filepath.Join(t.TempDir(), cacheFilename)))
	t.Cleanup(ns.timingCache.waitPendingSaves)

	return ns
}

// snapshotBuilder accumulates ActivitySnapshot values for test rendering
// without a subscriber. Replaces the old pattern of mutating node.Status/
// node.Symbol directly on the shared *Activity pointer.
type snapshotBuilder struct {
	snaps map[ActivityID]ActivitySnapshot
}

func newSnapshotBuilder() *snapshotBuilder {
	return &snapshotBuilder{snaps: make(map[ActivityID]ActivitySnapshot)}
}

func (b *snapshotBuilder) set(id ActivityID, label string, status ActivityStatus, elapsed time.Duration) {
	b.snaps[id] = ActivitySnapshot{
		Label:          label,
		Status:         status,
		Symbol:         status.GetSymbol(),
		Color:          status.GetColor(),
		CurrentElapsed: elapsed,
	}
}

func (b *snapshotBuilder) snapshot(id ActivityID) ActivitySnapshot {
	return b.snaps[id]
}

func (b *snapshotBuilder) has(id ActivityID) bool {
	_, ok := b.snaps[id]
	return ok
}
