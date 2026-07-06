package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestCQRS_WriteASCII_NoErrorOnSuccess(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Root")
	child := output.NewTreeNode("child", "Child")
	root.AddChild(child)

	var buf strings.Builder

	err := WriteASCII(&buf, root)
	if err != nil {
		t.Fatalf("WriteASCII should return nil on success, got: %v", err)
	}

	if !strings.Contains(buf.String(), "Root") {
		t.Errorf("expected output to contain 'Root', got %q", buf.String())
	}
}

func TestCQRS_RenderMarkdown(t *testing.T) {
	t.Parallel()

	root := output.NewTreeNode("root", "Root")
	child1 := output.NewTreeNode("c1", "Child 1")
	child2 := output.NewTreeNode("c2", "Child 2")
	grandchild := output.NewTreeNode("gc", "Grandchild")
	child1.AddChild(grandchild)
	root.AddChild(child1)
	root.AddChild(child2)

	got, err := RenderMarkdown(root)
	if err != nil {
		t.Fatalf("RenderMarkdown failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (root + 2 children + 1 grandchild), got %d:\n%s", len(lines), got)
	}

	if !strings.Contains(lines[0], "- Root") {
		t.Errorf("first line should be '- Root', got %q", lines[0])
	}

	if !strings.Contains(lines[1], "  - Child 1") {
		t.Errorf("second line should be indented '  - Child 1', got %q", lines[1])
	}

	if !strings.Contains(lines[2], "    - Grandchild") {
		t.Errorf("grandchild should be double-indented '    - Grandchild', got %q", lines[2])
	}
}

func TestCQRS_WriteMarkdown_NilRoot(t *testing.T) {
	t.Parallel()

	var buf strings.Builder

	err := WriteMarkdown(&buf, nil)
	if err != nil {
		t.Fatalf("WriteMarkdown(nil) should not error, got: %v", err)
	}

	if buf.Len() != 0 {
		t.Errorf("expected empty output for nil root, got %q", buf.String())
	}
}

func ExampleRenderMarkdown() {
	root := output.NewTreeNode("build", "Build")
	compile := output.NewTreeNode("compile", "Compile")
	test := output.NewTreeNode("test", "Test")

	root.AddChild(compile)
	root.AddChild(test)

	out, _ := RenderMarkdown(root)
	fmt.Println(out)
	// Output:
	// - Build
	//   - Compile
	//   - Test
}
