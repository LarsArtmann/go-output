package graph

import (
	"io"

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
// Writes directly from the Graph's nodes and edges — no DOTRenderer intermediary.
func WriteDOT(w io.Writer, g output.Graph, opts ...DOTOption) error {
	cfg := defaultDOTConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	return output.WriteRenderedRaw(w, "write dot", renderDOTString(g.Nodes(), g.Edges(), cfg))
}

// RenderDOT renders a Graph as a DOT string.
func RenderDOT(g output.Graph, opts ...DOTOption) (string, error) {
	return output.RenderFromWrite(func(w io.Writer) error {
		return WriteDOT(w, g, opts...)
	})
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

	return output.WriteRenderedRawFrom(w, r.Render, "write mermaid", "render mermaid")
}

// RenderMermaid renders a Graph as a Mermaid string.
func RenderMermaid(g output.Graph, opts ...MermaidOption) (string, error) {
	return output.RenderFromWrite(func(w io.Writer) error {
		return WriteMermaid(w, g, opts...)
	})
}
