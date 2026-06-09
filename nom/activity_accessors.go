package nom
import (
)
// GetActivities returns a copy of all activity display states.
func (ns *NOMStyleSubscriber) GetActivities() map[ActivityID]*ActivityDisplayState {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	// Return a deep copy to prevent external modification
	result := make(map[ActivityID]*ActivityDisplayState)
	for id, activity := range ns.activities {
		result[id] = activity.Copy()
	}
	return result
}
// GetActivity returns a specific activity display state.
func (ns *NOMStyleSubscriber) GetActivity(
	activityID ActivityID,
) *ActivityDisplayState {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	activity, exists := ns.activities[activityID]
	if !exists {
		return nil
	}
	return activity.Copy()
}
