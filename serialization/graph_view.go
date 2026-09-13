package serialization

import (
	"github.com/larsartmann/go-output"
)

type graphView struct {
	Nodes []graphNodeView `json:"nodes" yaml:"nodes" toml:"nodes"`
	Edges []graphEdgeView `json:"edges" yaml:"edges" toml:"edges"`
}

type graphNodeView struct {
	ID       string            `json:"id"                 yaml:"id"                 toml:"id"`
	Label    string            `json:"label"              yaml:"label"              toml:"label"`
	Shape    string            `json:"shape,omitempty"    yaml:"shape,omitempty"    toml:"shape,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty" toml:"metadata,omitempty"`
}

type graphEdgeView struct {
	From  string `json:"from"            yaml:"from"            toml:"from"`
	To    string `json:"to"              yaml:"to"              toml:"to"`
	Label string `json:"label,omitempty" yaml:"label,omitempty" toml:"label,omitempty"`
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
