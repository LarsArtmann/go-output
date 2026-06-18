package nom

import (
	"context"
	"sync"
)

// MultiSubscriber fans out events to multiple EventSubscribers, like
// io.MultiWriter for event streams. Each subscriber receives every event;
// errors from one subscriber do not prevent others from receiving the event.
//
// Example:
//
//	multi := nom.NewMultiSubscriber(
//	    nom.NewNOMStyleSubscriber(),  // drives the TUI
//	    logSubscriber,                // logs to disk
//	    metricsSubscriber,            // emits metrics
//	)
//	multi.OnEvent(ctx, event) // all three receive it
type MultiSubscriber struct {
	mu          sync.RWMutex
	subscribers []EventSubscriber
}

// NewMultiSubscriber creates a fanout subscriber that dispatches to all
// provided subscribers. Nil subscribers are silently skipped.
func NewMultiSubscriber(subscribers ...EventSubscriber) *MultiSubscriber {
	subs := make([]EventSubscriber, 0, len(subscribers))
	for _, s := range subscribers {
		if s != nil {
			subs = append(subs, s)
		}
	}
	return &MultiSubscriber{subscribers: subs}
}

// Add appends a subscriber to the fanout. Safe to call concurrently.
func (m *MultiSubscriber) Add(s EventSubscriber) {
	if s == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribers = append(m.subscribers, s)
}

// OnEvent dispatches the event to all subscribers. Each subscriber receives
// the event regardless of whether a prior subscriber returned an error.
// Returns the first error encountered, or nil if all succeeded.
func (m *MultiSubscriber) OnEvent(ctx context.Context, event Event) error {
	m.mu.RLock()
	subs := make([]EventSubscriber, len(m.subscribers))
	copy(subs, m.subscribers)
	m.mu.RUnlock()

	var firstErr error
	for _, s := range subs {
		if err := s.OnEvent(ctx, event); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Subscribers returns a snapshot of the current subscriber list.
func (m *MultiSubscriber) Subscribers() []EventSubscriber {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]EventSubscriber, len(m.subscribers))
	copy(out, m.subscribers)
	return out
}
