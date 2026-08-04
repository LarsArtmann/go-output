package nom

import (
	"fmt"
	"image/color"
	"sort"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// minChildrenForPartialCollapse is the minimum number of children a phase must
// have before the partial-collapse-under-pressure logic kicks in. Phases with
// fewer children are always expanded — they fit comfortably in any viewport.
const minChildrenForPartialCollapse = 4

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

	if dt.renderMode == RenderModeLayered {
		return dt.renderLayered(snapshots, maxHeight, maxWidth)
	}

	visible := dt.collectVisibleNodes(snapshots, maxHeight)

	if len(visible) == 0 {
		return MsgNoActivities
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

	var criticalPath map[ActivityID]bool
	if dt.showCriticalPath {
		criticalPath = dt.computeCriticalPath(snapshots)
	}

	var visible []VisibleEntry

	sortedRoots := dt.sortRootsByPriority(dt.roots, snapshots, criticalPath)

	for _, root := range sortedRoots {
		dt.walkSubtree(root, "", true, true, &visible, snapshots, maxHeight, criticalPath)

		if len(visible) >= maxHeight {
			break
		}
	}

	return visible
}

// sortRootsByPriority sorts root nodes by display priority using the same
// sortKey logic as childPriority: failed > running > pending > completed,
// with longer-elapsed steps surfacing first within each status group. This
// ensures running steps are always at the top of the tree and visible even
// when many completed roots would otherwise consume the viewport budget.
func (dt *DependencyTree) sortRootsByPriority(
	roots []*ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
	criticalPath map[ActivityID]bool,
) []*ActivityNode {
	if len(roots) <= 1 {
		return roots
	}

	sorted := make([]*ActivityNode, len(roots))
	copy(sorted, roots)

	sort.SliceStable(sorted, func(i, j int) bool {
		ki := sortKeyForNode(sorted[i], snapshots, criticalPath)
		kj := sortKeyForNode(sorted[j], snapshots, criticalPath)
		return ki.less(kj)
	})

	return sorted
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
	criticalPath map[ActivityID]bool,
) {
	if len(*visible) >= maxHeight {
		return
	}

	entry := VisibleEntry{
		Node:           node,
		Prefix:         prefix,
		Connector:      "",
		IsRoot:         isRoot,
		OnCriticalPath: dt.showCriticalPath && criticalPath[node.ID],
		Convergence:    dt.showConvergence && len(node.Deps) > 1,
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

	// Partial phase collapse under height pressure: when a phase has
	// running/pending children but the remaining viewport can't fit them all,
	// show a summary line with progress counts instead of expanding every
	// child. This prevents one large fan-out phase (e.g. golangci-lint with
	// 30 modules) from monopolizing the viewport, keeping other concurrent
	// work visible. Independent of collapseCompletedPhases — driven purely
	// by viewport pressure.
	if snap.IsPhase() && len(node.Children) >= minChildrenForPartialCollapse {
		remaining := maxHeight - len(*visible)
		nonTerminal := countNonTerminalChildren(snapshots, node.Children)
		if nonTerminal > remaining {
			pc := computePartialPhaseCounts(snapshots, node.Children)
			(*visible)[len(*visible)-1].PhaseCounts = &pc
			return
		}
	}

	children := dt.childPriority(node, snapshots, criticalPath)
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
			criticalPath,
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
	return formatActivityLabelWithOptions(snap, ThemeDefault, labelOptions{})
}

// labelOptions carries optional rendering modifiers for an activity label.
// Zero value means "no extra markers" and is used by the backward-compatible
// formatActivityLabel helper and by RenderNode.
type labelOptions struct {
	// OnCriticalPath renders a ◆ prefix to highlight the longest estimated-time
	// path through the DAG.
	OnCriticalPath bool
	// Convergence renders a ◇ prefix for nodes with multiple incoming
	// dependencies (fan-in points).
	Convergence bool
	// ShowCategory renders a dim [tag] prefix when the snapshot has a
	// non-empty Category.
	ShowCategory bool
}

// formatActivityLabelWithOptions is the configurable version of
// formatActivityLabel. It supports critical-path and convergence markers while
// preserving the original symbol/label/timing layout.
//
//nolint:cyclop // mirrors original formatActivityLabel logic plus two optional markers
func formatActivityLabelWithOptions(
	snap ActivitySnapshot,
	theme Theme,
	opts labelOptions,
) (display string, c color.Color) { //nolint:cyclop // mirrors original formatActivityLabel logic plus two optional markers
	symbol := snap.Symbol
	c = snap.Color

	if snap.IsPhase() {
		symbol = SymbolPhase
		c = theme.Colors.Phase
	}

	var markers []string

	// Category tag: a dim [tag] prefix when enabled and the snapshot carries
	// a non-empty Category.
	if opts.ShowCategory && snap.Category != "" {
		markers = append(markers, lipgloss.NewStyle().Faint(true).Render("["+string(snap.Category)+"]"))

		// Apply category tint: override the status color with the theme's
		// category color when one is defined. This makes all activities in
		// the same category visually group by color, not just by [tag] prefix.
		if tinted := theme.CategoryColor(snap.Category); tinted != nil {
			c = tinted
		}
	}

	if opts.Convergence {
		markers = append(markers, string(SymbolConvergence))
	}

	if opts.OnCriticalPath {
		markers = append(markers, string(SymbolCritical))
	}

	if len(markers) > 0 {
		display = strings.Join(markers, " ") + " "
	}

	display += fmt.Sprintf("%s %s", symbol, snap.Label)

	timingInfo := FormatActivityNodeTiming(
		snap.Status,
		snap.CurrentElapsed,
		snap.EstimatedTime,
	)

	if timingInfo != "" {
		// Slow-step escalation: dim yellow for >10s, dim red for >30s.
		// This surfaces performance bottlenecks without the user scanning every line.
		timingStyle := slowStepStyle(snap, theme)

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

		display += lipgloss.NewStyle().Foreground(theme.Colors.Fallback).Render(suffix)
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
func slowStepStyle(snap ActivitySnapshot, theme Theme) lipgloss.Style {
	d := snap.CurrentElapsed
	if snap.Status == ActivityStatusPending {
		d = snap.EstimatedTime
	}

	switch {
	case d >= slowThresholdRed:
		return lipgloss.NewStyle().Foreground(theme.Colors.Failed).Faint(true)
	case d >= slowThresholdYellow:
		return lipgloss.NewStyle().Foreground(theme.Colors.Running).Faint(true)
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

// computePartialPhaseCounts computes aggregate counts for a phase's children
// regardless of their status. Unlike computePhaseCounts, this includes running
// and pending children — used for partial collapse under viewport pressure.
func computePartialPhaseCounts(
	snapshots map[ActivityID]ActivitySnapshot,
	children []*ActivityNode,
) PhaseCounts {
	var pc PhaseCounts

	for _, child := range children {
		snap := lookupSnapshot(snapshots, child.ID)

		switch snap.Status {
		case ActivityStatusCompleted:
			pc.Completed++
		case ActivityStatusFailed:
			pc.Failed++
		case ActivityStatusRunning:
			pc.Running++
		default:
			pc.Pending++
		}

		if snap.CurrentElapsed > pc.MaxElapsed {
			pc.MaxElapsed = snap.CurrentElapsed
		}
	}

	return pc
}

// countNonTerminalChildren returns how many children are running or pending.
func countNonTerminalChildren(
	snapshots map[ActivityID]ActivitySnapshot,
	children []*ActivityNode,
) int {
	count := 0

	for _, child := range children {
		snap := lookupSnapshot(snapshots, child.ID)
		if snap.Status != ActivityStatusCompleted && snap.Status != ActivityStatusFailed {
			count++
		}
	}

	return count
}

// formatCollapsedPhaseLabel builds the display string for a collapsed phase.
// For fully terminal phases: "◈ Code Formatting  6/6 · 4.1s".
// For partially running phases (viewport pressure collapse):
// "⏵ golangci-lint  15/30 · 12.3s" — shows progress while indicating work
// is still in flight.
func formatCollapsedPhaseLabel(snap ActivitySnapshot, pc PhaseCounts, theme Theme) (display string, c color.Color) {
	if pc.IsPartial() {
		symbol := SymbolRunning
		c = theme.Colors.Running

		if pc.Failed > 0 {
			c = theme.Colors.Failed
		}

		done := pc.Completed + pc.Failed
		display = fmt.Sprintf("%s %s  %d/%d", symbol, snap.Label, done, pc.Total())
	} else {
		symbol := SymbolPhase
		c = theme.Colors.Phase

		if pc.Failed > 0 {
			c = theme.Colors.Failed
		}

		display = fmt.Sprintf("%s %s  %d/%d", symbol, snap.Label, pc.Completed, pc.Total())
	}

	if pc.MaxElapsed > 0 {
		display += " · " + FormatDuration(pc.MaxElapsed)
	}

	return display, c
}

func (dt *DependencyTree) renderLine( //nolint:cyclop // label + sub-line rendering with multiple optional overlays
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
		activityDisplay, color = formatCollapsedPhaseLabel(snap, *entry.PhaseCounts, dt.theme)
	} else {
		activityDisplay, color = formatActivityLabelWithOptions(snap, dt.theme, labelOptions{
			OnCriticalPath: entry.OnCriticalPath,
			Convergence:    entry.Convergence,
			ShowCategory:   dt.showCategory,
		})
	}

	// Option B: extra dependencies sub-line. When showExtraDeps is enabled,
	// nodes with multiple deps show a dim "↳ Compile, Lint" sub-line beneath
	// the label — replacing the old inline "←" suffix. When disabled (the
	// default, matching nom), extra deps are silently absorbed into the tree.
	if dt.showExtraDeps {
		if subLine := dt.formatExtraDeps(node, snapshots); subLine != "" {
			activityDisplay += "\n" + subLine
		}
	}

	// Blockage view: pending nodes show which incomplete dependencies are
	// preventing them from starting. This is the first DAG-only feature — it
	// uses the full Deps set, not just the display parent.
	if dt.showBlockage && snap.Status == ActivityStatusPending {
		if subLine := dt.formatBlockage(node, snapshots); subLine != "" {
			activityDisplay += "\n" + subLine
		}
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

	// LayerNodes holds a single wrapped row of activities in layered mode.
	// When non-empty, Node is nil and this entry renders as one horizontal row.
	LayerNodes []*ActivityNode

	// LayerHeader is a synthetic layer header or separator line. When non-empty,
	// Node and LayerNodes are nil.
	LayerHeader string

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

	// OnCriticalPath is true when the node lies on the longest estimated-time
	// path through the DAG. Rendered with a ◆ prefix in tree mode.
	OnCriticalPath bool

	// Convergence is true when the node has more than one incoming dependency
	// (fan-in point). Rendered with a ◇ prefix in tree mode.
	Convergence bool
}

// PhaseCounts holds aggregate status counts for a collapsed phase's children.
// MaxElapsed tracks the longest child duration (not the sum), since DAG steps
// run in parallel — summing would over-report wall-clock time.
type PhaseCounts struct {
	Completed  int
	Failed     int
	Running    int
	Pending    int
	MaxElapsed time.Duration
}

// Total returns the total number of children accounted for.
func (pc PhaseCounts) Total() int {
	return pc.Completed + pc.Failed + pc.Running + pc.Pending
}

// IsPartial returns true when the phase still has running or pending children,
// meaning the collapse was triggered by viewport pressure rather than all
// children reaching terminal state.
func (pc PhaseCounts) IsPartial() bool {
	return pc.Running > 0 || pc.Pending > 0
}

// ContainsNode reports whether this visible entry represents the given activity
// ID. It checks real nodes, collapsed phase entries, and layered-mode rows.
func (entry VisibleEntry) ContainsNode(id ActivityID) bool {
	if entry.Node != nil && entry.Node.ID == id {
		return true
	}

	for _, node := range entry.LayerNodes {
		if node.ID == id {
			return true
		}
	}

	return false
}

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

	if dt.renderMode == RenderModeLayered {
		return dt.collectLayeredEntries(snapshots, maxHeight)
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

// RenderVisibleEntry renders a single visible entry — a real activity
// node, a synthetic collapse marker, or a layered-mode header/row — using
// immutable snapshot data. This is the marker-aware primitive renderers
// should call per line.
func (dt *DependencyTree) RenderVisibleEntry(
	entry VisibleEntry,
	snapshots map[ActivityID]ActivitySnapshot,
	maxWidth int,
) string {
	if dt.renderMode == RenderModeLayered || len(entry.LayerNodes) > 0 || entry.LayerHeader != "" {
		return dt.renderLayeredLine(entry, snapshots, maxWidth)
	}

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

	if dt.renderMode == RenderModeLayered {
		return dt.renderLayeredLine(VisibleEntry{LayerNodes: []*ActivityNode{node}}, snapshots, 0)
	}

	snap := lookupSnapshot(snapshots, node.ID)
	display, color := formatActivityLabelWithOptions(snap, dt.theme, labelOptions{})

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
		Color:  ThemeDefault.Colors.Pending,
	}
}

// formatExtraDeps renders the Option B sub-line for a node's non-display-parent
// dependencies. Returns "" if the node has no extra deps. The output is a dim
// "↳ Compile, Lint" line meant to be appended after a "\n" separator.
func (dt *DependencyTree) formatExtraDeps(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) string {
	extraDeps := node.ExtraDeps()
	if len(extraDeps) == 0 {
		return ""
	}

	depNames := make([]string, 0, len(extraDeps))

	for _, depID := range extraDeps {
		depSnap := lookupSnapshot(snapshots, depID)

		name := depSnap.Label
		if name == "" {
			name = depID.String()
		}

		depNames = append(depNames, name)
	}

	return lipgloss.NewStyle().
		Faint(true).
		Render(fmt.Sprintf("  %s %s", SymbolDeps, strings.Join(depNames, ", ")))
}

// formatBlockage renders a dim sub-line listing incomplete dependencies for a
// pending node. Returns "" when the node has no incomplete deps or is not
// pending. Each incomplete dep is shown with its status symbol and, if
// running, its elapsed time, so the user sees why the node is stuck.
//
// Example output: "  ⊘ blocked by Compile (running · 5s), Lint (pending)".
func (dt *DependencyTree) formatBlockage(
	node *ActivityNode,
	snapshots map[ActivityID]ActivitySnapshot,
) string {
	if len(node.Deps) == 0 {
		return ""
	}

	var blockedBy []string

	for _, depID := range node.Deps {
		depSnap := lookupSnapshot(snapshots, depID)
		if depSnap.Status == ActivityStatusCompleted {
			continue
		}

		name := depSnap.Label
		if name == "" {
			name = depID.String()
		}

		statusLabel := depSnap.Status.String()
		if depSnap.Status == ActivityStatusRunning {
			statusLabel = fmt.Sprintf("%s · %s", depSnap.Status.String(), FormatDuration(depSnap.CurrentElapsed))
		}

		blockedBy = append(blockedBy, fmt.Sprintf("%s (%s)", name, statusLabel))
	}

	if len(blockedBy) == 0 {
		return ""
	}

	return lipgloss.NewStyle().
		Faint(true).
		Render(fmt.Sprintf("  %s blocked by %s", SymbolBlocked, strings.Join(blockedBy, ", ")))
}
