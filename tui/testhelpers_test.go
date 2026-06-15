package tui

import (
	"fmt"

	"github.com/larsartmann/go-output/nom"
)

func newTestTree(nodeCount int) *nom.DependencyTree {
	tree := nom.NewDependencyTree()

	tree.AddActivity(nom.ActivityID("root"), "Root", nil)

	for i := 1; i < nodeCount; i++ {
		id := nom.ActivityID(fmt.Sprintf("step-%d", i))
		name := fmt.Sprintf("Step %d", i)
		tree.AddActivity(id, name, []nom.ActivityID{"root"})
	}

	_ = tree.GetRootNodes()

	return tree
}
