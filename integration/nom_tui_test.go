package integration

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/go-output/nom"
	"github.com/larsartmann/go-output/tui"
)

func TestNOMSubscriber_Integration(t *testing.T) {
	t.Parallel()

	t.Run("full workflow lifecycle through events", func(t *testing.T) {
		t.Parallel()

		subscriber := nom.NewNOMSubscriber(nom.WithCachePath(filepath.Join(t.TempDir(), "nom-timing.csv")))
		ctx := context.Background()

		if err := subscriber.OnEvent(ctx, nom.WorkflowStarted{
			ID:   nom.NewWorkflowID("ci-1"),
			Name: nom.NewWorkflowName("CI Pipeline"),
		}); err != nil {
			t.Fatalf("workflow.started: %v", err)
		}

		if !subscriber.IsWorkflowRunning() {
			t.Error("workflow should be running")
		}

		if subscriber.GetWorkflowName() != "CI Pipeline" {
			t.Errorf("workflow name = %q, want %q", subscriber.GetWorkflowName(), "CI Pipeline")
		}

		if err := subscriber.OnEvent(ctx, nom.ActivityStarted{
			ID:   nom.NewActivityID("build"),
			Name: nom.NewActivityName("Build Project"),
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

		tree := subscriber.DependencyTree()
		if tree == nil {
			t.Fatal("dependency tree should exist")
		}

		if err := subscriber.OnEvent(ctx, nom.ActivityCompleted{
			ID:       nom.NewActivityID("build"),
			Name:     nom.NewActivityName("Build Project"),
			Duration: 5 * time.Second,
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

		rendered := tree.RenderWithSnapshots(nil, 10, 0)
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

		tree.AddActivity(nom.NewActivityID("root"), nil)
		tree.AddActivity(nom.NewActivityID("child1"), []nom.ActivityID{"root"})
		tree.AddActivity(nom.NewActivityID("child2"), []nom.ActivityID{"root"})
		tree.AddActivity(nom.NewActivityID("grandchild"), []nom.ActivityID{"child1"})

		snaps := map[nom.ActivityID]nom.ActivitySnapshot{
			nom.NewActivityID("root"): {
				Label:  "Root Task",
				Status: nom.ActivityStatusCompleted,
				Symbol: nom.ActivityStatusCompleted.GetSymbol(),
				Color:  nom.ActivityStatusCompleted.GetColor(),
			},
			nom.NewActivityID("child1"): {
				Label:  "Child 1",
				Status: nom.ActivityStatusRunning,
				Symbol: nom.ActivityStatusRunning.GetSymbol(),
				Color:  nom.ActivityStatusRunning.GetColor(),
			},
			nom.NewActivityID("child2"): {
				Label:  "Child 2",
				Status: nom.ActivityStatusPending,
				Symbol: nom.ActivityStatusPending.GetSymbol(),
				Color:  nom.ActivityStatusPending.GetColor(),
			},
			nom.NewActivityID("grandchild"): {
				Label:  "Grandchild",
				Status: nom.ActivityStatusFailed,
				Symbol: nom.ActivityStatusFailed.GetSymbol(),
				Color:  nom.ActivityStatusFailed.GetColor(),
			},
		}

		rendered := tree.RenderWithSnapshots(snaps, 10, 0)
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

	// Every Report* call lazily starts a REAL tea.Program on os.Stdout
	// (see tui/lifecycle.go ensureStarted), and the reporter exposes no state
	// getters outside the tui package. Cross-package smoke coverage of the
	// Report* surface therefore lives in tui/reporter_test.go, which uses the
	// package-internal newTestReporter() to suppress program startup.
	//
	// The one observable cross-package surface is the embedded nom subscriber:
	// it must be live, so events fired through it are visible to whatever
	// renders NOM mode inside the TUI.
	t.Run("reporter subscriber is live for nom events", func(t *testing.T) {
		t.Parallel()

		reporter := tui.NewBubbleTeaProgressReporter()
		defer reporter.Stop()

		sub := reporter.Subscriber()
		if sub == nil {
			t.Fatal("reporter.Subscriber() returned nil")
		}

		ctx := context.Background()
		fireWorkflowStarted(t, sub, ctx, "w1", "Pipeline")
		startActivity(t, sub, ctx, "build", "Build")
		completeActivity(t, sub, ctx, "build", "Build", 2*time.Second)

		snaps := sub.SnapshotActivities()

		buildSnap, ok := snaps[nom.NewActivityID("build")]
		if !ok {
			t.Fatal("build activity missing from reporter subscriber snapshots")
		}

		if buildSnap.Status != nom.ActivityStatusCompleted {
			t.Errorf("build status = %v, want Completed", buildSnap.Status)
		}
	})
}

func TestNOMSubscriber_RenderNodeVisibleNodes_Integration(t *testing.T) {
	t.Parallel()

	subscriber := nom.NewNOMSubscriber(nom.WithCachePath(filepath.Join(t.TempDir(), "nom-timing.csv")))
	ctx := context.Background()

	fireWorkflowStarted(t, subscriber, ctx, "w1", "Pipeline")

	startActivity(t, subscriber, ctx, "build", "Build")
	startActivity(t, subscriber, ctx, "test", "Test")
	completeActivity(t, subscriber, ctx, "build", "Build", 3*time.Second)

	tree := subscriber.DependencyTree()
	if tree == nil {
		t.Fatal("dependency tree should exist")
	}

	snaps := subscriber.SnapshotActivities()

	visible := tree.VisibleNodesWithSnapshots(snaps, 10)
	if len(visible) == 0 {
		t.Fatal("VisibleNodes should return nodes")
	}

	for _, node := range visible {
		rendered := tree.RenderNode(node, snaps)
		if rendered == "" {
			t.Errorf("RenderNode(%s) returned empty string", node.ID)
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

	buildSnap := snaps[nom.NewActivityID("build")]
	if buildSnap.Status != nom.ActivityStatusCompleted {
		t.Errorf("build status = %v, want Completed", buildSnap.Status)
	}

	testSnap := snaps[nom.NewActivityID("test")]
	if testSnap.Status != nom.ActivityStatusRunning {
		t.Errorf("test status = %v, want Running", testSnap.Status)
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

// mustUpdateActivityStatus is removed — ActivityNode no longer stores Activity fields.
// Use SnapshotActivities + RenderWithSnapshots, or send events through the subscriber.

// startActivity sends an activity.started event, failing the test on error.
func startActivity(t *testing.T, sub *nom.NOMSubscriber, ctx context.Context, id, name string) {
	t.Helper()

	if err := sub.OnEvent(ctx, nom.ActivityStarted{
		ID:   nom.NewActivityID(id),
		Name: nom.NewActivityName(name),
	}); err != nil {
		t.Fatalf("activity.started(%s): %v", id, err)
	}
}

// registerActivity sends an activity.registered event, failing the test on error.
func registerActivity(t *testing.T, sub *nom.NOMSubscriber, ctx context.Context, id, name string, deps ...string) {
	t.Helper()

	depIDs := make([]nom.ActivityID, 0, len(deps))
	for _, dep := range deps {
		depIDs = append(depIDs, nom.NewActivityID(dep))
	}

	if err := sub.OnEvent(ctx, nom.ActivityRegistered{
		ID:   nom.NewActivityID(id),
		Name: nom.NewActivityName(name),
		Deps: depIDs,
	}); err != nil {
		t.Fatalf("activity.registered(%s): %v", id, err)
	}
}

// fireWorkflowStarted sends a workflow.started event, failing the test on error.
func fireWorkflowStarted(t *testing.T, sub *nom.NOMSubscriber, ctx context.Context, id, name string) {
	t.Helper()

	if err := sub.OnEvent(ctx, nom.WorkflowStarted{
		ID:   nom.NewWorkflowID(id),
		Name: nom.NewWorkflowName(name),
	}); err != nil {
		t.Fatalf("workflow.started(%s): %v", id, err)
	}
}

// completeActivity sends an activity.completed event, failing the test on error.
func completeActivity(
	t *testing.T,
	sub *nom.NOMSubscriber,
	ctx context.Context,
	id, name string,
	duration time.Duration,
) {
	t.Helper()

	if err := sub.OnEvent(ctx, nom.ActivityCompleted{
		ID:       nom.NewActivityID(id),
		Name:     nom.NewActivityName(name),
		Duration: duration,
	}); err != nil {
		t.Fatalf("activity.completed(%s): %v", id, err)
	}
}

func TestNOM_LayeredMode_Integration(t *testing.T) {
	t.Parallel()

	sub := nom.NewNOMSubscriber(nom.WithCachePath(filepath.Join(t.TempDir(), "nom-timing.csv")))
	ctx := context.Background()

	fireWorkflowStarted(t, sub, ctx, "wf1", "CI Pipeline")

	// Build a multi-layer DAG via events with explicit deps.
	registerActivity(t, sub, ctx, "setup", "Setup")
	registerActivity(t, sub, ctx, "compile", "Compile", "setup")
	registerActivity(t, sub, ctx, "lint", "Lint", "setup")
	registerActivity(t, sub, ctx, "test", "Test", "compile")

	startActivity(t, sub, ctx, "setup", "Setup")
	startActivity(t, sub, ctx, "compile", "Compile")
	startActivity(t, sub, ctx, "lint", "Lint")
	startActivity(t, sub, ctx, "test", "Test")

	tree := sub.DependencyTree()
	tree.SetRenderMode(nom.RenderModeLayered)

	completeActivity(t, sub, ctx, "setup", "Setup", 2*time.Second)

	snaps := sub.SnapshotActivities()
	rendered := tree.RenderWithSnapshots(snaps, 20, 100)

	if rendered == "" {
		t.Fatal("layered render should not be empty")
	}

	if !strings.Contains(rendered, "Setup") {
		t.Errorf("expected Setup in layered output:\n%s", rendered)
	}

	if !strings.Contains(rendered, "Compile") {
		t.Errorf("expected Compile in layered output:\n%s", rendered)
	}

	// DOT export via AllNodes.
	allNodes := tree.AllNodes()
	if len(allNodes) == 0 {
		t.Fatal("AllNodes should return nodes")
	}

	hasDeps := false

	for _, n := range allNodes {
		if len(n.Deps) > 0 {
			hasDeps = true

			break
		}
	}

	if !hasDeps {
		t.Error("expected at least one node with deps in the DAG")
	}

	// DAGSummary should produce a non-empty structural summary.
	summary := tree.DAGSummaryWithSnapshots(snaps)
	if summary.Nodes == 0 {
		t.Error("DAGSummary should report non-zero nodes")
	}

	summaryStr := summary.String()
	if !strings.Contains(summaryStr, "nodes") {
		t.Errorf("DAGSummary.String() should contain 'nodes': %q", summaryStr)
	}
}
