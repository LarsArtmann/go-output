package d2_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/testhelpers"
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

	fmt.Println(testhelpers.MustRender(diagram))
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

	fmt.Println(testhelpers.MustRender(diagram))
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

	fmt.Println(testhelpers.MustRender(diagram))
}
