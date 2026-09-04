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
//
// Finish intentionally does NOT take or render the workflow error: the
// failed activity's ⚠ symbol and error annotation in the final tree
// already visualize failure, and the calling application owns the richer
// post-run summary (typically stderr + exit code). A nil subscriber makes
// Finish a no-op, matching Draw.
func (r *InlineRenderer) Finish() {
	if r.subscriber == nil {
		return
	}

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
	r.lastTreeFrame = ""
	r.lastStableTreeFrame = ""

	// Render from immutable snapshot (same race-free path as Draw).
	if final, ok := r.subscriber.RenderSnapshot(0, 0); ok && final != MsgNoActivities {
		r.write(final + "\n")
	}
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

// styleGroup appends text to groups, applying faint styling when color is enabled.
func styleGroup(groups []string, text string, noColor bool, faint lipgloss.Style) []string {
	if noColor {
		return append(groups, text)
	}

	return append(groups, faint.Render(text))
}

// renderSummary builds a one-line NOM-style summary bar. It takes the
// already-snapshotted startTime so it does not re-acquire tickMu (Draw, its
// only caller, already holds renderMu and snapshotted the full config).
// When noColor is false, counts are colored + bold, secondary metrics are
// faint, groups are separated by dim │, and the border box is dimmed.
func (r *InlineRenderer) renderSummary(startTime time.Time, noColor bool) string {
	counts := r.subscriber.GetActivityCounts()

	faint := lipgloss.NewStyle().Faint(true)
	bold := lipgloss.NewStyle().Bold(true)

	var groups []string

	// PRIMARY: activity counts (colored + bold when color enabled).
	if noColor {
		if s := counts.Summary(); s != "" {
			groups = append(groups, s)
		}
	} else {
		if s := counts.SummaryColored(r.subscriber.GetThemeColors()); s != "" {
			groups = append(groups, bold.Render(s))
		}
	}

	// SECONDARY: elapsed time.
	if !startTime.IsZero() {
		groups = styleGroup(groups, FormatDuration(time.Since(startTime)), noColor, faint)
	}

	// SECONDARY: optional segments (ETA, critical path, parallelism, DAG summary).
	for _, seg := range r.optionalSummarySegments(startTime) {
		groups = styleGroup(groups, seg, noColor, faint)
	}

	if len(groups) == 0 {
		return ""
	}

	// SECONDARY: total count + completion percent.
	totalStr := fmt.Sprintf("%s%d", SymbolTotal, counts.Total())
	if counts.Total() > 0 {
		totalStr += fmt.Sprintf(" (%d%%)", counts.CompletionPercent())
	}

	groups = styleGroup(groups, totalStr, noColor, faint)

	// Join groups — dim │ separator between groups when colored.
	var summary string
	if noColor {
		summary = strings.Join(groups, " ")
	} else {
		summary = strings.Join(groups, " "+faint.Render("│")+" ")
	}

	visualWidth := ansi.StringWidth(summary)
	border := strings.Repeat("─", max(visualWidth+2, 3))

	if noColor {
		return fmt.Sprintf("╭%s╮\n│ %s │\n╰%s╯", border, summary, border)
	}

	// Colored mode: dim border characters, summary keeps its own styling.
	topLine := faint.Render("╭" + border + "╮")
	sideBar := faint.Render("│")
	botLine := faint.Render("╰" + border + "╯")

	return topLine + "\n" + sideBar + " " + summary + " " + sideBar + "\n" + botLine
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
