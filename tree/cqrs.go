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
	if err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

// RenderASCII renders a TreeNode as an ASCII tree string.
func RenderASCII(root *output.TreeNode, opts ...Option) (string, error) {
	var buf strings.Builder
	if err := WriteASCII(&buf, root, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteMarkdown writes a TreeNode as a nested Markdown bullet list to the provided writer.
func WriteMarkdown(w io.Writer, root *output.TreeNode) error {
	if root == nil {
		return nil
	}

	return writeMarkdownNode(w, root, 0)
}

func writeMarkdownNode(w io.Writer, node *output.TreeNode, depth int) error {
	indent := strings.Repeat("  ", depth)

	label := rootLabel(node)

	if _, err := fmt.Fprintf(w, "%s- %s\n", indent, label); err != nil {
		return fmt.Errorf("write markdown node: %w", err)
	}

	for _, child := range node.Children {
		if err := writeMarkdownNode(w, child, depth+1); err != nil {
			return err
		}
	}

	return nil
}

func rootLabel(node *output.TreeNode) string {
	if !node.Label.IsZero() {
		return node.Label.Get()
	}

	if !node.ID.IsZero() {
		return node.ID.Get()
	}

	return ""
}

// RenderMarkdown renders a TreeNode as a nested Markdown bullet list string.
func RenderMarkdown(root *output.TreeNode) (string, error) {
	var buf strings.Builder
	if err := WriteMarkdown(&buf, root); err != nil {
		return "", err
	}

	return buf.String(), nil
}
