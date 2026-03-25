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

// Brand types for D2 diagram elements.
type (
	D2NodeIDBrand    struct{}
	D2NodeID         = BrandedID[D2NodeIDBrand]
	D2NodeLabelBrand struct{}
	D2NodeLabel      = BrandedID[D2NodeLabelBrand]
	D2EdgeFromBrand  struct{}
	D2EdgeToBrand    struct{}
	D2EdgeLabelBrand struct{}
	D2EdgeID         = BrandedID[D2EdgeFromBrand]
)

// Brand types for DOT/Graphviz diagram elements.
type (
	DOTGraphIDBrand   struct{}
	DOTGraphID        = BrandedID[DOTGraphIDBrand]
	DOTNodeIDBrand    struct{}
	DOTNodeID         = BrandedID[DOTNodeIDBrand]
	DOTEdgeFromBrand  struct{}
	DOTEdgeToBrand    struct{}
	DOTEdgeLabelBrand struct{}
)

// Brand types for TreeNode elements.
type (
	TreeNodeIDBrand    struct{}
	TreeNodeID         = BrandedID[TreeNodeIDBrand]
	TreeNodeLabelBrand struct{}
	TreeNodeLabel      = BrandedID[TreeNodeLabelBrand]
	TreeParentIDBrand  struct{}
	TreeParentID       = BrandedID[TreeParentIDBrand]
)

// Brand types for GraphNode/GraphEdge elements.
type (
	GraphNodeIDBrand    struct{}
	GraphNodeID         = BrandedID[GraphNodeIDBrand]
	GraphNodeLabelBrand struct{}
	GraphNodeLabel      = BrandedID[GraphNodeLabelBrand]
	GraphEdgeFromBrand  struct{}
	GraphEdgeToBrand    struct{}
	GraphEdgeLabelBrand struct{}
)

// Brand types for Mermaid diagram elements.
type (
	MermaidNodeIDBrand   struct{}
	MermaidNodeID        = BrandedID[MermaidNodeIDBrand]
	MermaidParentIDBrand struct{}
	MermaidParentID      = BrandedID[MermaidParentIDBrand]
)

// Brand types for HTML rendering.
type (
	HTMLTitleBrand struct{}
	HTMLTitle      = BrandedID[HTMLTitleBrand]
)
