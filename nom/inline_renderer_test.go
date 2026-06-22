package nom

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
)

func TestInlineRenderer_FirstRender_NoAnsiEscapes(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.Draw()

	output := buf.String()
	if output == "" {
		t.Fatal("first render should produce output")
	}

	if strings.Contains(output, "\033[A") {
		t.Error("first render should not move cursor up")
	}
}

func TestInlineRenderer_SubsequentRender_MovesCursor(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := newInlineTestRenderer(sub, &buf, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.Draw()
	buf.Reset()

	// Change tree state so the frame differs — frame diffing skips identical frames.
	sendActivityCompleted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"), time.Second)

	renderer.Draw()

	output := buf.String()
	if !strings.Contains(output, "\033[") || !strings.Contains(output, "A") {
		t.Errorf("second render should move cursor up, got:\n%q", output)
	}
}

func TestInlineRenderer_Finish_ClearsFrame(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.Draw()
	buf.Reset()

	renderer.Finish(nil)

	output := buf.String()
	if !strings.Contains(output, "\033[") {
		t.Error("Finish should clear the previous frame using ANSI escapes")
	}

	if !strings.Contains(output, "completed successfully") {
		t.Errorf("Finish should print success message, got:\n%s", output)
	}
}

func TestInlineRenderer_Finish_WithError(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")

	renderer.Draw()
	buf.Reset()

	renderer.Finish(errors.New("test failure"))

	output := buf.String()
	if !strings.Contains(output, "failed") {
		t.Errorf("Finish with error should print failure, got:\n%s", output)
	}
}

func TestInlineRenderer_SummaryBar(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 20)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.SetStartTime(time.Now().Add(-5 * time.Second))
	renderer.Draw()

	output := buf.String()
	if !strings.Contains(output, "╭") {
		t.Errorf("render should include summary box, got:\n%s", output)
	}

	if !strings.Contains(output, string(SymbolTiming)) {
		t.Errorf("summary should include timing, got:\n%s", output)
	}

	if !strings.Contains(output, "%") {
		t.Errorf("summary should include completion percentage, got:\n%s", output)
	}
}

func TestInlineRenderer_NilSubscriber(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	renderer := NewInlineRenderer(nil, &buf, 10)

	renderer.Draw()

	if buf.String() != "" {
		t.Error("render with nil subscriber should produce no output")
	}
}

func TestInlineRenderer_EmptyTree(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)

	renderer.Draw()

	if buf.String() != "" {
		t.Error("render with empty tree should produce no output")
	}
}

func TestInlineRenderer_StartStop_PeriodicRender(t *testing.T) {
	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := newInlineTestRenderer(sub, &buf, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.Start(ctx, 50*time.Millisecond)
	time.Sleep(180 * time.Millisecond)
	renderer.Stop()

	output := buf.String()

	// Frame diffing: the tree state is stable across ticks (no status changes,
	// no SetStartTime so no elapsed-time changes), so only the FIRST tick emits.
	// Subsequent ticks see an identical frame and correctly skip — this is the
	// core fix for the repetition bug. We verify at least 1 render happened.
	renders := strings.Count(output, "Step 1")
	if renders < 1 {
		t.Errorf("expected at least 1 periodic render, got %d renders", renders)
	}
}

func TestInlineRenderer_StartStop_Idempotent(t *testing.T) {
	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)

	ctx := context.Background()

	renderer.Start(ctx, 50*time.Millisecond)
	renderer.Start(ctx, 50*time.Millisecond)

	renderer.Stop()
	renderer.Stop()
}

func TestInlineRenderer_Refresh_TriggersRender(t *testing.T) {
	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)
	renderer.renderNotify = make(chan struct{}, 1)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.Start(ctx, time.Hour) // very long interval

	buf.Reset()
	renderer.Refresh()

	waitForRender(t, renderer, "refresh")

	renderer.Stop()

	if !strings.Contains(buf.String(), "Step 1") {
		t.Errorf("Refresh should trigger a render, got:\n%s", buf.String())
	}
}

func TestInlineRenderer_MaxHeightZero_UsesFallback(t *testing.T) {
	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 0)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")

	for i := range 60 {
		sendActivityStarted(t, sub, ctx,
			ActivityID(fmt.Sprintf("step%d", i)),
			ActivityName(fmt.Sprintf("Step %d", i)))
	}

	renderer.Draw()

	output := buf.String()

	lines := strings.Count(output, "\n") + 1
	if lines > 56 {
		t.Errorf("expected tree capped to ~50 lines with 2-line summary, got %d lines", lines)
	}
}

func TestInlineRenderer_EndToEnd_Lifecycle(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := newInlineTestRenderer(sub, &buf, 10)
	renderer.renderNotify = make(chan struct{}, 1)

	ctx := context.Background()

	// Start the workflow
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-e2e"), WorkflowName("E2E Test"))

	// Register activities with mixed states
	sendActivityStarted(t, sub, ctx, ActivityID("build"), ActivityName("Build Project"))
	sendActivityCompleted(t, sub, ctx, ActivityID("build"), ActivityName("Build Project"), 2*time.Second)
	sendActivityStarted(t, sub, ctx, ActivityID("test"), ActivityName("Run Tests"))

	renderer.SetStartTime(time.Now().Add(-3 * time.Second))

	// Start the background renderer
	renderer.Start(ctx, time.Hour)

	// Trigger initial render via refresh
	renderer.Refresh()

	waitForRender(t, renderer, "initial render")

	buf.Reset()

	// Change tree state so frame diffing emits a new frame.
	sendActivityCompleted(t, sub, ctx, ActivityID("test"), ActivityName("Run Tests"), time.Second)

	// Trigger a refresh
	renderer.Refresh()

	waitForRender(t, renderer, "refresh render")

	output := buf.String()

	if !strings.Contains(output, "Build Project") {
		t.Errorf("render should contain completed activity, got:\n%s", output)
	}

	if !strings.Contains(output, "Run Tests") {
		t.Errorf("render should contain running activity, got:\n%s", output)
	}

	if !strings.Contains(output, "%") {
		t.Errorf("summary should show completion percentage, got:\n%s", output)
	}

	// Stop background loop before finish to avoid concurrent tree access
	renderer.Stop()

	// Finish the workflow
	renderer.Finish(nil)

	finalOutput := buf.String()

	if !strings.Contains(finalOutput, "completed successfully") {
		t.Errorf("finish should print success message, got:\n%s", finalOutput[-min(len(finalOutput), 200):])
	}
}

// waitForRender blocks until the inline renderer signals that a render
// completed, or fails the test if the timeout elapses.
func waitForRender(t *testing.T, renderer *InlineRenderer, label string) {
	t.Helper()

	select {
	case <-renderer.renderNotify:
	case <-time.After(time.Second):
		t.Fatalf("%s did not complete within 1s", label)
	}
}

func TestInlineRenderer_CIMode_PlainTextNoCursorCodes(t *testing.T) {
	// detectPlainTextForWriter is evaluated at construction, so CI must be set first.
	// Not parallel: the env var affects renderer construction.
	t.Setenv("CI", "true")

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)

	if !renderer.plainText {
		t.Fatal("expected plainText=true when CI=true")
	}

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-ci"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("CI Step"))

	renderer.Draw()

	output := buf.String()

	if strings.Contains(output, ansi.CursorUp1) {
		t.Errorf("CI plain mode must not emit cursor-up codes, got:\n%q", output)
	}

	if strings.Contains(output, ansi.SetSynchronizedOutputMode) ||
		strings.Contains(output, ansi.ResetSynchronizedOutputMode) {
		t.Errorf("CI plain mode must not emit sync-region codes, got:\n%q", output)
	}

	// Plain mode still renders the activity label.
	if !strings.Contains(output, "CI Step") {
		t.Errorf("CI plain mode should still render the activity label, got:\n%q", output)
	}
}

func TestInlineRenderer_SetMaxHeight_TakesEffect(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 20)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-h"), "")

	steps := []string{"alpha", "bravo", "charlie", "delta", "echo"}
	for _, name := range steps {
		sendActivityStarted(t, sub, ctx, ActivityID(name), ActivityName(name))
	}

	// Cap to 2 lines; only the highest-priority activities should surface.
	renderer.SetMaxHeight(2)
	renderer.Draw()

	capped := buf.String()
	if capped == "" {
		t.Fatal("expected non-empty output after SetMaxHeight")
	}

	// Raise the cap and redraw — more activities should now be visible.
	renderer.SetMaxHeight(20)
	buf.Reset()
	renderer.Draw()

	full := buf.String()

	// Count how many distinct activity labels appear in each render.
	cappedVisible, fullVisible := 0, 0

	for _, name := range steps {
		if strings.Contains(capped, name) {
			cappedVisible++
		}

		if strings.Contains(full, name) {
			fullVisible++
		}
	}

	if cappedVisible >= fullVisible {
		t.Errorf("SetMaxHeight(2) should show fewer activities than SetMaxHeight(20): "+
			"capped=%d full=%d\ncapped:\n%q\nfull:\n%q", cappedVisible, fullVisible, capped, full)
	}
}
