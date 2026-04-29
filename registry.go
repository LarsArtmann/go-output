// Package output provides a registry for dynamically registering renderer factories.
//
// The registry is an opt-in plugin system. By default, formats are used directly
// via their constructors (NewMarkdownTable, NewD2Diagram, etc.). Use Register/Create
// when you need runtime-dispatchable format selection, e.g.:
//
//	output.Register(output.FormatJSON, func() output.Renderer { return myCustomJSONRenderer })
//	renderer, _ := output.Create(output.FormatJSON)
//
// The registry is thread-safe.
package output

import (
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
func Unregister(format Format) {
	regMu.Lock()
	defer regMu.Unlock()

	delete(registry, format)
}

// ErrNoRendererRegistered is returned when no renderer is registered for a format.
var ErrNoRendererRegistered = errors.New("no renderer registered for format")

// Create returns a new Renderer instance for the given format.
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
func RegisteredFormats() []Format {
	regMu.RLock()
	defer regMu.RUnlock()

	formats := make([]Format, 0, len(registry))
	for f := range registry {
		formats = append(formats, f)
	}

	slices.SortFunc(formats, func(a, b Format) int {
		if a.String() < b.String() {
			return -1
		}

		if a.String() > b.String() {
			return 1
		}

		return 0
	})

	return formats
}

// IsRegistered returns true if a format has a registered renderer.
func IsRegistered(format Format) bool {
	regMu.RLock()
	defer regMu.RUnlock()

	_, exists := registry[format]

	return exists
}
