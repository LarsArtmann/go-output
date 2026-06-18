package nom

import (
	"context"
	"testing"
	"time"
)

// TestDiagramExport_StatusShapes verifies distinct NodeShape per status.
func TestDiagramExport_StatusShapes(t *testing.T) {
	t.Parallel()

	subscriber := NewNOMStyleSubscriber()
	ctx := context.Background()

	diagramFireWorkflow(t, subscriber, ctx, "wf", "Workflow")
	diagramFireActivity(t, subscriber, ctx, "build", "Build")
	diagramFireActivity(t, subscriber, ctx, "test", "Test")
	diagramFireActivity(t, subscriber, ctx, "lint", "Lint")

	// All start as Running; transition to different terminal states.
	// testSetStatus mutates the shared *Activity pointer, so Store() sees it.
	dt := subscriber.GetDependencyTree()
	now := time.Now()
	testSetStatus(dt, ActivityID("test"), ActivityStatusCompleted, now)
	testSetStatus(dt, ActivityID("lint"), ActivityStatusFailed, now)

	nodes := subscriber.Store().Nodes()

	shapes := make(map[string]string)
	for _, n := range nodes {
		shapes[n.ID.Get()] = string(n.Shape)
	}

	if shapes["build"] == shapes["test"] {
		t.Error("running and completed should have different shapes")
	}

	if shapes["test"] == shapes["lint"] {
		t.Error("completed and failed should have different shapes")
	}
}

// TestDiagramExport_EdgeStructure verifies dependency chain projection.
func TestDiagramExport_EdgeStructure(t *testing.T) {
	t.Parallel()

	subscriber := NewNOMStyleSubscriber()
	ctx := context.Background()

	diagramFireWorkflow(t, subscriber, ctx, "wf", "Workflow")
	diagramFireActivity(t, subscriber, ctx, "a", "Alpha")
	diagramFireActivity(t, subscriber, ctx, "b", "Beta", "a")
	diagramFireActivity(t, subscriber, ctx, "c", "Gamma", "b")

	reader := subscriber.Store()

	edges := reader.Edges()
	if len(edges) != 2 {
		t.Fatalf("Edges = %d, want 2 (a→b, b→c)", len(edges))
	}

	// Verify "a" is the root (no incoming edges).
	for _, e := range edges {
		if e.To.Get() == "a" {
			t.Error("node 'a' should be a root with no incoming edges")
		}
	}
}

// TestDiagramExport_SubscriberProjection proves the subscriber's Store()
// correctly projects live progress state as output.GraphNode/Edge slices.
func TestDiagramExport_SubscriberProjection(t *testing.T) {
	t.Parallel()

	subscriber := NewNOMStyleSubscriber()
	ctx := context.Background()

	diagramFireWorkflow(t, subscriber, ctx, "build", "CI Build")
	diagramFireActivity(t, subscriber, ctx, "compile", "Compile Sources")
	diagramFireActivity(t, subscriber, ctx, "test", "Run Tests", "compile")
	diagramFireActivity(t, subscriber, ctx, "deploy", "Deploy", "test")

	reader := subscriber.Store()

	nodes := reader.Nodes()
	if len(nodes) != 3 {
		t.Fatalf("Nodes() = %d nodes, want 3", len(nodes))
	}

	for _, n := range nodes {
		if n.ID.Get() == "" {
			t.Error("empty node ID")
		}

		if n.Shape == "" {
			t.Errorf("node %q has empty Shape", n.ID.Get())
		}

		if n.Style.Fill == "" {
			t.Errorf("node %q has empty Style.Fill", n.ID.Get())
		}
	}

	edges := reader.Edges()
	if len(edges) != 2 {
		t.Fatalf("Edges() = %d, want 2 (compile→test, test→deploy)", len(edges))
	}

	counts := subscriber.GetActivityCounts()
	if counts.Running != 3 {
		t.Errorf("running = %d, want 3", counts.Running)
	}
}

type diagramTestEvent struct {
	eventType    string
	wID          WorkflowID
	wName        WorkflowName
	aID          ActivityID
	aName        ActivityName
	dependencies []ActivityID
}

func (e *diagramTestEvent) GetEventType() string          { return e.eventType }
func (e *diagramTestEvent) GetWorkflowID() WorkflowID     { return e.wID }
func (e *diagramTestEvent) GetWorkflowName() WorkflowName { return e.wName }
func (e *diagramTestEvent) GetActivityID() ActivityID     { return e.aID }
func (e *diagramTestEvent) GetActivityName() ActivityName { return e.aName }
func (e *diagramTestEvent) GetDependencies() []ActivityID { return e.dependencies }

// diagramFireWorkflow fires a workflow.started event using diagramTestEvent.
func diagramFireWorkflow(t *testing.T, ns *NOMStyleSubscriber, ctx context.Context, id, name string) {
	t.Helper()

	_ = ns.OnEvent(ctx, &diagramTestEvent{
		eventType: EventWorkflowStarted,
		wID:       WorkflowID(id),
		wName:     WorkflowName(name),
	})
}

// diagramFireActivity fires an activity.started event with optional dependencies.
func diagramFireActivity(t *testing.T, ns *NOMStyleSubscriber, ctx context.Context, id, name string, deps ...string) {
	t.Helper()

	dependencies := make([]ActivityID, len(deps))
	for i, dep := range deps {
		dependencies[i] = ActivityID(dep)
	}

	_ = ns.OnEvent(ctx, &diagramTestEvent{
		eventType:    EventActivityStarted,
		aID:          ActivityID(id),
		aName:        ActivityName(name),
		dependencies: dependencies,
	})
}
