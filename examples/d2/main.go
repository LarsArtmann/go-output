// Package main demonstrates D2 diagram features of the go-output library.
package main

import (
	"fmt"

	"github.com/larsartmann/go-output"
)

func main() {
	diagram := output.NewD2Diagram().
		SetDirection(output.D2DirRight).
		SetTitle("Microservice Architecture").
		AddClass("service", output.D2NodeStyle{
			Fill:     "lightblue",
			Stroke:   "navy",
			FontSize: 16,
		}).
		AddClass("database", output.D2NodeStyle{
			Fill:   "lightyellow",
			Stroke: "goldenrod",
		}).
		AddTable("users", []output.D2Column{
			{Name: "id", Type: "serial", Constraint: output.D2ConstraintPrimary},
			{Name: "email", Type: "varchar(255)", Constraint: output.D2ConstraintUnique},
			{Name: "created_at", Type: "timestamp"},
		}).
		AddTable("orders", []output.D2Column{
			{Name: "id", Type: "serial", Constraint: output.D2ConstraintPrimary},
			{Name: "user_id", Type: "int", Constraint: output.D2ConstraintForeign},
			{Name: "total", Type: "decimal(10,2)"},
		})

	diagram.
		AddNodeWithShape("api", "API Gateway", output.D2ShapeHexagon).
		AddNodeWithShape("auth", "Auth Service", output.D2ShapeCircle).
		AddNodeWithShape("orders-svc", "Order Service", output.D2ShapeRectangle).
		AddNodeWithShape("cache", "Redis Cache", output.D2ShapeCylinder).
		AddNode(output.D2Node{
			ID:    output.NewBrandedID[output.D2NodeIDBrand]("db"),
			Label: output.NewBrandedID[output.D2NodeLabelBrand]("PostgreSQL"),
			Shape: output.D2ShapeCylinder,
			Class: "database",
		})

	diagram.
		AddLabeledEdge("api", "auth", "authenticate").
		AddLabeledEdge("api", "orders-svc", "route").
		AddEdge(output.D2Edge{
			From:        output.NewBrandedID[output.D2NodeIDBrand]("orders-svc"),
			To:          output.NewBrandedID[output.D2NodeIDBrand]("cache"),
			TargetArrow: output.D2ArrowArrow,
		}).
		AddEdgeSimple("orders-svc", "db").
		AddEdge(output.D2Edge{
			From:        output.NewBrandedID[output.D2NodeIDBrand]("users"),
			To:          output.NewBrandedID[output.D2NodeIDBrand]("db"),
			TargetArrow: output.D2ArrowCFMany,
			Label:       output.NewBrandedID[output.D2NodeLabelBrand]("stores"),
		})

	fmt.Println(diagram.Render())
}
