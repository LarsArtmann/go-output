package plantuml

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
	"github.com/larsartmann/go-output/testhelpers/graphtest"
)

// errWriteFailed aliases testhelpers.ErrWrite for backward-compatible error checks.
var errWriteFailed = testhelpers.ErrWrite

func TestPlantUMLDiagramRender(t *testing.T) {
	t.Parallel()

	t.Run("empty diagram", func(t *testing.T) {
		t.Parallel()

		d := NewPlantUMLDiagram()

		out, err := d.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		testhelpers.AssertAllContained(t, out, "@startuml", "@enduml")
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

		testhelpers.AssertAllContained(t, out, "Service A", "Service B", "calls")
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
			testhelpers.AssertEqual(t, "sanitizePlantUMLID", tt.input, got, tt.want)
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

		testhelpers.AssertAllContained(t, out, "@startuml")
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

		testhelpers.AssertAllContained(t, out, "@startuml")
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

		testhelpers.AssertAllContained(t, buf.String(), "@startuml", "Alpha")
	})

	t.Run("propagates writer error", func(t *testing.T) {
		t.Parallel()

		data := output.NewTableData([]string{"Name"})
		data.AddRow([]string{"Alpha"})

		err := output.RenderTableData(
			data,
			output.FormatPlantUML,
			output.RenderOptions{Writer: &testhelpers.ErrorWriter{}},
		)
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

	testhelpers.AssertAllContained(t, out, "Parent", "Child", "Grandchild")
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

//nolint:exhaustruct // Test files use partial struct initialization
func TestPlantUMLDiagramWithNodeStyle(t *testing.T) {
	t.Parallel()

	d := NewPlantUMLDiagram()
	d.AddNode(output.GraphNode{
		ID:    output.NewBrandedID[output.GraphNodeIDBrand]("svc"),
		Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Service"),
		Style: output.NodeStyle{
			Fill:      "#e8a838",
			Stroke:    "#4a4030",
			FontColor: "#14110d",
		},
	})

	out, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertAllContained(
		t, out,
		"#e8a838;line:#4a4030;text:#14110d",
		"Service",
	)
}

func TestPlantUMLDiagramNoStyleNoColorSpec(t *testing.T) {
	t.Parallel()

	d := NewPlantUMLDiagram()
	d.AddNode(graphtest.NewTestNode("svc", "Plain"))

	out, err := d.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(out, "#[") {
		t.Errorf("Node without style should not emit color spec, got: %s", out)
	}

	testhelpers.AssertContains(t, out, "[Plain]", "Node label should still render")
}

// TestPlantUMLNodeStyleEscapesInjection verifies that malicious style values
// (semicolons, newlines, brackets) are escaped through the PlantUML render
// pipeline. A semicolon in a style value could inject additional PlantUML
// attributes; a newline could inject arbitrary PlantUML syntax.
func TestPlantUMLNodeStyleEscapesInjection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
	}{
		{"semicolon in Fill injects attribute", "red;line:evil"},
		{"newline in Fill", "red\n@startuml"},
		{"double quote in Fill", `red"; injected`},
		{"semicolon in Stroke", "#000;line:evil"},
		{"newline in FontColor", "#fff\n@enduml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := NewPlantUMLDiagram()
			d.AddNode(output.GraphNode{ //nolint:exhaustruct // Test uses minimal fields
				ID:    output.NewBrandedID[output.GraphNodeIDBrand]("n"),
				Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Test"),
				Style: output.NodeStyle{
					Fill:      tt.value,
					Stroke:    tt.value,
					FontColor: tt.value,
				},
			})

			out, err := d.Render()
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if strings.Contains(out, tt.value) {
				t.Errorf("raw malicious value %q leaked unescaped into PlantUML output", tt.value)
			}
		})
	}
}
