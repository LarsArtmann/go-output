package output

import (
	"fmt"
	"sync"
)

// RendererFactory is a function that creates a Renderer instance.
type RendererFactory func() Renderer

var (
	registry = make(map[Format]RendererFactory)
	regMu    sync.RWMutex
)

// Register registers a renderer factory for a format.
func Register(format Format, factory RendererFactory) error {
	regMu.Lock()
	defer regMu.Unlock()

	if _, exists := registry[format]; exists {
		return fmt.Errorf("format %q is already registered", format)
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

// Create returns a new Renderer instance for the given format.
func Create(format Format) (Renderer, error) {
	regMu.RLock()
	factory, exists := registry[format]
	regMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no renderer registered for format %q", format)
	}
	return factory(), nil
}

// RegisteredFormats returns a list of all registered formats.
func RegisteredFormats() []Format {
	regMu.RLock()
	defer regMu.RUnlock()

	formats := make([]Format, 0, len(registry))
	for f := range registry {
		formats = append(formats, f)
	}
	return formats
}

// IsRegistered returns true if a format has a registered renderer.
func IsRegistered(format Format) bool {
	regMu.RLock()
	defer regMu.RUnlock()
	_, exists := registry[format]
	return exists
}
