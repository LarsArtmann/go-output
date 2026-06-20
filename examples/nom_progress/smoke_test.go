package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/larsartmann/go-output/nom"
)

// TestNomProgressExampleSmokeTest guards against the v0.15.0 regression class:
// rendering with snapshots must produce non-blank, labelled output. The deleted
// footgun wrappers (RenderString etc.) silently rendered blanks on nil
// snapshots — this test fails loudly if the snapshot path ever regresses.
func TestNomProgressExampleSmokeTest(t *testing.T) {
	subscriber := nom.NewNOMStyleSubscriber()
	ctx := context.Background()

	activities := []struct {
		id, name string
		status   nom.ActivityStatus
		delay    time.Duration
	}{
		{"fetch", "Fetch Dependencies", nom.ActivityStatusCompleted, 10 * time.Millisecond},
		{"compile", "Compile Sources", nom.ActivityStatusCompleted, 20 * time.Millisecond},
		{"test", "Run Tests", nom.ActivityStatusRunning, 30 * time.Millisecond},
		{"lint", "Lint Code", nom.ActivityStatusFailed, 0},
		{"package", "Package Binary", nom.ActivityStatusPending, 0},
	}

	for _, a := range activities {
		if err := subscriber.OnEvent(ctx, nom.ActivityStarted{
			ID:   nom.NewActivityID(a.id),
			Name: nom.NewActivityName(a.name),
		}); err != nil {
			t.Fatalf("OnEvent(started %q): %v", a.id, err)
		}

		switch a.status {
		case nom.ActivityStatusCompleted:
			if err := subscriber.OnEvent(ctx, nom.ActivityCompleted{
				ID:       nom.NewActivityID(a.id),
				Name:     nom.NewActivityName(a.name),
				Duration: a.delay,
			}); err != nil {
				t.Fatalf("OnEvent(completed %q): %v", a.id, err)
			}
		case nom.ActivityStatusFailed:
			if err := subscriber.OnEvent(ctx, nom.ActivityFailed{
				ID:   nom.NewActivityID(a.id),
				Name: nom.NewActivityName(a.name),
			}); err != nil {
				t.Fatalf("OnEvent(failed %q): %v", a.id, err)
			}
		}
	}

	subscriber.UpdateRunningActivityElapsed()
	snaps := subscriber.SnapshotActivities()

	if len(snaps) == 0 {
		t.Fatal("SnapshotActivities() returned empty map — activities were not registered")
	}

	tree := subscriber.GetDependencyTree()

	// Full-height render: must contain every activity label.
	full := tree.RenderWithSnapshots(snaps, 20, 0)
	assertNonBlank(t, "full render", full)

	for _, a := range activities {
		if !strings.Contains(full, a.name) {
			t.Errorf("full render missing activity label %q\noutput:\n%s", a.name, full)
		}
	}

	// Height-limited render: must still be non-blank (failed/running surface first).
	limited := tree.RenderWithSnapshots(snaps, 3, 0)
	assertNonBlank(t, "limited render", limited)

	// Summary must be non-empty.
	summary := subscriber.GetActivityCounts().Summary()
	if strings.TrimSpace(summary) == "" {
		t.Error("ActivityCounts.Summary() is blank")
	}
}

func assertNonBlank(t *testing.T, label, output string) {
	t.Helper()

	stripped := strings.TrimSpace(output)
	if stripped == "" {
		t.Errorf("%s produced blank output", label)
	}
}
