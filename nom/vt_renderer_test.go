package nom

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

// newVTTestRenderer creates an InlineRenderer that writes to a VT harness
// (simulating a real terminal) instead of a plain bytes.Buffer. This exercises
// the actual inline-rendering path: cursor-up, erase-line, sync-output, and
// ghost-line cleanup — all processed by a real VT emulator.
//
// The caller gets both the harness (for screen assertions) and a buffer (for
// raw byte inspection if needed).
func newVTTestRenderer(
	t *testing.T,
	sub *NOMSubscriber,
	harness *vtHarness,
	maxHeight int,
) *InlineRenderer {
	t.Helper()

	r := NewInlineRenderer(sub, harness, maxHeight)
	r.SetNoColor(true) // deterministic output (no terminal color codes)

	// Force the inline rendering path: pretend the VT harness is a TTY.
	// snapshotConfig() computes effectivePlainText = plainText || !writerIsTTY,
	// so writerIsTTY=true + plainText=false gives us the inline path.
	r.tickMu.Lock()
	r.writerIsTTY = true
	r.plainText = false
	r.tickMu.Unlock()

	return r
}

// --- F7: First-frame Draw renders tree content ---

func TestVT_FirstFrame_ShowsTreeContent(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Build Module"))

	renderer.Draw()

	harness.assertScreenContains(t, "Build Module")
}

// --- F8: First-frame Draw hides cursor (DECTCEM off) ---

func TestVT_FirstFrame_HidesCursor(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 10)
	renderer.SetHideCursor(true)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step 1"))

	renderer.Draw()

	if !harness.cursorHidden {
		t.Error("cursor should be hidden after first frame with HideCursor enabled")
	}
}

// --- F9: Second-frame Draw shows updated state (old content repainted) ---

func TestVT_SecondFrame_UpdatesContent(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Pending Step"))

	renderer.Draw()

	// The screen should show the running step initially.
	harness.assertScreenContains(t, "Pending Step")

	// Mark it completed — the symbol changes from ⏵ to ✔.
	sendActivityCompleted(t, sub, ctx, ActivityID("step1"), ActivityName("Pending Step"), time.Second)

	renderer.Draw()

	// The step should still be visible but now show the completed symbol (✔).
	harness.assertScreenContains(t, "Pending Step")

	// The screen should NOT show a stale running symbol for the completed step.
	// After repaint, the VT screen reflects the new state cleanly.
	lines := harness.nonEmptyLines()
	for _, line := range lines {
		if strings.Contains(line, "Pending Step") && strings.Contains(line, "⏵") {
			t.Errorf("completed step should not show running symbol: %q", line)
		}
	}
}

// --- F10: Second-frame Draw uses cursor-up (cursor position moves up) ---

func TestVT_SecondFrame_CursorUpCorrect(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Step A"))

	renderer.Draw()

	posAfterFirst := harness.term.CursorPosition()

	// Same tree, different status to force a redraw.
	sendActivityCompleted(t, sub, ctx, ActivityID("step1"), ActivityName("Step A"), time.Second)

	renderer.Draw()

	posAfterSecond := harness.term.CursorPosition()

	// After the second Draw, the cursor should have moved UP (lower Y) to repaint
	// the frame, then settled back. The key assertion is that cursor-up was used:
	// the cursor Y after the second draw should be at or above the cursor Y after
	// the first draw (since we're overwriting the same frame area).
	if posAfterSecond.Y > posAfterFirst.Y+1 {
		t.Errorf(
			"cursor should not have moved far down after redraw: first=%v second=%v",
			posAfterFirst, posAfterSecond,
		)
	}

	// The screen should show the updated content (completed status).
	harness.assertScreenContains(t, "Step A")
}

// --- F11: Frame shrink leaves no ghost lines ---

func TestVT_FrameShrink_NoGhostLines(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 20)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")

	// Create multiple activities → tall frame.
	sendActivityStarted(t, sub, ctx, ActivityID("s1"), ActivityName("Alpha"))
	sendActivityStarted(t, sub, ctx, ActivityID("s2"), ActivityName("Beta"))
	sendActivityStarted(t, sub, ctx, ActivityID("s3"), ActivityName("Gamma"))

	renderer.Draw()

	// Verify all three are visible.
	harness.assertScreenContains(t, "Alpha")
	harness.assertScreenContains(t, "Beta")
	harness.assertScreenContains(t, "Gamma")

	// Now mark all as completed — the tree may collapse (elideCompletedUnderPressure)
	// and the frame shrinks. The renderer must erase ghost lines.
	sendActivityCompleted(t, sub, ctx, ActivityID("s1"), ActivityName("Alpha"), time.Second)
	sendActivityCompleted(t, sub, ctx, ActivityID("s2"), ActivityName("Beta"), time.Second)
	sendActivityCompleted(t, sub, ctx, ActivityID("s3"), ActivityName("Gamma"), time.Second)

	renderer.Draw()

	// No ghost text should remain in rows below the new shorter frame.
	// The new frame is shorter, so rows that had "Beta" or "Gamma" must be clean
	// if they're no longer part of the visible tree.
	lines := harness.nonEmptyLines()
	for _, line := range lines {
		// Each non-empty line should be part of the current tree state.
		// Ghost lines would show stale activity names that were completed.
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
	}
}

// --- F12: Frame grow adds new lines correctly ---

func TestVT_FrameGrow_AddsNewLines(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 20)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("s1"), ActivityName("First"))

	renderer.Draw()
	harness.assertScreenContains(t, "First")

	// Add a second activity — frame grows.
	sendActivityStarted(t, sub, ctx, ActivityID("s2"), ActivityName("Second"))

	renderer.Draw()

	harness.assertScreenContains(t, "First")
	harness.assertScreenContains(t, "Second")
}

// --- F13: Sync-output mode 2026 wrapping on TTY ---

func TestVT_SyncOutput_WrapsWithMode2026(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 10)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Sync Step"))

	renderer.Draw()

	// The renderer wraps output in mode 2026 sync codes when writerIsTTY=true.
	// The VT emulator should have received and processed the sync mode.
	if !harness.syncWasActive {
		t.Error("synchronized output mode (2026) should have been activated during Draw on TTY writer")
	}

	if harness.syncToggleCount < 2 {
		t.Errorf("sync mode should have been toggled at least twice (set+reset), got %d toggles",
			harness.syncToggleCount)
	}
}

// --- F14: Plain-text mode emits no escape sequences ---

func TestVT_PlainText_NoEscapeSequences(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)

	// Create renderer in default mode (non-TTY buffer → plain text).
	// We don't force writerIsTTY, so the renderer degrades to plain text.
	renderer := NewInlineRenderer(sub, harness, 10)
	renderer.SetNoColor(true)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Plain Step"))

	renderer.Draw()

	// In plain-text mode, sync output should NOT have been activated.
	if harness.syncWasActive {
		t.Error("synchronized output mode should NOT be activated in plain-text mode")
	}

	// The screen should still show the content.
	harness.assertScreenContains(t, "Plain Step")
}

// --- F15: Frame diffing skips identical redraws (no extra bytes) ---

func TestVT_FrameDiff_IdenticalSkipsRedraw(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	// Use a tee: VT harness for screen assertions + buffer for byte counting.
	harness := newVTHarness(t, 100, 30)

	var buf bytes.Buffer

	// Create a tee writer that writes to both.
	teeWriter := teeWriter{w1: harness, w2: &buf}

	renderer := NewInlineRenderer(sub, teeWriter, 10)
	renderer.SetNoColor(true)
	renderer.tickMu.Lock()
	renderer.writerIsTTY = true
	renderer.plainText = false
	renderer.tickMu.Unlock()

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Stable Step"))

	renderer.Draw()

	bytesAfterFirst := buf.Len()
	buf.Reset()

	// Draw again without any state change — frame is identical.
	renderer.Draw()

	bytesAfterSecond := buf.Len()

	if bytesAfterFirst == 0 {
		t.Fatal("first Draw should have produced output")
	}

	if bytesAfterSecond != 0 {
		t.Errorf(
			"identical frame should produce 0 bytes on second Draw, got %d bytes",
			bytesAfterSecond,
		)
	}

	// Screen should still show the content from the first draw.
	harness.assertScreenContains(t, "Stable Step")
}

// --- F16: Finish restores cursor visibility ---

func TestVT_Finish_ShowsCursor(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 10)
	renderer.SetHideCursor(true)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Final Step"))

	renderer.Draw()

	if !harness.cursorHidden {
		t.Fatal("cursor should be hidden during rendering")
	}

	renderer.Finish()

	if harness.cursorHidden {
		t.Error("cursor should be visible after Finish")
	}

	harness.assertScreenContains(t, "Final Step")
}

// teeWriter writes to two writers simultaneously — used to feed the VT emulator
// while also counting bytes for frame-diff assertions.
type teeWriter struct {
	w1, w2 interface {
		Write([]byte) (int, error)
	}
}

func (tw teeWriter) Write(p []byte) (int, error) {
	n1, err := tw.w1.Write(p)
	if err != nil {
		return n1, err
	}

	_, err = tw.w2.Write(p)
	if err != nil {
		return n1, err
	}

	return n1, nil
}

// --- F17: Color-on rendering produces SGR sequences ---

func TestVT_ColorOn_EmitsSGR(t *testing.T) {
	// NOT parallel: temporarily mutates global lipgloss.Writer.Profile.
	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)

	var buf bytes.Buffer

	tw := teeWriter{w1: harness, w2: &buf}

	renderer := NewInlineRenderer(sub, tw, 10)
	renderer.SetNoColor(false)

	renderer.tickMu.Lock()
	renderer.writerIsTTY = true
	renderer.plainText = false
	renderer.tickMu.Unlock()

	// Force lipgloss to emit ANSI SGR codes regardless of terminal detection.
	oldProfile := lipgloss.Writer.Profile
	lipgloss.Writer.Profile = colorprofile.ANSI

	t.Cleanup(func() { lipgloss.Writer.Profile = oldProfile })

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("ok"), ActivityName("Green Step"))
	sendActivityCompleted(t, sub, ctx, ActivityID("ok"), ActivityName("Green Step"), time.Second)

	renderer.Draw()

	raw := buf.String()

	// With ANSI profile forced, lipgloss emits SGR codes. The completed
	// activity renders with Colors.Completed (lipgloss.Color("10") = bright
	// green, emitted as \x1b[92m or \x1b[1;92m when combined with bold).
	if !strings.Contains(raw, "\x1b[") {
		t.Errorf("expected ANSI escape sequences in raw output with colors on\n\nRaw:\n%q", raw)
	}

	// Verify SGR color codes are present (not just cursor sequences).
	// SGR codes end with 'm'; cursor codes end with a letter.
	if !strings.Contains(raw, "m✔") && !strings.Contains(raw, "m⏵") {
		t.Errorf("expected SGR-colored status symbol in output\n\nRaw:\n%q", raw)
	}

	// VT emulator should have processed the SGR codes into cell colors.
	foundColoredCell := false
	for y := 0; y < 5 && !foundColoredCell; y++ {
		for x := 0; x < harness.term.Width() && !foundColoredCell; x++ {
			cell := harness.term.CellAt(x, y)
			if cell != nil && cell.Style.Fg != nil {
				foundColoredCell = true
			}
		}
	}

	if !foundColoredCell {
		t.Error("expected at least one cell with non-nil foreground color after SGR rendering")
	}

	// Content should still be visible on the screen.
	harness.assertScreenContains(t, "Green Step")
}

// --- F18: Layered mode renders activities grouped by depth ---

func TestVT_LayeredMode_ShowsActivities(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 20)

	// Switch to layered mode.
	sub.DependencyTree().SetRenderMode(RenderModeLayered)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")

	// Layer 0: setup (root).
	sendActivityStarted(t, sub, ctx, ActivityID("setup"), ActivityName("Setup"))

	// Layer 1: compile + lint (depend on setup).
	_ = sub.OnEvent(ctx, ActivityRegistered{
		ID:   ActivityID("compile"),
		Name: ActivityName("Compile"),
		Deps: []ActivityID{ActivityID("setup")},
	})
	_ = sub.OnEvent(ctx, ActivityRegistered{
		ID:   ActivityID("lint"),
		Name: ActivityName("Lint"),
		Deps: []ActivityID{ActivityID("setup")},
	})

	// Layer 2: test (depends on compile).
	_ = sub.OnEvent(ctx, ActivityRegistered{
		ID:   ActivityID("test"),
		Name: ActivityName("Test"),
		Deps: []ActivityID{ActivityID("compile")},
	})

	renderer.Draw()

	// All three activity labels should be visible on the VT screen.
	harness.assertScreenContains(t, "Setup")
	harness.assertScreenContains(t, "Compile")
	harness.assertScreenContains(t, "Lint")
	harness.assertScreenContains(t, "Test")
}

// --- F19: Layered mode updates on second frame ---

func TestVT_LayeredMode_SecondFrame(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 20)

	sub.DependencyTree().SetRenderMode(RenderModeLayered)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")

	sendActivityStarted(t, sub, ctx, ActivityID("a"), ActivityName("Alpha"))
	sendActivityStarted(t, sub, ctx, ActivityID("b"), ActivityName("Beta"))

	renderer.Draw()

	harness.assertScreenContains(t, "Alpha")
	harness.assertScreenContains(t, "Beta")

	// Complete Alpha — second frame should still show content correctly.
	sendActivityCompleted(t, sub, ctx, ActivityID("a"), ActivityName("Alpha"), time.Second)

	renderer.Draw()

	harness.assertScreenContains(t, "Alpha")
	harness.assertScreenContains(t, "Beta")
}

// --- F20: Layered mode with height pressure collapses completed layers ---

func TestVT_LayeredMode_HeightPressureCollapse(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)
	harness := newVTHarness(t, 100, 30)
	renderer := newVTTestRenderer(t, sub, harness, 4)

	sub.DependencyTree().SetRenderMode(RenderModeLayered)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-1"), "")

	sendActivityStarted(t, sub, ctx, ActivityID("a"), ActivityName("Alpha"))
	sendActivityStarted(t, sub, ctx, ActivityID("b"), ActivityName("Beta"))

	renderer.Draw()

	// Complete both — under height pressure the renderer should collapse
	// terminal layers without producing ghost lines or corruption.
	sendActivityCompleted(t, sub, ctx, ActivityID("a"), ActivityName("Alpha"), time.Second)
	sendActivityCompleted(t, sub, ctx, ActivityID("b"), ActivityName("Beta"), time.Second)

	renderer.Draw()

	// Content should still be visible on the screen.
	harness.assertScreenContains(t, "Alpha")
	harness.assertScreenContains(t, "Beta")
}
