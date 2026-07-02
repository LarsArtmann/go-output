package daghtml

// Node represents a single node in the DAG visualization.
type Node struct {
	// ID is the unique identifier for this node. Edge From/To values reference this.
	ID string `json:"id"`

	// Label is the display text shown inside the node rectangle.
	// May include emoji or icon characters prefixed to the name.
	Label string `json:"label"`

	// Color is the CSS color or var() reference used for the node's border
	// and left accent bar. Example: "var(--success)" or "#e8a838".
	Color string `json:"color"`

	// Tooltip is optional multi-line text shown in the SVG <title> element
	// (appears as a native browser tooltip on hover). Use " | " as a field
	// separator to match the existing visual convention.
	Tooltip string `json:"tooltip,omitempty"`

	// Error, when true, renders a small red dot in the top-right corner of
	// the node to flag failed or problematic nodes.
	Error bool `json:"error,omitempty"`
}

// Edge represents a directed edge in the DAG. The visual direction is From → To,
// meaning From is the upstream (dependency) and To is the downstream (dependent).
type Edge struct {
	// From is the ID of the source (upstream) node.
	From string `json:"from"`

	// To is the ID of the target (downstream) node.
	To string `json:"to"`
}

// DAG is the complete graph data model consumed by the HTML renderer.
// Nodes with duplicate IDs are deduplicated (first occurrence wins).
// Edges referencing unknown node IDs are silently dropped.
// Duplicate edges (same From/To pair) are deduplicated.
type DAG struct {
	// Nodes is the list of all nodes in the graph.
	Nodes []Node `json:"nodes"`

	// Edges is the list of directed edges connecting nodes.
	Edges []Edge `json:"edges"`
}

// NodeCount returns the number of nodes in the DAG.
func (d DAG) NodeCount() int { return len(d.Nodes) }

// EdgeCount returns the number of edges in the DAG.
func (d DAG) EdgeCount() int { return len(d.Edges) }

// IsEmpty returns true if the DAG has no nodes.
func (d DAG) IsEmpty() bool { return len(d.Nodes) == 0 }
