package nom

import (
	"fmt"
	"os"
	"strings"
	"time"

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

	// Snapshot all config once so SetStartTime/SetAppName/SetNoColor/SetHideCursor
	// can't race the reads below.
	cfg := r.snapshotConfig()

	if r.prevLines > 0 {
		r.write(ansi.CursorUp(r.prevLines) + "\r")

		for range r.prevLines {
			r.write(ansiClearLine)
			r.write("\n")
		}

		r.write(ansi.CursorUp(r.prevLines))
		r.write("\r")
		r.prevLines = 0
	}

	// Drain any pending external log lines so they appear above the final
	// tree, not lost when the frame is cleared.
	if len(r.pendingLines) > 0 {
		for _, line := range r.pendingLines {
			r.write(line)
			r.write("\n")
		}

		r.pendingLines = nil
	}

	if cfg.hideCursor {
		r.write(ansi.ShowCursor)
	}

	// Reset frame cache so a new workflow can render from scratch.
	r.lastFrame = ""

	// Render from immutable snapshot (same race-free path as Draw).
	if final, ok := r.subscriber.RenderSnapshot(0, 0); ok && final != MsgNoActivities {
		r.write(final + "\n")
	}

	// The completion line (e.g. "BuildFlow completed successfully after 1m15s.")
	// is intentionally NOT printed here — the calling application provides its
	// own post-run summary which is more detailed (auto-fixes, artifacts, etc.).
	// Printing both would produce duplicate/mangled output on the terminal.
}

// effectiveMaxHeight returns the given maxHeight if set, otherwise detects
// terminal height from writer fd, falling back to 50.
func (r *InlineRenderer) effectiveMaxHeight(maxHeight int) int {
	if maxHeight > 0 {
		return maxHeight
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

// renderSummary builds a one-line NOM-style summary bar. It takes the
// already-snapshotted startTime so it does not re-acquire tickMu (Draw, its
// only caller, already holds renderMu and snapshotted the full config).
func (r *InlineRenderer) renderSummary(startTime time.Time) string {
	counts := r.subscriber.GetActivityCounts()

	var parts []string

	countsStr := counts.Summary()
	if countsStr != "" {
		parts = append(parts, countsStr)
	}

	if !startTime.IsZero() {
		elapsed := time.Since(startTime)
		parts = append(parts, FormatDuration(elapsed))
	}

	parts = append(parts, r.optionalSummarySegments(startTime)...)

	if len(parts) == 0 {
		return ""
	}

	summary := strings.Join(parts, " ") + fmt.Sprintf(" %s%d", SymbolTotal, counts.Total())

	if counts.Total() > 0 {
		summary += fmt.Sprintf(" (%d%%)", counts.CompletionPercent())
	}

	visualWidth := ansi.StringWidth(summary)
	border := strings.Repeat("─", max(visualWidth+2, 3))

	return fmt.Sprintf("╭%s╮\n│ %s │\n╰%s╯", border, summary, border)
}

// optionalSummarySegments returns the optional summary bar segments: estimated
// remaining time, critical-path ETA, parallelism meter, and DAG structural
// summary. Each segment is only included when its toggle is on and the value
// is non-zero/non-empty. Invoked under renderMu.
func (r *InlineRenderer) optionalSummarySegments(_ time.Time) []string {
	var segments []string

	if r.estimatedRemaining != nil {
		if remaining := r.estimatedRemaining(); remaining > 0 {
			segments = append(segments, "~"+FormatDuration(remaining)+" left")
		}
	}

	if r.showCriticalPathETA {
		if remaining := r.criticalPathRemaining(); remaining > 0 {
			segments = append(segments, "~"+FormatDuration(remaining)+" critical")
		}
	}

	if r.subscriber.showParallelism {
		parallelism := r.subscriber.ParallelismStats()
		if parallelism.Possible > 0 {
			segments = append(segments, parallelism.String())
		}
	}

	if r.showDAGSummary {
		if dagStr := r.dagSummarySegment(); dagStr != "" {
			segments = append(segments, dagStr)
		}
	}

	return segments
}

// criticalPathRemaining returns the longest remaining-time path through the DAG,
// or 0 when there is no dependency tree. Invoked under renderMu.
func (r *InlineRenderer) criticalPathRemaining() time.Duration {
	tree := r.subscriber.DependencyTree()
	if tree == nil {
		return 0
	}

	snapshots := r.subscriber.SnapshotActivities()

	return tree.EstimatedCriticalPathRemaining(snapshots)
}

// dagSummarySegment returns the structural DAG summary string for the summary
// bar (e.g. "4 nodes · 4 edges · 3 layers"). Returns empty when there is no
// dependency tree. Invoked under renderMu.
func (r *InlineRenderer) dagSummarySegment() string {
	tree := r.subscriber.DependencyTree()
	if tree == nil {
		return ""
	}

	snapshots := r.subscriber.SnapshotActivities()

	return tree.DAGSummaryWithSnapshots(snapshots).String()
}
