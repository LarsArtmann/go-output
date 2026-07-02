package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/larsartmann/go-output/nom"
)

// exportDOTToTemp writes a DOT representation of the dependency tree to a
// temporary file and returns the file path. The DOT format is suitable for
// Graphviz rendering (dot -Tsvg output.dot -o output.svg).
func exportDOTToTemp(tree *nom.DependencyTree) string {
	if tree == nil {
		return ""
	}

	tmpDir := os.TempDir()
	path := filepath.Join(tmpDir, "go-output-dag.dot")

	content := treeToDOT(tree)

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return ""
	}

	return path
}

// treeToDOT converts the dependency tree to a DOT graph string.
func treeToDOT(tree *nom.DependencyTree) string {
	var b []byte

	b = append(b, "digraph dag {\n"...)
	b = append(b, "  rankdir=TB;\n"...)
	b = append(b, "  node [shape=box, style=rounded];\n"...)

	for _, node := range tree.AllNodes() {
		label := string(node.ID)
		b = append(b, fmt.Sprintf("  %q [label=%q];\n", string(node.ID), label)...)

		for _, dep := range node.Deps {
			b = append(b, fmt.Sprintf("  %q -> %q;\n", string(dep), string(node.ID))...)
		}
	}

	b = append(b, "}\n"...)

	return string(b)
}
