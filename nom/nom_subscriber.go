package nom

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/larsartmann/go-output"
)

// ActivityReader is the read-only contract for diagram export.
// NOMStyleSubscriber satisfies it via Store(), so any output.GraphRenderer
// (DOT, Mermaid, D2, PlantUML) can consume live progress state.
type ActivityReader interface {
	Nodes() []output.GraphNode
	Edges() []output.GraphEdge
}

// NOMStyleSubscriber implements EventSubscriber to provide NOM-style visualization.
type NOMStyleSubscriber struct {
	mu             sync.RWMutex
	activities     map[ActivityID]*Activity
	dependencyTree *DependencyTree
	timingCache    *TimingCache
	// counts is an incrementally-maintained cache of activity status counts.
	// Updated via applyDelta on every state transition and on creation/removal,
	// so GetActivityCounts is O(1) instead of O(n) per render frame. This is
	// the subscriber-level aggregate (the root of a per-subtree monoid);
	// per-node subtree counts are not maintained because no consumer needs them.
	counts       ActivityCounts
	workflowID   WorkflowID
	workflowName WorkflowName
	startTime    time.Time
	// showParallelism enables the "parallel: N/M possible" segment in the
	// inline renderer summary bar.
	showParallelism bool

	isRunning bool
	enabled   bool
	theme     Theme
}

// SubscriberOption configures a NOMStyleSubscriber at construction time.
type SubscriberOption func(*NOMStyleSubscriber)

// WithCachePath overrides the default timing-cache file path
// (~/.cache/nom-timing.csv). Tests inject a temp directory so the suite never
// reads or writes the real home directory, keeping it hermetic.
func WithCachePath(path string) SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.timingCache = NewTimingCache(withFilePath(path))
	}
}

// WithCollapseCompletedPhases enables phase-aware subtree collapsing in the
// dependency tree renderer. When all direct children of a phase are in
// terminal state (completed/failed), they are hidden and the phase renders
// a summary like "◈ Code Formatting  6/6 · 4.1s" instead of expanding every
// child. Consumers with many categories benefit from this to avoid walls of
// identical green checkmarks.
func WithCollapseCompletedPhases() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.collapseCompletedPhases = true
	}
}

// WithShowExtraDeps enables the Option B rendering mode: nodes with multiple
// dependencies show a dim "↳ Compile, Lint" sub-line beneath the label,
// making non-display-parent deps visible without cluttering the label.
// When disabled (the default, matching nom), extra deps are absorbed
// silently into the tree structure.
func WithShowExtraDeps() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.showExtraDeps = true
	}
}

// WithShowCriticalPath enables a ◆ prefix on nodes that lie on the longest
// estimated-time path through the dependency DAG. Off by default to preserve
// nom-compatible output.
func WithShowCriticalPath() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.showCriticalPath = true
	}
}

// WithShowConvergence enables a ◇ prefix on nodes with multiple incoming
// dependencies (DAG fan-in points). Off by default to preserve nom-compatible
// output.
func WithShowConvergence() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.showConvergence = true
	}
}

// WithShowBlockage enables a dim sub-line beneath pending nodes that lists
// their incomplete dependencies and current status. Off by default to avoid
// cluttering the tree when the parent is already visible.
func WithShowBlockage() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.showBlockage = true
	}
}

// WithRenderMode selects how the dependency tree is rendered: tree mode
// (default) or layered mode. Layered mode groups activities by DAG depth and
// renders each layer horizontally, making parallel work explicit.
func WithRenderMode(mode RenderMode) SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.renderMode = mode
	}
}

// WithShowParallelism enables a "parallel: N/M possible" segment in the inline
// renderer summary bar. Off by default to keep the summary compact.
func WithShowParallelism() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.showParallelism = true
	}
}

// WithShowCategory enables a dim [tag] prefix and category-tint coloring on
// activities that have a non-empty Category. Off by default to preserve
// nom-compatible output.
func WithShowCategory() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.showCategory = true
	}
}

// WithHideFutureLayers collapses layers where all nodes are still pending
// into a single "N pending" summary line, reducing visual noise from deep
// DAGs where only the first few layers are active. Layered mode only.
func WithHideFutureLayers() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.hideFutureLayers = true
	}
}

// WithCollapseCategories groups sibling activities by category and collapses
// them into a summary line (e.g. "3 build tasks"). Requires WithShowCategory
// to be useful, since categories must be set on activities.
func WithCollapseCategories() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.dependencyTree.collapseCategories = true
	}
}

// WithAutoTheme detects the terminal's background brightness via the
// COLORFGBG environment variable and selects an appropriate theme:
// dark backgrounds get ThemeDefault, light backgrounds get ThemeHighContrast.
// If COLORFGBG is unset (most modern terminals), ThemeDefault is used.
//
// COLORFGBG format is "fg;bg" where values 0-7 are dark and 8-15 are light.
// This is the same convention used by vim, tmux, and other terminal tools.
func WithAutoTheme() SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		theme := detectAutoTheme()
		ns.theme = theme
		ns.dependencyTree.theme = theme
	}
}

// detectAutoTheme inspects COLORFGBG and returns the best-matching theme.
func detectAutoTheme() Theme {
	// COLORFGBG is "fg;bg" — e.g. "15;0" means light-on-dark.
	fgBg := os.Getenv("COLORFGBG")
	if fgBg == "" {
		return ThemeDefault
	}

	parts := strings.Split(fgBg, ";")
	if len(parts) < 2 {
		return ThemeDefault
	}

	bg, err := strconv.Atoi(parts[1])
	if err != nil {
		return ThemeDefault
	}

	// 0-7 = dark background → use default (dark-optimized) theme.
	// 8-15 = light background → use high contrast theme.
	if bg >= 0 && bg <= 7 {
		return ThemeDefault
	}

	return ThemeHighContrast
}

// WithTheme sets the visual theme used for status symbols and colors.
// If not supplied, ThemeDefault is used.
func WithTheme(theme Theme) SubscriberOption {
	return func(ns *NOMStyleSubscriber) {
		ns.theme = theme
		ns.dependencyTree.theme = theme
	}
}

// NewNOMStyleSubscriber creates a new NOM-style subscriber.
func NewNOMStyleSubscriber(opts ...SubscriberOption) *NOMStyleSubscriber {
	ns := &NOMStyleSubscriber{
		activities:     make(map[ActivityID]*Activity),
		dependencyTree: NewDependencyTree(),
		timingCache:    NewTimingCache(),
		isRunning:      false,
		enabled:        true,
		theme:          ThemeDefault,
	}

	for _, opt := range opts {
		opt(ns)
	}

	return ns
}

// Theme returns the subscriber's active visual theme.
func (ns *NOMStyleSubscriber) Theme() Theme {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.theme
}

// DependencyTree returns the subscriber's dependency tree. Renderers and
// integration tests use this to inspect or configure display state directly.
func (ns *NOMStyleSubscriber) DependencyTree() *DependencyTree {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	return ns.dependencyTree
}

// Store returns an ActivityReader for diagram export. The projection is
// computed on-demand from the subscriber's current state — no bridge sync,
// no third state copy, always current.
//
// Example:
//
//	dot := graph.NewDOTRenderer()
//	dot.SetNodes(subscriber.Store().Nodes())
//	dot.SetEdges(subscriber.Store().Edges())
//	diagram, _ := dot.Render()
func (ns *NOMStyleSubscriber) Store() ActivityReader {
	return &subscriberView{ns: ns}
}

// subscriberView adapts NOMStyleSubscriber to the ActivityReader interface.
// It projects the subscriber's Activity map to GraphNode/Edge slices on-demand
// under the subscriber's read lock.
type subscriberView struct {
	ns *NOMStyleSubscriber
}

// Nodes projects all activities as output.GraphNode values for diagram export.
// Shape and Style are derived from Status at projection time (not cached on
// Activity), keeping the domain model decoupled from the graph framework.
func (v *subscriberView) Nodes() []output.GraphNode {
	v.ns.mu.RLock()
	defer v.ns.mu.RUnlock()

	out := make([]output.GraphNode, 0, len(v.ns.activities))
	for _, a := range v.ns.activities {
		out = append(out, output.GraphNode{
			ID:    a.ID,
			Label: a.Label,
			Shape: a.Status.NodeShape(),
			Style: a.Status.GraphStyle(),
		})
	}

	return out
}

// Edges projects the DAG's dependency edges for diagram export. Each edge
// goes FROM a dependency TO the dependent node — matching the tree's
// parent→child direction. Unlike the old implementation (which only walked
// display-tree Children), this iterates node.Deps and captures ALL edges,
// including non-display-parent dependencies.
//
// Lock ordering: acquires ns.mu.RLock first, then tree.mu.RLock.
// This ordering (subscriber → tree) is consistent across all code paths
// that nest both locks. Never reverse this order — it would deadlock.
func (v *subscriberView) Edges() []output.GraphEdge {
	v.ns.mu.RLock()
	defer v.ns.mu.RUnlock()

	tree := v.ns.dependencyTree

	tree.mu.RLock()
	defer tree.mu.RUnlock()

	var edges []output.GraphEdge

	for _, node := range tree.nodes {
		toID := string(node.ID)
		for _, depID := range node.Deps {
			edges = append(edges, output.GraphEdge{
				From: output.NewBrandedID[output.GraphNodeIDBrand](string(depID)),
				To:   output.NewBrandedID[output.GraphNodeIDBrand](toID),
			})
		}
	}

	return edges
}
