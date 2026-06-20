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

// ActivityNode represents a node in the dependency tree.
// Holds ONLY tree structure + an ID for snapshot lookup. The mutable Activity
// fields (Status/Symbol/Color/timing) are never read from the node —
// rendering uses ActivitySnapshot values taken under the subscriber's lock.
// This makes the render-vs-event-handler data race unrepresentable.
type ActivityNode struct {
	// ID identifies the activity for snapshot lookup at render time.
	ID ActivityID

	// Tree structure
	Parent           *ActivityNode
	Children         []*ActivityNode
	SecondaryParents []ActivityID // Non-primary dependencies (for display only)
	Depth            int
	// Display state
	IsRoot      bool
	IsDisplayed bool
}

// DependencyTree manages the hierarchical structure of activities and their dependencies.
type DependencyTree struct {
	mu     sync.RWMutex
	nodes  map[ActivityID]*ActivityNode // All nodes by activity ID
	roots  []*ActivityNode              // Root nodes (no dependencies)
	loaded bool                         // Whether tree has been built
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

// hasSecondaryParent returns true if this node already has the given activity ID
// as a secondary parent.
func (n *ActivityNode) hasSecondaryParent(id ActivityID) bool {
	return slices.Contains(n.SecondaryParents, id)
}
