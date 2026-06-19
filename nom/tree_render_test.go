package nom

import (
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestDependencyTree_EnsureBuild(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)

	dt.EnsureBuild()

	roots := dt.GetRootNodes()
	if len(roots) != 1 {
		t.Errorf("expected 1 root after EnsureBuild, got %d", len(roots))
	}
}

func TestDependencyTree_TreePrefix_RootNode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), NewActivity("root", "Root"), nil)
	dt.AddActivity(ActivityID("child"), NewActivity("child", "Child"), []ActivityID{"root"})
	testSetStatus(dt, ActivityID("root"), ActivityStatusRunning, time.Now())
	testSetStatus(dt, ActivityID("child"), ActivityStatusPending, time.Time{})

	got := dt.RenderString(10)
	if got == "" {
		t.Error("Render() should produce output")
	}
}

func TestDependencyTree_Render_PausedStatus(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
	testSetStatus(dt, ActivityID("a"), ActivityStatusPaused, time.Now())

	got := dt.RenderString(10)
	if got == "" {
		t.Error("Render() should produce output for paused status")
	}
}

func TestDependencyTree_Render_FailedPriority(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), NewActivity("a", "A"), nil)
	dt.AddActivity(ActivityID("b"), NewActivity("b", "B"), nil)
	dt.AddActivity(ActivityID("c"), NewActivity("c", "C"), nil)
	testSetStatus(dt, ActivityID("a"), ActivityStatusCompleted, time.Now())
	testSetStatus(dt, ActivityID("b"), ActivityStatusFailed, time.Now())
	testSetStatus(dt, ActivityID("c"), ActivityStatusPending, time.Time{})

	got := dt.RenderString(3)
	if got == "" {
		t.Error("Render() should produce output")
	}
}

func TestDependencyTree_AddActivity_WithNonExistentDependency(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	err := dt.AddActivity(ActivityID("child"), NewActivity("child", "Child"), []ActivityID{"nonexistent"})
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
	dt.AddActivity(ActivityID("a"), NewActivity("a", "Original"), nil)
	dt.AddActivity(ActivityID("a"), NewActivity("a", "Updated"), nil)

	node := dt.GetNode(ActivityID("a"))
	testhelpers.AssertEqual(t, "ActivityName", "", node.Label.Get(), "Updated")
}

func TestDependencyTree_Render_SecondaryDependencies(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), NewActivity("phase", "Phase"), nil)
	dt.AddActivity(ActivityID("step1"), NewActivity("step1", "Step1"), []ActivityID{"phase"})
	dt.AddActivity(ActivityID("step2"), NewActivity("step2", "Step2"), []ActivityID{"phase", "step1"})

	got := dt.RenderString(10)
	if got == "" {
		t.Fatal("Render() should produce output")
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
	dt.AddActivity(ActivityID("phase:build"), NewActivity("phase:build", "Build"), nil)
	dt.AddActivity(ActivityID("compile"), NewActivity("compile", "Compile"), []ActivityID{"phase:build"})

	now := time.Now()
	testSetStatus(dt, ActivityID("phase:build"), ActivityStatusRunning, now)

	got := dt.RenderString(10)
	if got == "" {
		t.Fatal("Render() should produce output")
	}

	// Phase node should use phase symbol
	if !strings.Contains(got, string(SymbolPhase)) {
		t.Errorf("render should contain phase symbol %q, got:\n%s", SymbolPhase, got)
	}
}

func TestDependencyTree_Render_PriorityOrdering(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase:build"), NewActivity("phase:build", "Build Phase"), nil)
	dt.AddActivity(ActivityID("compile"), NewActivity("compile", "Compile"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("test"), NewActivity("test", "Run Tests"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("lint"), NewActivity("lint", "Lint Code"), []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("deploy"), NewActivity("deploy", "Deploy"), []ActivityID{"phase:build"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	setStatusWithElapsed(
		dt,
		ActivityID("phase:build"),
		ActivityStatusRunning,
		SymbolRunning,
		Colors.Running,
		now,
		5*time.Second,
	)
	setStatusWithElapsed(
		dt,
		ActivityID("compile"),
		ActivityStatusCompleted,
		SymbolCompleted,
		Colors.Completed,
		now,
		2*time.Second,
	)
	setStatusWithElapsed(
		dt,
		ActivityID("test"),
		ActivityStatusRunning,
		SymbolRunning,
		Colors.Running,
		now,
		3*time.Second,
	)
	setStatusWithElapsed(dt, ActivityID("lint"), ActivityStatusPending, SymbolPaused, Colors.Paused, time.Time{}, 0)
	setStatusWithElapsed(
		dt,
		ActivityID("deploy"),
		ActivityStatusFailed,
		SymbolFailed,
		Colors.Failed,
		now,
		1*time.Second,
	)

	got := dt.RenderString(10)

	// Children should appear in priority order: failed, running, pending, completed.
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

	// With limited height, failed and running must survive over completed.
	limited := dt.RenderString(3)
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
	dt.AddActivity(ActivityID("a"), NewActivity("a", "This is a very long activity name that will not fit"), nil)

	now := time.Now()
	testSetStatus(dt, ActivityID("a"), ActivityStatusRunning, now)

	wide := dt.RenderWithWidth(10, 80)
	if strings.Contains(wide, "…") {
		t.Errorf("wide render should not truncate, got:\n%s", wide)
	}

	narrow := dt.RenderWithWidth(10, 20)
	if !strings.Contains(narrow, "…") {
		t.Errorf("narrow render should truncate with ellipsis, got:\n%s", narrow)
	}
}

// TestDependencyTree_RenderWithWidth_DeepNestingFitsMaxWidth is a regression
// test: when the tree-drawing prefix (indentation + ├──/└──) alone exceeds
// maxWidth, the renderer used to emit the full prefix and overflow the terminal.
// Now every rendered line must fit within maxWidth.
func TestDependencyTree_RenderWithWidth_DeepNestingFitsMaxWidth(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), NewActivity("root", "Root"), nil)
	dt.AddActivity(ActivityID("c1"), NewActivity("c1", "Child1"), []ActivityID{"root"})
	dt.AddActivity(ActivityID("c2"), NewActivity("c2", "Child2"), []ActivityID{"c1"})
	dt.AddActivity(ActivityID("c3"), NewActivity("c3", "Child3"), []ActivityID{"c2"})
	dt.AddActivity(ActivityID("c4"), NewActivity("c4", "Child4"), []ActivityID{"c3"})

	for _, maxW := range []int{80, 40, 30, 20, 15, 10, 5, 3} {
		got := dt.RenderWithWidth(20, maxW)
		for line := range strings.SplitSeq(got, "\n") {
			w := VisibleWidth(line)
			if w > maxW {
				t.Errorf("maxWidth=%d: line visible width %d exceeds limit: %q",
					maxW, w, StripANSI(line))
			}
		}
	}
}
