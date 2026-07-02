package nom

import (
	"image/color"
	"strings"
	"sync"

	"github.com/larsartmann/go-output"
)

// StatusDef describes the visual and diagram-export properties of a single
// activity status. Custom CI states (e.g. "skipped", "cached", "warning") are
// registered as StatusDef values and receive an ActivityStatus ID.
type StatusDef struct {
	Name       string
	Symbol     Symbol
	Color      color.Color
	Interest   int
	NodeShape  output.NodeShape
	GraphStyle output.GraphStyle
}

// statusRegistry is a thread-safe, open registry of activity statuses. Core
// statuses (pending, running, completed, failed) are pre-registered at IDs
// 0-3 for backward compatibility; additional statuses are allocated IDs
// starting from 4.
type statusRegistry struct {
	mu     sync.RWMutex
	byID   map[int]StatusDef
	byName map[string]int
	nextID int
}

func newStatusRegistry() *statusRegistry {
	r := &statusRegistry{
		byID:   make(map[int]StatusDef),
		byName: make(map[string]int),
		nextID: 0,
	}

	r.registerLocked(
		"pending",
		SymbolPending,
		Colors.Pending,
		2,
		output.NodeShapeEllipse,
		output.GraphStyle{
			Fill:      "#e5e7eb",
			Stroke:    "#9ca3af",
			FontColor: "#374151",
		},
	)

	r.registerLocked(
		"running",
		SymbolRunning,
		Colors.Running,
		1,
		output.NodeShapeBox,
		output.GraphStyle{
			Fill:      "#16a34a",
			Stroke:    "#15803d",
			FontColor: "#ffffff",
		},
	)

	r.registerLocked(
		"completed",
		SymbolCompleted,
		Colors.Completed,
		3,
		output.NodeShapeRect, //nolint:staticcheck
		output.GraphStyle{
			Fill:      "#6b7280",
			Stroke:    "#4b5563",
			FontColor: "#ffffff",
		},
	)

	r.registerLocked(
		"failed",
		SymbolFailed,
		Colors.Failed,
		0,
		output.NodeShapeDiamond,
		output.GraphStyle{
			Fill:      "#dc2626",
			Stroke:    "#991b1b",
			FontColor: "#ffffff",
		},
	)

	return r
}

//nolint:gochecknoglobals // package-level registry is the design intent
var globalRegistry = newStatusRegistry()

// RegisterStatus adds a new status to the global registry and returns its
// ActivityStatus ID. If a status with the same name already exists, its
// existing ID is returned. Name matching is case-insensitive.
func RegisterStatus(
	name string,
	symbol Symbol,
	c color.Color,
	interest int,
	shape output.NodeShape,
	style output.GraphStyle,
) ActivityStatus {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	if id, exists := globalRegistry.byName[normalizeStatusName(name)]; exists {
		return ActivityStatus(id)
	}

	id := globalRegistry.nextID
	globalRegistry.nextID++

	globalRegistry.byID[id] = StatusDef{
		Name:       name,
		Symbol:     symbol,
		Color:      c,
		Interest:   interest,
		NodeShape:  shape,
		GraphStyle: style,
	}
	globalRegistry.byName[normalizeStatusName(name)] = id

	return ActivityStatus(id)
}

// LookupStatus returns the definition for a registered status. The second
// return value is false for unknown IDs.
func LookupStatus(id ActivityStatus) (StatusDef, bool) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	def, ok := globalRegistry.byID[int(id)]

	return def, ok
}

// AllActivityStatuses returns the registered ActivityStatus IDs in ascending
// order. It is the dynamic equivalent of the former hardcoded slice and is
// updated whenever a new status is registered.
func AllActivityStatuses() []ActivityStatus {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	return globalRegistry.allStatusesLocked()
}

// AllRegisteredStatuses returns every status currently in the registry, in
// ascending ID order.
func AllRegisteredStatuses() []StatusDef {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	return globalRegistry.allRegisteredStatusesLocked()
}

// allStatusesLocked returns the registered ActivityStatus IDs in ascending
// order. The caller must hold the registry lock.
func (r *statusRegistry) allStatusesLocked() []ActivityStatus {
	out := make([]ActivityStatus, 0, len(r.byID))

	for id := range r.nextID {
		if _, ok := r.byID[id]; ok {
			out = append(out, ActivityStatus(id))
		}
	}

	return out
}

// allRegisteredStatusesLocked returns the registered StatusDef values in
// ascending ID order. The caller must hold the registry lock.
func (r *statusRegistry) allRegisteredStatusesLocked() []StatusDef {
	out := make([]StatusDef, 0, len(r.byID))

	for id := range r.nextID {
		if def, ok := r.byID[id]; ok {
			out = append(out, def)
		}
	}

	return out
}

// normalizeStatusName returns a lowercase canonical form for status-name lookup.
func normalizeStatusName(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// registerLocked registers a status under the next available ID. The caller
// must hold the write lock.
func (r *statusRegistry) registerLocked(
	name string,
	symbol Symbol,
	c color.Color,
	interest int,
	shape output.NodeShape,
	style output.GraphStyle,
) {
	id := r.nextID
	r.nextID++

	r.byID[id] = StatusDef{
		Name:       name,
		Symbol:     symbol,
		Color:      c,
		Interest:   interest,
		NodeShape:  shape,
		GraphStyle: style,
	}
	r.byName[normalizeStatusName(name)] = id
}
