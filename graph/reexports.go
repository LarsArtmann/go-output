package graph

import "github.com/larsartmann/go-output"

// GraphNodeID is a branded identifier for graph nodes.
type GraphNodeID = output.GraphNodeID

// GraphNodeLabel is a branded identifier for graph node labels.
type GraphNodeLabel = output.GraphNodeLabel

// NewGraphNodeID creates a new branded graph node ID.
//
// Deprecated: Use output.NewBrandedID[output.GraphNodeIDBrand](id) directly.
func NewGraphNodeID(id string) GraphNodeID {
	return output.NewBrandedID[output.GraphNodeIDBrand](id)
}

// NewGraphNodeLabel creates a new branded graph node label.
//
// Deprecated: Use output.NewBrandedID[output.GraphNodeLabelBrand](label) directly.
func NewGraphNodeLabel(label string) GraphNodeLabel {
	return output.NewBrandedID[output.GraphNodeLabelBrand](label)
}
