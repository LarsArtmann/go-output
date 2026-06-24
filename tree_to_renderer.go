package output

// TreeToRenderer runs the standard tree entrypoint prelude used by every
// sub-module's public XFromTree: instantiate the renderer, skip population
// when root is nil, otherwise call the renderer-specific addNodes with the
// given zero parent ID (every XFromTree starts a fresh subtree). Centralising
// this in root means each sub-module's entrypoint is a single line and the
// constructor + nil-check + populate prelude lives in exactly one place.
//
// The ID type is renderer-specific: graph and plantuml use TreeNodeID; d2 and
// mermaid use plain string. Both shapes share this single generic helper.
func TreeToRenderer[R, ID any](
	newRenderer func() R,
	addNodes func(R, *TreeNode, ID),
	root *TreeNode,
) R {
	r := newRenderer()
	if root != nil {
		addNodes(r, root, *new(ID))
	}

	return r
}
