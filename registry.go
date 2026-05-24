// Package output provides consistent output formatting for CLI applications.
package output

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
	"sync"
)

// ErrFormatAlreadyRegistered is returned when a format is already registered.
var ErrFormatAlreadyRegistered = errors.New("format already registered")

// RendererFactory is a function that creates a Renderer instance.
type RendererFactory func() Renderer

var (
	//nolint:gochecknoglobals // Registry map for format-to-renderer factory mapping.
	registry = make(map[Format]RendererFactory)
	//nolint:gochecknoglobals // Mutex protects concurrent access to registry.
	regMu sync.RWMutex
)

// Register registers a renderer factory for a format.
//
// Deprecated: Use format constructors directly. This function may be removed in a future version.
func Register(format Format, factory RendererFactory) error {
	regMu.Lock()
	defer regMu.Unlock()

	if _, exists := registry[format]; exists {
		return fmt.Errorf(
			"register factory %v for format %q: %w",
			factory,
			format,
			ErrFormatAlreadyRegistered,
		)
	}

	registry[format] = factory

	return nil
}

// Unregister removes a format from the registry.
//
// Deprecated: Use format constructors directly.
func Unregister(format Format) {
	regMu.Lock()
	defer regMu.Unlock()

	delete(registry, format)
}

// ErrNoRendererRegistered is returned when no renderer is registered for a format.
var ErrNoRendererRegistered = errors.New("no renderer registered for format")

// Create returns a new Renderer instance for the given format.
//
// Deprecated: Use format constructors directly. This function may be removed in a future version.
func Create(format Format) (Renderer, error) {
	regMu.RLock()

	factory, exists := registry[format]

	regMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %q", ErrNoRendererRegistered, format)
	}

	return factory(), nil
}

// RegisteredFormats returns a sorted list of all registered formats.
//
// Deprecated: Use format constructors directly.
func RegisteredFormats() []Format {
	regMu.RLock()
	defer regMu.RUnlock()

	formats := make([]Format, 0, len(registry))
	for f := range registry {
		formats = append(formats, f)
	}

	slices.SortFunc(formats, func(a, b Format) int {
		return cmp.Compare(a.String(), b.String())
	})

	return formats
}

// IsRegistered returns true if a format has a registered renderer.
//
// Deprecated: Use format constructors directly.
func IsRegistered(format Format) bool {
	regMu.RLock()
	defer regMu.RUnlock()

	_, exists := registry[format]

	return exists
}
