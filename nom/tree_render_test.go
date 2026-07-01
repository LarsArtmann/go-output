package nom

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDependencyTree_TreePrefix_RootNode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), nil)
	dt.AddActivity(ActivityID("child"), []ActivityID{"root"})

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
	dt.AddActivity(ActivityID("a"), nil)

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
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("b"), nil)
	dt.AddActivity(ActivityID("c"), nil)

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

	err := dt.AddActivity(ActivityID("child"), []ActivityID{"nonexistent"})
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
	dt.AddActivity(ActivityID("a"), nil)
	dt.AddActivity(ActivityID("a"), nil)

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
	dt.AddActivity(ActivityID("phase"), nil)
	dt.AddActivity(ActivityID("step1"), []ActivityID{"phase"})
	dt.AddActivity(ActivityID("step2"), []ActivityID{"phase", "step1"})

	snaps := newSnapshotBuilder()
	snaps.set(ActivityID("phase"), "Phase", ActivityStatusPending, 0)
	snaps.set(ActivityID("step1"), "Step1", ActivityStatusPending, 0)
	snaps.set(ActivityID("step2"), "Step2", ActivityStatusPending, 0)

	got := dt.RenderWithSnapshots(snaps.snaps, 10, 0)
	if got == "" {
		t.Fatal("Render should produce output")
	}

	if !strings.Contains(got, "←") {
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
	dt.AddActivity(ActivityID("build"), nil)
	dt.AddActivity(ActivityID("compile"), []ActivityID{"build"})

	snaps := newSnapshotBuilder()
	snaps.setPhase(ActivityID("build"), "Build", ActivityStatusRunning, 0)
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
	dt.AddActivity(ActivityID("build"), nil)
	dt.AddActivity(ActivityID("compile"), []ActivityID{"build"})
	dt.AddActivity(ActivityID("test"), []ActivityID{"build"})
	dt.AddActivity(ActivityID("lint"), []ActivityID{"build"})
	dt.AddActivity(ActivityID("deploy"), []ActivityID{"build"})

	snaps := newSnapshotBuilder()
	snaps.setPhase(ActivityID("build"), "Build Phase", ActivityStatusRunning, 5*time.Second)
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
	dt.AddActivity(ActivityID("a"), nil)

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
	dt.AddActivity(ActivityID("root"), nil)
	dt.AddActivity(ActivityID("c1"), []ActivityID{"root"})
	dt.AddActivity(ActivityID("c2"), []ActivityID{"c1"})
	dt.AddActivity(ActivityID("c3"), []ActivityID{"c2"})
	dt.AddActivity(ActivityID("c4"), []ActivityID{"c3"})

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

func TestRenderWithSnapshots_CollapseMarkerUnderHeightPressure(t *testing.T) {
	t.Parallel()

	tree := NewDependencyTree()
	tree.AddActivity(ActivityID("root"), nil)

	snaps := newSnapshotBuilder()

	snaps.set(ActivityID("root"), "Root", ActivityStatusRunning, 0)

	for i := range 6 {
		id := ActivityID(fmt.Sprintf("c%d", i))
		tree.AddActivity(id, []ActivityID{"root"})

		status := ActivityStatusCompleted
		if i >= 4 { // last two stay running
			status = ActivityStatusRunning
		}

		snaps.set(id, fmt.Sprintf("Child %d", i), status, 0)
	}

	if err := tree.Build(); err != nil {
		t.Fatalf("build: %v", err)
	}

	// maxHeight=4: root + 2 running + 1 collapse-marker line for the 4 completed.
	rendered := tree.RenderWithSnapshots(snaps.snaps, 4, 0)

	// The marker must be a dedicated tree line with the ellipsis + exact count.
	// Asserting the glyph alone is weak (could match elsewhere); assert the
	// full "⋯ 4 completed" phrase appears on its own line with a connector.
	var markerFound bool

	for line := range strings.SplitSeq(rendered, "\n") {
		if strings.Contains(line, "⋯") && strings.Contains(line, "4 completed") &&
			(strings.Contains(line, "├──") || strings.Contains(line, "└──")) {
			markerFound = true
			break
		}
	}

	if !markerFound {
		t.Errorf("expected a tree line like '├── ⋯ 4 completed'\ngot:\n%s", rendered)
	}

	// And the 4 completed children themselves must be elided (not rendered).
	for i := range 4 {
		if strings.Contains(rendered, fmt.Sprintf("Child %d", i)) {
			t.Errorf("completed child %d should be elided under height pressure\ngot:\n%s",
				i, rendered)
		}
	}
}
