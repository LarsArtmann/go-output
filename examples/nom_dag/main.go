package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/larsartmann/go-output/nom"
)

func main() {
	ctx := context.Background()

	subscriber := nom.NewNOMSubscriber(
		nom.WithRenderMode(nom.RenderModeLayered),
		nom.WithShowCategory(),
		nom.WithTheme(nom.ThemeDracula),
	)

	// Every event goes through one checked send helper — ignoring OnEvent
	// errors hides subscriber rejections in real pipelines.
	send := func(evt nom.Event) {
		if err := subscriber.OnEvent(ctx, evt); err != nil {
			fmt.Fprintf(os.Stderr, "event error: %v\n", err)
		}
	}

	send(nom.WorkflowStarted{
		ID:   nom.NewWorkflowID("dag-demo"),
		Name: nom.NewWorkflowName("DAG Build Pipeline"),
	})

	// Layer 0: root setup phase.
	send(nom.ActivityRegistered{
		ID:       nom.NewActivityID("setup"),
		Name:     nom.NewActivityName("Setup"),
		Kind:     nom.ActivityKindPhase,
		Category: nom.ActivityCategory("infra"),
	})

	// Layer 1: compile + lint depend on setup.
	send(nom.ActivityRegistered{
		ID:       nom.NewActivityID("compile"),
		Name:     nom.NewActivityName("Compile"),
		Deps:     []nom.ActivityID{nom.NewActivityID("setup")},
		Category: nom.ActivityCategory("build"),
	})
	send(nom.ActivityRegistered{
		ID:       nom.NewActivityID("lint"),
		Name:     nom.NewActivityName("Lint"),
		Deps:     []nom.ActivityID{nom.NewActivityID("setup")},
		Category: nom.ActivityCategory("build"),
	})

	// Layer 2: test depends on compile, deploy depends on lint.
	send(nom.ActivityRegistered{
		ID:       nom.NewActivityID("test"),
		Name:     nom.NewActivityName("Test"),
		Deps:     []nom.ActivityID{nom.NewActivityID("compile")},
		Category: nom.ActivityCategory("test"),
	})
	send(nom.ActivityRegistered{
		ID:       nom.NewActivityID("deploy"),
		Name:     nom.NewActivityName("Deploy"),
		Deps:     []nom.ActivityID{nom.NewActivityID("lint")},
		Category: nom.ActivityCategory("deploy"),
	})

	// Start the workflow.
	startAndProgress(send)

	// Render the final tree.
	snaps := subscriber.SnapshotActivities()
	rendered := subscriber.DependencyTree().RenderWithSnapshots(snaps, 20, 80)
	fmt.Println(rendered)

	// Print DAG structural summary.
	summary := subscriber.DependencyTree().DAGSummaryWithSnapshots(snaps)
	fmt.Printf("\nDAG: %s\n", summary.String())
}

func startAndProgress(send func(evt nom.Event)) {
	send(nom.ActivityStarted{ID: nom.NewActivityID("setup"), Name: nom.NewActivityName("Setup")})
	time.Sleep(100 * time.Millisecond)
	send(nom.ActivityCompleted{
		ID:       nom.NewActivityID("setup"),
		Name:     nom.NewActivityName("Setup"),
		Duration: 100 * time.Millisecond,
	})

	// Compile + lint run in parallel.
	send(nom.ActivityStarted{ID: nom.NewActivityID("compile"), Name: nom.NewActivityName("Compile")})
	send(nom.ActivityStarted{ID: nom.NewActivityID("lint"), Name: nom.NewActivityName("Lint")})
	time.Sleep(100 * time.Millisecond)
	send(nom.ActivityCompleted{
		ID:       nom.NewActivityID("compile"),
		Name:     nom.NewActivityName("Compile"),
		Duration: 200 * time.Millisecond,
	})
	send(nom.ActivityCompleted{
		ID:       nom.NewActivityID("lint"),
		Name:     nom.NewActivityName("Lint"),
		Duration: 150 * time.Millisecond,
	})

	// Test runs after compile.
	send(nom.ActivityStarted{ID: nom.NewActivityID("test"), Name: nom.NewActivityName("Test")})
	time.Sleep(50 * time.Millisecond)
	send(nom.ActivityCompleted{
		ID:       nom.NewActivityID("test"),
		Name:     nom.NewActivityName("Test"),
		Duration: 100 * time.Millisecond,
	})

	// Deploy runs after lint.
	send(nom.ActivityStarted{ID: nom.NewActivityID("deploy"), Name: nom.NewActivityName("Deploy")})
	time.Sleep(50 * time.Millisecond)
	send(nom.ActivityCompleted{
		ID:       nom.NewActivityID("deploy"),
		Name:     nom.NewActivityName("Deploy"),
		Duration: 80 * time.Millisecond,
	})

	send(nom.WorkflowCompleted{})
}
