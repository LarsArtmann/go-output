package nom

import (
	"image/color"
	"time"

	"github.com/larsartmann/go-output"
)

// Activity is the unified source of truth for a workflow activity's identity,
// lifecycle, and NOM-display state. It is a pure domain type — it no longer
// embeds output.GraphNode, so it carries no diagram-export concerns (Shape,
// Style, Metadata). Those are projected from Status at export time by
// subscriberView.Nodes(), keeping the domain model decoupled from the graph
// rendering framework.
type Activity struct {
	// ID and Label are the activity's identity, expressed via the root
	// package's branded types so they flow directly into GraphNode projection.
	ID    output.GraphNodeID
	Label output.GraphNodeLabel

	// Kind is the topological classification (Task or Phase). Set at
	// construction, never changes. Distinct from Status, which is the
	// lifecycle state. Phase nodes render with SymbolPhase/Colors.Phase.
	Kind ActivityKind
	// Status is the typed domain state (pending/running/completed/failed).
	Status ActivityStatus
	// StartTime is when the activity transitioned to Running (zero if not started).
	StartTime time.Time
	// EndTime is when the activity transitioned to Completed or Failed.
	EndTime time.Time
	// EstimatedTime is the predicted duration from the timing cache (zero if unknown).
	EstimatedTime time.Duration
	// Err holds the failure error if Status == ActivityStatusFailed.
	Err error
	// Symbol is the NOM-style display symbol (cached from Status for rendering).
	Symbol Symbol
	// Color is the lipgloss terminal color (cached from Status for rendering).
	Color color.Color
	// CurrentElapsed is updated periodically for running activities.
	CurrentElapsed time.Duration
	// Host optionally names where the activity runs (e.g. a build machine).
	// Populated from ActivityStarted.Host; rendered when non-empty.
	Host string
	// Download optionally tracks byte-progress for the activity. Populated from
	// ActivityStarted.Download; rendered as a progress bar when active.
	Download DownloadProgress
}

// NewActivity creates a Task Activity with branded ID/Label and default
// visual style derived from the pending status. Use NewPhase for grouping nodes.
func NewActivity(id, name string) *Activity {
	return newActivity(id, name, ActivityKindTask)
}

// NewPhase creates a Phase Activity — a grouping node whose children are the
// real deliverables. Phase nodes render with SymbolPhase/Colors.Phase
// regardless of their lifecycle Status. The kind is fixed at construction and
// never changes.
func NewPhase(id, name string) *Activity {
	return newActivity(id, name, ActivityKindPhase)
}

// newActivity is the shared constructor; Kind selects Task vs Phase rendering.
func newActivity(id, name string, kind ActivityKind) *Activity {
	a := &Activity{
		Kind:   kind,
		Status: ActivityStatusPending,
	}
	a.ID = output.NewBrandedID[output.GraphNodeIDBrand](id)
	a.Label = output.NewBrandedID[output.GraphNodeLabelBrand](name)
	a.applyDisplayStyle()

	return a
}

// SetRunning transitions the activity to running and stamps StartedAt.
func (a *Activity) SetRunning() {
	a.Status = ActivityStatusRunning
	a.StartTime = time.Now()
	a.EndTime = time.Time{}
	a.Err = nil
	a.applyDisplayStyle()
}

// SetCompleted transitions the activity to completed, stamps EndedAt, and
// finalizes CurrentElapsed so renderers show the total run duration.
func (a *Activity) SetCompleted() {
	a.Status = ActivityStatusCompleted

	a.EndTime = time.Now()
	if !a.StartTime.IsZero() {
		a.CurrentElapsed = a.EndTime.Sub(a.StartTime)
	}

	a.applyDisplayStyle()
}

// SetFailed transitions the activity to failed, records the error, stamps
// EndedAt, and finalizes CurrentElapsed.
func (a *Activity) SetFailed(err error) {
	a.Status = ActivityStatusFailed
	a.Err = err

	a.EndTime = time.Now()
	if !a.StartTime.IsZero() {
		a.CurrentElapsed = a.EndTime.Sub(a.StartTime)
	}

	a.applyDisplayStyle()
}

// SetEstimatedTime sets the predicted duration from the timing cache.
func (a *Activity) SetEstimatedTime(d time.Duration) {
	a.EstimatedTime = d
}

// IsRunning returns true if the activity is currently in the running state.
func (a *Activity) IsRunning() bool { return a.Status == ActivityStatusRunning }

// IsCompleted returns true if the activity has completed successfully.
func (a *Activity) IsCompleted() bool { return a.Status == ActivityStatusCompleted }

// IsFailed returns true if the activity has failed.
func (a *Activity) IsFailed() bool { return a.Status == ActivityStatusFailed }

// IsPhase reports whether this activity is a Phase grouping node (Kind ==
// ActivityKindPhase), as opposed to a concrete Task.
func (a *Activity) IsPhase() bool { return a.Kind.IsPhase() }

// Copy creates a shallow copy of the Activity. All fields are value types
// or immutable, so a shallow copy is sufficient.
func (a *Activity) Copy() *Activity {
	cpy := *a
	return &cpy
}

// applyDisplayStyle caches the NOM terminal Symbol and Color from the current
// Status. Diagram-export Shape/Style are NOT cached here — they are projected
// from Status at export time by subscriberView.Nodes(), keeping Activity
// decoupled from the graph rendering framework.
func (a *Activity) applyDisplayStyle() {
	a.Symbol = a.Status.GetSymbol()
	a.Color = a.Status.GetColor()
}
