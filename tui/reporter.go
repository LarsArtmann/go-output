package tui

import (
	"context"
	"fmt"
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
type BubbleTeaProgressReporter struct {
	mu      sync.RWMutex
	model   *ProgressModel
	program *tea.Program
	started bool
}

// NewProgressModel creates a new ProgressModel with default initialization.
func NewProgressModel() *ProgressModel {
	return &ProgressModel{
		messages:       make([]string, 0),
		steps:          make([]ProgressStep, 0),
		startTime:      time.Now(),
		lastUpdate:     time.Now(),
		workflowState:  WorkflowStateIdle,
		displayMode:    DisplayModeUniversal,
		activities:     make(map[nom.ActivityID]*nom.ActivityDisplayState),
		dependencyTree: nom.NewDependencyTree(),
		timingCache:    nom.NewTimingCache(),
		nomSubscriber:  nom.NewNOMStyleSubscriber(),
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
		model:   NewProgressModel(),
		started: false,
	}
}

// Subscriber returns the internal NOM subscriber used by the TUI model.
// Send events to this subscriber to populate the NOM-style dependency tree.
func (pr *BubbleTeaProgressReporter) Subscriber() *nom.NOMStyleSubscriber {
	return pr.model.nomSubscriber
}

// SetCancelFunc sets the context cancellation function called on ctrl+c.
// This allows the TUI to cancel the running workflow when the user presses ctrl+c.
func (pr *BubbleTeaProgressReporter) SetCancelFunc(fn context.CancelFunc) {
	pr.model.cancelFunc = fn
}

// SetDisplayMode switches the rendering mode between NOM and Universal.
// Must be called before Start(). DisplayModeNOM renders the dependency tree;
// DisplayModeUniversal renders step-based progress.
func (pr *BubbleTeaProgressReporter) SetDisplayMode(mode DisplayMode) {
	pr.model.displayMode = mode
}

// transitionWorkflowState safely transitions the workflow to a new state.
func (pr *BubbleTeaProgressReporter) transitionWorkflowState(newState WorkflowState) bool {
	if pr.model.workflowState.CanTransitionTo(newState) {
		pr.model.workflowState = newState
		return true
	}

	return false
}

// isWorkflowActive returns true if the workflow is in a state that accepts updates.
func (pr *BubbleTeaProgressReporter) isWorkflowActive() bool {
	return pr.model.workflowState.CanAcceptUpdates()
}

// sendToProgram sends a progress update message to the Bubble Tea program
// This helper eliminates double-locking bug by locking once per call.
func (pr *BubbleTeaProgressReporter) sendToProgram(msg ProgressUpdateMsg) {
	pr.mu.RLock()
	defer pr.mu.RUnlock()

	if pr.program != nil {
		pr.program.Send(msg)
	}
}

// ReportError transitions the workflow to error state (not part of universal-workflow interface).
func (pr *BubbleTeaProgressReporter) ReportError(err error) {
	pr.ensureStarted()
	// Transition to error state
	if pr.transitionWorkflowState(WorkflowStateErrored) {
		// Update current message to show the error
		pr.model.currentMessage = fmt.Sprintf("Error: %v", err)
		pr.sendToProgram(ProgressUpdateMsg{
			Type:    MessageUpdate,
			Message: pr.model.currentMessage,
		})
	}
}

// ============================================================================
// UNIVERSAL-WORKFLOW PROGRESSREPORTER INTERFACE IMPLEMENTATION
// ============================================================================
// ReportProgress implements universal-workflow ProgressReporter interface
// Reports workflow execution progress percentage (0.0 to 100.0).
func (pr *BubbleTeaProgressReporter) ReportProgress(percent float64) {
	pr.ensureStarted()
	// Transition to running state if this is the first progress report
	if pr.model.workflowState == WorkflowStateIdle {
		pr.transitionWorkflowState(WorkflowStateRunning)
	}
	// Only accept updates if workflow is active
	if !pr.isWorkflowActive() {
		return
	}

	pr.model.currentProgress = percent
	// Check for completion and transition state
	if percent >= 100.0 {
		pr.transitionWorkflowState(WorkflowStateCompleted)
	}

	pr.sendToProgram(ProgressUpdateMsg{
		Type:     ProgressUpdate,
		Progress: percent,
	})
}

// ReportMessage implements universal-workflow ProgressReporter interface
// Reports a progress message.
func (pr *BubbleTeaProgressReporter) ReportMessage(message string) {
	pr.ensureStarted()
	// Transition to running state if this is the first message
	if pr.model.workflowState == WorkflowStateIdle {
		pr.transitionWorkflowState(WorkflowStateRunning)
	}
	// Only accept updates if workflow is active
	if !pr.isWorkflowActive() {
		return
	}

	pr.model.currentMessage = message
	pr.model.messages = append(pr.model.messages, fmt.Sprintf("[%s] %s",
		time.Now().Format("15:04:05"), message))
	pr.sendToProgram(ProgressUpdateMsg{
		Type:    MessageUpdate,
		Message: message,
	})
}

// ReportStep implements universal-workflow ProgressReporter interface
// Reports step-based progress with current/total counters.
func (pr *BubbleTeaProgressReporter) ReportStep(current, total uint, message string) {
	pr.ensureStarted()
	// Transition to running state if this is the first step
	if pr.model.workflowState == WorkflowStateIdle {
		pr.transitionWorkflowState(WorkflowStateRunning)
	}
	// Only accept updates if workflow is active
	if !pr.isWorkflowActive() {
		return
	}
	// Update or create step
	stepFound := false

	for i := range pr.model.steps {
		if pr.model.steps[i].Message == message || pr.model.steps[i].IsActive {
			pr.model.steps[i].Current = current
			pr.model.steps[i].Total = total
			pr.model.steps[i].Message = message

			pr.model.steps[i].IsActive = current < total
			if current >= total && pr.model.steps[i].CompletedAt == nil {
				now := time.Now()
				pr.model.steps[i].CompletedAt = &now
			}

			stepFound = true

			break
		}
	}

	if !stepFound {
		step := ProgressStep{
			Current:   current,
			Total:     total,
			Message:   message,
			StartTime: time.Now(),
			IsActive:  current < total,
		}
		pr.model.steps = append(pr.model.steps, step)
	}

	pr.sendToProgram(ProgressUpdateMsg{
		Type:    StepUpdate,
		Current: current,
		Total:   total,
		Message: message,
	})
}
