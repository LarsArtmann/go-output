package graph

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// DOTOption configures DOT rendering.
type DOTOption func(*dotConfig)

type dotConfig struct {
	directed bool
	graphID  string
	rankdir  RankDir
	splines  SplineStyle
	nodesep  string
	ranksep  string
}

// WithDirected sets whether the graph is directed (digraph) or undirected (graph).
func WithDirected(directed bool) DOTOption {
	return func(c *dotConfig) { c.directed = directed }
}

// WithGraphID sets the DOT graph identifier.
func WithGraphID(id string) DOTOption {
	return func(c *dotConfig) { c.graphID = id }
}

// WithDOTRankDir sets the rank direction.
func WithDOTRankDir(dir RankDir) DOTOption {
	return func(c *dotConfig) { c.rankdir = dir }
}

// WithDOTSplines sets the edge routing style.
func WithDOTSplines(style SplineStyle) DOTOption {
	return func(c *dotConfig) { c.splines = style }
}

// WithNodeSep sets the minimum space between nodes in the same rank.
func WithNodeSep(sep string) DOTOption {
	return func(c *dotConfig) {
		if isValidNumericSep(sep) {
			c.nodesep = sep
		}
	}
}

// WithRankSep sets the minimum space between ranks.
func WithRankSep(sep string) DOTOption {
	return func(c *dotConfig) {
		if isValidNumericSep(sep) {
			c.ranksep = sep
		}
	}
}

//nolint:goconst // default sep values are fine to repeat
func defaultDOTConfig() dotConfig {
	return dotConfig{
		directed: true,
		graphID:  "G",
		rankdir:  RankDirTB,
		splines:  SplineOrtho,
		nodesep:  defaultNodeSep,
		ranksep:  defaultRankSep,
	}
}

// WriteDOT writes a Graph as DOT format to the provided writer.
func WriteDOT(w io.Writer, g output.Graph, opts ...DOTOption) error {
	cfg := defaultDOTConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	r := newDOTRenderer(cfg.directed)
	r.graphID = cfg.graphID
	r.rankdir = cfg.rankdir
	r.splines = cfg.splines
	r.nodesep = cfg.nodesep
	r.ranksep = cfg.ranksep
	r.SetNodes(g.Nodes())
	r.SetEdges(g.Edges())

	out, err := r.Render()
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, out)
	if err != nil {
		return fmt.Errorf("write dot output: %w", err)
	}

	return nil
}

// RenderDOT renders a Graph as a DOT string.
func RenderDOT(g output.Graph, opts ...DOTOption) (string, error) {
	var buf strings.Builder
	if err := WriteDOT(&buf, g, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// MermaidOption configures Mermaid rendering.
type MermaidOption func(*mermaidConfig)

type mermaidConfig struct {
	codeFence bool
}

// WithCodeFence enables or disables the ```mermaid code fence wrapper.
func WithCodeFence(enabled bool) MermaidOption {
	return func(c *mermaidConfig) { c.codeFence = enabled }
}

// WriteMermaid writes a Graph as Mermaid format to the provided writer.
func WriteMermaid(w io.Writer, g output.Graph, opts ...MermaidOption) error {
	cfg := mermaidConfig{codeFence: true}
	for _, opt := range opts {
		opt(&cfg)
	}

	r := NewMermaidRenderer()
	r.codeFence = cfg.codeFence
	r.SetNodes(g.Nodes())
	r.SetEdges(g.Edges())

	out, err := r.Render()
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, out)
	if err != nil {
		return fmt.Errorf("write mermaid output: %w", err)
	}

	return nil
}

// RenderMermaid renders a Graph as a Mermaid string.
func RenderMermaid(g output.Graph, opts ...MermaidOption) (string, error) {
	var buf strings.Builder
	if err := WriteMermaid(&buf, g, opts...); err != nil {
		return "", err
	}

	return buf.String(), nil
}
