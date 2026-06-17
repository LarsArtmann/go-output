package d2

import "github.com/larsartmann/go-output"

// newD2NodeID creates a D2NodeID with the given value. Test helper to
// eliminate `output.NewBrandedID[output.D2NodeIDBrand](...)` boilerplate.
func newD2NodeID(id string) output.D2NodeID {
	return output.NewBrandedID[output.D2NodeIDBrand](id)
}

// newD2NodeLabel creates a D2NodeLabel with the given value. Test helper
// for the same reason as newD2NodeID.
func newD2NodeLabel(label string) output.D2NodeLabel {
	return output.NewBrandedID[output.D2NodeLabelBrand](label)
}
