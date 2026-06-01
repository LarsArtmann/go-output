package d2_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewD2Diagram() {
	diagram := d2.NewD2Diagram().
		SetDirection(d2.D2DirRight).
		SetTitle("System Architecture").
		AddNodeSimple("frontend", "Frontend App").
		AddNodeSimple("api", "API Server").
		AddNodeSimple("db", "Database")

	diagram.AddEdgeSimple("frontend", "api")
	diagram.AddEdgeSimple("api", "db")

	result, err := diagram.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewD2Diagram_tables() {
	diagram := d2.NewD2Diagram()

	diagram.
		AddTable("accounts", []d2.D2Column{
			{Name: "id", Type: "serial", Constraint: d2.D2ConstraintPrimary},
			{Name: "email", Type: "varchar", Constraint: d2.D2ConstraintUnique},
		}).
		AddTable("profiles", []d2.D2Column{
			{Name: "id", Type: "bigint", Constraint: d2.D2ConstraintPrimary},
			{Name: "account_id", Type: "bigint", Constraint: d2.D2ConstraintForeign},
			{Name: "bio", Type: "text"},
			{Name: "avatar_url", Type: "varchar"},
		})

	diagram.AddLabeledEdge("accounts", "profiles", "has many")

	result, err := diagram.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewD2Diagram_styledNodes() {
	diagram := d2.NewD2Diagram()
	diagram.AddNode(d2.D2Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand]("server"),
		Label: output.NewBrandedID[output.D2NodeLabelBrand]("Web Server"),
		Shape: d2.D2ShapeHexagon,
		Style: d2.D2NodeStyle{
			Fill: "#E0F0FF",
			D2StrokeStyle: d2.D2StrokeStyle{
				Stroke: "#0066CC",
			},
			Shadow: true,
		},
	})

	result, err := diagram.Render()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		return
	}

	fmt.Println(result)
}
