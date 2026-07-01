package main

import (
	"fmt"
	"time"

	"github.com/larsartmann/go-output/nom"
	"github.com/larsartmann/go-output/tui"
)

func main() {
	reporter := tui.NewBubbleTeaProgressReporter()

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

	reporter.ReportProgress(100.0)
	reporter.ReportMessage("CI pipeline complete!")

	fmt.Println("TUI progress reporter demo complete.")
	fmt.Println("In a real application, this would show a rich terminal UI.")

	fmt.Println("\n=== NOM Activity Display ===")

	subscriber := nom.NewNOMStyleSubscriber()
	fmt.Printf("Subscriber enabled: %v\n", subscriber.IsEnabled())
	fmt.Printf("Timing cache path: %s\n", subscriber.GetTimingCache().GetFilePath())

	fmt.Println("\n=== NOM Symbols ===")
	fmt.Printf("Running: %s  Completed: %s  Failed: %s  Pending: %s\n",
		nom.SymbolRunning, nom.SymbolCompleted, nom.SymbolFailed, nom.SymbolPending)
	fmt.Printf("Download: %s  Upload: %s  Average: %s\n",
		nom.SymbolDownload, nom.SymbolUpload, nom.SymbolAverage)

	fmt.Println("\n=== Format Duration Examples ===")
	fmt.Printf("500ms: %s\n", nom.FormatDuration(500*time.Millisecond))
	fmt.Printf("1.5s:  %s\n", nom.FormatDuration(1500*time.Millisecond))
	fmt.Printf("2m30s: %s\n", nom.FormatDuration(150*time.Second))
}
