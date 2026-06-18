package nom

import (
	"context"
	"testing"

	"github.com/larsartmann/go-output"
)

// TestDiagramExport_SubscriberProjection proves the subscriber's Store()
// correctly projects live progress state as output.GraphNode/Edge slices.
func TestDiagramExport_SubscriberProjection(t *testing.T) {
	t.Parallel()

	subscriber := NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = subscriber.OnEvent(ctx, &diagramTestEvent{
		eventType: EventWorkflowStarted,
		wID:       WorkflowID("build"),
		wName:     WorkflowName("CI Build"),
	})
	_ = subscriber.OnEvent(ctx, &diagramTestEvent{
		eventType: EventActivityStarted,
		aID:       ActivityID("compile"),
		aName:     ActivityName("Compile Sources"),
	})
	_ = subscriber.OnEvent(ctx, &diagramTestEvent{
		eventType:    EventActivityStarted,
		aID:          ActivityID("test"),
		aName:        ActivityName("Run Tests"),
		dependencies: []ActivityID{ActivityID("compile")},
	})
	_ = subscriber.OnEvent(ctx, &diagramTestEvent{
		eventType:    EventActivityStarted,
		aID:          ActivityID("deploy"),
		aName:        ActivityName("Deploy"),
		dependencies: []ActivityID{ActivityID("test")},
	})

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

// TestDiagramExport_StatusShapes verifies distinct NodeShape per status.
func TestDiagramExport_StatusShapes(t *testing.T) {
	t.Parallel()

	store := NewActivityStore()

	running := NewActivity("build", "Build")
	running.SetRunning()
	store.Upsert(running)

	completed := NewActivity("test", "Test")
	completed.SetCompleted()
	store.Upsert(completed)

	failed := NewActivity("lint", "Lint")
	failed.SetFailed(errTestFailure)
	store.Upsert(failed)

	nodes := store.Nodes()
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

	store := NewActivityStore()
	store.Upsert(NewActivity("a", "Alpha"))
	store.Upsert(NewActivity("b", "Beta"))
	store.Upsert(NewActivity("c", "Gamma"))

	store.AddEdge(
		output.NewBrandedID[output.GraphNodeIDBrand]("a"),
		output.NewBrandedID[output.GraphNodeIDBrand]("b"),
	)
	store.AddEdge(
		output.NewBrandedID[output.GraphNodeIDBrand]("b"),
		output.NewBrandedID[output.GraphNodeIDBrand]("c"),
	)

	edges := store.Edges()
	if len(edges) != 2 {
		t.Fatalf("Edges = %d, want 2", len(edges))
	}

	roots := store.Roots()
	if len(roots) != 1 || roots[0].Get() != "a" {
		t.Errorf("roots = %v, want [a]", roots)
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
