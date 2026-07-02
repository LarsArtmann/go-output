package nom

import (
	"testing"
	"time"
)

func TestComputeCriticalPath_Diamond(t *testing.T) {
	t.Parallel()

	// Diamond: A → B, A → C, B → D, C → D.
	// A is root, B/C are intermediate, D is the sink.
	// All estimates are 1s. The critical path includes A and D (both paths to D
	// are equal, so the back-track marks both B and C as well because they both
	// satisfy depTotal + weight == nodeTotal).
	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("b"), []ActivityID{"a"})
	dt.AddActivity(ActivityID("c"), []ActivityID{"a"})
	dt.AddActivity(ActivityID("d"), []ActivityID{"b", "c"})

	snaps := newSnapshotBuilder()
	snaps.setWithEstimate(ActivityID("a"), "A", ActivityStatusCompleted, time.Second, time.Second)
	snaps.setWithEstimate(ActivityID("b"), "B", ActivityStatusCompleted, time.Second, time.Second)
	snaps.setWithEstimate(ActivityID("c"), "C", ActivityStatusCompleted, time.Second, time.Second)
	snaps.setWithEstimate(ActivityID("d"), "D", ActivityStatusPending, 0, time.Second)

	dt.mu.RLock()
	cp := dt.computeCriticalPath(snaps.snaps)
	dt.mu.RUnlock()

	if !cp[ActivityID("a")] {
		t.Errorf("expected a on critical path")
	}

	if !cp[ActivityID("d")] {
		t.Errorf("expected d on critical path")
	}
}

func TestComputeCriticalPath_LongestBranchWins(t *testing.T) {
	t.Parallel()

	// A → B → D (3s path) and A → C (1s path). D is the critical leaf.
	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("b"), []ActivityID{"a"})
	dt.AddActivity(ActivityID("c"), []ActivityID{"a"})
	dt.AddActivity(ActivityID("d"), []ActivityID{"b"})

	snaps := newSnapshotBuilder()
	snaps.setWithEstimate(ActivityID("a"), "A", ActivityStatusCompleted, time.Second, time.Second)
	snaps.setWithEstimate(ActivityID("b"), "B", ActivityStatusCompleted, 2*time.Second, 2*time.Second)
	snaps.setWithEstimate(ActivityID("c"), "C", ActivityStatusCompleted, time.Second, time.Second)
	snaps.setWithEstimate(ActivityID("d"), "D", ActivityStatusPending, 0, time.Second)

	dt.mu.RLock()
	cp := dt.computeCriticalPath(snaps.snaps)
	dt.mu.RUnlock()

	if !cp[ActivityID("a")] || !cp[ActivityID("b")] || !cp[ActivityID("d")] {
		t.Errorf("expected a, b, d on critical path, got %v", cp)
	}

	if cp[ActivityID("c")] {
		t.Errorf("c should NOT be on the critical path")
	}
}

func TestEstimatedCriticalPathRemaining(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("b"), []ActivityID{"a"})
	dt.AddActivity(ActivityID("c"), []ActivityID{"a"})
	dt.AddActivity(ActivityID("d"), []ActivityID{"b", "c"})

	snaps := newSnapshotBuilder()
	snaps.setWithEstimate(ActivityID("a"), "A", ActivityStatusCompleted, time.Second, time.Second)
	snaps.setWithEstimate(ActivityID("b"), "B", ActivityStatusRunning, 500*time.Millisecond, 2*time.Second)
	snaps.setWithEstimate(ActivityID("c"), "C", ActivityStatusPending, 0, time.Second)
	snaps.setWithEstimate(ActivityID("d"), "D", ActivityStatusPending, 0, time.Second)

	// Longest remaining path: b (1.5s remaining) + d (1s) = 2.5s, or c (1s) + d (1s) = 2s.
	want := 2500 * time.Millisecond
	if got := dt.EstimatedCriticalPathRemaining(snaps.snaps); got != want {
		t.Errorf("EstimatedCriticalPathRemaining = %v, want %v", got, want)
	}
}
