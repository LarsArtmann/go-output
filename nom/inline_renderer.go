package nom

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/term"

	"github.com/larsartmann/go-output/envdetect"
)

const (
	ansiCursorUp    = "\033[A"
	ansiCursorUpN   = "\033[%dA"
	ansiClearLine   = "\r\033[K"
	ansiClearScreen = "\033[2J\033[H"
	ansiHideCursor  = "\033[?25l"
	ansiShowCursor  = "\033[?25h"
	ansiSyncBegin   = "\033[?2026h"
	ansiSyncEnd     = "\033[?2026l"
)

// InlineRenderer renders the NOM dependency tree to a writer using ANSI
// escape codes for in-place updates — no alt-screen takeover.
//
// Inspired by how nix-output-monitor renders under nh: each render call
// moves the cursor up past the previously drawn lines, clears them, and
// redraws the current tree state. The terminal scrolls naturally, so
// output from previous steps remains visible above.
type InlineRenderer struct {
	subscriber *NOMStyleSubscriber
	writer     io.Writer
	prevLines  int
	maxHeight  int
	hideCursor bool
	startTime  time.Time
	appName    string
	noColor    bool

	tickMu       sync.RWMutex
	renderMu     sync.Mutex // serializes Draw/Finish terminal writes + prevLines
	cancelFn     context.CancelFunc
	tickerDone   chan struct{}
	refreshChan  chan struct{}
	renderNotify chan struct{} // test hook: signaled after each render if non-nil
}

// rendererConfig is an immutable snapshot of every configurable InlineRenderer
// field, captured under tickMu.RLock in one read. Draw/Finish/renderSummary all
// consume a single snapshot instead of doing scattered, piecemeal field reads
// (which were both racy for lock-free fields like maxHeight and redundant —
// renderSummary re-acquired the lock Draw already held).
type rendererConfig struct {
	hideCursor bool
	noColor    bool
	appName    string
	startTime  time.Time
	maxHeight  int
}

// snapshotConfig returns an immutable copy of all renderer configuration under
// a single tickMu.RLock. Callers then read freely without holding the lock,
// so the potentially slow terminal write (under renderMu) never blocks setters.
func (r *InlineRenderer) snapshotConfig() rendererConfig {
	r.tickMu.RLock()
	defer r.tickMu.RUnlock()

	return rendererConfig{
		hideCursor: r.hideCursor,
		noColor:    r.noColor,
		appName:    r.appName,
		startTime:  r.startTime,
		maxHeight:  r.maxHeight,
	}
}

// NewInlineRenderer creates an inline renderer bound to the given subscriber and writer.
// maxHeight caps the tree height; 0 means unlimited.
func NewInlineRenderer(subscriber *NOMStyleSubscriber, writer io.Writer, maxHeight int) *InlineRenderer {
	return &InlineRenderer{
		subscriber: subscriber,
		writer:     writer,
		maxHeight:  maxHeight,
		hideCursor: true,
		appName:    "Workflow",
		noColor:    detectNoColor(),
	}
}

// SetHideCursor controls whether the cursor is hidden during rendering (default: true).
// Must be called before Start(); concurrent calls after Start race with the render loop.
func (r *InlineRenderer) SetHideCursor(hide bool) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.hideCursor = hide
}

// SetNoColor forces colorless output. By default, colors are enabled unless
// NO_COLOR, TERM=dumb, or lacking a terminal is detected.
// Must be called before Start(); concurrent calls after Start race with the render loop.
func (r *InlineRenderer) SetNoColor(noColor bool) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.noColor = noColor
}

// SetAppName sets the application name for success/failure messages (default: "Workflow").
// Must be called before Start(); concurrent calls after Start race with the render loop.
func (r *InlineRenderer) SetAppName(name string) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.appName = name
}

// detectNoColor reports whether color output should be suppressed in the
// NOM inline renderer. The CI and NO_COLOR portions delegate to
// envdetect so root and nom stay aligned. The terminal fallback uses
// term.IsTerminal directly because stdoutIsTerminal closures from root
// are not available here.
func detectNoColor() bool {
	if envdetect.IsNoColor() || envdetect.IsCI() {
		return true
	}

	//nolint:gosec // File descriptors are always small positive integers.
	return !term.IsTerminal(int(os.Stdout.Fd()))
}

// SetStartTime sets the workflow start time for elapsed display.
// Must be called before Start(); concurrent calls after Start race with the render loop.
func (r *InlineRenderer) SetStartTime(t time.Time) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.startTime = t
}

// Render redraws the dependency tree in-place. On the first call it just prints.
// On subsequent calls it moves the cursor up to overwrite the previous frame.
//
// Draw renders one frame to the configured io.Writer. Unlike output.Renderer.Render()
// (which returns (string, error)), Draw writes directly to the writer and returns nothing —
// it is an incremental terminal redraw, not a one-shot format render.
// Split-brain M4 resolved: Render() is now reserved for the output.Renderer contract.
func (r *InlineRenderer) Draw() {
	if r.subscriber == nil {
		return
	}

	r.renderMu.Lock()
	defer r.renderMu.Unlock()

	cfg := r.snapshotConfig()

	r.subscriber.UpdateRunningActivityElapsed()

	maxH := r.effectiveMaxHeight(cfg.maxHeight)
	maxW := r.effectiveMaxWidth()

	// Render from an immutable snapshot — no lock held during the tree walk,
	// so event handlers can proceed concurrently without racing the render.
	frame, hasTree := r.subscriber.RenderSnapshot(maxH, maxW)
	if !hasTree || frame == msgNoActivitiesToDisplay {
		return
	}

	summary := r.renderSummary(cfg.startTime)
	if summary != "" {
		frame += "\n" + summary
	}

	// Count physical lines including wrapping so the next redraw lands correctly.
	physicalLines := PhysicalLineCount(frame, maxW)
	if physicalLines == 0 {
		return
	}

	var output string

	prevLines := r.prevLines

	if prevLines == 0 {
		if cfg.hideCursor {
			output = ansiHideCursor
		}

		output += frame + "\n"
	} else {
		// Move back to the top of the previous frame and repaint every line.
		var b strings.Builder

		fmt.Fprintf(&b, ansiCursorUpN, prevLines)
		b.WriteString("\r")

		frameLines := strings.Split(frame, "\n")

		for i, line := range frameLines {
			b.WriteString(ansiClearLine)
			b.WriteString(line)

			if i < len(frameLines)-1 {
				b.WriteString("\n")
			}
		}

		// When the frame shrank, the previous (taller) frame's tail is still on
		// screen below the new content. Wipe those leftover lines and park the
		// cursor just beneath the new frame so the next redraw lines up. Without
		// this, shrinking trees (e.g. completed children pruned under height
		// pressure) leave ghost lines behind.
		if extra := prevLines - physicalLines; extra > 0 {
			for range extra {
				b.WriteString("\n")
				b.WriteString(ansiClearLine)
			}

			fmt.Fprintf(&b, ansiCursorUpN, extra)
			b.WriteString("\r")
		}

		b.WriteString("\n")

		output = b.String()
	}

	r.write(ansiSyncBegin + output + ansiSyncEnd)

	r.prevLines = physicalLines
}

// write writes a string to the renderer's writer, ignoring errors.
// Terminal output is best-effort: a broken pipe should not crash the render loop.
func (r *InlineRenderer) write(s string) {
	_, _ = fmt.Fprint(r.writer, s)
}

// Start begins periodic background rendering at the given interval.
// Call Stop to terminate the background goroutine.
// The context can be used to cancel independently of Stop.
func (r *InlineRenderer) Start(ctx context.Context, interval time.Duration) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	if r.cancelFn != nil {
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	r.cancelFn = cancel
	r.tickerDone = make(chan struct{})
	r.refreshChan = make(chan struct{}, 1)

	go r.refreshLoop(ctx, interval)
}

// Refresh signals the background loop to redraw as soon as possible.
// Safe to call even when Start has not been called.
func (r *InlineRenderer) Refresh() {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	if r.refreshChan == nil {
		return
	}

	select {
	case r.refreshChan <- struct{}{}:
	default:
	}
}

// refreshLoop drives the periodic redraw. It renders on every tick, on explicit
// refresh signals, and at least once per second so elapsed timers stay current.
func (r *InlineRenderer) refreshLoop(ctx context.Context, interval time.Duration) {
	defer close(r.tickerDone)

	if interval <= 0 {
		interval = 100 * time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	maxFrame := time.NewTicker(time.Second)
	defer maxFrame.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-r.refreshChan:
			r.renderAndNotify()
		case <-ticker.C:
			r.renderAndNotify()
		case <-maxFrame.C:
			r.renderAndNotify()
		}
	}
}

// renderAndNotify calls Draw and, if a renderNotify channel is set,
// sends a non-blocking signal. This provides deterministic synchronization
// for tests without affecting production behavior.
func (r *InlineRenderer) renderAndNotify() {
	r.Draw()

	if r.renderNotify != nil {
		select {
		case r.renderNotify <- struct{}{}:
		default:
		}
	}
}

// Stop terminates the background refresh goroutine started by Start.
// It is safe to call Stop even if Start was never called.
//
// The cancel + wait happens OUTSIDE tickMu: the refresh loop's Draw() takes
// tickMu.RLock(), so holding the write lock across <-tickerDone would deadlock
// whenever Stop races a render already in flight (the loop could never acquire
// RLock to finish, so tickerDone would never close).
func (r *InlineRenderer) Stop() {
	r.tickMu.Lock()
	cancelFn := r.cancelFn
	done := r.tickerDone
	r.tickMu.Unlock()

	if cancelFn == nil {
		return
	}

	cancelFn()
	<-done

	r.tickMu.Lock()
	r.cancelFn = nil
	r.tickerDone = nil
	r.refreshChan = nil
	r.tickMu.Unlock()
}
