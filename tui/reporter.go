package tui

import (
	"context"
	"sync"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/larsartmann/go-output/nom"
)

// ============================================================================
// BUBBLE TEA PROGRESS REPORTER IMPLEMENTATION
// ============================================================================
// ============================================================================
// BUBBLE TEA PROGRESS REPORTER CORE
// ============================================================================
// BubbleTeaProgressReporter implements universal-workflow's ProgressReporter interface
// using Bubble Tea for rich TUI display inspired by `nh darwin switch`.
//
// Concurrency model:
//   - The reporter owns a workflowState field (protected by mu) used solely for
//     decision logic (should this update be accepted?).
//   - The model owns its own workflowState on the TUI goroutine for rendering.
//   - ALL model field mutations flow through send(), which dispatches to either
//     prog.Send (production: queued for TUI goroutine) or model.Update (test
//     mode: synchronous). The reporter NEVER touches model fields directly.
//   - This eliminates the data race where the caller goroutine and the TUI
//     goroutine both mutated shared model fields.
type BubbleTeaProgressReporter struct {
	mu            sync.RWMutex
	model         *ProgressModel
	program       *tea.Program
	started       bool
	workflowState workflowState // reporter-authoritative state for decision logic
}

// NewProgressModel creates a new ProgressModel with default initialization.
func NewProgressModel() *ProgressModel {
	return &ProgressModel{
		steps:          make([]progressStep, 0),
		startTime:      time.Now(),
		lastUpdate:     time.Now(),
		workflowState:  workflowStateIdle,
		displayMode:    DisplayModeUniversal,
		dependencyTree: nom.NewDependencyTree(),
		nomSubscriber:  nom.NewNOMSubscriber(),
		theme:          nom.ThemeDefault,
	}
}

// ============================================================================
// CONSTRUCTOR
// ============================================================================
// NewBubbleTeaProgressReporter creates a new TUI progress reporter that implements
// the universal-workflow ProgressReporter interface.
func NewBubbleTeaProgressReporter() *BubbleTeaProgressReporter {
	// Don't start the program yet - wait for first progress report
	return &BubbleTeaProgressReporter{
		model:         NewProgressModel(),
		workflowState: workflowStateIdle,
	}
}

// Subscriber returns the internal NOM subscriber used by the TUI model.
// Send events to this subscriber to populate the NOM-style dependency tree.
// This is the public API for consumers (e.g. BuildFlow's ProgressBridge) that
// dispatch lifecycle events to the visual renderer.
func (pr *BubbleTeaProgressReporter) Subscriber() *nom.NOMSubscriber {
	return pr.model.nomSubscriber
}

// SetCancelFunc sets the context cancellation function called on ctrl+c.
// This allows the TUI to cancel the running workflow when the user presses ctrl+c.
// Must be called before Start().
func (pr *BubbleTeaProgressReporter) SetCancelFunc(fn context.CancelFunc) {
	pr.mu.Lock()
	pr.model.cancelFunc = fn
	pr.mu.Unlock()
}

// SetDisplayMode switches the rendering mode between NOM and Universal.
// Must be called before Start(). DisplayModeNOM renders the dependency tree;
// DisplayModeUniversal renders step-based progress. Resets scrollOffset to
// prevent a stale position from one mode carrying into the other.
func (pr *BubbleTeaProgressReporter) SetDisplayMode(mode DisplayMode) {
	pr.mu.Lock()
	pr.model.displayMode = mode
	pr.model.scrollOffset = 0
	pr.mu.Unlock()
}

// transitionWorkflowState transitions the reporter's authoritative workflow
// state and sends a stateTransitionMsg so the model stays in sync.
func (pr *BubbleTeaProgressReporter) transitionWorkflowState(newState workflowState) bool {
	pr.mu.Lock()

	ok := pr.workflowState.canTransitionTo(newState)
	if ok {
		pr.workflowState = newState
	}
	pr.mu.Unlock()

	if ok {
		pr.send(stateTransitionMsg{NewState: newState})
	}

	return ok
}

// isWorkflowActive returns true if the workflow is in a state that accepts updates.
func (pr *BubbleTeaProgressReporter) isWorkflowActive() bool {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	return pr.workflowState.canAcceptUpdates()
}

// send dispatches a message to the TUI program if running, or processes it
// synchronously on the model when no program is attached (test/pre-start mode).
// The exclusive lock serializes model.Update calls in test mode, preventing
// concurrent field mutations. In production, prog.Send is a non-blocking
// channel write so the lock contention is negligible.
func (pr *BubbleTeaProgressReporter) send(msg tea.Msg) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if pr.program != nil {
		pr.program.Send(msg)
		return
	}

	// Test/pre-start mode: apply directly to the model so tests can verify state.
	pr.model.Update(msg)
}

// ReportError transitions the workflow to error state (not part of universal-workflow interface).
func (pr *BubbleTeaProgressReporter) ReportError(err error) {
	pr.ensureStarted()

	pr.mu.Lock()

	canError := pr.workflowState.canTransitionTo(workflowStateErrored)
	if canError {
		pr.workflowState = workflowStateErrored
	}
	pr.mu.Unlock()

	if canError {
		pr.send(errorMsg{Err: err})
	}
}

// ============================================================================
// UNIVERSAL-WORKFLOW PROGRESSREPORTER INTERFACE IMPLEMENTATION
// ============================================================================
// ensureStartedAndActive ensures the program is started, transitions to Running
// if currently Idle, and returns true only if the workflow accepts updates.
func (pr *BubbleTeaProgressReporter) ensureStartedAndActive() bool {
	pr.ensureStarted()

	pr.mu.Lock()

	needStart := pr.workflowState == workflowStateIdle
	if needStart {
		pr.workflowState = workflowStateRunning
	}

	active := pr.workflowState.canAcceptUpdates()
	pr.mu.Unlock()

	if needStart {
		pr.send(stateTransitionMsg{NewState: workflowStateRunning})
	}

	return active
}

// ReportProgress implements universal-workflow ProgressReporter interface
// Reports workflow execution progress percentage (0.0 to 100.0).
func (pr *BubbleTeaProgressReporter) ReportProgress(percent float64) {
	if !pr.ensureStartedAndActive() {
		return
	}

	pr.mu.Lock()
	if percent >= 100.0 && pr.workflowState.canTransitionTo(workflowStateCompleted) {
		pr.workflowState = workflowStateCompleted
	}
	pr.mu.Unlock()

	// The model self-transitions to Completed when it receives progress >= 100.
	pr.send(progressUpdateMsg{
		Type:     progressUpdate,
		Progress: percent,
	})
}

// ReportMessage implements universal-workflow ProgressReporter interface
// Reports a progress message.
func (pr *BubbleTeaProgressReporter) ReportMessage(message string) {
	if !pr.ensureStartedAndActive() {
		return
	}

	pr.send(progressUpdateMsg{
		Type:    messageUpdate,
		Message: message,
	})
}

// ReportStep implements universal-workflow ProgressReporter interface
// Reports step-based progress with current/total counters.
func (pr *BubbleTeaProgressReporter) ReportStep(current, total uint, message string) {
	if !pr.ensureStartedAndActive() {
		return
	}

	pr.send(stepUpdateMsg{
		Current: current,
		Total:   total,
		Message: message,
	})
}
