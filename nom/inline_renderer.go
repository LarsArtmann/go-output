package nom

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"

	output "github.com/larsartmann/go-output"
)

// ANSI escape sequences from github.com/charmbracelet/x/ansi — the proper
// charm library for terminal control codes. These replace the previous
// hand-rolled octal string constants (\033[...]) with the canonical,
// well-named, well-tested constants from the charm ecosystem.
//
// Reference: https://github.com/charmbracelet/x/ansi
const (
	// ansiClearLine clears the current line and returns cursor to column 0.
	// Combines \r (carriage return) with ansi.EraseLineRight (\x1b[K).
	ansiClearLine = "\r" + ansi.EraseLineRight
)

// InlineRenderer renders the NOM dependency tree to a writer using ANSI
// escape codes for in-place updates — no alt-screen takeover.
//
// Inspired by how nix-output-monitor renders under nh: each render call
// moves the cursor up past the previously drawn lines, clears them, and
// redraws the current tree state. The terminal scrolls naturally, so
// output from previous steps remains visible above.
//
// Frame diffing: Draw() compares the new frame against the last frame
// written. If nothing changed, zero bytes are emitted — no cursor-up,
// no clear, no repaint. This eliminates the repetition bug where every
// 200ms tick appended a full copy of the tree on terminals that don't
// support synchronized output (mode 2026). This mirrors bubbletea v2's
// cursedRenderer viewEquals() early-exit pattern.
type InlineRenderer struct {
	subscriber *NOMStyleSubscriber
	writer     io.Writer
	prevLines  int
	maxHeight  int
	hideCursor bool
	startTime  time.Time
	appName    string
	noColor    bool
	// plainText is set once at construction when the output is not an
	// interactive terminal (CI, piped, redirected). In plain mode Draw skips
	// cursor-manipulation/sync escape codes that corrupt non-interactive logs,
	// appending each frame on its own line instead.
	plainText bool

	// writerIsTTY is detected once at construction by checking whether the
	// writer is a real terminal (via term.IsTerminal on the writer's FD).
	// This gates synchronized-output (mode 2026) wrapping: sync codes are
	// only emitted when the writer is confirmed to be a TTY, preventing
	// corruption on pipes, buffers, and non-TTY writers.
	writerIsTTY bool

	// lastFrame is the last frame string written to the terminal. Draw()
	// compares the new frame against this; if identical, it skips the write
	// entirely (zero bytes). This is the core fix for the repetition bug.
	//
	// lastFramePlain tracks the plainText mode that was in effect when
	// lastFrame was stored. Draw() compares both frame content AND mode —
	// if the mode changed (e.g. SetPlainText called), the diff invalidates
	// even though the frame string is identical, because the OUTPUT differs
	// (plain text appends; inline mode wraps in cursor codes).
	//
	// Both fields are guarded exclusively by renderMu. SetPlainText signals
	// a mode change via snapshotConfig (under tickMu), never by writing
	// lastFrame directly — that would cross the two-mutex boundary and race.
	lastFrame      string
	lastFramePlain bool

	// estimatedRemaining is an optional callback that returns the estimated
	// total remaining time for the workflow (sum of pending-step estimates).
	// When non-nil and returning > 0, the summary bar renders "~Xm left".
	// This decouples the estimate source (BuildFlow's SQLite store) from the
	// renderer — the callback is invoked each frame under renderMu.
	estimatedRemaining func() time.Duration

	// showCriticalPathETA enables a "~Xm critical" segment in the summary bar
	// sourced from the DAG's longest remaining-time path. Off by default to
	// preserve the standard "~Xm left" behavior.
	showCriticalPathETA bool

	// showDAGSummary enables a structural summary segment ("N nodes · M edges ·
	// K layers") in the summary bar. Off by default.
	showDAGSummary bool

	// pendingLines holds external output (e.g. application log lines) that
	// should appear ABOVE the tree on the next Draw. This mirrors how
	// nix-output-manager (nom) interleaves nix build logs above the progress
	// tree: each frame clears the old tree, prints accumulated log lines,
	// then draws the tree below them.
	//
	// Guarded by renderMu (same lock as Draw/Finish) so EnqueueLines and Draw
	// never race on the slice.
	pendingLines []string

	tickMu       sync.RWMutex
	renderMu     sync.Mutex // serializes Draw/Finish terminal writes + prevLines
	cancelFn     context.CancelFunc
	tickerDone   chan struct{}
	refreshChan  chan struct{}
	renderNotify chan struct{} // test hook: signaled after each render if non-nil
	sigwinchStop chan struct{}

	// Dead-writer detection: when consecutive write errors exceed this
	// threshold, the renderer stops writing (no-op write) to avoid spinning
	// on a broken pipe. Guarded by renderMu.
	consecutiveWriteErrors int
	writerDead             bool
}

const maxConsecutiveWriteErrors = 10

// rendererConfig is an immutable snapshot of every configurable InlineRenderer
// field, captured under tickMu.RLock in one read. Draw/Finish/renderSummary all
// consume a single snapshot instead of doing scattered, piecemeal field reads
// (which were both racy for lock-free fields like maxHeight and redundant —
// renderSummary re-acquired the lock Draw already held).
type rendererConfig struct {
	hideCursor  bool
	noColor     bool
	plainText   bool
	appName     string
	startTime   time.Time
	maxHeight   int
	writerIsTTY bool
}

// snapshotConfig returns an immutable copy of all renderer configuration under
// a single tickMu.RLock. Callers then read freely without holding the lock,
// so the potentially slow terminal write (under renderMu) never blocks setters.
//
// Type safety: plainText is computed authoritatively here — if the writer is
// not a TTY, plainText is ALWAYS true regardless of what SetPlainText(false)
// was called with. This makes the impossible state (plain=false on a non-TTY
// writer, which would emit cursor codes to a pipe) unrepresentable.
func (r *InlineRenderer) snapshotConfig() rendererConfig {
	r.tickMu.RLock()
	defer r.tickMu.RUnlock()

	// Authoritative plainText: a non-TTY writer can NEVER use inline cursor
	// codes. SetPlainText(false) on a non-TTY is a no-op — the only valid
	// direction is degradation (TTY → plain), never upgrade (pipe → inline).
	effectivePlainText := r.plainText || !r.writerIsTTY

	return rendererConfig{
		hideCursor:  r.hideCursor,
		noColor:     r.noColor,
		plainText:   effectivePlainText,
		appName:     r.appName,
		startTime:   r.startTime,
		maxHeight:   r.maxHeight,
		writerIsTTY: r.writerIsTTY,
	}
}

// NewInlineRenderer creates an inline renderer bound to the given subscriber and writer.
// maxHeight caps the tree height; 0 means unlimited.
func NewInlineRenderer(subscriber *NOMStyleSubscriber, writer io.Writer, maxHeight int) *InlineRenderer {
	return &InlineRenderer{
		subscriber:  subscriber,
		writer:      writer,
		maxHeight:   maxHeight,
		hideCursor:  true,
		appName:     "Workflow",
		noColor:     detectNoColorForWriter(writer),
		plainText:   detectPlainTextForWriter(writer),
		writerIsTTY: writerIsTerminal(writer),
	}
}

// writerIsTerminal reports whether the given writer is a real terminal (TTY).
// Uses term.IsTerminal (ioctl TCGETS) — the same detection method as
// charmbracelet/termenv and bubbletea. Only *os.File writers can be TTYs;
// buffers, pipes, and other writers return false.
func writerIsTerminal(writer io.Writer) bool {
	if f, ok := writer.(*os.File); ok {
		//nolint:gosec // File descriptors are always small positive integers.
		return term.IsTerminal(int(f.Fd()))
	}

	return false
}

// detectPlainTextForWriter reports whether the renderer should emit plain,
// append-only output (no cursor/sync escape codes). Returns true when:
//   - Running under CI (output.IsCI), OR
//   - The writer is not a terminal (pipe, buffer, redirect)
//
// This replaces the old detectPlainText() which only checked CI, missing
// the common case of piped/redirected output where ANSI redraw codes
// produce the repetition bug.
func detectPlainTextForWriter(writer io.Writer) bool {
	if output.IsCI() {
		return true
	}

	return !writerIsTerminal(writer)
}

// detectNoColorForWriter reports whether color output should be suppressed.
// Checks NO_COLOR/TERM=dumb env vars via envdetect, then falls back to
// checking whether the WRITER (not os.Stdout) is a terminal. This fixes
// the mismatch where the old detectNoColor() checked os.Stdout but the
// renderer writes to an arbitrary writer (e.g. os.Stderr in BuildFlow).
func detectNoColorForWriter(writer io.Writer) bool {
	if output.IsNoColor() || output.IsCI() {
		return true
	}

	return !writerIsTerminal(writer)
}

// SetHideCursor controls whether the cursor is hidden during rendering (default: true).
// Thread-safe: may be called before or during the render loop.
func (r *InlineRenderer) SetHideCursor(hide bool) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.hideCursor = hide
}

// SetNoColor forces colorless output. By default, colors are enabled unless
// NO_COLOR, TERM=dumb, or lacking a terminal is detected.
// Thread-safe: may be called before or during the render loop.
func (r *InlineRenderer) SetNoColor(noColor bool) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.noColor = noColor
}

// SetAppName sets the application name for success/failure messages (default: "Workflow").
// Thread-safe: may be called before or during the render loop.
func (r *InlineRenderer) SetAppName(name string) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.appName = name
}

// SetStartTime sets the workflow start time for elapsed display.
// Thread-safe: may be called before or during the render loop.
func (r *InlineRenderer) SetStartTime(t time.Time) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.startTime = t
}

// SetMaxHeight updates the cap on rendered tree height (0 = unlimited).
// Safe to call concurrently with the render loop: Draw reads maxHeight via the
// lock-protected snapshotConfig(), so a resize takes effect on the next frame.
func (r *InlineRenderer) SetMaxHeight(maxHeight int) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.maxHeight = maxHeight
}

// SetPlainText forces plain, append-only output (no cursor/sync escape codes).
// By default, plainText is auto-detected at construction via detectPlainTextForWriter()
// (true under CI or when the writer is not a terminal).
//
// Type safety: if the writer is not a TTY, calling SetPlainText(false) is a
// no-op — a non-TTY writer can NEVER use inline cursor codes. This makes the
// impossible state (plain=false on a pipe) unrepresentable. SetPlainText only
// works in the degradation direction: TTY → plain.
//
// Thread-safe: may be called before or during the render loop.
func (r *InlineRenderer) SetPlainText(plain bool) {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	r.plainText = plain
}

// SetEstimatedRemainingFunc sets an optional callback that returns the
// estimated total remaining time for the workflow. When non-nil and returning
// a positive duration, the summary bar renders "~Xm left" after the elapsed
// time. This decouples the estimate source (e.g. BuildFlow's SQLite timing
// store) from the renderer.
//
// Thread-safe: may be called before or during the render loop. The callback
// itself is invoked under renderMu (not tickMu), so it must not acquire any
// renderer lock.
func (r *InlineRenderer) SetEstimatedRemainingFunc(fn func() time.Duration) {
	r.renderMu.Lock()
	r.estimatedRemaining = fn
	r.renderMu.Unlock()
}

// SetShowCriticalPathETA enables a "~Xm critical" segment in the summary bar
// sourced from the DAG's longest remaining-time path. When disabled (default),
// the summary bar only shows the remaining-time callback or the elapsed time.
//
// Thread-safe: may be called before or during the render loop.
func (r *InlineRenderer) SetShowCriticalPathETA(show bool) {
	r.renderMu.Lock()
	r.showCriticalPathETA = show
	r.renderMu.Unlock()
}

// SetShowDAGSummary enables a structural DAG summary segment ("N nodes · M
// edges · K layers") in the summary bar. Off by default.
//
// Thread-safe: may be called before or during the render loop.
func (r *InlineRenderer) SetShowDAGSummary(show bool) {
	r.renderMu.Lock()
	r.showDAGSummary = show
	r.renderMu.Unlock()
}

// EnqueueLines queues external output lines (e.g. application log messages)
// to be printed ABOVE the tree on the next Draw call. This mirrors how
// nix-output-manager (nom) interleaves nix build logs above the progress tree:
// each frame clears the old tree region, prints accumulated log lines, then
// draws the fresh tree below them.
//
// Thread-safe: acquires renderMu so it never races with Draw. Safe to call
// before or during the render loop. Lines are drained atomically on the next
// Draw — no partial interleaving.
func (r *InlineRenderer) EnqueueLines(lines []string) {
	if len(lines) == 0 {
		return
	}

	r.renderMu.Lock()
	r.pendingLines = append(r.pendingLines, lines...)
	r.renderMu.Unlock()
}

// HasPendingLines reports whether there are queued external lines waiting
// for the next Draw. Thread-safe.
func (r *InlineRenderer) HasPendingLines() bool {
	r.renderMu.Lock()
	defer r.renderMu.Unlock()

	return len(r.pendingLines) > 0
}

// Draw renders one frame to the configured io.Writer.
//
// Frame diffing: if the frame content is identical to the last frame written
// (same tree state, same elapsed-time rounding), Draw emits ZERO bytes — no
// cursor-up, no clear, no sync-output wrapping. This eliminates the repetition
// issue where every 200ms tick appended a full copy of the tree on terminals
// that don't support synchronized output. This mirrors bubbletea v2's
// cursedRenderer viewEquals() early-exit pattern.
//
// Split-brain M4 resolved: Render() is reserved for the output.Renderer contract.
func (r *InlineRenderer) Draw() {
	if r.subscriber == nil {
		return
	}

	r.renderMu.Lock()
	defer r.renderMu.Unlock()

	cfg := r.snapshotConfig()

	maxH := r.effectiveMaxHeight(cfg.maxHeight)
	maxW := r.effectiveMaxWidth()

	// Render from an immutable snapshot — no lock held during the tree walk,
	// so event handlers can proceed concurrently without racing the render.
	frame, hasTree := r.subscriber.RenderSnapshot(maxH, maxW)

	pending := r.pendingLines
	hasPending := len(pending) > 0

	// If there is nothing to render at all, bail out.
	if (!hasTree || frame == msgNoActivitiesToDisplay) && !hasPending {
		return
	}

	// If there is no tree but we have pending log lines, drain them and
	// return. This handles the rare case where logs arrive before any step
	// is registered.
	if !hasTree || frame == msgNoActivitiesToDisplay {
		r.drainPendingPlain(pending, cfg)
		return
	}

	summary := r.renderSummary(cfg.startTime)
	if summary != "" {
		frame += "\n" + summary
	}

	// Frame diffing: skip the write entirely if nothing changed since the
	// last frame AND there are no pending log lines.
	if !hasPending && frame == r.lastFrame && cfg.plainText == r.lastFramePlain {
		return
	}

	if cfg.plainText {
		r.drawPlainText(frame, pending, hasPending, cfg)
		return
	}

	r.drawInline(frame, pending, hasPending, maxW, cfg)
}

// drawPlainText handles CI/non-terminal output: appends pending log lines
// and the tree frame (only if changed) without any cursor manipulation.
// This prevents massive tree repetition when logs arrive frequently but the
// tree state hasn't changed.
func (r *InlineRenderer) drawPlainText(frame string, pending []string, hasPending bool, cfg rendererConfig) {
	treeChanged := frame != r.lastFrame

	if hasPending || treeChanged {
		var b strings.Builder

		for _, line := range pending {
			b.WriteString(line)
			b.WriteByte('\n')
		}

		if treeChanged {
			b.WriteString(frame)
			b.WriteByte('\n')
		}

		r.write(b.String())
	}

	r.pendingLines = nil
	r.lastFrame = frame
	r.lastFramePlain = cfg.plainText
}

// drawInline handles terminal inline rendering: cursor-up/clear-line redraw
// with synchronized output (mode 2026) and nom-style log interleaving.
func (r *InlineRenderer) drawInline(frame string, pending []string, hasPending bool, maxW int, cfg rendererConfig) {
	physicalLines := PhysicalLineCount(frame, maxW)
	if physicalLines == 0 && !hasPending {
		return
	}

	var output string

	if hasPending {
		// nom-style interleaving: clear the old tree frame, write log lines
		// above, then draw the tree as a fresh frame below them.
		output = buildLogAndTreeOutput(frame, pending, r.prevLines, physicalLines, cfg.hideCursor)
	} else {
		output = buildRedrawOutput(frame, r.prevLines, physicalLines, cfg.hideCursor)
	}

	r.prevLines = physicalLines

	if cfg.writerIsTTY {
		r.write(ansi.SetModeSynchronizedOutput + output + ansi.ResetModeSynchronizedOutput)
	} else {
		r.write(output)
	}

	r.pendingLines = nil
	r.lastFrame = frame
	r.lastFramePlain = cfg.plainText
}

// drainPendingPlain writes pending log lines directly to the writer without
// any cursor manipulation. Used when there is no tree frame to coordinate
// with (e.g. logs arriving before any step is registered).
func (r *InlineRenderer) drainPendingPlain(pending []string, cfg rendererConfig) {
	if len(pending) == 0 {
		return
	}

	var b strings.Builder
	for _, line := range pending {
		b.WriteString(line)
		b.WriteByte('\n')
	}

	if cfg.writerIsTTY {
		r.write(ansi.SetModeSynchronizedOutput + b.String() + ansi.ResetModeSynchronizedOutput)
	} else {
		r.write(b.String())
	}

	r.pendingLines = nil
	r.prevLines = 0
}

// buildLogAndTreeOutput assembles the ANSI payload for a frame that includes
// external log lines above the tree. It:
//  1. Clears the previous tree frame (cursor up + clear each line)
//  2. Returns the cursor to the top of the cleared region
//  3. Writes the pending log lines (they scroll naturally if the terminal
//     is at the bottom — old content scrolls into scrollback)
//  4. Writes the tree as a fresh frame (no cursor-up, just append)
//
// The net visual effect is: [logs above] [tree below], exactly like nom.
func buildLogAndTreeOutput(frame string, pending []string, prevLines, physicalLines int, hideCursor bool) string {
	var b strings.Builder

	// 1. Clear old tree frame if one exists.
	if prevLines > 0 {
		b.WriteString(ansi.CursorUp(prevLines))
		b.WriteString("\r")

		for range prevLines {
			b.WriteString(ansiClearLine)
			b.WriteString("\n")
		}

		// Return to the top of the cleared region so pending lines fill
		// the blank space rather than appending below it.
		b.WriteString(ansi.CursorUp(prevLines))
		b.WriteString("\r")
	} else if hideCursor {
		// First frame ever + pending logs: hide cursor (normally done in
		// buildRedrawOutput's first-frame branch).
		b.WriteString(ansi.HideCursor)
	}

	// 2. Write pending log lines.
	for _, line := range pending {
		b.WriteString(line)
		b.WriteString("\n")
	}

	// 3. Write tree as a fresh frame below the logs.
	b.WriteString(frame)
	b.WriteString("\n")

	// 4. When the tree shrank, wipe leftover lines from the taller previous
	// frame so pruned subtrees leave no ghosts — same logic as buildRedrawOutput.
	if extra := prevLines - physicalLines; extra > 0 {
		for range extra {
			b.WriteString("\n")
			b.WriteString(ansiClearLine)
		}

		b.WriteString(ansi.CursorUp(extra))
		b.WriteString("\r")
	}

	return b.String()
}

// buildRedrawOutput assembles the ANSI payload for one in-place redraw. On the
// first frame (prevLines == 0) it just emits the frame (plus an optional cursor
// hide); on subsequent frames it scrolls up, repaints each line, and wipes any
// leftover lines when the frame shrank (so pruned subtrees leave no ghosts).
//
// All escape sequences use github.com/charmbracelet/x/ansi constants:
//   - ansi.CursorUp(n): move cursor up n lines
//   - ansi.EraseLineRight: clear from cursor to end of line (\x1b[K)
//   - ansi.HideCursor / ansi.ShowCursor: cursor visibility
func buildRedrawOutput(frame string, prevLines, physicalLines int, hideCursor bool) string {
	if prevLines == 0 {
		output := ""
		if hideCursor {
			output = ansi.HideCursor
		}

		return output + frame + "\n"
	}

	// Move back to the top of the previous frame and repaint every line.
	var b strings.Builder

	b.WriteString(ansi.CursorUp(prevLines))
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

		b.WriteString(ansi.CursorUp(extra))
		b.WriteString("\r")
	}

	b.WriteString("\n")

	return b.String()
}

// write writes a string to the renderer's writer. When consecutive I/O errors
// exceed maxConsecutiveWriteErrors, the writer is marked dead and subsequent
// writes become silent no-ops — this prevents spinning on a broken pipe.
func (r *InlineRenderer) write(s string) {
	if r.writerDead {
		return
	}

	if _, err := fmt.Fprint(r.writer, s); err != nil {
		r.consecutiveWriteErrors++
		if r.consecutiveWriteErrors >= maxConsecutiveWriteErrors {
			r.writerDead = true
		}

		return
	}

	r.consecutiveWriteErrors = 0
}

// Start begins periodic background rendering at the given interval.
// Call Stop to terminate the background goroutine.
// The context can be used to cancel independently of Stop.
//
// Also starts a SIGWINCH listener for terminal resize handling: on resize,
// the lastFrame cache is invalidated so the next tick forces a full redraw.
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

	// SIGWINCH listener: on terminal resize, invalidate the frame cache so
	// the next Draw() emits a full redraw (width/height may have changed).
	// This mirrors bubbletea's listenForResize → resize() → Erase() pattern.
	r.sigwinchStop = make(chan struct{})
	go r.listenForResize(ctx)
}

// listenForResize listens for SIGWINCH (terminal resize) signals and
// invalidates the frame cache so the next Draw() forces a full redraw.
func (r *InlineRenderer) listenForResize(ctx context.Context) {
	// Capture locally so Stop() can safely nil r.sigwinchStop after we exit.
	done := r.sigwinchStop
	defer close(done)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGWINCH)
	defer signal.Stop(sigChan)

	for {
		select {
		case <-ctx.Done():
			return
		case <-sigChan:
			// Invalidate the frame cache: the terminal dimensions changed,
			// so even if the tree state is the same, the rendered frame
			// (wrapping, truncation) may differ. Force a full redraw.
			r.renderMu.Lock()
			r.lastFrame = ""
			r.renderMu.Unlock()

			r.Refresh()
		}
	}
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
//
// Panic recovery: if Draw() panics, the deferred recover ensures the cursor is
// restored (shown) before the goroutine exits. This prevents the terminal from
// being left in a broken state (hidden cursor) after a crash — mirroring
// bubbletea's recoverFromPanic → restoreTerminalState() pattern.
func (r *InlineRenderer) refreshLoop(ctx context.Context, interval time.Duration) {
	// Capture the done channel locally so Stop() can safely nil r.tickerDone
	// after we exit — the deferred close reads this local, not the field.
	done := r.tickerDone
	defer close(done)
	// Panic recovery: restore terminal state on crash. The cursor may have
	// been hidden by Draw(); show it so the terminal isn't left broken.
	defer func() {
		if rec := recover(); rec != nil {
			r.renderMu.Lock()
			if r.prevLines > 0 {
				r.write(ansi.ShowCursor)
				r.prevLines = 0
			}

			r.lastFrame = ""
			r.renderMu.Unlock()
		}
	}()

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

	r.tickMu.RLock()
	notify := r.renderNotify
	r.tickMu.RUnlock()

	if notify != nil {
		select {
		case notify <- struct{}{}:
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
	sigStop := r.sigwinchStop
	r.tickMu.Unlock()

	if cancelFn == nil {
		return
	}

	cancelFn()
	<-done

	if sigStop != nil {
		<-sigStop
	}

	r.tickMu.Lock()
	r.cancelFn = nil
	r.tickerDone = nil
	r.sigwinchStop = nil
	r.refreshChan = nil
	r.tickMu.Unlock()
}

// CompletionResult is the final status passed to RenderCompletion.
// It carries enough information for a one-line pass/fail summary.
type CompletionResult struct {
	// Success is true when the workflow completed without errors.
	Success bool
	// Elapsed is the total wall-clock duration of the workflow.
	Elapsed time.Duration
	// TotalSteps is the number of steps that were executed.
	TotalSteps int
	// FailedSteps is the number of steps that failed (0 when Success is true).
	FailedSteps int
}

// RenderCompletion renders a final one-line completion summary to the writer.
// It clears any remaining progress frame, shows the cursor, and prints the
// result. Call this after Stop() to produce the final output line.
//
// The summary line format:
//
//	✓ Workflow completed (42 steps, 12.3s)
//	✗ Workflow failed (3/42 steps failed, 45.6s)
func (r *InlineRenderer) RenderCompletion(result CompletionResult) {
	r.renderMu.Lock()
	defer r.renderMu.Unlock()

	// Restore cursor if it was hidden during rendering.
	if r.prevLines > 0 {
		r.write(ansi.ShowCursor)
		r.prevLines = 0
	}

	r.lastFrame = ""

	status := "✓"
	if !result.Success {
		status = "✗"
	}

	var detail string

	if result.Success {
		detail = fmt.Sprintf("%d steps, %s", result.TotalSteps, formatDuration(result.Elapsed))
	} else {
		detail = fmt.Sprintf("%d/%d steps failed, %s",
			result.FailedSteps, result.TotalSteps, formatDuration(result.Elapsed))
	}

	cfg := r.snapshotConfig()

	r.write(fmt.Sprintf("%s %s completed (%s)\n", status, cfg.appName, detail))
}

// formatDuration renders a time.Duration in a compact human-readable form.
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}

	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}

	return d.Round(time.Second).String()
}
