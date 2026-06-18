package nom

// GetActivities returns a copy of all activities.
func (ns *NOMStyleSubscriber) GetActivities() map[ActivityID]*Activity {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	result := make(map[ActivityID]*Activity, len(ns.activities))
	for id, activity := range ns.activities {
		result[id] = activity.Copy()
	}

	return result
}

// GetActivity returns a specific activity.
func (ns *NOMStyleSubscriber) GetActivity(
	activityID ActivityID,
) *Activity {
	ns.mu.RLock()
	defer ns.mu.RUnlock()

	activity, exists := ns.activities[activityID]
	if !exists {
		return nil
	}

	return activity.Copy()
}
