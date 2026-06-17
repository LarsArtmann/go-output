package output

import (
	id "github.com/larsartmann/go-branded-id"
)

// BrandedID re-export for backward compatibility.
// Use id.ID[Brand, string] for new code.
type BrandedID[Brand any] = id.ID[Brand, string]

// NewBrandedID creates a new branded ID from a string value.
func NewBrandedID[Brand any](value string) id.ID[Brand, string] {
	return id.NewID[Brand](value)
}

// D2NodeID is a branded identifier for D2 diagram nodes.
//
// NOTE(split-brain m6): This type is re-exported in d2/d2.go as `type D2NodeID = output.D2NodeID`.
// The brand struct D2NodeIDBrand also lives here in root, not in d2/ where it logically belongs.
// In the next minor version, consider moving D2-specific brands to the d2 module.
type D2NodeID = id.ID[D2NodeIDBrand, string]

// D2NodeLabel is a branded identifier for D2 diagram node labels.
type D2NodeLabel = id.ID[D2NodeLabelBrand, string]

// TreeNodeID is a branded identifier for tree nodes.
type TreeNodeID = id.ID[TreeNodeIDBrand, string]

// TreeNodeLabel is a branded identifier for tree node labels.
type TreeNodeLabel = id.ID[TreeNodeLabelBrand, string]

// GraphNodeID is a branded identifier for graph nodes.
type GraphNodeID = id.ID[GraphNodeIDBrand, string]

// GraphNodeLabel is a branded identifier for graph node labels.
type GraphNodeLabel = id.ID[GraphNodeLabelBrand, string]

// D2NodeIDBrand is the brand type for D2 node IDs.
type D2NodeIDBrand struct{}

// D2NodeLabelBrand is the brand type for D2 node labels.
type D2NodeLabelBrand struct{}

// TreeNodeIDBrand is the brand type for tree node IDs.
type TreeNodeIDBrand struct{}

// TreeNodeLabelBrand is the brand type for tree node labels.
type TreeNodeLabelBrand struct{}

// GraphNodeIDBrand is the brand type for graph node IDs.
type GraphNodeIDBrand struct{}

// GraphNodeLabelBrand is the brand type for graph node labels.
type GraphNodeLabelBrand struct{}
