package output

import "fmt"

// BrandedID is a type-safe identifier with a phantom brand type.
// This provides compile-time safety to prevent mixing different ID types.
type BrandedID[Brand any] struct {
	value string
}

// NewBrandedID creates a new branded ID from a string value.
func NewBrandedID[Brand any](value string) BrandedID[Brand] {
	return BrandedID[Brand]{value: value}
}

// Get returns the underlying string value.
func (id BrandedID[B]) Get() string {
	return id.value
}

// String returns the string representation.
func (id BrandedID[B]) String() string {
	return id.value
}

// IsEmpty returns true if the ID is empty.
func (id BrandedID[B]) IsEmpty() bool {
	return id.value == ""
}

// MarshalText implements encoding.TextMarshaler.
func (id BrandedID[B]) MarshalText() ([]byte, error) {
	return []byte(id.value), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (id *BrandedID[B]) UnmarshalText(text []byte) error {
	id.value = string(text)

	return nil
}

// Format implements fmt.Formatter.
func (id BrandedID[B]) Format(s fmt.State, verb rune) {
	format := "%s"

	if verb == 'v' {
		if s.Flag('#') {
			format = "%s{%q}"
		}
	}

	_, _ = fmt.Fprintf(s, format, id.value)
}

// D2NodeIDBrand is the brand type for D2 node IDs.
type D2NodeIDBrand struct{}

// D2NodeID is a branded identifier for D2 diagram nodes.
type D2NodeID = BrandedID[D2NodeIDBrand]

// D2NodeLabelBrand is the brand type for D2 node labels.
type D2NodeLabelBrand struct{}

// D2NodeLabel is a branded identifier for D2 diagram node labels.
type D2NodeLabel = BrandedID[D2NodeLabelBrand]

// D2EdgeFromBrand is the brand type for D2 edge source IDs.
type D2EdgeFromBrand struct{}

// D2EdgeToBrand is the brand type for D2 edge target IDs.
type D2EdgeToBrand struct{}

// D2EdgeLabelBrand is the brand type for D2 edge labels.
type D2EdgeLabelBrand struct{}

// D2EdgeID is a branded identifier for D2 diagram edges.
type D2EdgeID = BrandedID[D2EdgeFromBrand]

// DOTGraphIDBrand is the brand type for DOT graph IDs.
type DOTGraphIDBrand struct{}

// DOTGraphID is a branded identifier for DOT/Graphviz graphs.
type DOTGraphID = BrandedID[DOTGraphIDBrand]

// DOTNodeIDBrand is the brand type for DOT node IDs.
type DOTNodeIDBrand struct{}

// DOTNodeID is a branded identifier for DOT diagram nodes.
type DOTNodeID = BrandedID[DOTNodeIDBrand]

// DOTEdgeFromBrand is the brand type for DOT edge source IDs.
type DOTEdgeFromBrand struct{}

// DOTEdgeToBrand is the brand type for DOT edge target IDs.
type DOTEdgeToBrand struct{}

// DOTEdgeLabelBrand is the brand type for DOT edge labels.
type DOTEdgeLabelBrand struct{}

// TreeNodeIDBrand is the brand type for tree node IDs.
type TreeNodeIDBrand struct{}

// TreeNodeID is a branded identifier for tree nodes.
type TreeNodeID = BrandedID[TreeNodeIDBrand]

// TreeNodeLabelBrand is the brand type for tree node labels.
type TreeNodeLabelBrand struct{}

// TreeNodeLabel is a branded identifier for tree node labels.
type TreeNodeLabel = BrandedID[TreeNodeLabelBrand]

// TreeParentIDBrand is the brand type for tree parent IDs.
type TreeParentIDBrand struct{}

// TreeParentID is a branded identifier for tree parent nodes.
type TreeParentID = BrandedID[TreeParentIDBrand]

// GraphNodeIDBrand is the brand type for graph node IDs.
type GraphNodeIDBrand struct{}

// GraphNodeID is a branded identifier for graph nodes.
type GraphNodeID = BrandedID[GraphNodeIDBrand]

// GraphNodeLabelBrand is the brand type for graph node labels.
type GraphNodeLabelBrand struct{}

// GraphNodeLabel is a branded identifier for graph node labels.
type GraphNodeLabel = BrandedID[GraphNodeLabelBrand]

// GraphEdgeFromBrand is the brand type for graph edge source IDs.
type GraphEdgeFromBrand struct{}

// GraphEdgeToBrand is the brand type for graph edge target IDs.
type GraphEdgeToBrand struct{}

// GraphEdgeLabelBrand is the brand type for graph edge labels.
type GraphEdgeLabelBrand struct{}

// MermaidNodeIDBrand is the brand type for Mermaid node IDs.
type MermaidNodeIDBrand struct{}

// MermaidNodeID is a branded identifier for Mermaid diagram nodes.
type MermaidNodeID = BrandedID[MermaidNodeIDBrand]

// MermaidParentIDBrand is the brand type for Mermaid parent IDs.
type MermaidParentIDBrand struct{}

// MermaidParentID is a branded identifier for Mermaid parent nodes.
type MermaidParentID = BrandedID[MermaidParentIDBrand]

// HTMLTitleBrand is the brand type for HTML title IDs.
type HTMLTitleBrand struct{}

// HTMLTitle is a branded identifier for HTML titles.
type HTMLTitle = BrandedID[HTMLTitleBrand]
