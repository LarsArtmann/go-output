package nom

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// TestInlineRenderer_FullWorkflowLifecycle is an end-to-end integration test
// exercising the entire pipeline: workflow start -> activity registration with
// dependencies -> start/complete/fail transitions -> inline renders between
// phases -> final Finish. It guards against interaction bugs between the event
// subscriber, the dependency tree, and the inline renderer that unit tests on
// isolated pieces can miss.
func TestInlineRenderer_FullWorkflowLifecycle(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	renderer := NewInlineRenderer(sub, &buf, 20)
	renderer.SetNoColor(true)

	ctx := context.Background()

	// 1. Start the workflow.
	if err := sendWorkflowStarted(sub, ctx, WorkflowID("wf-int"), "Release Pipeline"); err != nil {
		t.Fatalf("workflow.started: %v", err)
	}

	// 2. Register a phase with three dependent activities.
	registerPhase(sub, ctx, ActivityID("build"), ActivityName("Build"))
	registerActivity(sub, ctx, ActivityID("compile"), ActivityName("Compile"), "build")
	registerActivity(sub, ctx, ActivityID("test"), ActivityName("Test"), "build")
	registerActivity(sub, ctx, ActivityID("deploy"), ActivityName("Deploy"), "build")

	renderer.Draw()

	firstFrame := buf.String()
	if !strings.Contains(firstFrame, "Build") {
		t.Errorf("first frame missing phase label; got:\n%s", firstFrame)
	}

	// 3. Run and complete compile + test, then fail deploy.
	buf.Reset()
	sendActivityStarted(t, sub, ctx, ActivityID("compile"), ActivityName("Compile"))
	sendActivityCompleted(t, sub, ctx, ActivityID("compile"), ActivityName("Compile"), 2*time.Second)
	sendActivityStarted(t, sub, ctx, ActivityID("test"), ActivityName("Test"))
	sendActivityCompleted(t, sub, ctx, ActivityID("test"), ActivityName("Test"), 5*time.Second)
	sendActivityStarted(t, sub, ctx, ActivityID("deploy"), ActivityName("Deploy"))

	renderer.Draw()

	midFrame := buf.String()
	for _, want := range []string{"Compile", "Test", "Deploy"} {
		if !strings.Contains(midFrame, want) {
			t.Errorf("mid frame missing %q; got:\n%s", want, midFrame)
		}
	}

	// 4. Finalize: fail deploy, finish the renderer with an error.
	if err := sub.OnEvent(ctx, ActivityFailed{
		ID:   ActivityID("deploy"),
		Name: ActivityName("Deploy"),
		Err:  errDeployFailed,
	}); err != nil {
		t.Fatalf("activity.failed: %v", err)
	}

	buf.Reset()
	renderer.Finish()

	final := buf.String()
	if final == "" {
		t.Fatal("Finish produced no output")
	}
	// The final static frame must reflect the terminal state of every activity.
	for _, want := range []string{"Compile", "Test", "Deploy"} {
		if !strings.Contains(final, want) {
			t.Errorf("final frame missing %q; got:\n%s", want, final)
		}
	}

	// Counts must be coherent: 2 completed, 1 failed, 0 running.
	counts := sub.GetActivityCounts()
	if counts.Completed != 2 {
		t.Errorf("completed = %d, want 2", counts.Completed)
	}

	if counts.Failed != 1 {
		t.Errorf("failed = %d, want 1", counts.Failed)
	}

	if counts.Running != 0 {
		t.Errorf("running = %d, want 0 after Finish", counts.Running)
	}
}

var errDeployFailed = strError("deploy failed")

type strError string

func (e strError) Error() string { return string(e) }
