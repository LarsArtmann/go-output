package main

import (
	"context"
	"errors"
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
		id     string
		name   string
		status nom.ActivityStatus
		delay  time.Duration
	}{
		{"fetch", "Fetch Dependencies", nom.ActivityStatusCompleted, 50 * time.Millisecond},
		{"compile", "Compile Sources", nom.ActivityStatusCompleted, 100 * time.Millisecond},
		{"test", "Run Tests", nom.ActivityStatusRunning, 200 * time.Millisecond},
		{"lint", "Lint Code", nom.ActivityStatusFailed, 0},
		{"package", "Package Binary", nom.ActivityStatusPending, 0},
	}

	for _, a := range activities {
		_ = subscriber.OnEvent(ctx, &workflowEvent{
			eventType: "activity.started",
			aID:       nom.NewActivityID(a.id),
			aName:     nom.NewActivityName(a.name),
		})

		time.Sleep(a.delay)

		if a.status == nom.ActivityStatusCompleted {
			_ = subscriber.OnEvent(ctx, &workflowEvent{
				eventType: "activity.completed",
				aID:       nom.NewActivityID(a.id),
				aName:     nom.NewActivityName(a.name),
				duration:  a.delay,
			})
		} else if a.status == nom.ActivityStatusFailed {
			_ = subscriber.OnEvent(ctx, &workflowEvent{
				eventType: "activity.failed",
				aID:       nom.NewActivityID(a.id),
				aName:     nom.NewActivityName(a.name),
				err:       errors.New("lint check failed"),
			})
		}
	}

	subscriber.UpdateRunningActivityElapsed()
	subscriber.SyncActivityTimingToTree()

	tree := subscriber.GetDependencyTree()

	fmt.Println("=== NOM Dependency Tree (priority-ordered) ===")
	fmt.Println("Failed/Running activities appear first when height is limited.")

	fmt.Println(tree.Render(3))
	fmt.Println()
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
