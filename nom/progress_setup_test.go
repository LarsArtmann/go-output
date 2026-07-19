package nom

import (
	"context"
	"testing"
)

// setupProgress returns a NOMSubscriber with one activity already running in a
// parent workflow, and the caller-supplied progress message applied. It is
// the canonical setup for tests asserting on ActivityProgress behaviour,
// sharing IDs/name with the rest of the suite. Tests with different activity
// IDs/names pass their own via setupProgressWith.
func setupProgress(t *testing.T, message string) (*NOMSubscriber, context.Context) {
	return setupProgressWith(t, "step1", "go-mod-tidy", message)
}

func setupProgressWith(t *testing.T, id, name, message string) (*NOMSubscriber, context.Context) {
	t.Helper()

	ns := NewNOMSubscriber()
	ctx := context.Background()

	emit(ctx, ns, WorkflowStarted{ID: WorkflowID("wf"), Name: WorkflowName("test")}, t)
	emit(ctx, ns, ActivityStarted{ID: ActivityID(id), Name: ActivityName(name)}, t)
	emit(ctx, ns, ActivityProgress{ID: ActivityID(id), Name: ActivityName(name), Message: message}, t)

	return ns, ctx
}

// emit dispatches ev to ns, failing t if the subscriber rejects it. Wraps the
// repeated _ = ns.OnEvent(ctx, ...) boilerplate that previously dominated
// every progress test's setup.
func emit(ctx context.Context, ns *NOMSubscriber, ev Event, t *testing.T) {
	t.Helper()

	if err := ns.OnEvent(ctx, ev); err != nil {
		t.Fatalf("setup event %T: %v", ev, err)
	}
}
