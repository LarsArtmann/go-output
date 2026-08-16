package main

import (
	"context"
	"fmt"
	"time"

	"github.com/larsartmann/go-output/nom"
	"github.com/larsartmann/go-output/tui"
)

func main() {
	reporter := tui.NewBubbleTeaProgressReporter()

	// The first Report* call lazily starts a real Bubble Tea program that
	// takes over the terminal until Stop() is called.
	reporter.ReportMessage("Starting CI pipeline...")
	time.Sleep(50 * time.Millisecond)

	steps := []struct {
		message string
		total   uint
	}{
		{"Fetch Dependencies", 10},
		{"Compile Sources", 50},
		{"Run Tests", 20},
		{"Package Binary", 5},
	}

	for _, step := range steps {
		reporter.ReportMessage(step.message)

		for i := uint(1); i <= step.total; i++ {
			reporter.ReportStep(i, step.total, step.message)
			reporter.ReportProgress(float64(i) / float64(step.total) * 25.0)
		}
	}

	// NOM mode: the reporter embeds a live nom subscriber. Switching the
	// display mode and firing sealed events renders the NOM dependency tree
	// in the same TUI.
	reporter.SetDisplayMode(tui.DisplayModeNOM)

	sub := reporter.Subscriber()
	ctx := context.Background()

	send := func(evt nom.Event) {
		if err := sub.OnEvent(ctx, evt); err != nil {
			fmt.Printf("event error: %v\n", err)
		}
	}

	send(nom.WorkflowStarted{ID: nom.NewWorkflowID("demo"), Name: nom.NewWorkflowName("CI Pipeline")})
	send(nom.ActivityStarted{ID: nom.NewActivityID("build"), Name: nom.NewActivityName("Build Module")})
	send(nom.ActivityProgress{ID: nom.NewActivityID("build"), Name: nom.NewActivityName("Build Module"), Message: "compiling main.go"})
	time.Sleep(100 * time.Millisecond)
	send(nom.ActivityCompleted{ID: nom.NewActivityID("build"), Name: nom.NewActivityName("Build Module"), Duration: 100 * time.Millisecond})

	// Stop flushes the timing cache and shuts the TUI down gracefully.
	// Always call it before exiting — otherwise the program goroutine leaks.
	reporter.Stop()

	fmt.Println("TUI progress reporter demo complete.")

	fmt.Println("\n=== NOM reference ===")
	fmt.Printf("Symbols — Running: %s  Completed: %s  Failed: %s  Pending: %s\n",
		nom.SymbolRunning, nom.SymbolCompleted, nom.SymbolFailed, nom.SymbolPending)

	fmt.Println("\n=== Format Duration Examples ===")
	fmt.Printf("500ms: %s\n", nom.FormatDuration(500*time.Millisecond))
	fmt.Printf("1.5s:  %s\n", nom.FormatDuration(1500*time.Millisecond))
	fmt.Printf("2m30s: %s\n", nom.FormatDuration(1500*time.Second))
}
