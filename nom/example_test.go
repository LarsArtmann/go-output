package nom_test

import (
	"context"
	"fmt"
	"time"

	"github.com/larsartmann/go-output/nom"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewNOMStyleSubscriber() {
	ctx := context.Background()
	sub := nom.NewNOMStyleSubscriber()

	_ = sub.OnEvent(ctx, nom.WorkflowStarted{
		ID:   nom.NewWorkflowID("build"),
		Name: nom.NewWorkflowName("Build"),
	})
	_ = sub.OnEvent(ctx, nom.ActivityStarted{
		ID:   nom.NewActivityID("compile"),
		Name: nom.NewActivityName("Compile"),
	})
	_ = sub.OnEvent(ctx, nom.ActivityCompleted{
		ID:       nom.NewActivityID("compile"),
		Name:     nom.NewActivityName("Compile"),
		Duration: 2 * time.Second,
	})

	snaps := sub.SnapshotActivities()
	fmt.Println(sub.GetDependencyTree().RenderWithSnapshots(snaps, 20, 0))

	counts := sub.GetActivityCounts()
	fmt.Printf("Completed: %d\n", counts.Completed)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleFormatDuration() {
	fmt.Println(nom.FormatDuration(500 * time.Millisecond))
	fmt.Println(nom.FormatDuration(2 * time.Second))
	fmt.Println(nom.FormatDuration(90 * time.Second))
	fmt.Println(nom.FormatDuration(2 * time.Hour))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleActivityStatus_String() {
	fmt.Println(nom.ActivityStatusPending)
	fmt.Println(nom.ActivityStatusRunning)
	fmt.Println(nom.ActivityStatusCompleted)
	fmt.Println(nom.ActivityStatusFailed)
}
