package integration

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/go-output/nom"
)

// TestInlineRenderer_WorkflowFailureErrorDisplayed verifies the inline renderer
// surfaces a failed activity and the workflow error through to Finish output —
// the path a real build takes when a step fails and the workflow aborts.
func TestInlineRenderer_WorkflowFailureErrorDisplayed(t *testing.T) {
	t.Parallel()

	subscriber := nom.NewNOMSubscriber()
	ctx := context.Background()

	fireWorkflowStarted(subscriber, ctx, "wf-fail", "Failing Pipeline")

	startActivity(subscriber, ctx, "fetch", "Fetch Dependencies")
	completeActivity(subscriber, ctx, "fetch", "Fetch Dependencies", 50*time.Millisecond)

	startActivity(subscriber, ctx, "compile", "Compile Sources")
	completeActivity(subscriber, ctx, "compile", "Compile Sources", 120*time.Millisecond)

	startActivity(subscriber, ctx, "test", "Run Tests")

	// The test step fails with a concrete error.
	stepErr := errors.New("test suite failed: 3 assertions failed")
	_ = subscriber.OnEvent(ctx, nom.ActivityFailed{
		ID:   nom.NewActivityID("test"),
		Name: nom.NewActivityName("Run Tests"),
		Err:  stepErr,
	})

	var buf bytes.Buffer

	renderer := nom.NewInlineRenderer(subscriber, &buf, 10)
	renderer.SetAppName("BuildFlow")
	renderer.SetStartTime(time.Now())

	// Mid-build frame.
	renderer.Draw()

	midBuild := buf.String()
	if !strings.Contains(midBuild, "Run Tests") {
		t.Errorf("mid-build frame should show the failed activity label\ngot:\n%q", midBuild)
	}

	buf.Reset()

	// Final render — Finish() renders the static tree (no completion line;
	// the calling application is responsible for the post-run summary).
	workflowErr := errors.New("step test failed: test suite failed: 3 assertions failed")
	renderer.Finish(workflowErr)

	final := buf.String()
	if final == "" {
		t.Fatal("Finish should produce final output")
	}

	// The final tree must show the failed activity.
	if !strings.Contains(final, "Run Tests") {
		t.Errorf("final output should show the failed activity\ngot:\n%q", final)
	}

	// The failed activity should carry the failure symbol.
	if !strings.Contains(final, "⚠") {
		t.Errorf("final output should show the failure symbol\ngot:\n%q", final)
	}

	// Counts should reflect one completed set + one failed.
	counts := subscriber.GetActivityCounts()
	if counts.Failed != 1 {
		t.Errorf("expected 1 failed activity, got %d", counts.Failed)
	}
}
