package nom

import (
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
	// StartedAt is when the activity transitioned to Running (zero if not started).
	StartedAt time.Time
	// EndedAt is when the activity transitioned to Completed or Failed.
	EndedAt time.Time
	// Estimate is the predicted duration from the timing cache (zero if unknown).
	Estimate time.Duration
	// Err holds the failure error if Status == ActivityStatusFailed.
	Err error
	// OperationType labels the activity for prefix symbols ("download", "upload", "").
	OperationType string
}

// NewActivity creates an Activity with a branded GraphNode ID and default
// visual style derived from the pending status.
func NewActivity(id, name string) *Activity {
	a := &Activity{
		Status: ActivityStatusPending,
	}
	a.ID = output.NewBrandedID[output.GraphNodeIDBrand](id)
	a.Label = output.NewBrandedID[output.GraphNodeLabelBrand](name)
	a.applyVisualStyle()
	return a
}

// SetRunning transitions the activity to running and stamps StartedAt.
func (a *Activity) SetRunning() {
	a.Status = ActivityStatusRunning
	a.StartedAt = time.Now()
	a.EndedAt = time.Time{}
	a.Err = nil
	a.applyVisualStyle()
}

// SetCompleted transitions the activity to completed and stamps EndedAt.
func (a *Activity) SetCompleted() {
	a.Status = ActivityStatusCompleted
	a.EndedAt = time.Now()
	a.applyVisualStyle()
}

// SetFailed transitions the activity to failed, records the error, and stamps EndedAt.
func (a *Activity) SetFailed(err error) {
	a.Status = ActivityStatusFailed
	a.Err = err
	a.EndedAt = time.Now()
	a.applyVisualStyle()
}

// SetEstimatedTime sets the predicted duration from the timing cache.
func (a *Activity) SetEstimatedTime(d time.Duration) {
	a.Estimate = d
}

// Elapsed returns the duration since StartedAt if running, or EndedAt-StartedAt
// if finished. Returns zero if not started or if StartTime is zero.
func (a *Activity) Elapsed() time.Duration {
	if a.StartedAt.IsZero() {
		return 0
	}
	if a.Status == ActivityStatusRunning {
		return time.Since(a.StartedAt)
	}
	if !a.EndedAt.IsZero() {
		return a.EndedAt.Sub(a.StartedAt)
	}
	return time.Since(a.StartedAt)
}

// IsRunning returns true if the activity is currently in the running state.
func (a *Activity) IsRunning() bool { return a.Status == ActivityStatusRunning }

// IsCompleted returns true if the activity has completed successfully.
func (a *Activity) IsCompleted() bool { return a.Status == ActivityStatusCompleted }

// IsFailed returns true if the activity has failed.
func (a *Activity) IsFailed() bool { return a.Status == ActivityStatusFailed }

// applyVisualStyle sets the GraphNode Shape and Style from the current Status,
// so diagram export (DOT/Mermaid/D2) automatically reflects the activity state.
func (a *Activity) applyVisualStyle() {
	a.Shape = a.Status.NodeShape()
	a.Style = a.Status.GraphStyle()
}
