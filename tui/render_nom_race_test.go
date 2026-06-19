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
// NOMStyleSubscriber.WithTreeRLock, this trip -race and produced garbled frames.
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

	_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
		eventType: nom.EventWorkflowStarted,
		wID:       nom.WorkflowID("wf-tui-race"),
		wName:     nom.WorkflowName("Race"),
	})

	const groupID = nom.ActivityID("group")

	ids := make([]nom.ActivityID, 10)

	for i := range ids {
		ids[i] = nom.ActivityID(fmt.Sprintf("act%d", i))
		_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
			eventType:    nom.EventActivityRegistered,
			aID:          ids[i],
			aName:        nom.ActivityName(fmt.Sprintf("Act %d", i)),
			dependencies: []nom.ActivityID{groupID},
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
					_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
						eventType: nom.EventActivityStarted,
						aID:       id,
						aName:     nom.ActivityName("Act"),
					})
				case 1:
					_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
						eventType: nom.EventActivityCompleted,
						aID:       id,
						aName:     nom.ActivityName("Act"),
						duration:  3 * time.Millisecond,
					})
				case 2:
					_ = model.nomSubscriber.OnEvent(ctx, &testEvent{
						eventType: nom.EventActivityFailed,
						aID:       id,
						aName:     nom.ActivityName("Act"),
						err:       errors.New("boom"),
					})
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
