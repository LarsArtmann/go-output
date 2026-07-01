package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/larsartmann/go-output/nom"
)

var errLintFailed = errors.New("lint check failed")

func main() {
	subscriber := nom.NewNOMStyleSubscriber()
	ctx := context.Background()

	send := func(evt nom.Event) {
		if err := subscriber.OnEvent(ctx, evt); err != nil {
			log.Fatalf("event dispatch failed: %v", err)
		}
	}

	send(nom.WorkflowStarted{
		ID:   nom.NewWorkflowID("build-42"),
		Name: nom.NewWorkflowName("CI Pipeline"),
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
		send(nom.ActivityStarted{
			ID:   nom.NewActivityID(a.id),
			Name: nom.NewActivityName(a.name),
		})

		time.Sleep(a.delay)

		switch a.status {
		case nom.ActivityStatusCompleted:
			send(nom.ActivityCompleted{
				ID:       nom.NewActivityID(a.id),
				Name:     nom.NewActivityName(a.name),
				Duration: a.delay,
			})
		case nom.ActivityStatusFailed:
			send(nom.ActivityFailed{
				ID:   nom.NewActivityID(a.id),
				Name: nom.NewActivityName(a.name),
				Err:  errLintFailed,
			})
		case nom.ActivityStatusPending, nom.ActivityStatusRunning:
			// left as-is for the demo
		}
	}

	// --- Demonstrate the v0.21.0 BuildFlow integration events ---

	// Sub-step progress: a running activity reports what it's doing inside.
	send(nom.ActivityProgress{
		ID:      nom.NewActivityID("test"),
		Name:    nom.NewActivityName("Run Tests"),
		Message: "Running package [2/26]: nom",
	})

	// Retry visibility: the failed lint step is retried (transitions back to
	// running, renders ⟳1 with the reason).
	send(nom.ActivityRetrying{
		ID:      nom.NewActivityID("lint"),
		Name:    nom.NewActivityName("Lint Code"),
		Attempt: 1,
		Reason:  "timeout",
	})

	// External estimate injection: pending steps get predicted durations from an
	// external store (e.g. BuildFlow's SQLite timing cache).
	subscriber.SetEstimatedTime(nom.NewActivityID("package"), 8*time.Second)

	// Subscriber-owned remaining-time computation (powers the "~Xm left" summary).
	remaining := subscriber.EstimatedTotalRemaining()
	fmt.Printf("Estimated time remaining: %s\n", nom.FormatDuration(remaining))

	snaps := subscriber.SnapshotActivities()

	fmt.Println("=== NOM Dependency Tree (priority-ordered) ===")
	fmt.Println("Failed/Running activities appear first when height is limited.")

	fmt.Println(subscriber.GetDependencyTree().RenderWithSnapshots(snaps, 3, 0))
	fmt.Println()
	fmt.Println(subscriber.GetDependencyTree().RenderWithSnapshots(snaps, 20, 0))

	counts := subscriber.GetActivityCounts()
	fmt.Printf("\nRunning: %d, Completed: %d, Failed: %d, Pending: %d\n",
		counts.Running, counts.Completed, counts.Failed, counts.Pending)

	fmt.Printf("Summary: %s\n", counts.Summary())

	send(nom.WorkflowCompleted{
		ID: nom.NewWorkflowID("build-42"),
	})

	fmt.Println("\nWorkflow completed!")
}
