package nom

import (
	"slices"
	"sync"
)

const (
	// msgNoActivitiesToDisplay is shown when dependency tree is empty.
	msgNoActivitiesToDisplay = "No activities to display"
)

// ActivityNode represents a node in the dependency tree.
type ActivityNode struct {
	// Shared display state (synced from ActivityDisplayState)
	DisplayState

	// Core activity information
	ActivityID   ActivityID
	ActivityName string
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
	order     []ActivityID             // Display order (smart filtered)
	loaded    bool                     // Whether tree has been built
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
		ActivityID:   id,
		ActivityName: name,
		DisplayState: DisplayState{
			Status: ActivityStatusPending,
			Symbol: SymbolPaused,
			Color:  ColorPaused,
		},
		Children:    make([]*ActivityNode, 0),
		IsDisplayed: true,
	}
}

// hasChild returns true if this node already has a child with the given activity ID.
func (n *ActivityNode) hasChild(id ActivityID) bool {
	return slices.ContainsFunc(n.Children, func(c *ActivityNode) bool {
		return c.ActivityID == id
	})
}

// hasSecondaryParent returns true if this node already has the given activity ID
// as a secondary parent.
func (n *ActivityNode) hasSecondaryParent(id ActivityID) bool {
	return slices.Contains(n.SecondaryParents, id)
}
