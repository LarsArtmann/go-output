package nom

import (
	"testing"
	"time"

	"github.com/larsartmann/go-output"
)

// assertActivityShape checks that a.Shape equals want, failing with label.
func assertActivityShape(t *testing.T, a *Activity, want output.NodeShape, label string) {
	t.Helper()

	if a.Shape != want {
		t.Errorf("Shape = %v, want %s", a.Shape, label)
	}
}

// assertActivityFill checks that a.Style.Fill equals want.
func assertActivityFill(t *testing.T, a *Activity, want string) {
	t.Helper()

	if a.Style.Fill != want {
		t.Errorf("Fill = %q, want %q", a.Style.Fill, want)
	}
}

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

	if a.StartTime.IsZero() {
		t.Error("StartTime should be set")
	}

	assertActivityShape(t, a, output.NodeShapeBox, "Box (running)")
	assertActivityFill(t, a, "#16a34a")
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

	if a.EndTime.IsZero() {
		t.Error("EndTime should be set")
	}

	if !a.StartTime.Before(a.EndTime) {
		t.Error("StartTime should be before EndTime")
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

	assertActivityShape(t, a, output.NodeShapeDiamond, "Diamond (failed)")
	assertActivityFill(t, a, "#dc2626")
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
