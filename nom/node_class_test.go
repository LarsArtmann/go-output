package nom

import "testing"

func TestActivityNode_Class(t *testing.T) {
	t.Parallel()

	tree := NewDependencyTree()
	tree.AddActivity(ActivityID("root"), nil)
	tree.AddActivity(ActivityID("mid"), []ActivityID{"root"})
	tree.AddActivity(ActivityID("leaf"), []ActivityID{"mid"})

	if err := tree.Build(); err != nil {
		t.Fatalf("build: %v", err)
	}

	tests := []struct {
		id    ActivityID
		class NodeClass
	}{
		{ActivityID("root"), NodeClassRoot},
		{ActivityID("mid"), NodeClassTwig},
		{ActivityID("leaf"), NodeClassLeaf},
	}

	for _, tc := range tests {
		node := tree.GetNode(tc.id)
		if node == nil {
			t.Fatalf("node %s not found", tc.id)
		}

		if got := node.Class(); got != tc.class {
			t.Errorf("%s: class = %s, want %s", tc.id, got, tc.class)
		}
	}
}
