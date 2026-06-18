package tui

import (
	"log/slog"

	tea "charm.land/bubbletea/v2"
)

// ============================================================================
// BUBBLE TEA LIFECYCLE METHODS
// ============================================================================
// ensureStarted starts the TUI program if not already started.
func (pr *BubbleTeaProgressReporter) ensureStarted() {
	pr.mu.RLock()

	if !pr.started {
		pr.mu.RUnlock()
		pr.mu.Lock()
		if !pr.started {
			// Start TUI program
			pr.program = tea.NewProgram(pr.model)
			// Run program in goroutine - Bubble Tea handles graceful shutdown internally
			go func() {
				if _, err := pr.program.Run(); err != nil {
					// TUI errors are non-fatal, continue execution
					slog.Error("TUI program error", "error", err)
				}
			}()

			pr.started = true
		}
		pr.mu.Unlock()
	} else {
		pr.mu.RUnlock()
	}
}

// Stop gracefully stops the TUI progress display with completion state.
func (pr *BubbleTeaProgressReporter) Stop() {
	pr.mu.Lock()

	canComplete := pr.workflowState.CanTransitionTo(WorkflowStateCompleted)
	if canComplete {
		pr.workflowState = WorkflowStateCompleted
	}

	nomMode := pr.model.displayMode == DisplayModeNOM
	prog := pr.program
	pr.mu.Unlock()

	// The model self-transitions to Completed from the 100% progress update.
	pr.send(ProgressUpdateMsg{
		Type:     ProgressUpdate,
		Progress: 100.0,
	})

	// In NOM mode with a real program, signal graceful shutdown.
	if prog != nil && nomMode {
		prog.Send(tea.Quit)
	}
}

// Start initializes and starts the TUI program
// This is called explicitly in NOM mode before workflow execution begins
// In normal mode, the TUI starts lazily on the first progress report.
func (pr *BubbleTeaProgressReporter) Start() {
	pr.ensureStarted()
}
