package nom

import "time"

// testSetStatus is a test helper that mutates the shared Activity pointer
// directly, replacing the old UpdateActivityStatus API. Symbol/color are
// derived from status via applyVisualStyle.
func testSetStatus(dt *DependencyTree, id ActivityID, status ActivityStatus, startTime time.Time) {
	node := dt.GetNode(id)
	if node == nil {
		return
	}

	node.Status = status
	node.applyVisualStyle()
	node.StartTime = startTime
}
