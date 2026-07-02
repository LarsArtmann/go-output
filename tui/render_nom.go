package tui

import (
	"strings"
	"time"

	"charm.land/lipgloss/v2"

	"github.com/larsartmann/go-output/nom"
)

// filterCriticalPath returns only entries whose nodes are on the critical path.
func filterCriticalPath(
	entries []nom.VisibleEntry,
	criticalIDs map[nom.ActivityID]bool,
) []nom.VisibleEntry {
	var filtered []nom.VisibleEntry

	for _, entry := range entries {
		// Keep layer headers and separators for context.
		if entry.LayerHeader != "" {
			filtered = append(filtered, entry)
			continue
		}

		for _, node := range entry.LayerNodes {
			if criticalIDs[node.ID] {
				filtered = append(filtered, entry)
				break
			}
		}

		if entry.Node != nil && criticalIDs[entry.Node.ID] {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

// renderNOMStyle creates a NOM-style display with dependency tree.
func (m *ProgressModel) renderNOMStyle() string {
	m.treeStartLine = 0

	var contentSections []string

	if m.dependencyTree != nil {
		treeSection := m.renderDependencyTree()
		if treeSection != "" {
			contentSections = append(contentSections, treeSection)
		}
	}

	result := m.assembleProgressSections(
		"🔄 NOM Workflow Execution",
		contentSections,
		m.renderNOMSummaryBar,
	)
	m.treeStartLine = 2

	return result
}

// renderDependencyTree renders the activity dependency tree in NOM style.
//
// Uses immutable ActivitySnapshot data for all field reads — no shared
// *Activity pointer is accessed, so event handlers (SetRunning/SetCompleted/
// SetFailed) mutating activities on dispatcher goroutines cannot race the
// render. The snapshot is taken under the subscriber's read lock (brief),
// then the tree walk reads only immutable data (no lock held).
//
// When scrollOffset > 0 the visible window is selected at the entry level
// (before rendering), avoiding the O(n) render+string-clip of the previous
// approach that rendered the entire tree then split the string to extract
// the visible lines.
func (m *ProgressModel) renderDependencyTree() string {
	if m.dependencyTree == nil || m.nomSubscriber == nil {
		return ""
	}

	snapshots := m.nomSubscriber.SnapshotActivities()
	tree := m.dependencyTree

	treeHeight := m.height - chromeLines
	if treeHeight <= 0 {
		treeHeight = defaultTreeHeight
	}

	var entries []nom.VisibleEntry

	if m.scrollOffset > 0 {
		// Scroll path: collect all entries (no height-pressure collapsing),
		// then select the visible window by offset.
		allEntries := tree.VisibleEntriesWithSnapshots(snapshots, 0)

		if len(allEntries) == 0 {
			m.visibleEntries = nil
			m.scrollOffset = 0

			return msgNoActivitiesToDisplay
		}

		// Clamp scrollOffset to valid range. scrollToBottomSentinel maps to
		// the last page.
		maxOffset := max(0, len(allEntries)-treeHeight)
		if m.scrollOffset > maxOffset {
			m.scrollOffset = maxOffset
		}

		end := min(m.scrollOffset+treeHeight, len(allEntries))
		entries = allEntries[m.scrollOffset:end]
	} else {
		entries = tree.VisibleEntriesWithSnapshots(snapshots, treeHeight)
	}

	if len(entries) == 0 {
		m.visibleEntries = nil

		return msgNoActivitiesToDisplay
	}

	// Critical-path filter: when enabled, show only entries whose nodes are on
	// the longest estimated-time path through the DAG.
	if m.criticalPathFilter && m.nomSubscriber != nil {
		criticalIDs := m.dependencyTree.CriticalPathIDs(snapshots)
		if len(criticalIDs) > 0 {
			entries = filterCriticalPath(entries, criticalIDs)
		}
	}

	m.visibleEntries = entries

	lines := make([]string, 0, len(entries))

	for _, entry := range entries {
		line := tree.RenderVisibleEntry(entry, snapshots, m.width)

		// Only real activity nodes and layered-mode rows are selectable; collapse-
		// marker lines (entry.Node == nil and no LayerNodes) never get the
		// selection highlight.
		if m.selectedNode != "" && entry.ContainsNode(m.selectedNode) {
			line = lipgloss.NewStyle().
				Background(colors.selectBG).
				Foreground(colors.selectFG).
				Render(line)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderNOMSummaryBar creates the NOM-style summary bar, colored by workflow state.
func (m *ProgressModel) renderNOMSummaryBar() string {
	counts := m.getActivityCounts()
	elapsed := time.Since(m.startTime)
	remaining := m.estimatedRemaining()
	summary := buildNOMSummary(counts, elapsed, remaining)
	baseStyle := createSummaryStyle(m.theme.Colors.Info)

	switch m.workflowState {
	case workflowStateIdle, workflowStateRunning:
		return baseStyle.Render(summary)
	case workflowStateCompleted:
		return baseStyle.Foreground(m.theme.Colors.Completed).Render("✅ " + summary)
	case workflowStateErrored:
		return baseStyle.Foreground(m.theme.Colors.Failed).Render("❌ " + summary)
	default:
		return baseStyle.Render(summary)
	}
}

// getActivityCounts delegates to the subscriber for counts.
func (m *ProgressModel) getActivityCounts() nom.ActivityCounts {
	return m.nomSubscriber.GetActivityCounts()
}

// estimatedRemaining delegates to the subscriber's projected remaining time
// (sum of pending/running activity estimates). Returns 0 when there is no
// subscriber or no unfinished estimated work, in which case "~Xm left" is
// omitted from the summary bar.
func (m *ProgressModel) estimatedRemaining() time.Duration {
	if m.nomSubscriber == nil {
		return 0
	}

	return m.nomSubscriber.EstimatedTotalRemaining()
}

func (m *ProgressModel) renderHelpOverlay() string {
	helpStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Foreground(colors.helpFG).
		Background(colors.helpBG)

	shortcuts := []string{
		"  Keyboard Shortcuts",
		"",
		"  j / ↓       Scroll down",
		"  k / ↑       Scroll up",
		"  pgdown      Scroll half page",
		"  pgup        Scroll half page up",
		"  g / Home    Scroll to top",
		"  G / End     Scroll to bottom",
		"  L           Toggle tree / layered mode",
		"  C           Toggle critical-path-only filter",
		"  D           Export DAG to DOT file",
		"  ?           Toggle this help",
		"  q / ctrl+c  Quit",
	}

	helpText := helpStyle.Render(strings.Join(shortcuts, "\n"))

	width := m.width
	height := m.height

	if width <= 0 {
		width = defaultHelpWidth
	}

	if height <= 0 {
		height = defaultHelpHeight
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, helpText)
}
