package tui

import (
	"context"
	"time"

	"github.com/larsartmann/go-output/nom"
)

// ============================================================================
// STRONG TYPES FOR WORKFLOW STATE MANAGEMENT
// ============================================================================
// workflowState represents the current state of the workflow execution.
// It is intentionally SEPARATE from nom.ActivityStatus (which tracks a single
// activity's lifecycle) — do NOT merge them. They share Running and Completed
// values but model different domains: a workflow Errored while individual
// activities report Failed; a workflow has no Paused state. Split-brain m3:
// documented to prevent a future "unify these enums" mistake.
//
// This makes invalid states unrepresentable through strong typing.
type workflowState uint8

const (
	workflowStateIdle workflowState = iota
	workflowStateRunning
	workflowStateCompleted
	workflowStateErrored
)

// workflowState String Constants.
const (
	workflowStateStringIdle      = "idle"
	workflowStateStringRunning   = "running"
	workflowStateStringCompleted = "completed"
	workflowStateStringErrored   = "errored"
	// workflowStateStringUnknown is the fallback for unrecognized workflow states.
	// Mirrors nom.StatusStringUnknown — both use "unknown" for the same purpose.
	workflowStateStringUnknown = "unknown"
)

// String returns the string representation of the workflow state.
func (ws workflowState) String() string {
	switch ws {
	case workflowStateIdle:
		return workflowStateStringIdle
	case workflowStateRunning:
		return workflowStateStringRunning
	case workflowStateCompleted:
		return workflowStateStringCompleted
	case workflowStateErrored:
		return workflowStateStringErrored
	default:
		return workflowStateStringUnknown
	}
}

// canAcceptUpdates returns whether this state allows progress updates.
func (ws workflowState) canAcceptUpdates() bool {
	switch ws {
	case workflowStateIdle, workflowStateRunning:
		return true
	case workflowStateCompleted, workflowStateErrored:
		return false
	default:
		return false
	}
}

// canAcceptTicks returns whether this state allows timer ticks.
func (ws workflowState) canAcceptTicks() bool {
	return ws == workflowStateIdle || ws == workflowStateRunning
}

// canTransitionTo checks if a transition to another state is valid.
func (ws workflowState) canTransitionTo(newState workflowState) bool {
	switch ws {
	case workflowStateIdle:
		return newState == workflowStateRunning
	case workflowStateRunning:
		return newState == workflowStateCompleted || newState == workflowStateErrored
	case workflowStateCompleted, workflowStateErrored:
		return false // Terminal states
	default:
		return false
	}
}

// progressStep represents a step-based progress item.
type progressStep struct {
	Current     uint
	Total       uint
	Message     string
	StartTime   time.Time
	CompletedAt *time.Time
}

// isActive returns true if the step is still in progress (not completed).
// Derived from CompletedAt to make the impossible state
// (CompletedAt != nil && active) unrepresentable.
func (s progressStep) isActive() bool {
	return s.CompletedAt == nil
}

// ProgressModel holds the state for the Bubble Tea progress display.
type ProgressModel struct {
	// Core progress data
	currentProgress float64
	steps           []progressStep
	// Display state
	startTime    time.Time
	lastUpdate   time.Time
	width        int
	height       int
	scrollOffset int
	// Status tracking with strong types
	workflowState  workflowState
	currentMessage string
	// Display mode determines which visualization style to render
	// Using enum prevents invalid states (e.g., NOM mode without subscriber)
	displayMode DisplayMode
	// NOM-style visualization fields
	// dependencyTree is a cached *pointer* to the subscriber's tree (not a copy).
	// Refreshed each tick via syncNOMSubscriber(). Safe to read between ticks.
	// Activity counts are fetched on-demand via GetActivityCounts() — no cache needed.
	dependencyTree *nom.DependencyTree
	nomSubscriber  *nom.NOMSubscriber
	showHelp       bool
	cancelFunc     context.CancelFunc
	selectedNode   nom.ActivityID
	visibleEntries []nom.VisibleEntry
	treeStartLine  int
	theme          nom.Theme
	// criticalPathFilter narrows the NOM view to critical-path nodes only.
	criticalPathFilter bool
	// dotExportPath holds the last DOT export path (for user feedback).
	dotExportPath string
}
