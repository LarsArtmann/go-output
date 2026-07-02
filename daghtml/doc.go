// Package daghtml renders directed acyclic graphs (DAGs) as self-contained,
// interactive HTML pages with an embedded Sugiyama layered layout algorithm.
//
// The visualization includes:
//   - SVG-based graph with cubic bezier edges
//   - Sugiyama layout (Kahn's rank assignment + 4-pass barycenter crossing reduction)
//   - Pan, zoom (wheel + pinch), and click-to-highlight connected subgraph
//   - Touch support (1-finger pan, 2-finger pinch-zoom)
//   - Layer depth labels
//   - Node tooltips with arbitrary metadata
//   - Dark theme with CSS custom properties (overridable)
//
// # Quick start
//
//	dag := daghtml.DAG{
//	    Nodes: []daghtml.Node{
//	        {ID: "a", Label: "Fetch", Color: "var(--success)"},
//	        {ID: "b", Label: "Parse", Color: "var(--accent)"},
//	        {ID: "c", Label: "Save", Color: "var(--error)", Error: true},
//	    },
//	    Edges: []daghtml.Edge{
//	        {From: "a", To: "b"},
//	        {From: "b", To: "c"},
//	    },
//	}
//	html, err := daghtml.Render(dag, daghtml.WithTitle("My Pipeline"))
//
// # Embedding in a host page
//
// For embedding just the graph section in an existing dashboard, use
// [GraphHTML] instead of [Render]. The host page must define the CSS
// custom properties (--bg, --surface, --accent, etc.) or include the
// graph CSS via [StyleSheet].
package daghtml
