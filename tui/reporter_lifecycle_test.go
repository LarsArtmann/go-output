package tui

import (
	"testing"

	"github.com/larsartmann/go-output/nom"
)

func TestBubbleTeaProgressReporter_Subscriber(t *testing.T) {
	t.Parallel()

	pr := NewBubbleTeaProgressReporter()

	sub := pr.Subscriber()
	if sub == nil {
		t.Fatal("Subscriber() should not return nil")
	}

	if sub != pr.model.nomSubscriber {
		t.Error("Subscriber() should return the model's nomSubscriber")
	}
}

func TestBubbleTeaProgressReporter_SetCancelFunc(t *testing.T) {
	t.Parallel()

	pr := NewBubbleTeaProgressReporter()

	called := false

	pr.SetCancelFunc(func() { called = true })

	if pr.model.cancelFunc == nil {
		t.Fatal("SetCancelFunc should assign the cancel function")
	}

	// Invoke the cancel function to verify it is the one we set.
	pr.model.cancelFunc()

	if !called {
		t.Error("cancel function should have been called")
	}
}

func TestBubbleTeaProgressReporter_SetCancelFunc_Nil(t *testing.T) {
	t.Parallel()

	pr := NewBubbleTeaProgressReporter()

	// Setting nil is a valid operation (e.g., clearing the cancel func).
	pr.SetCancelFunc(nil)

	if pr.model.cancelFunc != nil {
		t.Error("SetCancelFunc(nil) should assign nil")
	}
}

func TestBubbleTeaProgressReporter_SetDisplayMode(t *testing.T) {
	t.Parallel()

	pr := NewBubbleTeaProgressReporter()

	pr.SetDisplayMode(DisplayModeNOM)

	if pr.model.displayMode != DisplayModeNOM {
		t.Errorf("displayMode = %v, want %v", pr.model.displayMode, DisplayModeNOM)
	}

	pr.SetDisplayMode(DisplayModeUniversal)

	if pr.model.displayMode != DisplayModeUniversal {
		t.Errorf("displayMode = %v, want %v", pr.model.displayMode, DisplayModeUniversal)
	}
}

func TestBubbleTeaProgressReporter_EnsureStarted_Idempotent(t *testing.T) {
	t.Parallel()

	pr := NewBubbleTeaProgressReporter()

	// First call: starts the TUI.
	pr.ensureStarted()

	// Capture started state — second call should be a no-op.
	startedAfterFirst := pr.started
	if !startedAfterFirst {
		t.Error("ensureStarted() should set started=true on first call")
	}

	// Second call should not panic, even when started is already true.
	pr.ensureStarted()

	if !pr.started {
		t.Error("started should remain true after second ensureStarted()")
	}

	// Verify Subscriber is reachable through the started reporter.
	if pr.Subscriber() == nil {
		t.Error("Subscriber() should still return a valid subscriber")
	}
}

func TestBubbleTeaProgressReporter_StartNoOpIfAlreadyStarted(t *testing.T) {
	t.Parallel()

	pr := NewBubbleTeaProgressReporter()

	// Pre-start the reporter.
	pr.started = true

	// Start() should be a no-op when already started (no panic, no double-launch).
	pr.Start()

	if !pr.started {
		t.Error("Start() should not unset started flag")
	}
}

func TestNewProgressModel_HasNOMSubscriber(t *testing.T) {
	t.Parallel()

	model := NewProgressModel()
	if model.nomSubscriber == nil {
		t.Fatal("NewProgressModel should initialize nomSubscriber")
	}
}

// Sanity check: nom is used to verify reporter integrates with nom subscriber.
var _ = nom.NewNOMStyleSubscriber
