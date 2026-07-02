package nom

import (
	"cmp"
	"errors"
	"slices"
	"sync"
)

// ErrActivityNotFound is returned when an activity cannot be found in the tree.
var ErrActivityNotFound = errors.New("activity not found")

const (
	// msgNoActivitiesToDisplay is shown when dependency tree is empty.
	// Mirrors tui.MsgNoActivitiesToDisplay — identical string, kept separate
	// because nom and tui are distinct modules. See split-brain M3.
	msgNoActivitiesToDisplay = "No activities to display"
)

// ActivityNode represents a node in the dependency DAG.
// Holds ONLY structural data + an ID for snapshot lookup. The mutable Activity
// fields (Status/Symbol/Color/timing) are never read from the node —
// rendering uses ActivitySnapshot values taken under the subscriber's lock.
// This makes the render-vs-event-handler data race unrepresentable.
type ActivityNode struct {
	// ID identifies the activity for snapshot lookup at render time.
	ID ActivityID

	// Deps holds ALL dependency IDs — the true DAG edges (source of truth).
	// Each entry means "this node depends on depID" (depID must complete
	// before this node can start). Unlike the old SecondaryParents, these
	// are not display-only: they drive edge export, cycle detection, and
	// depth computation.
	Deps []ActivityID

	// Display layout — computed at Build() time, NOT at AddActivity time.
	// The tree walk follows Parent/Children, not Deps.
	Parent      *ActivityNode
	Children    []*ActivityNode
	Depth       int
	IsRoot      bool
	IsDisplayed bool
}

// DependencyTree manages the hierarchical structure of activities and their dependencies.
type DependencyTree struct {
	mu     sync.RWMutex
	nodes  map[ActivityID]*ActivityNode // All nodes by activity ID
	roots  []*ActivityNode              // Root nodes (no dependencies)
	loaded bool                         // Whether tree has been built

	// collapseCompletedPhases enables phase-aware subtree collapsing: when a
	// phase node has ALL direct children in terminal state (completed/failed),
	// the children are hidden and the phase renders a summary line like
	// "◈ Code Formatting  6/6 · 4.1s". Disabled by default; consumers with
	// many categories (e.g. BuildFlow) enable it to avoid walls of green.
	collapseCompletedPhases bool

	// showExtraDeps enables a dim sub-line beneath nodes with multiple
	// dependencies, showing the non-display-parent deps as "↳ Compile, Lint".
	// When false (default, matching nom's behavior), extra deps are silently
	// absorbed into the tree — each node appears once, under its deepest
	// parent, with no annotation.
	showExtraDeps bool

	// showCriticalPath enables a ◆ prefix on nodes that lie on the longest
	// estimated-time path through the DAG.
	showCriticalPath bool

	// showConvergence enables a ◇ prefix on nodes with multiple incoming
	// dependencies (fan-in points in the DAG).
	showConvergence bool

	// showBlockage enables a dim sub-line beneath pending nodes that lists
	// incomplete dependencies and their status, making blockers explicit.
	showBlockage bool
}

// NewDependencyTree creates a new dependency tree.
func NewDependencyTree() *DependencyTree {
	return &DependencyTree{
		nodes:  make(map[ActivityID]*ActivityNode),
		roots:  make([]*ActivityNode, 0),
		loaded: false,
	}
}

// newActivityNode creates a structural placeholder node for a dependency
// that has not been registered with a full Activity yet (e.g. a parent ID
// referenced before its own activity.started event arrives).
func newActivityNode(id ActivityID, _ string) *ActivityNode {
	return &ActivityNode{
		ID:          id,
		Children:    make([]*ActivityNode, 0),
		IsDisplayed: true,
	}
}

// nodeHasID returns a predicate matching nodes whose ID equals the given value.
func nodeHasID(id ActivityID) func(*ActivityNode) bool {
	return func(c *ActivityNode) bool {
		return c.ID == id
	}
}

// sortNodesByID sorts a slice of ActivityNode in place by ID for deterministic display order.
func sortNodesByID(nodes []*ActivityNode) {
	slices.SortStableFunc(nodes, func(a, b *ActivityNode) int {
		return cmp.Compare(string(a.ID), string(b.ID))
	})
}

// hasChild returns true if this node already has a child with the given activity ID.
func (n *ActivityNode) hasChild(id ActivityID) bool {
	return slices.ContainsFunc(n.Children, nodeHasID(id))
}

// removeChild removes the child with the given activity ID, if present.
func (n *ActivityNode) removeChild(id ActivityID) {
	n.Children = slices.DeleteFunc(n.Children, nodeHasID(id))
}

// hasDep returns true if this node already records the given dependency ID.
func (n *ActivityNode) hasDep(id ActivityID) bool {
	return slices.Contains(n.Deps, id)
}

// ExtraDeps returns dependency IDs that are NOT the display parent — i.e.,
// the "extra" edges beyond the primary tree edge. Used by the Option B
// renderer to show "↳ Compile, Lint" beneath a multi-dependency node.
func (n *ActivityNode) ExtraDeps() []ActivityID {
	var extra []ActivityID

	for _, depID := range n.Deps {
		if n.Parent == nil || depID != n.Parent.ID {
			extra = append(extra, depID)
		}
	}

	return extra
}

// NodeClass classifies a tree node by its structural position, mirroring NOM's
// mapRootsTwigsAndLeaves. Roots anchor the tree, leaves are terminal activities
// (the actual deliverables), and twigs are the intermediaries between them.
type NodeClass string

const (
	NodeClassRoot NodeClass = "root"
	NodeClassTwig NodeClass = "twig"
	NodeClassLeaf NodeClass = "leaf"
)

// Class returns the node's structural classification. A root has no parent
// (or IsRoot); a leaf has no children; everything else is a twig.
func (n *ActivityNode) Class() NodeClass {
	switch {
	case n.IsRoot || n.Parent == nil:
		return NodeClassRoot
	case len(n.Children) == 0:
		return NodeClassLeaf
	default:
		return NodeClassTwig
	}
}
