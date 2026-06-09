package tui

import (
	"testing"
)

func TestWorkflowState_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state WorkflowState
		want  string
	}{
		{WorkflowStateIdle, "idle"},
		{WorkflowStateRunning, "running"},
		{WorkflowStateCompleted, "completed"},
		{WorkflowStateErrored, "errored"},
		{WorkflowState(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			got := tt.state.String()
			if got != tt.want {
				t.Errorf("WorkflowState(%d).String() = %q, want %q", tt.state, got, tt.want)
			}
		})
	}
}

func TestWorkflowState_CanAcceptUpdates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state WorkflowState
		want  bool
	}{
		{WorkflowStateIdle, true},
		{WorkflowStateRunning, true},
		{WorkflowStateCompleted, false},
		{WorkflowStateErrored, false},
		{WorkflowState(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			t.Parallel()

			got := tt.state.CanAcceptUpdates()
			if got != tt.want {
				t.Errorf("CanAcceptUpdates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowState_CanAcceptTicks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		state WorkflowState
		want  bool
	}{
		{WorkflowStateIdle, true},
		{WorkflowStateRunning, true},
		{WorkflowStateCompleted, false},
		{WorkflowStateErrored, false},
	}

	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			t.Parallel()

			got := tt.state.CanAcceptTicks()
			if got != tt.want {
				t.Errorf("CanAcceptTicks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowState_CanTransitionTo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		from WorkflowState
		to   WorkflowState
		want bool
	}{
		{WorkflowStateIdle, WorkflowStateRunning, true},
		{WorkflowStateIdle, WorkflowStateCompleted, false},
		{WorkflowStateIdle, WorkflowStateErrored, false},
		{WorkflowStateRunning, WorkflowStateCompleted, true},
		{WorkflowStateRunning, WorkflowStateErrored, true},
		{WorkflowStateRunning, WorkflowStateIdle, false},
		{WorkflowStateCompleted, WorkflowStateRunning, false},
		{WorkflowStateCompleted, WorkflowStateIdle, false},
		{WorkflowStateErrored, WorkflowStateRunning, false},
		{WorkflowState(99), WorkflowStateRunning, false},
	}

	for _, tt := range tests {
		t.Run(tt.from.String()+"->"+tt.to.String(), func(t *testing.T) {
			t.Parallel()

			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.want {
				t.Errorf("CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowStateStringConstants(t *testing.T) {
	t.Parallel()

	if WorkflowStateStringIdle != "idle" {
		t.Errorf("WorkflowStateStringIdle = %q, want %q", WorkflowStateStringIdle, "idle")
	}

	if WorkflowStateStringRunning != "running" {
		t.Errorf("WorkflowStateStringRunning = %q, want %q", WorkflowStateStringRunning, "running")
	}

	if WorkflowStateStringCompleted != "completed" {
		t.Errorf("WorkflowStateStringCompleted = %q, want %q", WorkflowStateStringCompleted, "completed")
	}

	if WorkflowStateStringErrored != "errored" {
		t.Errorf("WorkflowStateStringErrored = %q, want %q", WorkflowStateStringErrored, "errored")
	}

	if WorkflowStateStringUnknown != "unknown" {
		t.Errorf("WorkflowStateStringUnknown = %q, want %q", WorkflowStateStringUnknown, "unknown")
	}
}
