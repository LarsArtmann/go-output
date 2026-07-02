package nom

import (
	"bytes"
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

// newInlineTestRenderer creates an InlineRenderer configured for unit testing
// against a bytes.Buffer. Production detection treats non-TTY buffers as
// plain-text (no cursor codes). Tests that need to exercise the inline
// rendering path (cursor-up, clear-line, sync-output codes) must bypass the
// authoritative plainText gate by setting writerIsTTY=true directly — this
// simulates what a real terminal would provide, without requiring an actual TTY.
func newInlineTestRenderer(sub *NOMStyleSubscriber, buf *bytes.Buffer, maxHeight int) *InlineRenderer {
	r := NewInlineRenderer(sub, buf, maxHeight)
	r.SetNoColor(true) // deterministic output (no terminal color codes)

	// Force the inline rendering path for testing: pretend the buffer is a TTY.
	// snapshotConfig() computes effectivePlainText = plainText || !writerIsTTY,
	// so writerIsTTY=true + plainText=false (default from NewInlineRenderer with
	// non-CI env) gives us the inline path.
	r.tickMu.Lock()
	r.writerIsTTY = true
	r.plainText = false
	r.tickMu.Unlock()

	return r
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
		Kind:           ActivityKindTask,
		Label:          label,
		Status:         status,
		Symbol:         status.GetSymbol(),
		Color:          status.GetColor(),
		CurrentElapsed: elapsed,
	}
}

func (b *snapshotBuilder) setWithEstimate(
	id ActivityID,
	label string,
	status ActivityStatus,
	elapsed, estimated time.Duration,
) {
	b.set(id, label, status, elapsed)
	s := b.snaps[id]
	s.EstimatedTime = estimated
	b.snaps[id] = s
}

func (b *snapshotBuilder) setPhase(id ActivityID, label string, status ActivityStatus, elapsed time.Duration) {
	b.set(id, label, status, elapsed)
	s := b.snaps[id]
	s.Kind = ActivityKindPhase
	b.snaps[id] = s
}

func (b *snapshotBuilder) setCategory(id ActivityID, label string, status ActivityStatus, category ActivityCategory) {
	b.set(id, label, status, 0)
	s := b.snaps[id]
	s.Category = category
	b.snaps[id] = s
}
