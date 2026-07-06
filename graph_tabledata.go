package output

import (
	"fmt"
	"strings"
)

// GraphNodeLabelFunc is a function that formats a cell value with its header into a label.
type GraphNodeLabelFunc func(header, cell string) string

// DefaultGraphNodeLabel returns a label in the format "header: cell".
func DefaultGraphNodeLabel(header, cell string) string {
	return fmt.Sprintf("%s: %s", header, cell)
}

// TreeNodeIDFunc resolves a TreeNode's ID for a specific graph format.
type TreeNodeIDFunc func(*TreeNode) string

// NodesFromTable creates GraphNodes from Table using the provided label function.
func NodesFromTable(data *Table, labelFn GraphNodeLabelFunc) []GraphNode {
	if data == nil {
		return nil
	}

	nodes := make([]GraphNode, 0, len(data.Rows))

	for i, row := range data.Rows {
		var labelParts []string

		for j, cell := range row {
			if j < len(data.Headers) {
				labelParts = append(labelParts, labelFn(data.Headers[j], cell))
			} else {
				labelParts = append(labelParts, cell)
			}
		}

		label := strings.Join(labelParts, "\n")

		//nolint:exhaustruct // Uses defaults for optional fields
		nodes = append(nodes, GraphNode{
			ID:    NewBrandedID[GraphNodeIDBrand](fmt.Sprintf("row%d", i)),
			Label: NewBrandedID[GraphNodeLabelBrand](label),
		})
	}

	return nodes
}

// AddRowEdges adds edges from data.CreateRowEdges() to the graph.
func (m *GraphRendererState) AddRowEdges(data *Table) {
	for _, edge := range data.CreateRowEdges() {
		//nolint:exhaustruct // Uses defaults for optional fields
		m.edges = append(m.edges, GraphEdge{
			From: edge.From,
			To:   edge.To,
		})
	}
}

// SetNodesFromTable creates nodes from Table, applies per-node modifications,
// adds them to the graph, and adds row edges.
func (m *GraphRendererState) SetNodesFromTable(
	data *Table,
	modifyNode func(i int, n *GraphNode),
) {
	nodes := NodesFromTable(data, DefaultGraphNodeLabel)

	for i := range nodes {
		if modifyNode != nil {
			modifyNode(i, &nodes[i])
		}
	}

	m.nodes = append(m.nodes, nodes...)
	m.AddRowEdges(data)
}
