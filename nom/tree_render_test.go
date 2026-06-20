package nom

import (
	"strings"
	"testing"
	"time"
)

func TestDependencyTree_TreePrefix_RootNode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), nil, nil)
	dt.AddActivity(ActivityID("child"), nil, []ActivityID{"root"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root", ActivityStatusRunning, 0)
	snaps.set(ActivityID("child"), "Child", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
	if got == "" {
		t.Error("Render should produce output")
	}
}

func TestDependencyTree_Render_PendingStatus(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil, nil)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("a"), "A", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
	if got == "" {
		t.Error("Render should produce output for pending status")
	}
}

func TestDependencyTree_Render_FailedPriority(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil, nil)
	dt.AddActivity(ActivityID("b"), nil, nil)
	dt.AddActivity(ActivityID("c"), nil, nil)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("a"), "A", ActivityStatusCompleted, 0)
	snaps.set(ActivityID("b"), "B", ActivityStatusFailed, 0)
	snaps.set(ActivityID("c"), "C", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 3, 0)
	if got == "" {
		t.Error("Render should produce output")
	}
}

func TestDependencyTree_AddActivity_WithNonExistentDependency(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	err := dt.AddActivity(ActivityID("child"), nil, []ActivityID{"nonexistent"})
	if err != nil {
		t.Fatalf("AddActivity() error: %v", err)
	}

	parent := dt.GetNode(ActivityID("nonexistent"))
	if parent == nil {
		t.Error("nonexistent dependency should be auto-created")
	}

	child := dt.GetNode(ActivityID("child"))
	assertChildParentID(t, child, "nonexistent")
}

func TestDependencyTree_AddActivity_UpdateExisting(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil, nil)
	dt.AddActivity(ActivityID("a"), nil, nil)

	node := dt.GetNode(ActivityID("a"))
	if node == nil {
		t.Fatal("node should exist")
	}
	if node.ID != ActivityID("a") {
		t.Errorf("ID = %q, want %q", node.ID, "a")
	}
}

func TestDependencyTree_Render_SecondaryDependencies(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), nil, nil)
	dt.AddActivity(ActivityID("step1"), nil, []ActivityID{"phase"})
	dt.AddActivity(ActivityID("step2"), nil, []ActivityID{"phase", "step1"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("phase"), "Phase", ActivityStatusPending, 0)
	snaps.set(ActivityID("step1"), "Step1", ActivityStatusPending, 0)
	snaps.set(ActivityID("step2"), "Step2", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
	if got == "" {
		t.Fatal("Render should produce output")
	}

	if !strings.Contains(got, "depends on") {
		t.Errorf("render should contain dependency annotation for secondary deps, got:\n%s", got)
	}

	step2Node := dt.GetNode(ActivityID("step2"))
	if len(step2Node.SecondaryParents) != 1 || step2Node.SecondaryParents[0] != ActivityID("step1") {
		t.Errorf("SecondaryParents = %v, want [step1]", step2Node.SecondaryParents)
	}
}

func TestDependencyTree_Render_PhaseStyling(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase:build"), nil, nil)
	dt.AddActivity(ActivityID("compile"), nil, []ActivityID{"phase:build"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("phase:build"), "Build", ActivityStatusRunning, 0)
	snaps.set(ActivityID("compile"), "Compile", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
	if got == "" {
		t.Fatal("Render should produce output")
	}

	if !strings.Contains(got, string(SymbolPhase)) {
		t.Errorf("render should contain phase symbol %q, got:\n%s", SymbolPhase, got)
	}
}

func TestDependencyTree_Render_PriorityOrdering(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase:build"), nil, nil)
	dt.AddActivity(ActivityID("compile"), nil, []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("test"), nil, []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("lint"), nil, []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("deploy"), nil, []ActivityID{"phase:build"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("phase:build"), "Build Phase", ActivityStatusRunning, 5*time.Second)
	snaps.set(ActivityID("compile"), "Compile", ActivityStatusCompleted, 2*time.Second)
	snaps.set(ActivityID("test"), "Run Tests", ActivityStatusRunning, 3*time.Second)
	snaps.set(ActivityID("lint"), "Lint Code", ActivityStatusPending, 0)
	snaps.set(ActivityID("deploy"), "Deploy", ActivityStatusFailed, 1*time.Second)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)

	failedIdx := strings.Index(got, "Deploy")
	runningIdx := strings.Index(got, "Run Tests")
	pendingIdx := strings.Index(got, "Lint Code")
	completedIdx := strings.Index(got, "Compile")

	for _, idx := range []int{failedIdx, runningIdx, pendingIdx, completedIdx} {
		if idx == -1 {
			t.Fatalf("expected all four activities in output, got:\n%s", got)
		}
	}

	if failedIdx >= runningIdx || runningIdx >= pendingIdx || pendingIdx >= completedIdx {
		t.Errorf("activities not in priority order; indices failed=%d running=%d pending=%d completed=%d\n%s",
			failedIdx, runningIdx, pendingIdx, completedIdx, got)
	}

	limited := dt.RenderWithSnapshots(snaps.snaps, 3, 0)
	if !strings.Contains(limited, "Deploy") {
		t.Errorf("limited render should keep failed activity, got:\n%s", limited)
	}
	if !strings.Contains(limited, "Run Tests") {
		t.Errorf("limited render should keep running activity, got:\n%s", limited)
	}
	if strings.Contains(limited, "Compile") {
		t.Errorf("limited render should drop completed activity, got:\n%s", limited)
	}
}

func TestDependencyTree_RenderWithWidth_TruncatesLongNames(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), nil, nil)

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("a"), "This is a very long activity name that will not fit", ActivityStatusRunning, 0)

	wide := dt.RenderWithSnapshots(snaps.snaps, 10, 80)
	if strings.Contains(wide, "…") {
		t.Errorf("wide render should not truncate, got:\n%s", wide)
	}

	narrow := dt.RenderWithSnapshots(snaps.snaps, 10, 20)
	if !strings.Contains(narrow, "…") {
		t.Errorf("narrow render should truncate with ellipsis, got:\n%s", narrow)
	}
}

func TestDependencyTree_RenderWithWidth_DeepNestingFitsMaxWidth(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), nil, nil)
	dt.AddActivity(ActivityID("c1"), nil, []ActivityID{"root"})
	dt.AddActivity(ActivityID("c2"), nil, []ActivityID{"c1"})
	dt.AddActivity(ActivityID("c3"), nil, []ActivityID{"c2"})
	dt.AddActivity(ActivityID("c4"), nil, []ActivityID{"c3"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("root"), "Root", ActivityStatusPending, 0)
	snaps.set(ActivityID("c1"), "Child1", ActivityStatusPending, 0)
	snaps.set(ActivityID("c2"), "Child2", ActivityStatusPending, 0)
	snaps.set(ActivityID("c3"), "Child3", ActivityStatusPending, 0)
	snaps.set(ActivityID("c4"), "Child4", ActivityStatusPending, 0)

	for _, maxW := range []int{80, 40, 30, 20, 15, 10, 5, 3} {
		got := dt.RenderWithSnapshots(snaps.snaps, 20, maxW)
		for line := range strings.SplitSeq(got, "\n") {
			w := VisibleWidth(line)
			if w > maxW {
				t.Errorf("maxWidth=%d: line visible width %d exceeds limit: %q",
					maxW, w, StripANSI(line))
			}
		}
	}
}
