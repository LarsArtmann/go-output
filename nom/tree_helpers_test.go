package nom

import "testing"

// assertChildParentID fails the test if the child node's parent ID does not
// match the expected value. Used by dependency-tree tests to verify auto- and
// explicit-parent links.
func assertChildParentID(t *testing.T, child *ActivityNode, wantParentID string) {
	t.Helper()

	if child == nil || child.Parent == nil {
		t.Fatalf("child or parent is nil")

		return
	}

	if got := child.Parent.ID.Get(); got != wantParentID {
		t.Errorf("child's parent ID = %q, want %q", got, wantParentID)
	}
}
