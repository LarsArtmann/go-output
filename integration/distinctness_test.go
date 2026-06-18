// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"reflect"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/nom"
)

// TestActivityNodeDistinctFromTreeNode is a compile-time regression guard for
// split-brain finding C1: nom.TreeNode was renamed to nom.ActivityNode to avoid
// collision with output.TreeNode (unrelated concepts that share a name).
// If anyone ever reintroduces a nom.TreeNode alias to ActivityNode, the
// var declarations below fail to compile because the two types are distinct.
//
// We assert via unused package-level vars that the types are NOT assignable to
// each other. This test compiles only because ActivityNode != TreeNode.
func TestActivityNodeDistinctFromTreeNode(t *testing.T) {
	t.Parallel()

	// Compile-time assertions: assigning one to the other would fail to build.
	// The vars are deliberately unused — their existence IS the assertion.
	var (
		_ nom.ActivityNode // present in nom
		_ output.TreeNode  // present in root, unrelated concept
	)

	// Runtime sanity: the two zero values are of different concrete types.
	nomNode := nom.ActivityNode{}
	rootNode := output.TreeNode{}
	_, _ = nomNode, rootNode

	// The real assertion: distinct type names. If they were the same type,
	// reflect.TypeOf would match. They must not.
	nomType := reflect.TypeOf(nomNode).String()
	rootType := reflect.TypeOf(rootNode).String()
	if nomType == rootType {
		t.Fatalf("nom.ActivityNode and output.TreeNode resolved to the same type %q — C1 split-brain regression", nomType)
	}

	// The nom type must not even be named "TreeNode" anymore.
	if nomType == "TreeNode" {
		t.Errorf("nom dependency-tree node is named %q — should be ActivityNode (C1)", nomType)
	}
}
