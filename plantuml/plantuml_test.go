package plantuml

import (
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

func TestPlantUMLFromTree(t *testing.T) {
	t.Parallel()

	t.Run("nil root", func(t *testing.T) {
		t.Parallel()

		d := PlantUMLFromTree(nil)
		if d == nil {
			t.Fatal("PlantUMLFromTree(nil) should return non-nil renderer")
		}
	})

	t.Run("with tree", func(t *testing.T) {
		t.Parallel()

		root := output.NewTreeNode("root", "Root")
		child := output.NewTreeNode("child", "Child")
		root.AddChild(child)

		d := PlantUMLFromTree(root)

		out, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertAllContained(t, out, "Root")
	})
}
