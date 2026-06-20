package nom

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
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

	renderer := NewInlineRenderer(sub, &buf, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.Draw()
	buf.Reset()

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

	renderer := NewInlineRenderer(sub, &buf, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.Start(ctx, 50*time.Millisecond)
	time.Sleep(180 * time.Millisecond)
	renderer.Stop()

	output := buf.String()

	renders := strings.Count(output, "Step 1")
	if renders < 2 {
		t.Errorf("expected at least 2 periodic renders, got %d renders", renders)
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

	renderer := NewInlineRenderer(sub, &buf, 10)
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
	// detectPlainText() is evaluated at construction, so CI must be set first.
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
	buf.Reset()
	renderer.Draw() // second draw — in plain mode, no cursor-up/sync codes

	output := buf.String()

	if strings.Contains(output, ansiCursorUp) {
		t.Errorf("CI plain mode must not emit cursor-up codes, got:\n%q", output)
	}

	if strings.Contains(output, ansiSyncBegin) || strings.Contains(output, ansiSyncEnd) {
		t.Errorf("CI plain mode must not emit sync-region codes, got:\n%q", output)
	}

	// Plain mode still renders the activity label.
	if !strings.Contains(output, "CI Step") {
		t.Errorf("CI plain mode should still render the activity label, got:\n%q", output)
	}
}
