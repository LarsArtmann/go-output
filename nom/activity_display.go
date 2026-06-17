package nom

import (
	"image/color"
	"time"
)

// DisplayState holds the shared display fields synchronized between
// ActivityDisplayState and ActivityNode via SyncActivityTimingToTree.
type DisplayState struct {
	Status         ActivityStatus
	Symbol         string
	Color          color.Color
	StartTime      time.Time
	EstimatedTime  time.Duration
	CurrentElapsed time.Duration
}

// ============================================================================
// ACTIVITY DISPLAY STATE
// ============================================================================
// ActivityDisplayState represents the display state of an activity in NOM visualization.
type ActivityDisplayState struct {
	// Shared display state (synced to ActivityNode)
	DisplayState

	// Core activity information
	ActivityID   ActivityID
	ActivityName ActivityName
	// Timing information (ActivityDisplayState-specific)
	EndTime time.Time
	// Operation type for prefix symbol
	OperationType string // download, upload, copy, move, delete, ""
	// Error information
	Error error
	// Dependencies for tree rendering
	Dependencies []string
}

// NewActivityDisplayState creates a new ActivityDisplayState.
func NewActivityDisplayState(
	activityID ActivityID,
	activityName ActivityName,
) *ActivityDisplayState {
	return &ActivityDisplayState{
		ActivityID:   activityID,
		ActivityName: activityName,
		DisplayState: DisplayState{
			Status: ActivityStatusPending,
			Symbol: SymbolPaused,
			Color:  ColorPaused,
		},
		EndTime:       time.Time{},
		OperationType: "",
		Error:         nil,
		Dependencies:  make([]string, 0),
	}
}

// SetRunning marks activity as running.
func (ads *ActivityDisplayState) SetRunning() {
	ads.Status = ActivityStatusRunning
	ads.Symbol = SymbolRunning
	ads.Color = ColorRunning
	ads.StartTime = time.Now()
}

// SetCompleted marks activity as completed.
func (ads *ActivityDisplayState) SetCompleted() {
	ads.Status = ActivityStatusCompleted
	ads.Symbol = SymbolCompleted
	ads.Color = ColorCompleted
	ads.calculateElapsedTime()
}

// SetFailed marks activity as failed.
func (ads *ActivityDisplayState) SetFailed(err error) {
	ads.Status = ActivityStatusFailed
	ads.Symbol = SymbolFailed
	ads.Color = ColorFailed
	ads.calculateElapsedTime()
	ads.Error = err
}

// calculateElapsedTime calculates elapsed time if start time is set.
func (ads *ActivityDisplayState) calculateElapsedTime() {
	if ads.StartTime.IsZero() {
		return
	}

	ads.EndTime = time.Now()
	ads.CurrentElapsed = ads.EndTime.Sub(ads.StartTime)
}

// SetEstimatedTime sets estimated duration from cache.
func (ads *ActivityDisplayState) SetEstimatedTime(duration time.Duration) {
	ads.EstimatedTime = duration
}

func (ads *ActivityDisplayState) setOperationType(operationType string) {
	ads.OperationType = operationType
}

func (ads *ActivityDisplayState) addDependency(dep string) {
	ads.Dependencies = append(ads.Dependencies, dep)
}

// IsRunning returns true if activity is currently running.
func (ads *ActivityDisplayState) IsRunning() bool {
	return ads.Status == ActivityStatusRunning
}

// IsCompleted returns true if activity is completed.
func (ads *ActivityDisplayState) IsCompleted() bool {
	return ads.Status == ActivityStatusCompleted
}

// IsFailed returns true if activity has failed.
func (ads *ActivityDisplayState) IsFailed() bool {
	return ads.Status == ActivityStatusFailed
}

// Copy creates a deep copy of ActivityDisplayState.
// This ensures external modifications don't affect the original.
func (ads *ActivityDisplayState) Copy() *ActivityDisplayState {
	return &ActivityDisplayState{
		ActivityID:    ads.ActivityID,
		ActivityName:  ads.ActivityName,
		DisplayState:  ads.DisplayState,
		EndTime:       ads.EndTime,
		OperationType: ads.OperationType,
		Error:         ads.Error,
		Dependencies:  append([]string{}, ads.Dependencies...),
	}
}
