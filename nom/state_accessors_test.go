package nom

import (
	"context"
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
)

func TestParallelismStats_String(t *testing.T) {
	t.Parallel()

	ps := ParallelismStats{Running: 2, Possible: 5}
	want := "parallel: 2/5 possible"

	if got := ps.String(); got != want {
		t.Errorf("ParallelismStats.String() = %q, want %q", got, want)
	}
}

func TestNOMSubscriber_ParallelismStats_Serial(t *testing.T) {
	t.Parallel()

	ns := NewNOMSubscriber()
	ctx := context.Background()

	ns.OnEvent(ctx, ActivityRegistered{ID: "a", Name: "a"})
	ns.OnEvent(ctx, ActivityRegistered{ID: "b", Name: "b"})
	ns.OnEvent(ctx, ActivityStarted{ID: "a", Name: "a"})

	stats := ns.ParallelismStats()

	testhelpers.AssertEqual(t, "Running", stats, stats.Running, 1)
	testhelpers.AssertEqual(t, "Possible", stats, stats.Possible, 1)
}

func TestNOMSubscriber_ParallelismStats_Diamond(t *testing.T) {
	t.Parallel()

	ns := NewNOMSubscriber()
	ctx := context.Background()

	// Diamond: a -> b, a -> c, b -> d, c -> d
	ns.OnEvent(ctx, ActivityRegistered{ID: "a", Name: "a"})
	ns.OnEvent(ctx, ActivityRegistered{ID: "b", Name: "b", Deps: []ActivityID{"a"}})
	ns.OnEvent(ctx, ActivityRegistered{ID: "c", Name: "c", Deps: []ActivityID{"a"}})
	ns.OnEvent(ctx, ActivityRegistered{ID: "d", Name: "d", Deps: []ActivityID{"b", "c"}})

	// All pending: only root 'a' is ready.
	stats := ns.ParallelismStats()
	testhelpers.AssertEqual(t, "Running", stats, stats.Running, 0)
	testhelpers.AssertEqual(t, "Possible", stats, stats.Possible, 1)

	// Start a and b.
	ns.OnEvent(ctx, ActivityStarted{ID: "a", Name: "a"})
	ns.OnEvent(ctx, ActivityStarted{ID: "b", Name: "b"})

	stats = ns.ParallelismStats()
	testhelpers.AssertEqual(t, "Running", stats, stats.Running, 2)
	testhelpers.AssertEqual(t, "Possible", stats, stats.Possible, 0)

	// Complete a and b; c becomes ready.
	ns.OnEvent(ctx, ActivityCompleted{ID: "a", Name: "a"})
	ns.OnEvent(ctx, ActivityCompleted{ID: "b", Name: "b"})

	stats = ns.ParallelismStats()
	testhelpers.AssertEqual(t, "Running", stats, stats.Running, 0)
	testhelpers.AssertEqual(t, "Possible", stats, stats.Possible, 1)
}

func TestNOMSubscriber_ParallelismStats_AllDone(t *testing.T) {
	t.Parallel()

	ns := NewNOMSubscriber()
	ctx := context.Background()

	ns.OnEvent(ctx, ActivityRegistered{ID: "a", Name: "a"})
	ns.OnEvent(ctx, ActivityStarted{ID: "a", Name: "a"})
	ns.OnEvent(ctx, ActivityCompleted{ID: "a", Name: "a"})

	stats := ns.ParallelismStats()
	testhelpers.AssertEqual(t, "Running", stats, stats.Running, 0)
	testhelpers.AssertEqual(t, "Possible", stats, stats.Possible, 0)
}
