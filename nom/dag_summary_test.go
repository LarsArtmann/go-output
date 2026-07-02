package nom

import (
	"strings"
	"testing"
	"time"
)

func TestDAGSummary_Empty(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	summary := dt.DAGSummary()

	if summary.Nodes != 0 {
		t.Errorf("Nodes = %d, want 0", summary.Nodes)
	}

	if summary.String() != "0 nodes · 0 edges" {
		t.Errorf("String() = %q, want %q", summary.String(), "0 nodes · 0 edges")
	}
}

func TestDAGSummary_BasicDAG(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root"), nil)
	_ = dt.AddActivity(ActivityID("a"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("b"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("leaf"), []ActivityID{"a", "b"})

	_ = dt.GetRootNodes()

	summary := dt.DAGSummary()

	if summary.Nodes != 4 {
		t.Errorf("Nodes = %d, want 4", summary.Nodes)
	}

	if summary.Edges != 4 {
		t.Errorf("Edges = %d, want 4", summary.Edges)
	}

	if summary.MaxDepth != 2 {
		t.Errorf("MaxDepth = %d, want 2", summary.MaxDepth)
	}

	if summary.Roots != 1 {
		t.Errorf("Roots = %d, want 1", summary.Roots)
	}

	if summary.MaxWidth < 2 {
		t.Errorf("MaxWidth = %d, want >= 2", summary.MaxWidth)
	}

	s := summary.String()
	if !strings.Contains(s, "4 nodes") {
		t.Errorf("String missing node count: %q", s)
	}

	if !strings.Contains(s, "4 edges") {
		t.Errorf("String missing edge count: %q", s)
	}

	if !strings.Contains(s, "layers") {
		t.Errorf("String missing layers: %q", s)
	}
}

func TestDAGSummary_WithSnapshots(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("a"), nil)
	_ = dt.AddActivity(ActivityID("b"), []ActivityID{"a"})

	_ = dt.GetRootNodes()

	snaps := newSnapshotBuilder()
	snaps.setWithEstimate(ActivityID("a"), "A", ActivityStatusRunning, 0, 10*time.Second)
	snaps.setWithEstimate(ActivityID("b"), "B", ActivityStatusPending, 0, 5*time.Second)

	summary := dt.DAGSummaryWithSnapshots(snaps.snaps)

	if summary.Critical <= 0 {
		t.Errorf("Critical = %v, want > 0", summary.Critical)
	}
}

func TestDAGSummary_PhasesCount(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	_ = dt.AddActivity(ActivityID("root"), nil)
	_ = dt.AddActivity(ActivityID("a"), []ActivityID{"root"})
	_ = dt.AddActivity(ActivityID("b"), []ActivityID{"root"})

	_ = dt.GetRootNodes()

	snaps := newSnapshotBuilder()
	snaps.setPhase(ActivityID("root"), "Root Phase", ActivityStatusRunning, 0)
	snaps.set(ActivityID("a"), "A", ActivityStatusPending, 0)
	snaps.set(ActivityID("b"), "B", ActivityStatusPending, 0)

	summary := dt.DAGSummaryWithSnapshots(snaps.snaps)

	if summary.Phases != 1 {
		t.Errorf("Phases = %d, want 1", summary.Phases)
	}

	s := summary.String()
	if !strings.Contains(s, "1 phases") {
		t.Errorf("String() missing phase count: %q", s)
	}
}

func TestRenderMode_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		mode RenderMode
		want string
	}{
		{RenderModeTree, "tree"},
		{RenderModeLayered, "layered"},
		{RenderMode(999), "unknown"},
	}

	for _, tc := range tests {
		got := tc.mode.String()
		if got != tc.want {
			t.Errorf("%d.String() = %q, want %q", tc.mode, got, tc.want)
		}
	}
}
