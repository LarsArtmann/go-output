package nom

import (
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
	workflowID     WorkflowID
	workflowName   WorkflowName
	startTime      time.Time
	isRunning      bool
	enabled        bool
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

// NewNOMStyleSubscriber creates a new NOM-style subscriber.
func NewNOMStyleSubscriber(opts ...SubscriberOption) *NOMStyleSubscriber {
	ns := &NOMStyleSubscriber{
		activities:     make(map[ActivityID]*Activity),
		dependencyTree: NewDependencyTree(),
		timingCache:    NewTimingCache(),
		isRunning:      false,
		enabled:        true,
	}

	for _, opt := range opts {
		opt(ns)
	}

	return ns
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
// Since Activity embeds output.GraphNode, the projection is a direct copy —
// Shape/Style are always in sync with Status via applyVisualStyle().
func (v *subscriberView) Nodes() []output.GraphNode {
	v.ns.mu.RLock()
	defer v.ns.mu.RUnlock()

	out := make([]output.GraphNode, 0, len(v.ns.activities))
	for _, a := range v.ns.activities {
		out = append(out, a.GraphNode)
	}

	return out
}

// Edges projects the dependency tree's edges for diagram export.
//
// Lock ordering: acquires ns.mu.RLock first, then tree.mu.RLock.
// This ordering (subscriber → tree) is consistent across all code paths
// that nest both locks. Never reverse this order — it would deadlock.
func (v *subscriberView) Edges() []output.GraphEdge {
	v.ns.mu.RLock()
	defer v.ns.mu.RUnlock()

	// Derive edges from the dependency tree's parent-child relationships.
	tree := v.ns.dependencyTree

	tree.mu.RLock()
	defer tree.mu.RUnlock()

	var edges []output.GraphEdge

	for _, node := range tree.nodes {
		parentID := string(node.ID)
		for _, child := range node.Children {
			edges = append(edges, output.GraphEdge{
				From: output.NewBrandedID[output.GraphNodeIDBrand](parentID),
				To:   output.NewBrandedID[output.GraphNodeIDBrand](string(child.ID)),
			})
		}
	}

	return edges
}
