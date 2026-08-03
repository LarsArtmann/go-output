package nom

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// TestInlineRenderer_PlainText_NoTreeReprintOnTimerTick verifies the core fix
// for the plainText (non-TTY/CI) tree repetition bug: when a step is running
// and has a live elapsed timer that changes every tick, the tree must NOT be
// re-appended on every Draw(). Only structural changes (step start/complete/fail)
// should trigger a tree reprint.
//
// Before the fix (stripTimings), the tree frame contained live timers like
// "12.3s" that changed every 100ms, making treeChanged always true.
func TestInlineRenderer_PlainText_NoTreeReprintOnTimerTick(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	// NewInlineRenderer with a bytes.Buffer → not a TTY → plainText mode.
	renderer := NewInlineRenderer(sub, &buf, 10)
	renderer.SetNoColor(true)

	if renderer.writerIsTTY {
		t.Fatal("buffer writer should not be detected as TTY")
	}

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-plaintext-tick"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Build"))
	renderer.SetStartTime(time.Now().Add(-5 * time.Second))

	// First Draw — emits the tree (first frame).
	renderer.Draw()

	firstOutput := buf.String()
	if firstOutput == "" {
		t.Fatal("first Draw should emit output")
	}

	treeCount := strings.Count(firstOutput, "Build")
	if treeCount < 1 {
		t.Fatalf("first Draw should contain the tree with 'Build' step, got:\n%s", firstOutput)
	}

	buf.Reset()

	// Second Draw — the elapsed timer changed but tree structure didn't.
	// stripTimings should make the stable comparison match → no tree reprint.
	renderer.Draw()

	secondOutput := buf.String()
	if secondOutput != "" {
		t.Errorf("second Draw with only timer change should emit zero bytes in plainText, got:\n%q", secondOutput)
	}

	// Third Draw — still only timer change. Still zero.
	renderer.Draw()

	if buf.String() != "" {
		t.Errorf("third Draw with only timer change should emit zero bytes in plainText, got:\n%q", buf.String())
	}
}

// TestInlineRenderer_PlainText_PendingLinesWithoutTreeReprint verifies that
// pending log lines are printed independently from the tree: EnqueueLines +
// Draw should output the log lines but NOT reprint the tree (since no
// structural change occurred).
func TestInlineRenderer_PlainText_PendingLinesWithoutTreeReprint(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)
	renderer.SetNoColor(true)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-plaintext-pending"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Build"))

	// First Draw — establishes the tree.
	renderer.Draw()
	buf.Reset()

	// Enqueue subprocess log lines and Draw.
	renderer.EnqueueLines([]string{"nix> evaluating flake...", "nix> building derivation..."})
	renderer.Draw()

	output := buf.String()

	// Pending lines should appear.
	if !strings.Contains(output, "nix> evaluating flake") {
		t.Errorf("pending log lines should appear in output, got:\n%q", output)
	}

	if !strings.Contains(output, "nix> building derivation") {
		t.Errorf("second pending line should appear in output, got:\n%q", output)
	}

	// The tree should NOT be reprinted (no structural change).
	if strings.Contains(output, "Build") {
		t.Errorf("tree should not be reprinted when only pending lines arrived, got:\n%q", output)
	}
}

// TestInlineRenderer_PlainText_StructuralChangeReprintsTree verifies that
// despite the stripTimings optimization, genuine structural changes still
// trigger a tree reprint in plainText mode.
func TestInlineRenderer_PlainText_StructuralChangeReprintsTree(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 10)
	renderer.SetNoColor(true)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-plaintext-structural"), "")
	sendActivityStarted(t, sub, ctx, ActivityID("step1"), ActivityName("Build"))

	// First Draw — tree appears.
	renderer.Draw()
	buf.Reset()

	// Complete the step — structural change.
	sendActivityCompleted(t, sub, ctx, ActivityID("step1"), ActivityName("Build"), 3*time.Second)

	renderer.Draw()

	output := buf.String()
	if output == "" {
		t.Fatal("Draw after structural change should emit output")
	}

	if !strings.Contains(output, "Build") {
		t.Errorf("Draw after structural change should contain the tree, got:\n%q", output)
	}
}

// TestStripTimings verifies that stripTimings replaces timing values with
// placeholders while leaving non-timing content (step names, symbols) intact.
func TestStripTimings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "seconds with decimal",
			input: "⏵ nix-build 12.3s",
			want:  "⏵ nix-build T",
		},
		{
			name:  "milliseconds",
			input: "✔ buf-format 456ms",
			want:  "✔ buf-format T",
		},
		{
			name:  "no timing value",
			input: "✔ buf-format",
			want:  "✔ buf-format",
		},
		{
			name:  "multiple timings",
			input: "⏵ step1 1.5s │ step2 2.3s",
			want:  "⏵ step1 T │ step2 T",
		},
		{
			name:  "estimate prefix",
			input: "~2.0s left",
			want:  "~T left",
		},
		{
			name:  "whole seconds",
			input: "✔ done 5s",
			want:  "✔ done T",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := stripTimings(tt.input)
			if got != tt.want {
				t.Errorf("stripTimings(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
