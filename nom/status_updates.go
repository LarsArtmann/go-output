package nom

import (
	"errors"
	"fmt"
	"image/color"
	"time"
)

// ErrActivityNotFound is returned when an activity cannot be found.
var ErrActivityNotFound = errors.New("activity not found")

// UpdateActivityStatus updates the status of an activity in the tree.
func (dt *DependencyTree) UpdateActivityStatus(
	activityID ActivityID,
	status ActivityStatus,
	symbol string,
	color color.Color,
	startTime time.Time,
	estimatedTime time.Duration,
) error {
	dt.mu.Lock()
	defer dt.mu.Unlock()

	node, exists := dt.nodes[activityID]
	if !exists {
		return fmt.Errorf(
			"failed to update activity %s: status=%s, symbol=%s, estimatedTime=%s: %w",
			activityID,
			status,
			symbol,
			estimatedTime,
			ErrActivityNotFound,
		)
	}

	node.Status = status
	node.Symbol = symbol
	node.Color = color
	node.StartTime = startTime
	node.EstimatedTime = estimatedTime

	return nil
}
