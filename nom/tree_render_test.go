package nom

import (
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
