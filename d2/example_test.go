package d2_test

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/testhelpers"
)

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewDiagram() {
	diagram := d2.NewDiagram().
		SetDirection(d2.DirRight).
		SetTitle("System Architecture").
		AddNodeSimple("frontend", "Frontend App").
		AddNodeSimple("api", "API Server").
		AddNodeSimple("db", "Database")

	diagram.AddEdgeSimple("frontend", "api")
	diagram.AddEdgeSimple("api", "db")

	fmt.Println(testhelpers.MustRender(diagram))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewDiagram_tables() {
	diagram := d2.NewDiagram()

	diagram.
		AddTable("accounts", []d2.Column{
			{Name: "id", Type: "serial", Constraint: d2.ConstraintPrimary},
			{Name: "email", Type: "varchar", Constraint: d2.ConstraintUnique},
		}).
		AddTable("profiles", []d2.Column{
			{Name: "id", Type: "bigint", Constraint: d2.ConstraintPrimary},
			{Name: "account_id", Type: "bigint", Constraint: d2.ConstraintForeign},
			{Name: "bio", Type: "text"},
			{Name: "avatar_url", Type: "varchar"},
		})

	diagram.AddLabeledEdge("accounts", "profiles", "has many")

	fmt.Println(testhelpers.MustRender(diagram))
}

//nolint:testableexamples // Demonstration example, output is dynamic
func ExampleNewDiagram_styledNodes() {
	diagram := d2.NewDiagram()
	diagram.AddNode(d2.Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand]("server"),
		Label: output.NewBrandedID[output.D2NodeLabelBrand]("Web Server"),
		Shape: d2.ShapeHexagon,
		Style: d2.NodeStyle{
			Fill: "#E0F0FF",
			StrokeStyle: d2.StrokeStyle{
				Stroke: "#0066CC",
			},
			Shadow: true,
		},
	})

	fmt.Println(testhelpers.MustRender(diagram))
}
