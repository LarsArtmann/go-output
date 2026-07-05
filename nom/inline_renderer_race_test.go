package nom

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// safeBuffer is a goroutine-safe bytes.Buffer wrapper so that the only
// unsynchronized state a concurrent Draw() can race on is the renderer's own
// fields (prevLines), isolating the data-race under test.
type safeBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *safeBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	n, err := s.buf.Write(p)

	return n, err //nolint:wrapcheck // bytes.Buffer.Write never errors
}

func (s *safeBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.buf.String()
}

// TestInlineRenderer_ConcurrentDraw_NoRace drives Draw() from many goroutines
// at the same time as the background refresh loop. Run with -race: it must be
// data-race clean. Draw() mutates prevLines and writes to the writer with no
// lock, so before the fix the race detector flags concurrent prevLines access.
func TestInlineRenderer_ConcurrentDraw_NoRace(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf safeBuffer

	renderer := NewInlineRenderer(sub, &buf, 20)

	ctx := context.Background()

	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-race"), "")
	for i := range 10 {
		sendActivityStarted(t, sub, ctx,
			ActivityID(fmt.Sprintf("step%d", i)),
			ActivityName(fmt.Sprintf("Step %d", i)))
	}

	renderer.Start(ctx, 5*time.Millisecond) // fast background loop
	t.Cleanup(renderer.Stop)

	var wg sync.WaitGroup

	for range 8 {
		wg.Go(func() {
			for range 40 {
				renderer.Draw()
			}
		})
	}

	wg.Wait()
}

// TestInlineRenderer_FrameShrink_ClearsStaleLines renders a tall frame, then a
// shorter one (completed children get pruned under height pressure), and
// asserts that the second redraw clears EVERY line of the previous frame.
//
// Before the fix, Draw() moved the cursor up by the old line count but only
// cleared/rewrote the new (smaller) line count, leaving ghost lines from the
// previous taller frame visible on screen.
func TestInlineRenderer_FrameShrink_ClearsStaleLines(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf bytes.Buffer

	// Height cap prunes completed children once they all finish.
	renderer := newInlineTestRenderer(sub, &buf, 4)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-shrink"), "")

	const groupID = ActivityID("group")
	registerActivity(sub, ctx, groupID, ActivityName("Group"))

	children := make([]string, 0, 6)
	for i := range 6 {
		children = append(children, fmt.Sprintf("child%d", i))
		registerActivity(sub, ctx,
			ActivityID(children[i]),
			ActivityName(fmt.Sprintf("Child %d", i)),
			groupID)
	}

	renderer.SetStartTime(time.Now())

	// First frame: group + capped children => tall.
	renderer.Draw()

	if clears := strings.Count(buf.String(), ansiClearLine); clears != 0 {
		t.Fatalf("first draw must not clear any lines, got %d", clears)
	}

	buf.Reset()

	// Complete every child so elideCompletedUnderPressure prunes them.
	for _, c := range children {
		sendActivityCompleted(t, sub, ctx,
			ActivityID(c), ActivityName("Child"), time.Second)
	}

	// Second frame: group only => shorter.
	renderer.Draw()

	output := buf.String()

	// Invariant: however many lines the redraw moved the cursor UP, it must
	// CLEAR at least that many — otherwise the previous taller frame leaves
	// ghost lines on screen. The cursor-up count is the previous prevLines.
	movedUp := cursorUpLines(output)
	clears := strings.Count(output, ansiClearLine)

	if movedUp <= 0 {
		t.Fatalf("setup error: redraw did not move the cursor up; output:\n%q", output)
	}

	if clears < movedUp {
		t.Errorf("redraw moved cursor up %d lines but only cleared %d — ghost lines remain\noutput:\n%q",
			movedUp, clears, output)
	}
}

// TestInlineRenderer_Stop_DoesNotDeadlockWithRenderLoop stress-calls Stop while
// the background loop is actively rendering. The old Stop() held tickMu (write)
// across <-tickerDone, but Draw() needs tickMu.RLock — so any Stop that raced an
// in-flight render deadlocked forever (tickerDone never closed). This must finish.
func TestInlineRenderer_Stop_DoesNotDeadlockWithRenderLoop(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf safeBuffer

	renderer := NewInlineRenderer(sub, &buf, 20)

	ctx := context.Background()

	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-stop"), "")
	for i := range 5 {
		sendActivityStarted(t, sub, ctx,
			ActivityID(fmt.Sprintf("step%d", i)),
			ActivityName(fmt.Sprintf("Step %d", i)))
	}

	renderer.Start(ctx, time.Millisecond) // render every 1ms

	// Keep the loop saturated with refresh requests while we stop it.
	stopDone := make(chan struct{})

	go func() {
		defer close(stopDone)

		for range 200 {
			renderer.Refresh()
			runtime.Gosched()
		}

		renderer.Stop()
	}()

	select {
	case <-stopDone:
	case <-time.After(5 * time.Second):
		t.Fatal("Stop() deadlocked against the render loop")
	}
}

// TestInlineRenderer_RenderRacingActivityMutation reproduces the core garbled-
// output bug: the dependency tree embeds the shared *Activity pointer, so event
// handlers (SetRunning/SetCompleted/SetFailed) mutate Activity fields while the
// render loop reads them. Before the fix, rendering was NOT under the subscriber
// lock, so -race flagged reads of Status/Symbol/Color mid-write and frames came
// out garbled/inconsistently sorted. RenderTree now takes ns.mu.RLock.
func TestInlineRenderer_RenderRacingActivityMutation(t *testing.T) {
	t.Parallel()

	sub := newTestSubscriber(t)

	var buf safeBuffer

	renderer := NewInlineRenderer(sub, &buf, 30)
	renderer.SetNoColor(true)

	ctx := context.Background()
	_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-mut"), "")

	const groupID = ActivityID("g")
	registerActivity(sub, ctx, groupID, ActivityName("Group"))

	ids := make([]ActivityID, 0, 12)
	for i := range 12 {
		ids = append(ids, ActivityID(fmt.Sprintf("act%d", i)))
		registerActivity(sub, ctx, ids[i], ActivityName(fmt.Sprintf("Act %d", i)), groupID)
	}

	renderer.SetStartTime(time.Now())
	renderer.Start(ctx, time.Millisecond) // render every 1ms
	t.Cleanup(renderer.Stop)

	// Hammer activity lifecycle events from many goroutines while rendering.
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)

		go func(id ActivityID) {
			defer wg.Done()

			for j := range 20 {
				switch j % 3 {
				case 0:
					sendActivityStarted(t, sub, ctx, id, ActivityName("Act"))
				case 1:
					sendActivityCompleted(t, sub, ctx, id, ActivityName("Act"), 5*time.Millisecond)
				case 2:
					_ = sub.OnEvent(ctx, ActivityFailed{
						ID:   id,
						Name: ActivityName("Act"),
						Err:  errors.New("boom"),
					})
				}
			}
		}(id)
	}

	// Also poke Draw directly to maximize render/mutation overlap.

	wg.Go(func() {
		for range 80 {
			renderer.Draw()
			runtime.Gosched()
		}
	})

	wg.Wait()
}

var cursorUpRe = regexp.MustCompile(`\x1b\[(\d+)A`)

// TestInlineRenderer_Finish_RacingSetters drives Finish() while another
// goroutine flips the config setters (SetStartTime/SetAppName/SetNoColor/
// SetHideCursor). All setters take tickMu.Lock; Finish read startTime/appName/
// noColor/hideCursor unlocked, which -race flagged. The snapshot-under-RLock
// at the top of Finish closes it.
func TestInlineRenderer_Finish_RacingSetters(t *testing.T) {
	t.Parallel()

	for range 20 {
		sub := newTestSubscriber(t)

		var buf safeBuffer

		renderer := NewInlineRenderer(sub, &buf, 10)

		ctx := context.Background()
		_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-fin"), "")
		sendActivityStarted(t, sub, ctx, ActivityID("s"), ActivityName("S"))
		renderer.SetStartTime(time.Now())
		renderer.Draw() // establish a frame so Finish has prevLines to clear

		stop := make(chan struct{})

		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					renderer.SetStartTime(time.Now())
					renderer.SetAppName("X")
					renderer.SetNoColor(true)
					renderer.SetHideCursor(false)
				}
			}
		}()

		renderer.Finish(errors.New("done"))
		close(stop)
	}
}

// TestInlineRenderer_RenderSummary_RacingSetStartTime drives Draw() (which
// calls renderSummary) while another goroutine flips SetStartTime. Before the
// fix, renderSummary read r.startTime.IsZero() WITHOUT the lock and then
// re-read r.startTime UNDER the lock — a classic TOCTOU data race. The snapshot
// is now taken once under tickMu.RLock.
func TestInlineRenderer_RenderSummary_RacingSetStartTime(t *testing.T) {
	t.Parallel()

	for range 20 {
		sub := newTestSubscriber(t)

		var buf safeBuffer

		renderer := NewInlineRenderer(sub, &buf, 10)

		ctx := context.Background()
		_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-st"), "")
		sendActivityStarted(t, sub, ctx, ActivityID("s"), ActivityName("S"))

		stop := make(chan struct{})

		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					renderer.SetStartTime(time.Now())
					renderer.SetStartTime(time.Time{})
				}
			}
		}()

		for range 100 {
			renderer.Draw()
			runtime.Gosched()
		}

		close(stop)
	}
}

// cursorUpLines returns N from the first "\033[NA" escape in output, or 0.
func cursorUpLines(output string) int {
	m := cursorUpRe.FindStringSubmatch(output)
	if m == nil {
		return 0
	}

	var n int
	for _, r := range m[1] {
		n = n*10 + int(r-'0')
	}

	return n
}

// TestInlineRenderer_RenderCompletion_RacingSetters drives RenderCompletion
// (which reads appName via snapshotConfig) while another goroutine flips
// SetAppName/SetMaxHeight/SetPlainText. Before the C1 fix, RenderCompletion
// read r.appName directly without any lock. Run with -race.
func TestInlineRenderer_RenderCompletion_RacingSetters(t *testing.T) {
	t.Parallel()

	for range 20 {
		sub := newTestSubscriber(t)

		var buf safeBuffer

		renderer := NewInlineRenderer(sub, &buf, 10)

		ctx := context.Background()
		_ = sendWorkflowStarted(sub, ctx, WorkflowID("wf-rc"), "")
		sendActivityStarted(t, sub, ctx, ActivityID("s"), ActivityName("S"))
		renderer.SetStartTime(time.Now())
		renderer.Draw() // establish a frame so RenderCompletion has prevLines to clear

		stop := make(chan struct{})

		go func() {
			for {
				select {
				case <-stop:
					return
				default:
					renderer.SetAppName("App-A")
					renderer.SetAppName("App-B")
					renderer.SetMaxHeight(5)
					renderer.SetPlainText(true)
					renderer.SetPlainText(false)
				}
			}
		}()

		for range 50 {
			renderer.RenderCompletion(CompletionResult{
				Success:    true,
				TotalSteps: 1,
				Elapsed:    5 * time.Second,
			})
			renderer.RenderCompletion(CompletionResult{
				Success:     false,
				FailedSteps: 1,
				TotalSteps:  1,
				Elapsed:     5 * time.Second,
			})
		}

		close(stop)
	}
}
