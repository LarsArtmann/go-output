package nom
import (
	"sync"
	"time"
)
// NOMStyleSubscriber implements EventSubscriber to provide NOM-style visualization.
type NOMStyleSubscriber struct {
	mu sync.RWMutex
	// Activity tracking
	activities     map[ActivityID]*ActivityDisplayState
	dependencyTree *DependencyTree
	timingCache    *TimingCache
	// Workflow state
	workflowID   WorkflowID
	workflowName WorkflowName
	startTime    time.Time
	isRunning    bool
	// Configuration
	enabled bool
}
// NewNOMStyleSubscriber creates a new NOM-style subscriber.
func NewNOMStyleSubscriber() *NOMStyleSubscriber {
	return &NOMStyleSubscriber{
		activities:     make(map[ActivityID]*ActivityDisplayState),
		dependencyTree: NewDependencyTree(),
		timingCache:    NewTimingCache(),
		isRunning:      false,
		enabled:        true,
	}
}
