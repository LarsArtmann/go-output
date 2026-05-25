package serialization

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestTOMLTreeRenderer(t *testing.T) {
	t.Parallel()

	t.Run("renders tree", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTreeRenderer()
		root := output.NewTreeNode("root", "Root")
		child := output.NewTreeNode("child1", "Child")
		child.Metadata = map[string]string{"key": "value"}
		root.AddChild(child)
		r.SetRoot(root)

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if !strings.Contains(out, "Root") {
			t.Error("TOML tree output should contain 'Root'")
		}

		if !strings.Contains(out, "Child") {
			t.Error("TOML tree output should contain 'Child'")
		}
	})

	t.Run("nil root returns empty", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTreeRenderer()

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if out != "" {
			t.Errorf("expected empty, got %q", out)
		}
	})

	t.Run("deep tree", func(t *testing.T) {
		t.Parallel()

		r := NewTOMLTreeRenderer()
		root := output.NewTreeNode("root", "Root")
		child := output.NewTreeNode("child", "Child")
		grandchild := output.NewTreeNode("gc", "Grandchild")
		child.AddChild(grandchild)
		root.AddChild(child)
		r.SetRoot(root)

		out, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if !strings.Contains(out, "Grandchild") {
			t.Error("TOML tree output should contain 'Grandchild'")
		}
	})
}
