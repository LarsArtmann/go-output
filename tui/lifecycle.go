package tui
import (
"log"
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
					log.Printf("TUI program error: %v", err)
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
	defer pr.mu.Unlock()
	if pr.program != nil {
		// Transition to completed state
		pr.transitionWorkflowState(WorkflowStateCompleted)
		pr.program.Send(ProgressUpdateMsg{
			Type:     ProgressUpdate,
			Progress: 100.0,
		})
		// In NOM mode, auto-exit after completion
		if pr.model.displayMode == DisplayModeNOM {
			pr.program.Send(tea.Quit)
		}
	}
}
// Start initializes and starts the TUI program
// This is called explicitly in NOM mode before workflow execution begins
// In normal mode, the TUI starts lazily on the first progress report.
func (pr *BubbleTeaProgressReporter) Start() {
	pr.ensureStarted()
}
