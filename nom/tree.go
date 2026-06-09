package nom
import (
	"image/color"
	"sync"
	"time"
)
const (
	// msgNoActivitiesToDisplay is shown when dependency tree is empty.
	msgNoActivitiesToDisplay = "No activities to display"
)
// TreeNode represents a node in the dependency tree.
type TreeNode struct {
	// Core activity information
	ActivityID   ActivityID
	ActivityName string
	Status       ActivityStatus
	Symbol       string
	Color        color.Color
	// Timing information
	StartTime      time.Time
	EstimatedTime  time.Duration
	CurrentElapsed time.Duration
	// Tree structure
	Parent   *TreeNode
	Children []*TreeNode
	Depth    int
	// Display state
	IsRoot      bool
	IsDisplayed bool
}
// DependencyTree manages the hierarchical structure of activities and their dependencies.
type DependencyTree struct {
	mu        sync.RWMutex
	buildOnce sync.Once
	nodes     map[ActivityID]*TreeNode // All nodes by activity ID
	roots     []*TreeNode                    // Root nodes (no dependencies)
	order     []ActivityID             // Display order (smart filtered)
	loaded    bool                           // Whether tree has been built
}
// NewDependencyTree creates a new dependency tree.
func NewDependencyTree() *DependencyTree {
	return &DependencyTree{
		nodes:  make(map[ActivityID]*TreeNode),
		roots:  make([]*TreeNode, 0),
		order:  make([]ActivityID, 0),
		loaded: false,
	}
}
// newTreeNode creates a new TreeNode with pending status.
func newTreeNode(id ActivityID, name string) *TreeNode {
	return &TreeNode{
		ActivityID:   id,
		ActivityName: name,
		Status:       ActivityStatusPending,
		Symbol:       SymbolPaused,
		Color:        ColorPaused,
		Children:     make([]*TreeNode, 0),
		IsDisplayed:  true,
	}
}
