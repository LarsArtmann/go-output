package serialization

import (
	"github.com/larsartmann/go-output"
)

type graphView struct {
	Nodes []graphNodeView `json:"nodes" toml:"nodes" yaml:"nodes"`
	Edges []graphEdgeView `json:"edges" toml:"edges" yaml:"edges"`
}

type graphNodeView struct {
	ID       string            `json:"id"                 toml:"id"                 yaml:"id"`
	Label    string            `json:"label"              toml:"label"              yaml:"label"`
	Shape    string            `json:"shape,omitempty"    toml:"shape,omitempty"    yaml:"shape,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty" toml:"metadata,omitempty" yaml:"metadata,omitempty"`
}

type graphEdgeView struct {
	From  string `json:"from"            toml:"from"            yaml:"from"`
	To    string `json:"to"              toml:"to"              yaml:"to"`
	Label string `json:"label,omitempty" toml:"label,omitempty" yaml:"label,omitempty"`
}

func buildGraphView(mixin output.GraphBuilder) graphView {
	graph := graphView{
		Nodes: make([]graphNodeView, 0, len(mixin.Nodes())),
		Edges: make([]graphEdgeView, 0, len(mixin.Edges())),
	}

	for _, node := range mixin.Nodes() {
		graph.Nodes = append(graph.Nodes, graphNodeView{
			ID:       node.ID.Get(),
			Label:    node.Label.Get(),
			Shape:    string(node.Shape),
			Metadata: node.Metadata,
		})
	}

	for _, edge := range mixin.Edges() {
		graph.Edges = append(graph.Edges, graphEdgeView{
			From:  edge.From.Get(),
			To:    edge.To.Get(),
			Label: brandedEdgeLabel(edge.Label),
		})
	}

	return graph
}
