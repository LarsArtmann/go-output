package nom

import (
	"image/color"
	"maps"
	"time"

	"github.com/larsartmann/go-output"
)

// Activity is the unified source of truth for a workflow activity's identity,
// visual representation, and temporal state. It embeds output.GraphNode so
// the same instance can be consumed by any output.GraphRenderer (DOT, Mermaid,
// D2, PlantUML) for diagram export.
type Activity struct {
	output.GraphNode

	// Status is the typed domain state (pending/running/completed/failed/paused).
	Status ActivityStatus
	// StartTime is when the activity transitioned to Running (zero if not started).
	StartTime time.Time
	// EndTime is when the activity transitioned to Completed or Failed.
	EndTime time.Time
	// EstimatedTime is the predicted duration from the timing cache (zero if unknown).
	EstimatedTime time.Duration
	// Err holds the failure error if Status == ActivityStatusFailed.
	Err error
	// OperationType labels the activity for prefix symbols ("download", "upload", "").
	OperationType string
	// Symbol is the NOM-style display symbol (cached from Status for rendering).
	Symbol string
	// Color is the lipgloss terminal color (cached from Status for rendering).
	Color color.Color
	// CurrentElapsed is updated periodically for running activities.
	CurrentElapsed time.Duration
	// Dependencies lists parent activity IDs (for tree rendering).
	Dependencies []string
}

// NewActivity creates an Activity with a branded GraphNode ID and default
// visual style derived from the pending status.
func NewActivity(id, name string) *Activity {
	a := &Activity{
		Status:       ActivityStatusPending,
		Dependencies: make([]string, 0),
	}
	a.ID = output.NewBrandedID[output.GraphNodeIDBrand](id)
	a.Label = output.NewBrandedID[output.GraphNodeLabelBrand](name)
	a.applyVisualStyle()

	return a
}

// SetRunning transitions the activity to running and stamps StartedAt.
func (a *Activity) SetRunning() {
	a.Status = ActivityStatusRunning
	a.StartTime = time.Now()
	a.EndTime = time.Time{}
	a.Err = nil
	a.applyVisualStyle()
}

// SetCompleted transitions the activity to completed, stamps EndedAt, and
// finalizes CurrentElapsed so renderers show the total run duration.
func (a *Activity) SetCompleted() {
	a.Status = ActivityStatusCompleted
	a.EndTime = time.Now()
	if !a.StartTime.IsZero() {
		a.CurrentElapsed = a.EndTime.Sub(a.StartTime)
	}
	a.applyVisualStyle()
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
	a.applyVisualStyle()
}

// SetEstimatedTime sets the predicted duration from the timing cache.
func (a *Activity) SetEstimatedTime(d time.Duration) {
	a.EstimatedTime = d
}

// Elapsed returns the duration since StartTime if running, or EndTime-StartTime
// if finished. Returns zero if not started or if StartTime is zero.
func (a *Activity) Elapsed() time.Duration {
	if a.StartTime.IsZero() {
		return 0
	}

	if a.Status == ActivityStatusRunning {
		return time.Since(a.StartTime)
	}

	if !a.EndTime.IsZero() {
		return a.EndTime.Sub(a.StartTime)
	}

	return time.Since(a.StartTime)
}

// IsRunning returns true if the activity is currently in the running state.
func (a *Activity) IsRunning() bool { return a.Status == ActivityStatusRunning }

// IsCompleted returns true if the activity has completed successfully.
func (a *Activity) IsCompleted() bool { return a.Status == ActivityStatusCompleted }

// IsFailed returns true if the activity has failed.
func (a *Activity) IsFailed() bool { return a.Status == ActivityStatusFailed }

// Copy creates a deep copy of the Activity.
func (a *Activity) Copy() *Activity {
	cpy := *a // shallow copy is fine for value types
	if a.Dependencies != nil {
		cpy.Dependencies = append([]string{}, a.Dependencies...)
	}

	if a.Metadata != nil {
		cpy.Metadata = make(map[string]string, len(a.Metadata))
		maps.Copy(cpy.Metadata, a.Metadata)
	}

	return &cpy
}

// applyVisualStyle sets the GraphNode Shape and Style from the current Status,
// so diagram export (DOT/Mermaid/D2) automatically reflects the activity state.
func (a *Activity) applyVisualStyle() {
	a.Shape = a.Status.NodeShape()
	a.Style = a.Status.GraphStyle()
	a.Symbol = a.Status.GetSymbol()
	a.Color = a.Status.GetColor()
}

// SetPaused transitions the activity to paused.
func (a *Activity) SetPaused() {
	a.Status = ActivityStatusPaused
	a.applyVisualStyle()
}

func (a *Activity) setOperationType(operationType string) {
	a.OperationType = operationType
}

func (a *Activity) addDependency(dep string) {
	a.Dependencies = append(a.Dependencies, dep)
}
