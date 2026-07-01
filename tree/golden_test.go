package tree

import (
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/larsartmann/go-output"
)

func TestGolden_Tree_SimpleHierarchy(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Project")
	src := output.NewTreeNode("src", "src")
	src.AddChild(output.NewTreeNode("main", "main.go"))
	src.AddChild(output.NewTreeNode("utils", "utils.go"))
	docs := output.NewTreeNode("docs", "docs")
	docs.AddChild(output.NewTreeNode("readme", "README.md"))
	docs.AddChild(output.NewTreeNode("changelog", "CHANGELOG.md"))
	root.AddChild(src)
	root.AddChild(docs)
	root.AddChild(output.NewTreeNode("gomod", "go.mod"))

	r := NewASCIITreeRenderer()
	r.SetColorMode(output.ColorModeNever)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_Tree_DeepNesting(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "root")
	l1 := output.NewTreeNode("l1", "level1")
	l2 := output.NewTreeNode("l2", "level2")
	l3 := output.NewTreeNode("l3", "level3")
	l4 := output.NewTreeNode("l4", "level4")
	l3.AddChild(l4)
	l2.AddChild(l3)
	l1.AddChild(l2)
	root.AddChild(l1)

	r := NewASCIITreeRenderer()
	r.SetColorMode(output.ColorModeNever)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_Tree_SingleNode(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("lonely", "lonely")

	r := NewASCIITreeRenderer()
	r.SetColorMode(output.ColorModeNever)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_Tree_MixedBranching(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("pipe", "build-pipeline")
	compile := output.NewTreeNode("compile", "compile")
	compile.AddChild(output.NewTreeNode("lint", "lint"))
	compile.AddChild(output.NewTreeNode("vet", "vet"))
	testNode := output.NewTreeNode("test", "test")
	testNode.AddChild(output.NewTreeNode("unit", "unit"))
	testNode.AddChild(output.NewTreeNode("integ", "integration"))
	testNode.AddChild(output.NewTreeNode("e2e", "e2e"))
	compile.AddChild(testNode)
	root.AddChild(compile)
	root.AddChild(output.NewTreeNode("pkg", "package"))
	deploy := output.NewTreeNode("deploy", "deploy")
	deploy.AddChild(output.NewTreeNode("staging", "staging"))
	deploy.AddChild(output.NewTreeNode("prod", "production"))
	root.AddChild(deploy)

	r := NewASCIITreeRenderer()
	r.SetColorMode(output.ColorModeNever)
	r.SetRoot(root)

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}
