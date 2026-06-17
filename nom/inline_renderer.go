package nom

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"
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
	cancelFn     context.CancelFunc
	tickerDone   chan struct{}
	refreshChan  chan struct{}
	renderNotify chan struct{} // test hook: signaled after each render if non-nil
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

// detectNoColor checks if color output should be suppressed.
// Mirrors color.go isNoColor()+isCI()+isStdoutTerminal() — keep in sync.
func detectNoColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return true
	}

	if os.Getenv("TERM") == "dumb" {
		return true
	}

	if os.Getenv("CI") != "" ||
		os.Getenv("GITHUB_ACTIONS") != "" ||
		os.Getenv("GITLAB_CI") != "" ||
		os.Getenv("JENKINS_URL") != "" ||
		os.Getenv("BUILDKITE") != "" {
		return true
	}

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
// TODO(split-brain M4): This method has a different signature than output.Renderer.Render()
// (void return vs (string, error)). Consider renaming to Draw() in the next minor version
// to reserve Render() for the output.Renderer contract project-wide.
func (r *InlineRenderer) Render() {
	if r.subscriber == nil {
		return
	}

	// Snapshot config under RLock to avoid racing with setters.
	r.tickMu.RLock()
	hideCursor := r.hideCursor
	r.tickMu.RUnlock()

	r.subscriber.UpdateRunningActivityElapsed()
	r.subscriber.SyncActivityTimingToTree()

	tree := r.subscriber.GetDependencyTree()
	if tree == nil {
		return
	}

	maxH := r.effectiveMaxHeight()
	maxW := r.effectiveMaxWidth()

	frame := tree.RenderWithWidth(maxH, maxW)

	if frame == msgNoActivitiesToDisplay {
		return
	}

	summary := r.renderSummary()
	if summary != "" {
		frame += "\n" + summary
	}

	// Count physical lines including wrapping so the next redraw lands correctly.
	physicalLines := PhysicalLineCount(frame, maxW)
	if physicalLines == 0 {
		return
	}

	var output string

	if r.prevLines == 0 {
		if hideCursor {
			output = ansiHideCursor
		}
	} else {
		output = fmt.Sprintf(ansiCursorUpN, r.prevLines) + "\r"

		frameLines := strings.Split(frame, "\n")

		var rebuilt strings.Builder

		for i, line := range frameLines {
			rebuilt.WriteString(ansiClearLine)
			rebuilt.WriteString(line)

			if i < len(frameLines)-1 {
				rebuilt.WriteString("\n")
			}
		}

		frame = rebuilt.String()
	}

	output += frame + "\n"

	r.write(ansiSyncBegin + output + ansiSyncEnd)

	r.prevLines = physicalLines
}

// write writes a string to the renderer's writer, ignoring errors.
// Terminal output is best-effort: a broken pipe should not crash the render loop.
func (r *InlineRenderer) write(s string) {
	_, _ = fmt.Fprint(r.writer, s)
}

// writef writes a formatted string to the renderer's writer, ignoring errors.
func (r *InlineRenderer) writef(format string, args ...any) {
	_, _ = fmt.Fprintf(r.writer, format, args...)
}

// Finish clears the in-place frame and prints the final static tree.
// Call this once when the workflow completes. It stops the background
// refresh loop before rendering to avoid concurrent tree access.
func (r *InlineRenderer) Finish(workflowErr error) {
	r.Stop()

	tree := r.subscriber.GetDependencyTree()

	if r.prevLines > 0 {
		r.write(fmt.Sprintf(ansiCursorUpN, r.prevLines) + "\r")

		for range r.prevLines {
			r.write(ansiClearLine)
			r.write("\n")
		}

		r.write(fmt.Sprintf(ansiCursorUpN, r.prevLines))
		r.write("\r")
		r.prevLines = 0
	}

	if r.hideCursor {
		r.write(ansiShowCursor)
	}

	if tree != nil {
		final := tree.Render(0)
		if final != msgNoActivitiesToDisplay {
			r.write(final + "\n")
		}
	}

	elapsed := time.Since(r.startTime)

	elapsedStr := ""
	if !r.startTime.IsZero() {
		elapsedStr = " after " + FormatDuration(elapsed)
	}

	status := "completed successfully" + elapsedStr + "."
	colorCode := "\033[32m" // green

	if workflowErr != nil {
		status = "failed: " + workflowErr.Error() + elapsedStr
		colorCode = "\033[31m" // red
	}

	if r.noColor {
		r.writef("%s %s\n", r.appName, status)
	} else {
		r.writef("%s%s %s\033[0m\n", colorCode, r.appName, status)
	}
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

// renderAndNotify calls Render and, if a renderNotify channel is set,
// sends a non-blocking signal. This provides deterministic synchronization
// for tests without affecting production behavior.
func (r *InlineRenderer) renderAndNotify() {
	r.Render()

	if r.renderNotify != nil {
		select {
		case r.renderNotify <- struct{}{}:
		default:
		}
	}
}

// Stop terminates the background refresh goroutine started by Start.
// It is safe to call Stop even if Start was never called.
func (r *InlineRenderer) Stop() {
	r.tickMu.Lock()
	defer r.tickMu.Unlock()

	if r.cancelFn != nil {
		r.cancelFn()
		<-r.tickerDone
		r.cancelFn = nil
		r.tickerDone = nil
		r.refreshChan = nil
	}
}

// effectiveMaxHeight returns the maxHeight if set, otherwise detects
// terminal height from writer fd, falling back to 50.
func (r *InlineRenderer) effectiveMaxHeight() int {
	if r.maxHeight > 0 {
		return r.maxHeight
	}

	if f, ok := r.writer.(*os.File); ok {
		if _, height, err := term.GetSize(int(f.Fd())); err == nil && height > 4 {
			return height - 4
		}
	}

	return 50
}

// effectiveMaxWidth returns the terminal width if the writer is a terminal,
// otherwise a sensible default for non-terminal writers.
func (r *InlineRenderer) effectiveMaxWidth() int {
	return GetTerminalWidth(r.writer)
}

// renderSummary builds a one-line NOM-style summary bar.
func (r *InlineRenderer) renderSummary() string {
	running, completed, failed, pending := r.subscriber.GetActivityCounts()

	var parts []string

	if running > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolRunning, running))
	}

	if completed > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolCompleted, completed))
	}

	if failed > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolFailed, failed))
	}

	if pending > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolPaused, pending))
	}

	total := running + completed + failed + pending

	if !r.startTime.IsZero() {
		r.tickMu.RLock()
		startTime := r.startTime
		r.tickMu.RUnlock()
		elapsed := time.Since(startTime)
		parts = append(parts, fmt.Sprintf("%s%s", SymbolTiming, FormatDuration(elapsed)))
	}

	if len(parts) == 0 {
		return ""
	}

	summary := strings.Join(parts, " ") + fmt.Sprintf(" %s%d", SymbolTotal, total)

	if total > 0 {
		pct := (completed + failed) * 100 / total
		summary += fmt.Sprintf(" (%d%%)", pct)
	}

	visualWidth := ansi.StringWidth(summary)
	border := strings.Repeat("─", max(visualWidth+2, 3))

	return fmt.Sprintf("╭%s╮\n│ %s │\n╰%s╯", border, summary, border)
}
