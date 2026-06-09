package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/larsartmann/go-output/nom"
	"github.com/larsartmann/go-output/tui"
)

func TestNOMSubscriber_Integration(t *testing.T) {
	t.Parallel()

	t.Run("full workflow lifecycle through events", func(t *testing.T) {
		t.Parallel()

		subscriber := nom.NewNOMStyleSubscriber()
		ctx := context.Background()

		if err := subscriber.OnEvent(ctx, &nomTestEvent{
			eventType: "workflow.started",
			wID:       nom.NewWorkflowID("ci-1"),
			wName:     nom.NewWorkflowName("CI Pipeline"),
		}); err != nil {
			t.Fatalf("workflow.started: %v", err)
		}

		if !subscriber.IsWorkflowRunning() {
			t.Error("workflow should be running")
		}

		if subscriber.GetWorkflowName() != "CI Pipeline" {
			t.Errorf("workflow name = %q, want %q", subscriber.GetWorkflowName(), "CI Pipeline")
		}

		if err := subscriber.OnEvent(ctx, &nomTestEvent{
			eventType: "activity.started",
			aID:       nom.NewActivityID("build"),
			aName:     nom.NewActivityName("Build Project"),
		}); err != nil {
			t.Fatalf("activity.started: %v", err)
		}

		activity := subscriber.GetActivity(nom.NewActivityID("build"))
		if activity == nil {
			t.Fatal("build activity should exist")
		}

		if !activity.IsRunning() {
			t.Error("build activity should be running")
		}

		tree := subscriber.GetDependencyTree()
		if tree == nil {
			t.Fatal("dependency tree should exist")
		}

		subscriber.UpdateRunningActivityElapsed()
		subscriber.SyncActivityTimingToTree()

		if err := subscriber.OnEvent(ctx, &nomTestEvent{
			eventType: "activity.completed",
			aID:       nom.NewActivityID("build"),
			aName:     nom.NewActivityName("Build Project"),
			duration:  5 * time.Second,
		}); err != nil {
			t.Fatalf("activity.completed: %v", err)
		}

		activity = subscriber.GetActivity(nom.NewActivityID("build"))
		if !activity.IsCompleted() {
			t.Error("build activity should be completed")
		}

		running, completed, failed, pending := subscriber.GetActivityCounts()
		if running != 0 {
			t.Errorf("running = %d, want 0", running)
		}

		if completed != 1 {
			t.Errorf("completed = %d, want 1", completed)
		}

		if failed != 0 {
			t.Errorf("failed = %d, want 0", failed)
		}

		if pending != 0 {
			t.Errorf("pending = %d, want 0", pending)
		}

		rendered := tree.Render(10)
		if rendered == "" {
			t.Error("tree render should not be empty")
		}

		subscriber.Reset()

		if subscriber.IsWorkflowRunning() {
			t.Error("workflow should not be running after reset")
		}
	})
}

func TestNOMDependencyTree_Integration(t *testing.T) {
	t.Parallel()

	t.Run("multi-level dependency tree renders correctly", func(t *testing.T) {
		t.Parallel()

		tree := nom.NewDependencyTree()

		tree.AddActivity(nom.NewActivityID("root"), "Root Task", nil)
		tree.AddActivity(nom.NewActivityID("child1"), "Child 1", []nom.ActivityID{"root"})
		tree.AddActivity(nom.NewActivityID("child2"), "Child 2", []nom.ActivityID{"root"})
		tree.AddActivity(nom.NewActivityID("grandchild"), "Grandchild", []nom.ActivityID{"child1"})

		tree.UpdateActivityStatus(
			nom.NewActivityID("root"),
			nom.ActivityStatusCompleted,
			nom.SymbolCompleted,
			nom.ColorCompleted,
			time.Now(),
			0,
		) //nolint:errcheck // test setup
		tree.UpdateActivityStatus(
			nom.NewActivityID("child1"),
			nom.ActivityStatusRunning,
			nom.SymbolRunning,
			nom.ColorRunning,
			time.Now(),
			10*time.Second,
		) //nolint:errcheck
		tree.UpdateActivityStatus(
			nom.NewActivityID("child2"),
			nom.ActivityStatusPending,
			nom.SymbolPaused,
			nom.ColorPaused,
			time.Time{},
			0,
		) //nolint:errcheck
		tree.UpdateActivityStatus(
			nom.NewActivityID("grandchild"),
			nom.ActivityStatusFailed,
			nom.SymbolFailed,
			nom.ColorFailed,
			time.Now(),
			0,
		) //nolint:errcheck

		rendered := tree.Render(10)
		if rendered == "" {
			t.Error("tree render should not be empty")
		}

		roots := tree.GetRootNodes()
		if len(roots) != 1 {
			t.Errorf("expected 1 root, got %d", len(roots))
		}

		if !roots[0].IsRoot {
			t.Error("root node should have IsRoot = true")
		}
	})
}

func TestTUIProgressReporter_Integration(t *testing.T) {
	t.Parallel()

	t.Run("reporter state transitions via progress", func(t *testing.T) {
		t.Parallel()

		reporter := tui.NewBubbleTeaProgressReporter()

		reporter.ReportProgress(25.0)
		reporter.ReportProgress(50.0)

		reporter.ReportStep(1, 3, "Build")
		reporter.ReportStep(2, 3, "Build")
		reporter.ReportStep(3, 3, "Build")

		reporter.ReportMessage("Compiling sources")

		reporter.ReportProgress(100.0)
	})

	t.Run("reporter error transitions to errored", func(t *testing.T) {
		t.Parallel()

		reporter := tui.NewBubbleTeaProgressReporter()
		reporter.ReportProgress(50.0)
		reporter.ReportError(errors.New("disk full"))
	})
}

func TestNOMTimingCache_Integration(t *testing.T) {
	t.Parallel()

	t.Run("cache records and retrieves averages", func(t *testing.T) {
		t.Parallel()

		cache := nom.NewTimingCache()

		cache.Record("build", 5*time.Second)
		cache.Record("build", 7*time.Second)
		cache.Record("build", 3*time.Second)

		avg := cache.GetAverage("build")
		if avg != 5*time.Second {
			t.Errorf("average = %v, want %v", avg, 5*time.Second)
		}
	})

	t.Run("cache GetAll returns all activities", func(t *testing.T) {
		t.Parallel()

		cache := nom.NewTimingCache()
		cache.Record("build", 10*time.Second)
		cache.Record("test", 5*time.Second)

		all := cache.GetAll()
		if len(all) != 2 {
			t.Errorf("GetAll() returned %d entries, want 2", len(all))
		}
	})
}

type nomTestEvent struct {
	eventType string
	wID       nom.WorkflowID
	wName     nom.WorkflowName
	aID       nom.ActivityID
	aName     nom.ActivityName
	duration  time.Duration
	err       error
}

func (e *nomTestEvent) GetEventType() string              { return e.eventType }
func (e *nomTestEvent) GetWorkflowID() nom.WorkflowID     { return e.wID }
func (e *nomTestEvent) GetWorkflowName() nom.WorkflowName { return e.wName }
func (e *nomTestEvent) GetActivityID() nom.ActivityID     { return e.aID }
func (e *nomTestEvent) GetActivityName() nom.ActivityName { return e.aName }
func (e *nomTestEvent) GetDuration() time.Duration        { return e.duration }
func (e *nomTestEvent) GetError() error                   { return e.err }
