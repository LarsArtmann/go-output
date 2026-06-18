package nom

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output"
)

func TestNewActivity(t *testing.T) {
	t.Parallel()

	a := NewActivity("build", "Build Module")

	if a.ID.Get() != "build" {
		t.Errorf("ID = %q, want %q", a.ID.Get(), "build")
	}

	if a.Label.Get() != "Build Module" {
		t.Errorf("Label = %q, want %q", a.Label.Get(), "Build Module")
	}

	if a.Status != ActivityStatusPending {
		t.Errorf("Status = %v, want %v", a.Status, ActivityStatusPending)
	}

	if a.Shape != output.NodeShapeEllipse {
		t.Errorf("Shape = %v, want %v (pending=ellipse)", a.Shape, output.NodeShapeEllipse)
	}
}

func TestActivity_SetRunning(t *testing.T) {
	t.Parallel()

	a := NewActivity("test", "Run Tests")
	a.SetRunning()

	if !a.IsRunning() {
		t.Error("expected running")
	}

	if a.StartedAt.IsZero() {
		t.Error("StartedAt should be set")
	}

	if a.Shape != output.NodeShapeBox {
		t.Errorf("Shape = %v, want Box (running)", a.Shape)
	}

	if a.Style.Fill != "#16a34a" {
		t.Errorf("Fill = %q, want green", a.Style.Fill)
	}
}

func TestActivity_SetCompleted(t *testing.T) {
	t.Parallel()

	a := NewActivity("deploy", "Deploy")
	a.SetRunning()
	time.Sleep(time.Millisecond)
	a.SetCompleted()

	if !a.IsCompleted() {
		t.Error("expected completed")
	}

	if a.EndedAt.IsZero() {
		t.Error("EndedAt should be set")
	}

	if !a.StartedAt.Before(a.EndedAt) {
		t.Error("StartedAt should be before EndedAt")
	}

	if a.Shape != output.NodeShapeRect {
		t.Errorf("Shape = %v, want Rect (completed)", a.Shape)
	}
}

func TestActivity_SetFailed(t *testing.T) {
	t.Parallel()

	a := NewActivity("lint", "Lint Code")
	a.SetRunning()
	a.SetFailed(errTestFailure)

	if !a.IsFailed() {
		t.Error("expected failed")
	}

	if a.Err == nil {
		t.Error("Err should be set")
	}

	if a.Shape != output.NodeShapeDiamond {
		t.Errorf("Shape = %v, want Diamond (failed)", a.Shape)
	}

	if a.Style.Fill != "#dc2626" {
		t.Errorf("Fill = %q, want red", a.Style.Fill)
	}
}

func TestActivity_Elapsed(t *testing.T) {
	t.Parallel()

	a := NewActivity("task", "Task")
	if a.Elapsed() != 0 {
		t.Error("pending activity should have zero elapsed")
	}

	a.SetRunning()
	time.Sleep(time.Millisecond)

	elapsed := a.Elapsed()
	if elapsed <= 0 {
		t.Error("running activity should have positive elapsed")
	}

	a.SetCompleted()

	finalElapsed := a.Elapsed()
	if finalElapsed <= 0 {
		t.Error("completed activity should have positive elapsed")
	}
}

func TestActivityStatus_NodeShape(t *testing.T) {
	t.Parallel()

	cases := []struct {
		status ActivityStatus
		shape  output.NodeShape
	}{
		{ActivityStatusPending, output.NodeShapeEllipse},
		{ActivityStatusRunning, output.NodeShapeBox},
		{ActivityStatusCompleted, output.NodeShapeRect},
		{ActivityStatusFailed, output.NodeShapeDiamond},
		{ActivityStatusPaused, output.NodeShapeHexagon},
	}

	for _, tc := range cases {
		got := tc.status.NodeShape()
		if got != tc.shape {
			t.Errorf("%s.NodeShape() = %v, want %v", tc.status, got, tc.shape)
		}
	}
}

func TestActivityStatus_GraphStyle(t *testing.T) {
	t.Parallel()

	cases := []struct {
		status ActivityStatus
	}{
		{ActivityStatusPending},
		{ActivityStatusRunning},
		{ActivityStatusCompleted},
		{ActivityStatusFailed},
		{ActivityStatusPaused},
	}

	for _, tc := range cases {
		got := tc.status.GraphStyle()
		if tc.status == ActivityStatusFailed && got.Fill == "" {
			t.Errorf("%s.GraphStyle() should have non-empty Fill", tc.status)
		}

		if tc.status == ActivityStatusRunning && got.Fill == "" {
			t.Errorf("%s.GraphStyle() should have non-empty Fill", tc.status)
		}
	}
}

var errTestFailure = &testError{"test failure"}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
