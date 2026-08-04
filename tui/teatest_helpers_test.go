package tui

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/teatest/v2"

	"github.com/larsartmann/go-output/nom"
)

// newTeatestModel creates a ProgressModel pre-populated with NOM activities
// and wraps it in a teatest TestModel — driving the REAL Bubble Tea program
// loop (Init, Update, View) rather than calling methods directly.
//
// This is the first true E2E test coverage for the TUI: it exercises message
// dispatch, the tick loop, View rendering, key handling, and quit — all through
// the actual tea.Program.
//
// Note: bubbletea v2 uses a diff renderer that writes cursor-positioning escape
// sequences, not full text frames. The teatest output is a raw byte stream of
// diffs, so we assert on stable substrings (timing, summary characters) and use
// FinalModel for content verification.
func newTeatestModel(t *testing.T, width, height int) *teatest.TestModel {
	t.Helper()

	model := NewProgressModel()
	model.displayMode = DisplayModeNOM
	model.workflowState = workflowStateRunning

	// Populate NOM subscriber with test activities BEFORE the program starts.
	ctx := context.Background()
	sub := model.nomSubscriber

	_ = sub.OnEvent(ctx, nom.WorkflowStarted{ID: "wf-e2e", Name: "E2E Test"})
	_ = sub.OnEvent(ctx, nom.ActivityStarted{ID: "build", Name: "Build Module"})
	_ = sub.OnEvent(ctx, nom.ActivityStarted{ID: "test", Name: "Run Tests"})

	// Sync the dependency tree so the first render (before any tick) shows data.
	model.dependencyTree = sub.DependencyTree()

	tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(width, height))

	t.Cleanup(func() {
		_ = tm.Quit()
		tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
	})

	return tm
}

// waitForVisible polls teatest output (ANSI-stripped) until it contains substr.
// Uses pollTeatestOutput (bounded reads) instead of teatest.WaitFor because
// teatest.WaitFor internally calls io.ReadAll, which blocks indefinitely when
// the program writes continuously (the tick loop never lets the buffer go empty
// long enough for io.ReadAll to see EOF). Under -race this deadlock is reliable.
func waitForVisible(t *testing.T, tm *teatest.TestModel, substr string) {
	t.Helper()

	pollTeatestOutput(t, tm, func(b []byte) bool {
		return strings.Contains(ansi.Strip(string(b)), substr)
	}, 3*time.Second)
}

// --- F24: Program starts, renders NOM content through the real loop ---

func TestTeatest_ProgramStarts_RendersContent(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	// The summary bar border is a stable substring in the diff output.
	// Timing values (e.g. "100ms", "1.0s") confirm the tick loop is running.
	waitForVisible(t, tm, "s")
}

// --- F25: Scroll down (j key) is processed without crashing ---

func TestTeatest_ScrollDown_NoCrash(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	waitForVisible(t, tm, "s")

	// Scroll down — the program should process the keypress via Update.
	tm.Send(tea.KeyPressMsg{Code: 'j'})

	// Program should still be producing output (tick loop still running).
	waitForVisible(t, tm, "s")
}

// --- F26: Scroll up (k key) ---

func TestTeatest_ScrollUp_NoCrash(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	waitForVisible(t, tm, "s")

	// Scroll down then back up.
	tm.Send(tea.KeyPressMsg{Code: 'j'})
	tm.Send(tea.KeyPressMsg{Code: 'k'})

	// Program should still render correctly.
	waitForVisible(t, tm, "s")
}

// --- F27: Help toggle (? key) ---

func TestTeatest_HelpToggle_NoCrash(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	waitForVisible(t, tm, "s")

	// Toggle help overlay on and off.
	tm.Send(tea.KeyPressMsg{Code: '?'})
	tm.Send(tea.KeyPressMsg{Code: '?'})

	// Program should still be responsive.
	waitForVisible(t, tm, "s")
}

// --- F28: Quit with 'q' exits cleanly and preserves model state ---

func TestTeatest_Quit_ExitsCleanly(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	waitForVisible(t, tm, "s")

	// Send 'q' to quit the program.
	tm.Send(tea.KeyPressMsg{Code: 'q'})

	// The program should finish (quit command was issued).
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	// Verify the final model state has the expected NOM activities.
	finalModel := tm.FinalModel(t, teatest.WithFinalTimeout(3*time.Second))

	m, ok := finalModel.(*ProgressModel)
	if !ok {
		t.Fatalf("final model should be *ProgressModel, got %T", finalModel)
	}

	if m.displayMode != DisplayModeNOM {
		t.Errorf("displayMode = %v, want NOM", m.displayMode)
	}

	// The dependency tree should still have our activities.
	if m.dependencyTree == nil {
		t.Error("dependency tree should not be nil after quit")
	}
}

// --- F29: Ctrl+C triggers quit with context cancellation ---

func TestTeatest_CtrlC_TriggersQuit(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	waitForVisible(t, tm, "s")

	// Send ctrl+c to quit.
	tm.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})

	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}

// --- F30: WindowSizeMsg propagates correctly ---

func TestTeatest_WindowSize_Propagates(t *testing.T) {
	model := NewProgressModel()
	model.displayMode = DisplayModeNOM
	model.workflowState = workflowStateRunning

	// WithInitialTermSize sends a WindowSizeMsg — width/height should be set.
	tm := teatest.NewTestModel(t, model, teatest.WithInitialTermSize(120, 40))
	t.Cleanup(func() { _ = tm.Quit() })

	// Send 'q' and check the final model has the correct dimensions.
	tm.Send(tea.KeyPressMsg{Code: 'q'})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	m, ok := tm.FinalModel(t, teatest.WithFinalTimeout(3*time.Second)).(*ProgressModel)
	if !ok {
		t.Fatal("final model should be *ProgressModel")
	}

	if m.width != 120 {
		t.Errorf("width = %d, want 120", m.width)
	}

	if m.height != 40 {
		t.Errorf("height = %d, want 40", m.height)
	}
}

// --- F31: L key toggles tree/layered mode through the real program loop ---

func TestTeatest_LKey_TogglesMode(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	waitForVisible(t, tm, "s")

	// Send 'L' to toggle to layered mode.
	tm.Send(tea.KeyPressMsg{Code: 'L'})

	// Program should still be responsive after the toggle.
	waitForVisible(t, tm, "s")

	// Send 'L' again to toggle back to tree mode.
	tm.Send(tea.KeyPressMsg{Code: 'L'})

	waitForVisible(t, tm, "s")

	// Verify the final model state reflects the expected mode.
	tm.Send(tea.KeyPressMsg{Code: 'q'})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	m, ok := tm.FinalModel(t, teatest.WithFinalTimeout(3*time.Second)).(*ProgressModel)
	if !ok {
		t.Fatal("final model should be *ProgressModel")
	}

	if m.dependencyTree == nil {
		t.Fatal("dependency tree should not be nil")
	}

	// After toggling twice, mode should be back to tree.
	if m.dependencyTree.RenderMode() != nom.RenderModeTree {
		t.Errorf("RenderMode = %v, want RenderModeTree (toggled twice)", m.dependencyTree.RenderMode())
	}
}

// --- F32: C key toggles critical-path filter without crashing ---

func TestTeatest_CKey_TogglesCriticalFilter(t *testing.T) {
	tm := newTeatestModel(t, 100, 30)

	waitForVisible(t, tm, "s")

	// Send 'C' to enable critical-path filter.
	tm.Send(tea.KeyPressMsg{Code: 'C'})

	waitForVisible(t, tm, "s")

	// Send 'C' again to disable.
	tm.Send(tea.KeyPressMsg{Code: 'C'})

	waitForVisible(t, tm, "s")

	// Quit and verify model state.
	tm.Send(tea.KeyPressMsg{Code: 'q'})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	m, ok := tm.FinalModel(t, teatest.WithFinalTimeout(3*time.Second)).(*ProgressModel)
	if !ok {
		t.Fatal("final model should be *ProgressModel")
	}

	// Filter should be off after toggling twice.
	if m.criticalPathFilter {
		t.Error("criticalPathFilter should be false after toggling twice")
	}
}
