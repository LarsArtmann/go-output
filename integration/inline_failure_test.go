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

	subscriber := nom.NewNOMStyleSubscriber()
	ctx := context.Background()

	fireWorkflowStarted(subscriber, ctx, "wf-fail", "Failing Pipeline")

	startActivity(subscriber, ctx, "fetch", "Fetch Dependencies")
	completeActivity(subscriber, ctx, "fetch", "Fetch Dependencies", 50*time.Millisecond)

	startActivity(subscriber, ctx, "compile", "Compile Sources")
	completeActivity(subscriber, ctx, "compile", "Compile Sources", 120*time.Millisecond)

	startActivity(subscriber, ctx, "test", "Run Tests")

	// The test step fails with a concrete error.
	stepErr := errors.New("test suite failed: 3 assertions failed")
	_ = subscriber.OnEvent(ctx, &nomTestEvent{
		eventType: nom.EventActivityFailed,
		aID:       nom.NewActivityID("test"),
		aName:     nom.NewActivityName("Run Tests"),
		err:       stepErr,
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

	// Final render with the propagated workflow error.
	workflowErr := errors.New("step test failed: test suite failed: 3 assertions failed")
	renderer.Finish(workflowErr)

	final := buf.String()
	if final == "" {
		t.Fatal("Finish should produce final output")
	}

	// The final status line must report failure and carry the error text.
	if !strings.Contains(final, "failed") {
		t.Errorf("final output should report failure\ngot:\n%q", final)
	}

	if !strings.Contains(final, "3 assertions failed") {
		t.Errorf("final output should include the error message\ngot:\n%q", final)
	}

	if !strings.Contains(final, "BuildFlow") {
		t.Errorf("final output should include the app name\ngot:\n%q", final)
	}

	// Counts should reflect one completed set + one failed.
	counts := subscriber.GetActivityCounts()
	if counts.Failed != 1 {
		t.Errorf("expected 1 failed activity, got %d", counts.Failed)
	}
}
