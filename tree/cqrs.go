package tree

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// Option configures ASCII tree rendering.
type Option func(*Config)

// Config holds tree rendering configuration.
type Config struct {
	colorMode output.ColorMode
}

// WithColorMode sets the color output mode.
func WithColorMode(mode output.ColorMode) Option {
	return func(c *Config) { c.colorMode = mode }
}

// WriteASCII writes a TreeNode as an ASCII tree to the provided writer.
func WriteASCII(w io.Writer, root *output.TreeNode, opts ...Option) error {
	cfg := Config{colorMode: output.ColorModeAuto}
	for _, opt := range opts {
		opt(&cfg)
	}

	r := NewASCIITreeRenderer()
	r.SetColorMode(cfg.colorMode)
	r.SetRoot(root)

	out, err := r.Render()
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	_, err = io.WriteString(w, out)
	return fmt.Errorf("write output: %w", err)
}

// RenderASCII renders a TreeNode as an ASCII tree string.
func RenderASCII(root *output.TreeNode, opts ...Option) (string, error) {
	var buf strings.Builder
	if err := WriteASCII(&buf, root, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}
