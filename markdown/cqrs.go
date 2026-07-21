package markdown

import (
	"io"

	"github.com/larsartmann/go-output"
)

// Option configures Markdown table rendering.
type Option func(*Config)

// Config holds Markdown rendering configuration.
type Config struct {
	output.ColorConfig
}

// WithColorMode sets the color output mode.
func WithColorMode(mode output.ColorMode) Option {
	return func(c *Config) { c.ColorMode = mode }
}

// Write writes a Table as a Markdown table to the provided writer.
func Write(w io.Writer, data *output.Table, opts ...Option) error {
	cfg := Config{ColorConfig: output.DefaultColorConfig()}
	for _, opt := range opts {
		opt(&cfg)
	}

	m := NewMarkdownTableFromTable(data)
	m.SetColorMode(cfg.ColorMode)

	return output.WriteRenderedRawFrom(w, m.Render, "markdown", "markdown")
}

// Render renders a Table as a Markdown string.
func Render(data *output.Table, opts ...Option) (string, error) {
	return output.RenderFromWrite(func(w io.Writer) error {
		return Write(w, data, opts...)
	})
}
