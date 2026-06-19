package tui

import (
	"testing"
)

func TestWorkflowState_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state workflowState
		want  string
	}{
		{workflowStateIdle, "idle"},
		{workflowStateRunning, "running"},
		{workflowStateCompleted, "completed"},
		{workflowStateErrored, "errored"},
		{workflowState(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			got := tt.state.String()
			if got != tt.want {
				t.Errorf("workflowState(%d).String() = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

func TestWorkflowState_CanAccept(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state       workflowState
		wantUpdates bool
		wantTicks   bool
	}{
		{workflowStateIdle, true, true},
		{workflowStateRunning, true, true},
		{workflowStateCompleted, false, false},
		{workflowStateErrored, false, false},
		{workflowState(99), false, false},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			t.Parallel()

			if got := tt.state.canAcceptUpdates(); got != tt.wantUpdates {
				t.Errorf("canAcceptUpdates() = %v, want %v", got, tt.wantUpdates)
			}

			if got := tt.state.canAcceptTicks(); got != tt.wantTicks {
				t.Errorf("canAcceptTicks() = %v, want %v", got, tt.wantTicks)
			}
		})
	}
}

func TestWorkflowState_CanTransitionTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		from workflowState
		to   workflowState
		want bool
	}{
		{workflowStateIdle, workflowStateRunning, true},
		{workflowStateIdle, workflowStateCompleted, false},
		{workflowStateIdle, workflowStateErrored, false},
		{workflowStateRunning, workflowStateCompleted, true},
		{workflowStateRunning, workflowStateErrored, true},
		{workflowStateRunning, workflowStateIdle, false},
		{workflowStateCompleted, workflowStateRunning, false},
		{workflowStateCompleted, workflowStateIdle, false},
		{workflowStateErrored, workflowStateRunning, false},
		{workflowState(99), workflowStateRunning, false},
	}

	for _, tt := range tests {
		t.Run(tt.from.String()+"->"+tt.to.String(), func(t *testing.T) {
			t.Parallel()

			got := tt.from.canTransitionTo(tt.to)
			if got != tt.want {
				t.Errorf("canTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowStateStringConstants(t *testing.T) {
	t.Parallel()

	if workflowStateStringIdle != "idle" {
		t.Errorf("workflowStateStringIdle = %q, want %q", workflowStateStringIdle, "idle")
	}

	if workflowStateStringRunning != "running" {
		t.Errorf("workflowStateStringRunning = %q, want %q", workflowStateStringRunning, "running")
	}

	if workflowStateStringCompleted != "completed" {
		t.Errorf("workflowStateStringCompleted = %q, want %q", workflowStateStringCompleted, "completed")
	}

	if workflowStateStringErrored != "errored" {
		t.Errorf("workflowStateStringErrored = %q, want %q", workflowStateStringErrored, "errored")
	}

	if workflowStateStringUnknown != "unknown" {
		t.Errorf("workflowStateStringUnknown = %q, want %q", workflowStateStringUnknown, "unknown")
	}
}
