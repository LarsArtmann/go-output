package nom

import (
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"golang.org/x/term"
)

// Finish clears the in-place frame and prints the final static tree.
// Call this once when the workflow completes. It stops the background
// refresh loop before rendering to avoid concurrent tree access.
func (r *InlineRenderer) Finish(workflowErr error) {
	r.Stop()

	r.renderMu.Lock()
	defer r.renderMu.Unlock()

	// Snapshot config under tickMu.RLock so SetStartTime/SetAppName/SetNoColor/
	// SetHideCursor can't race the reads below (they all take tickMu.Lock).
	r.tickMu.RLock()
	startTime := r.startTime
	appName := r.appName
	noColor := r.noColor
	hideCursor := r.hideCursor
	r.tickMu.RUnlock()

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

	if hideCursor {
		r.write(ansiShowCursor)
	}

	// Render under the subscriber's read lock (same reason as Draw).
	if final, ok := r.subscriber.RenderTree(0, 0); ok && final != msgNoActivitiesToDisplay {
		r.write(final + "\n")
	}

	elapsed := time.Since(startTime)

	elapsedStr := ""
	if !startTime.IsZero() {
		elapsedStr = " after " + FormatDuration(elapsed)
	}

	status := "completed successfully" + elapsedStr + "."
	statusColor := Colors.Completed

	if workflowErr != nil {
		status = "failed: " + workflowErr.Error() + elapsedStr
		statusColor = Colors.Failed
	}

	line := fmt.Sprintf("%s %s", appName, status)
	if noColor {
		r.write(line + "\n")
	} else {
		r.write(lipgloss.NewStyle().Foreground(statusColor).Render(line) + "\n")
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
	counts := r.subscriber.GetActivityCounts()

	var parts []string

	if counts.Running > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolRunning, counts.Running))
	}

	if counts.Completed > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolCompleted, counts.Completed))
	}

	if counts.Failed > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolFailed, counts.Failed))
	}

	if counts.Pending > 0 {
		parts = append(parts, fmt.Sprintf("%s%d", SymbolPaused, counts.Pending))
	}

	total := counts.Total()

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
		pct := (counts.Completed + counts.Failed) * 100 / total
		summary += fmt.Sprintf(" (%d%%)", pct)
	}

	visualWidth := ansi.StringWidth(summary)
	border := strings.Repeat("─", max(visualWidth+2, 3))

	return fmt.Sprintf("╭%s╮\n│ %s │\n╰%s╯", border, summary, border)
}
