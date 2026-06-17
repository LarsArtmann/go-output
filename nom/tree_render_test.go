package nom

import (
	"strings"
	"testing"
	"time"
)

func TestDependencyTree_EnsureBuild(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), "A", nil)

	dt.EnsureBuild()

	roots := dt.GetRootNodes()
	if len(roots) != 1 {
		t.Errorf("expected 1 root after EnsureBuild, got %d", len(roots))
	}
}

func TestDependencyTree_TreePrefix_RootNode(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("root"), "Root", nil)
	dt.AddActivity(ActivityID("child"), "Child", []ActivityID{"root"})
	dt.UpdateActivityStatus(ActivityID("root"), ActivityStatusRunning, SymbolRunning, ColorRunning, time.Now(), 0)
	dt.UpdateActivityStatus(ActivityID("child"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)

	got := dt.Render(10)
	if got == "" {
		t.Error("Render() should produce output")
	}
}

func TestDependencyTree_Render_PausedStatus(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), "A", nil)
	dt.UpdateActivityStatus(ActivityID("a"), ActivityStatusPaused, SymbolPaused, ColorPaused, time.Now(), 0)

	got := dt.Render(10)
	if got == "" {
		t.Error("Render() should produce output for paused status")
	}
}

func TestDependencyTree_Render_FailedPriority(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), "A", nil)
	dt.AddActivity(ActivityID("b"), "B", nil)
	dt.AddActivity(ActivityID("c"), "C", nil)
	dt.UpdateActivityStatus(ActivityID("a"), ActivityStatusCompleted, SymbolCompleted, ColorCompleted, time.Now(), 0)
	dt.UpdateActivityStatus(ActivityID("b"), ActivityStatusFailed, SymbolFailed, ColorFailed, time.Now(), 0)
	dt.UpdateActivityStatus(ActivityID("c"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)

	got := dt.Render(3)
	if got == "" {
		t.Error("Render() should produce output")
	}
}

func TestDependencyTree_AddActivity_WithNonExistentDependency(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()

	err := dt.AddActivity(ActivityID("child"), "Child", []ActivityID{"nonexistent"})
	if err != nil {
		t.Fatalf("AddActivity() error: %v", err)
	}

	parent := dt.GetNode(ActivityID("nonexistent"))
	if parent == nil {
		t.Error("nonexistent dependency should be auto-created")
	}

	child := dt.GetNode(ActivityID("child"))
	if child.Parent == nil || child.Parent.ActivityID != ActivityID("nonexistent") {
		t.Error("child's parent should be auto-created dependency")
	}
}

func TestDependencyTree_AddActivity_UpdateExisting(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("a"), "Original", nil)
	dt.AddActivity(ActivityID("a"), "Updated", nil)

	node := dt.GetNode(ActivityID("a"))
	if node.ActivityName != "Updated" {
		t.Errorf("ActivityName = %q, want %q", node.ActivityName, "Updated")
	}
}

func TestDependencyTree_Render_SecondaryDependencies(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase"), "Phase", nil)
	dt.AddActivity(ActivityID("step1"), "Step1", []ActivityID{"phase"})
	dt.AddActivity(ActivityID("step2"), "Step2", []ActivityID{"phase", "step1"})

	got := dt.Render(10)
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
	dt.AddActivity(ActivityID("phase:build"), "Build", nil)
	dt.AddActivity(ActivityID("compile"), "Compile", []ActivityID{"phase:build"})

	now := time.Now()
	dt.UpdateActivityStatus(ActivityID("phase:build"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)

	got := dt.Render(10)
	if got == "" {
		t.Fatal("Render() should produce output")
	}

	// Phase node should use phase symbol
	if !strings.Contains(got, SymbolPhase) {
		t.Errorf("render should contain phase symbol %q, got:\n%s", SymbolPhase, got)
	}
}

func TestDependencyTree_Render_PriorityOrdering(t *testing.T) {
	t.Parallel()

	dt := NewDependencyTree()
	dt.AddActivity(ActivityID("phase:build"), "Build Phase", nil)
	dt.AddActivity(ActivityID("compile"), "Compile", []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("test"), "Run Tests", []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("lint"), "Lint Code", []ActivityID{"phase:build"})
	dt.AddActivity(ActivityID("deploy"), "Deploy", []ActivityID{"phase:build"})

	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	setStatusWithElapsed(
		dt,
		ActivityID("phase:build"),
		ActivityStatusRunning,
		SymbolRunning,
		ColorRunning,
		now,
		5*time.Second,
	)
	setStatusWithElapsed(
		dt,
		ActivityID("compile"),
		ActivityStatusCompleted,
		SymbolCompleted,
		ColorCompleted,
		now,
		2*time.Second,
	)
	setStatusWithElapsed(dt, ActivityID("test"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 3*time.Second)
	setStatusWithElapsed(dt, ActivityID("lint"), ActivityStatusPending, SymbolPaused, ColorPaused, time.Time{}, 0)
	setStatusWithElapsed(dt, ActivityID("deploy"), ActivityStatusFailed, SymbolFailed, ColorFailed, now, 1*time.Second)

	got := dt.Render(10)

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
	limited := dt.Render(3)
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
	dt.AddActivity(ActivityID("a"), "This is a very long activity name that will not fit", nil)

	now := time.Now()
	dt.UpdateActivityStatus(ActivityID("a"), ActivityStatusRunning, SymbolRunning, ColorRunning, now, 0)

	wide := dt.RenderWithWidth(10, 80)
	if strings.Contains(wide, "…") {
		t.Errorf("wide render should not truncate, got:\n%s", wide)
	}

	narrow := dt.RenderWithWidth(10, 20)
	if !strings.Contains(narrow, "…") {
		t.Errorf("narrow render should truncate with ellipsis, got:\n%s", narrow)
	}
}
