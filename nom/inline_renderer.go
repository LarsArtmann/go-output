package nom

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	ansiCursorUp    = "\033[A"
	ansiCursorUpN   = "\033[%dA"
	ansiClearLine   = "\r\033[K"
	ansiClearScreen = "\033[2J\033[H"
	ansiHideCursor  = "\033[?25l"
	ansiShowCursor  = "\033[?25h"
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
func (r *InlineRenderer) SetHideCursor(hide bool) {
	r.hideCursor = hide
}

// SetNoColor forces colorless output. By default, colors are enabled unless
// NO_COLOR, TERM=dumb, or lacking a terminal is detected.
func (r *InlineRenderer) SetNoColor(noColor bool) {
	r.noColor = noColor
}

// SetAppName sets the application name for success/failure messages (default: "Workflow").
func (r *InlineRenderer) SetAppName(name string) {
	r.appName = name
}

// detectNoColor checks environment variables for color suppression.
func detectNoColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return true
	}

	if os.Getenv("TERM") == "dumb" {
		return true
	}

	if os.Getenv("CI") != "" {
		return true
	}

	return false
}

// SetStartTime sets the workflow start time for elapsed display.
func (r *InlineRenderer) SetStartTime(t time.Time) {
	r.startTime = t
}

// Render redraws the dependency tree in-place. On the first call it just prints.
// On subsequent calls it moves the cursor up to overwrite the previous frame.
func (r *InlineRenderer) Render() {
	if r.subscriber == nil {
		return
	}

	r.subscriber.UpdateRunningActivityElapsed()
	r.subscriber.SyncActivityTimingToTree()

	tree := r.subscriber.GetDependencyTree()
	if tree == nil {
		return
	}

	maxH := r.maxHeight
	if maxH <= 0 {
		maxH = 50
	}

	frame := tree.Render(maxH)

	if frame == msgNoActivitiesToDisplay {
		return
	}

	summary := r.renderSummary()
	if summary != "" {
		frame += "\n" + summary
	}

	lines := strings.Count(frame, "\n") + 1

	var output string

	if r.prevLines == 0 {
		if r.hideCursor {
			output = ansiHideCursor
		}
	} else {
		output = fmt.Sprintf(ansiCursorUpN, r.prevLines) + "\r"

		frameLines := strings.Split(frame, "\n")
		var rebuilt strings.Builder

		for i, line := range frameLines {
			rebuilt.WriteString(ansiClearLine + line)
			if i < len(frameLines)-1 {
				rebuilt.WriteString("\n")
			}
		}

		frame = rebuilt.String()
	}

	output += frame + "\n"

	fmt.Fprint(r.writer, output)

	r.prevLines = lines
}

// Finish clears the in-place frame and prints the final static tree.
// Call this once when the workflow completes.
func (r *InlineRenderer) Finish(workflowErr error) {
	tree := r.subscriber.GetDependencyTree()

	if r.prevLines > 0 {
		fmt.Fprint(r.writer, fmt.Sprintf(ansiCursorUpN, r.prevLines)+"\r")

		for range r.prevLines {
			fmt.Fprint(r.writer, ansiClearLine+"\n")
		}

		fmt.Fprint(r.writer, fmt.Sprintf(ansiCursorUpN, r.prevLines)+"\r")
		r.prevLines = 0
	}

	if r.hideCursor {
		fmt.Fprint(r.writer, ansiShowCursor)
	}

	if tree != nil {
		final := tree.Render(0)
		if final != msgNoActivitiesToDisplay {
			fmt.Fprintln(r.writer, final)
		}
	}

	if r.noColor {
		if workflowErr != nil {
			fmt.Fprintf(r.writer, "%s failed: %v\n", r.appName, workflowErr)
		} else {
			fmt.Fprintf(r.writer, "%s completed successfully.\n", r.appName)
		}
	} else {
		if workflowErr != nil {
			fmt.Fprintf(r.writer, "\033[31m%s failed: %v\033[0m\n", r.appName, workflowErr)
		} else {
			fmt.Fprintf(r.writer, "\033[32m%s completed successfully.\033[0m\n", r.appName)
		}
	}
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
		elapsed := time.Since(r.startTime)
		parts = append(parts, fmt.Sprintf("%s%s", SymbolTiming, FormatDuration(elapsed)))
	}

	if len(parts) == 0 {
		return ""
	}

	summary := strings.Join(parts, " ") + fmt.Sprintf(" %s%d", SymbolTotal, total)

	border := strings.Repeat("─", max(len(summary)+2, 3))

	return fmt.Sprintf("╭%s╮\n│ %s │\n╰%s╯", border, summary, border)
}
