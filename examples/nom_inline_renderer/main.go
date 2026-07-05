// Package main demonstrates the NOM InlineRenderer for live terminal output.
// Unlike nom_progress (which renders a static snapshot) or tui_progress (which
// uses the full Bubble Tea TUI), this example shows the simplest live inline
// rendering path: Start the renderer, fire events from workers, Stop when done.
package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/larsartmann/go-output/nom"
)

type step struct {
	id, name string
	deps     []string
	duration time.Duration
}

const stepCompile = "compile"

func main() {
	subscriber := nom.NewNOMSubscriber()
	renderer := nom.NewInlineRenderer(subscriber, os.Stdout, 15)
	renderer.SetAppName("Demo Build")

	ctx := context.Background()
	renderer.Start(ctx, 100*time.Millisecond)

	send := func(evt nom.Event) {
		if err := subscriber.OnEvent(ctx, evt); err != nil {
			log.Fatalf("event: %v", err)
		}
	}

	send(nom.WorkflowStarted{
		ID:   nom.NewWorkflowID("demo"),
		Name: nom.NewWorkflowName("Demo Build"),
	})

	steps := []step{
		{"fetch", "Fetch Dependencies", nil, 300 * time.Millisecond},
		{"compile", "Compile Sources", []string{"fetch"}, 500 * time.Millisecond},
		{"test", "Run Tests", []string{stepCompile}, 400 * time.Millisecond},
		{"lint", "Lint Code", []string{stepCompile}, 200 * time.Millisecond},
	}

	send(nom.ActivityStarted{
		ID:   nom.NewActivityID("fetch"),
		Name: nom.NewActivityName("Fetch Dependencies"),
	})

	runStep := func(s step, wg *sync.WaitGroup, send func(nom.Event)) {
		defer wg.Done()

		depIDs := make([]nom.ActivityID, 0, len(s.deps))
		for _, d := range s.deps {
			depIDs = append(depIDs, nom.NewActivityID(d))
		}

		send(nom.ActivityStarted{
			ID:   nom.NewActivityID(s.id),
			Name: nom.NewActivityName(s.name),
			Deps: depIDs,
		})

		time.Sleep(s.duration)

		send(nom.ActivityCompleted{
			ID:       nom.NewActivityID(s.id),
			Name:     nom.NewActivityName(s.name),
			Duration: s.duration,
		})
	}

	var wg sync.WaitGroup

	for _, s := range steps[1:] {
		wg.Add(1)

		go runStep(s, &wg, send)
	}

	time.Sleep(300 * time.Millisecond)

	send(nom.ActivityCompleted{
		ID:       nom.NewActivityID("fetch"),
		Name:     nom.NewActivityName("Fetch Dependencies"),
		Duration: 300 * time.Millisecond,
	})

	wg.Wait()

	renderer.Stop()
	renderer.Finish(nil)

	counts := subscriber.GetActivityCounts()

	println() //nolint:forbidigo // blank line after inline rendering

	log.Printf("Build complete: %d steps succeeded\n", counts.Completed)
}
