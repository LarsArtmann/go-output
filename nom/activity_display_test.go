package nom

import (
	"errors"
	"testing"
	"time"
)

func TestNewActivityDisplayState(t *testing.T) {
	t.Parallel()

	t.Run("initializes with pending defaults", func(t *testing.T) {
		t.Parallel()

		ads := NewActivityDisplayState(ActivityID("build"), ActivityName("Build Project"))

		if ads.ActivityID != ActivityID("build") {
			t.Errorf("ActivityID = %q, want %q", ads.ActivityID, "build")
		}

		if ads.ActivityName != ActivityName("Build Project") {
			t.Errorf("ActivityName = %q, want %q", ads.ActivityName, "Build Project")
		}

		if ads.Status != ActivityStatusPending {
			t.Errorf("Status = %v, want Pending", ads.Status)
		}

		if ads.Symbol != SymbolPaused {
			t.Errorf("Symbol = %q, want %q", ads.Symbol, SymbolPaused)
		}

		if ads.Error != nil {
			t.Errorf("Error = %v, want nil", ads.Error)
		}

		if len(ads.Dependencies) != 0 {
			t.Errorf("Dependencies = %v, want empty", ads.Dependencies)
		}
	})
}

func TestActivityDisplayState_SetRunning(t *testing.T) {
	t.Parallel()

	ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
	ads.SetRunning()

	if ads.Status != ActivityStatusRunning {
		t.Errorf("Status = %v, want Running", ads.Status)
	}

	if ads.Symbol != SymbolRunning {
		t.Errorf("Symbol = %q, want %q", ads.Symbol, SymbolRunning)
	}

	if ads.StartTime.IsZero() {
		t.Error("StartTime should be set")
	}
}

func TestActivityDisplayState_SetCompleted(t *testing.T) {
	t.Parallel()

	ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
	ads.SetRunning()
	time.Sleep(time.Millisecond)
	ads.SetCompleted()

	if ads.Status != ActivityStatusCompleted {
		t.Errorf("Status = %v, want Completed", ads.Status)
	}

	if ads.Symbol != SymbolCompleted {
		t.Errorf("Symbol = %q, want %q", ads.Symbol, SymbolCompleted)
	}

	if ads.EndTime.IsZero() {
		t.Error("EndTime should be set after completion")
	}

	if ads.CurrentElapsed <= 0 {
		t.Error("CurrentElapsed should be positive after completion")
	}
}

func TestActivityDisplayState_SetFailed(t *testing.T) {
	t.Parallel()

	ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
	testErr := errors.New("something went wrong")

	ads.SetRunning()
	ads.SetFailed(testErr)

	if ads.Status != ActivityStatusFailed {
		t.Errorf("Status = %v, want Failed", ads.Status)
	}

	if ads.Symbol != SymbolFailed {
		t.Errorf("Symbol = %q, want %q", ads.Symbol, SymbolFailed)
	}

	if !errors.Is(ads.Error, testErr) {
		t.Errorf("Error = %v, want %v", ads.Error, testErr)
	}
}

func TestActivityDisplayState_Predicates(t *testing.T) {
	t.Parallel()

	t.Run("IsRunning", func(t *testing.T) {
		t.Parallel()

		ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
		if ads.IsRunning() {
			t.Error("new activity should not be running")
		}

		ads.SetRunning()

		if !ads.IsRunning() {
			t.Error("activity should be running after SetRunning")
		}
	})

	t.Run("IsCompleted", func(t *testing.T) {
		t.Parallel()

		ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
		if ads.IsCompleted() {
			t.Error("new activity should not be completed")
		}

		ads.SetCompleted()

		if !ads.IsCompleted() {
			t.Error("activity should be completed after SetCompleted")
		}
	})

	t.Run("IsFailed", func(t *testing.T) {
		t.Parallel()

		ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
		if ads.IsFailed() {
			t.Error("new activity should not be failed")
		}

		ads.SetFailed(errors.New("fail"))

		if !ads.IsFailed() {
			t.Error("activity should be failed after SetFailed")
		}
	})
}

func TestActivityDisplayState_SetEstimatedTime(t *testing.T) {
	t.Parallel()

	ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
	dur := 5 * time.Second
	ads.SetEstimatedTime(dur)

	if ads.EstimatedTime != dur {
		t.Errorf("EstimatedTime = %v, want %v", ads.EstimatedTime, dur)
	}
}

func TestActivityDisplayState_SetOperationType(t *testing.T) {
	t.Parallel()

	ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
	ads.setOperationType(OperationTypeDownload)

	if ads.OperationType != OperationTypeDownload {
		t.Errorf("OperationType = %q, want %q", ads.OperationType, OperationTypeDownload)
	}
}

func TestActivityDisplayState_AddDependency(t *testing.T) {
	t.Parallel()

	ads := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
	ads.addDependency("dep1")
	ads.addDependency("dep2")

	if len(ads.Dependencies) != 2 {
		t.Fatalf("Dependencies len = %d, want 2", len(ads.Dependencies))
	}

	if ads.Dependencies[0] != "dep1" || ads.Dependencies[1] != "dep2" {
		t.Errorf("Dependencies = %v, want [dep1 dep2]", ads.Dependencies)
	}
}

func TestActivityDisplayState_Copy(t *testing.T) {
	t.Parallel()

	original := NewActivityDisplayState(ActivityID("a"), ActivityName("A"))
	original.SetRunning()
	original.addDependency("dep1")

	copied := original.Copy()

	if copied.ActivityID != original.ActivityID {
		t.Error("Copy should preserve ActivityID")
	}

	if copied.ActivityName != original.ActivityName {
		t.Error("Copy should preserve ActivityName")
	}

	if copied.Status != original.Status {
		t.Error("Copy should preserve Status")
	}

	copied.addDependency("dep2")

	if len(original.Dependencies) == 2 {
		t.Error("modifying copy should not affect original (deep copy failed)")
	}
}
