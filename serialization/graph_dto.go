package serialization

import (
	"github.com/larsartmann/go-output"
)

type graphDTO struct {
	Nodes []graphNodeDTO `json:"nodes" yaml:"nodes"`
	Edges []graphEdgeDTO `json:"edges" yaml:"edges"`
}

type graphNodeDTO struct {
	ID       string            `json:"id"                 yaml:"id"`
	Label    string            `json:"label"              yaml:"label"`
	Shape    string            `json:"shape,omitempty"    yaml:"shape,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

type graphEdgeDTO struct {
	From  string `json:"from"            yaml:"from"`
	To    string `json:"to"              yaml:"to"`
	Label string `json:"label,omitempty" yaml:"label,omitempty"`
}

func buildGraphDTO(mixin output.GraphRendererMixin) graphDTO {
	graph := graphDTO{
		Nodes: make([]graphNodeDTO, 0, len(mixin.Nodes())),
		Edges: make([]graphEdgeDTO, 0, len(mixin.Edges())),
	}

	for _, node := range mixin.Nodes() {
		graph.Nodes = append(graph.Nodes, graphNodeDTO{
			ID:       node.ID.Get(),
			Label:    node.Label.Get(),
			Shape:    string(node.Shape),
			Metadata: node.Metadata,
		})
	}

	for _, edge := range mixin.Edges() {
		graph.Edges = append(graph.Edges, graphEdgeDTO{
			From:  edge.From.Get(),
			To:    edge.To.Get(),
			Label: brandedEdgeLabel(edge.Label),
		})
	}

	return graph
}
