package output

import (
	"testing"
)

func TestTreeBuilder_FluentAPI(t *testing.T) {
	root := NewTreeBuilder().
		SetRoot("build", "Build").
		AddChild("build", "compile", "Compile").
		AddChild("build", "lint", "Lint").
		AddChild("compile", "test", "Test").
		Build()

	if root == nil {
		t.Fatal("Build() returned nil root")
	}

	if root.ID.Get() != "build" {
		t.Errorf("expected root ID 'build', got %q", root.ID.Get())
	}

	if root.Label.Get() != "Build" {
		t.Errorf("expected root label 'Build', got %q", root.Label.Get())
	}

	if len(root.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(root.Children))
	}

	compile := root.Children[0]
	if compile.ID.Get() != "compile" {
		t.Errorf("expected first child ID 'compile', got %q", compile.ID.Get())
	}

	if len(compile.Children) != 1 {
		t.Fatalf("expected compile to have 1 child, got %d", len(compile.Children))
	}

	if compile.Children[0].ID.Get() != "test" {
		t.Errorf("expected grandchild ID 'test', got %q", compile.Children[0].ID.Get())
	}
}

func TestTreeBuilder_NilRoot(t *testing.T) {
	root := NewTreeBuilder().Build()

	if root != nil {
		t.Errorf("expected nil root when SetRoot not called, got %v", root)
	}
}

func TestTreeBuilder_UnknownParent(t *testing.T) {
	root := NewTreeBuilder().
		SetRoot("root", "Root").
		AddChild("nonexistent", "child", "Child").
		Build()

	if len(root.Children) != 0 {
		t.Errorf("expected 0 children when parent not found, got %d", len(root.Children))
	}
}
