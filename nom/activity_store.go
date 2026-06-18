package nom

import (
	"sync"

	"github.com/larsartmann/go-output"
)

// ActivityStore is the map-backed primary store for activities. It provides
// O(1) lookups by ID and projects to output.GraphNode/output.GraphEdge slices
// so any output.GraphRenderer (DOT, Mermaid, D2, PlantUML) can consume the
// current activity state for diagram export.
//
// Thread-safe via sync.RWMutex. All mutations go through the write lock;
// projections (Nodes/Edges/Roots/Counts) snapshot under the read lock.
type ActivityStore struct {
	mu         sync.RWMutex
	activities map[output.GraphNodeID]*Activity
	edges      []output.GraphEdge
}

// NewActivityStore creates an empty ActivityStore.
func NewActivityStore() *ActivityStore {
	return &ActivityStore{
		activities: make(map[output.GraphNodeID]*Activity),
		edges:      make([]output.GraphEdge, 0),
	}
}

// Upsert inserts or replaces the activity for its ID.
func (s *ActivityStore) Upsert(a *Activity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activities[a.ID] = a
}

// Get returns the activity for the given ID, or nil if not found.
func (s *ActivityStore) Get(id output.GraphNodeID) (*Activity, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.activities[id]
	return a, ok
}

// All returns a snapshot slice of all activities (unordered).
func (s *ActivityStore) All() []*Activity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Activity, 0, len(s.activities))
	for _, a := range s.activities {
		out = append(out, a)
	}
	return out
}

// AddEdge records a directed dependency edge (from → to).
func (s *ActivityStore) AddEdge(from, to output.GraphNodeID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.edges = append(s.edges, output.GraphEdge{From: from, To: to})
}

// Nodes projects all activities as output.GraphNode values (value copies).
// This is the integration point for output.GraphRenderer.SetNodes().
func (s *ActivityStore) Nodes() []output.GraphNode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]output.GraphNode, 0, len(s.activities))
	for _, a := range s.activities {
		out = append(out, a.GraphNode)
	}
	return out
}

// Edges returns a copy of all dependency edges.
// This is the integration point for output.GraphRenderer.SetEdges().
func (s *ActivityStore) Edges() []output.GraphEdge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]output.GraphEdge, len(s.edges))
	copy(out, s.edges)
	return out
}

// Roots returns IDs of nodes that have no incoming edges (no parent).
// These are the top-level activities in the dependency tree.
func (s *ActivityStore) Roots() []output.GraphNodeID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hasParent := make(map[output.GraphNodeID]bool, len(s.edges))
	for _, e := range s.edges {
		hasParent[e.To] = true
	}

	var roots []output.GraphNodeID
	for id := range s.activities {
		if !hasParent[id] {
			roots = append(roots, id)
		}
	}
	return roots
}

// Counts tallies activities by status category.
// Paused activities are counted as pending (same as the prior GetActivityCounts).
func (s *ActivityStore) Counts() (running, completed, failed, pending int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.activities {
		switch a.Status {
		case ActivityStatusRunning:
			running++
		case ActivityStatusCompleted:
			completed++
		case ActivityStatusFailed:
			failed++
		case ActivityStatusPending, ActivityStatusPaused:
			pending++
		}
	}
	return running, completed, failed, pending
}

// Size returns the number of activities.
func (s *ActivityStore) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.activities)
}

// Clear removes all activities and edges.
func (s *ActivityStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activities = make(map[output.GraphNodeID]*Activity)
	s.edges = make([]output.GraphEdge, 0)
}
