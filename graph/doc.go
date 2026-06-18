// Package graph provides DOT (Graphviz) and Mermaid diagram renderers for
// graph-structured data.
//
// Both DOTRenderer and MermaidRenderer implement the output.GraphRenderer
// interface and can be populated from output.TableData (via SetNodesFromTableData)
// or output.TreeNode (via AddTreeNodes helper).
//
// # Branded IDs
//
// GraphNodeID and GraphNodeLabel are type aliases re-exported from the root
// package (output.GraphNodeID / output.GraphNodeLabel). The canonical import path
// is the root package; graph.GraphNodeID exists only as an ergonomic convenience.
// There is exactly one definition. See split-brain m6.
//
// # Quick Start
//
//	renderer := graph.NewDOTRenderer()
//	renderer.SetNodes([]output.GraphNode{*output.NewGraphNode("a", "Node A")})
//	renderer.SetEdges([]output.GraphEdge{*output.NewGraphEdge("a", "b")})
//	output, _ := renderer.Render()
package graph
