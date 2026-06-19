package nom

import (
	"path/filepath"
	"testing"
	"time"
)

// newTestSubscriber returns a NOMStyleSubscriber whose timing cache is isolated
// to a per-test temp directory, so no test ever reads or writes the real
// ~/.cache/nom-timing.csv. It also drains pending async saves on cleanup.
func newTestSubscriber(t *testing.T) *NOMStyleSubscriber {
	t.Helper()
	ns := NewNOMStyleSubscriber(WithCachePath(filepath.Join(t.TempDir(), cacheFilename)))
	t.Cleanup(ns.timingCache.waitPendingSaves)

	return ns
}

// testSetStatus is a test helper that mutates the shared Activity pointer
// directly, replacing the old UpdateActivityStatus API. Symbol/color are
// derived from status via applyVisualStyle.
func testSetStatus(dt *DependencyTree, id ActivityID, status ActivityStatus, startTime time.Time) {
	node := dt.GetNode(id)
	if node == nil {
		return
	}

	node.Status = status
	node.applyVisualStyle()
	node.StartTime = startTime
}
