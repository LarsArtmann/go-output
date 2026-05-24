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

// NodesFromTableData creates GraphNodes from TableData using the provided label function.
func NodesFromTableData(data *TableData, labelFn GraphNodeLabelFunc) []GraphNode {
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
func (m *GraphRendererMixin) AddRowEdges(data *TableData) {
	for _, edge := range data.CreateRowEdges() {
		//nolint:exhaustruct // Uses defaults for optional fields
		m.edges = append(m.edges, GraphEdge{
			From: NewBrandedID[GraphNodeIDBrand](edge.From),
			To:   NewBrandedID[GraphNodeIDBrand](edge.To),
		})
	}
}

// SetNodesFromTableData creates nodes from TableData, applies per-node modifications,
// adds them to the graph, and adds row edges.
func (m *GraphRendererMixin) SetNodesFromTableData(
	data *TableData,
	modifyNode func(i int, n *GraphNode),
) {
	nodes := NodesFromTableData(data, DefaultGraphNodeLabel)

	for i := range nodes {
		modifyNode(i, &nodes[i])
	}

	m.nodes = append(m.nodes, nodes...)
	m.AddRowEdges(data)
}
