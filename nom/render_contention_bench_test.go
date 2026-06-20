package nom

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

// benchEvent is a minimal event for populating/mutating a subscriber in benchmarks.
type benchEvent struct {
	eventType string
	wid       WorkflowID
	aid       ActivityID
	aname     ActivityName
}

func (e *benchEvent) GetEventType() string          { return e.eventType }
func (e *benchEvent) GetWorkflowID() WorkflowID     { return e.wid }
func (e *benchEvent) GetWorkflowName() WorkflowName { return WorkflowName("bench") }
func (e *benchEvent) GetActivityID() ActivityID     { return e.aid }
func (e *benchEvent) GetActivityName() ActivityName { return e.aname }
func (e *benchEvent) GetDependencies() []ActivityID { return nil }
func (e *benchEvent) GetDuration() time.Duration    { return 0 }
func (e *benchEvent) GetError() error               { return nil }

// buildBenchSubscriber creates a NOMStyleSubscriber with count activities spread
// across completed/running/pending states, mirroring a realistic mid-build tree.
func buildBenchSubscriber(b *testing.B, count int) *NOMStyleSubscriber {
	b.Helper()

	sub := NewNOMStyleSubscriber()
	ctx := context.Background()
	wid := NewWorkflowID("bench-wf")

	_ = sub.OnEvent(ctx, &benchEvent{eventType: EventWorkflowStarted, wid: wid})

	for i := range count {
		id := NewActivityID(fmt.Sprintf("step-%04d", i))
		name := NewActivityName(fmt.Sprintf("Step %d", i))

		_ = sub.OnEvent(ctx, &benchEvent{
			eventType: EventActivityRegistered, wid: wid, aid: id, aname: name,
		})
		_ = sub.OnEvent(ctx, &benchEvent{
			eventType: EventActivityStarted, wid: wid, aid: id, aname: name,
		})

		var final string

		switch {
		case i < count/3:
			final = EventActivityCompleted
		case i < count*2/3:
			final = "" // leave running
		default:
			final = "" // leave running/pending
		}

		if final != "" {
			_ = sub.OnEvent(ctx, &benchEvent{eventType: final, wid: wid, aid: id, aname: name})
		}
	}

	return sub
}

// BenchmarkRenderUnderStepChurn (#8) measures snapshot+render throughput under
// parallel contention — many goroutines snapshotting and walking the tree at
// once, as happens when a high-churn build fires events concurrently with the
// render loop. The snapshot path acquires the subscriber RLock; this surfaces
// any lock-contention regression.
func BenchmarkRenderUnderStepChurn(b *testing.B) {
	for _, size := range []int{100, 300} {
		b.Run(fmt.Sprintf("%dNodes", size), func(b *testing.B) {
			sub := buildBenchSubscriber(b, size)
			tree := sub.GetDependencyTree()
			maxH := size / 4

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					snaps := sub.SnapshotActivities()
					if tree.RenderWithSnapshots(snaps, maxH, 80) == "" {
						b.Fatal("render produced no output")
					}
				}
			})
		})
	}
}

// BenchmarkSnapshotActivities_Parallel isolates the lock-acquiring snapshot
// cost (the only serialized part of the render path) under contention.
func BenchmarkSnapshotActivities_Parallel(b *testing.B) {
	sub := buildBenchSubscriber(b, 300)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if len(sub.SnapshotActivities()) == 0 {
				b.Fatal("snapshot empty")
			}
		}
	})
}

// BenchmarkInlineRenderer_Draw (#19) measures end-to-end frame rendering:
// snapshot + tree walk + summary + ANSI assembly + write, the hot path for
// each redraw tick.
func BenchmarkInlineRenderer_Draw(b *testing.B) {
	for _, size := range []int{50, 200} {
		b.Run(fmt.Sprintf("%dNodes", size), func(b *testing.B) {
			sub := buildBenchSubscriber(b, size)

			var buf bytes.Buffer

			r := NewInlineRenderer(sub, &buf, size/3)
			r.SetAppName("Bench")
			r.SetStartTime(time.Now())

			b.ResetTimer()

			for b.Loop() {
				buf.Reset()
				r.Draw()
			}
		})
	}
}

// BenchmarkInlineRenderer_DrawWithChurn measures Draw while a background
// goroutine continuously mutates activity state, exercising the render-vs-
// mutation snapshot isolation under realistic concurrent load.
func BenchmarkInlineRenderer_DrawWithChurn(b *testing.B) {
	const size = 200

	sub := buildBenchSubscriber(b, size)

	var buf bytes.Buffer

	r := NewInlineRenderer(sub, &buf, size/4)
	r.SetStartTime(time.Now())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	const mutators = 8
	wg.Add(mutators)

	for range mutators {
		go func() {
			defer wg.Done()

			i := 0

			for {
				select {
				case <-ctx.Done():
					return
				default:
					id := NewActivityID(fmt.Sprintf("step-%04d", i%size))
					_ = sub.OnEvent(ctx, &benchEvent{
						eventType: EventActivityStarted, aid: id, aname: ActivityName("churn"),
					})
					i++
				}
			}
		}()
	}

	b.ResetTimer()

	for b.Loop() {
		buf.Reset()
		r.Draw()
	}

	b.StopTimer()
	cancel()
	wg.Wait()
}
