package output

import (
	"fmt"
	"io"
)

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
	switch {
	case verb == 'v' && s.Flag('#'):
		_, _ = fmt.Fprintf(s, "BrandedID{%q}", id.value)
	default:
		_, _ = io.WriteString(s, id.value)
	}
}

// D2NodeIDBrand is the brand type for D2 node IDs.
type D2NodeIDBrand struct{}

// D2NodeID is a branded identifier for D2 diagram nodes.
type D2NodeID = BrandedID[D2NodeIDBrand]

// D2NodeLabelBrand is the brand type for D2 node labels.
type D2NodeLabelBrand struct{}

// D2NodeLabel is a branded identifier for D2 diagram node labels.
type D2NodeLabel = BrandedID[D2NodeLabelBrand]

// TreeNodeIDBrand is the brand type for tree node IDs.
type TreeNodeIDBrand struct{}

// TreeNodeID is a branded identifier for tree nodes.
type TreeNodeID = BrandedID[TreeNodeIDBrand]

// TreeNodeLabelBrand is the brand type for tree node labels.
type TreeNodeLabelBrand struct{}

// TreeNodeLabel is a branded identifier for tree node labels.
type TreeNodeLabel = BrandedID[TreeNodeLabelBrand]

// GraphNodeIDBrand is the brand type for graph node IDs.
type GraphNodeIDBrand struct{}

// GraphNodeID is a branded identifier for graph nodes.
type GraphNodeID = BrandedID[GraphNodeIDBrand]

// GraphNodeLabelBrand is the brand type for graph node labels.
type GraphNodeLabelBrand struct{}

// GraphNodeLabel is a branded identifier for graph node labels.
type GraphNodeLabel = BrandedID[GraphNodeLabelBrand]
