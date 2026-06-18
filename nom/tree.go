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
type ActivityNode struct {
	// Activity is the source of truth (embeds GraphNode + Status + timing).
	// Shared by pointer with the subscriber's activities map — mutations are
	// instantly visible to both with zero sync overhead.
	*Activity

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
	mu        sync.RWMutex
	buildOnce sync.Once
	nodes     map[ActivityID]*ActivityNode // All nodes by activity ID
	roots     []*ActivityNode              // Root nodes (no dependencies)
	order     []ActivityID                 // Display order (smart filtered)
	loaded    bool                         // Whether tree has been built
}

// NewDependencyTree creates a new dependency tree.
func NewDependencyTree() *DependencyTree {
	return &DependencyTree{
		nodes:  make(map[ActivityID]*ActivityNode),
		roots:  make([]*ActivityNode, 0),
		order:  make([]ActivityID, 0),
		loaded: false,
	}
}

// newActivityNode creates a new ActivityNode with pending status.
func newActivityNode(id ActivityID, name string) *ActivityNode {
	return &ActivityNode{
		Activity:    NewActivity(string(id), name),
		Children:    make([]*ActivityNode, 0),
		IsDisplayed: true,
	}
}

// nodeHasID returns a predicate matching nodes whose ID equals the given value.
func nodeHasID(id ActivityID) func(*ActivityNode) bool {
	return func(c *ActivityNode) bool {
		return c.ID.Get() == string(id)
	}
}

// sortNodesByID sorts a slice of ActivityNode in place by ID for deterministic display order.
func sortNodesByID(nodes []*ActivityNode) {
	slices.SortStableFunc(nodes, func(a, b *ActivityNode) int {
		return cmp.Compare(a.ID.Get(), b.ID.Get())
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
