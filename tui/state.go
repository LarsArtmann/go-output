package tui

import (
	"context"
	"time"

	"github.com/larsartmann/go-output/nom"
)

// ============================================================================
// STRONG TYPES FOR WORKFLOW STATE MANAGEMENT
// ============================================================================
// WorkflowState represents the current state of the workflow execution
// This makes invalid states unrepresentable through strong typing.
type WorkflowState uint8

const (
	WorkflowStateIdle WorkflowState = iota
	WorkflowStateRunning
	WorkflowStateCompleted
	WorkflowStateErrored
)

// WorkflowState String Constants.
const (
	WorkflowStateStringIdle      = "idle"
	WorkflowStateStringRunning   = "running"
	WorkflowStateStringCompleted = "completed"
	WorkflowStateStringErrored   = "errored"
	// WorkflowStateStringUnknown is the fallback for unrecognized workflow states.
	// Mirrors nom.StatusStringUnknown — both use "unknown" for the same purpose.
	WorkflowStateStringUnknown = "unknown"
)

// String returns the string representation of the workflow state.
func (ws WorkflowState) String() string {
	switch ws {
	case WorkflowStateIdle:
		return WorkflowStateStringIdle
	case WorkflowStateRunning:
		return WorkflowStateStringRunning
	case WorkflowStateCompleted:
		return WorkflowStateStringCompleted
	case WorkflowStateErrored:
		return WorkflowStateStringErrored
	default:
		return WorkflowStateStringUnknown
	}
}

// CanAcceptUpdates returns whether this state allows progress updates.
func (ws WorkflowState) CanAcceptUpdates() bool {
	switch ws {
	case WorkflowStateIdle, WorkflowStateRunning:
		return true
	case WorkflowStateCompleted, WorkflowStateErrored:
		return false
	default:
		return false
	}
}

// CanAcceptTicks returns whether this state allows timer ticks.
func (ws WorkflowState) CanAcceptTicks() bool {
	return ws == WorkflowStateIdle || ws == WorkflowStateRunning
}

// CanTransitionTo checks if a transition to another state is valid.
func (ws WorkflowState) CanTransitionTo(newState WorkflowState) bool {
	switch ws {
	case WorkflowStateIdle:
		return newState == WorkflowStateRunning
	case WorkflowStateRunning:
		return newState == WorkflowStateCompleted || newState == WorkflowStateErrored
	case WorkflowStateCompleted, WorkflowStateErrored:
		return false // Terminal states
	default:
		return false
	}
}

// ProgressStep represents a step-based progress item.
type ProgressStep struct {
	Current     uint
	Total       uint
	Message     string
	StartTime   time.Time
	CompletedAt *time.Time
}

// IsActive returns true if the step is still in progress (not completed).
// Derived from CompletedAt to make the impossible state
// (CompletedAt != nil && active) unrepresentable.
func (s ProgressStep) IsActive() bool {
	return s.CompletedAt == nil
}

// ProgressModel holds the state for the Bubble Tea progress display.
type ProgressModel struct {
	// Core progress data
	currentProgress float64
	steps           []ProgressStep
	// Display state
	startTime    time.Time
	lastUpdate   time.Time
	width        int
	height       int
	scrollOffset int
	// Status tracking with strong types
	workflowState  WorkflowState
	currentMessage string
	// Display mode determines which visualization style to render
	// Using enum prevents invalid states (e.g., NOM mode without subscriber)
	displayMode DisplayMode
	// NOM-style visualization fields
	// dependencyTree is a cached *pointer* to the subscriber's tree (not a copy).
	// Refreshed each tick via syncNOMSubscriber(). Safe to read between ticks.
	// Activity counts are fetched on-demand via GetActivityCounts() — no cache needed.
	dependencyTree *nom.DependencyTree
	nomSubscriber  *nom.NOMStyleSubscriber
	showHelp       bool
	cancelFunc     context.CancelFunc
	selectedNode   nom.ActivityID
	visibleNodes   []*nom.ActivityNode
	treeStartLine  int
}
