package nom

import (
	"sync"
	"time"
)

// NOMStyleSubscriber implements EventSubscriber to provide NOM-style visualization.
type NOMStyleSubscriber struct {
	mu sync.RWMutex
	// Activity tracking
	activities     map[ActivityID]*ActivityDisplayState
	store          *ActivityStore // projection layer for diagram export (DOT/Mermaid/D2)
	dependencyTree *DependencyTree
	timingCache    *TimingCache
	// Workflow state
	workflowID   WorkflowID
	workflowName WorkflowName
	startTime    time.Time
	isRunning    bool
	// Configuration
	enabled bool
}

// NewNOMStyleSubscriber creates a new NOM-style subscriber.
func NewNOMStyleSubscriber() *NOMStyleSubscriber {
	return &NOMStyleSubscriber{
		activities:     make(map[ActivityID]*ActivityDisplayState),
		store:          NewActivityStore(),
		dependencyTree: NewDependencyTree(),
		timingCache:    NewTimingCache(),
		isRunning:      false,
		enabled:        true,
	}
}

// Store returns the ActivityStore for diagram export. The store projects the
// current activity state as output.GraphNode/output.GraphEdge slices, consumable
// by any output.GraphRenderer (DOT, Mermaid, D2, PlantUML).
//
// Example:
//
//	dot := graph.NewDOTRenderer()
//	dot.SetNodes(subscriber.Store().Nodes())
//	dot.SetEdges(subscriber.Store().Edges())
//	diagram, _ := dot.Render()
func (ns *NOMStyleSubscriber) Store() *ActivityStore {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return ns.store
}
