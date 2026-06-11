package plantuml

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

func assertAllContained(t *testing.T, haystack string, needles ...string) {
	t.Helper()

	for _, n := range needles {
		if !strings.Contains(haystack, n) {
			t.Errorf("should contain %q", n)
		}
	}
}

// errWriter always returns an error from Write.
type errWriter struct{}

func (errWriter) Write(p []byte) (int, error) {
	return 0, errWriteFailed
}

var errWriteFailed = errors.New("write failed")

func TestPlantUMLDiagramRender(t *testing.T) {
	t.Parallel()

	t.Run("empty diagram", func(t *testing.T) {
		t.Parallel()

		d := NewPlantUMLDiagram()

		out, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertAllContained(t, out, "@startuml", "@enduml")
	})

	t.Run("with nodes and edges", func(t *testing.T) {
		t.Parallel()

		d := NewPlantUMLDiagram()
		d.AddNode(graphtest.NewTestNode("svc-a", "Service A"))
		d.AddNode(graphtest.NewTestNode("svc-b", "Service B"))
		d.AddEdge(graphtest.NewTestEdge("svc-a", "svc-b", "calls"))

		out, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertAllContained(t, out, "Service A", "Service B", "calls")
	})
}

func TestSanitizePlantUMLID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"hello world", "hello_world"},
		{"my-service", "my_service"},
		{"my-cool-service", "my_cool_service"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got := sanitizePlantUMLID(tt.input)
			if got != tt.want {
				t.Errorf("sanitizePlantUMLID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPlantUMLFromTableData(t *testing.T) {
	t.Parallel()

	t.Run("nil data", func(t *testing.T) {
		t.Parallel()

		d := PlantUMLFromTableData(nil)
		if d == nil {
			t.Fatal("PlantUMLFromTableData(nil) should return non-nil renderer")
		}

		out, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertAllContained(t, out, "@startuml")
	})

	t.Run("with data", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})
		data.AddRow([]string{"Beta"})

		d := PlantUMLFromTableData(data)

		out, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertAllContained(t, out, "@startuml")
	})
}

func TestRenderPlantUMLTableData(t *testing.T) {
	t.Parallel()

	t.Run("writes PlantUML output to writer", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		var buf bytes.Buffer

		err := output.RenderTableData(data, output.FormatPlantUML, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTableData plantuml: %v", err)
		}

		assertAllContained(t, buf.String(), "@startuml", "Alpha")
	})

	t.Run("propagates writer error", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		err := output.RenderTableData(data, output.FormatPlantUML, output.RenderOptions{Writer: &errWriter{}})
		if err == nil {
			t.Fatal("expected error from errWriter")
		}
	})
}

func TestPlantUMLFromTree_Branches(t *testing.T) {
	t.Parallel()

	// Build a multi-level tree: parent → child → grandchild
	parent := output.NewTreeNode("parent", "Parent")
	child := output.NewTreeNode("child", "Child")
	grandchild := output.NewTreeNode("grandchild", "Grandchild")
	child.AddChild(grandchild)
	parent.AddChild(child)

	d := PlantUMLFromTree(parent)

	out, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	assertAllContained(t, out, "Parent", "Child", "Grandchild")
}

func TestPlantUMLDiagramAddNodeExistingID(t *testing.T) {
	t.Parallel()

	// Adding two nodes with the same ID: the second call should replace
	// the first (matches AddNode semantics from the GraphRendererState).
	d := NewPlantUMLDiagram()
	d.AddNode(graphtest.NewTestNode("svc", "Original"))
	d.AddNode(graphtest.NewTestNode("svc", "Updated"))

	out, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(out, "Updated") {
		t.Errorf("expected 'Updated' to replace 'Original', got %q", out)
	}
}
