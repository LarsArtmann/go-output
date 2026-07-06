package tui

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/larsartmann/go-output/nom"
)

// TestProgressModel_RenderRacingActivityMutation reproduces the TUI twin of the
// inline-renderer race: the BubbleTea View() (render goroutine) walks the
// dependency tree reading the shared *Activity fields, while external event
// dispatchers mutate those same fields via OnEvent (SetRunning/SetCompleted/
// SetFailed) on step goroutines — exactly how BuildFlow's ProgressBridge drives
// the subscriber. Before routing renderDependencyTree through
// NOMSubscriber.WithTreeRLock, this trip -race and produced garbled frames.
//
// The model fields (width/height/displayMode) are set before spawning, and only
// the renderer goroutine calls View(), so the only concurrent access under test
// is the subscriber's *Activity state — isolating the bug.
func TestProgressModel_RenderRacingActivityMutation(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.displayMode = DisplayModeNOM
	model.width = 80
	model.height = 40

	ctx := context.Background()

	_ = model.nomSubscriber.OnEvent(ctx, nom.WorkflowStarted{
		ID:   nom.WorkflowID("wf-tui-race"),
		Name: nom.WorkflowName("Race"),
	})

	const groupID = nom.ActivityID("group")

	ids := make([]nom.ActivityID, 0, 10)

	for i := range 10 {
		ids = append(ids, nom.ActivityID(fmt.Sprintf("act%d", i)))
		_ = model.nomSubscriber.OnEvent(ctx, nom.ActivityRegistered{
			ID:   ids[i],
			Name: nom.ActivityName(fmt.Sprintf("Act %d", i)),
			Deps: []nom.ActivityID{groupID},
		})
	}

	// Sync once so the model caches the tree pointer (mirrors a real tick).
	if _, _ = model.Update(tickMsg(time.Now())); false {
		t.Fatal("unreachable")
	}

	var wg sync.WaitGroup

	// Mutators: hammer lifecycle events from many goroutines (like BuildFlow's
	// step hooks running on the workflow engine's worker pool).
	for _, id := range ids {
		wg.Add(1)

		go func(id nom.ActivityID) {
			defer wg.Done()

			for j := range 15 {
				switch j % 3 {
				case 0:
					_ = model.nomSubscriber.OnEvent(ctx, nom.ActivityStarted{ID: id, Name: nom.ActivityName("Act")})
				case 1:
					_ = model.nomSubscriber.OnEvent(
						ctx,
						nom.ActivityCompleted{ID: id, Name: nom.ActivityName("Act"), Duration: 3 * time.Millisecond},
					)
				case 2:
					_ = model.nomSubscriber.OnEvent(
						ctx,
						nom.ActivityFailed{ID: id, Name: nom.ActivityName("Act"), Err: errors.New("boom")},
					)
				}
			}
		}(id)
	}

	// Renderer: the View() runs on a single goroutine (as bubbletea guarantees).
	wg.Go(func() {
		for range 80 {
			_ = model.View()
		}
	})

	wg.Wait()
}
