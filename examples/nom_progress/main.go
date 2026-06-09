package main

import (
	"context"
	"fmt"
	"time"

	"github.com/larsartmann/go-output/nom"
)

type workflowEvent struct {
	eventType string
	wID       nom.WorkflowID
	wName     nom.WorkflowName
	aID       nom.ActivityID
	aName     nom.ActivityName
	duration  time.Duration
	err       error
}

func (e *workflowEvent) GetEventType() string { return e.eventType }

func (e *workflowEvent) GetWorkflowID() nom.WorkflowID { return e.wID }

func (e *workflowEvent) GetWorkflowName() nom.WorkflowName { return e.wName }

func (e *workflowEvent) GetActivityID() nom.ActivityID { return e.aID }

func (e *workflowEvent) GetActivityName() nom.ActivityName { return e.aName }

func (e *workflowEvent) GetDuration() time.Duration { return e.duration }

func (e *workflowEvent) GetError() error { return e.err }

func main() {
	subscriber := nom.NewNOMStyleSubscriber()
	ctx := context.Background()

	_ = subscriber.OnEvent(ctx, &workflowEvent{
		eventType: "workflow.started",
		wID:       nom.NewWorkflowID("build-42"),
		wName:     nom.NewWorkflowName("CI Pipeline"),
	})

	activities := []struct {
		id   string
		name string
	}{
		{"fetch", "Fetch Dependencies"},
		{"compile", "Compile Sources"},
		{"test", "Run Tests"},
		{"package", "Package Binary"},
	}

	for _, a := range activities {
		_ = subscriber.OnEvent(ctx, &workflowEvent{
			eventType: "activity.started",
			aID:       nom.NewActivityID(a.id),
			aName:     nom.NewActivityName(a.name),
		})

		time.Sleep(100 * time.Millisecond)

		_ = subscriber.OnEvent(ctx, &workflowEvent{
			eventType: "activity.completed",
			aID:       nom.NewActivityID(a.id),
			aName:     nom.NewActivityName(a.name),
			duration:  100 * time.Millisecond,
		})
	}

	subscriber.UpdateRunningActivityElapsed()
	subscriber.SyncActivityTimingToTree()

	tree := subscriber.GetDependencyTree()

	fmt.Println("=== NOM Dependency Tree ===")

	fmt.Println(tree.Render(20))

	running, completed, failed, pending := subscriber.GetActivityCounts()
	fmt.Printf("\nRunning: %d, Completed: %d, Failed: %d, Pending: %d\n",
		running, completed, failed, pending)

	summary := nom.GetActivitySummaryString(running, 0, 0, running+completed+failed+pending)
	fmt.Printf("Summary: %s\n", summary)

	_ = subscriber.OnEvent(ctx, &workflowEvent{
		eventType: "workflow.completed",
		wID:       nom.NewWorkflowID("build-42"),
	})

	fmt.Println("\nWorkflow completed!")
}
