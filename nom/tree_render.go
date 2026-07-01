package nom

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// RenderWithSnapshots generates the tree rendering using immutable activity
// snapshots instead of reading the shared *Activity pointer. This is the
// canonical render path: snapshots are taken under the subscriber's read lock
// once, then the tree walk reads only immutable data. A nil snapshots map is
// treated as "all activities pending with blank labels" — useful for tests
// that build trees structurally without a subscriber.
func (dt *DependencyTree) RenderWithSnapshots(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight, maxWidth int,
) string {
	if err := dt.ensureBuilt(); err != nil {
		return fmt.Sprintf("Error building tree: %v", err)
	}

	visible := dt.collectVisibleNodes(snapshots, maxHeight)

	if len(visible) == 0 {
		return msgNoActivitiesToDisplay
	}

	var lines []string

	for _, entry := range visible {
		lines = append(lines, dt.renderLine(entry, snapshots, maxWidth))
	}

	return strings.Join(lines, "\n")
}

func (dt *DependencyTree) collectVisibleNodes(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight int,
) []VisibleEntry {
	dt.mu.RLock()
	defer dt.mu.RUnlock()

	if len(dt.roots) == 0 {
		return nil
	}

	if maxHeight <= 0 {
		maxHeight = len(dt.nodes)
	}

	var visible []VisibleEntry

	for _, root := range dt.roots {
		dt.walkSubtree(root, "", true, true, &visible, snapshots, maxHeight)

		if len(visible) >= maxHeight {
			break
		}
	}

	return visible
}

func (dt *DependencyTree) elideCompletedUnderPressure(
	children []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight, visibleCount int,
) (active []*ActivityNode, collapsedCompleted int) {
	if maxHeight <= 0 {
		return children, 0
	}

	remaining := maxHeight - visibleCount
	if remaining >= len(children) {
		return children, 0
	}

	for _, child := range children {
		snap := lookupSnapshot(snapshots, child.ID)
		if snap.Status == ActivityStatusCompleted {
			collapsedCompleted++
			continue
		}

		active = append(active, child)
	}

	if collapsedCompleted == 0 {
		return children, 0
	}

	return active, collapsedCompleted
}

func (dt *DependencyTree) walkSubtree( //nolint:cyclop // DFS traversal with phase collapse — inherently branchy
	node *ActivityNode,
	prefix string,
	isLastSibling bool,
	isRoot bool,
	visible *[]VisibleEntry,
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight int,
) {
	if len(*visible) >= maxHeight {
		return
	}

	entry := VisibleEntry{
		Node:      node,
		Prefix:    prefix,
		Connector: "",
		IsRoot:    isRoot,
	}

	if !isRoot {
		if isLastSibling {
			entry.Connector = "└── "
		} else {
			entry.Connector = "├── "
		}
	}

	*visible = append(*visible, entry)

	// Phase-aware collapse: if this is a phase node and ALL direct children
	// are in terminal state (completed/failed), collapse the subtree to a
	// single summary line. This turns "Code Formatting" with 6 green
	// checkmarks into "◈ Code Formatting  6/6 · 4.1s".
	snap := lookupSnapshot(snapshots, node.ID)
	if snap.IsPhase() && dt.collapseCompletedPhases && len(node.Children) > 0 {
		if pc, ok := computePhaseCounts(snapshots, node.Children); ok {
			(*visible)[len(*visible)-1].PhaseCounts = &pc
			return
		}
	}

	children := dt.childPriority(node, snapshots)
	if len(children) == 0 {
		return
	}

	children, collapsedDone := dt.elideCompletedUnderPressure(children, snapshots, maxHeight, len(*visible))

	var childIndent string

	if isRoot {
		childIndent = ""
	} else if isLastSibling {
		childIndent = prefix + "    "
	} else {
		childIndent = prefix + "│   "
	}

	for i, child := range children {
		if len(*visible) >= maxHeight {
			return
		}

		dt.walkSubtree(
			child,
			childIndent,
			i == len(children)-1 && collapsedDone == 0,
			false,
			visible,
			snapshots,
			maxHeight,
		)
	}

	// If completed children were elided under height pressure, surface a faint
	// "⋯ N completed" marker so the collapsed work is visible, not silently gone.
	if collapsedDone > 0 && len(*visible) < maxHeight {
		appendCollapseMarker(visible, childIndent, collapsedDone, len(children) == 0)
	}
}

// appendCollapseMarker adds a synthetic "⋯ N completed" entry to the visible
// list when completed children were elided under height pressure.
func appendCollapseMarker(visible *[]VisibleEntry, indent string, collapsedDone int, noRemainingChildren bool) {
	connector := "├── "
	if noRemainingChildren {
		connector = "└── "
	}

	*visible = append(*visible, VisibleEntry{
		CollapsedCompleted: collapsedDone,
		CollapseIndent:     indent,
		Connector:          connector,
	})
}

// formatActivityLabel builds the core display string for a single activity
// snapshot: phase-aware symbol + label + timing info. Returns the unstyled
// display and the status-derived color for the caller to apply.
func formatActivityLabel(snap ActivitySnapshot) (display string, c color.Color) {
	symbol := snap.Symbol
	c = snap.Color

	if snap.IsPhase() {
		symbol = SymbolPhase
		c = Colors.Phase
	}

	display = fmt.Sprintf("%s %s", symbol, snap.Label)

	timingInfo := FormatActivityNodeTiming(
		snap.Status,
		snap.CurrentElapsed,
		snap.EstimatedTime,
	)

	if timingInfo != "" {
		// Slow-step escalation: dim yellow for >10s, dim red for >30s.
		// This surfaces performance bottlenecks without the user scanning every line.
		timingStyle := slowStepStyle(snap)

		// An unmodified lipgloss style renders as plain text (no ANSI codes),
		// so we can always use Render without a separate "is styled?" check.
		display += " " + timingStyle.Render(timingInfo)
	}

	// Optional host tag (dormant unless the event carried one).
	if snap.Host != "" {
		display += " @" + snap.Host
	}

	// Retry suffix: show ⟳N when the activity has been retried. Placed after
	// timing/host so the retry badge reads as a status annotation. An optional
	// reason (e.g. "timeout") renders as "⟳2 (timeout)".
	if snap.RetryCount > 0 {
		suffix := fmt.Sprintf(" %s%d", SymbolRetrying, snap.RetryCount)
		if snap.RetryReason != "" {
			suffix += fmt.Sprintf(" (%s)", snap.RetryReason)
		}

		display += lipgloss.NewStyle().Foreground(Colors.Info).Render(suffix)
	}

	// Optional download progress bar — only while the activity is actively
	// running; a completed/failed download no longer needs a live bar.
	if snap.Status == ActivityStatusRunning && snap.Download.HasDownload() {
		display += " " + formatDownloadBar(snap.Download, downloadBarWidth)
	}

	// Sub-step progress message (e.g. "Tidying module [2/26]"). Rendered as a
	// dim sub-line beneath the activity label — mirrors how nom shows per-
	// derivation download/build progress inline.
	if snap.Progress != "" && snap.Status == ActivityStatusRunning {
		display += "\n" + lipgloss.NewStyle().
			Faint(true).
			Render(fmt.Sprintf("  %s %s", SymbolProgress, snap.Progress))
	}

	return display, c
}

const downloadBarWidth = 10

// Slow-step color escalation thresholds for timing display.
// Steps exceeding these durations get dim yellow/red timing to surface
// performance bottlenecks without requiring the user to scan every line.
const (
	slowThresholdYellow = 10 * time.Second
	slowThresholdRed    = 30 * time.Second
)

// formatDownloadBar renders a compact NOM-style byte-progress bar like
// "▕████░░░░▏ 45%". When the total is unknown it shows transferred bytes only.
func formatDownloadBar(d DownloadProgress, width int) string {
	if width < 4 {
		width = 4
	}

	if d.Total <= 0 {
		return fmt.Sprintf("%s %s", SymbolDownload, formatBytes(d.Downloaded))
	}

	filled := min(int(d.Fraction()*float64(width)), width)

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	pct := int(d.Fraction() * 100)

	return fmt.Sprintf("▕%s▏ %d%%", bar, pct)
}

// slowStepStyle returns the lipgloss style for the slow-step escalation.
// Returns an empty style (plain text) for fast steps, faint yellow for
// >10s, and faint red for >30s. Considers both elapsed time (for running/
// completed steps) and estimated time (for pending steps).
func slowStepStyle(snap ActivitySnapshot) lipgloss.Style {
	d := snap.CurrentElapsed
	if snap.Status == ActivityStatusPending {
		d = snap.EstimatedTime
	}

	switch {
	case d >= slowThresholdRed:
		return lipgloss.NewStyle().Foreground(Colors.Failed).Faint(true)
	case d >= slowThresholdYellow:
		return lipgloss.NewStyle().Foreground(Colors.Running).Faint(true)
	default:
		return lipgloss.NewStyle()
	}
}

// formatBytes renders a byte count in a human-readable binary form (KiB, MiB…).
func formatBytes(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1fGiB", float64(b)/float64(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1fMiB", float64(b)/float64(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1fKiB", float64(b)/float64(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}

// computePhaseCounts checks if all direct children of a phase are in terminal
// state. Returns the aggregate counts and true if collapsible; returns false
// if any child is pending or running.
func computePhaseCounts(snapshots map[ActivityID]ActivitySnapshot, children []*ActivityNode) (PhaseCounts, bool) {
	var pc PhaseCounts

	for _, child := range children {
		snap := lookupSnapshot(snapshots, child.ID)

		switch snap.Status { //nolint:exhaustive // default handles pending/running by returning false
		case ActivityStatusCompleted:
			pc.Completed++
		case ActivityStatusFailed:
			pc.Failed++
		default:
			return PhaseCounts{}, false
		}

		if snap.CurrentElapsed > pc.MaxElapsed {
			pc.MaxElapsed = snap.CurrentElapsed
		}
	}

	return pc, true
}

// formatCollapsedPhaseLabel builds the display string for a collapsed phase:
// "◈ Code Formatting  6/6 · 4.1s". Uses the phase symbol and phase color,
// or the failed color if any children failed.
func formatCollapsedPhaseLabel(snap ActivitySnapshot, pc PhaseCounts) (display string, c color.Color) {
	symbol := SymbolPhase
	c = Colors.Phase

	if pc.Failed > 0 {
		c = Colors.Failed
	}

	display = fmt.Sprintf("%s %s  %d/%d", symbol, snap.Label, pc.Completed, pc.Total())

	if pc.MaxElapsed > 0 {
		display += " · " + FormatDuration(pc.MaxElapsed)
	}

	return display, c
}

func (dt *DependencyTree) renderLine(
	entry VisibleEntry,
	snapshots map[ActivityID]ActivitySnapshot,
	maxWidth int,
) string {
	// Synthetic collapse marker: rendered when completed children were elided
	// under height pressure. Shown faint so it reads as "hidden, not gone".
	if entry.CollapsedCompleted > 0 {
		marker := fmt.Sprintf("%s%s ⋯ %d completed", entry.CollapseIndent, entry.Connector, entry.CollapsedCompleted)
		rendered := lipgloss.NewStyle().Faint(true).Render(marker)

		if maxWidth > 0 && VisibleWidth(rendered) > maxWidth {
			rendered = TruncateVisible(rendered, maxWidth)
		}

		return rendered
	}

	node := entry.Node
	snap := lookupSnapshot(snapshots, node.ID)

	var activityDisplay string

	var color color.Color

	if entry.PhaseCounts != nil {
		activityDisplay, color = formatCollapsedPhaseLabel(snap, *entry.PhaseCounts)
	} else {
		activityDisplay, color = formatActivityLabel(snap)
	}

	if len(node.SecondaryParents) > 0 {
		depNames := make([]string, 0, len(node.SecondaryParents))

		for _, depID := range node.SecondaryParents {
			depSnap := lookupSnapshot(snapshots, depID)

			name := depSnap.Label
			if name == "" {
				name = depID.String()
			}

			depNames = append(depNames, name)
		}

		activityDisplay += lipgloss.NewStyle().
			Foreground(Colors.Info).
			Render(" ← " + strings.Join(depNames, ", "))
	}

	fullPrefix := entry.Prefix + entry.Connector

	if maxWidth > 0 {
		available := maxWidth - ansi.StringWidth(fullPrefix)
		activityDisplay = TruncateVisible(activityDisplay, available)
	}

	style := activityNodeStyle(color)
	// Roots anchor the tree — render them bold so the top-level activities
	// stand out from their twigs/leaves (NOM root/twig/leaf styling).
	if node.Class() == NodeClassRoot {
		style = style.Bold(true)
	}

	rendered := style.Render(fullPrefix + activityDisplay)

	if maxWidth > 0 && VisibleWidth(rendered) > maxWidth {
		rendered = TruncateVisible(rendered, maxWidth)
	}

	return rendered
}

// VisibleNodesWithSnapshots returns the ordered list of real tree nodes that
// would be displayed for the given maxHeight, in priority order. Uses snapshots
// for status-based sorting.
//
// Only REAL activity nodes are returned — synthetic collapse-marker lines
// (produced under height pressure when completed children are elided) are
// skipped, so this slice NEVER contains a nil entry. Callers that need the
// markers too (e.g. to render them) must use VisibleEntriesWithSnapshots.
func (dt *DependencyTree) VisibleNodesWithSnapshots(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight int,
) []*ActivityNode {
	entries := dt.VisibleEntriesWithSnapshots(snapshots, maxHeight)

	nodes := make([]*ActivityNode, 0, len(entries))

	for _, entry := range entries {
		if entry.Node != nil {
			nodes = append(nodes, entry.Node)
		}
	}

	return nodes
}

// VisibleEntry is a single renderable line of the dependency tree. Exactly one
// variant is meaningful:
//
//   - A real activity line: Node != nil and CollapsedCompleted == 0.
//   - A synthetic collapse marker, rendered when completed children are elided
//     under height pressure: Node == nil and CollapsedCompleted > 0.
//
// Exposing the marker explicitly (instead of smuggling a nil into a
// []*ActivityNode) makes the "node or marker" choice representable without nil
// dereferences at the call site.
type VisibleEntry struct {
	Node *ActivityNode

	Prefix    string
	Connector string
	IsRoot    bool

	// CollapsedCompleted > 0 marks a synthetic "⋯ N completed" line; Node is
	// nil in that case.
	CollapsedCompleted int
	CollapseIndent     string

	// PhaseCounts, when non-nil, indicates this entry is a collapsed phase
	// (all direct children are in terminal state). The children are not
	// expanded; instead the summary counts are rendered inline on the phase
	// label. This is the phase-aware collapse: a completed category like
	// "Code Formatting" shows as "◈ Code Formatting  6/6 · 4.1s" instead
	// of expanding all 6 child steps.
	PhaseCounts *PhaseCounts
}

// PhaseCounts holds aggregate status counts for a collapsed phase's children.
// MaxElapsed tracks the longest child duration (not the sum), since DAG steps
// run in parallel — summing would over-report wall-clock time.
type PhaseCounts struct {
	Completed  int
	Failed     int
	MaxElapsed time.Duration
}

// Total returns the total number of children accounted for.
func (pc PhaseCounts) Total() int { return pc.Completed + pc.Failed }

// VisibleEntriesWithSnapshots returns the renderable tree lines (real nodes AND
// collapse markers) in display order, capped at maxHeight. This is the
// marker-aware variant of VisibleNodesWithSnapshots and is what renderers that
// show the "⋯ N completed" line (e.g. the bubbletea TUI) must use.
func (dt *DependencyTree) VisibleEntriesWithSnapshots(
	snapshots map[ActivityID]ActivitySnapshot,
	maxHeight int,
) []VisibleEntry {
	if err := dt.ensureBuilt(); err != nil {
		return nil
	}

	return dt.collectVisibleNodes(snapshots, maxHeight)
}

// ensureBuilt triggers a Build() if the tree has not yet been built.
// Splitting this off lets renderers share the "build if needed, then render
// under the read lock" preamble without duplicating the lock dance.
func (dt *DependencyTree) ensureBuilt() error {
	dt.mu.RLock()
	needsBuild := !dt.loaded
	dt.mu.RUnlock()

	if !needsBuild {
		return nil
	}

	return dt.Build()
}

// RenderVisibleEntry renders a single visible entry — either a real activity
// node or a synthetic collapse marker — using immutable snapshot data. This is
// the marker-aware primitive renderers should call per line (it mirrors the
// inline renderer's renderLine, which always handled markers correctly).
func (dt *DependencyTree) RenderVisibleEntry(
	entry VisibleEntry,
	snapshots map[ActivityID]ActivitySnapshot,
	maxWidth int,
) string {
	return dt.renderLine(entry, snapshots, maxWidth)
}

// RenderNode renders a single node for external consumers (e.g., TUI mouse
// click highlight). Uses the snapshot for label/color/symbol.
//
// Returns "" for a nil node so a stray collapse-marker can never panic callers
// that still use the node-only API. Prefer RenderVisibleEntry for code that
// needs to render markers.
func (dt *DependencyTree) RenderNode(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) string {
	if node == nil {
		return ""
	}

	snap := lookupSnapshot(snapshots, node.ID)
	display, color := formatActivityLabel(snap)

	return activityNodeStyle(color).Render(display)
}

func activityNodeStyle(color color.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(color).
		Width(0).
		Inline(true)
}

// lookupSnapshot returns the snapshot for id, or a blank pending snapshot if
// the activity hasn't been registered yet (e.g. a structural placeholder for
// a dependency that hasn't received its activity.started event).
func lookupSnapshot(snapshots map[ActivityID]ActivitySnapshot, id ActivityID) ActivitySnapshot {
	if snapshots == nil {
		return blankActivitySnapshot()
	}

	if snap, ok := snapshots[id]; ok {
		return snap
	}

	return blankActivitySnapshot()
}

// blankActivitySnapshot returns the default snapshot for unregistered activities:
// pending status, empty label, pending symbol/color. Kept as a function (not a
// global var) so the snapshot is never accidentally mutated by callers.
func blankActivitySnapshot() ActivitySnapshot {
	return ActivitySnapshot{
		Label:  "",
		Status: ActivityStatusPending,
		Symbol: SymbolPending,
		Color:  Colors.Pending,
	}
}
