package nom
import (
	"sync"
)
// Clear removes all nodes from the tree.
func (dt *DependencyTree) Clear() {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	dt.nodes = make(map[ActivityID]*TreeNode)
	dt.roots = make([]*TreeNode, 0)
	dt.order = make([]ActivityID, 0)
	dt.loaded = false
	dt.buildOnce = sync.Once{}
}
