package plantuml

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestPlantUMLDiagramRender(t *testing.T) {
	t.Parallel()

	t.Run("empty diagram", func(t *testing.T) {
		t.Parallel()

		d := NewPlantUMLDiagram()

		out, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if !strings.Contains(out, "@startuml") {
			t.Error("PlantUML output should contain @startuml")
		}

		if !strings.Contains(out, "@enduml") {
			t.Error("PlantUML output should contain @enduml")
		}
	})

	t.Run("with nodes and edges", func(t *testing.T) {
		t.Parallel()

		d := NewPlantUMLDiagram()
		d.AddNode(output.GraphNode{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("svc-a"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Service A"),
		})
		d.AddNode(output.GraphNode{
			ID:    output.NewBrandedID[output.GraphNodeIDBrand]("svc-b"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Service B"),
		})
		d.AddEdge(output.GraphEdge{
			From:  output.NewBrandedID[output.GraphNodeIDBrand]("svc-a"),
			To:    output.NewBrandedID[output.GraphNodeIDBrand]("svc-b"),
			Label: output.NewBrandedID[output.GraphNodeLabelBrand]("calls"),
		})

		out, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if !strings.Contains(out, "Service A") {
			t.Error("output should contain 'Service A'")
		}

		if !strings.Contains(out, "Service B") {
			t.Error("output should contain 'Service B'")
		}

		if !strings.Contains(out, "calls") {
			t.Error("output should contain edge label 'calls'")
		}
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

		if !strings.Contains(out, "@startuml") {
			t.Error("should contain @startuml")
		}
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

		if !strings.Contains(out, "@startuml") {
			t.Error("should contain @startuml")
		}
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

		if !strings.Contains(out, "Root") {
			t.Error("should contain 'Root'")
		}
	})
}
