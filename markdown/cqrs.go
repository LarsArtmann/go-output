package markdown

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// Option configures Markdown table rendering.
type Option func(*Config)

// Config holds Markdown rendering configuration.
type Config struct {
	colorMode output.ColorMode
}

// WithColorMode sets the color output mode.
func WithColorMode(mode output.ColorMode) Option {
	return func(c *Config) { c.colorMode = mode }
}

// Write writes a Table as a Markdown table to the provided writer.
func Write(w io.Writer, data *output.Table, opts ...Option) error {
	cfg := Config{colorMode: output.ColorModeAuto}
	for _, opt := range opts {
		opt(&cfg)
	}

	m := NewMarkdownTableFromTable(data)
	m.SetColorMode(cfg.colorMode)

	out, err := m.Render()
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	_, err = io.WriteString(w, out)

	return fmt.Errorf("write output: %w", err)
}

// Render renders a Table as a Markdown string.
func Render(data *output.Table, opts ...Option) (string, error) {
	var buf strings.Builder
	if err := Write(&buf, data, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}
