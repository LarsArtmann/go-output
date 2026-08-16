package tree

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestASCIITreeRenderer(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()

	// Test empty tree
	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != "" {
		t.Errorf("Empty tree should return empty string, got: %s", got)
	}

	// Build a test tree
	root := output.NewTreeNode("root", "Root")
	child1 := output.NewTreeNode("child1", "Child 1")
	child2 := output.NewTreeNode("child2", "Child 2")
	child1.AddChild(output.NewTreeNode("leaf1", "Leaf 1"))
	child1.AddChild(output.NewTreeNode("leaf2", "Leaf 2"))
	root.AddChild(child1)
	root.AddChild(child2)

	renderer.SetRoot(root)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "Root", "Output should contain 'Root'")
	assertContains(t, out, "Child 1", "Output should contain 'Child 1'")
	assertContains(t, out, "Child 2", "Output should contain 'Child 2'")
	assertContains(t, out, "Leaf 1", "Output should contain 'Leaf 1'")
}

func TestASCIITreeRendererWithMetadata(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()

	node := output.NewTreeNode("node", "Node with Meta")
	node.Metadata["key"] = "value"
	node.Metadata["count"] = "42"

	renderer.SetRoot(node)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "Node with Meta", "Output should contain node label")
	assertContains(t, out, "key: value", "Output should contain metadata")
}

func TestTreeRendererFromTable(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Age", "City"})
	data.AddRow([]string{"Alice", "30", "NYC"})
	data.AddRow([]string{"Bob", "25", "LA"})

	renderer := TreeRendererFromTable(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, out, "Data", "Output should contain 'Data' as root")
	assertContains(t, out, "Headers", "Output should contain 'Headers' section")
	assertContains(t, out, "Rows", "Output should contain 'Rows' section")
	assertContains(t, out, "Alice", "Output should contain 'Alice' row")
}

func TestTreeColorModeNever(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()
	renderer.SetColorMode(output.ColorModeNever)

	root := output.NewTreeNode("root", "Root")
	root.Metadata["key"] = "value"
	renderer.SetRoot(root)

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, "\x1b[") {
		t.Errorf("ColorModeNever should produce no ANSI codes, got: %q", got)
	}

	assertContains(t, got, "Root", "should contain label even without colors")
	assertContains(t, got, "key: value", "should contain metadata even without colors")
}

func TestTreeColorModeAlways(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()
	renderer.SetColorMode(output.ColorModeAlways)

	root := output.NewTreeNode("root", "Root")
	child := output.NewTreeNode("child", "Child")
	root.AddChild(child)
	renderer.SetRoot(root)

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("ColorModeAlways should produce ANSI codes, got: %q", got)
	}

	assertContains(t, got, "Root", "should contain label")
	assertContains(t, got, "Child", "should contain child label")
}

func TestTreeColorModeDefault(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()
	if renderer.colorMode != output.ColorModeAuto {
		t.Errorf("default ColorMode = %v, want %v", renderer.colorMode, output.ColorModeAuto)
	}
}

func TestTreeColoredMetadata(t *testing.T) {
	t.Parallel()

	renderer := NewASCIITreeRenderer()
	renderer.SetColorMode(output.ColorModeAlways)

	node := output.NewTreeNode("node", "Node")
	node.Metadata["count"] = "42"
	node.Metadata["status"] = "active"
	renderer.SetRoot(node)

	got, err := renderer.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertContains(t, got, "count: 42", "should contain metadata")
	assertContains(t, got, "status: active", "should contain metadata")
}

func TestTreeMetadataDeterministicOrder(t *testing.T) {
	t.Parallel()

	build := func() string {
		root := &output.TreeNode{
			ID:    output.NewBrandedID[output.TreeNodeIDBrand]("root"),
			Label: output.NewBrandedID[output.TreeNodeLabelBrand]("Root"),
			Metadata: map[string]string{
				"zeta":  "last",
				"alpha": "first",
				"mid":   "middle",
			},
		}

		r := NewASCIITreeRenderer()
		r.SetRoot(root)

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		return out
	}

	first := build()

	// Re-render many times: map iteration order is random, so an unsorted
	// loop produces different (k: v) orders across calls.
	for i := 0; i < 20; i++ {
		if got := build(); got != first {
			t.Fatalf("metadata order nondeterministic between renders:\nfirst:\n%s\ngot:\n%s", first, got)
		}
	}

	want := "(alpha: first, mid: middle, zeta: last)"
	if !strings.Contains(first, want) {
		t.Errorf("metadata not in sorted key order, want %s in:\n%s", want, first)
	}
}
