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
// Canonical import path: github.com/larsartmann/go-output (root).
// The d2 module re-exports this as d2.D2NodeID for convenience.
type D2NodeID = id.ID[D2NodeIDBrand, string]

// D2NodeLabel is a branded identifier for D2 diagram node labels.
// Canonical import path: github.com/larsartmann/go-output (root).
// The d2 module re-exports this as d2.D2NodeLabel for convenience.
type D2NodeLabel = id.ID[D2NodeLabelBrand, string]

// TreeNodeID is a branded identifier for tree nodes.
type TreeNodeID = id.ID[TreeNodeIDBrand, string]

// TreeNodeLabel is a branded identifier for tree node labels.
type TreeNodeLabel = id.ID[TreeNodeLabelBrand, string]

// GraphNodeID is a branded identifier for graph nodes.
type GraphNodeID = id.ID[GraphNodeIDBrand, string]

// GraphNodeLabel is a branded identifier for graph node labels.
type GraphNodeLabel = id.ID[GraphNodeLabelBrand, string]

// D2NodeIDBrand is the brand type for D2 node IDs. Lives in root so the
// type alias D2NodeID is visible to both root and d2/ callers without a
// circular import. Split-brain finding m6: root is the canonical home;
// d2.D2NodeID is a convenience re-export (type alias), not a second definition.
type D2NodeIDBrand struct{}

// D2NodeLabelBrand is the brand type for D2 node labels. See D2NodeIDBrand
// for why brand types live in root rather than d2/.
type D2NodeLabelBrand struct{}

// TreeNodeIDBrand is the brand type for tree node IDs.
type TreeNodeIDBrand struct{}

// TreeNodeLabelBrand is the brand type for tree node labels.
type TreeNodeLabelBrand struct{}

// GraphNodeIDBrand is the brand type for graph node IDs.
type GraphNodeIDBrand struct{}

// GraphNodeLabelBrand is the brand type for graph node labels.
type GraphNodeLabelBrand struct{}
