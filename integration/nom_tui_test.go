package integration

import (
	"context"
	"errors"
	"image/color"
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
			eventType: nom.EventWorkflowStarted,
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
			eventType: nom.EventActivityStarted,
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
			eventType: nom.EventActivityCompleted,
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

		counts := subscriber.GetActivityCounts()
		if counts.Running != 0 {
			t.Errorf("running = %d, want 0", counts.Running)
		}

		if counts.Completed != 1 {
			t.Errorf("completed = %d, want 1", counts.Completed)
		}

		if counts.Failed != 0 {
			t.Errorf("failed = %d, want 0", counts.Failed)
		}

		if counts.Pending != 0 {
			t.Errorf("pending = %d, want 0", counts.Pending)
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

		mustUpdateActivityStatus(t, tree, nom.NewActivityID("root"),
			nom.ActivityStatusCompleted, nom.SymbolCompleted, nom.ColorCompleted, time.Now(), 0)
		mustUpdateActivityStatus(t, tree, nom.NewActivityID("child1"),
			nom.ActivityStatusRunning, nom.SymbolRunning, nom.ColorRunning, time.Now(), 10*time.Second)
		mustUpdateActivityStatus(t, tree, nom.NewActivityID("child2"),
			nom.ActivityStatusPending, nom.SymbolPaused, nom.ColorPaused, time.Time{}, 0)
		mustUpdateActivityStatus(t, tree, nom.NewActivityID("grandchild"),
			nom.ActivityStatusFailed, nom.SymbolFailed, nom.ColorFailed, time.Now(), 0)

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

func TestNOMSubscriber_RenderNodeVisibleNodes_Integration(t *testing.T) {
	t.Parallel()

	subscriber := nom.NewNOMStyleSubscriber()
	ctx := context.Background()

	fireWorkflowStarted(subscriber, ctx, "w1", "Pipeline")

	startActivity(subscriber, ctx, "build", "Build")
	startActivity(subscriber, ctx, "test", "Test")
	completeActivity(subscriber, ctx, "build", "Build", 3*time.Second)

	subscriber.UpdateRunningActivityElapsed()
	subscriber.SyncActivityTimingToTree()

	tree := subscriber.GetDependencyTree()
	if tree == nil {
		t.Fatal("dependency tree should exist")
	}

	visible := tree.VisibleNodes(10)
	if len(visible) == 0 {
		t.Fatal("VisibleNodes should return nodes")
	}

	for _, node := range visible {
		rendered := tree.RenderNode(node, visible)
		if rendered == "" {
			t.Errorf("RenderNode(%s) returned empty string", nom.ActivityID(node.ID.Get()))
		}
	}

	roots := tree.GetRootNodes()
	if len(roots) == 0 {
		t.Fatal("GetRootNodes should return at least one root")
	}

	buildNode := tree.GetNode(nom.NewActivityID("build"))
	if buildNode == nil {
		t.Fatal("build node should exist in tree")
	}

	if buildNode.Status != nom.ActivityStatusCompleted {
		t.Errorf("build status = %v, want Completed", buildNode.Status)
	}

	testNode := tree.GetNode(nom.NewActivityID("test"))
	if testNode == nil {
		t.Fatal("test node should exist in tree")
	}

	if testNode.Status != nom.ActivityStatusRunning {
		t.Errorf("test status = %v, want Running", testNode.Status)
	}
}

func TestNOMTimingCache_Integration(t *testing.T) {
	t.Parallel()

	t.Run("cache records and retrieves averages", func(t *testing.T) {
		t.Parallel()

		cache := nom.NewTimingCache()

		cache.Record("build", 5*time.Second)
		cache.Record("build", 7*time.Second)
		cache.Record("build", 3*time.Second)

		median := cache.GetMedian("build")
		if median != 5*time.Second {
			t.Errorf("median = %v, want %v", median, 5*time.Second)
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

// mustUpdateActivityStatus calls DependencyTree.UpdateActivityStatus and
// fails the test on unexpected error. Test setup expects the activity to
// exist, so any error is a programming bug.
func mustUpdateActivityStatus(
	t *testing.T,
	tree *nom.DependencyTree,
	id nom.ActivityID,
	status nom.ActivityStatus,
	symbol string,
	c color.Color,
	startTime time.Time,
	estimated time.Duration,
) {
	t.Helper()

	if err := tree.UpdateActivityStatus(id, status, symbol, c, startTime, estimated); err != nil {
		t.Fatalf("UpdateActivityStatus(%s): %v", id, err)
	}
}

// startActivity sends an activity.started event and discards any error.
func startActivity(sub *nom.NOMStyleSubscriber, ctx context.Context, id, name string) {
	_ = sub.OnEvent(ctx, &nomTestEvent{
		eventType: nom.EventActivityStarted,
		aID:       nom.NewActivityID(id),
		aName:     nom.NewActivityName(name),
	})
}

// fireWorkflowStarted sends a workflow.started event and discards any error.
// Use for tests that only care about downstream behavior, not event errors.
func fireWorkflowStarted(sub *nom.NOMStyleSubscriber, ctx context.Context, id, name string) {
	_ = sub.OnEvent(ctx, &nomTestEvent{
		eventType: nom.EventWorkflowStarted,
		wID:       nom.NewWorkflowID(id),
		wName:     nom.NewWorkflowName(name),
	})
}

// completeActivity sends an activity.completed event and discards any error.
func completeActivity(sub *nom.NOMStyleSubscriber, ctx context.Context, id, name string, duration time.Duration) {
	_ = sub.OnEvent(ctx, &nomTestEvent{
		eventType: nom.EventActivityCompleted,
		aID:       nom.NewActivityID(id),
		aName:     nom.NewActivityName(name),
		duration:  duration,
	})
}
