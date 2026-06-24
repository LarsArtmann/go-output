package nom

import (
	"strings"
	"testing"
	"time"
)

func TestComputePhaseCounts_MaxElapsedNotSum(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), nil)
	dt.AddActivity(ActivityID("a"), []ActivityID{"phase"})
	dt.AddActivity(ActivityID("b"), []ActivityID{"phase"})
	dt.AddActivity(ActivityID("c"), []ActivityID{"phase"})

	snaps := newSnapshotBuilder()
	snaps.setPhase(ActivityID("phase"), "Parallel Steps", ActivityStatusCompleted, 0)
	snaps.set(ActivityID("a"), "A", ActivityStatusCompleted, 5*time.Second)
	snaps.set(ActivityID("b"), "B", ActivityStatusCompleted, 10*time.Second)
	snaps.set(ActivityID("c"), "C", ActivityStatusCompleted, 3*time.Second)

	dt.collapseCompletedPhases = true

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)

	if !strings.Contains(got, "Parallel Steps") {
		t.Errorf("expected collapsed phase label in output, got:\n%s", got)
	}

	if strings.Contains(got, "18") {
		t.Errorf("expected max(10s) not sum(18s) in elapsed, got:\n%s", got)
	}

	if !strings.Contains(got, "10.0s") {
		t.Errorf("expected max elapsed (10.0s) in output, got:\n%s", got)
	}
}

func TestComputePhaseCounts_DoesNotCollapseWhileChildrenRunning(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), nil)
	dt.AddActivity(ActivityID("a"), []ActivityID{"phase"})
	dt.AddActivity(ActivityID("b"), []ActivityID{"phase"})

	snaps := newSnapshotBuilder()
	snaps.setPhase(ActivityID("phase"), "Active Phase", ActivityStatusRunning, 0)
	snaps.set(ActivityID("a"), "A", ActivityStatusFailed, 2*time.Second)
	snaps.set(ActivityID("b"), "B", ActivityStatusRunning, 1*time.Second)

	dt.collapseCompletedPhases = true

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)

	if !strings.Contains(got, "A") || !strings.Contains(got, "B") {
		t.Errorf("expected children visible (not collapsed) while some still running, got:\n%s", got)
	}
}

func TestComputePhaseCounts_CollapsesWithFailedChildren(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), nil)
	dt.AddActivity(ActivityID("a"), []ActivityID{"phase"})
	dt.AddActivity(ActivityID("b"), []ActivityID{"phase"})

	snaps := newSnapshotBuilder()
	snaps.setPhase(ActivityID("phase"), "Failed Phase", ActivityStatusCompleted, 0)
	snaps.set(ActivityID("a"), "A", ActivityStatusCompleted, 3*time.Second)
	snaps.set(ActivityID("b"), "B", ActivityStatusFailed, 7*time.Second)

	dt.collapseCompletedPhases = true

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)

	if !strings.Contains(got, "Failed Phase") {
		t.Errorf("expected collapsed phase label, got:\n%s", got)
	}

	if !strings.Contains(got, "1/2") {
		t.Errorf("expected 1/2 (completed/total) in collapsed label, got:\n%s", got)
	}
}
