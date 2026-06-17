package nom

import (
	"testing"
	"time"
)

func TestElideCompletedUnderPressure(t *testing.T) {
	t.Parallel()

	completedNode := &ActivityNode{
		ActivityID:   ActivityID("c1"),
		ActivityName: "Completed Step",
		DisplayState: DisplayState{
			Status:         ActivityStatusCompleted,
			CurrentElapsed: 5 * time.Second,
		},
	}
	runningNode := &ActivityNode{
		ActivityID:   ActivityID("r1"),
		ActivityName: "Running Step",
		DisplayState: DisplayState{
			Status:         ActivityStatusRunning,
			CurrentElapsed: 2 * time.Second,
		},
	}

	dt := NewDependencyTree()

	tests := []struct {
		name         string
		children     []*ActivityNode
		maxHeight    int
		visibleCount int
		wantLen      int
		wantIDs      []string
	}{
		{
			name:         "no height pressure keeps all children",
			children:     []*ActivityNode{completedNode, runningNode},
			maxHeight:    100,
			visibleCount: 1,
			wantLen:      2,
			wantIDs:      []string{"c1", "r1"},
		},
		{
			name:         "unlimited height keeps all children",
			children:     []*ActivityNode{completedNode, runningNode},
			maxHeight:    0,
			visibleCount: 1,
			wantLen:      2,
			wantIDs:      []string{"c1", "r1"},
		},
		{
			name:         "height pressure elides completed",
			children:     []*ActivityNode{completedNode, runningNode},
			maxHeight:    2,
			visibleCount: 1,
			wantLen:      1,
			wantIDs:      []string{"r1"},
		},
		{
			name:         "only completed children under pressure returns empty",
			children:     []*ActivityNode{completedNode},
			maxHeight:    1,
			visibleCount: 1,
			wantLen:      0,
		},
		{
			name:         "no children returns empty",
			children:     nil,
			maxHeight:    10,
			visibleCount: 0,
			wantLen:      0,
		},
		{
			name:         "all active children kept even under pressure",
			children:     []*ActivityNode{runningNode},
			maxHeight:    1,
			visibleCount: 1,
			wantLen:      1,
			wantIDs:      []string{"r1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := dt.elideCompletedUnderPressure(tt.children, tt.maxHeight, tt.visibleCount)

			if len(got) != tt.wantLen {
				t.Fatalf("got %d children, want %d", len(got), tt.wantLen)
			}

			for i, wantID := range tt.wantIDs {
				if i >= len(got) {
					break
				}

				if string(got[i].ActivityID) != wantID {
					t.Errorf("child[%d] = %s, want %s", i, got[i].ActivityID, wantID)
				}
			}
		})
	}
}

func TestActivityStatus_Interest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status ActivityStatus
		want   int
	}{
		{ActivityStatusFailed, 0},
		{ActivityStatusRunning, 1},
		{ActivityStatusPaused, 2},
		{ActivityStatusPending, 3},
		{ActivityStatusCompleted, 4},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			t.Parallel()

			got := tt.status.Interest()
			if got != tt.want {
				t.Errorf("%s.Interest() = %d, want %d", tt.status, got, tt.want)
			}
		})
	}

	// Verify ordering: lower interest = higher priority
	if ActivityStatusFailed.Interest() >= ActivityStatusRunning.Interest() {
		t.Error("Failed should have lower interest (higher priority) than Running")
	}

	if ActivityStatusRunning.Interest() >= ActivityStatusCompleted.Interest() {
		t.Error("Running should have lower interest than Completed")
	}
}
