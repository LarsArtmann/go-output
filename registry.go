package output

import (
	"sync"
)

// formatRegistry is a thread-safe map from Format to a generic value of
// type T. It deduplicates the registry boilerplate previously repeated for
// table-data marshalers, any-data marshalers, and format shape capabilities.
//
// T is typically a function type (TableRenderer, UnknownRenderer) or
// a value type ([]Shape).
type formatRegistry[T any] struct {
	mu    sync.RWMutex
	items map[Format]T
}

// newFormatRegistry creates an empty registry.
func newFormatRegistry[T any]() *formatRegistry[T] {
	return &formatRegistry[T]{
		items: make(map[Format]T),
	}
}

// register stores a value for the given format, replacing any prior
// registration.
func (r *formatRegistry[T]) register(format Format, value T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items[format] = value
}

// get returns the value for the given format, or the zero value and false
// if no value is registered.
func (r *formatRegistry[T]) get(format Format) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.items[format]

	return v, ok
}

// formats returns all formats currently registered, in unspecified order.
func (r *formatRegistry[T]) formats() []Format {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Format, 0, len(r.items))
	for f := range r.items {
		out = append(out, f)
	}

	return out
}
