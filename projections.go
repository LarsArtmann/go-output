package output

// TableToGraph converts a Table into a Graph by creating one node per row
// and connecting consecutive rows with directed edges. Each node is labeled
// with all cell values from its row. This is a pure projection — the input
// Table is not modified.
func TableToGraph(data *Table, labelFn ...GraphNodeLabelFunc) Graph {
	if data == nil || len(data.Rows) == 0 {
		return Graph{}
	}

	fn := DefaultGraphNodeLabel
	if len(labelFn) > 0 && labelFn[0] != nil {
		fn = labelFn[0]
	}

	b := NewGraphBuilder()
	b.SetNodes(NodesFromTable(data, fn))
	b.AddRowEdges(data)

	return b.Build()
}

// GraphToTree converts a Graph into a TreeNode tree. The first node with no
// incoming edges becomes the root; outgoing edges define parent→child
// relationships. If the graph contains cycles, each node is visited at most
// once (cycle guard). If the graph has multiple roots, only the first root's
// subtree is returned. This is a pure projection — the input Graph is not
// modified.
func GraphToTree(g Graph) *TreeNode {
	nodes := g.Nodes()
	edges := g.Edges()

	if len(nodes) == 0 {
		return nil
	}

	children := make(map[string][]string)
	hasIncoming := make(map[string]bool, len(nodes))

	for _, node := range nodes {
		hasIncoming[node.ID.Get()] = false
	}

	for _, edge := range edges {
		from := edge.From.Get()
		to := edge.To.Get()
		children[from] = append(children[from], to)
		hasIncoming[to] = true
	}

	var rootID string

	for _, node := range nodes {
		if !hasIncoming[node.ID.Get()] {
			rootID = node.ID.Get()
			break
		}
	}

	if rootID == "" {
		rootID = nodes[0].ID.Get()
	}

	nodeMap := make(map[string]*GraphNode, len(nodes))
	for i := range nodes {
		nodeMap[nodes[i].ID.Get()] = &nodes[i]
	}

	visited := make(map[string]bool)

	var build func(id string) *TreeNode

	build = func(id string) *TreeNode {
		if visited[id] {
			return nil
		}

		visited[id] = true

		node := nodeMap[id]
		if node == nil {
			return nil
		}

		treeNode := NewTreeNode(id, node.Label.Get())

		for _, childID := range children[id] {
			if child := build(childID); child != nil {
				treeNode.AddChild(child)
			}
		}

		return treeNode
	}

	return build(rootID)
}

// GraphToTable converts a Graph into a Table with columns "ID" and "Label",
// one row per node. Edges are not represented in the table format. This is
// a pure projection — the input Graph is not modified.
func GraphToTable(g Graph) *Table {
	nodes := g.Nodes()
	if len(nodes) == 0 {
		return nil
	}

	t := NewTable([]string{"ID", "Label"})

	for _, node := range nodes {
		t.AddRow([]string{node.ID.Get(), node.Label.Get()})
	}

	return t
}
